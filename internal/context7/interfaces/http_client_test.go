package interfaces

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockHTTPClient はHTTPクライアントのモックです
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do はHTTPリクエストを実行します
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// TestHTTPClient_Interface はHTTPClientインターフェースの実装をテストします
func TestHTTPClient_Interface(t *testing.T) {
	// DefaultHTTPClientがHTTPClientインターフェースを実装していることを確認
	var _ HTTPClient = &DefaultHTTPClient{}

	// MockHTTPClientがHTTPClientインターフェースを実装していることを確認
	var _ HTTPClient = &MockHTTPClient{}
}

// TestNewDefaultHTTPClient は新しいDefaultHTTPClientインスタンスの作成をテストします
func TestNewDefaultHTTPClient(t *testing.T) {
	client := NewDefaultHTTPClient()

	if client == nil {
		t.Fatal("NewDefaultHTTPClient should return a non-nil client")
	}

	if client.client == nil {
		t.Error("DefaultHTTPClient.client should not be nil")
	}
}

// TestDefaultHTTPClient_Do は DefaultHTTPClient の Do メソッドをテストします
func TestDefaultHTTPClient_Do(t *testing.T) {
	client := NewDefaultHTTPClient()

	// テスト用のHTTPリクエストを作成
	req, err := http.NewRequest("GET", "https://httpbin.org/get", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// リクエストを実行（実際のHTTPリクエストは行わず、構造をテスト）
	// 実際のネットワーク呼び出しを避けるため、モックを使用
	originalClient := client.client

	// テスト用のモックレスポンスを作成
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"test": "response"}`)),
		Header:     make(http.Header),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	// モッククライアントを設定
	client.client = &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				return mockResponse, nil
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("DefaultHTTPClient.Do should not return error: %v", err)
	}

	if resp == nil {
		t.Fatal("Response should not be nil")
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", resp.Header.Get("Content-Type"))
	}

	// レスポンスボディを読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	resp.Body.Close()

	expectedBody := `{"test": "response"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, string(body))
	}

	// 元のクライアントを復元
	client.client = originalClient
}

// TestDefaultHTTPClient_Do_WithPostRequest は POST リクエストをテストします
func TestDefaultHTTPClient_Do_WithPostRequest(t *testing.T) {
	client := NewDefaultHTTPClient()

	// POST リクエスト用のボディを作成
	requestBody := `{"query": "test"}`
	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewReader([]byte(requestBody)))
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// テスト用のモックレスポンスを作成
	mockResponse := &http.Response{
		StatusCode: 201,
		Body:       io.NopCloser(strings.NewReader(`{"created": true}`)),
		Header:     make(http.Header),
	}
	mockResponse.Header.Set("Content-Type", "application/json")

	// モッククライアントを設定
	client.client = &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				// リクエストの検証
				if req.Method != "POST" {
					t.Errorf("Expected method POST, got %s", req.Method)
				}

				if req.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", req.Header.Get("Content-Type"))
				}

				// リクエストボディの検証
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Errorf("Failed to read request body: %v", err)
				}

				if string(body) != requestBody {
					t.Errorf("Expected request body '%s', got '%s'", requestBody, string(body))
				}

				return mockResponse, nil
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("DefaultHTTPClient.Do should not return error: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("Expected status code 201, got %d", resp.StatusCode)
	}
}

// TestDefaultHTTPClient_Do_WithHeaders はヘッダー付きリクエストをテストします
func TestDefaultHTTPClient_Do_WithHeaders(t *testing.T) {
	client := NewDefaultHTTPClient()

	req, err := http.NewRequest("GET", "https://httpbin.org/headers", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// カスタムヘッダーを設定
	req.Header.Set("X-Custom-Header", "test-value")
	req.Header.Set("Authorization", "Bearer token123")

	// テスト用のモックレスポンスを作成
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"headers": {"X-Custom-Header": "test-value"}}`)),
		Header:     make(http.Header),
	}

	// モッククライアントを設定
	client.client = &http.Client{
		Transport: &MockRoundTripper{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				// ヘッダーの検証
				if req.Header.Get("X-Custom-Header") != "test-value" {
					t.Errorf("Expected X-Custom-Header 'test-value', got '%s'", req.Header.Get("X-Custom-Header"))
				}

				if req.Header.Get("Authorization") != "Bearer token123" {
					t.Errorf("Expected Authorization 'Bearer token123', got '%s'", req.Header.Get("Authorization"))
				}

				return mockResponse, nil
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("DefaultHTTPClient.Do should not return error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

// TestDefaultHTTPClient_MultipleInstances は複数のインスタンスが独立していることをテストします
func TestDefaultHTTPClient_MultipleInstances(t *testing.T) {
	client1 := NewDefaultHTTPClient()
	client2 := NewDefaultHTTPClient()

	if client1 == client2 {
		t.Error("NewDefaultHTTPClient should return different instances")
	}

	if client1.client == client2.client {
		t.Error("DefaultHTTPClient instances should have different http.Client instances")
	}
}

// MockRoundTripper はhttp.RoundTripperのモック実装です
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

// RoundTrip はHTTPリクエストを処理します
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

// TestMockHTTPClient_Do はMockHTTPClientの動作をテストします
func TestMockHTTPClient_Do(t *testing.T) {
	expectedResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("mock response")),
	}

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストの検証
			if req.Method != "GET" {
				t.Errorf("Expected method GET, got %s", req.Method)
			}

			if req.URL.String() != "https://example.com/test" {
				t.Errorf("Expected URL 'https://example.com/test', got '%s'", req.URL.String())
			}

			return expectedResponse, nil
		},
	}

	req, err := http.NewRequest("GET", "https://example.com/test", nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	resp, err := mockClient.Do(req)
	if err != nil {
		t.Fatalf("MockHTTPClient.Do should not return error: %v", err)
	}

	if resp != expectedResponse {
		t.Error("MockHTTPClient should return the expected response")
	}
}

// TestDefaultHTTPClient_NilRequest はnilリクエストの処理をテストします
func TestDefaultHTTPClient_NilRequest(t *testing.T) {
	client := NewDefaultHTTPClient()

	// nilリクエストを渡すとpanicが発生することを確認
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil request")
		}
	}()

	_, _ = client.Do(nil)
}
