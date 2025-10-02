package git_diff_recorder

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	config "github.com/landmaster135/devbox/internal/git_diff_recorder/config"
	usecases "github.com/landmaster135/devbox/internal/git_diff_recorder/usecases"
)

func handleGetGitDiff(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	gitDir, err := request.RequireString("git_dir")
	if err != nil {
		return nil, err
	}

	stagedOnly := request.GetBool("staged_only", false)

	// 設定を作成
	cfg := &config.Config{
		GitDir:     gitDir,
		StagedOnly: stagedOnly,
	}

	// GitDiffGeneratorServiceを作成
	service := usecases.NewGitDiffGeneratorService(gitDir, cfg)

	// 詳細差分を取得
	detailedDiff, err := service.GetCurrentDetailedDiff()
	if err != nil {
		return nil, fmt.Errorf("git差分の取得に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(detailedDiff), nil
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("git_diff_prompt",
		mcp.WithPromptDescription("This is a prompt for the Git diff recorder."),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for the Git diff recorder.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this great Git diff recorder well."),
				},
			},
		}, nil
	})
	return s
}

func BuildMcpServer() {
	s := server.NewMCPServer(
		"Git Diff Recorder",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	tool := mcp.NewTool("get_git_diff",
		mcp.WithDescription("Get Git diff from the specified directory"),
		mcp.WithString("git_dir",
			mcp.Required(),
			mcp.Description("Absolute path to the target Git directory"),
		),
		mcp.WithBoolean("staged_only",
			mcp.Description("Whether to get only staged diff (default: false)"),
		),
	)

	s.AddTool(tool, handleGetGitDiff)

	// プロンプト
	s = addPromptIntoServer(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
