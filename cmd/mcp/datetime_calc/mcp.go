package datetime_calc

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/datetime_calculator/usecases"
)

func handleDatetimeCalc(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, err := request.RequireString("operation")
	if err != nil {
		return nil, err
	}

	year_1, err := request.RequireFloat("year_1")
	if err != nil {
		return nil, err
	}

	month_1, err := request.RequireFloat("month_1")
	if err != nil {
		return nil, err
	}

	day_1, err := request.RequireFloat("day_1")
	if err != nil {
		return nil, err
	}

	hour_1, err := request.RequireFloat("hour_1")
	if err != nil {
		return nil, err
	}

	minute_1, err := request.RequireFloat("minute_1")
	if err != nil {
		return nil, err
	}

	second_1, err := request.RequireFloat("second_1")
	if err != nil {
		return nil, err
	}

	duration_of_year, err := request.RequireFloat("duration_of_year")
	if err != nil {
		return nil, err
	}

	duration_of_month, err := request.RequireFloat("duration_of_month")
	if err != nil {
		return nil, err
	}

	duration_of_day, err := request.RequireFloat("duration_of_day")
	if err != nil {
		return nil, err
	}

	duration_of_hour, err := request.RequireFloat("duration_of_hour")
	if err != nil {
		return nil, err
	}

	duration_of_minute, err := request.RequireFloat("duration_of_minute")
	if err != nil {
		return nil, err
	}

	duration_of_second, err := request.RequireFloat("duration_of_second")
	if err != nil {
		return nil, err
	}

	// DatetimeCalculatorServiceを初期化
	service := usecases.NewDatetimeCalculatorService()
	result, err := service.HandleDatetimeCalc(op, year_1, month_1, day_1, hour_1, minute_1, second_1, duration_of_year, duration_of_month, duration_of_day, duration_of_hour, duration_of_minute, duration_of_second)
	if err != nil {
		return nil, fmt.Errorf("日時計算に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a datetime calculator prompt"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for datetime calculator.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this accurate calculator well."),
				},
			},
		}, nil
	})
	return s
}

func BuildTimeCalculatorServer() {
	s := server.NewMCPServer(
		"Time Calculator",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	tool := mcp.NewTool("datetime_calc",
		mcp.WithDescription("Perform basic time calculations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to calculate datetime"),
			mcp.Enum("add", "subtract"),
			// mcp.Enum("add", "subtract"),
		),
		mcp.WithNumber("year_1",
			mcp.Required(),
			mcp.Description("Year of first datetime"),
		),
		mcp.WithNumber("month_1",
			mcp.Required(),
			mcp.Description("Month of first datetime"),
		),
		mcp.WithNumber("day_1",
			mcp.Required(),
			mcp.Description("Day of first datetime"),
		),
		mcp.WithNumber("hour_1",
			mcp.Required(),
			mcp.Description("Hour of first datetime"),
		),
		mcp.WithNumber("minute_1",
			mcp.Required(),
			mcp.Description("Minute of first datetime"),
		),
		mcp.WithNumber("second_1",
			mcp.Required(),
			mcp.Description("Second of first datetime"),
		),
		mcp.WithNumber("duration_of_year",
			mcp.Required(),
			mcp.Description("Duration of year"),
		),
		mcp.WithNumber("duration_of_month",
			mcp.Required(),
			mcp.Description("Duration of month"),
		),
		mcp.WithNumber("duration_of_day",
			mcp.Required(),
			mcp.Description("Duration of day"),
		),
		mcp.WithNumber("duration_of_hour",
			mcp.Required(),
			mcp.Description("Duration of hour"),
		),
		mcp.WithNumber("duration_of_minute",
			mcp.Required(),
			mcp.Description("Duration of minute"),
		),
		mcp.WithNumber("duration_of_second",
			mcp.Required(),
			mcp.Description("Duration of second"),
		),
	)

	s.AddTool(tool, handleDatetimeCalc)

	// プロンプト
	s = addPromptIntoServer(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
