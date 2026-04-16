package setgceinstancemetadatafromyaml

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params は YAML から metadata 設定コマンド生成に必要な値。
type Params struct {
	InstanceName     string
	Zone             string
	MetadataYAMLPath string
}

// Service は set-gce-instance-metadata-from-yaml operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build は YAML から metadata 設定コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.InstanceName) {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if common.IsBlank(params.Zone) {
		return "", fmt.Errorf("zone は必須です")
	}
	if common.IsBlank(params.MetadataYAMLPath) {
		return "", fmt.Errorf("metadata-yaml-path は必須です")
	}

	metadataValue, err := s.buildMetadataValue(params.MetadataYAMLPath)
	if err != nil {
		return "", err
	}
	if metadataValue == "" {
		return fmt.Sprintf(
			"echo %s",
			common.ShellQuote(fmt.Sprintf("[INFO] No valid metadata found in '%s'", params.MetadataYAMLPath)),
		), nil
	}

	return fmt.Sprintf(
		"gcloud compute instances add-metadata %s --zone=%s --metadata=%s",
		common.ShellQuote(params.InstanceName),
		common.ShellQuote(params.Zone),
		common.ShellQuote(metadataValue),
	), nil
}

func (s *Service) buildMetadataValue(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("metadata-yaml-path の読み取りに失敗しました: %w", err)
	}

	pairs := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := normalizeMetadataValue(parts[1])
		if key == "" {
			continue
		}
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("metadata-yaml-path の解析に失敗しました: %w", err)
	}
	return strings.Join(pairs, ","), nil
}

func normalizeMetadataValue(raw string) string {
	value := strings.TrimSpace(raw)
	if len(value) >= 2 {
		if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
			return value[1 : len(value)-1]
		}
	}
	return value
}
