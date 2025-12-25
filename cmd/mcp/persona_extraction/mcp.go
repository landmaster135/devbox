package persona_extraction

import (
	"context"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/landmaster135/devbox/internal/persona_extraction/domain"
	usecases "github.com/landmaster135/devbox/internal/persona_extraction/usecases"
)

const toolVersion = "0.1.0"

func handlePersonaExtraction(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	svc := usecases.NewPersonaService()
	result, err := svc.HandleExtraction(request.GetArguments())
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(result), nil
}

func createPersonaTool() mcp.Tool {
	return mcp.NewTool(
		"extract_personas",
		mcp.WithDescription("キャラクターごとのペルソナ情報を構造化し、要約を生成します"),
		mcp.WithString(
			"context",
			mcp.Description("シナリオや作品の背景メモ"),
		),
		mcp.WithArray(
			"characters",
			mcp.Required(),
			mcp.Description("1人以上のキャラクター設定を含む配列"),
			mcp.Items(domain.PersonaJSONSchema()),
		),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
	)
}

func BuildPersonaExtractionServer() {
	s := server.NewMCPServer(
		"Persona Extraction",
		toolVersion,
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithLogging(),
	)

	s.AddTool(createPersonaTool(), handlePersonaExtraction)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "persona extraction server error: %v\n", err)
	}
}
