package usecases

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Web Search Tests                                  ##
// #==============================================================#
func TestBraveSearchService_HandleWebSearch_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

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
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

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

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_HandleWebSearch_NoAPIKey(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定（APIキーなし）
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("")

	// テストの実行
	result, err := service.HandleWebSearch("test query", 5, 0)

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BRAVE_API_KEY environment variable is required")
	assert.Empty(t, result)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
}

func TestBraveSearchService_HandleWebSearch_HTTPError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	mockResponse := createMockResponse(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行
	result, err := service.HandleWebSearch("test query", 5, 0)

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error of Brave API")
	assert.Empty(t, result)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_HandleWebSearch_JSONDecodeError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	// 不正なJSONレスポンス
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		Header:     make(http.Header),
	}
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行
	result, err := service.HandleWebSearch("test query", 5, 0)

	// 結果の検証
	assert.Error(t, err)
	assert.Empty(t, result)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_HandleWebSearch_CountLimit(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	responseData := BraveWebResponse{}
	responseData.Web.Results = []BraveWebResult{
		{Title: "Test", Description: "Test", URL: "https://example.com"},
	}

	mockResponse := createMockResponse(http.StatusOK, responseData)
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		// count=20が設定されていることを確認（25を指定したが20に制限される）
		return req.URL.Query().Get("count") == "20"
	})).Return(mockResponse, nil)

	// テストの実行（制限を超えるcount=25を指定）
	result, err := service.HandleWebSearch("test query", 25, 0)

	// 結果の検証
	assert.NoError(t, err)
	assert.Contains(t, result, "Test")

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

// #==============================================================#
// ##          Rate Limiting Tests                               ##
// #==============================================================#
func TestBraveSearchService_checkRateLimit_Normal(t *testing.T) {
	mockRateLimiter := &MockRateLimiter{}
	service := NewBraveSearchServiceWithAllDependencies(&DefaultHTTPClient{}, &DefaultEnvironmentReader{}, mockRateLimiter)

	// モックの設定
	mockRateLimiter.On("CheckLimit").Return(nil)

	err := service.checkRateLimit()
	assert.NoError(t, err)

	// モックの検証
	mockRateLimiter.AssertExpectations(t)
}
