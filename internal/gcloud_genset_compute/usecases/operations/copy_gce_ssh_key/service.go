package copygcesshkey

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params は SSH 鍵コピーコマンド生成に必要な値。
type Params struct {
	InstanceName string
	Zone         string
	SSHKeyPath   string
}

// Service は copy-gce-ssh-key operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build は SSH 鍵コピーコマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.InstanceName) {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if common.IsBlank(params.Zone) {
		return "", fmt.Errorf("zone は必須です")
	}
	if common.IsBlank(params.SSHKeyPath) {
		return "", fmt.Errorf("ssh-key-path は必須です")
	}

	copyCommand := fmt.Sprintf(
		"gcloud compute scp %s %s --zone=%s --tunnel-through-iap",
		common.ShellQuoteSSHKeyPath(params.SSHKeyPath),
		common.ShellQuote(params.InstanceName+":/tmp"),
		common.ShellQuote(params.Zone),
	)

	commands := []string{
		common.BuildSSHAgentSetupCommand(params.SSHKeyPath),
		copyCommand,
	}
	return strings.Join(commands, " && \\\n"), nil
}
