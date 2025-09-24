package usecases

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPClient はHTTPリクエストを実行するためのインターフェース
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient は標準のhttp.Clientを使用する実装
type DefaultHTTPClient struct {
	client *http.Client
}

// NewDefaultHTTPClient は新しいDefaultHTTPClientを作成する
func NewDefaultHTTPClient() *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Do はHTTPリクエストを実行する
func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// ForecastResponse は5日間予報用の構造体
type ForecastResponse struct {
	List []struct {
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
	} `json:"list"`
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
		Coord   struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
	} `json:"city"`
}

// DayForecast は日別予報の構造体
type DayForecast struct {
	Date           time.Time
	MinTemp        float64
	MaxTemp        float64
	Weather        string
	Description    string
	EmojiOfWeather string
	Humidity       int
	Pressure       int
	WindSpeed      float64
	Details        []ForecastDetail
}

// ForecastDetail は3時間毎の詳細予報の構造体
type ForecastDetail struct {
	Time        time.Time
	Temp        float64
	Weather     string
	Description string
}

// getWeatherEmoji は天気の種類に基づいて適切な絵文字を返す
func getWeatherEmoji(weatherMain string) string {
	switch strings.ToLower(weatherMain) {
	case "clear":
		return "☀️"
	case "clouds":
		return "☁️"
	case "rain":
		return "☔️"
	case "drizzle":
		return "🌦️"
	case "thunderstorm":
		return "⛈️"
	case "snow":
		return "❄️"
	case "mist", "fog", "haze":
		return "🌫️"
	case "dust", "sand", "ash":
		return "🌪️"
	case "squall", "tornado":
		return "🌪️"
	default:
		return "🌤️"
	}
}

// WeatherService はOpenWeather APIを使用した天気予報サービス
type WeatherService struct {
	httpClient HTTPClient
	baseURL    string
}

// NewWeatherService は新しいWeatherServiceを作成する
func NewWeatherService() *WeatherService {
	return &WeatherService{
		httpClient: NewDefaultHTTPClient(),
		baseURL:    "https://api.openweathermap.org/data/2.5",
	}
}

// NewWeatherServiceWithHTTPClient はHTTPクライアントを注入したWeatherServiceを作成する
func NewWeatherServiceWithHTTPClient(httpClient HTTPClient) *WeatherService {
	return &WeatherService{
		httpClient: httpClient,
		baseURL:    "https://api.openweathermap.org/data/2.5",
	}
}

// GetForecast5Days は5日間3時間毎の予報を取得する
func (w *WeatherService) GetForecast5Days(apiKey, city string) (*ForecastResponse, error) {
	params := url.Values{}
	params.Add("q", city)
	params.Add("appid", apiKey)
	params.Add("units", "metric")
	params.Add("lang", "ja")

	requestURL := fmt.Sprintf("%s/forecast?%s", w.baseURL, params.Encode())

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %v", err)
	}

	resp, err := w.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("APIリクエストエラー: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: ステータスコード %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	var forecastData ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecastData); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		switch {
		case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
			return nil, fmt.Errorf("レスポンス読み取りエラー: %v", err)
		case errors.As(err, &syntaxErr), errors.As(err, &typeErr):
			return nil, fmt.Errorf("JSON解析エラー: %v", err)
		default:
			return nil, fmt.Errorf("レスポンス読み取りエラー: %v", err)
		}
	}

	return &forecastData, nil
}

// aggregateForecastByDays は3時間毎のデータを日別に集約する
func (w *WeatherService) aggregateForecastByDays(forecast *ForecastResponse, days int) []DayForecast {
	dayMap := make(map[string]*DayForecast)

	for _, item := range forecast.List {
		t := time.Unix(item.Dt, 0)
		dateKey := t.Format("2006-01-02")

		// 指定日数を超えた場合はスキップ
		if len(dayMap) >= days && dayMap[dateKey] == nil {
			continue
		}

		if dayMap[dateKey] == nil {
			dayMap[dateKey] = &DayForecast{
				Date:           time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
				MinTemp:        item.Main.TempMin,
				MaxTemp:        item.Main.TempMax,
				Weather:        item.Weather[0].Main,
				Description:    item.Weather[0].Description,
				EmojiOfWeather: getWeatherEmoji(item.Weather[0].Main),
				Humidity:       item.Main.Humidity,
				Pressure:       item.Main.Pressure,
				WindSpeed:      item.Wind.Speed,
				Details:        []ForecastDetail{},
			}
		}

		day := dayMap[dateKey]

		// 最低・最高気温を更新
		if item.Main.TempMin < day.MinTemp {
			day.MinTemp = item.Main.TempMin
		}
		if item.Main.TempMax > day.MaxTemp {
			day.MaxTemp = item.Main.TempMax
		}

		// 詳細情報を追加
		day.Details = append(day.Details, ForecastDetail{
			Time:        t,
			Temp:        item.Main.Temp,
			Weather:     item.Weather[0].Main,
			Description: item.Weather[0].Description,
		})
	}

	// 日付順でソート
	var result []DayForecast
	now := time.Now()
	for i := 0; i < days && i < len(dayMap); i++ {
		targetDate := now.AddDate(0, 0, i)
		dateKey := targetDate.Format("2006-01-02")
		if day, exists := dayMap[dateKey]; exists {
			result = append(result, *day)
		}
	}

	return result
}

// GetForecastByDays は指定した日数の予報を取得する
func (w *WeatherService) GetForecastByDays(apiKey, city string, days int) ([]DayForecast, error) {
	if days > 5 {
		return nil, fmt.Errorf("無料プランでは最大5日間の予報のみ取得可能です")
	}

	forecast, err := w.GetForecast5Days(apiKey, city)
	if err != nil {
		return nil, err
	}

	return w.aggregateForecastByDays(forecast, days), nil
}

// FormatForecastOutput は予報データを見やすい形式でフォーマットする
func (w *WeatherService) FormatForecastOutput(forecasts []DayForecast, city string, days int) string {
	var output strings.Builder

	output.WriteString(fmt.Sprintf("=== %s の%d日間天気予報 ===\n\n", city, days))

	for _, day := range forecasts {
		output.WriteString(fmt.Sprintf("📅 %s\n", day.Date.Format("2006年01月02日 (Mon)")))
		output.WriteString(fmt.Sprintf("🌡️  気温: %.1f°C ～ %.1f°C\n", day.MinTemp, day.MaxTemp))
		output.WriteString(fmt.Sprintf("☁️  天気: %s (%s) %s\n", day.Weather, day.Description, day.EmojiOfWeather))
		output.WriteString(fmt.Sprintf("💨 湿度: %d%% | 気圧: %d hPa | 風速: %.1f m/s\n",
			day.Humidity, day.Pressure, day.WindSpeed))

		output.WriteString("⏰ 3時間毎の詳細:\n")
		for _, detail := range day.Details {
			output.WriteString(fmt.Sprintf("   %s: %.1f°C %s\n",
				detail.Time.Format("15:04"), detail.Temp, detail.Description))
		}
		output.WriteString(strings.Repeat("-", 50) + "\n")
	}

	return output.String()
}

// HandleWeatherForecast は天気予報取得のメインハンドラー
func (w *WeatherService) HandleWeatherForecast(apiKey, city string, maxDays int) (string, error) {
	forecasts, err := w.GetForecastByDays(apiKey, city, maxDays)
	if err != nil {
		return "", err
	}

	return w.FormatForecastOutput(forecasts, city, maxDays), nil
}
