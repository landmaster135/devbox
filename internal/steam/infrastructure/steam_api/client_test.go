package steam_api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockHTTPClient はHTTPクライアントのモック
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do はHTTPリクエストを実行します
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		headers  map[string]string
		expected *Client
	}{
		{
			name:   "正常なクライアント作成",
			apiKey: "test-api-key",
			headers: map[string]string{
				"User-Agent": "test-agent",
			},
			expected: &Client{
				apiKey: "test-api-key",
				headers: map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
					"User-Agent":   "test-agent",
				},
			},
		},
		{
			name:    "ヘッダーなしでのクライアント作成",
			apiKey:  "test-api-key",
			headers: nil,
			expected: &Client{
				apiKey: "test-api-key",
				headers: map[string]string{
					"Content-Type": "application/json",
					"Accept":       "application/json",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, tt.headers)

			assert.Equal(t, tt.expected.apiKey, client.apiKey)
			assert.Equal(t, tt.expected.headers, client.headers)
			assert.NotNil(t, client.httpClient)
		})
	}
}

func TestClient_Request_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// モックレスポンスを設定
			responseBody := `{"response": {"success": true}}`
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}
			return resp, nil
		},
	}

	client := &Client{
		httpClient: mockHTTPClient,
		apiKey:     "test-api-key",
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	ctx := context.Background()
	params := map[string]any{
		"steamid": "76561198000000000",
	}

	result, err := client.Request(ctx, "GET", "/test/endpoint", params)

	assert.NoError(t, err)
	assert.NotNil(t, result)

	// レスポンスがJSONとしてパースされていることを確認
	resultMap, ok := result.(map[string]any)
	assert.True(t, ok)
	assert.Contains(t, resultMap, "response")
}

func TestClient_Request_HTTPError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// 404エラーレスポンスを設定
			responseBody := `{"error": "Not Found"}`
			resp := &http.Response{
				StatusCode: 404,
				Status:     "404 Not Found",
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}
			return resp, nil
		},
	}

	client := &Client{
		httpClient: mockHTTPClient,
		apiKey:     "test-api-key",
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	ctx := context.Background()
	params := map[string]any{
		"steamid": "76561198000000000",
	}

	result, err := client.Request(ctx, "GET", "/test/endpoint", params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "HTTP error: 404")
}

func TestClient_Request_SteamAPIError(t *testing.T) {
	callCount := 0
	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			callCount++
			// Steam APIエラーレスポンスを設定
			responseBody := `{"code": 403, "description": "Access Denied"}`
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}
			return resp, nil
		},
	}

	client := &Client{
		httpClient: mockHTTPClient,
		apiKey:     "test-api-key",
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	ctx := context.Background()
	params := map[string]any{
		"steamid": "76561198000000000",
	}

	result, err := client.Request(ctx, "GET", "/test/endpoint", params)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 1, callCount, "Should not retry on Steam API error")

	// 安全な型アサーション
	if steamErr, ok := err.(*SteamError); ok {
		assert.Equal(t, 403, steamErr.Code)
		assert.Equal(t, "Access Denied", steamErr.Description)
	} else {
		t.Errorf("Expected *SteamError, got %T: %v", err, err)
	}
}

func TestClient_RequestWithoutKey_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// モックレスポンスを設定
			responseBody := `{"success": true, "data": {"name": "Test Game"}}`
			resp := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(responseBody)),
			}
			return resp, nil
		},
	}

	client := &Client{
		httpClient: mockHTTPClient,
		apiKey:     "test-api-key",
		headers: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	ctx := context.Background()
	params := map[string]any{
		"appids": 440,
		"cc":     "US",
	}

	result, err := client.RequestWithoutKey(ctx, "GET", "https://store.steampowered.com/api/appdetails", params)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestNewSteamClient(t *testing.T) {
	apiKey := "test-api-key"
	headers := map[string]string{
		"User-Agent": "test-agent",
	}

	steamClient := NewSteamClient(apiKey, headers)

	assert.NotNil(t, steamClient)
	assert.NotNil(t, steamClient.Users)
	assert.NotNil(t, steamClient.Apps)
}
