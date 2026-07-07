package dump_binary

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	infrastructures "github.com/landmaster135/devbox/internal/postgresql/infrastructures"
	writer "github.com/landmaster135/devbox/internal/postgresql/usecases/writer"
)

// DumpResult は binary ダンプ処理の結果です。
type DumpResult struct {
	TableName   string `json:"table_name"`
	RecordCount int    `json:"record_count"`
	OutputPath  string `json:"output_path"`
	FileName    string `json:"file_name"`
	Format      string `json:"format"`
	ExecutedAt  string `json:"executed_at"`
}

// Dumper は pg_dump による binary ダンプを実行します。
type Dumper struct {
	commandExecutor infrastructures.CommandExecutor
	fileWriter      writer.FileWriter
	timezone        string
	retryConfig     RetryConfig
}

func NewDumper(timezone string) *Dumper {
	return NewDumperWithDependencies(
		infrastructures.NewOSCommandExecutor(),
		&writer.DefaultFileWriter{},
		timezone,
		RetryConfig{},
	)
}

func NewDumperWithDependencies(commandExecutor infrastructures.CommandExecutor, fileWriter writer.FileWriter, timezone string, retryConfig RetryConfig) *Dumper {
	if commandExecutor == nil {
		commandExecutor = infrastructures.NewOSCommandExecutor()
	}
	if fileWriter == nil {
		fileWriter = &writer.DefaultFileWriter{}
	}

	retryConfig = normalizeRetryConfig(retryConfig)

	return &Dumper{
		commandExecutor: commandExecutor,
		fileWriter:      fileWriter,
		timezone:        strings.TrimSpace(timezone),
		retryConfig:     retryConfig,
	}
}

func (d *Dumper) DumpDatabase(ctx context.Context, databaseURL, outputPath string, excludeTableData []string) (*DumpResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	databaseURL = strings.TrimSpace(databaseURL)
	if databaseURL == "" {
		return nil, errors.New("database-url が設定されていません")
	}
	pgDumpDatabaseURL := sanitizeDatabaseURLForPgDump(databaseURL)

	if outputPath == "" {
		outputPath = "."
	}

	cleanPath := filepath.Clean(outputPath)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("無効なパスが指定されました: %s", outputPath)
	}

	if outputPath != "." {
		if err := d.fileWriter.MkdirAll(outputPath, 0755); err != nil {
			return nil, fmt.Errorf("出力ディレクトリ作成エラー: %w", err)
		}
	}

	fileName := d.generateFileName(databaseURL)
	filePath := filepath.Join(outputPath, fileName)
	args, err := buildPgDumpArgs(pgDumpDatabaseURL, filePath, excludeTableData)
	if err != nil {
		return nil, err
	}

	for attempt := 1; attempt <= d.retryConfig.MaxAttempts; attempt++ {
		output, err := d.commandExecutor.Execute("pg_dump", args...)
		if err == nil {
			return &DumpResult{
				TableName:   "all_tables",
				RecordCount: 0,
				OutputPath:  outputPath,
				FileName:    fileName,
				Format:      "binary",
				ExecutedAt:  d.currentTime().Format("2006-01-02 15:04:05"),
			}, nil
		}

		commandOutput := strings.TrimSpace(string(output))
		if !isRetriablePgDumpError(err, commandOutput) || attempt == d.retryConfig.MaxAttempts {
			if commandOutput == "" {
				return nil, fmt.Errorf("pg_dump の実行に失敗しました (attempts=%d): %w", attempt, err)
			}
			return nil, fmt.Errorf("pg_dump の実行に失敗しました (attempts=%d): %w: %s", attempt, err, commandOutput)
		}

		if err := d.retryConfig.SleepWithContext(ctx, pgDumpRetryDelay(d.retryConfig, attempt)); err != nil {
			return nil, err
		}
	}

	return nil, errors.New("pg_dump の実行に失敗しました")
}

func buildPgDumpArgs(databaseURL, filePath string, excludeTableData []string) ([]string, error) {
	args := []string{"-Fc", "--dbname", databaseURL}
	for _, table := range excludeTableData {
		trimmedTable := strings.TrimSpace(table)
		if trimmedTable == "" {
			return nil, errors.New("exclude-table-data に空のテーブル名は指定できません")
		}
		args = append(args, "--exclude-table-data="+trimmedTable)
	}
	args = append(args, "-f", filePath)
	return args, nil
}

func (d *Dumper) generateFileName(databaseURL string) string {
	name := resolveDatabaseName(databaseURL)
	return fmt.Sprintf("%s_%s.dump", name, d.currentTime().Format("20060102_150405"))
}

func (d *Dumper) currentTime() time.Time {
	if d.timezone == "" {
		return time.Now()
	}
	loc, err := time.LoadLocation(d.timezone)
	if err != nil {
		return time.Now()
	}
	return time.Now().In(loc)
}

func resolveDatabaseName(databaseURL string) string {
	const defaultName = "database"

	u, err := url.Parse(databaseURL)
	if err != nil {
		return defaultName
	}

	name := strings.TrimSpace(strings.TrimPrefix(u.Path, "/"))
	if name == "" {
		return defaultName
	}

	return strings.ReplaceAll(name, "/", "_")
}

func sanitizeDatabaseURLForPgDump(databaseURL string) string {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return databaseURL
	}

	query := u.Query()
	query.Del("statement_cache_mode")
	u.RawQuery = query.Encode()

	return u.String()
}
