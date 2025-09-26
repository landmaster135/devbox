package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	// OperationDeleteInstance は Cloud SQL インスタンス削除フローを表す。
	OperationDeleteInstance = "delete-instance"
	// OperationPatchDeletionProtection は削除保護の有効・無効化処理を表す。
	OperationPatchDeletionProtection = "patch-deletion-protection"
	// OperationPatchActivationPolicy は起動ポリシーの変更処理を表す。
	OperationPatchActivationPolicy = "patch-activation-policy"
	// OperationStartInstance はインスタンスの起動処理を表す。
	OperationStartInstance = "start-instance"
	// OperationStopInstance はインスタンスの停止処理を表す。
	OperationStopInstance = "stop-instance"
)

var supportedOperations = []string{
	OperationDeleteInstance,
	OperationPatchActivationPolicy,
	OperationPatchDeletionProtection,
	OperationStartInstance,
	OperationStopInstance,
}

// Config は CLI から渡されるパラメータを保持する。
type Config struct {
	Operation              string
	Help                   bool
	InstanceName           string
	DeletionProtectionMode string
	ActivationPolicy       string
}

// ParseFlags は標準フラグパーサーで CLI 引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は注入されたパーサーを利用して CLI 引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")
	parser.StringVar(&cfg.InstanceName, "instance-name", "", "対象の Cloud SQL インスタンス名")
	parser.StringVar(&cfg.DeletionProtectionMode, "deletion-protection-mode", "", "削除保護の有効/無効化 (enable|disable)")
	parser.StringVar(&cfg.ActivationPolicy, "activation-policy", "", "起動ポリシー (always|never)")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	normalizeConfig(cfg)

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func normalizeConfig(cfg *Config) {
	cfg.Operation = strings.TrimSpace(cfg.Operation)
	cfg.InstanceName = strings.TrimSpace(cfg.InstanceName)
	cfg.DeletionProtectionMode = strings.ToLower(strings.TrimSpace(cfg.DeletionProtectionMode))
	cfg.ActivationPolicy = strings.ToLower(strings.TrimSpace(cfg.ActivationPolicy))
}

func validateConfig(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}
	if !isSupportedOperation(cfg.Operation) {
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	if requiresInstance(cfg.Operation) && cfg.InstanceName == "" {
		return fmt.Errorf("instance-name パラメータは必須です")
	}

	switch cfg.Operation {
	case OperationPatchDeletionProtection:
		if cfg.DeletionProtectionMode == "" {
			return fmt.Errorf("deletion-protection-mode パラメータは必須です")
		}
		if cfg.DeletionProtectionMode != "enable" && cfg.DeletionProtectionMode != "disable" {
			return fmt.Errorf("deletion-protection-mode には enable または disable を指定してください")
		}
	case OperationPatchActivationPolicy:
		if cfg.ActivationPolicy == "" {
			return fmt.Errorf("activation-policy パラメータは必須です")
		}
		if cfg.ActivationPolicy != "always" && cfg.ActivationPolicy != "never" {
			return fmt.Errorf("activation-policy には always または never を指定してください")
		}
	}

	return nil
}

func requiresInstance(operation string) bool {
	switch operation {
	case OperationDeleteInstance, OperationPatchDeletionProtection, OperationPatchActivationPolicy, OperationStartInstance, OperationStopInstance:
		return true
	default:
		return false
	}
}

func isSupportedOperation(operation string) bool {
	for _, candidate := range supportedOperations {
		if candidate == operation {
			return true
		}
	}
	return false
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Cloud SQL 関連の gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(supportedOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -help\n        このヘルプを表示\n")
	fmt.Fprintf(os.Stderr, "  -instance-name string\n        対象の Cloud SQL インスタンス名 (operation に応じて必須)\n\n")

	fmt.Fprintf(os.Stderr, "operation ごとの追加パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  %s: 追加パラメータ不要 (instance-name のみ)\n", OperationStartInstance)
	fmt.Fprintf(os.Stderr, "  %s: 追加パラメータ不要 (instance-name のみ)\n", OperationStopInstance)
	fmt.Fprintf(os.Stderr, "  %s: -deletion-protection-mode (enable|disable)\n", OperationPatchDeletionProtection)
	fmt.Fprintf(os.Stderr, "  %s: -activation-policy (always|never)\n", OperationPatchActivationPolicy)
	fmt.Fprintf(os.Stderr, "  %s: 内部で start-instance と patch-deletion-protection を順番に実行後に削除\n", OperationDeleteInstance)

	fmt.Fprintf(os.Stderr, "\n使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -instance-name=my-instance\n", os.Args[0], OperationStartInstance)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -instance-name=my-instance -deletion-protection-mode=disable\n", os.Args[0], OperationPatchDeletionProtection)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -instance-name=my-instance\n", os.Args[0], OperationDeleteInstance)
}

func init() {
	sort.Strings(supportedOperations)
}
