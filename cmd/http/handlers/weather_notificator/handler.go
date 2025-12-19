package weather_notificator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/landmaster135/devbox/internal/weather_notificator/usecases"
)

// WeatherNotificationRequest は天気通知リクエストの構造体
type WeatherNotificationRequest struct {
	APIKey     string `json:"api_key"`
	City       string `json:"city"`
	MaxDays    int    `json:"max_days"`
	WebhookURL string `json:"webhook_url"`
}

// WeatherNotificationResponse は天気通知レスポンスの構造体
type WeatherNotificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ErrorResponse はエラーレスポンスの構造体
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// WeatherNotificationHandler は天気通知のHTTPハンドラ
type WeatherNotificationHandler struct {
	service *usecases.WeatherNotificatorService
}

// NewWeatherNotificationHandler は新しいWeatherNotificationHandlerを作成する
func NewWeatherNotificationHandler() *WeatherNotificationHandler {
	return &WeatherNotificationHandler{
		service: usecases.NewWeatherNotificatorService(),
	}
}

// NewWeatherNotificationHandlerWithService はサービスを注入したWeatherNotificationHandlerを作成する
func NewWeatherNotificationHandlerWithService(service *usecases.WeatherNotificatorService) *WeatherNotificationHandler {
	return &WeatherNotificationHandler{
		service: service,
	}
}

// validateRequest はリクエストのバリデーションを行う
func (h *WeatherNotificationHandler) validateRequest(req *WeatherNotificationRequest) error {
	if req.APIKey == "" {
		return fmt.Errorf("API キーが指定されていません")
	}

	if req.City == "" {
		return fmt.Errorf("都市名が指定されていません")
	}

	if req.MaxDays <= 0 {
		return fmt.Errorf("最大日数は1以上である必要があります")
	}

	if req.MaxDays > 5 {
		return fmt.Errorf("最大日数は5日以下である必要があります（OpenWeather API制限）")
	}

	if req.WebhookURL == "" {
		return fmt.Errorf("Discord Webhook URLが指定されていません")
	}

	return nil
}

// sendErrorResponse はエラーレスポンスを送信する
func (h *WeatherNotificationHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Success: false,
		Error:   message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("エラーレスポンスの送信に失敗しました: %v", err)
	}
}

// sendSuccessResponse は成功レスポンスを送信する
func (h *WeatherNotificationHandler) sendSuccessResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := WeatherNotificationResponse{
		Success: true,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("成功レスポンスの送信に失敗しました: %v", err)
	}
}

// HandleWeatherNotification は天気通知のHTTPハンドラメソッド
func (h *WeatherNotificationHandler) HandleWeatherNotification(w http.ResponseWriter, r *http.Request) {
	// HTTPメソッドのチェック
	if r.Method != http.MethodPost {
		h.sendErrorResponse(w, http.StatusMethodNotAllowed, "POSTメソッドのみサポートしています")
		return
	}

	// Content-Typeのチェック
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		h.sendErrorResponse(w, http.StatusBadRequest, "Content-Typeはapplication/jsonである必要があります")
		return
	}

	// リクエストボディの解析
	var req WeatherNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("リクエストボディの解析に失敗しました: %v", err)
		h.sendErrorResponse(w, http.StatusBadRequest, "リクエストボディの形式が正しくありません")
		return
	}

	// リクエストのバリデーション
	if err := h.validateRequest(&req); err != nil {
		log.Printf("リクエストのバリデーションに失敗しました: %v", err)
		h.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// 天気通知の実行
	err := h.service.HandleWeatherNotification(req.APIKey, req.City, req.MaxDays, req.WebhookURL)
	if err != nil {
		log.Printf("天気通知の実行に失敗しました: %v", err)
		h.sendErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("天気通知の実行に失敗しました: %v", err))
		return
	}

	// 成功レスポンスの送信
	successMessage := fmt.Sprintf("✅ %sの%d日間天気予報をDiscordに送信しました", req.City, req.MaxDays)
	log.Printf("%s", successMessage)
	h.sendSuccessResponse(w, successMessage)
}
