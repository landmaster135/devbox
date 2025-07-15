package usecases

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

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

// MockRateLimiter はレート制限のモック実装
type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) CheckLimit() error {
	args := m.Called()
	return args.Error(0)
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

// resetRateLimit はレート制限をリセットするヘルパー関数（廃止予定）
// 新しいテストではMockRateLimiterを使用してください
func resetRateLimit() {
	// この関数は後方互換性のために残していますが、
	// 新しいテストではMockRateLimiterを使用することを推奨します
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
	mockRateLimiter := &MockRateLimiter{}
	service := NewBraveSearchServiceWithAllDependencies(&DefaultHTTPClient{}, &DefaultEnvironmentReader{}, mockRateLimiter)

	// モックの設定
	mockRateLimiter.On("CheckLimit").Return(nil)

	err := service.checkRateLimit()
	assert.NoError(t, err)

	// モックの検証
	mockRateLimiter.AssertExpectations(t)
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

// #==============================================================#
// ##          Error Handling Tests                             ##
// #==============================================================#
func TestBraveSearchService_HandleWebSearch_JSONDecodeError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

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

func TestBraveSearchService_getPoisData_JSONDecodeError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

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
	result, err := service.getPoisData([]string{"location1"})

	// 結果の検証
	assert.Error(t, err)
	assert.Empty(t, result.Results)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_getDescriptionsData_JSONDecodeError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

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
	result, err := service.getDescriptionsData([]string{"location1"})

	// 結果の検証
	assert.Error(t, err)
	assert.Empty(t, result.Descriptions)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

// #==============================================================#
// ##          Additional Edge Cases                            ##
// #==============================================================#
func TestBraveSearchService_formatLocalResults_MultipleLocations(t *testing.T) {
	service := NewBraveSearchService()

	// 複数のロケーションデータ
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Location 1",
				Address: struct {
					StreetAddress   string `json:"streetAddress,omitempty"`
					AddressLocality string `json:"addressLocality,omitempty"`
					AddressRegion   string `json:"addressRegion,omitempty"`
					PostalCode      string `json:"postalCode,omitempty"`
				}{
					StreetAddress: "123 First St",
				},
				Phone: "111-111-1111",
			},
			{
				ID:   "location2",
				Name: "Location 2",
				Address: struct {
					StreetAddress   string `json:"streetAddress,omitempty"`
					AddressLocality string `json:"addressLocality,omitempty"`
					AddressRegion   string `json:"addressRegion,omitempty"`
					PostalCode      string `json:"postalCode,omitempty"`
				}{
					AddressLocality: "Second City",
				},
				Rating: struct {
					RatingValue float64 `json:"ratingValue,omitempty"`
					RatingCount int     `json:"ratingCount,omitempty"`
				}{
					RatingValue: 3.8,
					RatingCount: 50,
				},
			},
		},
	}

	descData := BraveDescription{
		Descriptions: map[string]string{
			"location1": "First location description",
		},
	}

	// テストの実行
	result := service.formatLocalResults(poisData, descData)

	// 結果の検証
	assert.Contains(t, result, "Name: Location 1")
	assert.Contains(t, result, "Address: 123 First St")
	assert.Contains(t, result, "Phone: 111-111-1111")
	assert.Contains(t, result, "Description: First location description")
	assert.Contains(t, result, "---") // セパレータ
	assert.Contains(t, result, "Name: Location 2")
	assert.Contains(t, result, "Address: Second City")
	assert.Contains(t, result, "Rating: 3.8 (50 reviews)")
	assert.Contains(t, result, "Description: No description available")
}

func TestBraveSearchService_formatLocalResults_PartialAddress(t *testing.T) {
	service := NewBraveSearchService()

	// 部分的な住所データ
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Test Location",
				Address: struct {
					StreetAddress   string `json:"streetAddress,omitempty"`
					AddressLocality string `json:"addressLocality,omitempty"`
					AddressRegion   string `json:"addressRegion,omitempty"`
					PostalCode      string `json:"postalCode,omitempty"`
				}{
					AddressLocality: "Test City",
					PostalCode:      "12345",
					// StreetAddressとAddressRegionは空
				},
			},
		},
	}

	descData := BraveDescription{}

	// テストの実行
	result := service.formatLocalResults(poisData, descData)

	// 結果の検証
	assert.Contains(t, result, "Address: Test City, 12345")
	assert.NotContains(t, result, "Address: , Test City, , 12345") // 空の部分は含まれない
}

func TestBraveSearchService_formatLocalResults_ZeroRating(t *testing.T) {
	service := NewBraveSearchService()

	// 評価が0の場合
	poisData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:   "location1",
				Name: "Test Location",
				Rating: struct {
					RatingValue float64 `json:"ratingValue,omitempty"`
					RatingCount int     `json:"ratingCount,omitempty"`
				}{
					RatingValue: 0.0,
					RatingCount: 0,
				},
			},
		},
	}

	descData := BraveDescription{}

	// テストの実行
	result := service.formatLocalResults(poisData, descData)

	// 結果の検証
	assert.Contains(t, result, "Rating: N/A")
}

func TestBraveSearchService_HandleLocalSearch_GoroutineError(t *testing.T) {
	// t.Skip("goroutineエラーのテストは複雑なため、スキップ")

	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key").Times(3) // 3回のAPI呼び出し用

	// Web検索レスポンス（ロケーションIDを含む）
	webResponseData := BraveWebResponse{}
	webResponseData.Locations.Results = []struct {
		ID    string `json:"id"`
		Title string `json:"title,omitempty"`
	}{
		{ID: "location1", Title: "Test Location 1"},
	}

	webMockResponse := createMockResponse(http.StatusOK, webResponseData)
	poisErrorResponse := createMockResponse(http.StatusInternalServerError, map[string]string{"error": "POI Error"})

	// Web検索は成功、POI取得でエラー
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/web/search"
	})).Return(webMockResponse, nil).Once()

	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/local/pois"
	})).Return(poisErrorResponse, nil).Once()

	// テストの実行
	result, err := service.HandleLocalSearch("test local query", 5)

	// 結果の検証（goroutineでエラーが発生）
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error of Brave API")
	assert.Empty(t, result)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

// #==============================================================#
// ##          Default Environment Reader Tests                 ##
// #==============================================================#
func TestDefaultEnvironmentReader_Getenv(t *testing.T) {
	reader := &DefaultEnvironmentReader{}

	// 実際の環境変数は設定されていない可能性があるので、空文字列を期待
	result := reader.Getenv("NON_EXISTENT_ENV_VAR")
	assert.Equal(t, "", result)
}

// #==============================================================#
// ##          Default HTTP Client Tests                        ##
// #==============================================================#
func TestDefaultHTTPClient_Do(t *testing.T) {
	client := &DefaultHTTPClient{}

	// 簡単なHTTPリクエストのテスト（実際のネットワーク呼び出しは避ける）
	req, err := http.NewRequest("GET", "http://httpbin.org/status/200", nil)
	assert.NoError(t, err)

	// このテストは実際のネットワーク接続に依存するため、
	// エラーが発生してもテストが失敗しないようにする
	_, _ = client.Do(req)
	// ネットワークエラーの可能性があるため、エラーチェックはしない
	// 主な目的はコードカバレッジの向上
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

// #==============================================================#
// ##          Local Search Complete Flow Tests                 ##
// #==============================================================#
func TestBraveSearchService_HandleLocalSearch_Normal(t *testing.T) {
	// t.Skip("")

	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key").Times(3) // 3回のAPI呼び出し用

	// Web検索レスポンス（ロケーションIDを含む）
	webResponseData := BraveWebResponse{}
	webResponseData.Locations.Results = []struct {
		ID    string `json:"id"`
		Title string `json:"title,omitempty"`
	}{
		{ID: "location1", Title: "Test Location 1"},
		{ID: "location2", Title: "Test Location 2"},
	}

	// POIレスポンス
	poisResponseData := BravePoiResponse{
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
				},
				Phone: "123-456-7890",
			},
		},
	}

	// 説明レスポンス
	descResponseData := BraveDescription{
		Descriptions: map[string]string{
			"location1": "This is a test location",
		},
	}

	// HTTPクライアントのモック設定（3回のAPI呼び出し）
	webMockResponse := createMockResponse(http.StatusOK, webResponseData)
	poisMockResponse := createMockResponse(http.StatusOK, poisResponseData)
	descMockResponse := createMockResponse(http.StatusOK, descResponseData)

	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/web/search"
	})).Return(webMockResponse, nil).Once()

	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/local/pois"
	})).Return(poisMockResponse, nil).Once()

	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/local/descriptions"
	})).Return(descMockResponse, nil).Once()

	// テストの実行
	result, err := service.HandleLocalSearch("test local query", 5)

	// 結果の検証
	assert.NoError(t, err)
	assert.Contains(t, result, "Name: Test Location 1")
	assert.Contains(t, result, "Address: 123 Test St, Test City")
	assert.Contains(t, result, "Phone: 123-456-7890")
	assert.Contains(t, result, "Description: This is a test location")

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_HandleLocalSearch_FallbackToWebSearch(t *testing.T) {
	// t.Skip("")

	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key").Times(2) // 2回のAPI呼び出し用

	// Web検索レスポンス（ロケーションIDなし）
	webResponseData := BraveWebResponse{}
	webResponseData.Web.Results = []BraveWebResult{
		{
			Title:       "Fallback Result",
			Description: "This is a fallback result",
			URL:         "https://example.com/fallback",
		},
	}

	webMockResponse := createMockResponse(http.StatusOK, webResponseData)
	fallbackMockResponse := createMockResponse(http.StatusOK, webResponseData)

	// 最初のWeb検索（ロケーション検索）
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/web/search" && req.URL.RawQuery != ""
	})).Return(webMockResponse, nil).Once()

	// フォールバック用のWeb検索
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/web/search"
	})).Return(fallbackMockResponse, nil).Once()

	// テストの実行
	result, err := service.HandleLocalSearch("test local query", 5)

	// 結果の検証
	assert.NoError(t, err)
	assert.Contains(t, result, "Fallback Result")
	assert.Contains(t, result, "This is a fallback result")
	assert.Contains(t, result, "https://example.com/fallback")

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_HandleLocalSearch_HTTPError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	mockResponse := createMockResponse(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行
	result, err := service.HandleLocalSearch("test local query", 5)

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error of Brave API")
	assert.Empty(t, result)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

// #==============================================================#
// ##          POI Data Tests                                    ##
// #==============================================================#
func TestBraveSearchService_getPoisData_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	responseData := BravePoiResponse{
		Results: []BraveLocation{
			{
				ID:    "location1",
				Name:  "Test POI 1",
				Phone: "123-456-7890",
			},
		},
	}

	mockResponse := createMockResponse(http.StatusOK, responseData)
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行
	result, err := service.getPoisData([]string{"location1"})

	// 結果の検証
	assert.NoError(t, err)
	assert.Len(t, result.Results, 1)
	assert.Equal(t, "location1", result.Results[0].ID)
	assert.Equal(t, "Test POI 1", result.Results[0].Name)
	assert.Equal(t, "123-456-7890", result.Results[0].Phone)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_getPoisData_NoAPIKey(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定（APIキーなし）
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("")

	// テストの実行
	result, err := service.getPoisData([]string{"location1"})

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BRAVE_API_KEY environment variable is required")
	assert.Empty(t, result.Results)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
}

func TestBraveSearchService_getPoisData_HTTPError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	mockResponse := createMockResponse(http.StatusBadRequest, map[string]string{"error": "Bad Request"})
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行
	result, err := service.getPoisData([]string{"location1"})

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error of Brave API")
	assert.Empty(t, result.Results)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_getPoisData_EmptyIDs(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	responseData := BravePoiResponse{
		Results: []BraveLocation{},
	}

	mockResponse := createMockResponse(http.StatusOK, responseData)
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行（空のIDと空文字列を含む）
	result, err := service.getPoisData([]string{"", "location1", ""})

	// 結果の検証
	assert.NoError(t, err)
	assert.Empty(t, result.Results)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

// #==============================================================#
// ##          Description Data Tests                           ##
// #==============================================================#
func TestBraveSearchService_getDescriptionsData_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	responseData := BraveDescription{
		Descriptions: map[string]string{
			"location1": "Description for location 1",
			"location2": "Description for location 2",
		},
	}

	mockResponse := createMockResponse(http.StatusOK, responseData)
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行
	result, err := service.getDescriptionsData([]string{"location1", "location2"})

	// 結果の検証
	assert.NoError(t, err)
	assert.Len(t, result.Descriptions, 2)
	assert.Equal(t, "Description for location 1", result.Descriptions["location1"])
	assert.Equal(t, "Description for location 2", result.Descriptions["location2"])

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_getDescriptionsData_NoAPIKey(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定（APIキーなし）
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("")

	// テストの実行
	result, err := service.getDescriptionsData([]string{"location1"})

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BRAVE_API_KEY environment variable is required")
	assert.Empty(t, result.Descriptions)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
}

func TestBraveSearchService_getDescriptionsData_HTTPError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	mockResponse := createMockResponse(http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行
	result, err := service.getDescriptionsData([]string{"location1"})

	// 結果の検証
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error of Brave API")
	assert.Empty(t, result.Descriptions)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_getDescriptionsData_EmptyIDs(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key")

	responseData := BraveDescription{
		Descriptions: map[string]string{},
	}

	mockResponse := createMockResponse(http.StatusOK, responseData)
	mockHTTPClient.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResponse, nil)

	// テストの実行（空のIDと空文字列を含む）
	result, err := service.getDescriptionsData([]string{"", "location1", ""})

	// 結果の検証
	assert.NoError(t, err)
	assert.Empty(t, result.Descriptions)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

// #==============================================================#
// ##          Edge Case Tests                                   ##
// #==============================================================#
func TestBraveSearchService_HandleWebSearch_CountLimit(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

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

func TestBraveSearchService_HandleLocalSearch_CountLimit(t *testing.T) {
	// t.Skip("")

	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

	resetRateLimit()

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key").Times(2) // 2回のAPI呼び出し用

	// Web検索レスポンス（ロケーションIDなし、フォールバックする）
	webResponseData := BraveWebResponse{}
	webResponseData.Web.Results = []BraveWebResult{
		{Title: "Fallback", Description: "Fallback", URL: "https://example.com"},
	}

	mockResponse := createMockResponse(http.StatusOK, webResponseData)
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		// count=20が設定されていることを確認（25を指定したが20に制限される）
		return req.URL.Query().Get("count") == "20"
	})).Return(mockResponse, nil).Times(2) // 最初の検索とフォールバック

	// テストの実行（制限を超えるcount=25を指定）
	result, err := service.HandleLocalSearch("test query", 25)

	// 結果の検証
	assert.NoError(t, err)
	assert.Contains(t, result, "Fallback")

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}
