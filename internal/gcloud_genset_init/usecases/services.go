package usecases

import (
	"fmt"
	"strings"
)

// Service は gcloud 初期設定向けのコマンドを生成する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// AuthLoginParams は gcloud auth login コマンド生成時のパラメータを表す。
type AuthLoginParams struct {
	ProjectID      string
	AdditionalArgs string
}

// SetProjectConfigParams は gcloud config set project コマンド生成時のパラメータを表す。
type SetProjectConfigParams struct {
	ProjectID      string
	AdditionalArgs string
}

// BuildAuthLoginCommand は gcloud auth login コマンドを生成する。
func (s *Service) BuildAuthLoginCommand(params AuthLoginParams) (string, error) {
	projectID := strings.TrimSpace(params.ProjectID)
	if projectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}

	command := fmt.Sprintf("gcloud auth login %s", shellQuote(projectID))

	if trimmed := strings.TrimSpace(params.AdditionalArgs); trimmed != "" {
		command = fmt.Sprintf("%s %s", command, trimmed)
	}

	return command, nil
}

// BuildSetProjectConfigCommand は gcloud config set project コマンドを生成する。
func (s *Service) BuildSetProjectConfigCommand(params SetProjectConfigParams) (string, error) {
	projectID := strings.TrimSpace(params.ProjectID)
	if projectID == "" {
		return "", fmt.Errorf("project-id は必須です")
	}

	command := fmt.Sprintf("gcloud config set project %s", shellQuote(projectID))

	if trimmed := strings.TrimSpace(params.AdditionalArgs); trimmed != "" {
		command = fmt.Sprintf("%s %s", command, trimmed)
	}

	return command, nil
}

// PrintHighlightedCommand は生成したコマンドを見やすい形式で出力する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}
