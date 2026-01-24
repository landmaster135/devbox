package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"

	dbExecutor "github.com/landmaster135/devbox/internal/postgresql/domain/executor"
	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	dump "github.com/landmaster135/devbox/internal/postgresql/usecases/dump"
	metaFetch "github.com/landmaster135/devbox/internal/postgresql/usecases/meta_fetch"
	templateRenderer "github.com/landmaster135/devbox/internal/postgresql/usecases/template_renderer"
)

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#

// DatabaseExecutor はデータベース操作のインターフェースです
type DatabaseExecutor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Ping() error
	Close() error
}

// TemplateRenderer はテンプレート処理のインターフェースです
type TemplateRenderer interface {
	RenderTableDetail(detail *model.TableDetail) (string, error)
	RenderTableList(data model.ListTablesData) (string, error)
}

// JSONMarshaler はJSON変換のインターフェースです
type JSONMarshaler interface {
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
}

// DefaultJSONMarshaler は標準のjson.MarshalIndentを使用する実装
type DefaultJSONMarshaler struct{}

func (m *DefaultJSONMarshaler) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// #==============================================================#
// ##          PostgreSQLService                                 ##
// #==============================================================#

// PostgreSQLService はPostgreSQLデータベースとの接続を管理します
type PostgreSQLService struct {
	executor         DatabaseExecutor
	templateRenderer TemplateRenderer
	jsonMarshaler    JSONMarshaler
	tableDumper      *dump.TableDumper
	databaseURL      string
	resourceBase     string
}

// NewPostgreSQLService は新しいPostgreSQLServiceを作成します
func NewPostgreSQLService(databaseURL string) (*PostgreSQLService, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("データベース接続の作成に失敗しました: %w", err)
	}

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベースへの接続テストに失敗しました: %w", err)
	}

	// リソースベースURLを作成
	resourceBase, err := createResourceBaseURL(databaseURL)
	if err != nil {
		return nil, err
	}

	executor := &dbExecutor.DefaultDatabaseExecutor{DB: db}

	return &PostgreSQLService{
		executor:         executor,
		templateRenderer: &templateRenderer.DefaultTemplateRenderer{},
		jsonMarshaler:    &DefaultJSONMarshaler{},
		tableDumper:      dump.NewTableDumper(executor),
		databaseURL:      databaseURL,
		resourceBase:     resourceBase,
	}, nil
}

// NewPostgreSQLServiceWithDependencies はテスト用に依存性を注入できるPostgreSQLServiceを作成します
func NewPostgreSQLServiceWithDependencies(executor DatabaseExecutor, templateRenderer TemplateRenderer, jsonMarshaler JSONMarshaler, tableDumper *dump.TableDumper, databaseURL, resourceBase string) *PostgreSQLService {
	return &PostgreSQLService{
		executor:         executor,
		templateRenderer: templateRenderer,
		jsonMarshaler:    jsonMarshaler,
		tableDumper:      tableDumper,
		databaseURL:      databaseURL,
		resourceBase:     resourceBase,
	}
}

// Close はデータベース接続を閉じます
func (s *PostgreSQLService) Close() error {
	return s.executor.Close()
}

// createResourceBaseURL はリソースベースURLを作成します
func createResourceBaseURL(databaseURL string) (string, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", err
	}

	// 基本的なURL検証
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("無効なURL形式です: %s", databaseURL)
	}

	// プロトコルをpostgresに変更し、パスワードを削除
	u.Scheme = "postgres"
	if u.User != nil {
		u.User = url.User(u.User.Username())
	}

	return u.String(), nil
}

// #==============================================================#
// ##          Query Execution Methods                           ##
// #==============================================================#

// ExecuteQuery はSQL読み取り専用クエリを実行します
func (s *PostgreSQLService) ExecuteQuery(ctx context.Context, sqlQuery string) ([]map[string]any, error) {
	// トランザクションを開始（読み取り専用）
	tx, err := s.executor.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: true,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// クエリを実行
	rows, err := tx.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// カラム名を取得
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 結果を格納するスライス
	var result []map[string]any

	// 各行を処理
	for rows.Next() {
		// スキャン用のインターフェースのスライスを作成
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// 行をスキャン
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// 行データをマップに変換
		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			// バイト配列の場合は文字列に変換
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// #==============================================================#
// ##          Handler Methods                                   ##
// #==============================================================#

// HandleToQuery はSQL読み取り専用クエリを実行して、結果をJSON形式で返します
func (s *PostgreSQLService) HandleToQuery(ctx context.Context, sqlQuery string) ([]map[string]any, error) {
	return s.ExecuteQuery(ctx, sqlQuery)
}

// HandleToGetTableSchema はテーブルのスキーマ情報を取得して、結果をテキスト形式で返します
func (s *PostgreSQLService) HandleToGetTableSchema(ctx context.Context, tableName string) (string, error) {
	// テーブルの詳細情報を取得
	detail, err := metaFetch.GetTableDetail(ctx, s.executor, tableName)
	if err != nil {
		return "", err
	}

	// テンプレートを使用して結果をフォーマット
	return s.templateRenderer.RenderTableDetail(detail)
}

// HandleToListTables はデータベース内のテーブル一覧を取得して、結果をテキスト形式で返します
func (s *PostgreSQLService) HandleToListTables(ctx context.Context) (string, error) {
	// テーブル情報の取得
	tables, err := metaFetch.GetAllTableSummaries(ctx, s.executor)
	if err != nil {
		return "", fmt.Errorf("テーブル情報の取得に失敗しました: %w", err)
	}

	// テーブルが見つからない
	if len(tables) == 0 {
		return "データベース内にテーブルが存在しません。", nil
	}

	// 出力の作成
	return s.templateRenderer.RenderTableList(model.ListTablesData{
		Tables: tables,
	})
}

func (s *PostgreSQLService) HandleGetAllTableSummaries(ctx context.Context) (string, error) {
	// テーブル情報の取得
	tables, err := metaFetch.GetAllTableSummaries(ctx, s.executor)
	if err != nil {
		return "", fmt.Errorf("テーブル情報の取得に失敗しました: %w", err)
	}

	// テーブルが見つからない
	if len(tables) == 0 {
		return "データベース内にテーブルが存在しません。", nil
	}

	// 結果をJSON形式で標準出力に表示
	jsonResult, err := json.MarshalIndent(tables, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSONパースに失敗しました: %w", err)
	}

	// 出力の作成
	return string(jsonResult), nil
}

// HandleToGetTableSchemaMinimum はテーブルの最小限のスキーマ情報を取得して、結果をJSON形式で返します
func (s *PostgreSQLService) HandleToGetTableSchemaMinimum(ctx context.Context, tableName string) (string, error) {
	schema, err := metaFetch.GetTableSchemaMinimum(ctx, s.executor, tableName)
	if err != nil {
		return "", fmt.Errorf("テーブルスキーマの取得に失敗しました: %v\n", err)
	}

	// 結果をJSON形式で標準出力に表示
	jsonResult, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換に失敗しました: %v\n", err)
	}

	return string(jsonResult), nil
}

// HandleToListTablesMinimum はデータベース内のテーブル一覧を取得して、結果をJSON形式で返します
func (s *PostgreSQLService) HandleToListTablesMinimum(ctx context.Context) (string, error) {
	tables, err := metaFetch.GetTablesMinimum(ctx, s.executor)
	if err != nil {
		return "", fmt.Errorf("テーブル一覧の取得に失敗しました: %v\n", err)
	}

	// 結果をJSON形式で標準出力に表示
	jsonResult, err := json.MarshalIndent(tables, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換に失敗しました: %v\n", err)
	}

	return string(jsonResult), nil
}

// HandleToDumpTable はテーブルの全レコードをダンプして、結果をJSON形式で返します
func HandleToDumpTable(ctx context.Context, dbURL, tableName, outputPath, format string, limit *int) (string, error) {
	// デフォルト値を設定
	if outputPath == "" {
		outputPath = "."
	}
	if format == "" {
		format = "json"
	}

	// PostgreSQLサービスを初期化
	service, err := NewPostgreSQLService(dbURL)
	if err != nil {
		return "", fmt.Errorf("PostgreSQLサービスの初期化に失敗しました: %v", err)
	}
	defer service.Close()

	// ダンプオプションを作成
	options := dump.NewDumpOptions(tableName, outputPath, format, limit)

	// ダンプを実行
	result, err := service.tableDumper.DumpTable(ctx, options)
	if err != nil {
		return "", fmt.Errorf("テーブルダンプの実行に失敗しました: %v", err)
	}

	// 結果をJSON形式で標準出力に表示
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("結果のJSON変換に失敗しました: %v", err)
	}

	return string(jsonResult), nil
}

// HandleToDumpAllTables はデータベース内の全テーブルをダンプして、結果をJSON形式で返します
func HandleToDumpAllTables(ctx context.Context, dbURL string, outputPath, format string, limit *int, concurrency *int) (string, string, error) {
	// デフォルト値を設定
	if outputPath == "" {
		outputPath = "."
	}
	if format == "" {
		format = "json"
	}

	// PostgreSQLサービスを初期化
	service, err := NewPostgreSQLService(dbURL)
	if err != nil {
		return "", "", fmt.Errorf("PostgreSQLサービスの初期化に失敗しました: %v", err)
	}
	defer service.Close()

	// 全テーブルダンプを実行
	jsonResult, minJSONResult, err := service.tableDumper.DumpAllTablesAndOutputJSON(ctx, outputPath, format, limit, concurrency)
	if err != nil {
		return "", "", err
	}

	return jsonResult, minJSONResult, nil
}
