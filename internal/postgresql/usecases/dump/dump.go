package dump

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	pq "github.com/lib/pq"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	metaFetch "github.com/landmaster135/devbox/internal/postgresql/usecases/meta_fetch"
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

func NewDumpResult(tableName string, recordCount int, outputPath, fileName, format, executedAt string) *DumpResult {
	return &DumpResult{
		TableName:   tableName,
		RecordCount: recordCount,
		OutputPath:  outputPath,
		FileName:    fileName,
		Format:      format,
		ExecutedAt:  executedAt,
	}
}

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
		if part == "" || !metaFetch.IdentifierPattern.MatchString(part) {
			return "", nil, fmt.Errorf("テーブル名に使用できない文字が含まれています: %s", tableName)
		}
		quotedParts[i] = pq.QuoteIdentifier(part)
	}

	return strings.Join(quotedParts, "."), parts, nil
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

	result = NewDumpResult(
		options.TableName,
		recordCount,
		options.OutputPath,
		fileName,
		options.Format,
		time.Now().Format("2006-01-02 15:04:05"),
	)

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

	if len(parts) == 2 && parts[0] != model.DefaultTableSchema {
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
			if err := writer.writeBatch(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := writer.writeBatch(batch); err != nil {
			return err
		}
	}

	return rows.Err()
}

type tableDataWriter interface {
	writeBatch(rows []map[string]any) error
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
