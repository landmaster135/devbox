package context7

import (
	"context"
	"encoding/json"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/landmaster135/devbox/internal/context7/domain/models"
	"github.com/landmaster135/devbox/internal/context7/usecases"
)

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#

// handleSearchLibrary はライブラリ名を検索してContext7互換IDを取得します
func handleSearchLibrary(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	libraryName, err := request.RequireString("library_name")
	if err != nil {
		return nil, err
	}

	// Context7Serviceを初期化
	service := usecases.NewContext7ServiceWithHTTPClient()

	// ライブラリを検索
	searchResponse, err := service.ResolveLibraryID(libraryName)
	if err != nil {
		return nil, fmt.Errorf("ライブラリ検索に失敗しました: %v", err)
	}

	// 結果をJSON形式で返す
	jsonResult, err := json.Marshal(searchResponse)
	if err != nil {
		return nil, fmt.Errorf("検索結果のJSON変換に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// handleGetLibraryDocs はライブラリIDからドキュメントを取得します
func handleGetLibraryDocs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	libraryID, err := request.RequireString("library_id")
	if err != nil {
		return nil, err
	}

	// オプションパラメータを取得
	topic := request.GetString("topic", "")
	tokens := request.GetInt("tokens", models.DefaultTokens)

	// Context7Serviceを初期化
	service := usecases.NewContext7ServiceWithHTTPClient()

	// ライブラリIDの形式を検証
	if err := service.ValidateLibraryID(libraryID); err != nil {
		return nil, fmt.Errorf("ライブラリID検証エラー: %v", err)
	}

	// ドキュメントオプションを設定
	options := models.DocOptions{
		Topic:  topic,
		Tokens: int(tokens),
	}

	// デフォルトトークン数を設定
	if options.Tokens <= 0 {
		options.Tokens = models.DefaultTokens
	}

	// ドキュメントを取得
	docs, err := service.GetLibraryDocs(libraryID, options)
	if err != nil {
		return nil, fmt.Errorf("ドキュメント取得に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(docs), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#

// setContext7Server はContext7関連のツールをサーバーに追加します
func setContext7Server(s *server.MCPServer) *server.MCPServer {
	// ライブラリ検索ツール
	searchTool := mcp.NewTool(
		"search_library",
		mcp.WithDescription("Context7でライブラリを検索してContext7互換IDを取得します"),
		mcp.WithString("library_name",
			mcp.Required(),
			mcp.Description("検索するライブラリ名（例: react, next.js）"),
		),
	)
	s.AddTool(searchTool, handleSearchLibrary)

	// ドキュメント取得ツール
	docsTool := mcp.NewTool(
		"get_library_docs",
		mcp.WithDescription("Context7互換ライブラリIDからドキュメントを取得します"),
		mcp.WithString("library_id",
			mcp.Required(),
			mcp.Description("Context7互換ライブラリID（例: /facebook/react）"),
		),
		mcp.WithString("topic",
			mcp.Description("特定のトピックに焦点を当てる（例: hooks, routing）"),
		),
		mcp.WithNumber("tokens",
			mcp.Description("取得する最大トークン数（デフォルト: 10000）"),
		),
	)
	s.AddTool(docsTool, handleGetLibraryDocs)

	return s
}

// addPromptIntoServer はContext7用のプロンプトをサーバーに追加します
func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt(
		"context7_system_prompt",
		mcp.WithPromptDescription("Context7ライブラリドキュメント取得用のシステムプロンプト"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "Context7を使用してライブラリドキュメントを効果的に取得するためのシステムプロンプト",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("Context7を使用して最新のライブラリドキュメントを取得し、開発者に有用な情報を提供します。まずライブラリを検索してIDを取得し、その後ドキュメントを取得してください。"),
				},
			},
		}, nil
	})
	return s
}

// createContext7Server はContext7用のMCPサーバーを作成します
func createContext7Server() *server.MCPServer {
	s := server.NewMCPServer(
		"Context7 Library Documentation",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setContext7Server(s)
	s = addPromptIntoServer(s)
	return s
}

// BuildContext7Server はContext7用のMCPサーバーを起動します
func BuildContext7Server() {
	s := createContext7Server()
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
