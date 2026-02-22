package web_clipper

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/landmaster135/devbox/internal/web_clipper/config"
	"github.com/landmaster135/devbox/internal/web_clipper/usecases"
)

func createConfigFromRequest(request mcp.CallToolRequest) (*config.Config, error) {
	targetTitle, err := request.RequireString("target_title")
	if err != nil {
		return nil, err
	}

	targetURL, err := request.RequireString("target_url")
	if err != nil {
		return nil, err
	}

	outFilePath, err := request.RequireString("out_file_path")
	if err != nil {
		return nil, err
	}

	topHeadingLevel, err := request.RequireInt("top_heading_level")
	if err != nil {
		return nil, err
	}

	srcMarkdownContent := request.GetString("src_markdown_content", "")
	srcMarkdownFile := request.GetString("src_markdown_file", "")

	cfg, err := config.NewConfig(
		config.OperationPatchMarkdown,
		targetTitle,
		targetURL,
		srcMarkdownContent,
		srcMarkdownFile,
		outFilePath,
		topHeadingLevel,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("MCPリクエストの検証に失敗しました: %v", err)
	}

	return cfg, nil
}

func handlePatchMarkdown(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cfg, err := createConfigFromRequest(request)
	if err != nil {
		return nil, err
	}

	service := usecases.NewService(nil)
	result, err := service.PatchMarkdown(
		cfg.TargetTitle,
		cfg.TargetURL,
		cfg.SrcMarkdownContent,
		cfg.SrcMarkdownFile,
		cfg.OutFilePath,
		cfg.TopHeadingLevel,
	)
	if err != nil {
		return nil, fmt.Errorf("MarkdownへのWeb記事リンク挿入に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt(
		"web_clipper_prompt",
		mcp.WithPromptDescription("Prompt for the web clipper tool."),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for web clipping markdown patch.",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleAssistant,
					Content: mcp.NewTextContent(
						"You can patch markdown by inserting a web article link under a specific heading level.",
					),
				},
			},
		}, nil
	})
	return s
}

func BuildWebClipperServer() {
	s := createWebClipperServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func createWebClipperServer() *server.MCPServer {
	s := server.NewMCPServer(
		"Web Clipper",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	s = setWebClipperTool(s)
	s = addPromptIntoServer(s)
	return s
}

func setWebClipperTool(s *server.MCPServer) *server.MCPServer {
	tool := mcp.NewTool(
		"patch_markdown",
		mcp.WithDescription("Patch markdown with a web article link for summary documents"),
		mcp.WithString(
			"target_title",
			mcp.Required(),
			mcp.Description("The article title for the markdown link"),
		),
		mcp.WithString(
			"target_url",
			mcp.Required(),
			mcp.Description("The article URL for the markdown link"),
		),
		mcp.WithString(
			"src_markdown_content",
			mcp.Description("Source markdown content. Cannot be specified with src_markdown_file"),
		),
		mcp.WithString(
			"src_markdown_file",
			mcp.Description("Source markdown file path. Cannot be specified with src_markdown_content"),
		),
		mcp.WithString(
			"out_file_path",
			mcp.Required(),
			mcp.Description("Output markdown file path"),
		),
		mcp.WithNumber(
			"top_heading_level",
			mcp.Required(),
			mcp.Description("Heading level (>=1) where the link line is inserted below the first heading"),
		),
	)

	s.AddTool(tool, handlePatchMarkdown)
	return s
}
