package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/open_weather_map/config"
	usecases "github.com/landmaster135/devbox/internal/open_weather_map/usecases"
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

	// WeatherServiceを初期化
	service := usecases.NewWeatherService()

	// 天気予報を取得
	result, err := service.HandleWeatherForecast(cfg.APIKey, cfg.City, cfg.MaxDays)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}
