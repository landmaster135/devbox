package dump

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
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

type tableDataWriter interface {
	writeBatch(rows []map[string]any) error
	Close() error
	RowsWritten() int
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

	if err := d.ensureOptions(ctx, options); err != nil {
		return nil, fmt.Errorf("オプションによる確保に失敗しました: %w", err)
	}

	rows, err := d.queryRows(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("データの取得に失敗しました: %w", err)
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
