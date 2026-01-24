package dump

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	writer "github.com/landmaster135/devbox/internal/postgresql/usecases/writer"
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

// DatabaseQueryExecutor はクエリ実行を抽象化するインターフェースです
type DatabaseQueryExecutor interface {
	QueryContextRows(ctx context.Context, query string, args ...any) (model.RowsInterface, error)
	QueryRowContextRow(ctx context.Context, query string, args ...any) model.RowInterface
}

// #==============================================================#
// ##          TableDumper                                       ##
// #==============================================================#

// TableDumper はテーブルダンプ機能を提供します
type TableDumper struct {
	executor   DatabaseQueryExecutor
	fileWriter writer.FileWriter
}

// NewTableDumper は新しいTableDumperを作成します
func NewTableDumper(executor DatabaseQueryExecutor) *TableDumper {
	return &TableDumper{
		executor:   executor,
		fileWriter: &writer.DefaultFileWriter{},
	}
}

// NewTableDumperWithDependencies はテスト用に依存性を注入できるTableDumperを作成します
func NewTableDumperWithDependencies(executor DatabaseQueryExecutor, fileWriter writer.FileWriter) *TableDumper {
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

	if err := streamRows(rows, columnOrder, writer); err != nil {
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

func (d *TableDumper) OutputResultIntoFile(jsonResult []byte, outputPath, format string) error {
	if err := d.ensureOutputDirectory(outputPath); err != nil {
		return fmt.Errorf("出力ディレクトリ作成エラー: %w", err)
	}

	fileName := d.generateResultFileName(format)
	filePath := filepath.Join(outputPath, fileName)

	file, err := d.fileWriter.Create(filePath)
	if err != nil {
		return fmt.Errorf("結果ファイル作成エラー: %w", err)
	}

	if _, err := d.fileWriter.Write(file, jsonResult); err != nil {
		_ = file.Close()
		return fmt.Errorf("結果ファイル書き込みエラー: %w", err)
	}

	if err := d.fileWriter.Close(file); err != nil {
		return fmt.Errorf("結果ファイルクローズエラー: %w", err)
	}

	return nil
}

func (d *TableDumper) newStreamWriter(format, filePath, tableName string, sortedColumns []string) (writer.TableDataWriter, error) {
	switch format {
	case "json":
		return writer.NewJSONStreamWriter(d.fileWriter, filePath)
	case "csv":
		return writer.NewCSVStreamWriter(d.fileWriter, filePath, sortedColumns)
	case "sql":
		return writer.NewSQLStreamWriter(d.fileWriter, filePath, tableName, sortedColumns)
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

func (d *TableDumper) generateResultFileName(format string) string {
	baseName := "results"
	timestamp := time.Now().Format("20060102_150405")
	extension := d.getFileExtension(format)
	return fmt.Sprintf("%s_%s.%s", timestamp, baseName, extension)
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
	case "markdown":
		return "md"
	default:
		return "txt"
	}
}

// streamRows は取得した行を一定件数ずつ writer に書き込みます
func streamRows(rows model.RowsInterface, columnOrder []string, dataWriter writer.TableDataWriter) error {
	if dataWriter == nil {
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
			if err := dataWriter.WriteBatch(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := dataWriter.WriteBatch(batch); err != nil {
			return err
		}
	}

	return rows.Err()
}
