package mark3labs

import (
	"context"
	"errors"
	"testing"

	"github.com/landmaster135/devbox/internal/mcp_infra_mark3labs/config"
)

func TestMark3labsAdapter_NewToolResultText_Normal(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter()
	result := adapter.NewToolResultText("ok")
	if result == nil {
		t.Fatal("result is nil")
	}
	if len(result.Content) != 1 {
		t.Fatalf("unexpected content length: got=%d want=%d", len(result.Content), 1)
	}

	textContent, ok := result.Content[0].(config.TextContent)
	if !ok {
		t.Fatalf("unexpected content type: %T", result.Content[0])
	}
	if textContent.Text != "ok" {
		t.Fatalf("unexpected text content: got=%q want=%q", textContent.Text, "ok")
	}
}

func TestMark3labsAdapter_AddToolAndPrompt_Normal(t *testing.T) {
	t.Parallel()

	adapter := NewAdapter()
	s := adapter.NewMCPServer(
		"test-server",
		"1.0.0",
		adapter.WithResourceCapabilities(true, true),
		adapter.WithPromptCapabilities(true),
		adapter.WithLogging(),
	)

	tool := adapter.NewTool(
		"echo",
		adapter.WithDescription("echo tool"),
		adapter.WithString(
			"message",
			adapter.Required(),
			adapter.Description("message"),
		),
	)
	s = adapter.AddTool(
		s,
		tool,
		func(ctx context.Context, request config.CallToolRequest) (*config.CallToolResult, error) {
			return adapter.NewToolResultText("echo"), nil
		},
	)

	prompt := adapter.NewPrompt(
		"echo_prompt",
		adapter.WithPromptDescription("echo prompt"),
	)
	s = adapter.AddPrompt(
		s,
		prompt,
		func(ctx context.Context, request config.GetPromptRequest) (*config.GetPromptResult, error) {
			return &config.GetPromptResult{
				Description: "prompt",
				Messages: []config.PromptMessage{
					{
						Role:    config.RoleAssistant,
						Content: adapter.NewTextContent("hello"),
					},
				},
			}, nil
		},
	)
	if s == nil {
		t.Fatal("server is nil")
	}
}

func TestMark3labsAdapter_ServeStdio_DelegatesFunction_Normal(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("stdio error")
	called := false
	adapter := NewAdapterWithServeStdioFunc(func(s *config.MCPServer) error {
		called = true
		if s == nil {
			t.Fatal("server is nil")
		}
		return expectedErr
	})

	s := adapter.NewMCPServer("test-server", "1.0.0")
	err := adapter.ServeStdio(s)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected error: got=%v want=%v", err, expectedErr)
	}
	if !called {
		t.Fatal("serveStdioFunc was not called")
	}
}
