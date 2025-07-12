package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

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

// TestHandleToGetUserRepositories_Additional はHandleToGetUserRepositories関数の追加テストケース
func TestHandleToGetUserRepositories_Additional(t *testing.T) {
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
			name:      "正常系 - 全パラメータ指定",
			username:  "test_user",
			sort:      "updated",
			direction: "desc",
			type_:     "owner",
			perPage:   50,
			page:      2,
			mockResponse: []map[string]interface{}{
				{
					"id":        float64(123456),
					"name":      "test-repo-1",
					"full_name": "test_user/test-repo-1",
					"private":   false,
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:      "正常系 - perPageが0の場合",
			username:  "test_user",
			sort:      "created",
			direction: "asc",
			type_:     "public",
			perPage:   0, // 0の場合はoptionsに追加されない
			page:      1,
			mockResponse: []map[string]interface{}{
				{
					"id":        float64(123457),
					"name":      "test-repo-2",
					"full_name": "test_user/test-repo-2",
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:      "正常系 - pageが0の場合",
			username:  "test_user",
			sort:      "pushed",
			direction: "",
			type_:     "",
			perPage:   30,
			page:      0, // 0の場合はoptionsに追加されない
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
			name:      "正常系 - 負の値のperPageとpage",
			username:  "test_user",
			sort:      "",
			direction: "",
			type_:     "",
			perPage:   -1, // 負の値の場合はoptionsに追加されない
			page:      -1, // 負の値の場合はoptionsに追加されない
			mockResponse: []map[string]interface{}{
				{
					"id":        float64(123459),
					"name":      "test-repo-4",
					"full_name": "test_user/test-repo-4",
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - ユーザーが存在しない",
			username:       "nonexistent_user",
			sort:           "",
			direction:      "",
			type_:          "",
			perPage:        30,
			page:           1,
			mockResponse:   nil,
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - 認証エラー",
			username:       "private_user",
			sort:           "",
			direction:      "",
			type_:          "",
			perPage:        30,
			page:           1,
			mockResponse:   nil,
			mockStatusCode: http.StatusUnauthorized,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - ネットワークエラー",
			username:       "test_user",
			sort:           "",
			direction:      "",
			type_:          "",
			perPage:        30,
			page:           1,
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      errors.New("ネットワーク接続エラー"),
			expectError:    true,
		},
		{
			name:           "異常系 - 不正なJSONレスポンス",
			username:       "test_user",
			sort:           "",
			direction:      "",
			type_:          "",
			perPage:        30,
			page:           1,
			mockResponse:   nil,
			mockStatusCode: http.StatusOK,
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
					if tc.name == "異常系 - 不正なJSONレスポンス" {
						responseBody = []byte("{invalid json}")
					} else if tc.mockResponse != nil {
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

				// JSONとしてパース可能かチェック
				var parsedResult []map[string]interface{}
				if err := json.Unmarshal([]byte(result), &parsedResult); err != nil {
					t.Errorf("結果のJSONパースに失敗しました: %v", err)
				}
			}
		})
	}
}
