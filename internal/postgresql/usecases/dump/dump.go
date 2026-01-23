package usecases

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	pq "github.com/lib/pq"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

var identifierPattern = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

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

func NewDumpOptions(tableName, outputPath, format string, limit *int) *DumpOptions {
	return &DumpOptions{
		TableName:  tableName,
		OutputPath: outputPath,
		Format:     format,
		Limit:      limit,
	}
}

// DumpResult はダンプ処理の結果を表します
type DumpResult struct {
	TableName   string `json:"table_name"`
	RecordCount int    `json:"record_count"`
	OutputPath  string `json:"output_path"`
	FileName    string `json:"file_name"`
	Format      string `json:"format"`
	ExecutedAt  string `json:"executed_at"`
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

// DumpTask は並行処理用のダンプタスクを表します
type DumpTask struct {
	Table   model.Table
	Options *DumpOptions
}

// DumpTaskResult は並行処理用のダンプタスク結果を表します
type DumpTaskResult struct {
	Success bool
	Result  *DumpResult
	Failed  *FailedDump
}

const defaultTableSchema = "public"

func quoteQualifiedTableName(tableName string) (string, []string, error) {
	if tableName == "" {
		return "", nil, errors.New("テーブル名が指定されていません")
	}

	parts := strings.Split(tableName, ".")
	if len(parts) > 2 {
		return "", nil, fmt.Errorf("サポートされていないテーブル識別子です: %s", tableName)
	}

	quotedParts := make([]string, len(parts))
	for i, part := range parts {
		if part == "" || !identifierPattern.MatchString(part) {
			return "", nil, fmt.Errorf("テーブル名に使用できない文字が含まれています: %s", tableName)
		}
		quotedParts[i] = pq.QuoteIdentifier(part)
	}

	return strings.Join(quotedParts, "."), parts, nil
}

func qualifyTableIdentifier(tableName string) (qualified string, schema string, name string, err error) {
	if tableName == "" {
		return "", "", "", errors.New("テーブル名が指定されていません")
	}

	parts := strings.Split(tableName, ".")
	if len(parts) > 2 {
		return "", "", "", fmt.Errorf("サポートされていないテーブル識別子です: %s", tableName)
	}

	schema = defaultTableSchema
	name = parts[len(parts)-1]
	if len(parts) == 2 {
		schema = parts[0]
	}

	if schema == "" {
		schema = defaultTableSchema
	}

	if !identifierPattern.MatchString(schema) {
		return "", "", "", fmt.Errorf("スキーマ名に使用できない文字が含まれています: %s", schema)
	}
	if !identifierPattern.MatchString(name) {
		return "", "", "", fmt.Errorf("テーブル名に使用できない文字が含まれています: %s", tableName)
	}

	qualified = fmt.Sprintf("%s.%s", pq.QuoteIdentifier(schema), pq.QuoteIdentifier(name))
	return qualified, schema, name, nil
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

// DatabaseQueryExecutor はクエリ実行を抽象化するインターフェースです
type DatabaseQueryExecutor interface {
	QueryContextRows(ctx context.Context, query string, args ...any) (model.RowsInterface, error)
	QueryRowContextRow(ctx context.Context, query string, args ...any) model.RowInterface
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
	executor   DatabaseQueryExecutor
	fileWriter FileWriter
}

// NewTableDumper は新しいTableDumperを作成します
func NewTableDumper(executor DatabaseQueryExecutor) *TableDumper {
	return &TableDumper{
		executor:   executor,
		fileWriter: &DefaultFileWriter{},
	}
}

// NewTableDumperWithDependencies はテスト用に依存性を注入できるTableDumperを作成します
func NewTableDumperWithDependencies(executor DatabaseQueryExecutor, fileWriter FileWriter) *TableDumper {
	return &TableDumper{
		executor:   executor,
		fileWriter: fileWriter,
	}
}

// DumpTable はテーブルの全レコードをダンプします
func (d *TableDumper) DumpTable(ctx context.Context, options *DumpOptions) (result *DumpResult, err error) {
	if options == nil {
		return nil, errors.New("オプションが指定されていません")
	}

	// パラメータ検証
	if err := d.validateOptions(options); err != nil {
		return nil, fmt.Errorf("オプション検証エラー: %w", err)
	}

	if err := d.ensureAllowedTable(ctx, options.TableName); err != nil {
		return nil, fmt.Errorf("テーブル検証エラー: %w", err)
	}

	// 出力ディレクトリを作成
	if err := d.ensureOutputDirectory(options.OutputPath); err != nil {
		return nil, fmt.Errorf("出力ディレクトリ作成エラー: %w", err)
	}

	// クエリを構築
	query, err := d.buildQuery(options)
	if err != nil {
		return nil, fmt.Errorf("クエリ構築エラー: %w", err)
	}

	rows, err := d.executor.QueryContextRows(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("データ取得エラー: %w", err)
	}
	defer rows.Close()

	columnOrder, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("カラム情報取得エラー: %w", err)
	}

	sortedColumns := make([]string, len(columnOrder))
	copy(sortedColumns, columnOrder)
	sort.Strings(sortedColumns)

	fileName := d.generateFileName(options.TableName, options.Format)
	filePath := filepath.Join(options.OutputPath, fileName)

	writer, err := d.newStreamWriter(options.Format, filePath, options.TableName, sortedColumns)
	if err != nil {
		return nil, fmt.Errorf("ファイル初期化エラー: %w", err)
	}

	defer func() {
		closeErr := writer.Close()
		if closeErr != nil {
			if err == nil {
				err = closeErr
			}
			result = nil
		}
	}()

	if err := d.streamRows(rows, columnOrder, writer); err != nil {
		return nil, fmt.Errorf("データ書き込みエラー: %w", err)
	}

	recordCount := writer.RowsWritten()

	result = &DumpResult{
		TableName:   options.TableName,
		RecordCount: recordCount,
		OutputPath:  options.OutputPath,
		FileName:    fileName,
		Format:      options.Format,
		ExecutedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	return result, nil
}

// validateOptions はオプションを検証します
func (d *TableDumper) validateOptions(options *DumpOptions) error {
	if options == nil {
		return errors.New("オプションが指定されていません")
	}

	if _, _, err := quoteQualifiedTableName(options.TableName); err != nil {
		return err
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

// ensureAllowedTable はテーブルが存在し、ホワイトリストに一致することを確認します
func (d *TableDumper) ensureAllowedTable(ctx context.Context, tableName string) error {
	_, parts, err := quoteQualifiedTableName(tableName)
	if err != nil {
		return err
	}

	if len(parts) == 2 && parts[0] != "public" {
		return fmt.Errorf("サポートされていないスキーマが指定されました: %s", parts[0])
	}

	tables, err := d.getAllTables(ctx)
	if err != nil {
		return fmt.Errorf("テーブル一覧取得エラー: %w", err)
	}

	targetName := parts[len(parts)-1]
	for _, tbl := range tables {
		if tbl.Name == targetName {
			return nil
		}
	}

	return fmt.Errorf("指定されたテーブルが存在しません: %s", tableName)
}

// buildQuery はダンプ用のクエリを構築します
func (d *TableDumper) buildQuery(options *DumpOptions) (string, error) {
	if options == nil {
		return "", errors.New("オプションが指定されていません")
	}

	quotedTable, _, err := quoteQualifiedTableName(options.TableName)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("SELECT * FROM %s", quotedTable)

	if options.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", *options.Limit)
	}

	return query, nil
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

// streamRows は取得した行を一定件数ずつ writer に書き込みます
func (d *TableDumper) streamRows(rows model.RowsInterface, columnOrder []string, writer tableDataWriter) error {
	if writer == nil {
		return errors.New("writerが初期化されていません")
	}

	batch := make([]map[string]any, 0, 1000)

	for rows.Next() {
		values := make([]any, len(columnOrder))
		valuePtrs := make([]any, len(columnOrder))
		for i := range columnOrder {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return err
		}

		row := make(map[string]any, len(columnOrder))
		for i, col := range columnOrder {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		batch = append(batch, row)
		if len(batch) == 1000 {
			if err := writer.WriteBatch(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := writer.WriteBatch(batch); err != nil {
			return err
		}
	}

	return rows.Err()
}

type tableDataWriter interface {
	WriteBatch(rows []map[string]any) error
	Close() error
	RowsWritten() int
}

func (d *TableDumper) newStreamWriter(format, filePath, tableName string, sortedColumns []string) (tableDataWriter, error) {
	switch format {
	case "json":
		return newJSONStreamWriter(d.fileWriter, filePath)
	case "csv":
		return newCSVStreamWriter(d.fileWriter, filePath, sortedColumns)
	case "sql":
		return newSQLStreamWriter(d.fileWriter, filePath, tableName, sortedColumns)
	default:
		return nil, fmt.Errorf("サポートされていないフォーマット: %s", format)
	}
}

func formatCSVValue(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		if len(v) == 0 {
			return ""
		}
		return base64.StdEncoding.EncodeToString(v)
	case time.Time:
		return v.Format(time.RFC3339Nano)
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(value)
	}
}

func formatSQLValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "NULL"
	case string:
		return pq.QuoteLiteral(v)
	case []byte:
		if len(v) == 0 {
			return "decode('', 'hex')"
		}
		hexStr := hex.EncodeToString(v)
		return fmt.Sprintf("decode('%s','hex')", hexStr)
	case time.Time:
		return pq.QuoteLiteral(v.Format(time.RFC3339Nano))
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case fmt.Stringer:
		return pq.QuoteLiteral(v.String())
	default:
		return pq.QuoteLiteral(fmt.Sprint(value))
	}
}

type jsonStreamWriter struct {
	file     *os.File
	rows     int
	closed   bool
	closeErr error
}

func newJSONStreamWriter(fileWriter FileWriter, filePath string) (*jsonStreamWriter, error) {
	file, err := fileWriter.Create(filePath)
	if err != nil {
		return nil, err
	}

	if _, err := file.WriteString("["); err != nil {
		_ = file.Close()
		return nil, err
	}

	return &jsonStreamWriter{file: file}, nil
}

func (w *jsonStreamWriter) WriteBatch(rows []map[string]any) error {
	if w.closed {
		return errors.New("既にクローズされたライターに書き込めません")
	}

	for _, row := range rows {
		rowJSON, err := json.MarshalIndent(row, "", "  ")
		if err != nil {
			return err
		}

		prefix := "\n"
		if w.rows > 0 {
			prefix = ",\n"
		}

		if _, err := w.file.WriteString(prefix); err != nil {
			return err
		}
		if _, err := w.file.WriteString("  "); err != nil {
			return err
		}
		if _, err := w.file.WriteString(strings.ReplaceAll(string(rowJSON), "\n", "\n  ")); err != nil {
			return err
		}

		w.rows++
	}

	return nil
}

func (w *jsonStreamWriter) Close() error {
	if w.closed {
		return w.closeErr
	}
	w.closed = true

	if w.file == nil {
		return nil
	}

	var err error
	if w.rows > 0 {
		_, err = w.file.WriteString("\n]")
	} else {
		_, err = w.file.WriteString("]")
	}

	if err != nil {
		w.closeErr = err
		_ = w.file.Close()
		return err
	}

	if closeErr := w.file.Close(); closeErr != nil {
		err = closeErr
	}

	w.closeErr = err
	return err
}

func (w *jsonStreamWriter) RowsWritten() int {
	return w.rows
}

type csvStreamWriter struct {
	file     *os.File
	writer   *csv.Writer
	headers  []string
	rows     int
	closed   bool
	closeErr error
}

func newCSVStreamWriter(fileWriter FileWriter, filePath string, headers []string) (*csvStreamWriter, error) {
	file, err := fileWriter.Create(filePath)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)

	if len(headers) > 0 {
		if err := writer.Write(headers); err != nil {
			writer.Flush()
			_ = file.Close()
			return nil, err
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			_ = file.Close()
			return nil, err
		}
	}

	return &csvStreamWriter{
		file:    file,
		writer:  writer,
		headers: headers,
	}, nil
}

func (w *csvStreamWriter) WriteBatch(rows []map[string]any) error {
	if w.closed {
		return errors.New("既にクローズされたライターに書き込めません")
	}

	for _, row := range rows {
		values := make([]string, len(w.headers))
		for i, header := range w.headers {
			values[i] = formatCSVValue(row[header])
		}

		if err := w.writer.Write(values); err != nil {
			return err
		}
		w.rows++
	}

	w.writer.Flush()
	if err := w.writer.Error(); err != nil {
		return err
	}

	return nil
}

func (w *csvStreamWriter) Close() error {
	if w.closed {
		return w.closeErr
	}
	w.closed = true

	if w.writer != nil {
		w.writer.Flush()
		if err := w.writer.Error(); err != nil {
			w.closeErr = err
			_ = w.file.Close()
			return err
		}
	}

	if w.file != nil {
		if err := w.file.Close(); err != nil {
			w.closeErr = err
			return err
		}
	}

	return nil
}

func (w *csvStreamWriter) RowsWritten() int {
	return w.rows
}

type sqlStreamWriter struct {
	file          *os.File
	tableName     string
	quotedTable   string
	columns       []string
	quotedColumns []string
	headerWritten bool
	rows          int
	closed        bool
	closeErr      error
	generatedAt   time.Time
}

func newSQLStreamWriter(fileWriter FileWriter, filePath, tableName string, columns []string) (*sqlStreamWriter, error) {
	file, err := fileWriter.Create(filePath)
	if err != nil {
		return nil, err
	}

	quotedTable, _, err := quoteQualifiedTableName(tableName)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	quotedColumns := make([]string, len(columns))
	for i, col := range columns {
		quotedColumns[i] = pq.QuoteIdentifier(col)
	}

	return &sqlStreamWriter{
		file:          file,
		tableName:     tableName,
		quotedTable:   quotedTable,
		columns:       columns,
		quotedColumns: quotedColumns,
		generatedAt:   time.Now(),
	}, nil
}

func (w *sqlStreamWriter) WriteBatch(rows []map[string]any) error {
	if w.closed {
		return errors.New("既にクローズされたライターに書き込めません")
	}

	if len(rows) == 0 {
		return nil
	}

	if !w.headerWritten {
		if _, err := fmt.Fprintf(w.file, "-- Table dump for %s\n", w.tableName); err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w.file, "-- Generated at %s\n\n", w.generatedAt.Format("2006-01-02 15:04:05")); err != nil {
			return err
		}
		w.headerWritten = true
	}

	for _, row := range rows {
		if _, err := fmt.Fprintf(w.file, "INSERT INTO %s (%s) VALUES (", w.quotedTable, strings.Join(w.quotedColumns, ", ")); err != nil {
			return err
		}

		values := make([]string, len(w.columns))
		for i, col := range w.columns {
			values[i] = formatSQLValue(row[col])
		}

		if _, err := w.file.WriteString(strings.Join(values, ", ")); err != nil {
			return err
		}
		if _, err := w.file.WriteString(");\n"); err != nil {
			return err
		}

		w.rows++
	}

	return nil
}

func (w *sqlStreamWriter) Close() error {
	if w.closed {
		return w.closeErr
	}
	w.closed = true

	if w.file == nil {
		return nil
	}

	if w.rows == 0 {
		if _, err := w.file.WriteString("-- No data to export\n"); err != nil {
			w.closeErr = err
			_ = w.file.Close()
			return err
		}
	}

	if err := w.file.Close(); err != nil {
		w.closeErr = err
		return err
	}

	return nil
}

func (w *sqlStreamWriter) RowsWritten() int {
	return w.rows
}

// #==============================================================#
// ##          DumpAllTables                                     ##
// #==============================================================#
// getAllTables はデータベース内の全テーブル一覧を取得します
func (d *TableDumper) getAllTables(ctx context.Context) ([]model.Table, error) {
	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = 'public'
		ORDER BY table_name
	`

	rows, err := d.executor.QueryContextRows(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []model.Table
	for rows.Next() {
		var table model.Table
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
	row := d.executor.QueryRowContextRow(ctx, query)
	if err := row.Scan(&dbName); err != nil {
		return "", err
	}

	return dbName, nil
}

// dumpSingleTable は単一テーブルのダンプを実行するサブルーチンです
func (d *TableDumper) dumpSingleTable(ctx context.Context, task DumpTask) DumpTaskResult {
	dumpResult, err := d.DumpTable(ctx, task.Options)
	if err != nil {
		return DumpTaskResult{
			Success: false,
			Failed: &FailedDump{
				TableName: task.Table.Name,
				Error:     err.Error(),
			},
		}
	}

	return DumpTaskResult{
		Success: true,
		Result:  dumpResult,
	}
}

// createInitialResult は初期の DumpAllTablesResult 構造体を作成します
func (d *TableDumper) createInitialResult(databaseName string, totalTables int) *DumpAllTablesResult {
	return &DumpAllTablesResult{
		DatabaseName: databaseName,
		TotalTables:  totalTables,
		Results:      []DumpResult{},
		FailedTables: []FailedDump{},
		ExecutedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
}

// collectDumpResults は並行処理の結果を収集して result に集約します
func (d *TableDumper) collectDumpResults(resultChan <-chan DumpTaskResult, result *DumpAllTablesResult) {
	var mu sync.Mutex
	for taskResult := range resultChan {
		mu.Lock()
		if taskResult.Success {
			result.Results = append(result.Results, *taskResult.Result)
		} else {
			result.FailedTables = append(result.FailedTables, *taskResult.Failed)
		}
		mu.Unlock()
	}
}

// executeConcurrentDumps は並行処理でテーブルダンプを実行します
func (d *TableDumper) executeConcurrentDumps(ctx context.Context, databaseName string, tables []model.Table, concurrency int, outputPath, format string, limit *int) (*DumpAllTablesResult, error) {
	if concurrency <= 0 {
		concurrency = 1
	}

	// 結果を格納する構造体を初期化
	result := d.createInitialResult(databaseName, len(tables))

	// テーブルが0個の場合は早期リターン
	if len(tables) == 0 {
		return result, nil
	}

	// 並行処理数を決定
	maxConcurrency := min(concurrency, len(tables))
	if maxConcurrency <= 0 {
		maxConcurrency = 1
	}

	// 並行処理用のチャネルを作成
	taskChan := make(chan DumpTask, len(tables))
	resultChan := make(chan DumpTaskResult, len(tables))

	// ワーカーgoroutineを起動
	var wg sync.WaitGroup
	for range maxConcurrency {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChan {
				select {
				case <-ctx.Done():
					// コンテキストがキャンセルされた場合
					resultChan <- DumpTaskResult{
						Success: false,
						Failed: &FailedDump{
							TableName: task.Table.Name,
							Error:     "コンテキストがキャンセルされました",
						},
					}
					return
				default:
					// 単一テーブルダンプを実行
					taskResult := d.dumpSingleTable(ctx, task)
					resultChan <- taskResult
				}
			}
		}()
	}

	// タスクをチャネルに送信
	for _, table := range tables {
		options := &DumpOptions{
			TableName:  table.Name,
			OutputPath: outputPath,
			Format:     format,
			Limit:      limit,
		}
		taskChan <- DumpTask{
			Table:   table,
			Options: options,
		}
	}
	close(taskChan)

	// 結果を収集するgoroutineを起動
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 結果を収集
	d.collectDumpResults(resultChan, result)

	return result, nil
}

// DumpAllTables はデータベース内の全テーブルを並行処理でダンプします
func (d *TableDumper) DumpAllTables(ctx context.Context, outputPath, format string, limit *int, concurrency *int) (*DumpAllTablesResult, error) {
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

	// 並行処理数のデフォルトと下限を補完
	effectiveConcurrency := 1
	if concurrency != nil && *concurrency > 0 {
		effectiveConcurrency = *concurrency
	}

	// 並行処理でダンプを実行
	return d.executeConcurrentDumps(ctx, databaseName, tables, effectiveConcurrency, outputPath, format, limit)
}
