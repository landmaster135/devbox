package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// TestListCommits はListCommitsメソッドをテストする
func TestListCommits(t *testing.T) {
	tests := []struct {
		name            string
		owner           string
		repo            string
		page            int
		perPage         int
		sha             string
		mockResponse    []map[string]interface{}
		mockStatusCode  int
		mockError       error
		expectError     bool
		expectedPage    int
		expectedPerPage int
	}{
		{
			name:            "正常系 - コミット一覧取得成功",
			owner:           "test_user",
			repo:            "test_repo",
			page:            1,
			perPage:         30,
			sha:             "",
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse: []map[string]interface{}{
				{
					"sha": "abc123def456",
					"commit": map[string]interface{}{
						"message": "最初のコミット",
						"author": map[string]interface{}{
							"name":  "Test User",
							"email": "test@example.com",
							"date":  "2023-01-01T12:00:00Z",
						},
					},
					"author": map[string]interface{}{
						"login": "test_user",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:            "正常系 - 特定のブランチのコミット一覧",
			owner:           "test_user",
			repo:            "test_repo",
			page:            1,
			perPage:         30,
			sha:             "develop",
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse: []map[string]interface{}{
				{
					"sha": "abc123def456",
					"commit": map[string]interface{}{
						"message": "developブランチのコミット",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:            "正常系 - pageが1未満の場合は1に設定される",
			owner:           "test_user",
			repo:            "test_repo",
			page:            0,
			perPage:         30,
			sha:             "",
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse: []map[string]interface{}{
				{
					"sha": "abc123def456",
					"commit": map[string]interface{}{
						"message": "テストコミット",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:            "正常系 - perPageが1未満の場合は30に設定される",
			owner:           "test_user",
			repo:            "test_repo",
			page:            1,
			perPage:         0,
			sha:             "",
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse: []map[string]interface{}{
				{
					"sha": "abc123def456",
					"commit": map[string]interface{}{
						"message": "テストコミット",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:            "異常系 - 認証エラー",
			owner:           "test_user",
			repo:            "test_repo",
			page:            1,
			perPage:         30,
			sha:             "",
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse:    nil,
			mockStatusCode:  http.StatusUnauthorized,
			mockError:       nil,
			expectError:     true,
		},
		{
			name:            "異常系 - ネットワークエラー",
			owner:           "test_user",
			repo:            "test_repo",
			page:            1,
			perPage:         30,
			sha:             "",
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse:    nil,
			mockStatusCode:  0,
			mockError:       errors.New("ネットワーク接続エラー"),
			expectError:     true,
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
					expectedPage := tc.page
					if expectedPage < 1 {
						expectedPage = 1
					}
					expectedPerPage := tc.perPage
					if expectedPerPage < 1 || expectedPerPage > 100 {
						expectedPerPage = 30
					}

					expectedURL := fmt.Sprintf("%s/repos/%s/%s/commits?page=%d&per_page=%d", apiBaseURL, tc.owner, tc.repo, expectedPage, expectedPerPage)
					if tc.sha != "" {
						expectedURL += fmt.Sprintf("&sha=%s", tc.sha)
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
			result, err := service.listCommits(tc.owner, tc.repo, tc.page, tc.perPage, tc.sha)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError && tc.mockResponse != nil {
				if len(result) != len(tc.mockResponse) {
					t.Errorf("期待された結果の長さ: %d, 実際: %d", len(tc.mockResponse), len(result))
				}
			}
		})
	}
}

// TestHandleToListCommits はHandleToListCommitsメソッドをテストする
func TestHandleToListCommits(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		page           int
		perPage        int
		sha            string
		mockResponse   []map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:    "正常系 - ハンドラー成功",
			owner:   "test_user",
			repo:    "test_repo",
			page:    1,
			perPage: 30,
			sha:     "",
			mockResponse: []map[string]interface{}{
				{
					"sha": "abc123def456",
					"commit": map[string]interface{}{
						"message": "テストコミット",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - APIエラー",
			owner:          "test_user",
			repo:           "test_repo",
			page:           1,
			perPage:        30,
			sha:            "",
			mockResponse:   nil,
			mockStatusCode: http.StatusUnauthorized,
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
			result, err := service.HandleToListCommits(tc.owner, tc.repo, tc.page, tc.perPage, tc.sha)

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
