package usecases

import (
	"fmt"
	"strings"
)

// Service は Document AI 関連の gcloud コマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// UndeployProcessorVersionParams はプロセッサバージョンのアンデプロイに必要なパラメータを表す。
type UndeployProcessorVersionParams struct {
	Region        string
	ProjectNumber string
	ProcessorID   string
	VersionID     string
}

// BuildUndeployProcessorVersionCommand は Document AI プロセッサバージョンをアンデプロイするためのコマンドを生成する。
func (s *Service) BuildUndeployProcessorVersionCommand(params UndeployProcessorVersionParams) (string, error) {
	region := strings.TrimSpace(params.Region)
	if region == "" {
		return "", fmt.Errorf("region は必須です")
	}

	projectNumber := strings.TrimSpace(params.ProjectNumber)
	if projectNumber == "" {
		return "", fmt.Errorf("project-number は必須です")
	}

	processorID := strings.TrimSpace(params.ProcessorID)
	if processorID == "" {
		return "", fmt.Errorf("processor-id は必須です")
	}

	versionID := strings.TrimSpace(params.VersionID)
	if versionID == "" {
		return "", fmt.Errorf("version-id は必須です")
	}

	endpoint := fmt.Sprintf("https://%s-documentai.googleapis.com/v1beta3/%s/locations/%s/processors/%s/processorVersions/%s:undeploy", region, projectNumber, region, processorID, versionID)

	command := fmt.Sprintf("curl -s -X POST -H \"Authorization: Bearer $(gcloud auth print-access-token)\" -H \"Content-Type: application/json\" \"%s\"", endpoint)

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
