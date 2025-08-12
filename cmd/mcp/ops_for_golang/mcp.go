package ops_for_golang

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/ops_for_golang/usecases"
)

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#
func handleTestCoverage(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	directory, err := request.RequireString("directory")
	if err != nil {
		return nil, err
	}

	grepPattern := request.GetString("grep_pattern", "")

	// GolangOpsServiceを初期化
	service := usecases.NewGolangOpsService()
	result, err := service.HandleTestCoverage(directory, grepPattern)
	if err != nil {
		return nil, fmt.Errorf("テストカバレッジの実行に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

func handleTestCoverageProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	directory, err := request.RequireString("directory")
	if err != nil {
		return nil, err
	}

	// GolangOpsServiceを初期化
	service := usecases.NewGolangOpsService()
	result, err := service.HandleTestCoverageProject(directory)
	if err != nil {
		return nil, fmt.Errorf("プロジェクト全体のテストカバレッジの実行に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

func handleCoverageFunc(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	coverageFile, err := request.RequireString("coverage_file")
	if err != nil {
		return nil, err
	}

	grepPattern := request.GetString("grep_pattern", "")

	// GolangOpsServiceを初期化
	service := usecases.NewGolangOpsService()
	result, err := service.HandleCoverageFunc(coverageFile, grepPattern)
	if err != nil {
		return nil, fmt.Errorf("カバレッジ関数情報の取得に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

func handleGoRun(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	executionFile, err := request.RequireString("execution_file")
	if err != nil {
		return nil, err
	}

	rootDirectory, err := request.RequireString("root_directory")
	if err != nil {
		return nil, err
	}

	parameters := request.GetString("parameters", "")

	// GolangOpsServiceを初期化
	service := usecases.NewGolangOpsService()
	result, err := service.HandleGoRun(executionFile, rootDirectory, parameters)
	if err != nil {
		return nil, fmt.Errorf("go runの実行に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#
func setGolangOpsServer(s *server.MCPServer) *server.MCPServer {
	// ツール: テストカバレッジ取得
	testCoverageTool := mcp.NewTool(
		"golang_test_coverage",
		mcp.WithDescription("Get test coverage for Go development (`go test -cover ./...`)"),
		mcp.WithString(
			"directory",
			mcp.Required(),
			mcp.Description("Absolute path to the target directory"),
		),
		mcp.WithString(
			"grep_pattern",
			mcp.Description("Regular expression pattern for output filtering (optional)"),
		),
	)
	s.AddTool(testCoverageTool, handleTestCoverage)

	// ツール: プロジェクト全体のテストカバレッジ取得
	testCoverageProjectTool := mcp.NewTool(
		"golang_test_coverage_project",
		mcp.WithDescription("Get project-wide test coverage and generate HTML report for Go development (requires go.mod in target directory)"),
		mcp.WithString(
			"directory",
			mcp.Required(),
			mcp.Description("Absolute path to the project directory (must contain go.mod)"),
		),
	)
	s.AddTool(testCoverageProjectTool, handleTestCoverageProject)

	// ツール: カバレッジファイルから関数情報取得
	coverageFuncTool := mcp.NewTool(
		"golang_coverage_func",
		mcp.WithDescription("Display function-level coverage information from coverage file (`go tool cover -func=coverage.out`)"),
		mcp.WithString(
			"coverage_file",
			mcp.Required(),
			mcp.Description("Absolute path to the coverage file (e.g., coverage.out). The directory containing  that coverage file also must contain go.mod."),
		),
		mcp.WithString(
			"grep_pattern",
			mcp.Description("Regular expression pattern for output filtering (optional)"),
		),
	)
	s.AddTool(coverageFuncTool, handleCoverageFunc)

	// ツール: go run実行
	goRunTool := mcp.NewTool(
		"golang_run",
		mcp.WithDescription("Execute CLI tools for Go development (`go run`)"),
		mcp.WithString(
			"execution_file",
			mcp.Required(),
			mcp.Description("Absolute path to the Go file to execute (e.g., /path/to/main.go)"),
		),
		mcp.WithString(
			"root_directory",
			mcp.Required(),
			mcp.Description("Absolute path to the root directory where go run command should be executed"),
		),
		mcp.WithString(
			"parameters",
			mcp.Description("Runtime parameters (optional)"),
		),
	)
	s.AddTool(goRunTool, handleGoRun)

	return s
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt(
		"golang_ops_prompt",
		mcp.WithPromptDescription("Go開発用操作ツールのシステムプロンプト"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "Go開発でよく使用するコマンドを効率的に実行するためのツールです。",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleAssistant,
					Content: mcp.NewTextContent(`このツールは、Go開発でよく実行するコマンドを自動化します。

利用可能な機能:
- golang_test_coverage: テストカバレッジ取得
- golang_test_coverage_project: プロジェクト全体のカバレッジ取得とHTMLレポート生成
- golang_coverage_func: カバレッジファイルから関数レベル情報表示
- golang_run: go runでCLIツール実行

grepパターンを使用して出力をフィルタリングすることも可能です。`),
				},
			},
		}, nil
	})
	return s
}

func createGolangOpsServer() *server.MCPServer {
	s := server.NewMCPServer(
		"Golang Operations",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setGolangOpsServer(s)
	s = addPromptIntoServer(s)
	return s
}

func BuildGolangOpsServer() {
	s := createGolangOpsServer()
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
