package dump

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

// #==============================================================#
// ##          Data Structures                                   ##
// #==============================================================#

// DumpAllTablesMinResult は全テーブルダンプの最小限の結果を表します
type DumpAllTablesMinResult struct {
	TotalTables      int    `json:"total_tables"`
	SuccessfulTables int    `json:"successful_tables"`
	FailedTables     int    `json:"failed_tables"`
	ExecutedAt       string `json:"executed_at"`
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

// #==============================================================#
// ##          DumpAllTables                                     ##
// #==============================================================#
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

// DumpAllTablesAndOutputJSON はデータベース内の全テーブルをダンプして結果を指定フォーマットで出力します
func (d *TableDumper) DumpAllTablesAndOutputJSON(ctx context.Context, outputPath, format string, limit *int, concurrency *int, resultFormat, heading string) (string, string, error) {
	resultFormat = strings.ToLower(resultFormat)
	if resultFormat != "json" && resultFormat != "markdown" {
		return "", "", fmt.Errorf("未対応の結果フォーマットです: %s", resultFormat)
	}

	result, err := d.DumpAllTables(ctx, outputPath, format, limit, concurrency)
	if err != nil {
		return "", "", fmt.Errorf("全テーブルダンプの実行に失敗しました: %v", err)
	}

	fullResult, minResult, err := formatDumpAllTablesResult(result, resultFormat, heading)
	if err != nil {
		return "", "", err
	}

	if err := d.OutputResultIntoFile([]byte(fullResult), outputPath, resultFormat); err != nil {
		return "", "", fmt.Errorf("結果のファイル出力に失敗しました: %v", err)
	}

	return fullResult, minResult, nil
}

func formatDumpAllTablesResult(result *DumpAllTablesResult, format, heading string) (string, string, error) {
	switch strings.ToLower(format) {
	case "json":
		return formatDumpAllTablesResultAsJSON(result)
	case "markdown":
		full, min := formatDumpAllTablesResultAsMarkdown(result, heading)
		return full, min, nil
	default:
		return "", "", fmt.Errorf("サポートされていない結果フォーマットです: %s", format)
	}
}

func formatDumpAllTablesResultAsJSON(result *DumpAllTablesResult) (string, string, error) {
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("結果のJSON変換に失敗しました: %v", err)
	}

	minResult := &DumpAllTablesMinResult{
		TotalTables:      result.TotalTables,
		SuccessfulTables: len(result.Results),
		FailedTables:     len(result.FailedTables),
		ExecutedAt:       result.ExecutedAt,
	}
	minJSONResult, err := json.MarshalIndent(minResult, "", "  ")
	if err != nil {
		return "", "", fmt.Errorf("最小限の結果のJSON変換に失敗しました: %v", err)
	}

	return string(jsonResult), string(minJSONResult), nil
}

func formatDumpAllTablesMinResultAsMarkdown(result *DumpAllTablesResult, title string) string {
	successCount := len(result.Results)
	failureCount := len(result.FailedTables)

	minBuilder := &strings.Builder{}
	fmt.Fprintf(minBuilder, "## %s\n\n", title)
	minBuilder.WriteString("```markdown\n")
	minBuilder.WriteString("| 項目 | 値 |\n")
	minBuilder.WriteString("| --- | --- |\n")
	fmt.Fprintf(minBuilder, "| Total tables | %d |\n", result.TotalTables)
	fmt.Fprintf(minBuilder, "| Successful tables | %d |\n", successCount)
	fmt.Fprintf(minBuilder, "| Failed tables | %d |\n", failureCount)
	fmt.Fprintf(minBuilder, "| Executed at | %s |", result.ExecutedAt)
	minBuilder.WriteString("```\n")

	return strings.TrimSpace(minBuilder.String())
}

func formatDumpAllTablesMainResultAsMarkdown(result *DumpAllTablesResult, title string) string {
	successCount := len(result.Results)
	failureCount := len(result.FailedTables)

	builder := &strings.Builder{}
	fmt.Fprintf(builder, "## %s\n\n", title)
	builder.WriteString("```markdown\n")
	builder.WriteString("| 項目 | 値 |\n")
	builder.WriteString("| --- | --- |\n")
	fmt.Fprintf(builder, "| Database | `%s` |\n", result.DatabaseName)
	fmt.Fprintf(builder, "| Total tables discovered | %d |\n", result.TotalTables)
	fmt.Fprintf(builder, "| Successful dumps | %d |\n", successCount)
	fmt.Fprintf(builder, "| Failed dumps | %d |\n", failureCount)
	fmt.Fprintf(builder, "| Executed at | %s |\n", result.ExecutedAt)
	builder.WriteString("```\n")

	builder.WriteString("\n### Successful Tables\n")
	builder.WriteString("```markdown\n")
	builder.WriteString("| Table | Rows | File | Format |\n")
	builder.WriteString("| --- | --- | --- | --- |\n")
	if successCount > 0 {
		for _, res := range result.Results {
			fmt.Fprintf(builder, "| `%s` | %d | %s | %s |\n", res.TableName, res.RecordCount, res.FileName, res.Format)
		}
	} else {
		builder.WriteString("| なし | 0 | - | - |\n")
	}
	builder.WriteString("```\n")

	builder.WriteString("\n### Failed Tables\n")
	builder.WriteString("```markdown\n")
	builder.WriteString("| Table | Error |\n")
	builder.WriteString("| --- | --- |\n")
	if failureCount > 0 {
		for _, failed := range result.FailedTables {
			fmt.Fprintf(builder, "| `%s` | %s |\n", failed.TableName, failed.Error)
		}
	} else {
		builder.WriteString("| なし | - |\n")
	}
	builder.WriteString("```\n")

	return strings.TrimSpace(builder.String())
}

func formatDumpAllTablesResultAsMarkdown(result *DumpAllTablesResult, heading string) (string, string) {

	title := strings.TrimSpace(heading)
	if title == "" {
		title = "PostgreSQL Dump Report"
	}

	mainBuilderString := formatDumpAllTablesMainResultAsMarkdown(result, heading)
	minBuilderString := formatDumpAllTablesMinResultAsMarkdown(result, heading)

	return strings.TrimSpace(mainBuilderString), minBuilderString
}
