package notion_sync

import (
	"context"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	config "github.com/landmaster135/devbox/internal/notion_sync/config"
	usecases "github.com/landmaster135/devbox/internal/notion_sync/usecases"
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
func createConfigFromEnv(envProvider EnvironmentProvider, operation, conID, pageID, markdownContent string, toggleH1, toggleH2, toggleH3 bool) (*config.Config, error) {
	token := envProvider.GetEnv("NOTION_INTEGRATION_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("環境変数NOTION_INTEGRATION_TOKENが設定されていません")
	}

	endpointURL := envProvider.GetEnv("NOTION_ENDPOINT_URL")
	if endpointURL == "" {
		return nil, fmt.Errorf("環境変数NOTION_ENDPOINT_URLが設定されていません")
	}

	return config.NewConfig(operation, token, conID, pageID, markdownContent, endpointURL, toggleH1, toggleH2, toggleH3)
}

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#

// handlePatchPage はNotionページにMarkdownコンテンツをパッチするハンドラー
func handlePatchPage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// パラメータを取得
	markdownContent, err := request.RequireString("markdown_content")
	if err != nil {
		return nil, err
	}

	// オプションパラメータを取得
	conID := request.GetString("con_id", "")
	pageID := request.GetString("page_id", "")
	toggleH1 := request.GetBool("toggle_h1", false)
	toggleH2 := request.GetBool("toggle_h2", false)
	toggleH3 := request.GetBool("toggle_h3", false)

	// 環境変数プロバイダーを作成
	envProvider := &StandardEnvironmentProvider{}

	// patch_pageツールでは常にpatch操作を使用
	operation := "patch"

	// 設定を作成
	cfg, err := createConfigFromEnv(envProvider, operation, conID, pageID, markdownContent, toggleH1, toggleH2, toggleH3)
	if err != nil {
		return nil, fmt.Errorf("設定の作成に失敗しました: %v", err)
	}

	// NotionSyncServiceを初期化
	service := usecases.NewNotionSyncService()

	// Notion同期を実行
	result, err := service.HandleNotionSync(cfg)
	if err != nil {
		return nil, fmt.Errorf("notion同期に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#

// setNotionSyncServer はNotion同期ツールを提供するMCPサーバを設定します
func setNotionSyncServer(s *server.MCPServer) *server.MCPServer {
	tool := mcp.NewTool(
		"patch_page",
		mcp.WithDescription("Patch Markdown content into Notion page"),
		mcp.WithString(
			"con_id",
			mcp.Description("con_id (excludes page_id)"),
		),
		mcp.WithString(
			"page_id",
			mcp.Description("page_id (excludes con_id)"),
		),
		mcp.WithString(
			"markdown_content",
			mcp.Required(),
			mcp.Description("Markdown content"),
		),
		mcp.WithBoolean(
			"toggle_h1",
			mcp.Description("Toggles heading1 or not"),
		),
		mcp.WithBoolean(
			"toggle_h2",
			mcp.Description("Toggles heading2 or not"),
		),
		mcp.WithBoolean(
			"toggle_h3",
			mcp.Description("Toggles heading3 or not"),
		),
	)
	s.AddTool(tool, handlePatchPage)

	return s
}

// addPromptIntoServer はプロンプトをサーバーに追加します
func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt(
		"notion_sync_prompt",
		mcp.WithPromptDescription("This is a Notion synchronization prompt"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for Notion synchronization.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You can use this tool to synchronize Markdown content with Notion pages."),
				},
			},
		}, nil
	})
	return s
}

// createNotionSyncServer はNotion同期MCPサーバーを作成します
func createNotionSyncServer() *server.MCPServer {
	s := server.NewMCPServer(
		"Notion Sync",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setNotionSyncServer(s)
	s = addPromptIntoServer(s)
	return s
}

// BuildNotionSyncServer はNotion同期MCPサーバーを構築して実行します
func BuildNotionSyncServer() {
	s := createNotionSyncServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
