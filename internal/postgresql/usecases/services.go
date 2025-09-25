package usecases

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
)

// #==============================================================#
// ##          Data Structures                                   ##
// #==============================================================#

// ColumnInfo はデータベースのカラム詳細情報を表します
type ColumnInfo struct {
	Name       string         `json:"column_name"`
	Type       string         `json:"data_type"`
	IsNullable string         `json:"is_nullable"`
	Default    sql.NullString `json:"column_default"`
	Comment    string         `json:"column_comment"`
}

// TableSummary はテーブルのサマリー情報を表します
type TableSummary struct {
	Name    string       `json:"table_name"`
	Comment string       `json:"table_comment"`
	PK      []string     `json:"primary_keys"`
	UK      []UniqueKey  `json:"unique_keys"`
	FK      []ForeignKey `json:"foreign_keys"`
}

// UniqueKey は一意キー制約を表します
type UniqueKey struct {
	Name    string   `json:"constraint_name"`
	Columns []string `json:"columns"`
}

// ForeignKey は外部キー制約を表します
type ForeignKey struct {
	Name       string   `json:"constraint_name"`
	Columns    []string `json:"columns"`
	RefTable   string   `json:"referenced_table"`
	RefColumns []string `json:"referenced_columns"`
}

// IndexInfo はインデックス情報を表します
type IndexInfo struct {
	Name    string   `json:"index_name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"is_unique"`
}

// TableDetail はテーブルの詳細情報を表します
type TableDetail struct {
	Name        string       `json:"table_name"`
	Comment     string       `json:"table_comment"`
	Columns     []ColumnInfo `json:"columns"`
	PrimaryKeys []string     `json:"primary_keys"`
	UniqueKeys  []UniqueKey  `json:"unique_keys"`
	ForeignKeys []ForeignKey `json:"foreign_keys"`
	Indexes     []IndexInfo  `json:"indexes"`
}

// Column はデータベースのカラム情報を表します（最小限の情報）
type Column struct {
	Name     string `json:"column_name"`
	DataType string `json:"data_type"`
}

// Table はデータベースのテーブル情報を表します
type Table struct {
	Name string `json:"table_name"`
}

// ListTablesData はテーブル一覧のテンプレートに渡すデータ構造
type ListTablesData struct {
	Tables []TableSummary
}

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
	RenderTableDetail(detail *TableDetail) (string, error)
	RenderTableList(data ListTablesData) (string, error)
}

// JSONMarshaler はJSON変換のインターフェースです
type JSONMarshaler interface {
	MarshalIndent(v any, prefix, indent string) ([]byte, error)
}

// #==============================================================#
// ##          Default Implementations                           ##
// #==============================================================#

// DefaultDatabaseExecutor は標準のsql.DBを使用する実装
type DefaultDatabaseExecutor struct {
	db *sql.DB
}

func (d *DefaultDatabaseExecutor) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

func (d *DefaultDatabaseExecutor) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

func (d *DefaultDatabaseExecutor) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, opts)
}

func (d *DefaultDatabaseExecutor) Ping() error {
	return d.db.Ping()
}

func (d *DefaultDatabaseExecutor) Close() error {
	return d.db.Close()
}

// QueryContextRows は新しいインターフェース用のメソッド
func (d *DefaultDatabaseExecutor) QueryContextRows(ctx context.Context, query string, args ...any) (RowsInterface, error) {
	rows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &SQLRowsWrapper{rows: rows}, nil
}

// QueryRowContextRow は新しいインターフェース用のメソッド
func (d *DefaultDatabaseExecutor) QueryRowContextRow(ctx context.Context, query string, args ...any) RowInterface {
	row := d.db.QueryRowContext(ctx, query, args...)
	return &SQLRowWrapper{row: row}
}

// SQLRowsWrapper は *sql.Rows を RowsInterface として扱うためのラッパー
type SQLRowsWrapper struct {
	rows *sql.Rows
}

func (w *SQLRowsWrapper) Columns() ([]string, error) {
	return w.rows.Columns()
}

func (w *SQLRowsWrapper) Next() bool {
	return w.rows.Next()
}

func (w *SQLRowsWrapper) Scan(dest ...any) error {
	return w.rows.Scan(dest...)
}

func (w *SQLRowsWrapper) Close() error {
	return w.rows.Close()
}

func (w *SQLRowsWrapper) Err() error {
	return w.rows.Err()
}

// SQLRowWrapper は *sql.Row を RowInterface として扱うためのラッパー
type SQLRowWrapper struct {
	row *sql.Row
}

func (w *SQLRowWrapper) Scan(dest ...any) error {
	return w.row.Scan(dest...)
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
	tableDumper      *TableDumper
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

	executor := &DefaultDatabaseExecutor{db: db}

	return &PostgreSQLService{
		executor:         executor,
		templateRenderer: &DefaultTemplateRenderer{},
		jsonMarshaler:    &DefaultJSONMarshaler{},
		tableDumper:      NewTableDumper(executor),
		databaseURL:      databaseURL,
		resourceBase:     resourceBase,
	}, nil
}

// NewPostgreSQLServiceWithDependencies はテスト用に依存性を注入できるPostgreSQLServiceを作成します
func NewPostgreSQLServiceWithDependencies(executor DatabaseExecutor, templateRenderer TemplateRenderer, jsonMarshaler JSONMarshaler, tableDumper *TableDumper, databaseURL, resourceBase string) *PostgreSQLService {
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
// ##          Table Information Methods                         ##
// #==============================================================#

// GetTablesMinimum はデータベース内のテーブル一覧を取得します
func (s *PostgreSQLService) GetTablesMinimum(ctx context.Context) ([]Table, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
	`

	rows, err := s.executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var table Table
		if err := rows.Scan(&table.Name); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

// GetTableSchemaMinimum はテーブルの最小限のスキーマ情報を取得します
func (s *PostgreSQLService) GetTableSchemaMinimum(ctx context.Context, tableName string) ([]Column, error) {
	_, schema, name, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_schema = $1
		  AND table_name = $2
	`

	rows, err := s.executor.QueryContext(ctx, query, schema, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var column Column
		if err := rows.Scan(&column.Name, &column.DataType); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

// fetchTableWithComments はテーブル名とコメントを取得します
func (s *PostgreSQLService) fetchTableWithComments(ctx context.Context, tableName string) (TableSummary, error) {
	_, schema, name, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return TableSummary{}, err
	}

	query := `
		SELECT t.table_name,
		       COALESCE(pg_catalog.obj_description(c.oid), '') AS table_comment
		FROM information_schema.tables t
		JOIN pg_catalog.pg_class c ON c.relname = t.table_name
		JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
		WHERE t.table_schema = $1
		  AND t.table_name = $2
		  AND n.nspname = $1
	`

	var table TableSummary
	err = s.executor.QueryRowContext(ctx, query, schema, name).Scan(&table.Name, &table.Comment)
	if err != nil {
		return TableSummary{}, err
	}

	return table, nil
}

// fetchPrimaryKeys はテーブルの主キーカラムを取得します
func (s *PostgreSQLService) fetchPrimaryKeys(ctx context.Context, tableName string) ([]string, error) {
	qualified, _, _, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT a.attname
		FROM pg_index i
		JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
		WHERE i.indrelid = $1::regclass
		AND i.indisprimary
		ORDER BY a.attnum
	`

	rows, err := s.executor.QueryContext(ctx, query, qualified)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var primaryKeys []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return primaryKeys, nil
}

// fetchUniqueKeys はテーブルの一意キー制約を取得します
func (s *PostgreSQLService) fetchUniqueKeys(ctx context.Context, tableName string) ([]UniqueKey, error) {
	query := `
		SELECT
			c.conname AS constraint_name,
			a.attname AS column_name
		FROM
			pg_constraint c
		JOIN
			pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
		WHERE
			c.conrelid = $1::regclass
			AND c.contype = 'u'
		ORDER BY
			c.conname, a.attnum
	`

	qualified, _, _, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	rows, err := s.executor.QueryContext(ctx, query, qualified)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uniqueKeys []UniqueKey
	var currentUniqueKey *UniqueKey
	for rows.Next() {
		var constraintName, columnName string
		if err := rows.Scan(&constraintName, &columnName); err != nil {
			return nil, err
		}

		if currentUniqueKey == nil || currentUniqueKey.Name != constraintName {
			uniqueKeys = append(uniqueKeys, UniqueKey{
				Name:    constraintName,
				Columns: []string{},
			})
			currentUniqueKey = &uniqueKeys[len(uniqueKeys)-1]
		}

		currentUniqueKey.Columns = append(currentUniqueKey.Columns, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return uniqueKeys, nil
}

// fetchForeignKeys はテーブルの外部キー制約を取得します
func (s *PostgreSQLService) fetchForeignKeys(ctx context.Context, tableName string) ([]ForeignKey, error) {
	query := `
		SELECT
			c.conname AS constraint_name,
			a.attname AS column_name,
			cl.relname AS referenced_table,
			af.attname AS referenced_column
		FROM
			pg_constraint c
		JOIN
			pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = ANY(c.conkey)
		JOIN
			pg_class cl ON cl.oid = c.confrelid
		JOIN
			pg_attribute af ON af.attrelid = c.confrelid AND af.attnum = ANY(c.confkey)
		WHERE
			c.conrelid = $1::regclass
			AND c.contype = 'f'
		ORDER BY
			c.conname, a.attnum
	`

	qualified, _, _, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	rows, err := s.executor.QueryContext(ctx, query, qualified)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foreignKeys []ForeignKey
	var currentFK *ForeignKey
	for rows.Next() {
		var constraintName, columnName, refTableName, refColumnName string
		if err := rows.Scan(&constraintName, &columnName, &refTableName, &refColumnName); err != nil {
			return nil, err
		}

		if currentFK == nil || currentFK.Name != constraintName {
			foreignKeys = append(foreignKeys, ForeignKey{
				Name:       constraintName,
				RefTable:   refTableName,
				Columns:    []string{},
				RefColumns: []string{},
			})
			currentFK = &foreignKeys[len(foreignKeys)-1]
		}

		currentFK.Columns = append(currentFK.Columns, columnName)
		currentFK.RefColumns = append(currentFK.RefColumns, refColumnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return foreignKeys, nil
}

// fetchTableColumns はテーブルのカラム情報を取得します
func (s *PostgreSQLService) fetchTableColumns(ctx context.Context, tableName string) ([]ColumnInfo, error) {
	_, schema, name, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			c.column_name,
			c.data_type,
			c.is_nullable,
			c.column_default,
			COALESCE(pg_catalog.col_description(
				pg_catalog.pg_class.oid,
				c.ordinal_position
			), '') AS column_comment
		FROM
			information_schema.columns c
		JOIN
			pg_catalog.pg_class ON pg_catalog.pg_class.relname = c.table_name
		WHERE
			c.table_schema = $1
			AND c.table_name = $2
		ORDER BY
			c.ordinal_position
	`

	rows, err := s.executor.QueryContext(ctx, query, schema, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		if err := rows.Scan(&col.Name, &col.Type, &col.IsNullable, &col.Default, &col.Comment); err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

// fetchTableIndexes はテーブルのインデックス情報を取得します
func (s *PostgreSQLService) fetchTableIndexes(ctx context.Context, tableName string) ([]IndexInfo, error) {
	_, schema, name, err := qualifyTableIdentifier(tableName)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			i.relname AS index_name,
			a.attname AS column_name,
			ix.indisunique AS is_unique
		FROM
			pg_index ix
		JOIN
			pg_class i ON i.oid = ix.indexrelid
		JOIN
			pg_attribute a ON a.attrelid = ix.indrelid AND a.attnum = ANY(ix.indkey)
		JOIN
			pg_class t ON t.oid = ix.indrelid
	JOIN
		pg_namespace n ON n.oid = t.relnamespace
		LEFT JOIN
			pg_constraint c ON c.conindid = ix.indexrelid
		WHERE
			t.relname = $1
			AND n.nspname = $2
			AND c.contype IS NULL
			AND NOT ix.indisprimary
		ORDER BY
			i.relname, a.attnum
	`

	rows, err := s.executor.QueryContext(ctx, query, name, schema)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []IndexInfo
	var currentIdx *IndexInfo
	for rows.Next() {
		var indexName, columnName string
		var isUnique bool
		if err := rows.Scan(&indexName, &columnName, &isUnique); err != nil {
			return nil, err
		}

		if currentIdx == nil || currentIdx.Name != indexName {
			indexes = append(indexes, IndexInfo{
				Name:    indexName,
				Unique:  isUnique,
				Columns: []string{},
			})
			currentIdx = &indexes[len(indexes)-1]
		}

		currentIdx.Columns = append(currentIdx.Columns, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return indexes, nil
}

// #==============================================================#
// ##          High-level Service Methods                        ##
// #==============================================================#

// GetTableDetail はテーブルの詳細情報を取得します
func (s *PostgreSQLService) GetTableDetail(ctx context.Context, tableName string) (*TableDetail, error) {
	// テーブル情報を取得
	tableInfo, err := s.fetchTableWithComments(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("テーブル情報の取得に失敗しました: %w", err)
	}

	// 主キー情報を取得
	primaryKeys, err := s.fetchPrimaryKeys(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("主キー情報の取得に失敗しました: %w", err)
	}

	// 一意キー情報を取得
	uniqueKeys, err := s.fetchUniqueKeys(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("一意キー情報の取得に失敗しました: %w", err)
	}

	// 外部キー情報を取得
	foreignKeys, err := s.fetchForeignKeys(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("外部キー情報の取得に失敗しました: %w", err)
	}

	// カラム情報を取得
	columns, err := s.fetchTableColumns(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("カラム情報の取得に失敗しました: %w", err)
	}

	// インデックス情報を取得
	indexes, err := s.fetchTableIndexes(ctx, tableName)
	if err != nil {
		return nil, fmt.Errorf("インデックス情報の取得に失敗しました: %w", err)
	}

	// 詳細情報を作成
	detail := &TableDetail{
		Name:        tableInfo.Name,
		Comment:     tableInfo.Comment,
		Columns:     columns,
		PrimaryKeys: primaryKeys,
		UniqueKeys:  uniqueKeys,
		ForeignKeys: foreignKeys,
		Indexes:     indexes,
	}

	return detail, nil
}

// GetAllTableSummaries はデータベース内の全てのテーブルのサマリー情報を取得します
func (s *PostgreSQLService) GetAllTableSummaries(ctx context.Context) ([]TableSummary, error) {
	// テーブル一覧を取得
	query := `
		SELECT t.table_name,
		       COALESCE(pg_catalog.obj_description(pg_catalog.pg_class.oid), '') AS table_comment
		FROM information_schema.tables t
		JOIN pg_catalog.pg_class ON pg_catalog.pg_class.relname = t.table_name
		WHERE t.table_schema = 'public'
		ORDER BY t.table_name
	`

	rows, err := s.executor.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []TableSummary
	for rows.Next() {
		var table TableSummary
		if err := rows.Scan(&table.Name, &table.Comment); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 各テーブルの追加情報を取得
	for i := range tables {
		// 主キー情報を取得
		tables[i].PK, err = s.fetchPrimaryKeys(ctx, tables[i].Name)
		if err != nil {
			return nil, err
		}

		// 一意キー情報を取得
		tables[i].UK, err = s.fetchUniqueKeys(ctx, tables[i].Name)
		if err != nil {
			return nil, err
		}

		// 外部キー情報を取得
		tables[i].FK, err = s.fetchForeignKeys(ctx, tables[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return tables, nil
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
	detail, err := s.GetTableDetail(ctx, tableName)
	if err != nil {
		return "", err
	}

	// テンプレートを使用して結果をフォーマット
	return s.templateRenderer.RenderTableDetail(detail)
}

// HandleToListTables はデータベース内のテーブル一覧を取得して、結果をテキスト形式で返します
func (s *PostgreSQLService) HandleToListTables(ctx context.Context) (string, error) {
	// テーブル情報の取得
	tables, err := s.GetAllTableSummaries(ctx)
	if err != nil {
		return "", fmt.Errorf("テーブル情報の取得に失敗しました: %w", err)
	}

	// テーブルが見つからない
	if len(tables) == 0 {
		return "データベース内にテーブルが存在しません。", nil
	}

	// 出力の作成
	return s.templateRenderer.RenderTableList(ListTablesData{
		Tables: tables,
	})
}

// HandleToGetTableSchemaMinimum はテーブルの最小限のスキーマ情報を取得して、結果をJSON形式で返します
func (s *PostgreSQLService) HandleToGetTableSchemaMinimum(ctx context.Context, tableName string) ([]Column, error) {
	return s.GetTableSchemaMinimum(ctx, tableName)
}

// HandleToListTablesMinimum はデータベース内のテーブル一覧を取得して、結果をJSON形式で返します
func (s *PostgreSQLService) HandleToListTablesMinimum(ctx context.Context) ([]Table, error) {
	return s.GetTablesMinimum(ctx)
}

// HandleToDumpTable はテーブルの全レコードをダンプして、結果をJSON形式で返します
func (s *PostgreSQLService) HandleToDumpTable(ctx context.Context, tableName, outputPath, format string, limit *int) (*DumpResult, error) {
	// デフォルト値を設定
	if outputPath == "" {
		outputPath = "."
	}
	if format == "" {
		format = "json"
	}

	// ダンプオプションを作成
	options := DumpOptions{
		TableName:  tableName,
		OutputPath: outputPath,
		Format:     format,
		Limit:      limit,
	}

	// ダンプを実行
	return s.tableDumper.DumpTable(ctx, options)
}

// HandleToDumpAllTables はデータベース内の全テーブルをダンプして、結果をJSON形式で返します
func (s *PostgreSQLService) HandleToDumpAllTables(ctx context.Context, outputPath, format string, limit *int, concurrency *int) (*DumpAllTablesResult, error) {
	// デフォルト値を設定
	if outputPath == "" {
		outputPath = "."
	}
	if format == "" {
		format = "json"
	}

	// 全テーブルダンプを実行
	return s.tableDumper.DumpAllTables(ctx, outputPath, format, limit, concurrency)
}
