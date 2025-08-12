package usecases

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/context7/domain/models"
	"github.com/landmaster135/devbox/internal/context7/interfaces"
)

// MockHTTPClient はHTTPクライアントのモックです
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do はHTTPリクエストを実行します
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestContext7Service_NewContext7Service は新しいContext7Serviceインスタンスの作成をテストします
func TestContext7Service_NewContext7Service(t *testing.T) {
	httpClient := interfaces.NewDefaultHTTPClient()
	service := NewContext7Service(httpClient)

	if service == nil {
		t.Fatal("NewContext7Service should return a non-nil service")
	}

	if service.httpClient != httpClient {
		t.Error("httpClient should be set correctly")
	}

	if service.baseURL != models.Context7APIBaseURL {
		t.Errorf("baseURL should be %s, got %s", models.Context7APIBaseURL, service.baseURL)
	}
}

// TestContext7Service_ResolveLibraryID_Normal は正常なライブラリ検索をテストします
func TestContext7Service_ResolveLibraryID_Normal(t *testing.T) {
	mockResponse := `{
		"results": [
			{
				"id": "/facebook/react",
				"title": "React",
				"description": "A JavaScript library for building user interfaces",
				"branch": "main",
				"lastUpdateDate": "2024-01-01",
				"state": "finalized",
				"totalTokens": 50000,
				"totalSnippets": 100,
				"totalPages": 50,
				"trustScore": 9
			}
		]
	}`

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストURLの検証
			expectedURL := "https://context7.com/api/v1/search?query=react"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	response, err := service.ResolveLibraryID("react")

	if err != nil {
		t.Fatalf("ResolveLibraryID should not return error: %v", err)
	}

	if response == nil {
		t.Fatal("Response should not be nil")
	}

	if len(response.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(response.Results))
	}

	result := response.Results[0]
	if result.ID != "/facebook/react" {
		t.Errorf("Expected ID '/facebook/react', got '%s'", result.ID)
	}

	if result.Title != "React" {
		t.Errorf("Expected title 'React', got '%s'", result.Title)
	}
}

// TestContext7Service_ResolveLibraryID_RateLimit はレート制限エラーをテストします
func TestContext7Service_ResolveLibraryID_RateLimit(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 429,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	response, err := service.ResolveLibraryID("react")

	if err != nil {
		t.Fatalf("ResolveLibraryID should not return error for rate limit: %v", err)
	}

	if response.Error == nil {
		t.Fatal("Response should contain error for rate limit")
	}

	expectedError := "レート制限に達しました。しばらく待ってから再試行してください。"
	if *response.Error != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, *response.Error)
	}
}

// TestContext7Service_GetLibraryDocs_Normal は正常なドキュメント取得をテストします
func TestContext7Service_GetLibraryDocs_Normal(t *testing.T) {
	mockResponse := "# React Documentation\n\nReact is a JavaScript library..."

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストURLの検証
			expectedURL := "https://context7.com/api/v1/facebook/react?tokens=10000&type=txt"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			// ヘッダーの検証
			if req.Header.Get("X-Context7-Source") != "go-cli-client" {
				t.Error("Expected X-Context7-Source header to be 'go-cli-client'")
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	options := models.DocOptions{
		Tokens: models.DefaultTokens,
	}

	docs, err := service.GetLibraryDocs("/facebook/react", options)

	if err != nil {
		t.Fatalf("GetLibraryDocs should not return error: %v", err)
	}

	if docs != mockResponse {
		t.Errorf("Expected docs '%s', got '%s'", mockResponse, docs)
	}
}

// TestContext7Service_GetLibraryDocs_WithOptions はオプション付きドキュメント取得をテストします
func TestContext7Service_GetLibraryDocs_WithOptions(t *testing.T) {
	mockResponse := "# React Hooks Documentation\n\nHooks are a new addition..."

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストURLの検証（トピックとトークン数を含む）
			expectedURL := "https://context7.com/api/v1/facebook/react?tokens=5000&topic=hooks&type=txt"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	options := models.DocOptions{
		Topic:  "hooks",
		Tokens: 5000,
	}

	docs, err := service.GetLibraryDocs("/facebook/react", options)

	if err != nil {
		t.Fatalf("GetLibraryDocs should not return error: %v", err)
	}

	if docs != mockResponse {
		t.Errorf("Expected docs '%s', got '%s'", mockResponse, docs)
	}
}

// TestContext7Service_GetLibraryDocs_NotFound は404エラーをテストします
func TestContext7Service_GetLibraryDocs_NotFound(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	options := models.DocOptions{}

	docs, err := service.GetLibraryDocs("/nonexistent/library", options)

	if err != nil {
		t.Fatalf("GetLibraryDocs should not return error for 404: %v", err)
	}

	expectedError := "指定されたライブラリIDのドキュメントが見つかりません。"
	if docs != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, docs)
	}
}

// TestContext7Service_FormatSearchResults_Normal は検索結果の整形をテストします
func TestContext7Service_FormatSearchResults_Normal(t *testing.T) {
	service := NewContext7Service(&MockHTTPClient{})

	trustScore := 9.0
	stars := 50000
	searchResponse := &models.SearchResponse{
		Results: []models.SearchResult{
			{
				ID:             "/facebook/react",
				Title:          "React",
				Description:    "A JavaScript library for building user interfaces",
				TotalSnippets:  100,
				TrustScore:     &trustScore,
				Stars:          &stars,
				Versions:       []string{"18.0.0", "17.0.0"},
				LastUpdateDate: "2024-01-01",
				State:          "finalized",
			},
		},
	}

	result := service.FormatSearchResults(searchResponse)

	// 結果に期待される内容が含まれているかチェック
	expectedContents := []string{
		"検索結果:",
		"1. React",
		"ID: /facebook/react",
		"説明: A JavaScript library for building user interfaces",
		"コードスニペット数: 100",
		"信頼スコア: 9.0",
		"スター数: 50000",
		"利用可能バージョン: 18.0.0, 17.0.0",
		"最終更新: 2024-01-01",
		"状態: finalized",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(result, expected) {
			t.Errorf("Result should contain '%s', but got: %s", expected, result)
		}
	}
}

// TestContext7Service_FormatSearchResults_Empty は空の検索結果をテストします
func TestContext7Service_FormatSearchResults_Empty(t *testing.T) {
	service := NewContext7Service(&MockHTTPClient{})

	searchResponse := &models.SearchResponse{
		Results: []models.SearchResult{},
	}

	result := service.FormatSearchResults(searchResponse)

	expected := "検索結果が見つかりませんでした。"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestContext7Service_FormatSearchResults_Error はエラーレスポンスをテストします
func TestContext7Service_FormatSearchResults_Error(t *testing.T) {
	service := NewContext7Service(&MockHTTPClient{})

	errorMsg := "API呼び出しが失敗しました"
	searchResponse := &models.SearchResponse{
		Results: []models.SearchResult{},
		Error:   &errorMsg,
	}

	result := service.FormatSearchResults(searchResponse)

	expected := "エラー: API呼び出しが失敗しました"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// TestContext7Service_ValidateLibraryID_Normal は正常なライブラリID検証をテストします
func TestContext7Service_ValidateLibraryID_Normal(t *testing.T) {
	service := NewContext7Service(&MockHTTPClient{})

	testCases := []string{
		"/facebook/react",
		"facebook/react",
		"/vercel/next.js",
		"mongodb/docs",
		"/supabase/supabase/v2.0.0",
	}

	for _, testCase := range testCases {
		err := service.ValidateLibraryID(testCase)
		if err != nil {
			t.Errorf("ValidateLibraryID should not return error for '%s': %v", testCase, err)
		}
	}
}

// TestContext7Service_ValidateLibraryID_Invalid は無効なライブラリID検証をテストします
func TestContext7Service_ValidateLibraryID_Invalid(t *testing.T) {
	service := NewContext7Service(&MockHTTPClient{})

	testCases := []struct {
		input    string
		expected string
	}{
		{"", "ライブラリIDが空です"},
		{"react", "ライブラリIDの形式が正しくありません。期待される形式: org/project または org/project/version"},
		{"/react", "ライブラリIDの形式が正しくありません。期待される形式: org/project または org/project/version"},
	}

	for _, testCase := range testCases {
		err := service.ValidateLibraryID(testCase.input)
		if err == nil {
			t.Errorf("ValidateLibraryID should return error for '%s'", testCase.input)
		}

		if err.Error() != testCase.expected {
			t.Errorf("Expected error '%s', got '%s'", testCase.expected, err.Error())
		}
	}
}

// TestContext7Service_NewContext7ServiceWithHTTPClient は NewContext7ServiceWithHTTPClient をテストします
func TestContext7Service_NewContext7ServiceWithHTTPClient(t *testing.T) {
	service := NewContext7ServiceWithHTTPClient()

	if service == nil {
		t.Fatal("NewContext7ServiceWithHTTPClient should return a non-nil service")
	}

	if service.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if service.baseURL != models.Context7APIBaseURL {
		t.Errorf("baseURL should be %s, got %s", models.Context7APIBaseURL, service.baseURL)
	}
}

// TestContext7Service_ResolveLibraryID_HTTPRequestError はHTTPリクエスト作成エラーをテストします
func TestContext7Service_ResolveLibraryID_HTTPRequestError(t *testing.T) {
	service := NewContext7Service(&MockHTTPClient{})

	// 無効なURLを設定してHTTPリクエスト作成エラーを発生させる
	service.baseURL = "://invalid-url"

	_, err := service.ResolveLibraryID("react")

	if err == nil {
		t.Fatal("ResolveLibraryID should return error for invalid URL")
	}

	if !strings.Contains(err.Error(), "検索URLの構築に失敗しました") {
		t.Errorf("Expected URL construction error, got: %v", err)
	}
}

// TestContext7Service_ResolveLibraryID_HTTPClientError はHTTPクライアントエラーをテストします
func TestContext7Service_ResolveLibraryID_HTTPClientError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	service := NewContext7Service(mockClient)
	_, err := service.ResolveLibraryID("react")

	if err == nil {
		t.Fatal("ResolveLibraryID should return error for HTTP client error")
	}

	if !strings.Contains(err.Error(), "HTTPリクエストの実行に失敗しました") {
		t.Errorf("Expected HTTP request error, got: %v", err)
	}
}

// TestContext7Service_ResolveLibraryID_ResponseBodyReadError はレスポンスボディ読み取りエラーをテストします
func TestContext7Service_ResolveLibraryID_ResponseBodyReadError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       &ErrorReader{},
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	_, err := service.ResolveLibraryID("react")

	if err == nil {
		t.Fatal("ResolveLibraryID should return error for response body read error")
	}

	if !strings.Contains(err.Error(), "レスポンスボディの読み取りに失敗しました") {
		t.Errorf("Expected response body read error, got: %v", err)
	}
}

// TestContext7Service_ResolveLibraryID_JSONParseError はJSONパースエラーをテストします
func TestContext7Service_ResolveLibraryID_JSONParseError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("invalid json")),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	_, err := service.ResolveLibraryID("react")

	if err == nil {
		t.Fatal("ResolveLibraryID should return error for JSON parse error")
	}

	if !strings.Contains(err.Error(), "JSONレスポンスのパースに失敗しました") {
		t.Errorf("Expected JSON parse error, got: %v", err)
	}
}

// TestContext7Service_ResolveLibraryID_ServerError はサーバーエラーをテストします
func TestContext7Service_ResolveLibraryID_ServerError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	response, err := service.ResolveLibraryID("react")

	if err != nil {
		t.Fatalf("ResolveLibraryID should not return error for server error: %v", err)
	}

	if response.Error == nil {
		t.Fatal("Response should contain error for server error")
	}

	expectedError := "API呼び出しが失敗しました (ステータス: 500)"
	if *response.Error != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, *response.Error)
	}
}

// TestContext7Service_GetLibraryDocs_HTTPRequestError はHTTPリクエスト作成エラーをテストします
func TestContext7Service_GetLibraryDocs_HTTPRequestError(t *testing.T) {
	service := NewContext7Service(&MockHTTPClient{})

	// 無効なURLを設定してHTTPリクエスト作成エラーを発生させる
	service.baseURL = "://invalid-url"

	options := models.DocOptions{}
	_, err := service.GetLibraryDocs("/facebook/react", options)

	if err == nil {
		t.Fatal("GetLibraryDocs should return error for invalid URL")
	}

	if !strings.Contains(err.Error(), "ドキュメントURLの構築に失敗しました") {
		t.Errorf("Expected URL construction error, got: %v", err)
	}
}

// TestContext7Service_GetLibraryDocs_HTTPClientError はHTTPクライアントエラーをテストします
func TestContext7Service_GetLibraryDocs_HTTPClientError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	service := NewContext7Service(mockClient)
	options := models.DocOptions{}
	_, err := service.GetLibraryDocs("/facebook/react", options)

	if err == nil {
		t.Fatal("GetLibraryDocs should return error for HTTP client error")
	}

	if !strings.Contains(err.Error(), "HTTPリクエストの実行に失敗しました") {
		t.Errorf("Expected HTTP request error, got: %v", err)
	}
}

// TestContext7Service_GetLibraryDocs_ResponseBodyReadError はレスポンスボディ読み取りエラーをテストします
func TestContext7Service_GetLibraryDocs_ResponseBodyReadError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       &ErrorReader{},
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	options := models.DocOptions{}
	_, err := service.GetLibraryDocs("/facebook/react", options)

	if err == nil {
		t.Fatal("GetLibraryDocs should return error for response body read error")
	}

	if !strings.Contains(err.Error(), "レスポンスボディの読み取りに失敗しました") {
		t.Errorf("Expected response body read error, got: %v", err)
	}
}

// TestContext7Service_GetLibraryDocs_RateLimit はレート制限エラーをテストします
func TestContext7Service_GetLibraryDocs_RateLimit(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 429,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	options := models.DocOptions{}
	docs, err := service.GetLibraryDocs("/facebook/react", options)

	if err != nil {
		t.Fatalf("GetLibraryDocs should not return error for rate limit: %v", err)
	}

	expectedError := "レート制限に達しました。しばらく待ってから再試行してください。"
	if docs != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, docs)
	}
}

// TestContext7Service_GetLibraryDocs_EmptyResponse は空のレスポンスをテストします
func TestContext7Service_GetLibraryDocs_EmptyResponse(t *testing.T) {
	testCases := []string{
		"",
		"No content available",
		"No context data available",
	}

	for _, responseText := range testCases {
		mockClient := &MockHTTPClient{
			DoFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(responseText)),
				}, nil
			},
		}

		service := NewContext7Service(mockClient)
		options := models.DocOptions{}
		docs, err := service.GetLibraryDocs("/facebook/react", options)

		if err != nil {
			t.Fatalf("GetLibraryDocs should not return error for empty response: %v", err)
		}

		expectedMessage := "このライブラリのドキュメントは利用できません。"
		if docs != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, docs)
		}
	}
}

// TestContext7Service_GetLibraryDocs_DefaultTokens はデフォルトトークン数をテストします
func TestContext7Service_GetLibraryDocs_DefaultTokens(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// クエリパラメータの検証
			expectedURL := "https://context7.com/api/v1/facebook/react?tokens=10000&type=txt"
			if req.URL.String() != expectedURL {
				t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
			}

			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("# React Documentation")),
			}, nil
		},
	}

	service := NewContext7Service(mockClient)
	options := models.DocOptions{} // Tokensを指定しない

	_, err := service.GetLibraryDocs("/facebook/react", options)
	if err != nil {
		t.Fatalf("GetLibraryDocs should not return error: %v", err)
	}
}

// TestContext7Service_FormatSearchResults_WithoutOptionalFields はオプションフィールドなしの検索結果をテストします
func TestContext7Service_FormatSearchResults_WithoutOptionalFields(t *testing.T) {
	service := NewContext7Service(&MockHTTPClient{})

	searchResponse := &models.SearchResponse{
		Results: []models.SearchResult{
			{
				ID:             "/test/library",
				Title:          "Test Library",
				Description:    "A test library",
				TotalSnippets:  50,
				LastUpdateDate: "2024-01-01",
				State:          "finalized",
				// TrustScore, Stars, Versions は nil
			},
		},
	}

	result := service.FormatSearchResults(searchResponse)

	// 必須フィールドが含まれていることを確認
	expectedContents := []string{
		"検索結果:",
		"1. Test Library",
		"ID: /test/library",
		"説明: A test library",
		"コードスニペット数: 50",
		"最終更新: 2024-01-01",
		"状態: finalized",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(result, expected) {
			t.Errorf("Result should contain '%s', but got: %s", expected, result)
		}
	}

	// オプションフィールドが含まれていないことを確認
	unexpectedContents := []string{
		"信頼スコア:",
		"スター数:",
		"利用可能バージョン:",
	}

	for _, unexpected := range unexpectedContents {
		if strings.Contains(result, unexpected) {
			t.Errorf("Result should not contain '%s', but got: %s", unexpected, result)
		}
	}
}

// ErrorReader はio.ReadAllでエラーを発生させるためのモック
type ErrorReader struct{}

func (e *ErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (e *ErrorReader) Close() error {
	return nil
}
