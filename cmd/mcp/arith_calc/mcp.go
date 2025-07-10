package arith_calc

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases"
)

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#
func handleToCalculate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, err := request.RequireString("operation")
	if err != nil {
		return nil, err
	}

	x, err := request.RequireFloat("x")
	if err != nil {
		return nil, err
	}

	y, err := request.RequireFloat("y")
	if err != nil {
		return nil, err
	}

	// CalculatorServiceを初期化
	service := usecases.NewCalculatorService()
	result, err := service.HandleToCalculate(op, x, y)
	if err != nil {
		return nil, fmt.Errorf("パラメータを用いた算術計算に失敗しました: %v", err)
	}

	return mcp.FormatNumberResult(result), nil
}

func handleToCalculateWithArray(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, err := request.RequireString("operation")
	if err != nil {
		return nil, err
	}

	numbers, err := request.RequireFloatSlice("numbers")
	if err != nil {
		return nil, err
	}

	service := usecases.NewCalculatorService()
	result, err := service.HandleToCalculateWithArray(op, numbers)
	if err != nil {
		return nil, fmt.Errorf("配列を用いた算術計算に失敗しました: %v", err)
	}

	return mcp.FormatNumberResult(result), nil
}

func handleToEvaluateLineCount(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	// 必須パラメータの取得
	filePath, err := request.RequireString("file_path")
	if err != nil {
		return nil, err
	}

	threshold, err := request.RequireFloat("threshold")
	if err != nil {
		return nil, err
	}

	// FileEvaluatorServiceを初期化
	service := usecases.NewFileEvaluatorService()
	jsonResult, err := service.HandleToEvaluateLineCount(filePath, int(threshold))
	if err != nil {
		return nil, fmt.Errorf("配列を用いた算術計算に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#
func setTwoNumbersInputtingCalcServer(s *server.MCPServer) *server.MCPServer {
	tool := mcp.NewTool(
		"calculate",
		mcp.WithDescription("Perform basic arithmetic calculations with two numbers"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The arithmetic operation to perform"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)
	s.AddTool(tool, handleToCalculate)

	toolWithArray := mcp.NewTool(
		"calculate_with_multiple_numbers",
		mcp.WithDescription("Perform basic arithmetic calculations with multiple numbers"),
		mcp.WithString(
			"operation",
			mcp.Required(),
			mcp.Description("The arithmetic operation to perform with multiple numbers"),
			mcp.Enum("sum"),
		),
		mcp.WithArray(
			"numbers",
			mcp.Required(),
			mcp.Description("Multiple numbers"),
		),
	)
	s.AddTool(toolWithArray, handleToCalculateWithArray)

	return s
}

// setFileLineCountEvaluatorServer はファイルの行数評価ツールを提供するMCPサーバを設定します
func setFileLineCountEvaluatorServer(s *server.MCPServer) *server.MCPServer {

	// ツール: ファイルの行数評価
	tool := mcp.NewTool(
		"evaluate_line_count",
		mcp.WithDescription("ファイルの行数が指定された閾値より大きいかどうかを評価します"),
		mcp.WithString(
			"file_path",
			mcp.Required(),
			mcp.Description("評価するファイルの絶対パス"),
		),
		mcp.WithNumber(
			"threshold",
			mcp.Required(),
			mcp.Description("比較する行数の閾値"),
		),
	)
	s.AddTool(tool, handleToEvaluateLineCount)

	return s
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt(
		"system_prompt_01",
		mcp.WithPromptDescription("This is an arithmetic calculator prompt"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for arithmetic calculator.",
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

func createArithCalcServer() *server.MCPServer {
	s := server.NewMCPServer(
		"Arithmetic Calculator",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setTwoNumbersInputtingCalcServer(s)
	s = setFileLineCountEvaluatorServer(s)
	s = addPromptIntoServer(s)
	return s
}

func BuildArithCalculatorServer() {
	s := createArithCalcServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
