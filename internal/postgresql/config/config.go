package config

import (
	"flag"
	"fmt"
	"os"
)

// Config はCLIツールの設定を表します
type Config struct {
	Operation   string // 必須: "dump"
	DatabaseURL string // 必須: PostgreSQL接続URL
	TableName   string // 必須: ダンプするテーブル名
	OutputPath  string // オプション: 出力ディレクトリパス
	Format      string // オプション: 出力フォーマット (json, csv, sql)
	Limit       *int   // オプション: 最大レコード数
	Help        bool   // ヘルプフラグ
}

// ParseFlags はコマンドライン引数を解析してConfigを返します
func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: dump)")
	flag.StringVar(&cfg.DatabaseURL, "database-url", "", "PostgreSQLデータベース接続URL (必須)")
	flag.StringVar(&cfg.TableName, "table-name", "", "ダンプするテーブル名 (必須)")
	flag.StringVar(&cfg.OutputPath, "output-path", "", "出力ディレクトリパス (オプション、デフォルト: カレントディレクトリ)")
	flag.StringVar(&cfg.Format, "format", "json", "出力フォーマット (オプション: json, csv, sql、デフォルト: json)")

	var limitValue int
	flag.IntVar(&limitValue, "limit", 0, "最大レコード数 (オプション)")

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

	if c.Operation != "dump" {
		return fmt.Errorf("未対応の操作です: %s (対応操作: dump)", c.Operation)
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("--database-url は必須です")
	}

	if c.TableName == "" {
		return fmt.Errorf("--table-name は必須です")
	}

	// フォーマットの検証
	validFormats := map[string]bool{
		"json": true,
		"csv":  true,
		"sql":  true,
	}
	if !validFormats[c.Format] {
		return fmt.Errorf("未対応のフォーマットです: %s (対応フォーマット: json, csv, sql)", c.Format)
	}

	// limitの検証
	if c.Limit != nil && *c.Limit < 0 {
		return fmt.Errorf("--limit は正の数である必要があります: %d", *c.Limit)
	}

	return nil
}

// PrintUsage は使用方法を表示します
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "PostgreSQL CLIツール\n\n")
	fmt.Fprintf(os.Stderr, "使用方法:\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli [オプション]\n\n")
	fmt.Fprintf(os.Stderr, "必須オプション:\n")
	fmt.Fprintf(os.Stderr, "  --operation string      実行する操作 (dump)\n")
	fmt.Fprintf(os.Stderr, "  --database-url string   PostgreSQLデータベース接続URL\n")
	fmt.Fprintf(os.Stderr, "  --table-name string     ダンプするテーブル名\n\n")
	fmt.Fprintf(os.Stderr, "オプション:\n")
	fmt.Fprintf(os.Stderr, "  --output-path string    出力ディレクトリパス (デフォルト: カレントディレクトリ)\n")
	fmt.Fprintf(os.Stderr, "  --format string         出力フォーマット: json, csv, sql (デフォルト: json)\n")
	fmt.Fprintf(os.Stderr, "  --limit int             最大レコード数\n")
	fmt.Fprintf(os.Stderr, "  --help                  このヘルプを表示\n\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  # 基本的な使用例\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump --database-url=\"postgres://user:pass@localhost/db\" --table-name=users\n\n")
	fmt.Fprintf(os.Stderr, "  # 全オプション指定\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump --database-url=\"postgres://user:pass@localhost/db\" --table-name=users --format=json --output-path=/tmp --limit=100\n\n")
}
