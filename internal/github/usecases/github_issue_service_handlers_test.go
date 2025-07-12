package usecases

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

// TestHandleToCreateIssue はHandleToCreateIssue関数をテストする
func TestHandleToCreateIssue(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		title          string
		body           string
		labels         []interface{}
		assignees      []interface{}
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:  "正常系 - 基本的なイシュー作成",
			owner: "test_user",
			repo:  "test_repo",
			title: "テストイシュー",
			body:  "これはテストイシューです",
			labels: []interface{}{
				"bug",
				"enhancement",
			},
			assignees: []interface{}{
				"test_user",
			},
			mockResponse: map[string]interface{}{
				"id":     float64(123456),
				"number": float64(1),
				"title":  "テストイシュー",
				"body":   "これはテストイシューです",
				"state":  "open",
				"labels": []interface{}{
					map[string]interface{}{
						"name": "bug",
					},
				},
			},
			mockStatusCode: http.StatusCreated,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:      "正常系 - タイトルのみ",
			owner:     "test_user",
			repo:      "test_repo",
			title:     "シンプルなイシュー",
			body:      "",
			labels:    nil,
			assignees: nil,
			mockResponse: map[string]interface{}{
				"id":     float64(123457),
				"number": float64(2),
				"title":  "シンプルなイシュー",
				"state":  "open",
			},
			mockStatusCode: http.StatusCreated,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - 認証エラー",
			owner:          "test_user",
			repo:           "test_repo",
			title:          "テストイシュー",
			body:           "",
			labels:         nil,
			assignees:      nil,
			mockResponse:   nil,
			mockStatusCode: http.StatusUnauthorized,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - リポジトリが存在しない",
			owner:          "nonexistent",
			repo:           "nonexistent",
			title:          "テストイシュー",
			body:           "",
			labels:         nil,
			assignees:      nil,
			mockResponse:   nil,
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - ネットワークエラー",
			owner:          "test_user",
			repo:           "test_repo",
			title:          "テストイシュー",
			body:           "",
			labels:         nil,
			assignees:      nil,
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
					if tc.mockResponse != nil {
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
			result, err := service.HandleToCreateIssue(tc.owner, tc.repo, tc.title, tc.body, tc.labels, tc.assignees)

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
				var parsedResult map[string]interface{}
				if err := json.Unmarshal([]byte(result), &parsedResult); err != nil {
					t.Errorf("結果のJSONパースに失敗しました: %v", err)
				}
			}
		})
	}
}

// TestHandleToListIssues はHandleToListIssues関数をテストする
func TestHandleToListIssues(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		state          string
		sort           string
		direction      string
		perPage        int
		page           int
		mockResponse   []map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:      "正常系 - 基本的なイシュー一覧取得",
			owner:     "test_user",
			repo:      "test_repo",
			state:     "open",
			sort:      "created",
			direction: "desc",
			perPage:   30,
			page:      1,
			mockResponse: []map[string]interface{}{
				{
					"id":     float64(123456),
					"number": float64(1),
					"title":  "テストイシュー1",
					"state":  "open",
				},
				{
					"id":     float64(123457),
					"number": float64(2),
					"title":  "テストイシュー2",
					"state":  "open",
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:      "正常系 - パラメータなし",
			owner:     "test_user",
			repo:      "test_repo",
			state:     "",
			sort:      "",
			direction: "",
			perPage:   30,
			page:      1,
			mockResponse: []map[string]interface{}{
				{
					"id":     float64(123456),
					"number": float64(1),
					"title":  "テストイシュー1",
					"state":  "open",
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - 認証エラー",
			owner:          "test_user",
			repo:           "test_repo",
			state:          "",
			sort:           "",
			direction:      "",
			perPage:        30,
			page:           1,
			mockResponse:   nil,
			mockStatusCode: http.StatusUnauthorized,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - リポジトリが存在しない",
			owner:          "nonexistent",
			repo:           "nonexistent",
			state:          "",
			sort:           "",
			direction:      "",
			perPage:        30,
			page:           1,
			mockResponse:   nil,
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - ネットワークエラー",
			owner:          "test_user",
			repo:           "test_repo",
			state:          "",
			sort:           "",
			direction:      "",
			perPage:        30,
			page:           1,
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
					if tc.mockResponse != nil {
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
			result, err := service.HandleToListIssues(tc.owner, tc.repo, tc.state, tc.sort, tc.direction, tc.perPage, tc.page)

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

// TestHandleToUpdateIssue はHandleToUpdateIssue関数をテストする
func TestHandleToUpdateIssue(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		issueNumber    int
		title          string
		body           string
		state          string
		labels         []interface{}
		assignees      []interface{}
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:        "正常系 - 完全なイシュー更新",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			title:       "更新されたイシュー",
			body:        "これは更新されたイシューです",
			state:       "closed",
			labels: []interface{}{
				"bug",
				"fixed",
			},
			assignees: []interface{}{
				"test_user",
			},
			mockResponse: map[string]interface{}{
				"id":     float64(123456),
				"number": float64(1),
				"title":  "更新されたイシュー",
				"body":   "これは更新されたイシューです",
				"state":  "closed",
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:        "正常系 - 部分的な更新",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			title:       "新しいタイトル",
			body:        "",
			state:       "",
			labels:      nil,
			assignees:   nil,
			mockResponse: map[string]interface{}{
				"id":     float64(123456),
				"number": float64(1),
				"title":  "新しいタイトル",
				"state":  "open",
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - 認証エラー",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    1,
			title:          "更新されたイシュー",
			body:           "",
			state:          "",
			labels:         nil,
			assignees:      nil,
			mockResponse:   nil,
			mockStatusCode: http.StatusUnauthorized,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - イシューが存在しない",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    999,
			title:          "更新されたイシュー",
			body:           "",
			state:          "",
			labels:         nil,
			assignees:      nil,
			mockResponse:   nil,
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - ネットワークエラー",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    1,
			title:          "更新されたイシュー",
			body:           "",
			state:          "",
			labels:         nil,
			assignees:      nil,
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
					if tc.mockResponse != nil {
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
			result, err := service.HandleToUpdateIssue(tc.owner, tc.repo, tc.issueNumber, tc.title, tc.body, tc.state, tc.labels, tc.assignees)

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
				var parsedResult map[string]interface{}
				if err := json.Unmarshal([]byte(result), &parsedResult); err != nil {
					t.Errorf("結果のJSONパースに失敗しました: %v", err)
				}
			}
		})
	}
}

// TestHandleToAddIssueComment はHandleToAddIssueComment関数をテストする
func TestHandleToAddIssueComment(t *testing.T) {
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
			name:        "正常系 - コメント追加成功",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			body:        "これはテストコメントです",
			mockResponse: map[string]interface{}{
				"id":   float64(123456),
				"body": "これはテストコメントです",
				"user": map[string]interface{}{
					"login": "test_user",
				},
				"created_at": "2023-01-01T12:00:00Z",
			},
			mockStatusCode: http.StatusCreated,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:        "正常系 - 長いコメント",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			body:        "これは非常に長いテストコメントです。複数行にわたって記述されており、様々な情報を含んでいます。\n\n- 項目1\n- 項目2\n- 項目3\n\n詳細な説明がここに続きます。",
			mockResponse: map[string]interface{}{
				"id":   float64(123457),
				"body": "これは非常に長いテストコメントです。複数行にわたって記述されており、様々な情報を含んでいます。\n\n- 項目1\n- 項目2\n- 項目3\n\n詳細な説明がここに続きます。",
				"user": map[string]interface{}{
					"login": "test_user",
				},
			},
			mockStatusCode: http.StatusCreated,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - 認証エラー",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    1,
			body:           "これはテストコメントです",
			mockResponse:   nil,
			mockStatusCode: http.StatusUnauthorized,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - イシューが存在しない",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    999,
			body:           "これはテストコメントです",
			mockResponse:   nil,
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - 空のコメント本文",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    1,
			body:           "",
			mockResponse:   nil,
			mockStatusCode: http.StatusUnprocessableEntity,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - ネットワークエラー",
			owner:          "test_user",
			repo:           "test_repo",
			issueNumber:    1,
			body:           "これはテストコメントです",
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
					if tc.mockResponse != nil {
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
			result, err := service.HandleToAddIssueComment(tc.owner, tc.repo, tc.issueNumber, tc.body)

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
				var parsedResult map[string]interface{}
				if err := json.Unmarshal([]byte(result), &parsedResult); err != nil {
					t.Errorf("結果のJSONパースに失敗しました: %v", err)
				}
			}
		})
	}
}
