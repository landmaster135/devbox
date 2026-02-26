package config

import (
	"context"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

type Request = mcp.Request
type CallToolParams = mcp.CallToolParams
type CallToolRequest = mcp.CallToolRequest
type CallToolResult = mcp.CallToolResult

type GetPromptRequest = mcp.GetPromptRequest
type GetPromptResult = mcp.GetPromptResult
type PromptMessage = mcp.PromptMessage
type Prompt = mcp.Prompt
type PromptOption = mcp.PromptOption

type Role = mcp.Role

const (
	RoleAssistant = mcp.RoleAssistant
	RoleUser      = mcp.RoleUser
)

type TextContent = mcp.TextContent
type Tool = mcp.Tool
type ToolOption = mcp.ToolOption
type PropertyOption = mcp.PropertyOption

type MCPServer = server.MCPServer
type ServerOption = server.ServerOption

type ToolHandlerFunc = func(context.Context, CallToolRequest) (*CallToolResult, error)
type PromptHandlerFunc = func(context.Context, GetPromptRequest) (*GetPromptResult, error)
