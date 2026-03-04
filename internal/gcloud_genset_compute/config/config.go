package config

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var zonePattern = regexp.MustCompile(`^[a-z0-9-]+$`)

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
	// OperationListDiskTypes はディスクタイプ一覧を取得して表示する操作。
	OperationListDiskTypes = "list-disk-types"
	// OperationListMachineTypes はマシンタイプ一覧を取得して表示する操作。
	OperationListMachineTypes = "list-machine-types"
	// OperationStartGCEInstance はインスタンス起動コマンドを生成する操作。
	OperationStartGCEInstance = "start-gce-instance"
	// OperationStopGCEInstance はインスタンス停止コマンドを生成する操作。
	OperationStopGCEInstance = "stop-gce-instance"
	// OperationRebootGCEInstance はインスタンス再起動コマンドを生成する操作。
	OperationRebootGCEInstance = "reboot-gce-instance"
	// OperationDeleteGCEInstance はインスタンス削除コマンドを生成する操作。
	OperationDeleteGCEInstance = "delete-gce-instance"
	// OperationCopyGCESSHKey は SSH 鍵コピーコマンドを生成する操作。
	OperationCopyGCESSHKey = "copy-gce-ssh-key"
	// OperationConnectGCEInstance はインスタンス SSH 接続コマンドを生成する操作。
	OperationConnectGCEInstance = "connect-gce-instance"
	// OperationSetupGCEFirewallAndSSH は firewall 作成 + SSH 鍵コピー + SSH 接続コマンドを生成する操作。
	OperationSetupGCEFirewallAndSSH = "setup-gce-firewall-and-ssh"
	// OperationSetGCEInstanceMetadataFromYAML は YAML からインスタンス metadata 設定コマンドを生成する操作。
	OperationSetGCEInstanceMetadataFromYAML = "set-gce-instance-metadata-from-yaml"
	// OperationAddStartupScriptToGCEInstance はスタートアップスクリプト登録コマンドを生成する操作。
	OperationAddStartupScriptToGCEInstance = "add-startup-script-to-gce-instance"
	// OperationCreateGCEInstanceWithStartupScript はインスタンス作成 + スタートアップスクリプト登録コマンドを生成する操作。
	OperationCreateGCEInstanceWithStartupScript = "create-gce-instance-with-startup-script"
	// OperationCreateGCEInstanceAndConfigure はインスタンス作成 + metadata 設定 + スタートアップスクリプト登録コマンドを生成する操作。
	OperationCreateGCEInstanceAndConfigure = "create-gce-instance-and-configure"
)

var validOperations = []string{
	OperationConnectGCEInstance,
	OperationCopyGCESSHKey,
	OperationDeleteGCEInstance,
	OperationCreateGCEIngressSSHFirewallRule,
	OperationCreateGCEIAPSSHFirewallRule,
	OperationCreateGCEInstance,
	OperationCreateGCERouterAndNAT,
	OperationCreateGCEInstanceAndConfigure,
	OperationCreateGCEInstanceWithStartupScript,
	OperationListDiskTypes,
	OperationListGCloudInstances,
	OperationListMachineTypes,
	OperationRebootGCEInstance,
	OperationSetGCEInstanceMetadataFromYAML,
	OperationSetupGCEFirewallAndSSH,
	OperationStartGCEInstance,
	OperationStopGCEInstance,
	OperationAddStartupScriptToGCEInstance,
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
	defaultSSHKeyPath         = "$HOME/.ssh/google_compute_engine"
	defaultMetadataYAMLPath   = "cmd/cli/gcloud-genset-compute/metadata/config/env.yml"
	defaultStartupScriptPath  = "cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh"
)

// Config は CLI 引数から得られる設定値を保持する。
type Config struct {
	Operation string
	Help      bool

	InstanceName      string
	Zone              string
	SSHKeyPath        string
	MachineType       string
	BootDiskSize      string
	BootDiskType      string
	MetadataYAMLPath  string
	StartupScriptPath string

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

	Zones          []string
	MinDiskSizeGiB int
	MaxDiskSizeGiB int
	zonesRaw       string
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
	parser.StringVar(&cfg.SSHKeyPath, "ssh-key-path", "", "SSH 秘密鍵ファイルパス (例: $HOME/.ssh/google_compute_engine)")
	parser.StringVar(&cfg.MachineType, "machine-type", "", "マシンタイプ (例: e2-medium)")
	parser.StringVar(&cfg.BootDiskSize, "boot-disk-size", "", "ブートディスクサイズ (例: 100GB)")
	parser.StringVar(&cfg.BootDiskType, "boot-disk-type", "", "ブートディスクタイプ (例: pd-balanced)")
	parser.StringVar(&cfg.MetadataYAMLPath, "metadata-yaml-path", "", "metadata 設定用 YAML ファイルパス")
	parser.StringVar(&cfg.StartupScriptPath, "startup-script-path", "", "スタートアップスクリプトファイルパス")

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
	parser.StringVar(&cfg.zonesRaw, "zones", "", "対象ゾーン一覧 (カンマ区切り, 例: asia-southeast3-a,asia-southeast3-b)")
	parser.IntVar(&cfg.MinDiskSizeGiB, "min-disk-size-gib", 0, "最小ディスクサイズ (GiB)")
	parser.IntVar(&cfg.MaxDiskSizeGiB, "max-disk-size-gib", 0, "最大ディスクサイズ (GiB)")

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
	cfg.SSHKeyPath = strings.TrimSpace(cfg.SSHKeyPath)
	cfg.MachineType = strings.TrimSpace(cfg.MachineType)
	cfg.BootDiskSize = strings.TrimSpace(cfg.BootDiskSize)
	cfg.BootDiskType = strings.TrimSpace(cfg.BootDiskType)
	cfg.MetadataYAMLPath = strings.TrimSpace(cfg.MetadataYAMLPath)
	cfg.StartupScriptPath = strings.TrimSpace(cfg.StartupScriptPath)
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
	cfg.zonesRaw = strings.TrimSpace(cfg.zonesRaw)
	cfg.Zones = parseCommaSeparatedValues(cfg.zonesRaw)
}

func parseCommaSeparatedValues(value string) []string {
	if value == "" {
		return []string{}
	}

	parts := strings.Split(value, ",")
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		normalized = append(normalized, trimmed)
	}

	return normalized
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
	case OperationCreateGCEInstanceWithStartupScript:
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
		if cfg.StartupScriptPath == "" {
			cfg.StartupScriptPath = defaultStartupScriptPath
		}
	case OperationCreateGCEInstanceAndConfigure:
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
		if cfg.MetadataYAMLPath == "" {
			cfg.MetadataYAMLPath = defaultMetadataYAMLPath
		}
		if cfg.StartupScriptPath == "" {
			cfg.StartupScriptPath = defaultStartupScriptPath
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
	case OperationCopyGCESSHKey:
		if cfg.Zone == "" {
			cfg.Zone = defaultZone
		}
		if cfg.SSHKeyPath == "" {
			cfg.SSHKeyPath = defaultSSHKeyPath
		}
	case OperationConnectGCEInstance:
		if cfg.Zone == "" {
			cfg.Zone = defaultZone
		}
	case OperationSetGCEInstanceMetadataFromYAML:
		if cfg.Zone == "" {
			cfg.Zone = defaultZone
		}
		if cfg.MetadataYAMLPath == "" {
			cfg.MetadataYAMLPath = defaultMetadataYAMLPath
		}
	case OperationAddStartupScriptToGCEInstance:
		if cfg.Zone == "" {
			cfg.Zone = defaultZone
		}
		if cfg.StartupScriptPath == "" {
			cfg.StartupScriptPath = defaultStartupScriptPath
		}
	case OperationSetupGCEFirewallAndSSH:
		if cfg.Zone == "" {
			cfg.Zone = defaultZone
		}
		if cfg.SSHKeyPath == "" {
			cfg.SSHKeyPath = defaultSSHKeyPath
		}
	}
}

func validateConfig(cfg *Config) error {
	switch cfg.Operation {
	case OperationCreateGCEInstance:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
	case OperationCreateGCEInstanceWithStartupScript:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
	case OperationCreateGCEInstanceAndConfigure:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
	case OperationCreateGCERouterAndNAT:
		if cfg.RouterName == "" {
			return fmt.Errorf("router-name は必須です")
		}
	case OperationStartGCEInstance, OperationStopGCEInstance, OperationRebootGCEInstance, OperationDeleteGCEInstance:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
		if cfg.Zone == "" {
			return fmt.Errorf("zone は必須です")
		}
	case OperationCopyGCESSHKey:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
		if cfg.Zone == "" {
			return fmt.Errorf("zone は必須です")
		}
		if cfg.SSHKeyPath == "" {
			return fmt.Errorf("ssh-key-path は必須です")
		}
	case OperationConnectGCEInstance:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
		if cfg.Zone == "" {
			return fmt.Errorf("zone は必須です")
		}
	case OperationSetGCEInstanceMetadataFromYAML:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
		if cfg.Zone == "" {
			return fmt.Errorf("zone は必須です")
		}
		if cfg.MetadataYAMLPath == "" {
			return fmt.Errorf("metadata-yaml-path は必須です")
		}
	case OperationAddStartupScriptToGCEInstance:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
		if cfg.Zone == "" {
			return fmt.Errorf("zone は必須です")
		}
		if cfg.StartupScriptPath == "" {
			return fmt.Errorf("startup-script-path は必須です")
		}
	case OperationSetupGCEFirewallAndSSH:
		if cfg.InstanceName == "" {
			return fmt.Errorf("instance-name は必須です")
		}
		if cfg.Zone == "" {
			return fmt.Errorf("zone は必須です")
		}
		if cfg.SSHKeyPath == "" {
			return fmt.Errorf("ssh-key-path は必須です")
		}
	case OperationListDiskTypes, OperationListMachineTypes:
		if cfg.MinDiskSizeGiB < 0 {
			return fmt.Errorf("min-disk-size-gib は0以上で指定してください")
		}
		if cfg.MaxDiskSizeGiB < 0 {
			return fmt.Errorf("max-disk-size-gib は0以上で指定してください")
		}
		if cfg.MinDiskSizeGiB > 0 && cfg.MaxDiskSizeGiB > 0 && cfg.MinDiskSizeGiB > cfg.MaxDiskSizeGiB {
			return fmt.Errorf("min-disk-size-gib は max-disk-size-gib 以下で指定してください")
		}
		for _, zone := range cfg.Zones {
			if !zonePattern.MatchString(zone) {
				return fmt.Errorf("zones の値が不正です: %s", zone)
			}
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
	fmt.Fprintf(os.Stderr, "Google Compute Engine 向け CLI ツール（操作によりコマンド生成または実行結果表示）\n\n")

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

	fmt.Fprintf(os.Stderr, "create-gce-instance-with-startup-script:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (default: %s)\n", defaultZone)
	fmt.Fprintf(os.Stderr, "  -machine-type string (default: %s)\n", defaultMachineType)
	fmt.Fprintf(os.Stderr, "  -boot-disk-size string (default: %s)\n", defaultBootDiskSize)
	fmt.Fprintf(os.Stderr, "  -boot-disk-type string (default: %s)\n", defaultBootDiskType)
	fmt.Fprintf(os.Stderr, "  -startup-script-path string (default: %s)\n\n", defaultStartupScriptPath)

	fmt.Fprintf(os.Stderr, "create-gce-instance-and-configure:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (default: %s)\n", defaultZone)
	fmt.Fprintf(os.Stderr, "  -machine-type string (default: %s)\n", defaultMachineType)
	fmt.Fprintf(os.Stderr, "  -boot-disk-size string (default: %s)\n", defaultBootDiskSize)
	fmt.Fprintf(os.Stderr, "  -boot-disk-type string (default: %s)\n", defaultBootDiskType)
	fmt.Fprintf(os.Stderr, "  -metadata-yaml-path string (default: %s)\n", defaultMetadataYAMLPath)
	fmt.Fprintf(os.Stderr, "  -startup-script-path string (default: %s)\n\n", defaultStartupScriptPath)

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

	fmt.Fprintf(os.Stderr, "list-disk-types:\n")
	fmt.Fprintf(os.Stderr, "  -zones string (カンマ区切り)\n")
	fmt.Fprintf(os.Stderr, "  -min-disk-size-gib int (0 で指定なし)\n")
	fmt.Fprintf(os.Stderr, "  -max-disk-size-gib int (0 で指定なし)\n")
	fmt.Fprintf(os.Stderr, "        ※ gcloud を実行して一覧結果を表示\n\n")

	fmt.Fprintf(os.Stderr, "list-machine-types:\n")
	fmt.Fprintf(os.Stderr, "  -zones string (カンマ区切り)\n")
	fmt.Fprintf(os.Stderr, "  -min-disk-size-gib int (最大永続ディスクサイズの下限, 0 で指定なし)\n")
	fmt.Fprintf(os.Stderr, "  -max-disk-size-gib int (最大永続ディスクサイズの上限, 0 で指定なし)\n")
	fmt.Fprintf(os.Stderr, "        ※ gcloud を実行して一覧結果を表示\n\n")

	fmt.Fprintf(os.Stderr, "start-gce-instance:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (必須)\n\n")

	fmt.Fprintf(os.Stderr, "stop-gce-instance:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (必須)\n\n")

	fmt.Fprintf(os.Stderr, "reboot-gce-instance:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (必須)\n\n")

	fmt.Fprintf(os.Stderr, "delete-gce-instance:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (必須)\n\n")

	fmt.Fprintf(os.Stderr, "copy-gce-ssh-key:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (default: %s)\n", defaultZone)
	fmt.Fprintf(os.Stderr, "  -ssh-key-path string (default: %s)\n\n", defaultSSHKeyPath)

	fmt.Fprintf(os.Stderr, "connect-gce-instance:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (default: %s)\n\n", defaultZone)

	fmt.Fprintf(os.Stderr, "set-gce-instance-metadata-from-yaml:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (default: %s)\n", defaultZone)
	fmt.Fprintf(os.Stderr, "  -metadata-yaml-path string (default: %s)\n\n", defaultMetadataYAMLPath)

	fmt.Fprintf(os.Stderr, "add-startup-script-to-gce-instance:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (default: %s)\n", defaultZone)
	fmt.Fprintf(os.Stderr, "  -startup-script-path string (default: %s)\n\n", defaultStartupScriptPath)

	fmt.Fprintf(os.Stderr, "setup-gce-firewall-and-ssh:\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string (必須)\n")
	fmt.Fprintf(os.Stderr, "  -zone string (default: %s)\n", defaultZone)
	fmt.Fprintf(os.Stderr, "  -ssh-key-path string (default: %s)\n\n", defaultSSHKeyPath)

	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=create-gce-instance -instance-name=my-vm -zone=us-central1-a -machine-type=e2-medium\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=list-disk-types -zones=asia-southeast3-a -min-disk-size-gib=4 -max-disk-size-gib=65536\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=list-machine-types -zones=asia-southeast3-a,asia-southeast3-b -min-disk-size-gib=100\n", os.Args[0])
}

func init() {
	sort.Strings(validOperations)
}
