package http_request

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

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

func BuildMcpServer() {
	s := server.NewMCPServer(
		"HTTP requestor",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
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
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		method, err := request.RequireString("method")
		if err != nil {
			return nil, err
		}
		url, err := request.RequireString("url")
		if err != nil {
			return nil, err
		}
		body := request.GetString("body", "")

		// Create and send request
		var req *http.Request
		var reqErr error
		if body != "" {
			req, reqErr = http.NewRequest(method, url, strings.NewReader(body))
		} else {
			req, reqErr = http.NewRequest(method, url, nil)
		}
		if reqErr != nil {
			return nil, fmt.Errorf("failed to create request: %v", reqErr)
		}

		client := &http.Client{}
		resp, respErr := client.Do(req)
		if respErr != nil {
			return nil, fmt.Errorf("request failed: %v", respErr)
		}
		defer resp.Body.Close()

		// Return response
		respBody, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, fmt.Errorf("failed to read response: %v", bodyErr)
		}

		return mcp.NewToolResultText(fmt.Sprintf("Status: %d\nBody: %s", resp.StatusCode, string(respBody))), nil
	})

	// プロンプト
	s = addPromptIntoServer(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
