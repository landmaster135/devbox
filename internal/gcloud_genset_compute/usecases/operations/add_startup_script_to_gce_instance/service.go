package addstartupscripttogceinstance

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params はスタートアップスクリプト登録コマンド生成に必要な値。
type Params struct {
	InstanceName      string
	Zone              string
	StartupScriptPath string
}

// Service は add-startup-script-to-gce-instance operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build はスタートアップスクリプト登録コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.InstanceName) {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if common.IsBlank(params.Zone) {
		return "", fmt.Errorf("zone は必須です")
	}
	if common.IsBlank(params.StartupScriptPath) {
		return "", fmt.Errorf("startup-script-path は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute instances add-metadata %s --zone=%s --metadata-from-file startup-script=%s",
		common.ShellQuote(params.InstanceName),
		common.ShellQuote(params.Zone),
		common.ShellQuote(params.StartupScriptPath),
	), nil
}
