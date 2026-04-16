package creategceiapsshfirewallrule

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params は IAP SSH 用 firewall rule 作成コマンド生成に必要な値。
type Params struct {
	RuleName     string
	Direction    string
	Action       string
	Rules        string
	SourceRanges string
}

// Service は create-gce-iap-ssh-firewall-rule operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build は IAP SSH 用 firewall rule 作成コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.RuleName) {
		return "", fmt.Errorf("rule-name は必須です")
	}
	if common.IsBlank(params.Direction) {
		return "", fmt.Errorf("direction は必須です")
	}
	if common.IsBlank(params.Action) {
		return "", fmt.Errorf("action は必須です")
	}
	if common.IsBlank(params.Rules) {
		return "", fmt.Errorf("rules は必須です")
	}
	if common.IsBlank(params.SourceRanges) {
		return "", fmt.Errorf("source-ranges は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute firewall-rules create %s --direction=%s --action=%s --rules=%s --source-ranges=%s",
		common.ShellQuote(params.RuleName),
		common.ShellQuote(params.Direction),
		common.ShellQuote(params.Action),
		common.ShellQuote(params.Rules),
		common.ShellQuote(params.SourceRanges),
	), nil
}
