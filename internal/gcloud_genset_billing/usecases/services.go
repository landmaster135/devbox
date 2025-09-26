package usecases

import (
	"fmt"
	"strings"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_billing/config"
)

// Service は gcloud billing 向けコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// ListBudgetsParams は予算リスト取得コマンドの生成に必要なパラメータを表す。
type ListBudgetsParams struct {
	BillingAccount string
	Limit          int
}

// ListProjectsParams は請求対象プロジェクト一覧コマンドの生成に必要なパラメータを表す。
type ListProjectsParams struct {
	BillingAccount string
	Limit          int
	Filter         string
}

// DescribeProjectParams は請求対象プロジェクト詳細コマンドの生成に必要なパラメータを表す。
type DescribeProjectParams struct {
	ProjectID string
}

// DescribeBudgetParams は予算詳細取得コマンドの生成に必要なパラメータを表す。
type DescribeBudgetParams struct {
	BillingAccount string
	BudgetID       string
}

// BuildCommand は設定された operation に応じた gcloud コマンドを生成する。
func (s *Service) BuildCommand(conf *cfg.Config) (string, error) {
	switch conf.Operation {
	case cfg.OperationListBudgets:
		return s.BuildListBudgetsCommand(ListBudgetsParams{
			BillingAccount: conf.BillingAccount,
			Limit:          conf.Limit,
		})
	case cfg.OperationListProjects:
		return s.BuildListProjectsCommand(ListProjectsParams{
			BillingAccount: conf.BillingAccount,
			Limit:          conf.Limit,
			Filter:         conf.Filter,
		})
	case cfg.OperationDescribeProject:
		return s.BuildDescribeProjectCommand(DescribeProjectParams{ProjectID: conf.ProjectID})
	case cfg.OperationDescribeBudget:
		return s.BuildDescribeBudgetCommand(DescribeBudgetParams{
			BillingAccount: conf.BillingAccount,
			BudgetID:       conf.BudgetID,
		})
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", conf.Operation)
	}
}

// BuildListBudgetsCommand は gcloud billing budgets list コマンドを生成する。
func (s *Service) BuildListBudgetsCommand(params ListBudgetsParams) (string, error) {
	if params.Limit <= 0 {
		return "", fmt.Errorf("limit は1以上で指定してください")
	}

	billingAccount := strings.TrimSpace(params.BillingAccount)
	if billingAccount != "" {
		return fmt.Sprintf("gcloud billing budgets list --billing-account=%s --limit=%d", shellQuote(billingAccount), params.Limit), nil
	}

	mainCommand := fmt.Sprintf("gcloud billing budgets list --billing-account=\"$billing_account\" --limit=%d", params.Limit)
	return buildCommandWithAutoBillingAccount(mainCommand), nil
}

// BuildListProjectsCommand は gcloud billing projects list コマンドを生成する。
func (s *Service) BuildListProjectsCommand(params ListProjectsParams) (string, error) {
	if params.Limit <= 0 {
		return "", fmt.Errorf("limit は1以上で指定してください")
	}

	billingAccount := strings.TrimSpace(params.BillingAccount)
	command := "gcloud billing projects list"

	if billingAccount != "" {
		command = fmt.Sprintf("%s --billing-account=%s", command, shellQuote(billingAccount))
	} else {
		command = fmt.Sprintf("%s --billing-account=\"$billing_account\"", command)
	}

	command = fmt.Sprintf("%s --limit=%d", command, params.Limit)

	if filter := strings.TrimSpace(params.Filter); filter != "" {
		command = fmt.Sprintf("%s --filter=%s", command, shellQuote(filter))
	}

	if billingAccount != "" {
		return command, nil
	}

	return buildCommandWithAutoBillingAccount(command), nil
}

// BuildDescribeProjectCommand は gcloud billing projects describe コマンドを生成する。
func (s *Service) BuildDescribeProjectCommand(params DescribeProjectParams) (string, error) {
	if projectID := strings.TrimSpace(params.ProjectID); projectID != "" {
		return fmt.Sprintf("gcloud billing projects describe %s", shellQuote(projectID)), nil
	}

	return "project_id=$(gcloud config get-value project 2>/dev/null); if [ -z \"$project_id\" ] || [ \"$project_id\" = \"(unset)\" ]; then echo \"現在の gcloud プロジェクトが設定されていません\" >&2; exit 1; fi; gcloud billing projects describe \"$project_id\"", nil
}

// BuildDescribeBudgetCommand は gcloud billing budgets describe コマンドを生成する。
func (s *Service) BuildDescribeBudgetCommand(params DescribeBudgetParams) (string, error) {
	budgetID := strings.TrimSpace(params.BudgetID)
	if budgetID == "" {
		return "", fmt.Errorf("budget-id は必須です")
	}

	billingAccount := strings.TrimSpace(params.BillingAccount)
	if billingAccount != "" {
		return fmt.Sprintf("gcloud billing budgets describe %s --billing-account=%s", shellQuote(budgetID), shellQuote(billingAccount)), nil
	}

	mainCommand := fmt.Sprintf("gcloud billing budgets describe %s --billing-account=\"$billing_account\"", shellQuote(budgetID))
	return buildCommandWithAutoBillingAccount(mainCommand), nil
}

// PrintHighlightedCommand は生成されたコマンドを強調表示して出力する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

func buildCommandWithAutoBillingAccount(mainCommand string) string {
	prelude := "billing_account=$(gcloud billing accounts list --format='value(name)' 2>/dev/null | head -n 1); if [ -z \"$billing_account\" ]; then echo \"請求アカウントが見つかりません\" >&2; exit 1; fi; "
	return prelude + mainCommand
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\\''")
	return "'" + escaped + "'"
}
