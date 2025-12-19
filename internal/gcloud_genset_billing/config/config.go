package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	// OperationListBudgets は請求予算一覧取得用コマンドを生成する操作を表す。
	OperationListBudgets = "list-budgets"
	// OperationListProjects は請求対象プロジェクト一覧取得用コマンドを生成する操作を表す。
	OperationListProjects = "list-projects"
	// OperationDescribeProject は請求対象プロジェクト詳細取得用コマンドを生成する操作を表す。
	OperationDescribeProject = "describe-project"
	// OperationDescribeBudget は請求予算詳細取得用コマンドを生成する操作を表す。
	OperationDescribeBudget = "describe-budget"

	defaultLimit = 10
)

var validOperations = []string{
	OperationListBudgets,
	OperationListProjects,
	OperationDescribeProject,
	OperationDescribeBudget,
}

// Config は CLI で指定されたパラメータを保持する。
type Config struct {
	Operation      string
	Help           bool
	Limit          int
	Filter         string
	BillingAccount string
	ProjectID      string
	BudgetID       string
}

// ParseFlags は標準のフラグパーサーで引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定のパーサーで引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{Limit: defaultLimit}

	parser.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(validOperations, ", ")))
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")
	parser.IntVar(&cfg.Limit, "limit", defaultLimit, "取得する件数の上限")
	parser.StringVar(&cfg.Filter, "filter", "", "gcloud billing projects list に渡すフィルター")
	parser.StringVar(&cfg.BillingAccount, "billing-account", "", "対象の請求アカウントID")
	parser.StringVar(&cfg.ProjectID, "project-id", "", "対象のプロジェクトID")
	parser.StringVar(&cfg.BudgetID, "budget-id", "", "対象の予算ID")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	cfg.Operation = strings.TrimSpace(cfg.Operation)
	cfg.BillingAccount = strings.TrimSpace(cfg.BillingAccount)
	cfg.ProjectID = strings.TrimSpace(cfg.ProjectID)
	cfg.BudgetID = strings.TrimSpace(cfg.BudgetID)
	cfg.Filter = strings.TrimSpace(cfg.Filter)

	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}

	if !isValidOperation(cfg.Operation) {
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	switch cfg.Operation {
	case OperationListBudgets:
		if cfg.Limit <= 0 {
			return fmt.Errorf("limit は1以上で指定してください")
		}
	case OperationListProjects:
		if cfg.Limit <= 0 {
			return fmt.Errorf("limit は1以上で指定してください")
		}
	case OperationDescribeBudget:
		if cfg.BudgetID == "" {
			return fmt.Errorf("budget-id パラメータは必須です")
		}
	}

	return nil
}

func isValidOperation(op string) bool {
	for _, v := range validOperations {
		if op == v {
			return true
		}
	}
	return false
}

// PrintUsage は CLI の使い方を標準エラーに出力する。
func PrintUsage() {
	sort.Strings(validOperations)
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud Billing 向け gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (%s)\n", strings.Join(validOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n")
	fmt.Fprintf(os.Stderr, "  -billing-account string\n")
	fmt.Fprintf(os.Stderr, "        対象の請求アカウントID (未指定時は自動取得)\n\n")

	fmt.Fprintf(os.Stderr, "list-budgets / list-projects 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -limit int\n")
	fmt.Fprintf(os.Stderr, "        取得する件数の上限 (デフォルト: %d)\n", defaultLimit)
	fmt.Fprintf(os.Stderr, "  -filter string\n")
	fmt.Fprintf(os.Stderr, "        list-projects 時のフィルター条件\n\n")

	fmt.Fprintf(os.Stderr, "describe-project 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -project-id string\n")
	fmt.Fprintf(os.Stderr, "        対象のプロジェクトID (未指定時は gcloud config を使用)\n\n")

	fmt.Fprintf(os.Stderr, "describe-budget 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -budget-id string\n")
	fmt.Fprintf(os.Stderr, "        対象の予算ID (必須)\n")
	fmt.Fprintf(os.Stderr, "  -billing-account string\n")
	fmt.Fprintf(os.Stderr, "        対象の請求アカウントID (未指定時は自動取得)\n")
}
