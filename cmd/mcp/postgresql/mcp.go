package postgresql

import (
	"context"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

const (
	version = "1.0.0"
)

// setPostgreSQLQueryServer は受け取ったMCPサーバにPostgreSQL用のツールを付与して、そのMCPサーバを返します。
func setPostgreSQLQueryServer(databaseURL, timezone string, s *server.MCPServer) *server.MCPServer {
	// PostgreSQLハンドラーを初期化
	handler, err := NewPostgreSQLMCPHandler(databaseURL, timezone)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create PostgreSQL handler: %v\n", err)
		os.Exit(1)
	}

	// ツール1: SQL読み取り専用クエリの実行
	queryTool := mcp.NewTool("query",
		mcp.WithDescription("Run a read-only SQL query"),
		mcp.WithString("sql",
			mcp.Required(),
			mcp.Description("SQL query to execute"),
		),
	)
	s.AddTool(queryTool, handler.HandleToQuery)

	// ツール2: テーブルスキーマの取得
	getTableSchemaTool := mcp.NewTool("get_table_schema",
		mcp.WithDescription("Get schema information for a database table"),
		mcp.WithString("table_name",
			mcp.Required(),
			mcp.Description("Name of the table"),
		),
	)
	s.AddTool(getTableSchemaTool, handler.HandleToGetTableSchema)

	// ツール3: テーブル一覧の取得（最小限の情報）
	listTablesMinimumTool := mcp.NewTool("list_tables_minimum",
		mcp.WithDescription("List all tables in the database"),
	)
	s.AddTool(listTablesMinimumTool, handler.HandleToListTablesMinimum)

	// ツール4: テーブルの最小限のスキーマ情報を取得
	getTableSchemaMinimumTool := mcp.NewTool("get_table_schema_minimum",
		mcp.WithDescription("Get minimal schema information for a database table (column names and data types only)"),
		mcp.WithString("table_name",
			mcp.Required(),
			mcp.Description("Name of the table"),
		),
	)
	s.AddTool(getTableSchemaMinimumTool, handler.HandleToGetTableSchemaMinimum)

	// ツール5: テーブル一覧の取得
	listTablesTool := mcp.NewTool("list_tables",
		mcp.WithDescription("List all tables in the database"),
	)
	s.AddTool(listTablesTool, handler.HandleToListTables)

	// ツール6: テーブルダンプ
	dumpTableTool := mcp.NewTool("dump_table",
		mcp.WithDescription("Dump all records from a specified table to a file"),
		mcp.WithString("table_name",
			mcp.Required(),
			mcp.Description("Name of the table to dump"),
		),
		mcp.WithString("output_path",
			mcp.Description("Absolute directory path to save the dump file (default: current directory)"),
		),
		mcp.WithString("format",
			mcp.Description("Output format: json, csv, or sql (default: json)"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of records to dump (optional)"),
		),
	)
	s.AddTool(dumpTableTool, handler.HandleToDumpTable)

	return s
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a prompt for the PostgreSQL client."),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for the PostgreSQL client.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this great client for PostgreSQL well."),
				},
			},
		}, nil
	})
	return s
}

func createPostgreSQLServer() *server.MCPServer {
	// 環境変数からデータベースURLを取得
	databaseURL := os.Getenv("POSTGRESQL_DATABASE_URL")
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "Error: POSTGRESQL_DATABASE_URL environment variable not set.")
		os.Exit(1)
	}
	timezone := os.Getenv("POSTGRESQL_TIMEZONE")

	// MCPサーバーを作成
	s := server.NewMCPServer(
		"PostgreSQL Database Server",
		version,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)
	s = setPostgreSQLQueryServer(databaseURL, timezone, s)

	// プロンプト
	s = addPromptIntoServer(s)

	return s
}

// BuildPostgreSQLServer はPostgreSQLのMCPサーバーを構築します
func BuildPostgreSQLServer() {
	s := createPostgreSQLServer()

	// サーバーを起動
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}
