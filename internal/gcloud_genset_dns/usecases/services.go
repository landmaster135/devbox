package usecases

import (
	"fmt"
	"strings"
)

// Service は Google Cloud DNS 向けの gcloud コマンドを生成するサービス。
type Service struct{}

// NewService は Service のインスタンスを生成する。
func NewService() *Service {
	return &Service{}
}

// ManagedZonesListParams は gcloud dns managed-zones list コマンド生成に必要なパラメータ。
type ManagedZonesListParams struct {
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

// BuildManagedZonesListCommand は gcloud dns managed-zones list コマンドを生成する。
func (s *Service) BuildManagedZonesListCommand(params ManagedZonesListParams) (string, error) {
	var commandParts []string
	commandParts = append(commandParts, "gcloud dns managed-zones list")

	if project := strings.TrimSpace(params.Project); project != "" {
		commandParts = append(commandParts, fmt.Sprintf("--project=%s", shellQuote(project)))
	}

	if format := strings.TrimSpace(params.Format); format != "" {
		commandParts = append(commandParts, fmt.Sprintf("--format=%s", shellQuote(format)))
	}

	if filter := strings.TrimSpace(params.Filter); filter != "" {
		commandParts = append(commandParts, fmt.Sprintf("--filter=%s", shellQuote(filter)))
	}

	if params.Limit < 0 {
		return "", fmt.Errorf("limit は0以上で指定してください")
	}
	if params.Limit > 0 {
		commandParts = append(commandParts, fmt.Sprintf("--limit=%d", params.Limit))
	}

	if params.PageSize < 0 {
		return "", fmt.Errorf("page-size は0以上で指定してください")
	}
	if params.PageSize > 0 {
		commandParts = append(commandParts, fmt.Sprintf("--page-size=%d", params.PageSize))
	}

	if sortBy := strings.TrimSpace(params.SortBy); sortBy != "" {
		commandParts = append(commandParts, fmt.Sprintf("--sort-by=%s", shellQuote(sortBy)))
	}

	if verbosity := strings.TrimSpace(params.Verbosity); verbosity != "" {
		commandParts = append(commandParts, fmt.Sprintf("--verbosity=%s", shellQuote(verbosity)))
	}

	if params.URI {
		commandParts = append(commandParts, "--uri")
	}

	if additional := strings.TrimSpace(params.AdditionalArgs); additional != "" {
		commandParts = append(commandParts, additional)
	}

	return strings.Join(commandParts, " "), nil
}

// PrintHighlightedCommand は生成したコマンドを視認しやすい形式で出力する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\\''")
	return "'" + escaped + "'"
}
