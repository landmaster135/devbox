package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/weather_notificator/config"
	usecases "github.com/landmaster135/devbox/internal/weather_notificator/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	// 天気通知を実行
	handleWeatherNotification(cfg)
}

// handleWeatherNotification は天気通知を処理する
func handleWeatherNotification(cfg *config.Config) {
	// WeatherNotificatorServiceを初期化
	service := usecases.NewWeatherNotificatorService()

	// 天気通知を実行
	err := service.HandleWeatherNotification(cfg.APIKey, cfg.City, cfg.MaxDays, cfg.WebhookURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 成功メッセージを出力
	fmt.Printf("✅ %sの%d日間天気予報をDiscordに送信しました\n", cfg.City, cfg.MaxDays)
}
