package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	discord "github.com/landmaster135/devbox/internal/discord_webhook/infrastructure/discord"
	discordUsecases "github.com/landmaster135/devbox/internal/discord_webhook/usecases"
	weatherUsecases "github.com/landmaster135/devbox/internal/open_weather_map/usecases"
)

// WeatherNotificatorService は天気通知サービス
type WeatherNotificatorService struct {
	weatherService *weatherUsecases.WeatherService
	discordService *discordUsecases.DiscordWebhookService
}

// NewWeatherNotificatorService は新しいWeatherNotificatorServiceを作成する
func NewWeatherNotificatorService() *WeatherNotificatorService {
	return &WeatherNotificatorService{
		weatherService: weatherUsecases.NewWeatherService(),
		discordService: discordUsecases.NewDefaultDiscordWebhookService(),
	}
}

// NewWeatherNotificatorServiceWithDependencies は依存関係を注入したWeatherNotificatorServiceを作成する
func NewWeatherNotificatorServiceWithDependencies(
	weatherService *weatherUsecases.WeatherService,
	discordService *discordUsecases.DiscordWebhookService,
) *WeatherNotificatorService {
	return &WeatherNotificatorService{
		weatherService: weatherService,
		discordService: discordService,
	}
}

// SendWeatherNotification は指定した都市の天気予報をDiscordに通知する
func (s *WeatherNotificatorService) SendWeatherNotification(ctx context.Context, apiKey, city string, maxDays int, webhookURL string) error {
	// 天気予報を取得
	forecasts, err := s.weatherService.GetForecastByDays(apiKey, city, maxDays)
	if err != nil {
		return fmt.Errorf("天気予報の取得に失敗しました: %w", err)
	}

	if len(forecasts) == 0 {
		return fmt.Errorf("天気予報データが取得できませんでした")
	}

	// 各日の天気予報をDiscordに送信
	for i, forecast := range forecasts {
		embedTitle := s.createEmbedTitle(city, forecast, i+1, maxDays)
		embedDescription := s.createEmbedDescription(forecast)
		embedColor := "orange"
		fields := s.createEmbedFields(forecast)

		// 新しい天気予報専用の通知を送信
		err := s.discordService.SendWeatherNotification(
			ctx,
			webhookURL,
			embedTitle,
			embedDescription,
			embedColor,
			fields,
		)
		if err != nil {
			return fmt.Errorf("Discord通知の送信に失敗しました（%d日目）: %w", i+1, err)
		}

		// 連続送信時の負荷軽減のため少し待機
		if i < len(forecasts)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	return nil
}

// createEmbedTitle はembedのタイトルを作成する
func (s *WeatherNotificatorService) createEmbedTitle(city string, forecast weatherUsecases.DayForecast, dayIndex, totalDays int) string {
	return fmt.Sprintf("%s %s の天気予報 (%d/%d日目)", forecast.EmojiOfWeather, city, dayIndex, totalDays)
}

// createEmbedDescription はembedの説明文を作成する
func (s *WeatherNotificatorService) createEmbedDescription(forecast weatherUsecases.DayForecast) string {
	return fmt.Sprintf("📅 %s", forecast.Date.Format("2006年01月02日 (Mon)"))
}

// createEmbedFields はembedのフィールドを作成する
func (s *WeatherNotificatorService) createEmbedFields(forecast weatherUsecases.DayForecast) []*discord.EmbedField {
	var fields []*discord.EmbedField

	// 基本天気情報（横並び表示）
	fields = append(fields, &discord.EmbedField{
		Name:   "🌡️ 気温",
		Value:  fmt.Sprintf("%.1f°C ～ %.1f°C", forecast.MinTemp, forecast.MaxTemp),
		Inline: true,
	})

	fields = append(fields, &discord.EmbedField{
		Name:   "☁️ 天気",
		Value:  fmt.Sprintf("%s %s", forecast.Weather, forecast.EmojiOfWeather),
		Inline: true,
	})

	fields = append(fields, &discord.EmbedField{
		Name:   "💧 湿度",
		Value:  fmt.Sprintf("%d%%", forecast.Humidity),
		Inline: true,
	})

	// 2行目の情報（横並び表示）
	fields = append(fields, &discord.EmbedField{
		Name:   "🌪️ 気圧",
		Value:  fmt.Sprintf("%d hPa", forecast.Pressure),
		Inline: true,
	})

	fields = append(fields, &discord.EmbedField{
		Name:   "💨 風速",
		Value:  fmt.Sprintf("%.1f m/s", forecast.WindSpeed),
		Inline: true,
	})

	// 空のフィールドで改行効果
	fields = append(fields, &discord.EmbedField{
		Name:   "\u200b", // 不可視文字
		Value:  "\u200b", // 不可視文字
		Inline: true,
	})

	// 3時間毎の詳細予報（縦表示）
	if len(forecast.Details) > 0 {
		var detailText strings.Builder
		for _, detail := range forecast.Details {
			detailText.WriteString(fmt.Sprintf("• %s: %.1f°C - %s\n",
				detail.Time.Format("15:04"), detail.Temp, detail.Description))
		}

		fields = append(fields, &discord.EmbedField{
			Name:   "⏰ 3時間毎の詳細予報",
			Value:  detailText.String(),
			Inline: false,
		})
	}

	return fields
}

// HandleWeatherNotification は天気通知のメインハンドラー
func (s *WeatherNotificatorService) HandleWeatherNotification(apiKey, city string, maxDays int, webhookURL string) error {
	ctx := context.Background()
	return s.SendWeatherNotification(ctx, apiKey, city, maxDays, webhookURL)
}
