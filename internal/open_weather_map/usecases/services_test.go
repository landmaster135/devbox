package usecases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// MockHTTPClient はテスト用のHTTPクライアントモック
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do はHTTPリクエストを実行する
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

// createMockForecastResponse はテスト用のモック予報レスポンスを作成する
func createMockForecastResponse() *ForecastResponse {
	now := time.Now()
	return &ForecastResponse{
		List: []struct {
			Dt   int64 `json:"dt"`
			Main struct {
				Temp      float64 `json:"temp"`
				FeelsLike float64 `json:"feels_like"`
				TempMin   float64 `json:"temp_min"`
				TempMax   float64 `json:"temp_max"`
				Pressure  int     `json:"pressure"`
				Humidity  int     `json:"humidity"`
			} `json:"main"`
			Weather []struct {
				Main        string `json:"main"`
				Description string `json:"description"`
				Icon        string `json:"icon"`
			} `json:"weather"`
			Wind struct {
				Speed float64 `json:"speed"`
				Deg   int     `json:"deg"`
			} `json:"wind"`
			DtTxt string `json:"dt_txt"`
		}{
			{
				Dt: now.Unix(),
				Main: struct {
					Temp      float64 `json:"temp"`
					FeelsLike float64 `json:"feels_like"`
					TempMin   float64 `json:"temp_min"`
					TempMax   float64 `json:"temp_max"`
					Pressure  int     `json:"pressure"`
					Humidity  int     `json:"humidity"`
				}{
					Temp:      20.5,
					FeelsLike: 22.0,
					TempMin:   18.0,
					TempMax:   23.0,
					Pressure:  1013,
					Humidity:  65,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
					Icon        string `json:"icon"`
				}{
					{
						Main:        "Clear",
						Description: "晴れ",
						Icon:        "01d",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 3.5,
					Deg:   180,
				},
				DtTxt: now.Format("2006-01-02 15:04:05"),
			},
		},
		City: struct {
			Name    string `json:"name"`
			Country string `json:"country"`
			Coord   struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"coord"`
		}{
			Name:    "Tokyo",
			Country: "JP",
			Coord: struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			}{
				Lat: 35.6762,
				Lon: 139.6503,
			},
		},
	}
}

// TestNewWeatherService_Normal は正常なWeatherServiceの作成をテストする
func TestNewWeatherService_Normal(t *testing.T) {
	service := NewWeatherService()

	if service == nil {
		t.Error("NewWeatherService() returned nil")
		return
	}

	if service.baseURL != "https://api.openweathermap.org/data/2.5" {
		t.Errorf("baseURL = %v, want https://api.openweathermap.org/data/2.5", service.baseURL)
	}

	if service.httpClient == nil {
		t.Error("httpClient is nil")
	}
}

// TestNewWeatherServiceWithHTTPClient_Normal はHTTPクライアント注入のテストを行う
func TestNewWeatherServiceWithHTTPClient_Normal(t *testing.T) {
	mockClient := &MockHTTPClient{}
	service := NewWeatherServiceWithHTTPClient(mockClient)

	if service == nil {
		t.Error("NewWeatherServiceWithHTTPClient() returned nil")
		return
	}

	if service.httpClient != mockClient {
		t.Error("httpClient was not injected correctly")
	}
}

// TestWeatherService_GetForecast5Days_Normal は正常な5日間予報取得をテストする
func TestWeatherService_GetForecast5Days_Normal(t *testing.T) {
	mockResponse := createMockForecastResponse()
	responseBody, _ := json.Marshal(mockResponse)

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			// リクエストURLの検証
			expectedParams := []string{"q=Tokyo%2CJP", "appid=test-api-key", "units=metric", "lang=ja"}
			for _, param := range expectedParams {
				if !strings.Contains(req.URL.RawQuery, param) {
					t.Errorf("Expected parameter %s not found in URL: %s", param, req.URL.RawQuery)
				}
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(responseBody)),
			}, nil
		},
	}

	service := NewWeatherServiceWithHTTPClient(mockClient)
	result, err := service.GetForecast5Days("test-api-key", "Tokyo,JP")

	if err != nil {
		t.Errorf("GetForecast5Days() error = %v, want nil", err)
		return
	}

	if result == nil {
		t.Error("GetForecast5Days() returned nil result")
		return
	}

	if result.City.Name != "Tokyo" {
		t.Errorf("City.Name = %v, want Tokyo", result.City.Name)
	}

	if len(result.List) != 1 {
		t.Errorf("len(List) = %v, want 1", len(result.List))
	}
}

// TestWeatherService_GetForecast5Days_APIError はAPIエラーのテストを行う
func TestWeatherService_GetForecast5Days_APIError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"cod":401,"message":"Invalid API key"}`)),
			}, nil
		},
	}

	service := NewWeatherServiceWithHTTPClient(mockClient)
	_, err := service.GetForecast5Days("invalid-key", "Tokyo,JP")

	if err == nil {
		t.Error("GetForecast5Days() error = nil, want error")
		return
	}

	if !strings.Contains(err.Error(), "APIエラー") {
		t.Errorf("Error message = %v, want to contain 'APIエラー'", err.Error())
	}
}

// TestWeatherService_GetForecast5Days_NetworkError はネットワークエラーのテストを行う
func TestWeatherService_GetForecast5Days_NetworkError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	service := NewWeatherServiceWithHTTPClient(mockClient)
	_, err := service.GetForecast5Days("test-key", "Tokyo,JP")

	if err == nil {
		t.Error("GetForecast5Days() error = nil, want error")
		return
	}

	if !strings.Contains(err.Error(), "APIリクエストエラー") {
		t.Errorf("Error message = %v, want to contain 'APIリクエストエラー'", err.Error())
	}
}

// TestWeatherService_GetForecast5Days_JSONError はJSON解析エラーのテストを行う
func TestWeatherService_GetForecast5Days_JSONError(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{invalid json}")),
			}, nil
		},
	}

	service := NewWeatherServiceWithHTTPClient(mockClient)
	_, err := service.GetForecast5Days("test-key", "Tokyo,JP")

	if err == nil {
		t.Fatal("GetForecast5Days() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "JSON解析エラー") {
		t.Errorf("Error message = %v, want to contain 'JSON解析エラー'", err.Error())
	}
}

// TestWeatherService_GetForecastByDays_Normal は正常な日数指定予報取得をテストする
func TestWeatherService_GetForecastByDays_Normal(t *testing.T) {
	mockResponse := createMockForecastResponse()
	responseBody, _ := json.Marshal(mockResponse)

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(responseBody)),
			}, nil
		},
	}

	service := NewWeatherServiceWithHTTPClient(mockClient)
	result, err := service.GetForecastByDays("test-key", "Tokyo,JP", 3)

	if err != nil {
		t.Errorf("GetForecastByDays() error = %v, want nil", err)
		return
	}

	if result == nil {
		t.Error("GetForecastByDays() returned nil result")
		return
	}

	if len(result) > 3 {
		t.Errorf("len(result) = %v, want <= 3", len(result))
	}
}

// TestWeatherService_GetForecastByDays_ExceedMaxDays は最大日数超過のテストを行う
func TestWeatherService_GetForecastByDays_ExceedMaxDays(t *testing.T) {
	service := NewWeatherService()
	_, err := service.GetForecastByDays("test-key", "Tokyo,JP", 6)

	if err == nil {
		t.Error("GetForecastByDays() error = nil, want error")
		return
	}

	if !strings.Contains(err.Error(), "無料プランでは最大5日間") {
		t.Errorf("Error message = %v, want to contain '無料プランでは最大5日間'", err.Error())
	}
}

// TestWeatherService_FormatForecastOutput_Normal は予報出力フォーマットのテストを行う
func TestWeatherService_FormatForecastOutput_Normal(t *testing.T) {
	service := NewWeatherService()

	forecasts := []DayForecast{
		{
			Date:           time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			MinTemp:        10.0,
			MaxTemp:        20.0,
			Weather:        "Clear",
			Description:    "晴れ",
			EmojiOfWeather: "☀️",
			Humidity:       65,
			Pressure:       1013,
			WindSpeed:      3.5,
			Details: []ForecastDetail{
				{
					Time:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
					Temp:        15.0,
					Weather:     "Clear",
					Description: "晴れ",
				},
			},
		},
	}

	result := service.FormatForecastOutput(forecasts, "Tokyo,JP", 1)

	if result == "" {
		t.Error("FormatForecastOutput() returned empty string")
		return
	}

	expectedStrings := []string{
		"Tokyo,JP の1日間天気予報",
		"2024年01月01日",
		"気温: 10.0°C ～ 20.0°C",
		"天気: Clear (晴れ) ☀️",
		"湿度: 65%",
		"気圧: 1013 hPa",
		"風速: 3.5 m/s",
		"3時間毎の詳細",
		"12:00: 15.0°C 晴れ",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(result, expected) {
			t.Errorf("FormatForecastOutput() result does not contain '%s'", expected)
		}
	}
}

// TestWeatherService_AggregateForecastByDays_UsesJST は日付集約がJSTで行われることをテストする
func TestWeatherService_AggregateForecastByDays_UsesJST(t *testing.T) {
	service := NewWeatherService()
	nowJST := time.Now().In(jstLocation)
	base := time.Date(nowJST.Year(), nowJST.Month(), nowJST.Day(), 6, 0, 0, 0, jstLocation)

	forecast := &ForecastResponse{
		List: []struct {
			Dt   int64 `json:"dt"`
			Main struct {
				Temp      float64 `json:"temp"`
				FeelsLike float64 `json:"feels_like"`
				TempMin   float64 `json:"temp_min"`
				TempMax   float64 `json:"temp_max"`
				Pressure  int     `json:"pressure"`
				Humidity  int     `json:"humidity"`
			} `json:"main"`
			Weather []struct {
				Main        string `json:"main"`
				Description string `json:"description"`
				Icon        string `json:"icon"`
			} `json:"weather"`
			Wind struct {
				Speed float64 `json:"speed"`
				Deg   int     `json:"deg"`
			} `json:"wind"`
			DtTxt string `json:"dt_txt"`
		}{
			{
				Dt: base.UTC().Unix(),
				Main: struct {
					Temp      float64 `json:"temp"`
					FeelsLike float64 `json:"feels_like"`
					TempMin   float64 `json:"temp_min"`
					TempMax   float64 `json:"temp_max"`
					Pressure  int     `json:"pressure"`
					Humidity  int     `json:"humidity"`
				}{
					Temp:      18.0,
					FeelsLike: 18.0,
					TempMin:   17.0,
					TempMax:   20.0,
					Pressure:  1012,
					Humidity:  60,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
					Icon        string `json:"icon"`
				}{
					{
						Main:        "Clouds",
						Description: "曇り",
						Icon:        "02d",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 2.5,
					Deg:   90,
				},
				DtTxt: base.Format("2006-01-02 15:04:05"),
			},
			{
				Dt: base.Add(12 * time.Hour).UTC().Unix(),
				Main: struct {
					Temp      float64 `json:"temp"`
					FeelsLike float64 `json:"feels_like"`
					TempMin   float64 `json:"temp_min"`
					TempMax   float64 `json:"temp_max"`
					Pressure  int     `json:"pressure"`
					Humidity  int     `json:"humidity"`
				}{
					Temp:      22.0,
					FeelsLike: 21.5,
					TempMin:   16.0,
					TempMax:   23.5,
					Pressure:  1010,
					Humidity:  55,
				},
				Weather: []struct {
					Main        string `json:"main"`
					Description string `json:"description"`
					Icon        string `json:"icon"`
				}{
					{
						Main:        "Clear",
						Description: "晴れ",
						Icon:        "01n",
					},
				},
				Wind: struct {
					Speed float64 `json:"speed"`
					Deg   int     `json:"deg"`
				}{
					Speed: 1.8,
					Deg:   110,
				},
				DtTxt: base.Add(12 * time.Hour).Format("2006-01-02 15:04:05"),
			},
		},
	}

	result := service.aggregateForecastByDays(forecast, 1)
	if len(result) != 1 {
		t.Fatalf("aggregateForecastByDays() len = %d, want 1", len(result))
	}

	day := result[0]
	if day.Date.Location() != jstLocation {
		t.Errorf("Date location = %v, want JST", day.Date.Location())
	}

	expectedDate := time.Date(base.Year(), base.Month(), base.Day(), 0, 0, 0, 0, jstLocation)
	if !day.Date.Equal(expectedDate) {
		t.Errorf("Date = %v, want %v", day.Date, expectedDate)
	}

	if len(day.Details) != 2 {
		t.Fatalf("Details len = %d, want 2", len(day.Details))
	}

	if day.Details[0].Time.Location() != jstLocation || day.Details[1].Time.Location() != jstLocation {
		t.Error("Detail times are not converted to JST")
	}

	if day.MinTemp != 16.0 {
		t.Errorf("MinTemp = %v, want 16.0", day.MinTemp)
	}
	if day.MaxTemp != 23.5 {
		t.Errorf("MaxTemp = %v, want 23.5", day.MaxTemp)
	}
}

// TestGetWeatherEmoji_Normal は天気絵文字マッピングのテストを行う
func TestGetWeatherEmoji_Normal(t *testing.T) {
	tests := []struct {
		weatherMain string
		expected    string
	}{
		{"Clear", "☀️"},
		{"clear", "☀️"},
		{"Clouds", "☁️"},
		{"clouds", "☁️"},
		{"Rain", "☔️"},
		{"rain", "☔️"},
		{"Drizzle", "🌦️"},
		{"drizzle", "🌦️"},
		{"Thunderstorm", "⛈️"},
		{"thunderstorm", "⛈️"},
		{"Snow", "❄️"},
		{"snow", "❄️"},
		{"Mist", "🌫️"},
		{"mist", "🌫️"},
		{"Fog", "🌫️"},
		{"fog", "🌫️"},
		{"Haze", "🌫️"},
		{"haze", "🌫️"},
		{"Dust", "🌪️"},
		{"dust", "🌪️"},
		{"Sand", "🌪️"},
		{"sand", "🌪️"},
		{"Ash", "🌪️"},
		{"ash", "🌪️"},
		{"Squall", "🌪️"},
		{"squall", "🌪️"},
		{"Tornado", "🌪️"},
		{"tornado", "🌪️"},
		{"Unknown", "🌤️"},
		{"", "🌤️"},
	}

	for _, tt := range tests {
		t.Run(tt.weatherMain, func(t *testing.T) {
			result := getWeatherEmoji(tt.weatherMain)
			if result != tt.expected {
				t.Errorf("getWeatherEmoji(%s) = %s, want %s", tt.weatherMain, result, tt.expected)
			}
		})
	}
}

// TestWeatherService_HandleWeatherForecast_Normal はメインハンドラーのテストを行う
func TestWeatherService_HandleWeatherForecast_Normal(t *testing.T) {
	mockResponse := createMockForecastResponse()
	responseBody, _ := json.Marshal(mockResponse)

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader(responseBody)),
			}, nil
		},
	}

	service := NewWeatherServiceWithHTTPClient(mockClient)
	result, err := service.HandleWeatherForecast("test-key", "Tokyo,JP", 3)

	if err != nil {
		t.Errorf("HandleWeatherForecast() error = %v, want nil", err)
		return
	}

	if result == "" {
		t.Error("HandleWeatherForecast() returned empty string")
		return
	}

	if !strings.Contains(result, "Tokyo,JP の3日間天気予報") {
		t.Errorf("HandleWeatherForecast() result does not contain expected header")
	}
}

// TestWeatherService_HandleWeatherForecast_Error はハンドラーでのエラー伝播をテストする
func TestWeatherService_HandleWeatherForecast_Error(t *testing.T) {
	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"cod":401,"message":"Invalid API key"}`)),
			}, nil
		},
	}

	service := NewWeatherServiceWithHTTPClient(mockClient)
	_, err := service.HandleWeatherForecast("bad-key", "Tokyo,JP", 3)

	if err == nil {
		t.Fatal("HandleWeatherForecast() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "APIエラー") {
		t.Errorf("Error message = %v, want to contain 'APIエラー'", err.Error())
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// TestHTTPClient_Do は標準HTTPクライアントのDoをテストする
func TestHTTPClient_Do(t *testing.T) {
	client := NewHTTPClient()
	client.client.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://example.com" {
			t.Errorf("unexpected request URL: %s", req.URL.String())
		}
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(strings.NewReader("ok")),
		}, nil
	})

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v, want nil", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusAccepted)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if string(body) != "ok" {
		t.Errorf("body = %s, want ok", string(body))
	}
}
