package server

import (
	"log"
	"net/http"

	weatherHandler "github.com/landmaster135/devbox/cmd/http/handlers/weather_notificator"
)

// setupRouter はHTTPルーターを設定して返す
func setupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// 天気通知ハンドラを初期化
	weatherNotificationHandler := weatherHandler.NewWeatherNotificationHandler()

	// ルートの設定
	mux.HandleFunc("/weather-notification", weatherNotificationHandler.HandleWeatherNotification)

	// 以降、ここにルートを登録

	// ヘルスチェックエンドポイント
	mux.HandleFunc("/health", handleHealth)

	// ルートエンドポイント（API情報）
	mux.HandleFunc("/", handleRoot)

	log.Println("ルーターが初期化されました")
	log.Println("利用可能なエンドポイント:")
	log.Println("  POST /weather-notification - 天気通知を送信")
	log.Println("  GET  /health               - ヘルスチェック")
	log.Println("  GET  /                     - API情報")

	return mux
}
