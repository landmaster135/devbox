package usecases

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          Web Search Tests                                  ##
// #==============================================================#
func TestBraveSearchService_HandleWebSearch_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定
	mockEnvReader.GetenvFunc = func(key string) string {
		if key == "BRAVE_API_KEY" {
			return "test-api-key"
		}
		return ""
	}

	// レスポンスデータの作成
	responseData := BraveWebResponse{}
	responseData.Web.Results = []BraveWebResult{
		{
			Title:       "Test Result 1",
			Description: "This is a test result 1",
			URL:         "https://example.com/1",
		},
		{
			Title:       "Test Result 2",
			Description: "This is a test result 2",
			URL:         "https://example.com/2",
		},
	}

	mockResponse := createMockResponse(http.StatusOK, responseData)
	mockHTTPClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return mockResponse, nil
	}

	// テストの実行
	result, err := service.HandleWebSearch("test query", 5, 0)

	// 結果の検証
	assert.NoError(t, err)
	assert.Contains(t, result, "Test Result 1")
	assert.Contains(t, result, "This is a test result 1")
	assert.Contains(t, result, "https://example.com/1")
	assert.Contains(t, result, "Test Result 2")
	assert.Contains(t, result, "This is a test result 2")
	assert.Contains(t, result, "https://example.com/2")
}

func TestBraveSearchService_HandleWebSearch_NoAPIKey(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定（APIキーなし）
	mockEnvReader.GetenvFunc = func(key string) string {
		return ""
	}

	// テストの実行
	result, err := service.HandleWebSearch("test query", 5, 0)

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BRAVE_API_KEY environment variable is required")
	assert.Empty(t, result)
}

func TestBraveSearchService_HandleWebSearch_HTTPError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定
	mockEnvReader.GetenvFunc = func(key string) string {
		if key == "BRAVE_API_KEY" {
			return "test-api-key"
		}
		return ""
	}

	mockResponse := createMockResponse(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	mockHTTPClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return mockResponse, nil
	}

	// テストの実行
	result, err := service.HandleWebSearch("test query", 5, 0)

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error of Brave API")
	assert.Empty(t, result)
}

func TestBraveSearchService_HandleWebSearch_JSONDecodeError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定
	mockEnvReader.GetenvFunc = func(key string) string {
		if key == "BRAVE_API_KEY" {
			return "test-api-key"
		}
		return ""
	}

	// 不正なJSONレスポンス
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		Header:     make(http.Header),
	}
	mockHTTPClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return mockResponse, nil
	}

	// テストの実行
	result, err := service.HandleWebSearch("test query", 5, 0)

	// 結果の検証
	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestBraveSearchService_HandleWebSearch_CountLimit(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定
	mockEnvReader.GetenvFunc = func(key string) string {
		if key == "BRAVE_API_KEY" {
			return "test-api-key"
		}
		return ""
	}

	responseData := BraveWebResponse{}
	responseData.Web.Results = []BraveWebResult{
		{Title: "Test", Description: "Test", URL: "https://example.com"},
	}

	mockResponse := createMockResponse(http.StatusOK, responseData)
	mockHTTPClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		// count=20が設定されていることを確認（25を指定したが20に制限される）
		if req.URL.Query().Get("count") != "20" {
			t.Errorf("Expected count=20, but got count=%s", req.URL.Query().Get("count"))
		}
		return mockResponse, nil
	}

	// テストの実行（制限を超えるcount=25を指定）
	result, err := service.HandleWebSearch("test query", 25, 0)

	// 結果の検証
	assert.NoError(t, err)
	assert.Contains(t, result, "Test")
}

// #==============================================================#
// ##          Rate Limiting Tests                               ##
// #==============================================================#
func TestBraveSearchService_checkRateLimit_Normal(t *testing.T) {
	mockRateLimiter := &MockRateLimiter{}
	service := NewBraveSearchServiceWithAllDependencies(&DefaultHTTPClient{}, &DefaultEnvironmentReader{}, mockRateLimiter)

	// モックの設定
	mockRateLimiter.CheckLimitFunc = func() error {
		return nil
	}

	err := service.checkRateLimit()
	assert.NoError(t, err)
}
