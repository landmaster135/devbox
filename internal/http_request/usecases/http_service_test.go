package services

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"

	models "github.com/landmaster135/devbox/internal/http_request/domain/models"
)

// ///////////////////////
// MockHTTPRepository はHTTPRepositoryのモック実装です
// ///////////////////////
type MockHTTPRepository struct {
	LoadJSONFileFunc func(filePath string) ([]byte, error)
	SendRequestFunc  func(request *models.HTTPRequest) (*models.HTTPResponse, error)
}

// LoadJSONFile はJSONファイルを読み込むモックメソッドです
func (m *MockHTTPRepository) LoadJSONFile(filePath string) ([]byte, error) {
	return m.LoadJSONFileFunc(filePath)
}

// SendRequest はHTTPリクエストを送信するモックメソッドです
func (m *MockHTTPRepository) SendRequest(request *models.HTTPRequest) (*models.HTTPResponse, error) {
	return m.SendRequestFunc(request)
}

// TestNewHTTPService はNewHTTPServiceメソッドのテストです
func TestNewHTTPService(t *testing.T) {
	// Arrange
	mockRepo := &MockHTTPRepository{}

	// Act
	service := NewHTTPService(mockRepo)

	// Assert
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
	if service.httpRepo != mockRepo {
		t.Errorf("Expected apiRepo to be %v, got %v", mockRepo, service.httpRepo)
	}
}

// TestSendRequestWithJSONFile はSendRequestWithJSONFileメソッドのテストです
func TestSendRequestWithJSONFile(t *testing.T) {
	// テストケース
	testCases := []struct {
		name         string
		url          string
		method       string
		jsonFilePath string
		mockJSON     []byte
		mockResponse *models.HTTPResponse
		mockError    error
		expectError  bool
	}{
		{
			name:         "正常系",
			url:          "https://api.example.com",
			method:       "POST",
			jsonFilePath: "test.json",
			mockJSON:     []byte(`{"test": "data"}`),
			mockResponse: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"success": true}`),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:         "JSONファイル読み込みエラー",
			url:          "https://api.example.com",
			method:       "POST",
			jsonFilePath: "invalid.json",
			mockJSON:     nil,
			mockResponse: nil,
			mockError:    errors.New("ファイルが見つかりません"),
			expectError:  true,
		},
		{
			name:         "HTTPリクエスト送信エラー",
			url:          "https://api.example.com",
			method:       "POST",
			jsonFilePath: "test.json",
			mockJSON:     []byte(`{"test": "data"}`),
			mockResponse: nil,
			mockError:    errors.New("接続エラー"),
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := &MockHTTPRepository{
				LoadJSONFileFunc: func(filePath string) ([]byte, error) {
					if tc.mockError != nil && filePath == "invalid.json" {
						return nil, tc.mockError
					}
					return tc.mockJSON, nil
				},
				SendRequestFunc: func(request *models.HTTPRequest) (*models.HTTPResponse, error) {
					// URLとメソッドが正しく設定されているか確認
					if request.URL != tc.url {
						t.Errorf("Expected URL %s, got %s", tc.url, request.URL)
					}
					if request.Method != tc.method {
						t.Errorf("Expected method %s, got %s", tc.method, request.Method)
					}

					// Content-Typeヘッダーが設定されているか確認
					if value, exists := request.Headers["Content-Type"]; !exists || value != "application/json" {
						t.Errorf("Expected Content-Type header with value application/json, got %s", value)
					}
					if value, exists := request.Headers["Accept"]; !exists || value != "application/json" {
						t.Errorf("Expected Accept header with value application/json, got %s", value)
					}

					if tc.name == "HTTPリクエスト送信エラー" {
						return nil, tc.mockError
					}

					return tc.mockResponse, nil
				},
			}

			// テスト対象のサービスを作成
			service := NewHTTPService(mockRepo)

			// メソッドを実行
			response, err := service.SendRequestWithJSONFile(tc.url, tc.method, tc.jsonFilePath)

			// エラー発生の期待値と実際の結果を比較
			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				// レスポンスが期待通りであることを確認
				if !reflect.DeepEqual(response, tc.mockResponse) {
					t.Errorf("Expected response %+v, got %+v", tc.mockResponse, response)
				}
			}
		})
	}
}

// TestSendRequestWithJSONFileAndHeaders_Encoding はエンコーディング指定のテストです
func TestSendRequestWithJSONFileAndHeaders_Encoding(t *testing.T) {
	// テストケース
	testCases := []struct {
		name             string
		url              string
		method           string
		jsonFilePath     string
		headers          map[string]string
		encoding         string
		expectedEncoding string
	}{
		{
			name:             "Shift_JIS指定",
			url:              "https://api.example.com",
			method:           "POST",
			jsonFilePath:     "test.json",
			headers:          map[string]string{"Accept": "application/json"},
			encoding:         "shift_jis",
			expectedEncoding: "shift_jis",
		},
		{
			name:             "UTF-8指定",
			url:              "https://api.example.com",
			method:           "POST",
			jsonFilePath:     "test.json",
			headers:          map[string]string{"Accept": "application/json"},
			encoding:         "utf-8",
			expectedEncoding: "utf-8",
		},
		{
			name:             "EUC-JP指定",
			url:              "https://api.example.com",
			method:           "POST",
			jsonFilePath:     "test.json",
			headers:          map[string]string{"Accept": "application/json"},
			encoding:         "euc-jp",
			expectedEncoding: "euc-jp",
		},
		{
			name:             "auto指定",
			url:              "https://api.example.com",
			method:           "POST",
			jsonFilePath:     "test.json",
			headers:          map[string]string{"Accept": "application/json"},
			encoding:         "auto",
			expectedEncoding: "auto",
		},
		{
			name:             "空文字列指定（autoとして扱われる）",
			url:              "https://api.example.com",
			method:           "POST",
			jsonFilePath:     "test.json",
			headers:          map[string]string{"Accept": "application/json"},
			encoding:         "",
			expectedEncoding: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := &MockHTTPRepository{
				LoadJSONFileFunc: func(filePath string) ([]byte, error) {
					return []byte(`{"test": "data"}`), nil
				},
				SendRequestFunc: func(request *models.HTTPRequest) (*models.HTTPResponse, error) {
					// エンコーディングが正しく設定されているか確認
					if request.Encoding != tc.expectedEncoding {
						t.Errorf("Expected encoding %s, got %s", tc.expectedEncoding, request.Encoding)
					}

					// URLとメソッドが正しく設定されているか確認
					if request.URL != tc.url {
						t.Errorf("Expected URL %s, got %s", tc.url, request.URL)
					}
					if request.Method != tc.method {
						t.Errorf("Expected method %s, got %s", tc.method, request.Method)
					}

					// モックレスポンスを返す
					return &models.HTTPResponse{
						StatusCode: 200,
						Headers:    map[string]string{"Content-Type": "application/json"},
						Body:       []byte(`{"success": true}`),
					}, nil
				},
			}

			// テスト対象のサービスを作成
			service := NewHTTPService(mockRepo)

			// メソッドを実行
			response, err := service.SendRequestWithJSONFileAndHeaders(tc.url, tc.method, tc.jsonFilePath, tc.headers, tc.encoding)

			// エラーがないことを確認
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// レスポンスが期待通りであることを確認
			expectedResponse := &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"success": true}`),
			}

			if !reflect.DeepEqual(response, expectedResponse) {
				t.Errorf("Expected response %+v, got %+v", expectedResponse, response)
			}
		})
	}
}

// TestSendRequestWithoutJSONFile_Encoding はエンコーディング指定のテストです
func TestSendRequestWithoutJSONFile_Encoding(t *testing.T) {
	// テストケース
	testCases := []struct {
		name             string
		url              string
		method           string
		headers          map[string]string
		encoding         string
		expectedEncoding string
	}{
		{
			name:             "Shift_JIS指定",
			url:              "https://example.com",
			method:           "GET",
			headers:          map[string]string{"Accept": "text/html"},
			encoding:         "shift_jis",
			expectedEncoding: "shift_jis",
		},
		{
			name:             "UTF-8指定",
			url:              "https://example.com",
			method:           "GET",
			headers:          map[string]string{"Accept": "text/html"},
			encoding:         "utf-8",
			expectedEncoding: "utf-8",
		},
		{
			name:             "EUC-JP指定",
			url:              "https://example.com",
			method:           "GET",
			headers:          map[string]string{"Accept": "text/html"},
			encoding:         "euc-jp",
			expectedEncoding: "euc-jp",
		},
		{
			name:             "auto指定",
			url:              "https://example.com",
			method:           "GET",
			headers:          map[string]string{"Accept": "text/html"},
			encoding:         "auto",
			expectedEncoding: "auto",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := &MockHTTPRepository{
				SendRequestFunc: func(request *models.HTTPRequest) (*models.HTTPResponse, error) {
					// エンコーディングが正しく設定されているか確認
					if request.Encoding != tc.expectedEncoding {
						t.Errorf("Expected encoding %s, got %s", tc.expectedEncoding, request.Encoding)
					}

					// URLとメソッドが正しく設定されているか確認
					if request.URL != tc.url {
						t.Errorf("Expected URL %s, got %s", tc.url, request.URL)
					}
					if request.Method != tc.method {
						t.Errorf("Expected method %s, got %s", tc.method, request.Method)
					}

					// ボディがnilであることを確認
					if request.Body != nil {
						t.Errorf("Expected body to be nil, got %v", request.Body)
					}

					// モックレスポンスを返す
					return &models.HTTPResponse{
						StatusCode: 200,
						Headers:    map[string]string{"Content-Type": "text/html; charset=shift_jis"},
						Body:       []byte(`<html><head><title>テストページ</title></head><body>こんにちは</body></html>`),
					}, nil
				},
			}

			// テスト対象のサービスを作成
			service := NewHTTPService(mockRepo)

			// メソッドを実行
			response, err := service.SendRequestWithoutJSONFile(tc.url, tc.method, tc.headers, tc.encoding)

			// エラーがないことを確認
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// レスポンスが期待通りであることを確認
			expectedResponse := &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "text/html; charset=shift_jis"},
				Body:       []byte(`<html><head><title>テストページ</title></head><body>こんにちは</body></html>`),
			}

			if !reflect.DeepEqual(response, expectedResponse) {
				t.Errorf("Expected response %+v, got %+v", expectedResponse, response)
			}
		})
	}
}

// TestSendRequestWithoutJSONFile はSendRequestWithoutJSONFileメソッドのテストです
func TestSendRequestWithoutJSONFile(t *testing.T) {
	// テストケース
	testCases := []struct {
		name         string
		url          string
		method       string
		headers      map[string]string
		mockResponse *models.HTTPResponse
		mockError    error
		expectError  bool
	}{
		{
			name:   "正常系 - GETリクエスト",
			url:    "https://api.example.com/users",
			method: "GET",
			headers: map[string]string{
				"Accept": "application/json",
			},
			mockResponse: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"users": []}`),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:   "正常系 - POSTリクエスト（ボディなし）",
			url:    "https://api.example.com/ping",
			method: "POST",
			headers: map[string]string{
				"Authorization": "Bearer token123",
			},
			mockResponse: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       []byte("pong"),
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:   "正常系 - DELETEリクエスト",
			url:    "https://api.example.com/users/123",
			method: "DELETE",
			headers: map[string]string{
				"Authorization": "Bearer token123",
			},
			mockResponse: &models.HTTPResponse{
				StatusCode: 204,
				Headers:    map[string]string{},
				Body:       []byte{},
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:   "エラー系 - APIリクエスト送信失敗",
			url:    "https://api.example.com/error",
			method: "GET",
			headers: map[string]string{
				"Accept": "application/json",
			},
			mockResponse: nil,
			mockError:    errors.New("接続エラー"),
			expectError:  true,
		},
		{
			name:    "正常系 - 空のヘッダー",
			url:     "https://api.example.com/public",
			method:  "GET",
			headers: map[string]string{},
			mockResponse: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"message": "public data"}`),
			},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := &MockHTTPRepository{
				SendRequestFunc: func(request *models.HTTPRequest) (*models.HTTPResponse, error) {
					// URLとメソッドが正しく設定されているか確認
					if request.URL != tc.url {
						t.Errorf("Expected URL %s, got %s", tc.url, request.URL)
					}
					if request.Method != tc.method {
						t.Errorf("Expected method %s, got %s", tc.method, request.Method)
					}

					// ヘッダーが正しく設定されているか確認
					for key, expectedValue := range tc.headers {
						if actualValue, exists := request.Headers[key]; !exists || actualValue != expectedValue {
							t.Errorf("Expected header %s with value %s, got %s", key, expectedValue, actualValue)
						}
					}

					// ボディがnilであることを確認
					if request.Body != nil {
						t.Errorf("Expected body to be nil, got %v", request.Body)
					}

					if tc.mockError != nil {
						return nil, tc.mockError
					}

					return tc.mockResponse, nil
				},
			}

			// テスト対象のサービスを作成
			service := NewHTTPService(mockRepo)

			// メソッドを実行
			response, err := service.SendRequestWithoutJSONFile(tc.url, tc.method, tc.headers, "auto")

			// エラー発生の期待値と実際の結果を比較
			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				// レスポンスが期待通りであることを確認
				if !reflect.DeepEqual(response, tc.mockResponse) {
					t.Errorf("Expected response %+v, got %+v", tc.mockResponse, response)
				}
			}
		})
	}
}

// MockJSONEncoder はJSONエンコーダーをモックするための構造体です
type MockJSONEncoder struct {
	EncodeFunc func(v interface{}) error
}

// Encode はJSONエンコードをモックするメソッドです
func (m *MockJSONEncoder) Encode(v interface{}) error {
	return m.EncodeFunc(v)
}

// TestFormatResponse はFormatResponseメソッドのテストです
func TestFormatResponse(t *testing.T) {
	// テストケース
	testCases := []struct {
		name           string
		response       *models.HTTPResponse
		expectedPrefix string
		expectError    bool
		notContains    []string
	}{
		{
			name: "正常系 - JSONレスポンス",
			response: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"success": true, "message": "OK"}`),
			},
			expectedPrefix: "Status: 200",
			expectError:    false,
		},
		{
			name: "正常系 - 空のボディ",
			response: &models.HTTPResponse{
				StatusCode: 204,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte{},
			},
			expectedPrefix: "Status: 204",
			expectError:    false,
		},
		{
			name: "正常系 - テキストレスポンス",
			response: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       []byte("Plain text response"),
			},
			expectedPrefix: "Status: 200",
			expectError:    false,
		},
		{
			name: "正常系 - 無効なJSONレスポンス（テキストとして処理）",
			response: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"invalid": json}`),
			},
			expectedPrefix: "Status: 200",
			expectError:    false,
		},
		{
			name: "正常系 - 空のヘッダー",
			response: &models.HTTPResponse{
				StatusCode: 500,
				Headers:    map[string]string{},
				Body:       []byte(`Internal Server Error`),
			},
			expectedPrefix: "Status: 500",
			expectError:    false,
		},
		{
			name: "正常系 - 複数のヘッダー",
			response: &models.HTTPResponse{
				StatusCode: 201,
				Headers: map[string]string{
					"Content-Type":  "application/json",
					"Location":      "/api/users/123",
					"Cache-Control": "no-cache",
				},
			},
			expectedPrefix: "Status: 201",
			expectError:    false,
		},
		{
			name: "正常系 - 大きなJSONレスポンス",
			response: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"users": [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}, {"id": 3, "name": "Charlie"}], "total": 3, "page": 1}`),
			},
			expectedPrefix: "Status: 200",
			expectError:    false,
		},
		{
			name: "正常系 - 特殊文字を含むレスポンス",
			response: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"message": "こんにちは世界", "emoji": "🌍", "special": "\"quoted\""}`),
			},
			expectedPrefix: "Status: 200",
			expectError:    false,
		},
		{
			name: "警告付きレスポンス",
			response: &models.HTTPResponse{
				StatusCode: 403,
				Headers:    map[string]string{"Content-Type": "text/html"},
				Body:       []byte("<html></html>"),
				Warnings:   []string{"Cloudflare warning detected"},
			},
			expectedPrefix: "Status: 403",
			expectError:    false,
		},
		{
			name: "HTMLボディから不要なタグを除去",
			response: &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "text/html"},
				Body: []byte(
					`<html><body><!-- comment --><header>brand</header><div data-testid="wrapper"><main class="content" style="color:red"><script>console.log('xss')</script><p>text</p><!--secret--><button>click</button><span data-allow-mismatch="true" data-allow-missmatch="true" data-testid="main-span">keep</span></main></div><footer>copyright</footer></body></html>`),
			},
			expectedPrefix: "Status: 200",
			expectError:    false,
			notContains: []string{
				"<script", "console.log", "class=", "style=", "<button", "<!--", "secret", "data-allow-missmatch", "data-allow-mismatch", "data-testid", "<header", "brand", "<footer", "copyright",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := &MockHTTPRepository{}

			// テスト対象のサービスを作成
			service := NewHTTPService(mockRepo)

			// メソッドを実行
			result, err := service.FormatResponse(tc.response)

			// エラー発生の期待値と実際の結果を比較
			if tc.expectError {
				if err == nil {
					t.Fatalf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				// 結果が期待通りであることを確認
				if result == "" {
					t.Error("Expected non-empty result")
				}

				// 期待されるプレフィックスで始まることを確認
				if !bytes.HasPrefix([]byte(result), []byte(tc.expectedPrefix)) {
					t.Errorf("Expected result to start with %s, got %s", tc.expectedPrefix, result)
				}

				// ヘッダー情報が含まれていることを確認
				for key, value := range tc.response.Headers {
					headerStr := key + ": " + value
					if !bytes.Contains([]byte(result), []byte(headerStr)) {
						t.Errorf("Expected result to contain header %s, got %s", headerStr, result)
					}
				}

				// ボディが空でない場合、ボディ情報が含まれていることを確認
				if len(tc.response.Body) > 0 {
					if !bytes.Contains([]byte(result), []byte("Body:")) {
						t.Errorf("Expected result to contain body section, got %s", result)
					}
				}

				if len(tc.response.Warnings) > 0 {
					for _, warning := range tc.response.Warnings {
						if !bytes.Contains([]byte(result), []byte(warning)) {
							t.Errorf("Expected warning %s to be present, got %s", warning, result)
						}
					}
				}

				for _, fragment := range tc.notContains {
					if bytes.Contains([]byte(result), []byte(fragment)) {
						t.Errorf("Expected result not to contain %q, but it was present. result: %s", fragment, result)
					}
				}
			}
		})
	}
}

func TestSanitizeHTMLBody(t *testing.T) {
	input := `<html><body><!--comment--><header>brand</header><div data-testid="wrapper"><main class="content" style="color:red"><script>console.log('xss')</script><p>text</p><!--secret--><button>click</button><span data-allow-mismatch="true" data-allow-missmatch="true" data-testid="main-span">keep</span></main></div><footer>copyright</footer></body></html>`
	got := sanitizeHTMLBody(input)
	for _, fragment := range []string{"<script", "console.log", "class=", "style=", "<button", "<!--", "secret", "data-allow-missmatch", "data-allow-mismatch", "data-testid", "<header", "brand", "<footer", "copyright"} {
		if strings.Contains(got, fragment) {
			t.Fatalf("sanitizeHTMLBody should remove %q but result was %s", fragment, got)
		}
	}
}

// TestSendRequestWithJSONFileAndHeaders はSendRequestWithJSONFileAndHeadersメソッドのテストです
func TestSendRequestWithJSONFileAndHeaders(t *testing.T) {
	// テストケース
	testCases := []struct {
		name           string
		url            string
		method         string
		jsonFilePath   string
		headers        map[string]string
		expectedHeader string
		expectedValue  string
	}{
		{
			name:           "認証トークンを含むリクエスト",
			url:            "https://api.example.com",
			method:         "POST",
			jsonFilePath:   "test.json",
			headers:        map[string]string{"Authorization": "Bearer test-token"},
			expectedHeader: "Authorization",
			expectedValue:  "Bearer test-token",
		},
		{
			name:           "認証トークンなしのリクエスト",
			url:            "https://api.example.com",
			method:         "GET",
			jsonFilePath:   "test.json",
			headers:        map[string]string{},
			expectedHeader: "Content-Type",
			expectedValue:  "application/json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックリポジトリの設定
			mockRepo := &MockHTTPRepository{
				LoadJSONFileFunc: func(filePath string) ([]byte, error) {
					return []byte(`{"test": "data"}`), nil
				},
				SendRequestFunc: func(request *models.HTTPRequest) (*models.HTTPResponse, error) {
					// ヘッダーが正しく設定されているか確認
					if value, exists := request.Headers[tc.expectedHeader]; !exists || value != tc.expectedValue {
						t.Errorf("Expected header %s with value %s, got %s", tc.expectedHeader, tc.expectedValue, value)
					}

					// URLとメソッドが正しく設定されているか確認
					if request.URL != tc.url {
						t.Errorf("Expected URL %s, got %s", tc.url, request.URL)
					}
					if request.Method != tc.method {
						t.Errorf("Expected method %s, got %s", tc.method, request.Method)
					}

					// モックレスポンスを返す
					return &models.HTTPResponse{
						StatusCode: 200,
						Headers:    map[string]string{"Content-Type": "application/json"},
						Body:       []byte(`{"success": true}`),
					}, nil
				},
			}

			// テスト対象のサービスを作成
			service := NewHTTPService(mockRepo)

			// メソッドを実行
			response, err := service.SendRequestWithJSONFileAndHeaders(tc.url, tc.method, tc.jsonFilePath, tc.headers, "auto")

			// エラーがないことを確認
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			// レスポンスが期待通りであることを確認
			expectedResponse := &models.HTTPResponse{
				StatusCode: 200,
				Headers:    map[string]string{"Content-Type": "application/json"},
				Body:       []byte(`{"success": true}`),
			}

			if !reflect.DeepEqual(response, expectedResponse) {
				t.Errorf("Expected response %+v, got %+v", expectedResponse, response)
			}
		})
	}
}
