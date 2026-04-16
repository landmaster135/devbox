package usecases

import (
	"github.com/landmaster135/devbox/internal/mcp_infra_mark3labs/config"
	"github.com/landmaster135/devbox/internal/mcp_infra_mark3labs/infrastructures/mark3labs"
)

type Request = config.Request
type CallToolParams = config.CallToolParams
type CallToolRequest = config.CallToolRequest
type CallToolResult = config.CallToolResult

type GetPromptRequest = config.GetPromptRequest
type GetPromptResult = config.GetPromptResult
type PromptMessage = config.PromptMessage
type MCPServer = config.MCPServer
type TextContent = config.TextContent

const (
	RoleAssistant = config.RoleAssistant
)

type Service struct {
	adapter mark3labs.Adapter
}

func NewService(adapter mark3labs.Adapter) *Service {
	resolvedAdapter := adapter
	if resolvedAdapter == nil {
		resolvedAdapter = mark3labs.NewAdapter()
	}

	return &Service{
		adapter: resolvedAdapter,
	}
}

func (s *Service) Mark3labs() mark3labs.Adapter {
	return s.adapter
}
