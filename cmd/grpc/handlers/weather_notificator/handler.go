package weather_notificator

import (
	"context"
	"fmt"
	"log"

	pb "github.com/landmaster135/devbox/cmd/grpc/proto/weather_notificator"
	usecases "github.com/landmaster135/devbox/internal/weather_notificator/usecases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WeatherNotificatorHandler はgRPCハンドラーを実装する構造体
type WeatherNotificatorHandler struct {
	pb.UnimplementedWeatherNotificatorServiceServer
	service *usecases.WeatherNotificatorService
}

// NewWeatherNotificatorHandler は新しいWeatherNotificatorHandlerを作成する
func NewWeatherNotificatorHandler() *WeatherNotificatorHandler {
	return &WeatherNotificatorHandler{
		service: usecases.NewWeatherNotificatorService(),
	}
}

// NewWeatherNotificatorHandlerWithService は依存関係を注入したWeatherNotificatorHandlerを作成する
func NewWeatherNotificatorHandlerWithService(service *usecases.WeatherNotificatorService) *WeatherNotificatorHandler {
	return &WeatherNotificatorHandler{
		service: service,
	}
}

// SendWeatherNotification は天気予報をDiscordに通知するgRPCメソッド
func (h *WeatherNotificatorHandler) SendWeatherNotification(ctx context.Context, req *pb.WeatherNotificationRequest) (*pb.WeatherNotificationResponse, error) {
	// リクエストのバリデーション
	if err := h.validateRequest(req); err != nil {
		log.Printf("リクエストバリデーションエラー: %v", err)
		return &pb.WeatherNotificationResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	// 天気通知サービスを呼び出し
	err := h.service.SendWeatherNotification(
		ctx,
		req.GetApiKey(),
		req.GetCity(),
		int(req.GetMaxDays()),
		req.GetWebhookUrl(),
	)

	if err != nil {
		log.Printf("天気通知処理エラー: %v", err)
		return &pb.WeatherNotificationResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	// 成功レスポンス
	successMessage := fmt.Sprintf("✅ %sの%d日間天気予報をDiscordに送信しました", req.GetCity(), req.GetMaxDays())
	log.Printf("天気通知成功: %s", successMessage)

	return &pb.WeatherNotificationResponse{
		Success: true,
		Message: successMessage,
	}, nil
}

// validateRequest はリクエストのバリデーションを行う
func (h *WeatherNotificatorHandler) validateRequest(req *pb.WeatherNotificationRequest) error {
	if req.GetApiKey() == "" {
		return fmt.Errorf("APIキーが指定されていません")
	}

	if req.GetCity() == "" {
		return fmt.Errorf("都市名が指定されていません")
	}

	if req.GetMaxDays() <= 0 || req.GetMaxDays() > 5 {
		return fmt.Errorf("最大日数は1-5の範囲で指定してください")
	}

	if req.GetWebhookUrl() == "" {
		return fmt.Errorf("discord Webhook URLが指定されていません")
	}

	return nil
}

// GetService はサービスインスタンスを取得する（テスト用）
func (h *WeatherNotificatorHandler) GetService() *usecases.WeatherNotificatorService {
	return h.service
}
