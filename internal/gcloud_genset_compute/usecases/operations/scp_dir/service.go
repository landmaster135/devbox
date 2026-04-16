package scpdir

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params は scp-dir operation のコマンド生成に必要な値。
type Params struct {
	InstanceName string
	Zone         string
	SrcDir       string
	DestDir      string
}

// Service は scp-dir operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build はローカルディレクトリを再帰コピーするコマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.InstanceName) {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if common.IsBlank(params.Zone) {
		return "", fmt.Errorf("zone は必須です")
	}
	if common.IsBlank(params.SrcDir) {
		return "", fmt.Errorf("src-dir は必須です")
	}
	if common.IsBlank(params.DestDir) {
		return "", fmt.Errorf("dest-dir は必須です")
	}

	normalizedDestDir := strings.TrimSpace(params.DestDir)
	normalizedDestDir = strings.TrimRight(normalizedDestDir, "/")
	if normalizedDestDir == "" {
		return "", fmt.Errorf("dest-dir は必須です")
	}

	destination := fmt.Sprintf("%s:%s/", params.InstanceName, normalizedDestDir)
	return fmt.Sprintf(
		"gcloud compute scp --recurse %s %s --zone=%s",
		common.ShellQuote(params.SrcDir),
		common.ShellQuote(destination),
		common.ShellQuote(params.Zone),
	), nil
}
