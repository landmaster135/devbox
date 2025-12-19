package config

import (
	"fmt"
	"os"
)

const (
	// OperationManagedZonesList は managed-zones list コマンド生成を表す。
	OperationManagedZonesList = "managed-zones-list"
)

// Config は CLI で指定されたパラメータを保持する。
type Config struct {
	Operation      string
	Help           bool
	Project        string
	Format         string
	Filter         string
	Limit          int
	PageSize       int
	SortBy         string
	Verbosity      string
	URI            bool
	AdditionalArgs string
}

// ParseFlags は標準のフラグパーサーで引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は注入されたフラグパーサーで引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: managed-zones-list)")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")
	parser.StringVar(&cfg.Project, "project", "", "対象のGCPプロジェクトID")
	parser.StringVar(&cfg.Format, "format", "", "出力フォーマット (yaml, json, csv 等)")
	parser.StringVar(&cfg.Filter, "filter", "", "結果をフィルタリングする条件")
	parser.IntVar(&cfg.Limit, "limit", 0, "表示する結果の最大数")
	parser.IntVar(&cfg.PageSize, "page-size", 0, "1ページあたりの結果数")
	parser.StringVar(&cfg.SortBy, "sort-by", "", "結果のソート基準")
	parser.StringVar(&cfg.Verbosity, "verbosity", "", "詳細レベル (debug, info, warning, error, critical, none)")
	parser.BoolVar(&cfg.URI, "uri", false, "URI 形式で出力")
	parser.StringVar(&cfg.AdditionalArgs, "additional-args", "", "gcloud コマンドに追加する引数")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if len(parser.Args()) > 0 {
		return nil, fmt.Errorf("未処理の位置引数があります: %v", parser.Args())
	}

	if cfg.Help {
		return cfg, nil
	}

	if cfg.Operation == "" {
		return nil, fmt.Errorf("operation パラメータは必須です")
	}

	if cfg.Limit < 0 {
		return nil, fmt.Errorf("limit パラメータは0以上で指定してください")
	}

	if cfg.PageSize < 0 {
		return nil, fmt.Errorf("page-size パラメータは0以上で指定してください")
	}

	switch cfg.Operation {
	case OperationManagedZonesList:
		// 追加の必須パラメータなし
	default:
		return nil, fmt.Errorf("未対応の操作です: %s", cfg.Operation)
	}

	return cfg, nil
}

// PrintUsage は CLI ツールの使用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud DNS 向け gcloud コマンド生成ツール\n\n")
	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (managed-zones-list)\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "managed-zones-list 操作用パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -project string\n")
	fmt.Fprintf(os.Stderr, "        対象のGCPプロジェクトID\n")
	fmt.Fprintf(os.Stderr, "  -format string\n")
	fmt.Fprintf(os.Stderr, "        出力フォーマット (yaml, json, csv 等)\n")
	fmt.Fprintf(os.Stderr, "  -filter string\n")
	fmt.Fprintf(os.Stderr, "        結果をフィルタリングする条件\n")
	fmt.Fprintf(os.Stderr, "  -limit int\n")
	fmt.Fprintf(os.Stderr, "        表示する結果の最大数 (0 で指定なし)\n")
	fmt.Fprintf(os.Stderr, "  -page-size int\n")
	fmt.Fprintf(os.Stderr, "        1ページあたりの結果数 (0 で指定なし)\n")
	fmt.Fprintf(os.Stderr, "  -sort-by string\n")
	fmt.Fprintf(os.Stderr, "        結果のソート基準\n")
	fmt.Fprintf(os.Stderr, "  -verbosity string\n")
	fmt.Fprintf(os.Stderr, "        詳細レベル (debug, info, warning, error, critical, none)\n")
	fmt.Fprintf(os.Stderr, "  -uri\n")
	fmt.Fprintf(os.Stderr, "        URI 形式で出力\n")
	fmt.Fprintf(os.Stderr, "  -additional-args string\n")
	fmt.Fprintf(os.Stderr, "        gcloud コマンドに追加する任意の引数\n\n")

	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=managed-zones-list\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=managed-zones-list -project=my-project -format=json\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=managed-zones-list -filter=\"name:example\" -limit=20 --uri\n", os.Args[0])
}
