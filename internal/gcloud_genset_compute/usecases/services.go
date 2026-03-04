package usecases

import (
	"fmt"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_compute/config"
	creategceiapsshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_iap_ssh_firewall_rule"
	creategceingresssshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_ingress_ssh_firewall_rule"
	creategceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance"
	creategcerouterandnat "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_router_and_nat"
	listgcloudinstances "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_gcloud_instances"
)

// Service は GCE 向け gcloud コマンド生成を担当する。
type Service struct {
	createGCEInstanceOperation               createGCEInstanceOperation
	createGCERouterAndNATOperation           createGCERouterAndNATOperation
	createGCEIAPSSHFirewallRuleOperation     createGCEIAPSSHFirewallRuleOperation
	createGCEIngressSSHFirewallRuleOperation createGCEIngressSSHFirewallRuleOperation
	listGCloudInstancesOperation             listGCloudInstancesOperation
}

// CreateGCEInstanceParams はインスタンス作成コマンド生成に必要な値。
type CreateGCEInstanceParams = creategceinstance.Params

// CreateGCERouterAndNATParams は Router/NAT 作成コマンド生成に必要な値。
type CreateGCERouterAndNATParams = creategcerouterandnat.Params

// CreateGCEIAPSSHFirewallRuleParams は IAP SSH 用 firewall rule 作成コマンド生成に必要な値。
type CreateGCEIAPSSHFirewallRuleParams = creategceiapsshfirewallrule.Params

// CreateGCEIngressSSHFirewallRuleParams は VPC 内 SSH 用 firewall rule 作成コマンド生成に必要な値。
type CreateGCEIngressSSHFirewallRuleParams = creategceingresssshfirewallrule.Params

// ListGCloudInstancesParams はインスタンス一覧コマンド生成に必要な値。
type ListGCloudInstancesParams = listgcloudinstances.Params

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return newServiceWithOperations(
		creategceinstance.NewService(),
		creategcerouterandnat.NewService(),
		creategceiapsshfirewallrule.NewService(),
		creategceingresssshfirewallrule.NewService(),
		listgcloudinstances.NewService(),
	)
}

// BuildCommand は operation に応じた gcloud コマンドを生成する。
func (s *Service) BuildCommand(conf *cfg.Config) (string, error) {
	switch conf.Operation {
	case cfg.OperationCreateGCEInstance:
		return s.BuildCreateGCEInstanceCommand(CreateGCEInstanceParams{
			InstanceName: conf.InstanceName,
			Zone:         conf.Zone,
			MachineType:  conf.MachineType,
			BootDiskSize: conf.BootDiskSize,
			BootDiskType: conf.BootDiskType,
		})
	case cfg.OperationCreateGCERouterAndNAT:
		return s.BuildCreateGCERouterAndNATCommand(CreateGCERouterAndNATParams{
			RouterName: conf.RouterName,
			Region:     conf.Region,
			Network:    conf.Network,
			NATName:    conf.NATName,
		})
	case cfg.OperationCreateGCEIAPSSHFirewallRule:
		return s.BuildCreateGCEIAPSSHFirewallRuleCommand(CreateGCEIAPSSHFirewallRuleParams{
			RuleName:     conf.RuleName,
			Direction:    conf.Direction,
			Action:       conf.Action,
			Rules:        conf.Rules,
			SourceRanges: conf.SourceRanges,
		})
	case cfg.OperationCreateGCEIngressSSHFirewallRule:
		return s.BuildCreateGCEIngressSSHFirewallRuleCommand(CreateGCEIngressSSHFirewallRuleParams{
			RuleName:     conf.RuleName,
			AllowRule:    conf.AllowRule,
			SourceRanges: conf.SourceRanges,
		})
	case cfg.OperationListGCloudInstances:
		return s.BuildListGCloudInstancesCommand(ListGCloudInstancesParams{
			Filter: conf.Filter,
			Format: conf.Format,
		})
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", conf.Operation)
	}
}

// BuildCreateGCEInstanceCommand はインスタンス作成コマンドを生成する。
func (s *Service) BuildCreateGCEInstanceCommand(params CreateGCEInstanceParams) (string, error) {
	return s.createGCEInstanceOperation.Build(params)
}

// BuildCreateGCERouterAndNATCommand は Router/NAT 作成コマンドを生成する。
func (s *Service) BuildCreateGCERouterAndNATCommand(params CreateGCERouterAndNATParams) (string, error) {
	return s.createGCERouterAndNATOperation.Build(params)
}

// BuildCreateGCEIAPSSHFirewallRuleCommand は IAP SSH 用 firewall rule 作成コマンドを生成する。
func (s *Service) BuildCreateGCEIAPSSHFirewallRuleCommand(params CreateGCEIAPSSHFirewallRuleParams) (string, error) {
	return s.createGCEIAPSSHFirewallRuleOperation.Build(params)
}

// BuildCreateGCEIngressSSHFirewallRuleCommand は VPC 内 SSH 用 firewall rule 作成コマンドを生成する。
func (s *Service) BuildCreateGCEIngressSSHFirewallRuleCommand(params CreateGCEIngressSSHFirewallRuleParams) (string, error) {
	return s.createGCEIngressSSHFirewallRuleOperation.Build(params)
}

// BuildListGCloudInstancesCommand はインスタンス一覧コマンドを生成する。
func (s *Service) BuildListGCloudInstancesCommand(params ListGCloudInstancesParams) (string, error) {
	return s.listGCloudInstancesOperation.Build(params)
}

// PrintHighlightedCommand は生成されたコマンドを見やすい形式で出力する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}
