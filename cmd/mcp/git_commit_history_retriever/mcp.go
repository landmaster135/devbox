package git_commit_history_retriever

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/landmaster135/devbox/internal/git_commit_history_retriever/config"
	"github.com/landmaster135/devbox/internal/git_commit_history_retriever/usecases"
)

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("git_commit_history_prompt",
		mcp.WithPromptDescription("This is a prompt for the Git commit history retriever."),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for the Git commit history retriever.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this great Git commit history retriever well."),
				},
			},
		}, nil
	})
	return s
}

// createConfigFromRequest はMCPリクエストからConfigを作成する
func createConfigFromRequest(request mcp.CallToolRequest) (*config.Config, error) {
	gitDir, err := request.RequireString("git_dir")
	if err != nil {
		return nil, err
	}

	keyword := request.GetString("keyword", "")
	since := request.GetString("since", "")
	until := request.GetString("until", "")

	c, err := config.NewConfig(gitDir, keyword, since, until); if err != nil{
		return nil, fmt.Errorf("MCPリクエストからConfigを作成するのに失敗しました: %v", err)
	}
	return c, nil
}

func BuildMcpServer() {
	s := server.NewMCPServer(
		"Git Commit History Retriever",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	// get_git_commit_history ツール
	historyTool := mcp.NewTool("get_git_commit_history",
		mcp.WithDescription("Get Git commit history from the specified directory"),
		mcp.WithString("git_dir",
			mcp.Required(),
			mcp.Description("Absolute path to the target Git directory"),
		),
		mcp.WithString("keyword",
			mcp.Description("Search keyword (optional)"),
		),
		mcp.WithString("since",
			mcp.Description("Start date (optional, YYYY-MM-DD format)"),
		),
		mcp.WithString("until",
			mcp.Description("End date (optional, YYYY-MM-DD format)"),
		),
	)

	s.AddTool(historyTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cfg, err := createConfigFromRequest(request)
		if err != nil {
			return nil, err
		}

		// GitCommitHistoryServiceを作成
		service := usecases.NewGitCommitHistoryService(cfg.GitDir, cfg)

		// コミット履歴を取得
		history, err := service.GetCommitHistory()
		if err != nil {
			return nil, fmt.Errorf("Gitコミット履歴の取得に失敗しました: %v", err)
		}

		return mcp.NewToolResultText(history), nil
	})

	// get_git_commit_history_with_details ツール
	detailsTool := mcp.NewTool("get_git_commit_history_with_details",
		mcp.WithDescription("Get Git commit history with detailed information from the specified directory"),
		mcp.WithString("git_dir",
			mcp.Required(),
			mcp.Description("Absolute path to the target Git directory"),
		),
		mcp.WithString("keyword",
			mcp.Description("Search keyword (optional)"),
		),
		mcp.WithString("since",
			mcp.Description("Start date (optional, YYYY-MM-DD format)"),
		),
		mcp.WithString("until",
			mcp.Description("End date (optional, YYYY-MM-DD format)"),
		),
	)

	s.AddTool(detailsTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		cfg, err := createConfigFromRequest(request)
		if err != nil {
			return nil, err
		}

		// GitCommitHistoryServiceを作成
		service := usecases.NewGitCommitHistoryService(cfg.GitDir, cfg)

		// コミット履歴と詳細を取得
		historyWithDetails, err := service.GetCommitHistoryWithDetails()
		if err != nil {
			return nil, fmt.Errorf("Gitコミット履歴と詳細の取得に失敗しました: %v", err)
		}

		return mcp.NewToolResultText(historyWithDetails), nil
	})

	// プロンプト
	s = addPromptIntoServer(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
