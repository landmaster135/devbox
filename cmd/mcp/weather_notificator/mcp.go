package weather_notificator

import (
	"context"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	config "github.com/landmaster135/devbox/internal/weather_notificator/config"
	usecases "github.com/landmaster135/devbox/internal/weather_notificator/usecases"
)

// #==============================================================#
// ##          Environment Provider Interface                   ##
// #==============================================================#

// EnvironmentProvider は環境変数を取得するためのインターフェース
type EnvironmentProvider interface {
	GetEnv(key string) string
}

// StandardEnvironmentProvider は標準の環境変数プロバイダー
type StandardEnvironmentProvider struct{}

// GetEnv は環境変数を取得する
func (p *StandardEnvironmentProvider) GetEnv(key string) string {
	return os.Getenv(key)
}

// #==============================================================#
// ##          Helper Functions                                 ##
// #==============================================================#

// createConfigFromEnv は環境変数とパラメータから設定を作成する
func createConfigFromEnv(envProvider EnvironmentProvider, city string, maxDays int) (*config.Config, error) {
	apiKey := envProvider.GetEnv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("環境変数OPENWEATHER_API_KEYが設定されていません")
	}

	webhookURL := envProvider.GetEnv("DISCORD_WEBHOOK_URL")
	if webhookURL == "" {
		return nil, fmt.Errorf("環境変数DISCORD_WEBHOOK_URLが設定されていません")
	}

	return config.NewConfig(apiKey, city, maxDays, webhookURL)
}

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#

// handleSendWeatherNotification は天気通知を送信するハンドラー
func handleSendWeatherNotification(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// パラメータを取得
	city, err := request.RequireString("city")
	if err != nil {
		return nil, err
	}

	maxDaysInt, err := request.RequireInt("max_days")
	if err != nil {
		return nil, err
	}

	// 環境変数プロバイダーを作成
	envProvider := &StandardEnvironmentProvider{}

	// 設定を作成
	cfg, err := createConfigFromEnv(envProvider, city, maxDaysInt)
	if err != nil {
		return nil, fmt.Errorf("設定の作成に失敗しました: %v", err)
	}

	// WeatherNotificatorServiceを初期化
	service := usecases.NewWeatherNotificatorService()

	// 天気通知を実行
	err = service.HandleWeatherNotification(cfg.APIKey, cfg.City, cfg.MaxDays, cfg.WebhookURL)
	if err != nil {
		return nil, fmt.Errorf("天気通知の送信に失敗しました: %v", err)
	}

	// 成功メッセージを作成
	result := fmt.Sprintf("✅ %sの%d日間天気予報をDiscordに送信しました", cfg.City, cfg.MaxDays)

	return mcp.NewToolResultText(result), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#

// setWeatherNotificatorServer は天気通知ツールを提供するMCPサーバを設定します
func setWeatherNotificatorServer(s *server.MCPServer) *server.MCPServer {
	tool := mcp.NewTool(
		"send_weather_notification",
		mcp.WithDescription("Send weather forecast notification to Discord for specified city"),
		mcp.WithString(
			"city",
			mcp.Required(),
			mcp.Description("City name (e.g., Tokyo, Osaka, New York)"),
		),
		mcp.WithNumber(
			"max_days",
			mcp.Required(),
			mcp.Description("Maximum number of days for forecast (1-5 days)"),
		),
	)
	s.AddTool(tool, handleSendWeatherNotification)

	return s
}

// addPromptIntoServer はプロンプトをサーバーに追加します
func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt(
		"weather_notification_prompt",
		mcp.WithPromptDescription("This is a weather notification prompt"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for weather notification.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You can use this tool to send weather forecast notifications to Discord. The tool requires OPENWEATHER_API_KEY and DISCORD_WEBHOOK_URL environment variables to be set."),
				},
			},
		}, nil
	})
	return s
}

// createWeatherNotificatorServer は天気通知MCPサーバーを作成します
func createWeatherNotificatorServer() *server.MCPServer {
	s := server.NewMCPServer(
		"Weather Notificator",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setWeatherNotificatorServer(s)
	s = addPromptIntoServer(s)
	return s
}

// BuildWeatherNotificatorServer は天気通知MCPサーバーを構築して実行します
func BuildWeatherNotificatorServer() {
	s := createWeatherNotificatorServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
