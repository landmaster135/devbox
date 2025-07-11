package usecases

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestNewGitHubRepositoryService はNewGitHubRepositoryService関数をテストする
func TestNewGitHubRepositoryService(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "正常系 - トークンあり",
			token: "test_token",
		},
		{
			name:  "正常系 - トークンなし",
			token: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := NewGitHubRepositoryService(tc.token)

			if service == nil {
				t.Fatal("NewGitHubRepositoryServiceがnilを返しました")
			}

			if service.clientService == nil {
				t.Fatal("clientServiceがnilです")
			}
		})
	}
}

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
			result, err := service.ListCommits(tc.owner, tc.repo, tc.page, tc.perPage, tc.sha)

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

// TestSearchRepositories はSearchRepositoriesメソッドをテストする
func TestSearchRepositories(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		page            int
		perPage         int
		mockResponse    map[string]interface{}
		mockStatusCode  int
		mockError       error
		expectError     bool
		expectedPage    int
		expectedPerPage int
	}{
		{
			name:            "正常系 - リポジトリ検索成功",
			query:           "test",
			page:            1,
			perPage:         30,
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse: map[string]interface{}{
				"total_count":        float64(2),
				"incomplete_results": false,
				"items": []interface{}{
					map[string]interface{}{
						"id":        float64(123456),
						"name":      "test-repo-1",
						"full_name": "test-user/test-repo-1",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:            "正常系 - pageが1未満の場合は1に設定される",
			query:           "test",
			page:            0,
			perPage:         30,
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse: map[string]interface{}{
				"total_count": float64(1),
				"items":       []interface{}{},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:            "異常系 - 認証エラー",
			query:           "test",
			page:            1,
			perPage:         30,
			expectedPage:    1,
			expectedPerPage: 30,
			mockResponse:    nil,
			mockStatusCode:  http.StatusUnauthorized,
			mockError:       nil,
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

					expectedURL := fmt.Sprintf("%s/search/repositories?q=%s&page=%d&per_page=%d", apiBaseURL, tc.query, expectedPage, expectedPerPage)
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
			result, err := service.SearchRepositories(tc.query, tc.page, tc.perPage)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError && tc.mockResponse != nil {
				if result["total_count"] != tc.mockResponse["total_count"] {
					t.Errorf("total_countが異なります。期待: %v, 実際: %v", tc.mockResponse["total_count"], result["total_count"])
				}
			}
		})
	}
}

// TestGetUserRepositories はGetUserRepositoriesメソッドをテストする
func TestGetUserRepositories(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		options        map[string]interface{}
		mockResponse   []map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:     "正常系 - オプションなし",
			username: "test_user",
			options:  map[string]interface{}{},
			mockResponse: []map[string]interface{}{
				{
					"id":        float64(123456),
					"name":      "test-repo-1",
					"full_name": "test_user/test-repo-1",
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:     "正常系 - すべてのオプション",
			username: "test_user",
			options: map[string]interface{}{
				"sort":      "updated",
				"direction": "asc",
				"per_page":  10,
				"page":      2,
				"type":      "owner",
			},
			mockResponse: []map[string]interface{}{
				{
					"id":        float64(123458),
					"name":      "test-repo-3",
					"full_name": "test_user/test-repo-3",
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - 認証エラー",
			username:       "test_user",
			options:        map[string]interface{}{},
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

					// リクエストの検証
					expectedBaseURL := fmt.Sprintf("%s/users/%s/repos", apiBaseURL, tc.username)
					if !strings.HasPrefix(req.URL.String(), expectedBaseURL) {
						t.Errorf("期待されたURLの基本部分: %s, 実際: %s", expectedBaseURL, req.URL.String())
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
			result, err := service.GetUserRepositories(tc.username, tc.options)

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

// TestHandleToSearchRepositories はHandleToSearchRepositoriesメソッドをテストする
func TestHandleToSearchRepositories(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		page           int
		perPage        int
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:    "正常系 - ハンドラー成功",
			query:   "test",
			page:    1,
			perPage: 30,
			mockResponse: map[string]interface{}{
				"total_count": float64(1),
				"items":       []interface{}{},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - APIエラー",
			query:          "test",
			page:           1,
			perPage:        30,
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
			result, err := service.HandleToSearchRepositories(tc.query, tc.page, tc.perPage)

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

// TestHandleToGetUserRepositories はHandleToGetUserRepositoriesメソッドをテストする
func TestHandleToGetUserRepositories(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		sort           string
		direction      string
		type_          string
		perPage        int
		page           int
		mockResponse   []map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:      "正常系 - ハンドラー成功",
			username:  "test_user",
			sort:      "",
			direction: "",
			type_:     "",
			perPage:   0,
			page:      0,
			mockResponse: []map[string]interface{}{
				{
					"id":   float64(123456),
					"name": "test-repo",
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - APIエラー",
			username:       "test_user",
			sort:           "",
			direction:      "",
			type_:          "",
			perPage:        0,
			page:           0,
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
			result, err := service.HandleToGetUserRepositories(tc.username, tc.sort, tc.direction, tc.type_, tc.perPage, tc.page)

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
