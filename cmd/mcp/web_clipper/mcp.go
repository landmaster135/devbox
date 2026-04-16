package web_clipper

import (
	"context"
	"fmt"

	mcpInfraUsecases "github.com/landmaster135/devbox/internal/mcp_infra_mark3labs/usecases"
	webClipperConfig "github.com/landmaster135/devbox/internal/web_clipper/config"
	webClipperUsecases "github.com/landmaster135/devbox/internal/web_clipper/usecases"
)

var mcpInfra = mcpInfraUsecases.NewService(nil).Mark3labs()

func createConfigFromRequest(request mcpInfraUsecases.CallToolRequest) (*webClipperConfig.Config, error) {
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

	cfg, err := webClipperConfig.NewConfig(
		webClipperConfig.OperationPatchMarkdown,
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

func handlePatchMarkdown(ctx context.Context, request mcpInfraUsecases.CallToolRequest) (*mcpInfraUsecases.CallToolResult, error) {
	cfg, err := createConfigFromRequest(request)
	if err != nil {
		return nil, err
	}

	service := webClipperUsecases.NewService(nil)
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

	return mcpInfra.NewToolResultText(result), nil
}

func addPromptIntoServer(s *mcpInfraUsecases.MCPServer) *mcpInfraUsecases.MCPServer {
	prompt := mcpInfra.NewPrompt(
		"web_clipper_prompt",
		mcpInfra.WithPromptDescription("Prompt for the web clipper tool."),
	)
	s = mcpInfra.AddPrompt(
		s,
		prompt,
		func(ctx context.Context, request mcpInfraUsecases.GetPromptRequest) (*mcpInfraUsecases.GetPromptResult, error) {
			return &mcpInfraUsecases.GetPromptResult{
				Description: "System prompt for web clipping markdown patch.",
				Messages: []mcpInfraUsecases.PromptMessage{
					{
						Role: mcpInfraUsecases.RoleAssistant,
						Content: mcpInfra.NewTextContent(
							"You can patch markdown by inserting a web article link under a specific heading level.",
						),
					},
				},
			}, nil
		},
	)
	return s
}

func BuildWebClipperServer() {
	s := createWebClipperServer()
	if err := mcpInfra.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func createWebClipperServer() *mcpInfraUsecases.MCPServer {
	s := mcpInfra.NewMCPServer(
		"Web Clipper",
		"1.0.0",
		mcpInfra.WithResourceCapabilities(true, true),
		mcpInfra.WithPromptCapabilities(true),
		mcpInfra.WithLogging(),
	)

	s = setWebClipperTool(s)
	s = addPromptIntoServer(s)
	return s
}

func setWebClipperTool(s *mcpInfraUsecases.MCPServer) *mcpInfraUsecases.MCPServer {
	tool := mcpInfra.NewTool(
		"patch_markdown",
		mcpInfra.WithDescription("Patch markdown with a web article link for summary documents"),
		mcpInfra.WithString(
			"target_title",
			mcpInfra.Required(),
			mcpInfra.Description("The article title for the markdown link"),
		),
		mcpInfra.WithString(
			"target_url",
			mcpInfra.Required(),
			mcpInfra.Description("The article URL for the markdown link"),
		),
		mcpInfra.WithString(
			"src_markdown_content",
			mcpInfra.Description("Source markdown content. Cannot be specified with src_markdown_file"),
		),
		mcpInfra.WithString(
			"src_markdown_file",
			mcpInfra.Description("Source markdown file path. Cannot be specified with src_markdown_content"),
		),
		mcpInfra.WithString(
			"out_file_path",
			mcpInfra.Required(),
			mcpInfra.Description("Output markdown file path"),
		),
		mcpInfra.WithNumber(
			"top_heading_level",
			mcpInfra.Required(),
			mcpInfra.Description("Heading level (>=1) where the link line is inserted below the first heading"),
		),
	)

	return mcpInfra.AddTool(s, tool, handlePatchMarkdown)
}
