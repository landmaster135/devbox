package server

import (
	"net/http"

	cronWorkflowHandler "github.com/landmaster135/devbox/cmd/http/handlers/cron_workflow"
	weatherHandler "github.com/landmaster135/devbox/cmd/http/handlers/weather_notificator"
	logging "github.com/landmaster135/devbox/internal/logging"
)

// setupRouter はHTTPルーターを設定して返す
func setupRouter(logger *logging.StructuredLogger) *http.ServeMux {
	mux := http.NewServeMux()
	routerLogger := logging.Ensure(logger).WithTags("router")

	// ハンドラを初期化
	weatherNotificationHandler := weatherHandler.NewWeatherNotificationHandler(routerLogger.WithTags("handler", "weather-notification"))
	cronWorkflowPageHandler := cronWorkflowHandler.NewHandler(routerLogger.WithTags("handler", "cron-workflow"))

	// ルートの設定
	mux.HandleFunc("/weather-notification", weatherNotificationHandler.HandleWeatherNotification)
	mux.HandleFunc(cronWorkflowHandler.BaseEndpoint, cronWorkflowPageHandler.HandleCronWorkflowPage)
	mux.HandleFunc(cronWorkflowHandler.ManualRunEndpoint, cronWorkflowPageHandler.HandleManualRun)

	// 以降、ここにルートを登録

	// ヘルスチェックエンドポイント
	mux.HandleFunc("/health", handleHealth(routerLogger.WithTags("endpoint", "health")))

	// ルートエンドポイント（API情報）
	mux.HandleFunc("/", handleRoot(routerLogger.WithTags("endpoint", "root")))

	routerLogger.Infof("ルーターが初期化されました")
	routerLogger.Infof("利用可能なエンドポイント:")
	routerLogger.Infof("  POST /weather-notification     - 天気通知を送信")
	routerLogger.Infof("  GET  /cron-workflow            - ワークフローダッシュボード")
	routerLogger.Infof("  POST /cron-workflow/manual-run - ワークフロー手動実行")
	routerLogger.Infof("  GET  /health               - ヘルスチェック")
	routerLogger.Infof("  GET  /                     - API情報")

	return mux
}
