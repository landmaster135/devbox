package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

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

// TestCreatePullRequestReview_ErrorCases はCreatePullRequestReview関数のエラーケースをテストする
func TestCreatePullRequestReview_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		pullNumber     int
		options        map[string]interface{}
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:       "正常系 - レビュー作成成功",
			owner:      "test_user",
			repo:       "test_repo",
			pullNumber: 1,
			options: map[string]interface{}{
				"body":  "これは良いプルリクエストです",
				"event": "APPROVE",
			},
			mockResponse: map[string]interface{}{
				"id":    float64(123456),
				"body":  "これは良いプルリクエストです",
				"state": "APPROVED",
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:       "異常系 - JSONマーシャリングエラー",
			owner:      "test_user",
			repo:       "test_repo",
			pullNumber: 1,
			options: map[string]interface{}{
				"body":     "テストレビュー",
				"callback": func() {}, // JSONにシリアライズできない値
			},
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - プルリクエストが存在しない",
			owner:          "test_user",
			repo:           "test_repo",
			pullNumber:     999,
			options:        map[string]interface{}{"body": "テストレビュー"},
			mockResponse:   nil,
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - 権限不足",
			owner:          "test_user",
			repo:           "test_repo",
			pullNumber:     1,
			options:        map[string]interface{}{"body": "テストレビュー"},
			mockResponse:   nil,
			mockStatusCode: http.StatusForbidden,
			mockError:      nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// JSONマーシャリングエラーのテストケース
			if tc.name == "異常系 - JSONマーシャリングエラー" {
				service := NewGitHubPullRequestService("test_token")
				_, err := service.CreatePullRequestReview(tc.owner, tc.repo, tc.pullNumber, tc.options)
				if !tc.expectError && err != nil {
					t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
				}
				if tc.expectError && err == nil {
					t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
				}
				return
			}

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

			// GitHubPullRequestServiceの作成
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := NewGitHubPullRequestServiceWithDependencies(clientService)

			// テスト対象の関数を実行
			result, err := service.CreatePullRequestReview(tc.owner, tc.repo, tc.pullNumber, tc.options)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError && tc.mockResponse != nil {
				compareMaps(t, tc.mockResponse, result)
			}
		})
	}
}

// TestMergePullRequest_ErrorCases はMergePullRequest関数のエラーケースをテストする
func TestMergePullRequest_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		pullNumber     int
		options        map[string]interface{}
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:       "正常系 - マージ成功",
			owner:      "test_user",
			repo:       "test_repo",
			pullNumber: 1,
			options: map[string]interface{}{
				"commit_title":   "Merge pull request #1",
				"commit_message": "詳細なマージメッセージ",
				"merge_method":   "merge",
			},
			mockResponse: map[string]interface{}{
				"sha":     "abc123def456",
				"merged":  true,
				"message": "Pull Request successfully merged",
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:       "異常系 - JSONマーシャリングエラー",
			owner:      "test_user",
			repo:       "test_repo",
			pullNumber: 1,
			options: map[string]interface{}{
				"commit_title": "Merge pull request #1",
				"callback":     func() {}, // JSONにシリアライズできない値
			},
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - マージできない状態",
			owner:          "test_user",
			repo:           "test_repo",
			pullNumber:     1,
			options:        map[string]interface{}{"merge_method": "merge"},
			mockResponse:   nil,
			mockStatusCode: http.StatusMethodNotAllowed,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - コンフリクトが存在",
			owner:          "test_user",
			repo:           "test_repo",
			pullNumber:     1,
			options:        map[string]interface{}{"merge_method": "merge"},
			mockResponse:   nil,
			mockStatusCode: http.StatusConflict,
			mockError:      nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// JSONマーシャリングエラーのテストケース
			if tc.name == "異常系 - JSONマーシャリングエラー" {
				service := NewGitHubPullRequestService("test_token")
				_, err := service.MergePullRequest(tc.owner, tc.repo, tc.pullNumber, tc.options)
				if !tc.expectError && err != nil {
					t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
				}
				if tc.expectError && err == nil {
					t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
				}
				return
			}

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

			// GitHubPullRequestServiceの作成
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := NewGitHubPullRequestServiceWithDependencies(clientService)

			// テスト対象の関数を実行
			result, err := service.MergePullRequest(tc.owner, tc.repo, tc.pullNumber, tc.options)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError && tc.mockResponse != nil {
				compareMaps(t, tc.mockResponse, result)
			}
		})
	}
}

// TestUpdateIssue_ErrorCases はUpdateIssue関数の追加エラーケースをテストする
func TestUpdateIssue_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		issueNumber    int
		options        map[string]interface{}
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:        "異常系 - 不正なJSONレスポンス",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			options: map[string]interface{}{
				"title": "更新されたイシュー",
			},
			mockResponse:   nil,
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:        "異常系 - ネットワークエラー",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			options: map[string]interface{}{
				"title": "更新されたイシュー",
			},
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      errors.New("ネットワーク接続エラー"),
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

			// GitHubIssueServiceの作成
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := NewGitHubIssueServiceWithDependencies(clientService)

			// テスト対象の関数を実行
			result, err := service.UpdateIssue(tc.owner, tc.repo, tc.issueNumber, tc.options)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError && tc.mockResponse != nil {
				compareMaps(t, tc.mockResponse, result)
			}
		})
	}
}

// TestAddIssueComment_ErrorCases はAddIssueComment関数の追加エラーケースをテストする
func TestAddIssueComment_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		issueNumber    int
		body           string
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:           "異常系 - 不正なJSONレスポンス",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    1,
			body:           "これはテストコメントです",
			mockResponse:   nil,
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - サーバーエラー",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    1,
			body:           "これはテストコメントです",
			mockResponse:   nil,
			mockStatusCode: http.StatusInternalServerError,
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

			// GitHubIssueServiceの作成
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := NewGitHubIssueServiceWithDependencies(clientService)

			// テスト対象の関数を実行
			result, err := service.AddIssueComment(tc.owner, tc.repo, tc.issueNumber, tc.body)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError && tc.mockResponse != nil {
				compareMaps(t, tc.mockResponse, result)
			}
		})
	}
}
