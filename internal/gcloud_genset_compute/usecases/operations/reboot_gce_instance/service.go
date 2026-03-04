package rebootgceinstance

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params はインスタンス再起動コマンド生成に必要な値。
type Params struct {
	InstanceName string
	Zone         string
}

// Service は reboot-gce-instance operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build はインスタンス再起動コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.InstanceName) {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if common.IsBlank(params.Zone) {
		return "", fmt.Errorf("zone は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute instances reset %s --zone=%s",
		common.ShellQuote(params.InstanceName),
		common.ShellQuote(params.Zone),
	), nil
}
