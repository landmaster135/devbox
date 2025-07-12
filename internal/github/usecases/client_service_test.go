package usecases

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

// TestNewGitHubClientService は NewGitHubClientService 関数をテストする
func TestNewGitHubClientService(t *testing.T) {
	// テストケース
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
			// 関数を実行
			service := NewGitHubClientService(tc.token)

			// 検証
			if service == nil {
				t.Fatal("サービスがnilです")
			}

			if service.token != tc.token {
				t.Errorf("期待されたトークン: %s, 実際: %s", tc.token, service.token)
			}

			if service.httpClient == nil {
				t.Fatal("HTTPクライアントがnilです")
			}

			if service.jsonMarshaler == nil {
				t.Fatal("JSONマーシャラーがnilです")
			}
		})
	}
}

// TestGitHubErrorError はGitHubError構造体のErrorメソッドをテストする
func TestGitHubErrorError(t *testing.T) {
	tests := []struct {
		name     string
		ghError  GitHubError
		expected string
	}{
		{
			name: "基本的なエラーメッセージ",
			ghError: GitHubError{
				Message:    "Not Found",
				StatusCode: 404,
			},
			expected: "GitHub API Error: Not Found (Status: 404)",
		},
		{
			name: "ドキュメントURLを含むエラー",
			ghError: GitHubError{
				Message:          "Validation Failed",
				DocumentationURL: "https://docs.github.com/rest/reference/issues",
				StatusCode:       422,
			},
			expected: "GitHub API Error: Validation Failed (Status: 422)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errorMsg := tt.ghError.Error()
			if errorMsg != tt.expected {
				t.Errorf("GitHubError.Error() = %v, want %v", errorMsg, tt.expected)
			}
		})
	}
}

// TestReturnJSONResult はReturnJSONResultメソッドをテストする
func TestReturnJSONResult(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name:        "正常系 - 有効なJSONデータ",
			input:       map[string]interface{}{"key": "value"},
			expectError: false,
		},
		{
			name: "異常系 - JSONにマーシャルできないデータ",
			input: func() interface{} {
				// JSONにマーシャルできない循環参照を持つデータ構造
				type Circular struct {
					Self *Circular
				}
				c := &Circular{}
				c.Self = c
				return c
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewGitHubClientService("test_token")
			result, err := service.ReturnJSONResult(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("ReturnJSONResult() エラーが期待されていましたが、エラーは発生しませんでした")
				}
				// エラーが発生した場合、resultは空文字列であるべき
				if result != "" {
					t.Errorf("ReturnJSONResult() エラー時に空文字列ではない結果が返されました: %v", result)
				}
			} else {
				if err != nil {
					t.Errorf("ReturnJSONResult() error = %v", err)
					return
				}
				// 正常系の場合、resultは空文字列ではないはず
				if result == "" {
					t.Errorf("ReturnJSONResult() 結果が空文字列です")
				}
			}
		})
	}
}

// TestDoRequest は DoRequest メソッドを直接テストする
func TestDoRequest(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		body           io.Reader
		token          string
		mockResponse   []byte
		mockStatusCode int
		mockError      error
		expectError    bool
	}{
		{
			name:           "正常系 - GETリクエスト",
			method:         "GET",
			url:            "https://api.github.com/user",
			body:           nil,
			token:          "testtoken",
			mockResponse:   []byte(`{"login": "testuser", "id": 12345}`),
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "正常系 - POSTリクエスト",
			method:         "POST",
			url:            "https://api.github.com/repos/owner/repo/issues",
			body:           strings.NewReader(`{"title": "Test"}`),
			token:          "testtoken",
			mockResponse:   []byte(`{"id": 123, "number": 1, "title": "Test"}`),
			mockStatusCode: http.StatusCreated,
			mockError:      nil,
			expectError:    false,
		},
		{
			name:           "異常系 - 無効なリクエスト",
			method:         "GET",
			url:            "://invalid-url",
			body:           nil,
			token:          "testtoken",
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      nil,
			expectError:    true,
		},
		{
			name:           "異常系 - リクエスト読み取りエラー",
			method:         "POST",
			url:            "https://api.github.com/repos/owner/repo/issues",
			body:           &errorReader{},
			token:          "testtoken",
			mockResponse:   nil,
			mockStatusCode: 0,
			mockError:      errors.New("リクエストボディ読み取りエラー"),
			expectError:    true,
		},
		{
			name:           "異常系 - レスポンス読み取りエラー",
			method:         "GET",
			url:            "https://api.github.com/user",
			body:           nil,
			token:          "testtoken",
			mockResponse:   nil,
			mockStatusCode: http.StatusOK,
			mockError:      nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 無効なURLテスト
			if tc.name == "異常系 - 無効なリクエスト" {
				service := NewGitHubClientService(tc.token)
				_, err := service.DoRequest(tc.method, tc.url, tc.body)
				if err == nil {
					t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
				}
				return
			}

			// モックHTTPクライアントの作成
			mockClient := &MockHTTPClient{
				DoFunc: func(req *http.Request) (*http.Response, error) {
					// errorReaderのテストでは常にエラーを返す
					if tc.name == "異常系 - リクエスト読み取りエラー" {
						return nil, tc.mockError
					}

					// メソッドの検証
					if req.Method != tc.method {
						t.Errorf("期待されたHTTPメソッド: %s, 実際: %s", tc.method, req.Method)
					}

					// URLの検証
					if req.URL.String() != tc.url {
						t.Errorf("期待されたURL: %s, 実際: %s", tc.url, req.URL.String())
					}

					// ヘッダーの検証
					if req.Header.Get("Accept") != "application/vnd.github.v3+json" {
						t.Errorf("期待されたAcceptヘッダー: application/vnd.github.v3+json, 実際: %s", req.Header.Get("Accept"))
					}

					if tc.token != "" && req.Header.Get("Authorization") != "token "+tc.token {
						t.Errorf("期待されたAuthorizationヘッダー: token %s, 実際: %s", tc.token, req.Header.Get("Authorization"))
					}

					if (tc.method == "POST" || tc.method == "PATCH" || tc.method == "PUT") && req.Header.Get("Content-Type") != "application/json" {
						t.Errorf("期待されたContent-Typeヘッダー: application/json, 実際: %s", req.Header.Get("Content-Type"))
					}

					// エラーのシミュレーション
					if tc.mockError != nil {
						return nil, tc.mockError
					}

					// レスポンス読み取りエラーのシミュレーション
					if tc.name == "異常系 - レスポンス読み取りエラー" {
						return &http.Response{
							StatusCode: tc.mockStatusCode,
							Body:       &errorReadCloser{},
						}, nil
					}

					// 通常のレスポンス
					return &http.Response{
						StatusCode: tc.mockStatusCode,
						Body:       io.NopCloser(bytes.NewReader(tc.mockResponse)),
					}, nil
				},
			}

			// GitHubClientServiceのhttpClientをモックに置き換える
			service := NewGitHubClientServiceWithDependencies(mockClient, tc.token, &DefaultJSONMarshaler{})

			// テスト対象の関数を実行
			data, err := service.DoRequest(tc.method, tc.url, tc.body)

			// エラーの検証
			if tc.expectError && err == nil {
				t.Error("エラーが期待されていましたが、エラーは発生しませんでした")
			}
			if !tc.expectError && err != nil {
				t.Errorf("エラーは期待されていませんでしたが、エラーが発生しました: %v", err)
			}

			// 正常系の場合、レスポンスを検証
			if !tc.expectError {
				if !bytes.Equal(data, tc.mockResponse) {
					t.Errorf("期待されたレスポンス: %s, 実際: %s", string(tc.mockResponse), string(data))
				}
			}
		})
	}
}
