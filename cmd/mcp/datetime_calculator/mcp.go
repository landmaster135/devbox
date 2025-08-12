package datetime_calculator

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

func handleTimeUnitSum(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	numbers, err := request.RequireFloatSlice("numbers")
	if err != nil {
		return nil, err
	}

	inputUnit, err := request.RequireString("input_unit")
	if err != nil {
		return nil, err
	}

	outputUnit, err := request.RequireString("output_unit")
	if err != nil {
		return nil, err
	}

	// DatetimeCalculatorServiceを初期化
	service := usecases.NewDatetimeCalculatorService()
	result, err := service.HandleTimeUnitSum(numbers, inputUnit, outputUnit)
	if err != nil {
		return nil, fmt.Errorf("時間単位合計計算に失敗しました: %v", err)
	}

	return mcp.FormatNumberResult(result), nil
}

func handleTimeExtractionFromSentence(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePath := request.GetString("file_path", "")
	textInput := request.GetString("text_input", "")
	outputUnit, err := request.RequireString("output_unit")
	if err != nil {
		return nil, err
	}

	// DatetimeCalculatorServiceを初期化
	service := usecases.NewDatetimeCalculatorService()
	result, err := service.HandleTimeExtraction(filePath, textInput, outputUnit)
	if err != nil {
		return nil, fmt.Errorf("時間抽出に失敗しました: %v", err)
	}

	return mcp.FormatNumberResult(result), nil
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

	datetimeCalcTool := mcp.NewTool("datetime_calc",
		mcp.WithDescription("Perform basic time calculations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to calculate datetime"),
			mcp.Enum("add", "subtract"),
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
	s.AddTool(datetimeCalcTool, handleDatetimeCalc)

	timeUnitSumTool := mcp.NewTool("time_unit_sum",
		mcp.WithDescription("Calculate sum of time values with unit conversion"),
		mcp.WithArray("numbers",
			mcp.Required(),
			mcp.Description("Array of numbers to sum"),
		),
		mcp.WithString("input_unit",
			mcp.Required(),
			mcp.Description("Input time unit"),
			mcp.Enum("year", "month", "day", "hour", "minute", "second"),
		),
		mcp.WithString("output_unit",
			mcp.Required(),
			mcp.Description("Output time unit"),
			mcp.Enum("year", "month", "day", "hour", "minute", "second"),
		),
	)
	s.AddTool(timeUnitSumTool, handleTimeUnitSum)

	// 新しいtime_extraction_from_sentenceツール
	timeExtractionTool := mcp.NewTool("time_extraction_from_sentence",
		mcp.WithDescription("Extract time values from text or file content"),
		mcp.WithString("file_path",
			mcp.Description("Absolute path to file (.md or .txt format)"),
		),
		mcp.WithString("text_input",
			mcp.Description("Text input to extract time from"),
		),
		mcp.WithString("output_unit",
			mcp.Required(),
			mcp.Description("Output time unit"),
			mcp.Enum("year", "month", "day", "hour", "minute", "second"),
		),
	)
	s.AddTool(timeExtractionTool, handleTimeExtractionFromSentence)

	// プロンプト
	s = addPromptIntoServer(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
