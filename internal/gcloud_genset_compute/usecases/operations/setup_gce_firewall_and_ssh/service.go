package setupgcefirewallandssh

import (
	"fmt"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

const (
	dateSuffixFormat    = "060102-150405"
	defaultIAPRuleName  = "allow-ssh-ingress-from-iap"
	defaultIngressRule  = "allow-ingress-ssh"
	defaultDirection    = "INGRESS"
	defaultAction       = "allow"
	defaultRules        = "tcp:22"
	defaultIAPSources   = "35.235.240.0/20"
	defaultAllowRule    = "tcp:22"
	defaultIngressCIDRs = "10.0.0.0/8"
)

// Params は firewall 作成 + SSH 鍵コピー + SSH 接続コマンド生成に必要な値。
type Params struct {
	InstanceName  string
	Zone          string
	SSHKeyPath    string
	CreatesSSHKey bool
	Forces        bool
}

// Service は setup-gce-firewall-and-ssh operation のコマンド生成を担当する。
type Service struct {
	now func() time.Time
}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return newServiceWithNow(time.Now)
}

func newServiceWithNow(now func() time.Time) *Service {
	return &Service{now: now}
}

// Build は firewall 作成 + SSH 鍵コピー + SSH 接続の複合コマンドを生成する。
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
	if params.Forces && !params.CreatesSSHKey {
		return "", fmt.Errorf("forces は creates-ssh-key=true の場合のみ指定できます")
	}

	dateSuffix := s.now().Format(dateSuffixFormat)
	iapRuleName := fmt.Sprintf("%s-%s", defaultIAPRuleName, dateSuffix)
	ingressRuleName := fmt.Sprintf("%s-%s", defaultIngressRule, dateSuffix)

	commands := make([]string, 0, 6)
	if params.CreatesSSHKey {
		commands = append(commands, common.BuildSSHKeyCreationCommand(params.SSHKeyPath, params.Forces))
	}
	commands = append(commands,
		common.BuildSSHAgentSetupCommand(params.SSHKeyPath),
		fmt.Sprintf(
			"gcloud compute firewall-rules create %s --direction=%s --action=%s --rules=%s --source-ranges=%s",
			common.ShellQuote(iapRuleName),
			common.ShellQuote(defaultDirection),
			common.ShellQuote(defaultAction),
			common.ShellQuote(defaultRules),
			common.ShellQuote(defaultIAPSources),
		),
		fmt.Sprintf(
			"gcloud compute firewall-rules create %s --allow=%s --source-ranges=%s",
			common.ShellQuote(ingressRuleName),
			common.ShellQuote(defaultAllowRule),
			common.ShellQuote(defaultIngressCIDRs),
		),
		fmt.Sprintf(
			"gcloud compute scp %s %s --zone=%s --tunnel-through-iap",
			common.ShellQuoteSSHKeyPath(params.SSHKeyPath),
			common.ShellQuote(params.InstanceName+":/tmp"),
			common.ShellQuote(params.Zone),
		),
		fmt.Sprintf(
			"gcloud compute ssh %s --zone=%s --tunnel-through-iap",
			common.ShellQuote(params.InstanceName),
			common.ShellQuote(params.Zone),
		),
	)

	return strings.Join(commands, " && \\\n"), nil
}
