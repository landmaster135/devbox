package usecases

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#
// MockHTTPClient はHTTPクライアントのモック実装
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockEnvironmentReader は環境変数読み取りのモック実装
type MockEnvironmentReader struct {
	mock.Mock
}

func (m *MockEnvironmentReader) Getenv(key string) string {
	args := m.Called(key)
	return args.String(0)
}

// #==============================================================#
// ##          Helper Functions                                  ##
// #==============================================================#
// createMockResponse はモックHTTPレスポンスを作成するヘルパー関数
func createMockResponse(statusCode int, body interface{}) *http.Response {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			panic(err)
		}
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     make(http.Header),
	}
}

// resetRateLimit はレート制限をリセットするヘルパー関数
func resetRateLimit() {
	requestCount.mu.Lock()
	defer requestCount.mu.Unlock()
	requestCount.second = 0
	requestCount.month = 0
	requestCount.lastReset = time.Now()
}

// #==============================================================#
// ##          BraveSearchService Tests                          ##
// #==============================================================#
func TestNewBraveSearchService_Normal(t *testing.T) {
	service := NewBraveSearchService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.httpClient)
	assert.NotNil(t, service.envReader)
}

func TestNewBraveSearchServiceWithDependencies_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}

	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)
	assert.NotNil(t, service)
	assert.Equal(t, mockHTTPClient, service.httpClient)
	assert.Equal(t, mockEnvReader, service.envReader)
}

// #==============================================================#
// ##          Rate Limiting Tests                               ##
// #==============================================================#
func TestBraveSearchService_checkRateLimit_Normal(t *testing.T) {
	service := NewBraveSearchService()
	resetRateLimit()

	err := service.checkRateLimit()
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount.second)
	assert.Equal(t, 1, requestCount.month)
}

func TestBraveSearchService_checkRateLimit_SecondLimitExceeded(t *testing.T) {
	service := NewBraveSearchService()
	resetRateLimit()

	// 秒間制限を超える
	requestCount.second = rateLimit.perSecond

	err := service.checkRateLimit()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
}

func TestBraveSearchService_checkRateLimit_MonthLimitExceeded(t *testing.T) {
	service := NewBraveSearchService()
	resetRateLimit()

	// 月間制限を超える
	requestCount.month = rateLimit.perMonth

	err := service.checkRateLimit()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
}

func TestBraveSearchService_checkRateLimit_SecondReset(t *testing.T) {
	service := NewBraveSearchService()
	resetRateLimit()

	// 秒間制限をリセット
	requestCount.lastReset = time.Now().Add(-2 * time.Second)
	requestCount.second = rateLimit.perSecond

	err := service.checkRateLimit()
	assert.NoError(t, err)
	assert.Equal(t, 1, requestCount.second)
}

// #==============================================================#
// ##          Web Search Tests                                  ##
// #==============================================================#
func TestBraveSearchService_HandleWebSearch_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

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

	resetRateLimit()

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

	resetRateLimit()

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

func TestBraveSearchService_HandleWebSearch_RateLimitExceeded(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// レート制限を超える
	requestCount.second = rateLimit.perSecond

	// テストの実行
	result, err := service.HandleWebSearch("test query", 5, 0)

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	assert.Empty(t, result)
}

// #==============================================================#
// ##          Local Search Tests                                ##
// #==============================================================#
func TestBraveSearchService_HandleLocalSearch_NoAPIKey(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定（APIキーなし）
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("")

	// テストの実行
	result, err := service.HandleLocalSearch("test local query", 5)

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BRAVE_API_KEY environment variable is required")
	assert.Empty(t, result)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
}

func TestBraveSearchService_HandleLocalSearch_RateLimitExceeded(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// レート制限を超える
	requestCount.second = rateLimit.perSecond

	// テストの実行
	result, err := service.HandleLocalSearch("test local query", 5)

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	assert.Empty(t, result)
}

// #==============================================================#
// ##          Utility Function Tests                            ##
// #==============================================================#
func TestBraveSearchService_getOrDefault_Normal(t *testing.T) {
	service := NewBraveSearchService()

	// 通常の値
	result := service.getOrDefault("test", "default")
	assert.Equal(t, "test", result)

	// 空の値
	result = service.getOrDefault("", "default")
	assert.Equal(t, "default", result)
}

func TestBraveSearchService_formatLocalResults_Normal(t *testing.T) {
	service := NewBraveSearchService()

	// テストデータの作成
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Test Location 1",
				Address: struct {
					StreetAddress   string `json:"streetAddress,omitempty"`
					AddressLocality string `json:"addressLocality,omitempty"`
					AddressRegion   string `json:"addressRegion,omitempty"`
					PostalCode      string `json:"postalCode,omitempty"`
				}{
					StreetAddress:   "123 Test St",
					AddressLocality: "Test City",
					AddressRegion:   "Test Region",
					PostalCode:      "12345",
				},
				Phone: "123-456-7890",
				Rating: struct {
					RatingValue float64 `json:"ratingValue,omitempty"`
					RatingCount int     `json:"ratingCount,omitempty"`
				}{
					RatingValue: 4.5,
					RatingCount: 100,
				},
				OpeningHours: []string{"Mon-Fri: 9AM-5PM", "Sat-Sun: Closed"},
				PriceRange:   "$$",
			},
		},
	}

	descData := BraveDescription{
		Descriptions: map[string]string{
			"location1": "This is a test location 1",
		},
	}

	// テストの実行
	result := service.formatLocalResults(poisData, descData)

	// 結果の検証
	assert.Contains(t, result, "Name: Test Location 1")
	assert.Contains(t, result, "Address: 123 Test St, Test City, Test Region, 12345")
	assert.Contains(t, result, "Phone: 123-456-7890")
	assert.Contains(t, result, "Rating: 4.5 (100 reviews)")
	assert.Contains(t, result, "Price Range: $$")
	assert.Contains(t, result, "Hours: Mon-Fri: 9AM-5PM, Sat-Sun: Closed")
	assert.Contains(t, result, "Description: This is a test location 1")
}

func TestBraveSearchService_formatLocalResults_EmptyResults(t *testing.T) {
	service := NewBraveSearchService()

	// 空のデータ
	emptyPoisData := BravePoiResponse{
		Results: []BraveLocation{},
	}
	descData := BraveDescription{}

	result := service.formatLocalResults(emptyPoisData, descData)
	assert.Equal(t, "No local results found", result)
}

func TestBraveSearchService_formatLocalResults_MissingData(t *testing.T) {
	service := NewBraveSearchService()

	// データが不完全な場合
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Test Location 1",
				// 他のフィールドは空
			},
		},
	}

	descData := BraveDescription{
		Descriptions: map[string]string{},
	}

	result := service.formatLocalResults(poisData, descData)

	// デフォルト値が使用されることを確認
	assert.Contains(t, result, "Name: Test Location 1")
	assert.Contains(t, result, "Address: N/A")
	assert.Contains(t, result, "Phone: N/A")
	assert.Contains(t, result, "Rating: N/A")
	assert.Contains(t, result, "Price Range: N/A")
	assert.Contains(t, result, "Hours: N/A")
	assert.Contains(t, result, "Description: No description available")
}
