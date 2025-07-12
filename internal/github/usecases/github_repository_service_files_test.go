package usecases

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestGetFileContents はGetFileContentsメソッドをテストする
func TestGetFileContents(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		path           string
		branch         string
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:   "正常系 - ファイル取得成功",
			owner:  "test_user",
			repo:   "test_repo",
			path:   "README.md",
			branch: "",
			mockResponse: map[string]interface{}{
				"name":     "README.md",
				"path":     "README.md",
				"type":     "file",
				"content":  base64.StdEncoding.EncodeToString([]byte("# テストリポジトリ\nこれはテスト用のREADMEファイルです。")),
				"encoding": "base64",
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:   "正常系 - 指定ブランチのファイル取得",
			owner:  "test_user",
			repo:   "test_repo",
			path:   "README.md",
			branch: "develop",
			mockResponse: map[string]interface{}{
				"name":     "README.md",
				"path":     "README.md",
				"type":     "file",
				"content":  base64.StdEncoding.EncodeToString([]byte("# 開発ブランチ\nこれは開発ブランチのREADMEファイルです。")),
				"encoding": "base64",
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - ファイルが存在しない",
			owner:          "test_user",
			repo:           "test_repo",
			path:           "nonexistent.md",
			branch:         "",
			mockResponse:   nil,
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// モックHTTPクライアントの作成
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if tc.mockError != nil {
						return nil, tc.mockError
					}

					// リクエストの検証
					expectedURL := fmt.Sprintf("%s/repos/%s/%s/contents/%s", apiBaseURL, tc.owner, tc.repo, tc.path)
					if tc.branch != "" {
						expectedURL += fmt.Sprintf("?ref=%s", tc.branch)
					}
					if req.URL.String() != expectedURL {
						t.Errorf("期待されたURL: %s, 実際: %s", expectedURL, req.URL.String())
					}

					// モックレスポンスの作成
					var responseBody []byte
					if tc.mockResponse != nil {
						responseBody, _ = json.Marshal(tc.mockResponse)
					}

					return &http.Response{
						StatusCode: tc.mockStatusCode,
						Body:       io.NopCloser(bytes.NewReader(responseBody)),
					}, nil
				},
			}

			// GitHubRepositoryServiceの作成
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := &GitHubRepositoryService{clientService: clientService}

			// テスト対象の関数を実行
			result, err := service.GetFileContents(tc.owner, tc.repo, tc.path, tc.branch)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError && tc.mockResponse != nil {
				// decoded_contentフィールドの検証
				if content, ok := tc.mockResponse["content"].(string); ok {
					expectedContent, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(content, "\n", ""))
					if string(expectedContent) != result["decoded_content"] {
						t.Errorf("decoded_contentの値が異なります。期待: %s, 実際: %s", string(expectedContent), result["decoded_content"])
					}
				}
			}
		})
	}
}

// TestHandleToGetFileContents はHandleToGetFileContentsメソッドをテストする
func TestHandleToGetFileContents(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		path           string
		branch         string
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:   "正常系 - ハンドラー成功",
			owner:  "test_user",
			repo:   "test_repo",
			path:   "README.md",
			branch: "",
			mockResponse: map[string]interface{}{
				"name":     "README.md",
				"content":  base64.StdEncoding.EncodeToString([]byte("テストファイル")),
				"encoding": "base64",
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - APIエラー",
			owner:          "test_user",
			repo:           "test_repo",
			path:           "README.md",
			branch:         "",
			mockResponse:   nil,
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// モックHTTPクライアントの作成
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					if tc.mockError != nil {
						return nil, tc.mockError
					}

					// モックレスポンスの作成
					var responseBody []byte
					if tc.mockResponse != nil {
						responseBody, _ = json.Marshal(tc.mockResponse)
					}

					return &http.Response{
						StatusCode: tc.mockStatusCode,
						Body:       io.NopCloser(bytes.NewReader(responseBody)),
					}, nil
				},
			}

			// GitHubRepositoryServiceの作成
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := &GitHubRepositoryService{clientService: clientService}

			// テスト対象の関数を実行
			result, err := service.HandleToGetFileContents(tc.owner, tc.repo, tc.path, tc.branch)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError {
				if result == "" {
					t.Error("結果が空文字列です")
				}
			}
		})
	}
}
