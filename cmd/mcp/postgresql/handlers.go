package postgresql

import (
	"context"
	"encoding/json"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"

	usecases "github.com/landmaster135/devbox/internal/postgresql/usecases"
)

// PostgreSQLMCPHandler はMCPリクエストを処理するハンドラーです
type PostgreSQLMCPHandler struct {
	url     string
	service *usecases.PostgreSQLService
}

// NewPostgreSQLMCPHandler は新しいPostgreSQLMCPHandlerを作成します
func NewPostgreSQLMCPHandler(databaseURL string) (*PostgreSQLMCPHandler, error) {
	service, err := usecases.NewPostgreSQLService(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("PostgreSQLサービスの作成に失敗しました: %w", err)
	}

	return &PostgreSQLMCPHandler{
		url:     databaseURL,
		service: service,
	}, nil
}

// Close はリソースを解放します
func (h *PostgreSQLMCPHandler) Close() error {
	return h.service.Close()
}

// HandleToQuery はSQL読み取り専用クエリを実行して、結果をJSON形式で返します
func (h *PostgreSQLMCPHandler) HandleToQuery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sqlQuery, err := request.RequireString("sql")
	if err != nil {
		return returnError(err)
	}

	result, err := h.service.HandleToQuery(ctx, sqlQuery)
	if err != nil {
		return returnError(err)
	}

	return returnJSONResult(result)
}

// HandleToGetTableSchema はテーブルのスキーマ情報を取得して、結果をテキスト形式で返します
func (h *PostgreSQLMCPHandler) HandleToGetTableSchema(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tableName, err := request.RequireString("table_name")
	if err != nil {
		return returnError(err)
	}

	result, err := h.service.HandleToGetTableSchema(ctx, tableName)
	if err != nil {
		return returnError(err)
	}

	return returnTextResult(result)
}

// HandleToListTables はデータベース内のテーブル一覧を取得して、結果をテキスト形式で返します
func (h *PostgreSQLMCPHandler) HandleToListTables(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := h.service.HandleToListTables(ctx)
	if err != nil {
		return returnError(err)
	}

	return returnTextResult(result)
}

// HandleToGetTableSchemaMinimum はテーブルの最小限のスキーマ情報を取得して、結果をJSON形式で返します
func (h *PostgreSQLMCPHandler) HandleToGetTableSchemaMinimum(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tableName, err := request.RequireString("table_name")
	if err != nil {
		return returnError(err)
	}

	result, err := h.service.HandleToGetTableSchemaMinimum(ctx, tableName)
	if err != nil {
		return returnError(err)
	}

	return returnTextResult(result)
}

// HandleToListTablesMinimum はデータベース内のテーブル一覧を取得して、結果をJSON形式で返します
func (h *PostgreSQLMCPHandler) HandleToListTablesMinimum(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result, err := h.service.HandleToListTablesMinimum(ctx)
	if err != nil {
		return returnError(err)
	}

	return returnTextResult(result)
}

// HandleToDumpTable はテーブルの全レコードをダンプして、結果をJSON形式で返します
func (h *PostgreSQLMCPHandler) HandleToDumpTable(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tableName, err := request.RequireString("table_name")
	if err != nil {
		return nil, err
	}

	outputPath := request.GetString("output_path", "")
	format := request.GetString("format", "json")

	var limit *int
	if limitValue := request.GetInt("limit", 0); limitValue > 0 {
		limit = &limitValue
	}

	result, err := usecases.HandleToDumpTable(ctx, h.url, tableName, outputPath, format, limit)
	if err != nil {
		return nil, fmt.Errorf("テーブルダンプの実行に失敗しました: %v", err)
	}

	return returnJSONResult(result)
}

// #==============================================================#
// ##          Helper Functions                                  ##
// #==============================================================#

// returnJSONResult は結果をJSON形式で返却します
func returnJSONResult(result interface{}) (*mcp.CallToolResult, error) {
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(jsonResult)), nil
}

// returnTextResult はテキスト結果を返却します
func returnTextResult(result string) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(result), nil
}

// returnError はエラーを返却します
func returnError(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(fmt.Sprintf("Error: %v", err)), fmt.Errorf("database error: %w", err)
}
