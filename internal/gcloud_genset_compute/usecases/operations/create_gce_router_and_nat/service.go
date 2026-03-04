package creategcerouterandnat

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params は Router/NAT 作成コマンド生成に必要な値。
type Params struct {
	RouterName string
	Region     string
	Network    string
	NATName    string
}

// Service は create-gce-router-and-nat operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build は Router/NAT 作成コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.RouterName) {
		return "", fmt.Errorf("router-name は必須です")
	}
	if common.IsBlank(params.Region) {
		return "", fmt.Errorf("region は必須です")
	}
	if common.IsBlank(params.Network) {
		return "", fmt.Errorf("network は必須です")
	}
	if common.IsBlank(params.NATName) {
		return "", fmt.Errorf("nat-name は必須です")
	}

	routerCommand := fmt.Sprintf(
		"gcloud compute routers create %s --region=%s --network=%s",
		common.ShellQuote(params.RouterName),
		common.ShellQuote(params.Region),
		common.ShellQuote(params.Network),
	)

	natCommand := fmt.Sprintf(
		"gcloud compute routers nats create %s --router=%s --region=%s --auto-allocate-nat-external-ips --nat-all-subnet-ip-ranges",
		common.ShellQuote(params.NATName),
		common.ShellQuote(params.RouterName),
		common.ShellQuote(params.Region),
	)

	return fmt.Sprintf("%s && \\\n%s", routerCommand, natCommand), nil
}
