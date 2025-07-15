package brave_search

import (
	"context"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/brave_search/usecases"
)

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#
func handleWebSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 引数の取得
	query, err := request.RequireString("query")
	if err != nil {
		return nil, err
	}

	// オプションパラメータの取得
	count := request.GetInt("count", 10)
	offset := request.GetInt("offset", 0)

	// BraveSearchServiceを初期化
	service := usecases.NewBraveSearchService()
	result, err := service.HandleWebSearch(query, count, offset)
	if err != nil {
		return nil, fmt.Errorf("web検索に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

func handleLocalSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 引数の取得
	query, err := request.RequireString("query")
	if err != nil {
		return nil, err
	}

	// オプションパラメータの取得
	count := request.GetInt("count", 5)

	// BraveSearchServiceを初期化
	service := usecases.NewBraveSearchService()
	result, err := service.HandleLocalSearch(query, count)
	if err != nil {
		return nil, fmt.Errorf("ローカル検索に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#
func setBraveWebSearchServer(s *server.MCPServer) *server.MCPServer {
	// Web検索ツールの定義
	webSearchTool := mcp.NewTool("brave_web_search",
		mcp.WithDescription("Performs a web search using the Brave Search API, ideal for general queries, news, articles, and online content. "+
			"Use this for broad information gathering, recent events, or when you need diverse web sources. "+
			"Supports pagination, content filtering, and freshness controls. "+
			"Maximum 20 results per request, with offset for pagination."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query (max 400 chars, 50 words)"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of results (1-20, default 10)"),
		),
		mcp.WithNumber("offset",
			mcp.Description("Pagination offset (max 9, default 0)"),
		),
	)
	s.AddTool(webSearchTool, handleWebSearch)

	return s
}

func setBraveLocalSearchServer(s *server.MCPServer) *server.MCPServer {
	// ローカル検索ツールの定義
	localSearchTool := mcp.NewTool("brave_local_search",
		mcp.WithDescription("Searches for local businesses and places using Brave's Local Search API. "+
			"Best for queries related to physical locations, businesses, restaurants, services, etc. "+
			"Returns detailed information including:\n"+
			"- Business names and addresses\n"+
			"- Ratings and review counts\n"+
			"- Phone numbers and opening hours\n"+
			"Use this when the query implies 'near me' or mentions specific locations. "+
			"Automatically falls back to web search if no local results are found."),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Local search query (e.g. 'pizza near Central Park')"),
		),
		mcp.WithNumber("count",
			mcp.Description("Number of results (1-20, default 5)"),
		),
	)
	s.AddTool(localSearchTool, handleLocalSearch)

	return s
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a search engine prompt"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for search engine.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this search engine well."),
				},
			},
		}, nil
	})
	return s
}

func createBraveSearchServer() *server.MCPServer {
	s := server.NewMCPServer(
		"Brave Search",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setBraveWebSearchServer(s)
	s = setBraveLocalSearchServer(s)
	s = addPromptIntoServer(s)
	return s
}

// MCPサーバを構築する関数
func BuildBraveSearchServer() {
	// APIキーのチェック
	apiKey := os.Getenv("BRAVE_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: BRAVE_API_KEY environment variable is required")
		os.Exit(1)
	}

	s := createBraveSearchServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
