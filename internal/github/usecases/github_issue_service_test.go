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

func TestAddToOptions(t *testing.T) {
	tests := []struct {
		name          string
		options       map[string]interface{}
		args          map[string]interface{}
		key           string
		expectedValue interface{}
		expectedExist bool
	}{
		{
			name:          "キーが存在する場合",
			options:       map[string]interface{}{},
			args:          map[string]interface{}{"key": "value"},
			key:           "key",
			expectedValue: "value",
			expectedExist: true,
		},
		{
			name:          "キーが存在しない場合",
			options:       map[string]interface{}{},
			args:          map[string]interface{}{"other": "value"},
			key:           "key",
			expectedValue: nil,
			expectedExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddToOptions(tt.options, tt.args, tt.key)
			value, exists := tt.options[tt.key]
			if exists != tt.expectedExist {
				t.Errorf("AddToOptions() key exists = %v, want %v", exists, tt.expectedExist)
			}
			if exists && value != tt.expectedValue {
				t.Errorf("AddToOptions() value = %v, want %v", value, tt.expectedValue)
			}
		})
	}
}

func TestNewGitHubIssueService(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{
			name:  "トークンあり",
			token: "testtoken",
		},
		{
			name:  "トークンなし",
			token: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := NewGitHubIssueService(tc.token)

			if service == nil {
				t.Fatal("サービスがnilです")
			}

			if service.clientService == nil {
				t.Fatal("クライアントサービスがnilです")
			}
		})
	}
}

func TestCreateIssue(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		options        map[string]interface{}
		mockResponse   map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:  "正常系 - イシュー作成成功",
			owner: "test_user",
			repo:  "test_repo",
			options: map[string]interface{}{
				"title": "テストイシュー",
				"body":  "これはテストイシューです",
			},
			mockResponse: map[string]interface{}{
				"id":     float64(123456),
				"number": float64(1),
				"title":  "テストイシュー",
				"body":   "これはテストイシューです",
				"state":  "open",
			},
			mockStatusCode: http.StatusCreated,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:  "異常系 - 認証エラー",
			owner: "test_user",
			repo:  "test_repo",
			options: map[string]interface{}{
				"title": "テストイシュー",
			},
			mockResponse: map[string]interface{}{
				"message":           "Bad credentials",
				"documentation_url": "https://docs.github.com/rest",
			},
			mockStatusCode: http.StatusUnauthorized,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:  "異常系 - リポジトリが存在しない",
			owner: "nonexistent",
			repo:  "nonexistent",
			options: map[string]interface{}{
				"title": "テストイシュー",
			},
			mockResponse: map[string]interface{}{
				"message":           "Not Found",
				"documentation_url": "https://docs.github.com/rest",
			},
			mockStatusCode: http.StatusNotFound,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:  "異常系 - 不正なJSONレスポンス",
			owner: "test_user",
			repo:  "test_repo",
			options: map[string]interface{}{
				"title": "テストイシュー",
			},
			mockResponse:   nil,
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:  "異常系 - ネットワークエラー",
			owner: "test_user",
			repo:  "test_repo",
			options: map[string]interface{}{
				"title": "テストイシュー",
			},
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      errors.New("ネットワーク接続エラー"),
			expectError:    true,
		},
		{
			name:  "異常系 - JSONマーシャリングエラー",
			owner: "test_user",
			repo:  "test_repo",
			options: map[string]interface{}{
				"title":    "テストイシュー",
				"callback": func() {}, // JSONにシリアライズできない値
			},
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// JSONマーシャリングエラーのテストケース
			if tc.name == "異常系 - JSONマーシャリングエラー" {
				service := NewGitHubIssueService("test_token")
				_, err := service.CreateIssue(tc.owner, tc.repo, tc.options)
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
					// ネットワークエラーのシミュレーション
					if tc.mockError != nil {
						return nil, tc.mockError
					}

					// リクエストの検証
					expectedURL := apiBaseURL + "/repos/" + tc.owner + "/" + tc.repo + "/issues"
					if req.URL.String() != expectedURL {
						t.Errorf("期待されたURL: %s, 実際: %s", expectedURL, req.URL.String())
					}

					if req.Method != "POST" {
						t.Errorf("期待されたHTTPメソッド: POST, 実際: %s", req.Method)
					}

					if req.Header.Get("Accept") != "application/vnd.github.v3+json" {
						t.Errorf("期待されたAcceptヘッダー: application/vnd.github.v3+json, 実際: %s", req.Header.Get("Accept"))
					}

					if req.Header.Get("Content-Type") != "application/json" {
						t.Errorf("期待されたContent-Typeヘッダー: application/json, 実際: %s", req.Header.Get("Content-Type"))
					}

					if req.Header.Get("Authorization") != "token test_token" {
						t.Errorf("期待されたAuthorizationヘッダー: token test_token, 実際: %s", req.Header.Get("Authorization"))
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

			// GitHubIssueServiceのhttpClientをモックに置き換える
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := NewGitHubIssueServiceWithDependencies(clientService)

			// テスト対象の関数を実行
			result, err := service.CreateIssue(tc.owner, tc.repo, tc.options)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError {
				compareMaps(t, tc.mockResponse, result)
			}
		})
	}
}

func TestListIssues(t *testing.T) {
	tests := []struct {
		name           string
		owner          string
		repo           string
		options        map[string]interface{}
		mockResponse   []map[string]interface{}
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:    "正常系 - イシュー一覧取得成功",
			owner:   "test_user",
			repo:    "test_repo",
			options: map[string]interface{}{},
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
					"state":  "closed",
				},
			},
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:  "正常系 - クエリパラメータあり",
			owner: "test_user",
			repo:  "test_repo",
			options: map[string]interface{}{
				"state":     "open",
				"sort":      "created",
				"direction": "desc",
				"per_page":  30,
				"page":      1,
			},
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
			name:    "異常系 - 認証エラー",
			owner:   "test_user",
			repo:    "test_repo",
			options: map[string]interface{}{},
			mockResponse: []map[string]interface{}{
				{
					"message":           "Bad credentials",
					"documentation_url": "https://docs.github.com/rest",
				},
			},
			mockStatusCode: http.StatusUnauthorized,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - 不正なJSONレスポンス",
			owner:          "test_user",
			repo:           "test_repo",
			options:        map[string]interface{}{},
			mockResponse:   nil,
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - ネットワークエラー",
			owner:          "test_user",
			repo:           "test_repo",
			options:        map[string]interface{}{},
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
					// ネットワークエラーのシミュレーション
					if tc.mockError != nil {
						return nil, tc.mockError
					}

					// リクエストの検証
					expectedBaseURL := apiBaseURL + "/repos/" + tc.owner + "/" + tc.repo + "/issues"
					actualURL := req.URL.String()

					// ベースURLの検証
					actualBaseURL := actualURL
					if strings.Contains(actualURL, "?") {
						actualBaseURL = actualURL[:strings.Index(actualURL, "?")]
					}

					if actualBaseURL != expectedBaseURL {
						t.Errorf("期待されたベースURL: %s, 実際: %s", expectedBaseURL, actualBaseURL)
					}

					// クエリパラメータの検証（順序に依存しない）
					if len(tc.options) > 0 {
						// 実際のクエリパラメータを解析
						actualQueryParams := make(map[string]string)
						if strings.Contains(actualURL, "?") {
							queryString := actualURL[strings.Index(actualURL, "?")+1:]
							queryParts := strings.Split(queryString, "&")
							for _, part := range queryParts {
								if strings.Contains(part, "=") {
									kv := strings.Split(part, "=")
									actualQueryParams[kv[0]] = kv[1]
								}
							}
						}

						// 期待されるクエリパラメータと比較
						for k, v := range tc.options {
							expectedValue := fmt.Sprintf("%v", v)
							actualValue, exists := actualQueryParams[k]
							if !exists {
								t.Errorf("クエリパラメータ %s が見つかりません", k)
							} else if actualValue != expectedValue {
								t.Errorf("クエリパラメータ %s の値が異なります。期待: %s, 実際: %s", k, expectedValue, actualValue)
							}
						}

						// 余分なクエリパラメータがないことを確認
						if len(actualQueryParams) != len(tc.options) {
							t.Errorf("クエリパラメータの数が異なります。期待: %d, 実際: %d", len(tc.options), len(actualQueryParams))
						}
					}

					if req.Method != "GET" {
						t.Errorf("期待されたHTTPメソッド: GET, 実際: %s", req.Method)
					}

					if req.Header.Get("Accept") != "application/vnd.github.v3+json" {
						t.Errorf("期待されたAcceptヘッダー: application/vnd.github.v3+json, 実際: %s", req.Header.Get("Accept"))
					}

					if req.Header.Get("Authorization") != "token test_token" {
						t.Errorf("期待されたAuthorizationヘッダー: token test_token, 実際: %s", req.Header.Get("Authorization"))
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

			// GitHubIssueServiceのhttpClientをモックに置き換える
			clientService := NewGitHubClientServiceWithDependencies(mockClient, "test_token", &DefaultJSONMarshaler{})
			service := NewGitHubIssueServiceWithDependencies(clientService)

			// テスト対象の関数を実行
			result, err := service.ListIssues(tc.owner, tc.repo, tc.options)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、結果を検証
			if !tc.expectError {
				if len(result) != len(tc.mockResponse) {
					t.Errorf("期待された結果の長さ: %d, 実際: %d", len(tc.mockResponse), len(result))
				}

				// 各アイテムを検証
				for i, expectedItem := range tc.mockResponse {
					if i >= len(result) {
						t.Errorf("インデックス %d の結果アイテムが見つかりません", i)
						continue
					}
					actualItem := result[i]
					compareMaps(t, expectedItem, actualItem)
				}
			}
		})
	}
}

func TestUpdateIssue(t *testing.T) {
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
			name:        "正常系 - イシュー更新成功",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			options: map[string]interface{}{
				"title": "更新されたイシュー",
				"body":  "これは更新されたイシューです",
				"state": "closed",
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
			name:        "異常系 - 認証エラー",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			options: map[string]interface{}{
				"title": "更新されたイシュー",
			},
			mockResponse: map[string]interface{}{
				"message":           "Bad credentials",
				"documentation_url": "https://docs.github.com/rest",
			},
			mockStatusCode: http.StatusUnauthorized,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:        "異常系 - JSONマーシャリングエラー",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			options: map[string]interface{}{
				"title":    "更新されたイシュー",
				"callback": func() {}, // JSONにシリアライズできない値
			},
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// JSONマーシャリングエラーのテストケース
			if tc.name == "異常系 - JSONマーシャリングエラー" {
				service := NewGitHubIssueService("test_token")
				_, err := service.UpdateIssue(tc.owner, tc.repo, tc.issueNumber, tc.options)
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
					// ネットワークエラーのシミュレーション
					if tc.mockError != nil {
						return nil, tc.mockError
					}

					// リクエストの検証
					expectedURL := fmt.Sprintf("%s/repos/%s/%s/issues/%d", apiBaseURL, tc.owner, tc.repo, tc.issueNumber)
					if req.URL.String() != expectedURL {
						t.Errorf("期待されたURL: %s, 実際: %s", expectedURL, req.URL.String())
					}

					if req.Method != "PATCH" {
						t.Errorf("期待されたHTTPメソッド: PATCH, 実際: %s", req.Method)
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

			// GitHubIssueServiceのhttpClientをモックに置き換える
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
			if !tc.expectError {
				compareMaps(t, tc.mockResponse, result)
			}
		})
	}
}

func TestAddIssueComment(t *testing.T) {
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
			},
			mockStatusCode: http.StatusCreated,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:        "異常系 - 認証エラー",
			owner:       "test_user",
			repo:        "test_repo",
			issueNumber: 1,
			body:        "これはテストコメントです",
			mockResponse: map[string]interface{}{
				"message":           "Bad credentials",
				"documentation_url": "https://docs.github.com/rest",
			},
			mockStatusCode: http.StatusUnauthorized,
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
					// ネットワークエラーのシミュレーション
					if tc.mockError != nil {
						return nil, tc.mockError
					}

					// リクエストの検証
					expectedURL := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", apiBaseURL, tc.owner, tc.repo, tc.issueNumber)
					if req.URL.String() != expectedURL {
						t.Errorf("期待されたURL: %s, 実際: %s", expectedURL, req.URL.String())
					}

					if req.Method != "POST" {
						t.Errorf("期待されたHTTPメソッド: POST, 実際: %s", req.Method)
					}

					// リクエストボディの検証
					body, _ := io.ReadAll(req.Body)
					var requestBody map[string]interface{}
					if err := json.Unmarshal(body, &requestBody); err != nil {
						t.Fatalf("リクエストボディのJSONパースに失敗しました: %v", err)
					}

					// bodyパラメータの検証
					if requestBody["body"] != tc.body {
						t.Errorf("期待されたbody: %s, 実際: %s", tc.body, requestBody["body"])
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

			// GitHubIssueServiceのhttpClientをモックに置き換える
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
			if !tc.expectError {
				compareMaps(t, tc.mockResponse, result)
			}
		})
	}
}
