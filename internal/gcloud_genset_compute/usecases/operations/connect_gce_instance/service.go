package connectgceinstance

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

const defaultSSHKeyPath = "$HOME/.ssh/google_compute_engine"

// Params は GCE SSH 接続コマンド生成に必要な値。
type Params struct {
	InstanceName  string
	Zone          string
	SSHKeyPath    string
	CreatesSSHKey bool
	Forces        bool
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
	if common.IsBlank(params.SSHKeyPath) {
		params.SSHKeyPath = defaultSSHKeyPath
	}
	if params.Forces && !params.CreatesSSHKey {
		return "", fmt.Errorf("forces は creates-ssh-key=true の場合のみ指定できます")
	}

	connectCommand := fmt.Sprintf(
		"gcloud compute ssh %s --zone=%s --tunnel-through-iap",
		common.ShellQuote(params.InstanceName),
		common.ShellQuote(params.Zone),
	)

	commands := make([]string, 0, 3)
	if params.CreatesSSHKey {
		commands = append(commands, common.BuildSSHKeyCreationCommand(params.SSHKeyPath, params.Forces))
	}
	commands = append(commands,
		common.BuildSSHAgentSetupCommand(params.SSHKeyPath),
		connectCommand,
	)
	return strings.Join(commands, " && \\\n"), nil
}
