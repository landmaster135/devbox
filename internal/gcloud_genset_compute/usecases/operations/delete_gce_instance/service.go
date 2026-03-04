package deletegceinstance

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params はインスタンス削除コマンド生成に必要な値。
type Params struct {
	InstanceName string
	Zone         string
}

// Service は delete-gce-instance operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build はインスタンス削除コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.InstanceName) {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if common.IsBlank(params.Zone) {
		return "", fmt.Errorf("zone は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute instances delete %s --zone=%s --quiet",
		common.ShellQuote(params.InstanceName),
		common.ShellQuote(params.Zone),
	), nil
}
