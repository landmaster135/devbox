package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	OperationListTables = "list-tables"
	FormatText          = "text"
	FormatJSON          = "json"
)

// Config は sqlite CLI の設定を保持します。
type Config struct {
	Operation string
	DBPath    string
	Format    string
	Help      bool
}

// ParseFlags は CLI 引数を解析します。
func ParseFlags() (*Config, error) {
	return ParseFlagsFromArgs(os.Args[1:])
}

// ParseFlagsFromArgs は指定引数を解析します。
func ParseFlagsFromArgs(args []string) (*Config, error) {
	flagSet := flag.NewFlagSet("sqlite", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	var (
		operation    string
		misspelledOp string
		dbPath       string
		format       string
		help         bool
	)

	flagSet.StringVar(&operation, "operation", "", "実行する操作 (list-tables)")
	flagSet.StringVar(&operation, "o", "", "実行する操作の短縮形")
	flagSet.StringVar(&misspelledOp, "opearation", "", "実行する操作 (deprecated: operation の誤記互換)")
	flagSet.StringVar(&dbPath, "db-path", "", "SQLite ファイルのパス")
	flagSet.StringVar(&format, "format", FormatText, "出力形式 (text, json)")
	flagSet.BoolVar(&help, "help", false, "ヘルプを表示")
	flagSet.BoolVar(&help, "h", false, "ヘルプを表示")

	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if help {
		return &Config{Help: true}, nil
	}

	if strings.TrimSpace(operation) == "" {
		operation = strings.TrimSpace(misspelledOp)
	}

	return NewConfig(operation, dbPath, format)
}

// NewConfig は Config を生成して妥当性を検証します。
func NewConfig(operation, dbPath, format string) (*Config, error) {
	op := strings.TrimSpace(operation)
	if op == "" {
		return nil, fmt.Errorf("--operation は必須です")
	}
	if op != OperationListTables {
		return nil, fmt.Errorf("未対応の operation です: %s", op)
	}

	path := strings.TrimSpace(dbPath)
	if path == "" {
		return nil, fmt.Errorf("--db-path は必須です")
	}

	fmtType := strings.TrimSpace(format)
	if fmtType == "" {
		fmtType = FormatText
	}
	if fmtType != FormatText && fmtType != FormatJSON {
		return nil, fmt.Errorf("--format は text または json を指定してください: %s", fmtType)
	}

	return &Config{
		Operation: op,
		DBPath:    path,
		Format:    fmtType,
	}, nil
}

// PrintUsage は標準エラー出力に利用方法を表示します。
func PrintUsage() {
	PrintUsageTo(os.Stderr)
}

// PrintUsageTo は指定 writer に利用方法を表示します。
func PrintUsageTo(w io.Writer) {
	fmt.Fprintf(w, "SQLite CLIツール\n\n")
	fmt.Fprintf(w, "使用方法:\n")
	fmt.Fprintf(w, "  go run ./cmd/cli/sqlite --operation=list-tables --db-path=./sample.db\n")
	fmt.Fprintf(w, "  go run ./cmd/cli/sqlite --opearation=list-tables --db-path=./sample.db  # 互換フラグ\n\n")
	fmt.Fprintf(w, "オプション:\n")
	fmt.Fprintf(w, "  --operation, -o  実行する操作 (list-tables) [必須]\n")
	fmt.Fprintf(w, "  --opearation     実行する操作 (deprecated: operation の誤記互換)\n")
	fmt.Fprintf(w, "  --db-path        SQLite ファイルのパス [必須]\n")
	fmt.Fprintf(w, "  --format         出力形式 (text, json) [デフォルト: text]\n")
	fmt.Fprintf(w, "  --help, -h       このヘルプを表示\n")
}
