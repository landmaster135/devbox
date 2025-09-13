package router

import (
	"encoding/json"
	"log"
	"net/http"
)

// APIInfo はAPI情報を表す構造体
type APIInfo struct {
	Name        string         `json:"name"`
	Version     string         `json:"version"`
	Description string         `json:"description"`
	Endpoints   []EndpointInfo `json:"endpoints"`
}

// EndpointInfo はエンドポイント情報を表す構造体
type EndpointInfo struct {
	Method      string         `json:"method"`
	Path        string         `json:"path"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Example     map[string]any `json:"example,omitempty"`
}

// createAPIInfo はAPI情報を構造体として作成する
func createAPIInfo() APIInfo {
	return APIInfo{
		Name:        "Weather Notification API",
		Version:     "1.0.0",
		Description: "天気予報をDiscordに通知するAPIサービス",
		Endpoints: []EndpointInfo{
			{
				Method:      "POST",
				Path:        "/weather-notification",
				Description: "指定した都市の天気予報をDiscordに送信",
				Parameters: map[string]any{
					"api_key":     "OpenWeather API キー（必須）",
					"city":        "都市名（必須）",
					"max_days":    "最大日数（必須、1-5日）",
					"webhook_url": "Discord Webhook URL（必須）",
				},
				Example: map[string]any{
					"api_key":     "your_openweather_api_key",
					"city":        "Tokyo",
					"max_days":    3,
					"webhook_url": "https://discord.com/api/webhooks/...",
				},
			},
			{
				Method:      "GET",
				Path:        "/health",
				Description: "APIのヘルスチェック",
			},
			{
				Method:      "GET",
				Path:        "/",
				Description: "API情報の取得",
			},
		},
	}
}

// handleHealth はヘルスチェックエンドポイントのハンドラ
func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GETメソッドのみサポートしています", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "message": "Weather Notification API is running"}`))
}

// handleRoot はルートエンドポイントのハンドラ
func handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GETメソッドのみサポートしています", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 構造体からAPI情報を作成
	apiInfo := createAPIInfo()

	// JSONに変換
	jsonData, err := json.MarshalIndent(apiInfo, "", "  ")
	if err != nil {
		log.Printf("JSON変換エラー: %v", err)
		http.Error(w, "内部サーバーエラー", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
