package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	// OperationCreateGCEInstance は VM インスタンス作成コマンドを生成する操作。
	OperationCreateGCEInstance = "create-gce-instance"
	// OperationCreateGCERouterAndNAT は Router/NAT 作成コマンドを生成する操作。
	OperationCreateGCERouterAndNAT = "create-gce-router-and-nat"
	// OperationCreateGCEIAPSSHFirewallRule は IAP SSH 用 firewall ルール作成コマンドを生成する操作。
	OperationCreateGCEIAPSSHFirewallRule = "create-gce-iap-ssh-firewall-rule"
	// OperationCreateGCEIngressSSHFirewallRule は VPC 内 SSH 用 firewall ルール作成コマンドを生成する操作。
	OperationCreateGCEIngressSSHFirewallRule = "create-gce-ingress-ssh-firewall-rule"
	// OperationListGCloudInstances はインスタンス一覧取得コマンドを生成する操作。
	OperationListGCloudInstances = "list-gcloud-instances"
)

var validOperations = []string{
	OperationCreateGCEIngressSSHFirewallRule,
	OperationCreateGCEIAPSSHFirewallRule,
	OperationCreateGCEInstance,
	OperationCreateGCERouterAndNAT,
	OperationListGCloudInstances,
}

const (
	defaultZone         = "us-central1-a"
	defaultMachineType  = "e2-medium"
	defaultBootDiskSize = "100GB"
	defaultBootDiskType = "pd-balanced"

	defaultRegion  = "us-central1"
	defaultNetwork = "default"
	defaultNATName = "nat1"

	defaultIAPRuleName     = "allow-ssh-ingress-from-iap"
	defaultDirection       = "INGRESS"
	defaultAction          = "allow"
	defaultRules           = "tcp:22"
	defaultIAPSourceRanges = "35.235.240.0/20"

	defaultIngressRuleName     = "allow-ingress-ssh"
	defaultAllowRule           = "tcp:22"
	defaultIngressSourceRanges = "10.0.0.0/8"

	defaultInstanceListFormat = "table(name, zone.basename(), scheduling.preemptible.yesno(yes=true, no=''), networkInterfaces.internal_ip():label=INTERNAL_IP, external_ip():label=EXTERNAL_IP, status)"
)

// Config は CLI 引数から得られる設定値を保持する。
type Config struct {
	Operation string
	Help      bool

	InstanceName string
	Zone         string
	MachineType  string
	BootDiskSize string
	BootDiskType string

	RouterName string
	Region     string
	Network    string
	NATName    string

	RuleName     string
	Direction    string
	Action       string
	Rules        string
	SourceRanges string
	AllowRule    string

	Filter string
	Format string
}

// ParseFlags は標準のフラグパーサーで引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたパーサーを用いて引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(validOperations, ", ")))
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")

	parser.StringVar(&cfg.InstanceName, "instance-name", "", "作成する GCE インスタンス名")
	parser.StringVar(&cfg.Zone, "zone", "", "インスタンスのゾーン (例: us-central1-a)")
	parser.StringVar(&cfg.MachineType, "machine-type", "", "マシンタイプ (例: e2-medium)")
	parser.StringVar(&cfg.BootDiskSize, "boot-disk-size", "", "ブートディスクサイズ (例: 100GB)")
	parser.StringVar(&cfg.BootDiskType, "boot-disk-type", "", "ブートディスクタイプ (例: pd-balanced)")

	parser.StringVar(&cfg.RouterName, "router-name", "", "作成する Cloud Router 名")
	parser.StringVar(&cfg.Region, "region", "", "リージョン (例: us-central1)")
	parser.StringVar(&cfg.Network, "network", "", "対象 VPC ネットワーク")
	parser.StringVar(&cfg.NATName, "nat-name", "", "作成する Cloud NAT 名")

	parser.StringVar(&cfg.RuleName, "rule-name", "", "作成する firewall rule 名")
	parser.StringVar(&cfg.Direction, "direction", "", "firewall rule の方向 (例: INGRESS)")
	parser.StringVar(&cfg.Action, "action", "", "firewall rule のアクション (例: allow)")
	parser.StringVar(&cfg.Rules, "rules", "", "許可するルール (例: tcp:22)")
	parser.StringVar(&cfg.SourceRanges, "source-ranges", "", "許可する送信元 CIDR")
	parser.StringVar(&cfg.AllowRule, "allow-rule", "", "--allow に渡すルール (例: tcp:22)")

	parser.StringVar(&cfg.Filter, "filter", "", "インスタンス一覧フィルタ")
	parser.StringVar(&cfg.Format, "format", "", "インスタンス一覧表示フォーマット")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	if len(parser.Args()) > 0 {
		return nil, fmt.Errorf("位置引数はサポートしていません: %s", strings.Join(parser.Args(), ", "))
	}

	normalizeConfig(cfg)

	if err := validateOperation(cfg.Operation); err != nil {
		return nil, err
	}

	applyDefaults(cfg)

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func normalizeConfig(cfg *Config) {
	cfg.Operation = strings.TrimSpace(cfg.Operation)
	cfg.InstanceName = strings.TrimSpace(cfg.InstanceName)
	cfg.Zone = strings.TrimSpace(cfg.Zone)
	cfg.MachineType = strings.TrimSpace(cfg.MachineType)
	cfg.BootDiskSize = strings.TrimSpace(cfg.BootDiskSize)
	cfg.BootDiskType = strings.TrimSpace(cfg.BootDiskType)
	cfg.RouterName = strings.TrimSpace(cfg.RouterName)
	cfg.Region = strings.TrimSpace(cfg.Region)
	cfg.Network = strings.TrimSpace(cfg.Network)
	cfg.NATName = strings.TrimSpace(cfg.NATName)
	cfg.RuleName = strings.TrimSpace(cfg.RuleName)
	cfg.Direction = strings.TrimSpace(cfg.Direction)
	cfg.Action = strings.TrimSpace(cfg.Action)
	cfg.Rules = strings.TrimSpace(cfg.Rules)
	cfg.SourceRanges = strings.TrimSpace(cfg.SourceRanges)
	cfg.AllowRule = strings.TrimSpace(cfg.AllowRule)
	cfg.Filter = strings.TrimSpace(cfg.Filter)
	cfg.Format = strings.TrimSpace(cfg.Format)
}

func validateOperation(operation string) error {
	if operation == "" {
		return fmt.Errorf("operation は必須です")
	}

	if !isValidOperation(operation) {
		return fmt.Errorf("未対応のoperationです: %s", operation)
	}

	return nil
}

func applyDefaults(cfg *Config) {
	switch cfg.Operation {
	case OperationCreateGCEInstance:
		if cfg.Zone == "" {
			cfg.Zone = defaultZone
		}
		if cfg.MachineType == "" {
			cfg.MachineType = defaultMachineType
		}
		if cfg.BootDiskSize == "" {
			cfg.BootDiskSize = defaultBootDiskSize
		}
		if cfg.BootDiskType == "" {
			cfg.BootDiskType = defaultBootDiskType
		}
	case OperationCreateGCERouterAndNAT:
		if cfg.Region == "" {
			cfg.Region = defaultRegion
		}
		if cfg.Network == "" {
			cfg.Network = defaultNetwork
		}
		if cfg.NATName == "" {
			cfg.NATName = defaultNATName
		}
	case OperationCreateGCEIAPSSHFirewallRule:
		if cfg.RuleName == "" {
			cfg.RuleName = defaultIAPRuleName
		}
		if cfg.Direction == "" {
			cfg.Direction = defaultDirection
		}
		if cfg.Action == "" {
			cfg.Action = defaultAction
		}
		if cfg.Rules == "" {
			cfg.Rules = defaultRules
		}
		if cfg.SourceRanges == "" {
			cfg.SourceRanges = defaultIAPSourceRanges
		}
	case OperationCreateGCEIngressSSHFirewallRule:
		if cfg.RuleName == "" {
			cfg.RuleName = defaultIngressRuleName
		}
		if cfg.AllowRule == "" {
			cfg.AllowRule = defaultAllowRule
		}
		if cfg.SourceRanges == "" {
			cfg.SourceRanges = defaultIngressSourceRanges
		}
	case OperationListGCloudInstances:
		if cfg.Format == "" {
			cfg.Format = defaultInstanceListFormat
		}
	}
}

func validateConfig(cfg *Config) error {
	switch cfg.Operation {
	case OperationCreateGCEInstance:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
	case OperationCreateGCERouterAndNAT:
		if cfg.RouterName == "" {
			return fmt.Errorf("router-name は必須です")
		}
	}

	return nil
}

func isValidOperation(operation string) bool {
	for _, candidate := range validOperations {
		if candidate == operation {
			return true
		}
	}
	return false
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Compute Engine 向け gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (%s)\n", strings.Join(validOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "create-gce-instance:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (default: %s)\n", defaultZone)
	fmt.Fprintf(os.Stderr, "  -machine-type string (default: %s)\n", defaultMachineType)
	fmt.Fprintf(os.Stderr, "  -boot-disk-size string (default: %s)\n", defaultBootDiskSize)
	fmt.Fprintf(os.Stderr, "  -boot-disk-type string (default: %s)\n\n", defaultBootDiskType)

	fmt.Fprintf(os.Stderr, "create-gce-router-and-nat:\n")
	fmt.Fprintf(os.Stderr, "  -router-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -region string (default: %s)\n", defaultRegion)
	fmt.Fprintf(os.Stderr, "  -network string (default: %s)\n", defaultNetwork)
	fmt.Fprintf(os.Stderr, "  -nat-name string (default: %s)\n\n", defaultNATName)

	fmt.Fprintf(os.Stderr, "create-gce-iap-ssh-firewall-rule:\n")
	fmt.Fprintf(os.Stderr, "  -rule-name string (default: %s)\n", defaultIAPRuleName)
	fmt.Fprintf(os.Stderr, "  -direction string (default: %s)\n", defaultDirection)
	fmt.Fprintf(os.Stderr, "  -action string (default: %s)\n", defaultAction)
	fmt.Fprintf(os.Stderr, "  -rules string (default: %s)\n", defaultRules)
	fmt.Fprintf(os.Stderr, "  -source-ranges string (default: %s)\n\n", defaultIAPSourceRanges)

	fmt.Fprintf(os.Stderr, "create-gce-ingress-ssh-firewall-rule:\n")
	fmt.Fprintf(os.Stderr, "  -rule-name string (default: %s)\n", defaultIngressRuleName)
	fmt.Fprintf(os.Stderr, "  -allow-rule string (default: %s)\n", defaultAllowRule)
	fmt.Fprintf(os.Stderr, "  -source-ranges string (default: %s)\n\n", defaultIngressSourceRanges)

	fmt.Fprintf(os.Stderr, "list-gcloud-instances:\n")
	fmt.Fprintf(os.Stderr, "  -filter string\n")
	fmt.Fprintf(os.Stderr, "  -format string (default: table)\n\n")

	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=create-gce-instance -instance-name=my-vm -zone=us-central1-a -machine-type=e2-medium\n", os.Args[0])
}

func init() {
	sort.Strings(validOperations)
}
