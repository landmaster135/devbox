package github

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/github/usecases"
)

type GithubHandler struct {
	Token string
}

func NewGithubHandler(token string) *GithubHandler {
	return &GithubHandler{
		Token: token,
	}
}

// handleToSearchCode はGitHub全体でコードを検索して、結果をJSON形式で返します
func (h *GithubHandler) handleToSearchCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return nil, err
	}

	// 数値オプションパラメータを取得
	page := request.GetInt("page", 1)
	perPage := request.GetInt("per_page", 30)

	// GitHubSearchServiceを初期化
	service, err := usecases.NewGitHubSearchService(h.Token)
	if err != nil {
		return nil, fmt.Errorf("GitHubSearchServiceの初期化に失敗しました: %v", err)
	}
	result, err := service.HandleToSearchCode(query, page, perPage)
	if err != nil {
		return nil, fmt.Errorf("GitHub検索に失敗しました: %v", err)
	}

	return returnJSONResult(result)
}

// SetGitHubGlobalServer は受け取ったMCPサーバにGitHubグローバル検索用のツールを付与して、そのMCPサーバを返します。
func SetGitHubGlobalServer(token string, s *server.MCPServer) *server.MCPServer {
	// ツール1: コード検索
	searchCodeTool := mcp.NewTool("search_code",
		mcp.WithDescription("Search for code across GitHub repositories"),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("Search query. For example: 'addClass in:file language:js repo:jquery/jquery'"),
		),
		mcp.WithNumber("page",
			mcp.Description("Page number (default: 1)"),
		),
		mcp.WithNumber("per_page",
			mcp.Description("Results per page (default: 30, max: 100)"),
		),
	)

	// トークンを使用したハンドラーを作成
	h := NewGithubHandler(token)
	s.AddTool(searchCodeTool, h.handleToSearchCode)

	return s
}
