package usecases

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// #==============================================================#
// ##          Data Structures                                   ##
// #==============================================================#

// DumpOptions はダンプ処理のオプションを表します
type DumpOptions struct {
	TableName  string
	OutputPath string
	Format     string
	Limit      *int
}

// DumpResult はダンプ処理の結果を表します
type DumpResult struct {
	TableName    string `json:"table_name"`
	RecordCount  int    `json:"record_count"`
	OutputPath   string `json:"output_path"`
	FileName     string `json:"file_name"`
	Format       string `json:"format"`
	ExecutedAt   string `json:"executed_at"`
}

// DumpAllTablesResult は全テーブルダンプの結果を表します
type DumpAllTablesResult struct {
	DatabaseName string       `json:"database_name"`
	TotalTables  int          `json:"total_tables"`
	Results      []DumpResult `json:"results"`
	FailedTables []FailedDump `json:"failed_tables"`
	ExecutedAt   string       `json:"executed_at"`
}

// FailedDump は失敗したダンプ情報を表します
type FailedDump struct {
	TableName string `json:"table_name"`
	Error     string `json:"error"`
}

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#

// FileWriter はファイル書き込み操作のインターフェースです
type FileWriter interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (*os.File, error)
}

// #==============================================================#
// ##          Default Implementations                           ##
// #==============================================================#

// DefaultFileWriter は標準のファイル操作を使用する実装
type DefaultFileWriter struct{}

func (w *DefaultFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (w *DefaultFileWriter) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (w *DefaultFileWriter) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// #==============================================================#
// ##          TableDumper                                       ##
// #==============================================================#

// TableDumper はテーブルダンプ機能を提供します
type TableDumper struct {
	executor   DatabaseExecutor
	fileWriter FileWriter
}

// NewTableDumper は新しいTableDumperを作成します
func NewTableDumper(executor DatabaseExecutor) *TableDumper {
	return &TableDumper{
		executor:   executor,
		fileWriter: &DefaultFileWriter{},
	}
}

// NewTableDumperWithDependencies はテスト用に依存性を注入できるTableDumperを作成します
func NewTableDumperWithDependencies(executor DatabaseExecutor, fileWriter FileWriter) *TableDumper {
	return &TableDumper{
		executor:   executor,
		fileWriter: fileWriter,
	}
}

// DumpTable はテーブルの全レコードをダンプします
func (d *TableDumper) DumpTable(ctx context.Context, options DumpOptions) (*DumpResult, error) {
	// パラメータ検証
	if err := d.validateOptions(options); err != nil {
		return nil, fmt.Errorf("オプション検証エラー: %w", err)
	}

	// 出力ディレクトリを作成
	if err := d.ensureOutputDirectory(options.OutputPath); err != nil {
		return nil, fmt.Errorf("出力ディレクトリ作成エラー: %w", err)
	}

	// クエリを構築
	query := d.buildQuery(options)

	// データを取得
	data, err := d.fetchTableData(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("データ取得エラー: %w", err)
	}

	// ファイル名を生成
	fileName := d.generateFileName(options.TableName, options.Format)
	filePath := filepath.Join(options.OutputPath, fileName)

	// フォーマットに応じてファイルに出力
	if err := d.writeDataToFile(filePath, data, options.Format, options.TableName); err != nil {
		return nil, fmt.Errorf("ファイル書き込みエラー: %w", err)
	}

	// 結果を作成
	result := &DumpResult{
		TableName:   options.TableName,
		RecordCount: len(data),
		OutputPath:  options.OutputPath,
		FileName:    fileName,
		Format:      options.Format,
		ExecutedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	return result, nil
}

// validateOptions はオプションを検証します
func (d *TableDumper) validateOptions(options DumpOptions) error {
	if options.TableName == "" {
		return fmt.Errorf("テーブル名が指定されていません")
	}

	if options.OutputPath == "" {
		options.OutputPath = "."
	}

	// パストラバーサル攻撃を防ぐ
	cleanPath := filepath.Clean(options.OutputPath)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("無効なパスが指定されました: %s", options.OutputPath)
	}

	// サポートされているフォーマットかチェック
	validFormats := map[string]bool{
		"json": true,
		"csv":  true,
		"sql":  true,
	}
	if !validFormats[options.Format] {
		return fmt.Errorf("サポートされていないフォーマットです: %s", options.Format)
	}

	if options.Limit != nil && *options.Limit <= 0 {
		return fmt.Errorf("limitは正の数である必要があります: %d", *options.Limit)
	}

	return nil
}

// buildQuery はダンプ用のクエリを構築します
func (d *TableDumper) buildQuery(options DumpOptions) string {
	query := fmt.Sprintf("SELECT * FROM %s", options.TableName)

	if options.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *options.Limit)
	}

	return query
}

// fetchTableData はテーブルデータを取得します
func (d *TableDumper) fetchTableData(ctx context.Context, query string) ([]map[string]interface{}, error) {
	rows, err := d.executor.QueryContext(ctx, query)
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
	var result []map[string]interface{}

	// 各行を処理
	for rows.Next() {
		// スキャン用のインターフェースのスライスを作成
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// 行をスキャン
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// 行データをマップに変換
		row := make(map[string]interface{})
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

// ensureOutputDirectory は出力ディレクトリが存在することを確認し、必要に応じて作成します
func (d *TableDumper) ensureOutputDirectory(outputPath string) error {
	if outputPath == "" || outputPath == "." {
		return nil
	}

	return d.fileWriter.MkdirAll(outputPath, 0755)
}

// generateFileName はファイル名を生成します
func (d *TableDumper) generateFileName(tableName, format string) string {
	timestamp := time.Now().Format("20060102_150405")
	extension := d.getFileExtension(format)
	return fmt.Sprintf("%s_%s.%s", tableName, timestamp, extension)
}

// getFileExtension はフォーマットに応じたファイル拡張子を返します
func (d *TableDumper) getFileExtension(format string) string {
	switch format {
	case "json":
		return "json"
	case "csv":
		return "csv"
	case "sql":
		return "sql"
	default:
		return "txt"
	}
}

// writeDataToFile はデータをファイルに書き込みます
func (d *TableDumper) writeDataToFile(filePath string, data []map[string]interface{}, format string, tableName string) error {
	switch format {
	case "json":
		return d.writeJSONFile(filePath, data)
	case "csv":
		return d.writeCSVFile(filePath, data)
	case "sql":
		return d.writeSQLFile(filePath, data, tableName)
	default:
		return fmt.Errorf("サポートされていないフォーマット: %s", format)
	}
}

// writeJSONFile はJSONファイルを書き込みます
func (d *TableDumper) writeJSONFile(filePath string, data []map[string]interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return d.fileWriter.WriteFile(filePath, jsonData, 0644)
}

// writeCSVFile はCSVファイルを書き込みます
func (d *TableDumper) writeCSVFile(filePath string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return d.fileWriter.WriteFile(filePath, []byte(""), 0644)
	}

	file, err := d.fileWriter.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// ヘッダーを書き込み
	var headers []string
	for key := range data[0] {
		headers = append(headers, key)
	}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// データを書き込み
	for _, row := range data {
		var values []string
		for _, header := range headers {
			value := row[header]
			if value == nil {
				values = append(values, "")
			} else {
				values = append(values, fmt.Sprintf("%v", value))
			}
		}
		if err := writer.Write(values); err != nil {
			return err
		}
	}

	return nil
}

// writeSQLFile はSQL INSERT文ファイルを書き込みます
func (d *TableDumper) writeSQLFile(filePath string, data []map[string]interface{}, tableName string) error {
	if len(data) == 0 {
		return d.fileWriter.WriteFile(filePath, []byte("-- No data to export\n"), 0644)
	}

	var sqlBuilder strings.Builder

	// ヘッダーコメント
	sqlBuilder.WriteString(fmt.Sprintf("-- Table dump for %s\n", tableName))
	sqlBuilder.WriteString(fmt.Sprintf("-- Generated at %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// カラム名を取得
	var columns []string
	for key := range data[0] {
		columns = append(columns, key)
	}

	// INSERT文を生成
	for _, row := range data {
		sqlBuilder.WriteString(fmt.Sprintf("INSERT INTO %s (", tableName))
		sqlBuilder.WriteString(strings.Join(columns, ", "))
		sqlBuilder.WriteString(") VALUES (")

		var values []string
		for _, col := range columns {
			value := row[col]
			if value == nil {
				values = append(values, "NULL")
			} else {
				// 文字列の場合はクォートで囲む
				switch v := value.(type) {
				case string:
					values = append(values, fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''")))
				default:
					values = append(values, fmt.Sprintf("%v", v))
				}
			}
		}

		sqlBuilder.WriteString(strings.Join(values, ", "))
		sqlBuilder.WriteString(");\n")
	}

	return d.fileWriter.WriteFile(filePath, []byte(sqlBuilder.String()), 0644)
}

// DumpAllTables はデータベース内の全テーブルをダンプします
func (d *TableDumper) DumpAllTables(ctx context.Context, outputPath, format string, limit *int) (*DumpAllTablesResult, error) {
	// デフォルト値を設定
	if outputPath == "" {
		outputPath = "."
	}
	if format == "" {
		format = "json"
	}

	// 出力ディレクトリを作成
	if err := d.ensureOutputDirectory(outputPath); err != nil {
		return nil, fmt.Errorf("出力ディレクトリ作成エラー: %w", err)
	}

	// 全テーブル一覧を取得
	tables, err := d.getAllTables(ctx)
	if err != nil {
		return nil, fmt.Errorf("テーブル一覧取得エラー: %w", err)
	}

	// データベース名を取得
	databaseName, err := d.getDatabaseName(ctx)
	if err != nil {
		return nil, fmt.Errorf("データベース名取得エラー: %w", err)
	}

	// 結果を格納する構造体を初期化
	result := &DumpAllTablesResult{
		DatabaseName: databaseName,
		TotalTables:  len(tables),
		Results:      []DumpResult{},
		FailedTables: []FailedDump{},
		ExecutedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	// 各テーブルをダンプ
	for _, table := range tables {
		options := DumpOptions{
			TableName:  table.Name,
			OutputPath: outputPath,
			Format:     format,
			Limit:      limit,
		}

		dumpResult, err := d.DumpTable(ctx, options)
		if err != nil {
			// エラーが発生した場合は失敗リストに追加
			result.FailedTables = append(result.FailedTables, FailedDump{
				TableName: table.Name,
				Error:     err.Error(),
			})
		} else {
			// 成功した場合は結果リストに追加
			result.Results = append(result.Results, *dumpResult)
		}
	}

	return result, nil
}

// getAllTables はデータベース内の全テーブル一覧を取得します
func (d *TableDumper) getAllTables(ctx context.Context) ([]Table, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`

	rows, err := d.executor.QueryContext(ctx, query)
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

// getDatabaseName は現在のデータベース名を取得します
func (d *TableDumper) getDatabaseName(ctx context.Context) (string, error) {
	query := "SELECT current_database()"

	var dbName string
	row := d.executor.QueryRowContext(ctx, query)
	if err := row.Scan(&dbName); err != nil {
		return "", err
	}

	return dbName, nil
}
