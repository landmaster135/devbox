package open_weather_map

import (
	"context"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	config "github.com/landmaster135/devbox/internal/open_weather_map/config"
	usecases "github.com/landmaster135/devbox/internal/open_weather_map/usecases"
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

	return config.NewConfig(apiKey, city, maxDays)
}

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#

// handleGetWeatherForecast は天気予報を取得するハンドラー
func handleGetWeatherForecast(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// WeatherServiceを初期化
	service := usecases.NewWeatherService()

	// 天気予報を取得
	result, err := service.HandleWeatherForecast(cfg.APIKey, cfg.City, cfg.MaxDays)
	if err != nil {
		return nil, fmt.Errorf("天気予報の取得に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#

// setOpenWeatherMapServer は天気予報取得ツールを提供するMCPサーバを設定します
func setOpenWeatherMapServer(s *server.MCPServer) *server.MCPServer {
	tool := mcp.NewTool(
		"get_weather_forecast",
		mcp.WithDescription("Get weather forecast for specified city using OpenWeather API"),
		mcp.WithString(
			"city",
			mcp.Required(),
			mcp.Description("City name (e.g., Tokyo,JP, London,UK, New York,US)"),
		),
		mcp.WithNumber(
			"max_days",
			mcp.Required(),
			mcp.Description("Maximum number of days for forecast (1-5 days)"),
		),
	)
	s.AddTool(tool, handleGetWeatherForecast)

	return s
}

// addPromptIntoServer はプロンプトをサーバーに追加します
func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt(
		"openweather_forecast_prompt",
		mcp.WithPromptDescription("This is an OpenWeather API forecast prompt"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for OpenWeather API forecast.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You can use this tool to get weather forecasts for any city using the OpenWeather API. The tool requires OPENWEATHER_API_KEY environment variable to be set. Specify city names in the format 'City,CountryCode' (e.g., Tokyo,JP, London,UK)."),
				},
			},
		}, nil
	})
	return s
}

// createOpenWeatherMapServer はOpenWeather API MCPサーバーを作成します
func createOpenWeatherMapServer() *server.MCPServer {
	s := server.NewMCPServer(
		"OpenWeather Map",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setOpenWeatherMapServer(s)
	s = addPromptIntoServer(s)
	return s
}

// BuildOpenWeatherMapServer はOpenWeather API MCPサーバーを構築して実行します
func BuildOpenWeatherMapServer() {
	s := createOpenWeatherMapServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
