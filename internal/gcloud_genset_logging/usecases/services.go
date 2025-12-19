package usecases

import (
	"fmt"
	"strings"
)

// Service はgcloud loggingコマンドを生成する
type Service struct{}

// NewService はServiceインスタンスを返す
func NewService() *Service {
	return &Service{}
}

// LoggingReadParams はlogging readコマンド生成時のパラメータを表す
type LoggingReadParams struct {
	Severity       string
	Limit          int
	Query          string
	ResourceType   string
	Filter         string
	AdditionalArgs string
}

// CreateSinkParams はlogging sinks createコマンド生成時のパラメータを表す
type CreateSinkParams struct {
	SinkName       string
	Destination    string
	LogFilter      string
	AdditionalArgs string
}

// BuildLoggingReadCommand はgcloud logging readコマンドを生成する
func (s *Service) BuildLoggingReadCommand(params LoggingReadParams) (string, error) {
	if params.Limit <= 0 {
		return "", fmt.Errorf("limit は1以上で指定してください")
	}

	filter, err := s.buildLoggingReadFilter(params)
	if err != nil {
		return "", err
	}

	command := fmt.Sprintf("gcloud logging read \"%s\" --limit=%d", filter, params.Limit)

	if trimmed := strings.TrimSpace(params.AdditionalArgs); trimmed != "" {
		command = fmt.Sprintf("%s %s", command, trimmed)
	}

	return command, nil
}

func (s *Service) buildLoggingReadFilter(params LoggingReadParams) (string, error) {
	if params.Filter != "" {
		return params.Filter, nil
	}

	var parts []string
	if params.Severity != "" {
		parts = append(parts, fmt.Sprintf("severity>=%s", params.Severity))
	}
	if params.ResourceType != "" {
		parts = append(parts, fmt.Sprintf("resource.type=%s", params.ResourceType))
	}
	if params.Query != "" {
		parts = append(parts, params.Query)
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("logging read用のフィルター条件を指定してください")
	}

	return strings.Join(parts, " AND "), nil
}

// BuildCreateSinkCommand はgcloud logging sinks createコマンドを生成する
func (s *Service) BuildCreateSinkCommand(params CreateSinkParams) (string, error) {
	if strings.TrimSpace(params.SinkName) == "" {
		return "", fmt.Errorf("sink-name は必須です")
	}
	if strings.TrimSpace(params.Destination) == "" {
		return "", fmt.Errorf("destination は必須です")
	}

	command := fmt.Sprintf("gcloud logging sinks create %s %s", params.SinkName, params.Destination)

	if params.LogFilter != "" {
		command = fmt.Sprintf("%s --log-filter=\"%s\"", command, params.LogFilter)
	}

	if trimmed := strings.TrimSpace(params.AdditionalArgs); trimmed != "" {
		command = fmt.Sprintf("%s %s", command, trimmed)
	}

	return command, nil
}

func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}
