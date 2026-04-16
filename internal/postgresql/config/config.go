package config

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

// Config はCLIツールの設定を表します
type Config struct {
	Operation     string // 必須: "dump", "dump-all-tables", "list-tables-minimum", "list-tables"
	DatabaseURL   string // 必須: PostgreSQL接続URL
	TableName     string // dump操作時のみ必須: ダンプするテーブル名
	OutputPath    string // オプション: 出力ディレクトリパス
	Format        string // オプション: 出力フォーマット (json, csv, sql, text, binary)
	Limit         *int   // オプション: 最大レコード数
	Concurrency   *int   // オプション: 並行処理数 (1-CPUコア数)
	ResultFormat  string // オプション: ダンプ結果のフォーマット (json, markdown)
	ResultHeading string // オプション: ダンプ結果の見出し (Markdown出力用)
	Timezone      string // オプション: タイムスタンプ・ファイル名で使用するタイムゾーン
	Help          bool   // ヘルプフラグ
}

// ParseFlags はコマンドライン引数を解析してConfigを返します
func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: dump, dump-all-tables, list-tables-minimum, list-tables)")
	flag.StringVar(&cfg.DatabaseURL, "database-url", "", "PostgreSQLデータベース接続URL (必須)")
	flag.StringVar(&cfg.TableName, "table-name", "", "ダンプするテーブル名 (dump操作時のみ必須)")
	flag.StringVar(&cfg.OutputPath, "output-path", "", "出力ディレクトリパス (オプション、デフォルト: カレントディレクトリ)")
	flag.StringVar(&cfg.Format, "format", "json", "出力フォーマット (オプション: json, csv, sql, text, binary、デフォルト: json)")
	flag.StringVar(&cfg.ResultFormat, "result-format", "json", "ダンプ結果のフォーマット (オプション: json, markdown、デフォルト: json)")
	flag.StringVar(&cfg.ResultHeading, "result-heading", "", "ダンプ結果サマリの見出し (Markdown出力時に利用)")
	flag.StringVar(&cfg.Timezone, "timezone", "", "タイムゾーン (例: Asia/Tokyo)。未指定の場合はシステムローカルを使用")

	var limitValue int
	flag.IntVar(&limitValue, "limit", 0, "最大レコード数 (オプション)")

	maxConcurrency := runtime.NumCPU()
	var concurrencyValue int
	flag.IntVar(&concurrencyValue, "concurrency", maxConcurrency, fmt.Sprintf("並行処理数 (オプション: 1-%d、デフォルト: CPUコア数)", maxConcurrency))

	flag.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	flag.Parse()

	// limitが指定されている場合のみポインタを設定
	if limitValue != 0 {
		cfg.Limit = &limitValue
	}

	// concurrencyの設定とバリデーション
	if concurrencyValue < 1 || concurrencyValue > maxConcurrency {
		return nil, fmt.Errorf("--concurrency は1以上かつ%d以下である必要があります: %d", maxConcurrency, concurrencyValue)
	}
	cfg.Concurrency = &concurrencyValue

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

	// フォーマットの検証
	validFormats := map[string]bool{
		"json":   true,
		"csv":    true,
		"sql":    true,
		"text":   true,
		"binary": true,
	}
	if !validFormats[c.Format] {
		return fmt.Errorf("未対応のフォーマットです: %s (対応フォーマット: json, csv, sql, text, binary)", c.Format)
	}

	// table-nameはdump操作時のみ必須
	if c.Operation == "dump" && c.TableName == "" {
		return fmt.Errorf("--table-name は必須です (dump操作時)")
	}

	// result-format の検証
	validResultFormats := map[string]bool{
		"json":     true,
		"markdown": true,
	}
	if !validResultFormats[c.ResultFormat] {
		return fmt.Errorf("未対応の結果フォーマットです: %s (対応フォーマット: json, markdown)", c.ResultFormat)
	}

	// 操作別のフォーマット制限
	if c.Operation == "list-tables-minimum" && c.Format != "json" {
		return fmt.Errorf("list-tables-minimum操作ではjsonフォーマットのみ対応しています")
	}
	if c.Format == "binary" && c.Operation != "dump-all-tables" {
		return fmt.Errorf("binaryフォーマットはdump-all-tables操作でのみ対応しています")
	}
	if c.Operation == "dump-all-tables" && c.Format == "binary" && c.TableName != "" {
		return fmt.Errorf("--table-name は指定できません (dump-all-tables操作で--format=binary時)")
	}
	if c.Operation == "dump-all-tables" && c.Format == "binary" && c.Limit != nil {
		return fmt.Errorf("dump-all-tables操作で--format=binaryの場合、--limit は指定できません")
	}

	// limitの検証
	if c.Limit != nil && *c.Limit < 0 {
		return fmt.Errorf("--limit は正の数もしくは0である必要があります: %d", *c.Limit)
	}

	if c.Timezone != "" {
		if _, err := time.LoadLocation(c.Timezone); err != nil {
			return fmt.Errorf("--timezone で指定されたタイムゾーンが無効です: %v", err)
		}
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
	fmt.Fprintf(os.Stderr, "  --format string         出力フォーマット: json, csv, sql, text, binary (デフォルト: json)\n")
	fmt.Fprintf(os.Stderr, "  --result-format string  ダンプ結果フォーマット: json, markdown (デフォルト: json)\n")
	fmt.Fprintf(os.Stderr, "  --result-heading string ダンプ結果サマリの見出し (Markdown時のみ)\n")
	fmt.Fprintf(os.Stderr, "  --timezone string       タイムスタンプ/ファイル名で利用するタイムゾーン (例: Asia/Tokyo)\n")
	fmt.Fprintf(os.Stderr, "  --limit int             最大レコード数\n")
	fmt.Fprintf(os.Stderr, "  --concurrency int       並行処理数: 1-CPUコア数 (デフォルト: CPUコア数)\n")
	fmt.Fprintf(os.Stderr, "  --help                  このヘルプを表示\n\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  # 単一テーブルダンプ\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump --database-url=\"postgres://user:pass@localhost/db\" --table-name=users\n\n")
	fmt.Fprintf(os.Stderr, "  # 全テーブルダンプ\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump-all-tables --database-url=\"postgres://user:pass@localhost/db\"\n\n")
	fmt.Fprintf(os.Stderr, "  # 全テーブルをbinary形式でダンプ\n")
	fmt.Fprintf(os.Stderr, "  postgresql-cli --operation=dump-all-tables --database-url=\"postgres://user:pass@localhost/db\" --format=binary\n\n")
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
