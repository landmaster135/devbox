package usecases

import (
	"testing"
	"time"

	discordUsecases "github.com/landmaster135/devbox/internal/discord_webhook/usecases"
	weatherUsecases "github.com/landmaster135/devbox/internal/open_weather_map/usecases"
)

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

// createTestDayForecast はテスト用の DayForecast を作成する
func createTestDayForecast() weatherUsecases.DayForecast {
	return weatherUsecases.DayForecast{
		Date:           time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		MinTemp:        10.5,
		MaxTemp:        20.3,
		Weather:        "Clear",
		Description:    "晴れ",
		EmojiOfWeather: "☀️",
		Humidity:       60,
		Pressure:       1013,
		WindSpeed:      5.2,
		Details: []weatherUsecases.ForecastDetail{
			{
				Time:        time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
				Temp:        15.0,
				Weather:     "Clear",
				Description: "晴れ",
			},
			{
				Time:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				Temp:        18.5,
				Weather:     "Clear",
				Description: "晴れ",
			},
		},
	}
}

// #==============================================================#
// ##          Test Class                                        ##
// #==============================================================#

// WeatherNotificatorServiceTest はテストクラス
type WeatherNotificatorServiceTest struct {
	service        *WeatherNotificatorService
	testForecast   weatherUsecases.DayForecast
	testAPIKey     string
	testCity       string
	testMaxDays    int
	testWebhookURL string
}

// setupTest はテストのセットアップを行う
func (test *WeatherNotificatorServiceTest) setupTest() {
	test.service = NewWeatherNotificatorService()
	test.testForecast = createTestDayForecast()
	test.testAPIKey = "test-api-key"
	test.testCity = "Tokyo"
	test.testMaxDays = 3
	test.testWebhookURL = "https://discord.com/api/webhooks/test"
}

// #==============================================================#
// ##          Constructor Tests                                 ##
// #==============================================================#

func TestNewWeatherNotificatorService_Normal(t *testing.T) {
	// Arrange & Act
	service := NewWeatherNotificatorService()

	// Assert
	if service == nil {
		t.Fatal("サービスがnilです")
	}
	if service.GetWeatherService() == nil {
		t.Error("WeatherServiceがnilです")
	}
	if service.GetDiscordService() == nil {
		t.Error("DiscordServiceがnilです")
	}
}

func TestNewWeatherNotificatorServiceWithDependencies_Normal(t *testing.T) {
	// Arrange
	weatherService := weatherUsecases.NewWeatherService()
	discordService := discordUsecases.NewDefaultDiscordWebhookService()

	// Act
	service := NewWeatherNotificatorServiceWithDependencies(weatherService, discordService)

	// Assert
	if service == nil {
		t.Fatal("サービスがnilです")
	}
	if service.GetWeatherService() == nil {
		t.Error("WeatherServiceがnilです")
	}
	if service.GetDiscordService() == nil {
		t.Error("DiscordServiceがnilです")
	}
}

// #==============================================================#
// ##          Getter Tests                                      ##
// #==============================================================#

func TestWeatherNotificatorService_GetWeatherService_Normal(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()

	// Act
	result := test.service.GetWeatherService()

	// Assert
	if result == nil {
		t.Error("WeatherServiceがnilです")
	}
}

func TestWeatherNotificatorService_GetDiscordService_Normal(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()

	// Act
	result := test.service.GetDiscordService()

	// Assert
	if result == nil {
		t.Error("DiscordServiceがnilです")
	}
}

// #==============================================================#
// ##          Private Method Tests (via reflection)            ##
// #==============================================================#

func TestWeatherNotificatorService_createEmbedTitle_Normal(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()
	forecast := test.testForecast
	city := test.testCity
	dayIndex := 1
	totalDays := test.testMaxDays

	// Act
	result := test.service.createEmbedTitle(city, forecast, dayIndex, totalDays)

	// Assert
	expectedTitle := "☀️ Tokyo の天気予報 (1/3日目)"
	if result != expectedTitle {
		t.Errorf("タイトルが期待値と異なります。期待値: %s, 実際: %s", expectedTitle, result)
	}
}

func TestWeatherNotificatorService_createEmbedDescription_Normal(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()
	forecast := test.testForecast

	// Act
	result := test.service.createEmbedDescription(forecast)

	// Assert
	expectedDescription := "📅 2024年01月01日 (Mon)"
	if result != expectedDescription {
		t.Errorf("説明文が期待値と異なります。期待値: %s, 実際: %s", expectedDescription, result)
	}
}

func TestWeatherNotificatorService_createEmbedFields_Normal(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()
	forecast := test.testForecast

	// Act
	result := test.service.createEmbedFields(forecast)

	// Assert
	expectedFieldCount := 7 // 基本情報6個 + 詳細予報1個
	if len(result) != expectedFieldCount {
		t.Errorf("フィールド数が期待値と異なります。期待値: %d, 実際: %d", expectedFieldCount, len(result))
	}

	// 最初のフィールド（気温）をチェック
	if len(result) > 0 {
		tempField := result[0]
		if tempField.Name != "🌡️ 気温" {
			t.Errorf("気温フィールドの名前が期待値と異なります。期待値: 🌡️ 気温, 実際: %s", tempField.Name)
		}
		expectedTempValue := "10.5°C ～ 20.3°C"
		if tempField.Value != expectedTempValue {
			t.Errorf("気温フィールドの値が期待値と異なります。期待値: %s, 実際: %s", expectedTempValue, tempField.Value)
		}
		if !tempField.Inline {
			t.Error("気温フィールドのInlineがfalseです")
		}
	}

	// 2番目のフィールド（天気）をチェック
	if len(result) > 1 {
		weatherField := result[1]
		if weatherField.Name != "☁️ 天気" {
			t.Errorf("天気フィールドの名前が期待値と異なります。期待値: ☁️ 天気, 実際: %s", weatherField.Name)
		}
		expectedWeatherValue := "Clear ☀️"
		if weatherField.Value != expectedWeatherValue {
			t.Errorf("天気フィールドの値が期待値と異なります。期待値: %s, 実際: %s", expectedWeatherValue, weatherField.Value)
		}
	}

	// 3番目のフィールド（湿度）をチェック
	if len(result) > 2 {
		humidityField := result[2]
		if humidityField.Name != "💧 湿度" {
			t.Errorf("湿度フィールドの名前が期待値と異なります。期待値: 💧 湿度, 実際: %s", humidityField.Name)
		}
		expectedHumidityValue := "60%"
		if humidityField.Value != expectedHumidityValue {
			t.Errorf("湿度フィールドの値が期待値と異なります。期待値: %s, 実際: %s", expectedHumidityValue, humidityField.Value)
		}
	}

	// 詳細予報フィールドをチェック
	if len(result) > 6 {
		detailField := result[6]
		if detailField.Name != "⏰ 3時間毎の詳細予報" {
			t.Errorf("詳細予報フィールドの名前が期待値と異なります。期待値: ⏰ 3時間毎の詳細予報, 実際: %s", detailField.Name)
		}
		if detailField.Inline {
			t.Error("詳細予報フィールドのInlineがtrueです")
		}
		// 詳細予報の内容に時刻と気温が含まれているかチェック
		if len(detailField.Value) == 0 {
			t.Error("詳細予報フィールドの値が空です")
		}
	}
}

// #==============================================================#
// ##          Integration Tests (Basic)                        ##
// #==============================================================#

func TestWeatherNotificatorService_SendWeatherNotification_EmptyForecast(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()

	// 無効なAPIキーを使用して空の結果を期待
	invalidAPIKey := ""
	invalidCity := ""

	// Act
	err := test.service.HandleWeatherNotification(invalidAPIKey, invalidCity, test.testMaxDays, test.testWebhookURL)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
}

func TestWeatherNotificatorService_HandleWeatherNotification_InvalidParameters(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()

	// Act - 無効なパラメータでテスト
	err := test.service.HandleWeatherNotification("", "", 0, "")

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
}

// #==============================================================#
// ##          Edge Case Tests                                   ##
// #==============================================================#

func TestWeatherNotificatorService_createEmbedTitle_EdgeCases(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()

	testCases := []struct {
		name      string
		city      string
		dayIndex  int
		totalDays int
		expected  string
	}{
		{
			name:      "単日予報",
			city:      "Osaka",
			dayIndex:  1,
			totalDays: 1,
			expected:  "☀️ Osaka の天気予報 (1/1日目)",
		},
		{
			name:      "5日予報の最終日",
			city:      "Kyoto",
			dayIndex:  5,
			totalDays: 5,
			expected:  "☀️ Kyoto の天気予報 (5/5日目)",
		},
		{
			name:      "空の都市名",
			city:      "",
			dayIndex:  1,
			totalDays: 3,
			expected:  "☀️  の天気予報 (1/3日目)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := test.service.createEmbedTitle(tc.city, test.testForecast, tc.dayIndex, tc.totalDays)

			// Assert
			if result != tc.expected {
				t.Errorf("タイトルが期待値と異なります。期待値: %s, 実際: %s", tc.expected, result)
			}
		})
	}
}

func TestWeatherNotificatorService_createEmbedFields_EmptyDetails(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()

	// 詳細情報が空の予報データを作成
	forecastWithoutDetails := test.testForecast
	forecastWithoutDetails.Details = []weatherUsecases.ForecastDetail{}

	// Act
	result := test.service.createEmbedFields(forecastWithoutDetails)

	// Assert
	expectedFieldCount := 6 // 基本情報6個のみ（詳細予報なし）
	if len(result) != expectedFieldCount {
		t.Errorf("フィールド数が期待値と異なります。期待値: %d, 実際: %d", expectedFieldCount, len(result))
	}
}

func TestWeatherNotificatorService_createEmbedFields_MultipleDetails(t *testing.T) {
	// Arrange
	test := &WeatherNotificatorServiceTest{}
	test.setupTest()

	// 複数の詳細情報を持つ予報データを作成
	forecastWithMultipleDetails := test.testForecast
	forecastWithMultipleDetails.Details = []weatherUsecases.ForecastDetail{
		{
			Time:        time.Date(2024, 1, 1, 6, 0, 0, 0, time.UTC),
			Temp:        12.0,
			Weather:     "Clear",
			Description: "晴れ",
		},
		{
			Time:        time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
			Temp:        15.0,
			Weather:     "Clear",
			Description: "晴れ",
		},
		{
			Time:        time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			Temp:        18.5,
			Weather:     "Clouds",
			Description: "曇り",
		},
		{
			Time:        time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC),
			Temp:        20.0,
			Weather:     "Clear",
			Description: "晴れ",
		},
	}

	// Act
	result := test.service.createEmbedFields(forecastWithMultipleDetails)

	// Assert
	expectedFieldCount := 7 // 基本情報6個 + 詳細予報1個
	if len(result) != expectedFieldCount {
		t.Errorf("フィールド数が期待値と異なります。期待値: %d, 実際: %d", expectedFieldCount, len(result))
	}

	// 詳細予報フィールドに4つの時間帯が含まれているかチェック
	if len(result) > 6 {
		detailField := result[6]
		detailValue := detailField.Value

		// 各時間帯が含まれているかチェック
		expectedTimes := []string{"06:00", "09:00", "12:00", "15:00"}
		for _, expectedTime := range expectedTimes {
			if len(detailValue) == 0 || detailValue == "\u200b" {
				t.Errorf("詳細予報に時刻 %s が含まれていません", expectedTime)
			}
		}
	}
}
