package creategceingresssshfirewallrule

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params は VPC 内 SSH 用 firewall rule 作成コマンド生成に必要な値。
type Params struct {
	RuleName     string
	AllowRule    string
	SourceRanges string
}

// Service は create-gce-ingress-ssh-firewall-rule operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build は VPC 内 SSH 用 firewall rule 作成コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.RuleName) {
		return "", fmt.Errorf("rule-name は必須です")
	}
	if common.IsBlank(params.AllowRule) {
		return "", fmt.Errorf("allow-rule は必須です")
	}
	if common.IsBlank(params.SourceRanges) {
		return "", fmt.Errorf("source-ranges は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute firewall-rules create %s --allow=%s --source-ranges=%s",
		common.ShellQuote(params.RuleName),
		common.ShellQuote(params.AllowRule),
		common.ShellQuote(params.SourceRanges),
	), nil
}
