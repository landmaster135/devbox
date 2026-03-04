package usecases

import (
	"fmt"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_compute/config"
	connectgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/connect_gce_instance"
	copygcesshkey "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/copy_gce_ssh_key"
	creategceiapsshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_iap_ssh_firewall_rule"
	creategceingresssshfirewallrule "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_ingress_ssh_firewall_rule"
	creategceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_instance"
	creategcerouterandnat "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/create_gce_router_and_nat"
	deletegceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/delete_gce_instance"
	listgcloudinstances "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/list_gcloud_instances"
	rebootgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/reboot_gce_instance"
	setupgcefirewallandssh "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/setup_gce_firewall_and_ssh"
	startgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/start_gce_instance"
	stopgceinstance "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/operations/stop_gce_instance"
)

// Service は GCE 向け gcloud コマンド生成を担当する。
type Service struct {
	createGCEInstanceOperation               createGCEInstanceOperation
	createGCERouterAndNATOperation           createGCERouterAndNATOperation
	createGCEIAPSSHFirewallRuleOperation     createGCEIAPSSHFirewallRuleOperation
	createGCEIngressSSHFirewallRuleOperation createGCEIngressSSHFirewallRuleOperation
	listGCloudInstancesOperation             listGCloudInstancesOperation
	startGCEInstanceOperation                startGCEInstanceOperation
	stopGCEInstanceOperation                 stopGCEInstanceOperation
	rebootGCEInstanceOperation               rebootGCEInstanceOperation
	deleteGCEInstanceOperation               deleteGCEInstanceOperation
	copyGCESSHKeyOperation                   copyGCESSHKeyOperation
	connectGCEInstanceOperation              connectGCEInstanceOperation
	setupGCEFirewallAndSSHOperation          setupGCEFirewallAndSSHOperation
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

// StartGCEInstanceParams はインスタンス起動コマンド生成に必要な値。
type StartGCEInstanceParams = startgceinstance.Params

// StopGCEInstanceParams はインスタンス停止コマンド生成に必要な値。
type StopGCEInstanceParams = stopgceinstance.Params

// RebootGCEInstanceParams はインスタンス再起動コマンド生成に必要な値。
type RebootGCEInstanceParams = rebootgceinstance.Params

// DeleteGCEInstanceParams はインスタンス削除コマンド生成に必要な値。
type DeleteGCEInstanceParams = deletegceinstance.Params

// CopyGCESSHKeyParams は SSH 鍵コピーコマンド生成に必要な値。
type CopyGCESSHKeyParams = copygcesshkey.Params

// ConnectGCEInstanceParams は GCE SSH 接続コマンド生成に必要な値。
type ConnectGCEInstanceParams = connectgceinstance.Params

// SetupGCEFirewallAndSSHParams は firewall 作成 + SSH 鍵コピー + SSH 接続コマンド生成に必要な値。
type SetupGCEFirewallAndSSHParams = setupgcefirewallandssh.Params

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return newServiceWithOperations(
		creategceinstance.NewService(),
		creategcerouterandnat.NewService(),
		creategceiapsshfirewallrule.NewService(),
		creategceingresssshfirewallrule.NewService(),
		listgcloudinstances.NewService(),
		startgceinstance.NewService(),
		stopgceinstance.NewService(),
		rebootgceinstance.NewService(),
		deletegceinstance.NewService(),
		copygcesshkey.NewService(),
		connectgceinstance.NewService(),
		setupgcefirewallandssh.NewService(),
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
	case cfg.OperationStartGCEInstance:
		return s.BuildStartGCEInstanceCommand(StartGCEInstanceParams{
			InstanceName: conf.InstanceName,
			Zone:         conf.Zone,
		})
	case cfg.OperationStopGCEInstance:
		return s.BuildStopGCEInstanceCommand(StopGCEInstanceParams{
			InstanceName: conf.InstanceName,
			Zone:         conf.Zone,
		})
	case cfg.OperationRebootGCEInstance:
		return s.BuildRebootGCEInstanceCommand(RebootGCEInstanceParams{
			InstanceName: conf.InstanceName,
			Zone:         conf.Zone,
		})
	case cfg.OperationDeleteGCEInstance:
		return s.BuildDeleteGCEInstanceCommand(DeleteGCEInstanceParams{
			InstanceName: conf.InstanceName,
			Zone:         conf.Zone,
		})
	case cfg.OperationCopyGCESSHKey:
		return s.BuildCopyGCESSHKeyCommand(CopyGCESSHKeyParams{
			InstanceName: conf.InstanceName,
			Zone:         conf.Zone,
			SSHKeyPath:   conf.SSHKeyPath,
		})
	case cfg.OperationConnectGCEInstance:
		return s.BuildConnectGCEInstanceCommand(ConnectGCEInstanceParams{
			InstanceName: conf.InstanceName,
			Zone:         conf.Zone,
		})
	case cfg.OperationSetupGCEFirewallAndSSH:
		return s.BuildSetupGCEFirewallAndSSHCommand(SetupGCEFirewallAndSSHParams{
			InstanceName: conf.InstanceName,
			Zone:         conf.Zone,
			SSHKeyPath:   conf.SSHKeyPath,
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

// BuildStartGCEInstanceCommand はインスタンス起動コマンドを生成する。
func (s *Service) BuildStartGCEInstanceCommand(params StartGCEInstanceParams) (string, error) {
	return s.startGCEInstanceOperation.Build(params)
}

// BuildStopGCEInstanceCommand はインスタンス停止コマンドを生成する。
func (s *Service) BuildStopGCEInstanceCommand(params StopGCEInstanceParams) (string, error) {
	return s.stopGCEInstanceOperation.Build(params)
}

// BuildRebootGCEInstanceCommand はインスタンス再起動コマンドを生成する。
func (s *Service) BuildRebootGCEInstanceCommand(params RebootGCEInstanceParams) (string, error) {
	return s.rebootGCEInstanceOperation.Build(params)
}

// BuildDeleteGCEInstanceCommand はインスタンス削除コマンドを生成する。
func (s *Service) BuildDeleteGCEInstanceCommand(params DeleteGCEInstanceParams) (string, error) {
	return s.deleteGCEInstanceOperation.Build(params)
}

// BuildCopyGCESSHKeyCommand は SSH 鍵コピーコマンドを生成する。
func (s *Service) BuildCopyGCESSHKeyCommand(params CopyGCESSHKeyParams) (string, error) {
	return s.copyGCESSHKeyOperation.Build(params)
}

// BuildConnectGCEInstanceCommand は GCE SSH 接続コマンドを生成する。
func (s *Service) BuildConnectGCEInstanceCommand(params ConnectGCEInstanceParams) (string, error) {
	return s.connectGCEInstanceOperation.Build(params)
}

// BuildSetupGCEFirewallAndSSHCommand は firewall 作成 + SSH 鍵コピー + SSH 接続コマンドを生成する。
func (s *Service) BuildSetupGCEFirewallAndSSHCommand(params SetupGCEFirewallAndSSHParams) (string, error) {
	return s.setupGCEFirewallAndSSHOperation.Build(params)
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
