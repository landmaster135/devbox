package timezone

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/timezone/usecases"
)

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a prompt for the timezone client."),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for the timezone client.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this great client for timezone well."),
				},
			},
		}, nil
	})
	return s
}

// BuildTimezoneServer はタイムゾーンMCPサーバーを構築する関数です
func BuildTimezoneServer() {
	s := server.NewMCPServer(
		"Timezone Service",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	// 現在時刻を取得するツール
	getCurrentTool := mcp.NewTool("get-current-timezone",
		mcp.WithDescription("Get the current time in the specified timezone"),
		mcp.WithString("timezone",
			mcp.Required(),
			mcp.Description("The timezone to get the current time for"),
		),
	)

	s.AddTool(getCurrentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		timezone, err := request.RequireString("timezone")
		if err != nil {
			return nil, err
		}

		service := usecases.NewTimezoneService()
		result, err := service.HandleGetCurrentTime(timezone)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	})

	// 時刻変換ツール
	convertTool := mcp.NewTool("convert-timezone",
		mcp.WithDescription("Convert time between timezones"),
		mcp.WithString("datetime",
			mcp.Required(),
			mcp.Description("The date and time to convert (e.g. 2025-03-31 12:00:00)"),
		),
		mcp.WithString("from_timezone",
			mcp.Required(),
			mcp.Description("The source timezone"),
		),
		mcp.WithString("to_timezone",
			mcp.Required(),
			mcp.Description("The target timezone"),
		),
	)

	s.AddTool(convertTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		datetime, err := request.RequireString("datetime")
		if err != nil {
			return nil, err
		}
		fromTimezone, err := request.RequireString("from_timezone")
		if err != nil {
			return nil, err
		}
		toTimezone, err := request.RequireString("to_timezone")
		if err != nil {
			return nil, err
		}

		service := usecases.NewTimezoneService()
		result, err := service.HandleConvertTime(datetime, fromTimezone, toTimezone)
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	})

	// 利用可能なタイムゾーンを取得するツール
	listTool := mcp.NewTool("list-available-timezones",
		mcp.WithDescription("List all available timezones"),
	)

	s.AddTool(listTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		service := usecases.NewTimezoneService()
		result, err := service.HandleListAvailableTimezones()
		if err != nil {
			return nil, err
		}

		return mcp.NewToolResultText(result), nil
	})

	// プロンプト
	s = addPromptIntoServer(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("サーバーエラー: %v\n", err)
	}
}
