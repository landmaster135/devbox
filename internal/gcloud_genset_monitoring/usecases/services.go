package usecases

import (
	"fmt"
	"strconv"
	"strings"
)

// Service は gcloud monitoring コマンドを生成する責務を持つ
type Service struct{}

// NewService は Service のインスタンスを返す
func NewService() *Service {
	return &Service{}
}

// ListDashboardsParams は dashboards list コマンド生成時の入力値を表す
type ListDashboardsParams struct {
	Project  string
	Filter   string
	Format   string
	PageSize int
	SortBy   string
	Limit    int
}

// DescribeDashboardParams は dashboards describe コマンド生成時の入力値を表す
type DescribeDashboardParams struct {
	DashboardID string
	Project     string
	Format      string
}

// ListSnoozesParams は snoozes list コマンド生成時の入力値を表す
type ListSnoozesParams struct {
	Project    string
	Filter     string
	Format     string
	PageSize   int
	SortBy     string
	Limit      int
	IncludeURI bool
}

// ListUptimeConfigsParams は uptime list-configs コマンド生成時の入力値を表す
type ListUptimeConfigsParams struct {
	Project    string
	Filter     string
	Format     string
	PageSize   int
	SortBy     string
	Limit      int
	IncludeURI bool
}

// BuildListDashboardsCommand は gcloud monitoring dashboards list コマンドを生成する
func (s *Service) BuildListDashboardsCommand(params ListDashboardsParams) (string, error) {
	if params.PageSize == 0 {
		return "", fmt.Errorf("page-size は1以上で指定してください")
	}
	if params.Limit == 0 {
		return "", fmt.Errorf("limit は1以上で指定してください")
	}

	options := buildCommonListOptions(listCommonOptions{
		Project:  params.Project,
		Filter:   params.Filter,
		Format:   params.Format,
		PageSize: params.PageSize,
		SortBy:   params.SortBy,
		Limit:    params.Limit,
	})

	return buildCommand("gcloud monitoring dashboards list", options), nil
}

// BuildDescribeDashboardCommand は gcloud monitoring dashboards describe コマンドを生成する
func (s *Service) BuildDescribeDashboardCommand(params DescribeDashboardParams) (string, error) {
	dashboardID := strings.TrimSpace(params.DashboardID)
	if dashboardID == "" {
		return "", fmt.Errorf("dashboard-id は必須です")
	}

	var options []string
	if project := strings.TrimSpace(params.Project); project != "" {
		options = append(options, fmt.Sprintf("--project=%s", project))
	}
	if format := strings.TrimSpace(params.Format); format != "" {
		options = append(options, fmt.Sprintf("--format=%s", strconv.Quote(format)))
	}

	command := fmt.Sprintf("gcloud monitoring dashboards describe %s", dashboardID)
	return buildCommand(command, options), nil
}

// BuildListSnoozesCommand は gcloud monitoring snoozes list コマンドを生成する
func (s *Service) BuildListSnoozesCommand(params ListSnoozesParams) (string, error) {
	if params.PageSize == 0 {
		return "", fmt.Errorf("page-size は1以上で指定してください")
	}
	if params.Limit == 0 {
		return "", fmt.Errorf("limit は1以上で指定してください")
	}

	options := buildCommonListOptions(listCommonOptions{
		Project:  params.Project,
		Filter:   params.Filter,
		Format:   params.Format,
		PageSize: params.PageSize,
		SortBy:   params.SortBy,
		Limit:    params.Limit,
	})

	if params.IncludeURI {
		options = append(options, "--uri")
	}

	return buildCommand("gcloud monitoring snoozes list", options), nil
}

// BuildListUptimeConfigsCommand は gcloud monitoring uptime list-configs コマンドを生成する
func (s *Service) BuildListUptimeConfigsCommand(params ListUptimeConfigsParams) (string, error) {
	if params.PageSize == 0 {
		return "", fmt.Errorf("page-size は1以上で指定してください")
	}
	if params.Limit == 0 {
		return "", fmt.Errorf("limit は1以上で指定してください")
	}

	options := buildCommonListOptions(listCommonOptions{
		Project:  params.Project,
		Filter:   params.Filter,
		Format:   params.Format,
		PageSize: params.PageSize,
		SortBy:   params.SortBy,
		Limit:    params.Limit,
	})

	if params.IncludeURI {
		options = append(options, "--uri")
	}

	return buildCommand("gcloud monitoring uptime list-configs", options), nil
}

// PrintHighlightedCommand は生成したコマンドを見やすい形で出力する
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

type listCommonOptions struct {
	Project  string
	Filter   string
	Format   string
	PageSize int
	SortBy   string
	Limit    int
}

func buildCommonListOptions(params listCommonOptions) []string {
	var options []string

	if project := strings.TrimSpace(params.Project); project != "" {
		options = append(options, fmt.Sprintf("--project=%s", project))
	}
	if filter := strings.TrimSpace(params.Filter); filter != "" {
		options = append(options, fmt.Sprintf("--filter=%s", filter))
	}
	if format := strings.TrimSpace(params.Format); format != "" {
		options = append(options, fmt.Sprintf("--format=%s", strconv.Quote(format)))
	}
	if params.PageSize > 0 {
		options = append(options, fmt.Sprintf("--page-size=%d", params.PageSize))
	}
	if sortBy := strings.TrimSpace(params.SortBy); sortBy != "" {
		options = append(options, fmt.Sprintf("--sort-by=%s", sortBy))
	}
	if params.Limit > 0 {
		options = append(options, fmt.Sprintf("--limit=%d", params.Limit))
	}

	return options
}

func buildCommand(base string, options []string) string {
	if len(options) == 0 {
		return base
	}
	return fmt.Sprintf("%s %s", base, strings.Join(options, " "))
}
