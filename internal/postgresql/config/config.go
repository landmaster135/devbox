package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

// Config はCLIツールの設定を表します
type Config struct {
	Operation   string // 必須: "dump", "dump-all-tables", "list-tables-minimum", "list-tables"
	DatabaseURL string // 必須: PostgreSQL接続URL
	TableName   string // dump操作時のみ必須: ダンプするテーブル名
	OutputPath  string // オプション: 出力ディレクトリパス
	Format      string // オプション: 出力フォーマット (json, csv, sql, text)
	Limit       *int   // オプション: 最大レコード数
	Concurrency *int   // オプション: 並行処理数 (1-10)
	Help        bool   // ヘルプフラグ
}

// ParseFlags はコマンドライン引数を解析してConfigを返します
func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: dump, dump-all-tables, list-tables-minimum, list-tables)")
	flag.StringVar(&cfg.DatabaseURL, "database-url", "", "PostgreSQLデータベース接続URL (必須)")
	flag.StringVar(&cfg.TableName, "table-name", "", "ダンプするテーブル名 (dump操作時のみ必須)")
	flag.StringVar(&cfg.OutputPath, "output-path", "", "出力ディレクトリパス (オプション、デフォルト: カレントディレクトリ)")
	flag.StringVar(&cfg.Format, "format", "json", "出力フォーマット (オプション: json, csv, sql, text、デフォルト: json)")

	var limitValue int
	flag.IntVar(&limitValue, "limit", 0, "最大レコード数 (オプション)")

	var concurrencyValue int
	// concurrencyの設定とバリデーション
	maxConcurrency := 10
	// デフォルト値を設定
	concurrency := min(runtime.NumCPU(), maxConcurrency)
	flag.IntVar(&concurrencyValue, "concurrency", 1, fmt.Sprintf("並行処理数 (オプション: 1-%d、デフォルト: CPUコア数)", concurrency))
	if concurrencyValue != 1 {
		// 指定された値のバリデーション
		if concurrencyValue < 1 || concurrencyValue > concurrency {
			return nil, fmt.Errorf("--concurrency は1以上かつ%d以下である必要があります: %d", maxConcurrency, concurrencyValue)
		}
		concurrency = concurrencyValue
	}
	cfg.Concurrency = &concurrency

	flag.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	flag.Parse()

	// limitが指定されている場合のみポインタを設定
	if limitValue != 0 {
		cfg.Limit = &limitValue
	}

	// ヘルプが要求された場合は検証をスキップ
	if cfg.Help {
		return cfg, nil
	}

	// 必須パラメータの検証
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate は設定の妥当性を検証します
func (c *Config) validate() error {
	if c.Operation == "" {
		return fmt.Errorf("--operation は必須です")
	}

	// 対応操作の検証
	validOperations := map[string]bool{
		"dump":                true,
		"dump-all-tables":     true,
		"list-tables-minimum": true,
		"list-tables":         true,
	}
	if !validOperations[c.Operation] {
		return fmt.Errorf("未対応の操作です: %s (対応操作: dump, dump-all-tables, list-tables-minimum, list-tables)", c.Operation)
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("--database-url は必須です")
	}

	// table-nameはdump操作時のみ必須
	if c.Operation == "dump" && c.TableName == "" {
		return fmt.Errorf("--table-name は必須です (dump操作時)")
	}

	// フォーマットの検証
	validFormats := map[string]bool{
		"json": true,
		"csv":  true,
		"sql":  true,
		"text": true,
	}
	if !validFormats[c.Format] {
		return fmt.Errorf("未対応のフォーマットです: %s (対応フォーマット: json, csv, sql, text)", c.Format)
	}

	// 操作別のフォーマット制限
	if c.Operation == "list-tables-minimum" && c.Format != "json" {
		return fmt.Errorf("list-tables-minimum操作ではjsonフォーマットのみ対応しています")
	}

	// limitの検証
	if c.Limit != nil && *c.Limit < 0 {
		return fmt.Errorf("--limit は正の数もしくは0である必要があります: %d", *c.Limit)
	}

	return nil
}

// PrintUsage は使用方法を表示します
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "PostgreSQL CLIツール\n\n")
	fmt.Fprintf(os.Stderr, "使用方法:\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli [オプション]\n\n")
	fmt.Fprintf(os.Stderr, "必須オプション:\n")
	fmt.Fprintf(os.Stderr, "  --operation string      実行する操作 (dump, dump-all-tables, list-tables-minimum, list-tables)\n")
	fmt.Fprintf(os.Stderr, "  --database-url string   PostgreSQLデータベース接続URL\n")
	fmt.Fprintf(os.Stderr, "  --table-name string     ダンプするテーブル名 (dump操作時のみ必須)\n\n")
	fmt.Fprintf(os.Stderr, "オプション:\n")
	fmt.Fprintf(os.Stderr, "  --output-path string    出力ディレクトリパス (デフォルト: カレントディレクトリ)\n")
	fmt.Fprintf(os.Stderr, "  --format string         出力フォーマット: json, csv, sql, text (デフォルト: json)\n")
	fmt.Fprintf(os.Stderr, "  --limit int             最大レコード数\n")
	fmt.Fprintf(os.Stderr, "  --concurrency int       並行処理数: 1-10 (デフォルト: CPUコア数、最大10)\n")
	fmt.Fprintf(os.Stderr, "  --help                  このヘルプを表示\n\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  # 単一テーブルダンプ\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump --database-url=\"postgres://user:pass@localhost/db\" --table-name=users\n\n")
	fmt.Fprintf(os.Stderr, "  # 全テーブルダンプ\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump-all-tables --database-url=\"postgres://user:pass@localhost/db\"\n\n")
	fmt.Fprintf(os.Stderr, "  # 全テーブルダンプ（CSV形式、出力先指定）\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump-all-tables --database-url=\"postgres://user:pass@localhost/db\" --format=csv --output-path=/tmp/dumps\n\n")
	fmt.Fprintf(os.Stderr, "  # 全テーブルダンプ（各テーブル最大1000件）\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump-all-tables --database-url=\"postgres://user:pass@localhost/db\" --limit=1000\n\n")
	fmt.Fprintf(os.Stderr, "  # 全テーブルダンプ（並行処理数3で実行）\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump-all-tables --database-url=\"postgres://user:pass@localhost/db\" --concurrency=3\n\n")
	fmt.Fprintf(os.Stderr, "  # テーブル一覧（最小限）\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=list-tables-minimum --database-url=\"postgres://user:pass@localhost/db\"\n\n")
	fmt.Fprintf(os.Stderr, "  # テーブル一覧（詳細）\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=list-tables --database-url=\"postgres://user:pass@localhost/db\"\n\n")
	fmt.Fprintf(os.Stderr, "  # テーブル一覧（テキスト形式）\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=list-tables --database-url=\"postgres://user:pass@localhost/db\" --format=text\n\n")
}
