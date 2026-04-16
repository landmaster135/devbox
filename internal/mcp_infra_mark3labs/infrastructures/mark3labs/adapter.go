package mark3labs

import (
	"github.com/landmaster135/devbox/internal/mcp_infra_mark3labs/config"
	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

type Adapter interface {
	NewToolResultText(text string) *config.CallToolResult
	NewTextContent(text string) config.TextContent
	NewPrompt(name string, opts ...config.PromptOption) config.Prompt
	WithPromptDescription(description string) config.PromptOption

	NewTool(name string, opts ...config.ToolOption) config.Tool
	WithDescription(description string) config.ToolOption
	WithString(name string, opts ...config.PropertyOption) config.ToolOption
	WithNumber(name string, opts ...config.PropertyOption) config.ToolOption
	Required() config.PropertyOption
	Description(description string) config.PropertyOption

	NewMCPServer(name, version string, opts ...config.ServerOption) *config.MCPServer
	WithResourceCapabilities(subscribe, listChanged bool) config.ServerOption
	WithPromptCapabilities(listChanged bool) config.ServerOption
	WithLogging() config.ServerOption

	AddPrompt(
		s *config.MCPServer,
		prompt config.Prompt,
		handler config.PromptHandlerFunc,
	) *config.MCPServer
	AddTool(
		s *config.MCPServer,
		tool config.Tool,
		handler config.ToolHandlerFunc,
	) *config.MCPServer

	ServeStdio(s *config.MCPServer) error
}

type Mark3labsAdapter struct {
	serveStdioFunc func(s *config.MCPServer) error
}

var _ Adapter = (*Mark3labsAdapter)(nil)

func NewAdapter() Adapter {
	return NewAdapterWithServeStdioFunc(nil)
}

func NewAdapterWithServeStdioFunc(serveStdioFunc func(s *config.MCPServer) error) Adapter {
	stdioFunc := serveStdioFunc
	if stdioFunc == nil {
		stdioFunc = func(s *config.MCPServer) error {
			return server.ServeStdio(s)
		}
	}

	return &Mark3labsAdapter{
		serveStdioFunc: stdioFunc,
	}
}

func (a *Mark3labsAdapter) NewToolResultText(text string) *config.CallToolResult {
	return mcp.NewToolResultText(text)
}

func (a *Mark3labsAdapter) NewTextContent(text string) config.TextContent {
	return mcp.NewTextContent(text)
}

func (a *Mark3labsAdapter) NewPrompt(name string, opts ...config.PromptOption) config.Prompt {
	return mcp.NewPrompt(name, opts...)
}

func (a *Mark3labsAdapter) WithPromptDescription(description string) config.PromptOption {
	return mcp.WithPromptDescription(description)
}

func (a *Mark3labsAdapter) NewTool(name string, opts ...config.ToolOption) config.Tool {
	return mcp.NewTool(name, opts...)
}

func (a *Mark3labsAdapter) WithDescription(description string) config.ToolOption {
	return mcp.WithDescription(description)
}

func (a *Mark3labsAdapter) WithString(name string, opts ...config.PropertyOption) config.ToolOption {
	return mcp.WithString(name, opts...)
}

func (a *Mark3labsAdapter) WithNumber(name string, opts ...config.PropertyOption) config.ToolOption {
	return mcp.WithNumber(name, opts...)
}

func (a *Mark3labsAdapter) Required() config.PropertyOption {
	return mcp.Required()
}

func (a *Mark3labsAdapter) Description(description string) config.PropertyOption {
	return mcp.Description(description)
}

func (a *Mark3labsAdapter) NewMCPServer(
	name,
	version string,
	opts ...config.ServerOption,
) *config.MCPServer {
	return server.NewMCPServer(name, version, opts...)
}

func (a *Mark3labsAdapter) WithResourceCapabilities(
	subscribe,
	listChanged bool,
) config.ServerOption {
	return server.WithResourceCapabilities(subscribe, listChanged)
}

func (a *Mark3labsAdapter) WithPromptCapabilities(listChanged bool) config.ServerOption {
	return server.WithPromptCapabilities(listChanged)
}

func (a *Mark3labsAdapter) WithLogging() config.ServerOption {
	return server.WithLogging()
}

func (a *Mark3labsAdapter) AddPrompt(
	s *config.MCPServer,
	prompt config.Prompt,
	handler config.PromptHandlerFunc,
) *config.MCPServer {
	s.AddPrompt(prompt, handler)
	return s
}

func (a *Mark3labsAdapter) AddTool(
	s *config.MCPServer,
	tool config.Tool,
	handler config.ToolHandlerFunc,
) *config.MCPServer {
	s.AddTool(tool, handler)
	return s
}

func (a *Mark3labsAdapter) ServeStdio(s *config.MCPServer) error {
	return a.serveStdioFunc(s)
}
