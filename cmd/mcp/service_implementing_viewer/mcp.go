package service_implementing_viewer

import (
	"context"
	"fmt"
	"strings"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/service_implementing_viewer/usecases"
)

func handleGetServiceImplementingStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	rootDir, err := request.RequireString("root_dir")
	if err != nil {
		return nil, err
	}

	targetDirsStr, err := request.RequireString("target_dirs")
	if err != nil {
		return nil, err
	}

	// target_dirsをカンマ区切りで分割
	targetDirs := parseTargetDirs(targetDirsStr)
	if len(targetDirs) == 0 {
		return nil, fmt.Errorf("target_dirsが空です")
	}

	// ServiceImplementingViewerServiceを初期化
	service := usecases.NewServiceImplementingViewerService(rootDir, targetDirs)
	result, _, err := service.GetServiceImplementingStatus()
	if err != nil {
		return nil, fmt.Errorf("サービス実装状況の取得に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

// parseTargetDirs はカンマ区切りの文字列をスライスに変換する
func parseTargetDirs(targetDirs string) []string {
	if targetDirs == "" {
		return []string{}
	}

	dirs := strings.Split(targetDirs, ",")
	result := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		trimmed := strings.TrimSpace(dir)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a service implementing viewer prompt"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for service implementing viewer.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this service implementing viewer well."),
				},
			},
		}, nil
	})
	return s
}

func BuildServiceImplementingViewerServer() {
	s := server.NewMCPServer(
		"Service Implementing Viewer",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	getServiceImplementingStatusTool := mcp.NewTool("get_service_implementing_status",
		mcp.WithDescription("Get service implementation status across different directories"),
		mcp.WithString("root_dir",
			mcp.Required(),
			mcp.Description("Absolute path to the root directory"),
		),
		mcp.WithString("target_dirs",
			mcp.Required(),
			mcp.Description("Target directories to check (comma-separated)"),
		),
	)
	s.AddTool(getServiceImplementingStatusTool, handleGetServiceImplementingStatus)

	// プロンプト
	s = addPromptIntoServer(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
