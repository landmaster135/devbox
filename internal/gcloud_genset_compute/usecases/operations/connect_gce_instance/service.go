package connectgceinstance

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

const defaultSSHKeyPath = "$HOME/.ssh/google_compute_engine"

// Params は GCE SSH 接続コマンド生成に必要な値。
type Params struct {
	InstanceName string
	Zone         string
}

// Service は connect-gce-instance operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build は GCE SSH 接続コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.InstanceName) {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if common.IsBlank(params.Zone) {
		return "", fmt.Errorf("zone は必須です")
	}

	connectCommand := fmt.Sprintf(
		"gcloud compute ssh %s --zone=%s --tunnel-through-iap",
		common.ShellQuote(params.InstanceName),
		common.ShellQuote(params.Zone),
	)

	commands := []string{
		common.BuildSSHAgentSetupCommand(defaultSSHKeyPath),
		connectCommand,
	}
	return strings.Join(commands, " && \\\n"), nil
}
