package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	OperationListDashboards    = "list-dashboards"
	OperationDescribeDashboard = "describe-dashboard"
	OperationListSnoozes       = "list-snoozes"
	OperationListUptimeConfigs = "list-uptime-configs"
)

// Config はCLIで指定されたパラメータを保持する
// operation は必須であり、その他のパラメータは操作内容に応じて利用される
type Config struct {
	Operation   string
	Help        bool
	Project     string
	Filter      string
	Format      string
	PageSize    int
	SortBy      string
	Limit       int
	URI         bool
	DashboardID string
}

// ParseFlags は標準のフラグパーサーで引数を解析する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたパーサーを使って引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{
		PageSize: -1,
		Limit:    -1,
	}

	parser.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: list-dashboards, describe-dashboard, list-snoozes, list-uptime-configs)")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")
	parser.StringVar(&cfg.Project, "project", "", "Google Cloud プロジェクトID")
	parser.StringVar(&cfg.Filter, "filter", "", "結果をフィルタリングする式")
	parser.StringVar(&cfg.Format, "format", "", "出力形式 (table, json, yaml など)")
	parser.IntVar(&cfg.PageSize, "page-size", -1, "1ページあたりの取得件数")
	parser.StringVar(&cfg.SortBy, "sort-by", "", "ソート対象のフィールド")
	parser.IntVar(&cfg.Limit, "limit", -1, "取得する最大件数")
	parser.BoolVar(&cfg.URI, "uri", false, "リソースの URI を表示 (list-snoozes/list-uptime-configs 操作用)")
	parser.StringVar(&cfg.DashboardID, "dashboard-id", "", "対象のダッシュボードID (describe-dashboard 操作用)")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	remaining := parser.Args()
	if len(remaining) > 0 {
		return nil, fmt.Errorf("未処理の引数があります: %v", remaining)
	}

	if cfg.Help {
		return cfg, nil
	}

	if cfg.Operation == "" {
		return nil, fmt.Errorf("operation パラメータは必須です")
	}

	if cfg.PageSize < -1 {
		return nil, fmt.Errorf("page-size パラメータは1以上を指定するか省略してください")
	}
	if cfg.Limit < -1 {
		return nil, fmt.Errorf("limit パラメータは1以上を指定するか省略してください")
	}

	switch cfg.Operation {
	case OperationListDashboards:
		if cfg.PageSize == 0 {
			return nil, fmt.Errorf("page-size パラメータは1以上で指定してください")
		}
		if cfg.Limit == 0 {
			return nil, fmt.Errorf("limit パラメータは1以上で指定してください")
		}
		if cfg.URI {
			return nil, fmt.Errorf("list-dashboards 操作では uri パラメータを指定できません")
		}
	case OperationDescribeDashboard:
		if strings.TrimSpace(cfg.DashboardID) == "" {
			return nil, fmt.Errorf("describe-dashboard 操作には dashboard-id パラメータが必須です")
		}
		if cfg.Filter != "" {
			return nil, fmt.Errorf("describe-dashboard 操作では filter パラメータを指定できません")
		}
		if cfg.PageSize > -1 {
			return nil, fmt.Errorf("describe-dashboard 操作では page-size パラメータを指定できません")
		}
		if cfg.SortBy != "" {
			return nil, fmt.Errorf("describe-dashboard 操作では sort-by パラメータを指定できません")
		}
		if cfg.Limit > -1 {
			return nil, fmt.Errorf("describe-dashboard 操作では limit パラメータを指定できません")
		}
		if cfg.URI {
			return nil, fmt.Errorf("describe-dashboard 操作では uri パラメータを指定できません")
		}
	case OperationListSnoozes, OperationListUptimeConfigs:
		if cfg.PageSize == 0 {
			return nil, fmt.Errorf("page-size パラメータは1以上で指定してください")
		}
		if cfg.Limit == 0 {
			return nil, fmt.Errorf("limit パラメータは1以上で指定してください")
		}
	default:
		return nil, fmt.Errorf("未対応の操作です: %s", cfg.Operation)
	}

	return cfg, nil
}

// PrintUsage はツールの使用方法を標準エラーへ表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud Monitoring 用の gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (list-dashboards | describe-dashboard | list-snoozes | list-uptime-configs)\n")
	fmt.Fprintf(os.Stderr, "  -project string\n")
	fmt.Fprintf(os.Stderr, "        Google Cloud プロジェクトID\n")
	fmt.Fprintf(os.Stderr, "  -format string\n")
	fmt.Fprintf(os.Stderr, "        出力形式 (table, json, yaml など)\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "list-dashboards 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -filter string\n")
	fmt.Fprintf(os.Stderr, "        ダッシュボード一覧のフィルター式\n")
	fmt.Fprintf(os.Stderr, "  -page-size int\n")
	fmt.Fprintf(os.Stderr, "        1ページあたりの取得件数\n")
	fmt.Fprintf(os.Stderr, "  -sort-by string\n")
	fmt.Fprintf(os.Stderr, "        並び替えに使用するフィールド\n")
	fmt.Fprintf(os.Stderr, "  -limit int\n")
	fmt.Fprintf(os.Stderr, "        取得する最大件数\n\n")

	fmt.Fprintf(os.Stderr, "describe-dashboard 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -dashboard-id string\n")
	fmt.Fprintf(os.Stderr, "        ダッシュボードID (必須)\n\n")

	fmt.Fprintf(os.Stderr, "list-snoozes / list-uptime-configs 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -filter string\n")
	fmt.Fprintf(os.Stderr, "        フィルター式\n")
	fmt.Fprintf(os.Stderr, "  -page-size int\n")
	fmt.Fprintf(os.Stderr, "        1ページあたりの取得件数\n")
	fmt.Fprintf(os.Stderr, "  -sort-by string\n")
	fmt.Fprintf(os.Stderr, "        並び替えに使用するフィールド\n")
	fmt.Fprintf(os.Stderr, "  -limit int\n")
	fmt.Fprintf(os.Stderr, "        取得する最大件数\n")
	fmt.Fprintf(os.Stderr, "  -uri\n")
	fmt.Fprintf(os.Stderr, "        リソース URI を表示\n\n")

	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=list-dashboards -project=my-project -format=json\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=describe-dashboard -project=my-project -dashboard-id=my-dashboard\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=list-snoozes -project=my-project -filter='displayName:maintenance' -uri\n", os.Args[0])
}
