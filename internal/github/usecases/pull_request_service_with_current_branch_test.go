package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

// createMockResponse はテスト用のHTTPレスポンスを作成します
func createMockResponse(statusCode int, data map[string]interface{}) (*http.Response, error) {
	responseBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil
}

// TestHandleToCreatePullRequestWithCurrentBranch は新しいメソッドのテストです
func TestHandleToCreatePullRequestWithCurrentBranch(t *testing.T) {
	tests := []struct {
		name              string
		owner             string
		repo              string
		title             string
		base              string
		body              string
		draft             bool
		repoPath          string
		mockCurrentBranch string
		mockError         error
		expectError       bool
		expectedBranch    string
	}{
		{
			name:              "正常系 - 現在のブランチを使用してプルリクエスト作成",
			owner:             "test_user",
			repo:              "test_repo",
			title:             "Test Pull Request",
			base:              "main",
			body:              "This is a test pull request",
			draft:             false,
			repoPath:          "/home/user/test_repo",
			mockCurrentBranch: "feature/test-branch",
			mockError:         nil,
			expectError:       false,
			expectedBranch:    "feature/test-branch",
		},
		{
			name:              "正常系 - ドラフトプルリクエスト作成",
			owner:             "test_user",
			repo:              "test_repo",
			title:             "Draft Pull Request",
			base:              "develop",
			body:              "",
			draft:             true,
			repoPath:          "/home/user/test_repo",
			mockCurrentBranch: "feature/draft-branch",
			mockError:         nil,
			expectError:       false,
			expectedBranch:    "feature/draft-branch",
		},
		{
			name:              "異常系 - Gitブランチ取得エラー",
			owner:             "test_user",
			repo:              "test_repo",
			title:             "Test Pull Request",
			base:              "main",
			body:              "This is a test pull request",
			draft:             false,
			repoPath:          "/invalid/path",
			mockCurrentBranch: "",
			mockError:         errors.New("repository path does not exist: /invalid/path"),
			expectError:       true,
			expectedBranch:    "",
		},
		{
			name:              "異常系 - 絶対パスでない場合",
			owner:             "test_user",
			repo:              "test_repo",
			title:             "Test Pull Request",
			base:              "main",
			body:              "This is a test pull request",
			draft:             false,
			repoPath:          "relative/path",
			mockCurrentBranch: "",
			mockError:         errors.New("repo_path must be an absolute path, got: relative/path"),
			expectError:       true,
			expectedBranch:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// モックHTTPクライアントの作成
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// プルリクエスト作成が成功した場合のモックレスポンス
					if tc.mockError == nil {
						return createMockResponse(200, map[string]interface{}{
							"id":     123,
							"number": 1,
							"title":  tc.title,
							"head": map[string]interface{}{
								"ref": tc.expectedBranch,
							},
							"base": map[string]interface{}{
								"ref": tc.base,
							},
							"draft": tc.draft,
						})
					}
					return nil, errors.New("mock error")
				},
			}

			// モックGitBranchProviderの作成
			mockGitBranchProvider := &MockGitBranchProvider{
				CurrentBranch: tc.mockCurrentBranch,
				Error:         tc.mockError,
			}

			// GitHubPullRequestServiceの作成
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := NewGitHubPullRequestServiceWithDependencies(clientService, mockGitBranchProvider)

			// テスト対象の関数を実行
			result, err := service.HandleToCreatePullRequestWithCurrentBranch(
				tc.owner, tc.repo, tc.title, tc.base, tc.body, tc.draft, tc.repoPath,
			)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError && result != "" {
				// JSONレスポンスが返されることを確認
				if len(result) == 0 {
					t.Error("結果が空です")
				}
			}
		})
	}
}

// TestDefaultGitBranchProvider_GetCurrentBranchFromPath はDefaultGitBranchProviderのテストです
func TestDefaultGitBranchProvider_GetCurrentBranchFromPath(t *testing.T) {
	tests := []struct {
		name          string
		absolutePath  string
		expectError   bool
		errorContains string
	}{
		{
			name:          "異常系 - 相対パス",
			absolutePath:  "relative/path",
			expectError:   true,
			errorContains: "repo_path must be an absolute path",
		},
		{
			name:          "異常系 - 存在しないパス",
			absolutePath:  "/nonexistent/path",
			expectError:   true,
			errorContains: "repository path does not exist",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			provider := &DefaultGitBranchProvider{}

			_, err := provider.GetCurrentBranchFromPath(tc.absolutePath)

			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}
			if tc.expectError && err != nil && tc.errorContains != "" {
				if !contains(err.Error(), tc.errorContains) {
					t.Errorf("エラーメッセージに '%s' が含まれていません。実際のエラー: %v", tc.errorContains, err)
				}
			}
		})
	}
}

// contains は文字列に部分文字列が含まれているかチェックします
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

// containsSubstring は文字列の中に部分文字列があるかチェックします
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
