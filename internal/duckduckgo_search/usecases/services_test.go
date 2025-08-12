package usecases

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// MockHTTPClient はHTTPクライアントのモック実装です
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestDuckDuckGoSearchService_cleanText はcleanTextメソッドのテストです
func TestDuckDuckGoSearchService_cleanText(t *testing.T) {
	service := NewDuckDuckGoSearchService()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTMLタグの除去",
			input:    "<p>Hello <b>World</b></p>",
			expected: "Hello World",
		},
		{
			name:     "余分な空白の除去",
			input:    "  Hello   World  ",
			expected: "Hello World",
		},
		{
			name:     "単体のエンティティテスト - &amp;",
			input:    "&amp;",
			expected: "&",
		},
		{
			name:     "単体のエンティティテスト - &lt;",
			input:    "&lt;",
			expected: "<",
		},
		{
			name:     "単体のエンティティテスト - &gt;",
			input:    "&gt;",
			expected: ">",
		},
		{
			name:     "単体のエンティティテスト - &quot;",
			input:    "&quot;",
			expected: "\"",
		},
		{
			name:     "単体のエンティティテスト - &#39;",
			input:    "&#39;",
			expected: "'",
		},
		{
			name:     "組み合わせテスト - エンティティのみ",
			input:    "Hello &amp; world &quot;test&quot; &#39;example&#39;",
			expected: "Hello & world \"test\" 'example'",
		},
		{
			name:     "組み合わせテスト - HTMLタグとエンティティ",
			input:    "<p>Hello &amp; <b>world</b> &quot;test&quot;</p>",
			expected: "Hello & world \"test\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.cleanText(tt.input)
			if result != tt.expected {
				t.Errorf("cleanText(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestDuckDuckGoSearchService_decodeURL はdecodeURLメソッドのテストです
func TestDuckDuckGoSearchService_decodeURL(t *testing.T) {
	service := NewDuckDuckGoSearchService()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "通常のURL",
			input:    "https://example.com",
			expected: "https://example.com",
		},
		{
			name:     "DuckDuckGoリダイレクトURL",
			input:    "/l/?uddg=https%3A//example.com",
			expected: "https://example.com",
		},
		{
			name:     "エンコードされたURL",
			input:    "https%3A//example.com",
			expected: "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := service.decodeURL(tt.input)
			if result != tt.expected {
				t.Errorf("decodeURL() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestRateLimiter_CheckRateLimit はRateLimiterのテストです
func TestRateLimiter_CheckRateLimit(t *testing.T) {
	rateLimiter := NewRateLimiter(1, 30)

	// 最初のリクエストは成功するはず
	err := rateLimiter.CheckRateLimit()
	if err != nil {
		t.Errorf("First request should succeed, got error: %v", err)
	}

	// レート制限に達したときはエラーが返されるはず
	err = rateLimiter.CheckRateLimit()
	if err == nil {
		t.Error("Rate limit should be exceeded")
	}
}

// TestDuckDuckGoSearchService_HandleWebSearch_Normal はHandleWebSearchメソッドの正常系テストです
func TestDuckDuckGoSearchService_HandleWebSearch_Normal(t *testing.T) {
	// モックHTTPクライアントを作成
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// モックレスポンスを作成
			mockHTML := `
				<html>
					<body>
						<a class="result__a" href="https://example.com">Example Title</a>
						<a class="result__snippet">Example description</a>
					</body>
				</html>
			`
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockHTML)),
			}
			return resp, nil
		},
	}

	// レート制限を緩く設定
	rateLimiter := NewRateLimiter(10, 100)
	service := NewDuckDuckGoSearchServiceWithDependencies(mockClient, rateLimiter)

	result, err := service.HandleWebSearch("test query", 1)
	if err != nil {
		t.Errorf("HandleWebSearch() error = %v", err)
	}

	if !strings.Contains(result, "Example Title") {
		t.Errorf("Result should contain 'Example Title', got: %s", result)
	}
}

// TestDuckDuckGoSearchService_HandleWebSearch_HTTPError はHandleWebSearchメソッドのHTTPエラーテストです
func TestDuckDuckGoSearchService_HandleWebSearch_HTTPError(t *testing.T) {
	// モックHTTPクライアントを作成（エラーを返す）
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("network error")
		},
	}

	// レート制限を緩く設定
	rateLimiter := NewRateLimiter(10, 100)
	service := NewDuckDuckGoSearchServiceWithDependencies(mockClient, rateLimiter)

	_, err := service.HandleWebSearch("test query", 1)
	if err == nil {
		t.Error("HandleWebSearch() should return error for network failure")
	}
}

// TestDuckDuckGoSearchService_HandleWebSearch_RateLimit はHandleWebSearchメソッドのレート制限テストです
func TestDuckDuckGoSearchService_HandleWebSearch_RateLimit(t *testing.T) {
	// モックHTTPクライアントを作成
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("<html></html>")),
			}
			return resp, nil
		},
	}

	// 厳しいレート制限を設定
	rateLimiter := NewRateLimiter(1, 1)
	service := NewDuckDuckGoSearchServiceWithDependencies(mockClient, rateLimiter)

	// 最初のリクエストは成功
	_, err := service.HandleWebSearch("test query", 1)
	if err != nil {
		t.Errorf("First request should succeed, got error: %v", err)
	}

	// 2回目のリクエストはレート制限でエラー
	_, err = service.HandleWebSearch("test query", 1)
	if err == nil {
		t.Error("Second request should fail due to rate limit")
	}
}

// TestDuckDuckGoSearchService_HandleInstantSearch_Normal はHandleInstantSearchメソッドの正常系テストです
func TestDuckDuckGoSearchService_HandleInstantSearch_Normal(t *testing.T) {
	// モックHTTPクライアントを作成
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// モックJSONレスポンスを作成
			mockJSON := `{
				"Abstract": "Test abstract",
				"AbstractSource": "Test source",
				"AbstractURL": "https://example.com"
			}`
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockJSON)),
			}
			return resp, nil
		},
	}

	// レート制限を緩く設定
	rateLimiter := NewRateLimiter(10, 100)
	service := NewDuckDuckGoSearchServiceWithDependencies(mockClient, rateLimiter)

	result, err := service.HandleInstantSearch("test query")
	if err != nil {
		t.Errorf("HandleInstantSearch() error = %v", err)
	}

	if !strings.Contains(result, "Test abstract") {
		t.Errorf("Result should contain 'Test abstract', got: %s", result)
	}
}

// TestDuckDuckGoSearchService_HandleInstantSearch_NoResults はHandleInstantSearchメソッドの結果なしテストです
func TestDuckDuckGoSearchService_HandleInstantSearch_NoResults(t *testing.T) {
	// モックHTTPクライアントを作成
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// 空のJSONレスポンスを作成
			mockJSON := `{}`
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader(mockJSON)),
			}
			return resp, nil
		},
	}

	// レート制限を緩く設定
	rateLimiter := NewRateLimiter(10, 100)
	service := NewDuckDuckGoSearchServiceWithDependencies(mockClient, rateLimiter)

	result, err := service.HandleInstantSearch("test query")
	if err != nil {
		t.Errorf("HandleInstantSearch() error = %v", err)
	}

	if result != "No instant answer found" {
		t.Errorf("Result should be 'No instant answer found', got: %s", result)
	}
}

// TestDuckDuckGoSearchService_parseSearchResults はparseSearchResultsメソッドのテストです
func TestDuckDuckGoSearchService_parseSearchResults(t *testing.T) {
	service := NewDuckDuckGoSearchService()

	tests := []struct {
		name        string
		html        string
		maxResults  int
		expectCount int
	}{
		{
			name:        "基本的なパターン",
			html:        `<a class="result__a" href="https://example.com">Example Title</a>`,
			maxResults:  10,
			expectCount: 1,
		},
		{
			name:        "結果なし",
			html:        `<html><body>No results</body></html>`,
			maxResults:  10,
			expectCount: 0,
		},
		{
			name: "複数の結果",
			html: `
				<a class="result__a" href="https://example1.com">Title 1</a>
				<a class="result__a" href="https://example2.com">Title 2</a>
			`,
			maxResults:  10,
			expectCount: 2,
		},
		{
			name: "最大結果数の制限",
			html: `
				<a class="result__a" href="https://example1.com">Title 1</a>
				<a class="result__a" href="https://example2.com">Title 2</a>
				<a class="result__a" href="https://example3.com">Title 3</a>
			`,
			maxResults:  2,
			expectCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.parseSearchResults(tt.html, tt.maxResults)
			if err != nil {
				t.Errorf("parseSearchResults() error = %v", err)
			}

			if len(results) != tt.expectCount {
				t.Errorf("parseSearchResults() returned %d results, want %d", len(results), tt.expectCount)
			}
		})
	}
}

// TestNewDuckDuckGoSearchService はNewDuckDuckGoSearchServiceのテストです
func TestNewDuckDuckGoSearchService(t *testing.T) {
	service := NewDuckDuckGoSearchService()
	if service == nil {
		t.Error("NewDuckDuckGoSearchService() should not return nil")
		return
	}

	if service.httpClient == nil {
		t.Error("httpClient should not be nil")
	}

	if service.rateLimiter == nil {
		t.Error("rateLimiter should not be nil")
	}
}

// TestNewDuckDuckGoSearchServiceWithDependencies はNewDuckDuckGoSearchServiceWithDependenciesのテストです
func TestNewDuckDuckGoSearchServiceWithDependencies(t *testing.T) {
	mockClient := &MockHTTPClient{}
	rateLimiter := NewRateLimiter(1, 30)

	service := NewDuckDuckGoSearchServiceWithDependencies(mockClient, rateLimiter)
	if service == nil {
		t.Error("NewDuckDuckGoSearchServiceWithDependencies() should not return nil")
		return
	}

	if service.httpClient != mockClient {
		t.Error("httpClient should be the provided mock client")
	}

	if service.rateLimiter != rateLimiter {
		t.Error("rateLimiter should be the provided rate limiter")
	}
}

// TestDefaultHTTPClient はDefaultHTTPClientのテストです
func TestDefaultHTTPClient(t *testing.T) {
	timeout := 10 * time.Second
	client := NewDefaultHTTPClient(timeout)

	if client == nil {
		t.Error("NewDefaultHTTPClient() should not return nil")
		return
	}

	if client.client.Timeout != timeout {
		t.Errorf("Client timeout should be %v, got %v", timeout, client.client.Timeout)
	}
}

// TestDefaultHTTPClient_Do はDefaultHTTPClientのDoメソッドのテストです
func TestDefaultHTTPClient_Do(t *testing.T) {
	client := NewDefaultHTTPClient(1 * time.Second)

	// 無効なURLでリクエストを作成（エラーが発生することを確認）
	req, err := http.NewRequest("GET", "invalid-url", nil)
	if err == nil {
		// 無効なURLでもリクエスト作成が成功する場合があるので、実際にDoを呼び出してエラーを確認
		_, err = client.Do(req)
		if err == nil {
			t.Error("Do() should return error for invalid URL")
		}
	}
}

// TestRateLimiter_CheckRateLimit_Reset はRateLimiterのリセット機能のテストです
func TestRateLimiter_CheckRateLimit_Reset(t *testing.T) {
	rateLimiter := NewRateLimiter(1, 30)

	// 最初のリクエストは成功
	err := rateLimiter.CheckRateLimit()
	if err != nil {
		t.Errorf("First request should succeed, got error: %v", err)
	}

	// 2回目のリクエストはレート制限でエラー
	err = rateLimiter.CheckRateLimit()
	if err == nil {
		t.Error("Second request should fail due to rate limit")
	}

	// 時間を進めるために少し待つ（実際のテストでは時間のモックを使用することが推奨）
	time.Sleep(1100 * time.Millisecond)

	// リセット後は再び成功するはず
	err = rateLimiter.CheckRateLimit()
	if err != nil {
		t.Errorf("Request after reset should succeed, got error: %v", err)
	}
}
