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
// ##          Local Search Tests                                ##
// #==============================================================#
func TestBraveSearchService_HandleLocalSearch_Normal(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	mockRateLimiter := &MockRateLimiter{}
	service := NewBraveSearchServiceWithAllDependencies(mockHTTPClient, mockEnvReader, mockRateLimiter)

	// レート制限を無効化
	mockRateLimiter.On("CheckLimit").Return(nil)

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
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	mockRateLimiter := &MockRateLimiter{}
	service := NewBraveSearchServiceWithAllDependencies(mockHTTPClient, mockEnvReader, mockRateLimiter)

	// レート制限を無効化
	mockRateLimiter.On("CheckLimit").Return(nil)

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

func TestBraveSearchService_HandleLocalSearch_NoAPIKey(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

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

func TestBraveSearchService_HandleLocalSearch_HTTPError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

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

func TestBraveSearchService_HandleLocalSearch_GoroutineError(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	mockRateLimiter := &MockRateLimiter{}
	service := NewBraveSearchServiceWithAllDependencies(mockHTTPClient, mockEnvReader, mockRateLimiter)

	// レート制限を無効化
	mockRateLimiter.On("CheckLimit").Return(nil)

	// モックの設定（goroutineでエラーが発生するため、すべてのAPI呼び出しが実行されるとは限らない）
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key").Maybe()

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
	descErrorResponse := createMockResponse(http.StatusInternalServerError, map[string]string{"error": "Description Error"})

	// Web検索は成功、POI取得とDescription取得でエラー（goroutineで並行実行されるため両方設定）
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/web/search"
	})).Return(webMockResponse, nil).Once()

	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/local/pois"
	})).Return(poisErrorResponse, nil).Maybe()

	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/local/descriptions"
	})).Return(descErrorResponse, nil).Maybe()

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

func TestBraveSearchService_HandleLocalSearch_CountLimit(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	mockRateLimiter := &MockRateLimiter{}
	service := NewBraveSearchServiceWithAllDependencies(mockHTTPClient, mockEnvReader, mockRateLimiter)

	// レート制限を無効化
	mockRateLimiter.On("CheckLimit").Return(nil)

	// モックの設定
	mockEnvReader.On("Getenv", "BRAVE_API_KEY").Return("test-api-key").Times(2) // 2回のAPI呼び出し用

	// Web検索レスポンス（ロケーションIDなし、フォールバックする）
	webResponseData := BraveWebResponse{}
	webResponseData.Web.Results = []BraveWebResult{
		{Title: "Fallback", Description: "Fallback", URL: "https://example.com"},
	}

	webMockResponse := createMockResponse(http.StatusOK, webResponseData)
	fallbackMockResponse := createMockResponse(http.StatusOK, webResponseData)

	// 最初のWeb検索（ロケーション検索）
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/web/search" && req.URL.Query().Get("count") == "20" && req.URL.RawQuery != ""
	})).Return(webMockResponse, nil).Once()

	// フォールバック用のWeb検索
	mockHTTPClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/res/v1/web/search" && req.URL.Query().Get("count") == "20"
	})).Return(fallbackMockResponse, nil).Once()

	// テストの実行（制限を超えるcount=25を指定）
	result, err := service.HandleLocalSearch("test query", 25)

	// 結果の検証
	assert.NoError(t, err)
	assert.Contains(t, result, "Fallback")

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

func TestBraveSearchService_getPoisData_JSONDecodeError(t *testing.T) {
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
	result, err := service.getPoisData([]string{"location1"})

	// 結果の検証
	assert.Error(t, err)
	assert.Empty(t, result.Results)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_getPoisData_EmptyIDs(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

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

func TestBraveSearchService_getDescriptionsData_JSONDecodeError(t *testing.T) {
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
	result, err := service.getDescriptionsData([]string{"location1"})

	// 結果の検証
	assert.Error(t, err)
	assert.Empty(t, result.Descriptions)

	// モックの検証
	mockEnvReader.AssertExpectations(t)
	mockHTTPClient.AssertExpectations(t)
}

func TestBraveSearchService_getDescriptionsData_EmptyIDs(t *testing.T) {
	mockHTTPClient := &MockHTTPClient{}
	mockEnvReader := &MockEnvironmentReader{}
	service := NewBraveSearchServiceWithDependencies(mockHTTPClient, mockEnvReader)

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
