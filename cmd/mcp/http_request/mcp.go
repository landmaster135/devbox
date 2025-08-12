package http_request

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/landmaster135/devbox/internal/http_request/domain/models"
	"github.com/landmaster135/devbox/internal/http_request/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/http_request/usecases/services"
)

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#
func handleHttpRequest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method, err := request.RequireString("method")
	if err != nil {
		return nil, err
	}

	url, err := request.RequireString("url")
	if err != nil {
		return nil, err
	}

	body := request.GetString("body", "")
	encoding := request.GetString("encoding", "auto")

	// 依存関係の注入
	apiRepo := repositories.NewHTTPRepository()
	apiService := services.NewHTTPService(apiRepo)

	// ヘッダーの準備
	headers := map[string]string{"Accept": "application/json"}

	var response *models.HTTPResponse
	if body != "" {
		// Content-Typeヘッダーを追加
		headers["Content-Type"] = "application/json"

		// リクエストを作成
		apiRequest := &models.HTTPRequest{
			URL:      url,
			Method:   method,
			Headers:  headers,
			Body:     []byte(body),
			Encoding: encoding,
		}

		// リクエストを送信
		response, err = apiRepo.SendRequest(apiRequest)
	} else {
		// JSONファイルなしでリクエストを送信（GETなど）
		response, err = apiService.SendRequestWithoutJSONFile(url, method, headers, encoding)
	}

	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの送信に失敗しました: %v", err)
	}

	// レスポンスを整形
	formattedResponse, err := apiService.FormatResponse(response)
	if err != nil {
		return nil, fmt.Errorf("レスポンスの整形に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(formattedResponse), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#
func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a prompt for the HTTP client."),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for the HTTP client.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this great HTTP client well."),
				},
			},
		}, nil
	})
	return s
}

func setHttpRequestServer(s *server.MCPServer) *server.MCPServer {
	tool := mcp.NewTool("http_request",
		mcp.WithDescription("Make HTTP requests to external APIs"),
		mcp.WithString("method",
			mcp.Required(),
			mcp.Description("HTTP method to use"),
			mcp.Enum("GET", "POST", "PUT", "DELETE"),
		),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to send the request to"),
			mcp.Pattern("^https?://.*"),
		),
		mcp.WithString("body",
			mcp.Description("Request body (for POST/PUT)"),
		),
		mcp.WithString("encoding",
			mcp.Description("Character encoding (shift_jis, utf-8, euc-jp, auto)"),
		),
	)
	s.AddTool(tool, handleHttpRequest)
	return s
}

func createHTTPRequestServer() *server.MCPServer {
	s := server.NewMCPServer(
		"HTTP requestor",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setHttpRequestServer(s)
	s = addPromptIntoServer(s)
	return s
}

func BuildMcpServer() {
	s := createHTTPRequestServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
