package duckduckgo_search

import (
	"context"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/duckduckgo_search/usecases"
)

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#

// handleWebSearch はWeb検索のMCPハンドラーです
func handleWebSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 引数の取得
	query, err := request.RequireString("query")
	if err != nil {
		return nil, err
	}

	// オプションパラメータの取得
	count := request.GetInt("count", 10)

	// DuckDuckGoSearchServiceを初期化
	service := usecases.NewDuckDuckGoSearchService()
	result, err := service.HandleWebSearch(query, count)
	if err != nil {
		return nil, fmt.Errorf("web検索に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

// handleInstantSearch はInstant Answer検索のMCPハンドラーです
func handleInstantSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 引数の取得
	query, err := request.RequireString("query")
	if err != nil {
		return nil, err
	}

	// DuckDuckGoSearchServiceを初期化
	service := usecases.NewDuckDuckGoSearchService()
	result, err := service.HandleInstantSearch(query)
	if err != nil {
		return nil, fmt.Errorf("instant answer検索に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a DuckDuckGo search engine prompt"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for DuckDuckGo search engine.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this DuckDuckGo search engine well for privacy-focused search results."),
				},
			},
		}, nil
	})
	return s
}

// MCPサーバを構築する関数
func BuildDuckDuckGoSearchServer() {
	// サーバの作成
	s := server.NewMCPServer(
		"DuckDuckGo Search",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	// Web検索ツールの定義
	webSearchTool := mcp.NewTool("duckduckgo_web_search",
		mcp.WithDescription("Performs a web search using DuckDuckGo search engine, focusing on privacy and unbiased results. "+
			"Use this for general information gathering, news, articles, and online content when you need privacy-focused search results. "+
			"Returns organic search results without tracking or personalization. "+
			"Maximum 20 results per request."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query (any keywords you would use in DuckDuckGo search)"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of results (1-20, default 10)"),
		),
	)

	// Instant Answer検索ツールの定義
	instantSearchTool := mcp.NewTool("duckduckgo_instant_search",
		mcp.WithDescription("Performs an instant answer search using DuckDuckGo's Instant Answer API. "+
			"Best for quick facts, definitions, calculations, and direct answers. "+
			"Returns structured information including abstracts, definitions, and related topics. "+
			"Use this when you need immediate factual information rather than web page results."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query for instant answers (e.g., 'weather Tokyo', 'what is Python', '2+2')"),
		),
	)

	// Web検索ツールのハンドラ
	s.AddTool(webSearchTool, handleWebSearch)

	// Instant Answer検索ツールのハンドラ
	s.AddTool(instantSearchTool, handleInstantSearch)

	// プロンプト
	s = addPromptIntoServer(s)

	// サーバの起動
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
