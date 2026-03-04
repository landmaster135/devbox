package usecases

import (
	"fmt"
	"strings"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_compute/config"
)

// Service は GCE 向け gcloud コマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// CreateGCEInstanceParams はインスタンス作成コマンド生成に必要な値。
type CreateGCEInstanceParams struct {
	InstanceName string
	Zone         string
	MachineType  string
	BootDiskSize string
	BootDiskType string
}

// CreateGCERouterAndNATParams は Router/NAT 作成コマンド生成に必要な値。
type CreateGCERouterAndNATParams struct {
	RouterName string
	Region     string
	Network    string
	NATName    string
}

// CreateGCEIAPSSHFirewallRuleParams は IAP SSH 用 firewall rule 作成コマンド生成に必要な値。
type CreateGCEIAPSSHFirewallRuleParams struct {
	RuleName     string
	Direction    string
	Action       string
	Rules        string
	SourceRanges string
}

// CreateGCEIngressSSHFirewallRuleParams は VPC 内 SSH 用 firewall rule 作成コマンド生成に必要な値。
type CreateGCEIngressSSHFirewallRuleParams struct {
	RuleName     string
	AllowRule    string
	SourceRanges string
}

// ListGCloudInstancesParams はインスタンス一覧コマンド生成に必要な値。
type ListGCloudInstancesParams struct {
	Filter string
	Format string
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
	if strings.TrimSpace(params.InstanceName) == "" {
		return "", fmt.Errorf("instance-name は必須です")
	}
	if strings.TrimSpace(params.Zone) == "" {
		return "", fmt.Errorf("zone は必須です")
	}
	if strings.TrimSpace(params.MachineType) == "" {
		return "", fmt.Errorf("machine-type は必須です")
	}
	if strings.TrimSpace(params.BootDiskSize) == "" {
		return "", fmt.Errorf("boot-disk-size は必須です")
	}
	if strings.TrimSpace(params.BootDiskType) == "" {
		return "", fmt.Errorf("boot-disk-type は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute instances create %s --zone=%s --machine-type=%s --no-address --boot-disk-size=%s --boot-disk-type=%s",
		shellQuote(params.InstanceName),
		shellQuote(params.Zone),
		shellQuote(params.MachineType),
		shellQuote(params.BootDiskSize),
		shellQuote(params.BootDiskType),
	), nil
}

// BuildCreateGCERouterAndNATCommand は Router/NAT 作成コマンドを生成する。
func (s *Service) BuildCreateGCERouterAndNATCommand(params CreateGCERouterAndNATParams) (string, error) {
	if strings.TrimSpace(params.RouterName) == "" {
		return "", fmt.Errorf("router-name は必須です")
	}
	if strings.TrimSpace(params.Region) == "" {
		return "", fmt.Errorf("region は必須です")
	}
	if strings.TrimSpace(params.Network) == "" {
		return "", fmt.Errorf("network は必須です")
	}
	if strings.TrimSpace(params.NATName) == "" {
		return "", fmt.Errorf("nat-name は必須です")
	}

	routerCommand := fmt.Sprintf(
		"gcloud compute routers create %s --region=%s --network=%s",
		shellQuote(params.RouterName),
		shellQuote(params.Region),
		shellQuote(params.Network),
	)

	natCommand := fmt.Sprintf(
		"gcloud compute routers nats create %s --router=%s --region=%s --auto-allocate-nat-external-ips --nat-all-subnet-ip-ranges",
		shellQuote(params.NATName),
		shellQuote(params.RouterName),
		shellQuote(params.Region),
	)

	return fmt.Sprintf("%s && \\\n%s", routerCommand, natCommand), nil
}

// BuildCreateGCEIAPSSHFirewallRuleCommand は IAP SSH 用 firewall rule 作成コマンドを生成する。
func (s *Service) BuildCreateGCEIAPSSHFirewallRuleCommand(params CreateGCEIAPSSHFirewallRuleParams) (string, error) {
	if strings.TrimSpace(params.RuleName) == "" {
		return "", fmt.Errorf("rule-name は必須です")
	}
	if strings.TrimSpace(params.Direction) == "" {
		return "", fmt.Errorf("direction は必須です")
	}
	if strings.TrimSpace(params.Action) == "" {
		return "", fmt.Errorf("action は必須です")
	}
	if strings.TrimSpace(params.Rules) == "" {
		return "", fmt.Errorf("rules は必須です")
	}
	if strings.TrimSpace(params.SourceRanges) == "" {
		return "", fmt.Errorf("source-ranges は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute firewall-rules create %s --direction=%s --action=%s --rules=%s --source-ranges=%s",
		shellQuote(params.RuleName),
		shellQuote(params.Direction),
		shellQuote(params.Action),
		shellQuote(params.Rules),
		shellQuote(params.SourceRanges),
	), nil
}

// BuildCreateGCEIngressSSHFirewallRuleCommand は VPC 内 SSH 用 firewall rule 作成コマンドを生成する。
func (s *Service) BuildCreateGCEIngressSSHFirewallRuleCommand(params CreateGCEIngressSSHFirewallRuleParams) (string, error) {
	if strings.TrimSpace(params.RuleName) == "" {
		return "", fmt.Errorf("rule-name は必須です")
	}
	if strings.TrimSpace(params.AllowRule) == "" {
		return "", fmt.Errorf("allow-rule は必須です")
	}
	if strings.TrimSpace(params.SourceRanges) == "" {
		return "", fmt.Errorf("source-ranges は必須です")
	}

	return fmt.Sprintf(
		"gcloud compute firewall-rules create %s --allow=%s --source-ranges=%s",
		shellQuote(params.RuleName),
		shellQuote(params.AllowRule),
		shellQuote(params.SourceRanges),
	), nil
}

// BuildListGCloudInstancesCommand はインスタンス一覧コマンドを生成する。
func (s *Service) BuildListGCloudInstancesCommand(params ListGCloudInstancesParams) (string, error) {
	if strings.TrimSpace(params.Format) == "" {
		return "", fmt.Errorf("format は必須です")
	}

	parts := []string{"gcloud", "compute", "instances", "list"}
	if strings.TrimSpace(params.Filter) != "" {
		parts = append(parts, "--filter="+shellQuote(params.Filter))
	}
	parts = append(parts, "--format="+shellQuote(params.Format))

	return strings.Join(parts, " "), nil
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

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}
