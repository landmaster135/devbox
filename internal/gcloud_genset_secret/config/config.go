package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	OperationCreateSecret               = "create-secret"
	OperationAddSecretVersion           = "add-secret-version"
	OperationCreateAndAddSecretVersion  = "create-and-add-secret-version"
	OperationAccessSecretVersion        = "access-secret-version"
	OperationUpdateSecretLabels         = "update-secret-labels"
	OperationUpdateSecretVersionAliases = "update-secret-version-aliases"

	replicationPolicyAutomatic   = "automatic"
	replicationPolicyUserManaged = "user-managed"

	defaultVersion = "latest"
)

// Config は CLI で指定されたパラメータを保持する。
type Config struct {
	Operation         string
	Help              bool
	SecretName        string
	ReplicationPolicy string
	Locations         string
	SecretValue       string
	Version           string
	Labels            string
	AliasOption       string
}

// ParseFlags は標準のフラグパーサーを使用して引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたフラグパーサーで引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{
		ReplicationPolicy: replicationPolicyAutomatic,
		Version:           defaultVersion,
	}

	parser.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: create-secret など)")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")
	parser.StringVar(&cfg.SecretName, "secret-name", "", "対象となるシークレット名")
	parser.StringVar(&cfg.ReplicationPolicy, "replication-policy", replicationPolicyAutomatic, "レプリケーションポリシー (automatic または user-managed)")
	parser.StringVar(&cfg.Locations, "locations", "", "レプリケーションポリシーが user-managed の場合に指定するロケーション (カンマ区切り)")
	parser.StringVar(&cfg.SecretValue, "secret-value", "", "シークレットに設定する値")
	parser.StringVar(&cfg.Version, "version", defaultVersion, "対象とするバージョン (デフォルト: latest)")
	parser.StringVar(&cfg.Labels, "labels", "", "更新するラベル (KEY=VALUE,...)")
	parser.StringVar(&cfg.AliasOption, "alias-option", "", "バージョンエイリアス更新オプション (--clear-version-aliases 等)")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	if strings.TrimSpace(cfg.Operation) == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}

	requireSecret := func() error {
		if strings.TrimSpace(cfg.SecretName) == "" {
			return fmt.Errorf("secret-name パラメータは必須です")
		}
		return nil
	}

	replicationPolicy := strings.TrimSpace(cfg.ReplicationPolicy)
	if replicationPolicy == "" {
		replicationPolicy = replicationPolicyAutomatic
	}
	cfg.ReplicationPolicy = replicationPolicy

	if replicationPolicy != replicationPolicyAutomatic && replicationPolicy != replicationPolicyUserManaged {
		return fmt.Errorf("replication-policy には automatic または user-managed を指定してください")
	}

	switch cfg.Operation {
	case OperationCreateSecret:
		if err := requireSecret(); err != nil {
			return err
		}
		if replicationPolicy == replicationPolicyUserManaged && strings.TrimSpace(cfg.Locations) == "" {
			return fmt.Errorf("replication-policy が user-managed の場合、locations パラメータは必須です")
		}
	case OperationAddSecretVersion:
		if err := requireSecret(); err != nil {
			return err
		}
		if cfg.SecretValue == "" {
			return fmt.Errorf("secret-value パラメータは必須です")
		}
	case OperationCreateAndAddSecretVersion:
		if err := requireSecret(); err != nil {
			return err
		}
		if cfg.SecretValue == "" {
			return fmt.Errorf("secret-value パラメータは必須です")
		}
		if replicationPolicy == replicationPolicyUserManaged && strings.TrimSpace(cfg.Locations) == "" {
			return fmt.Errorf("replication-policy が user-managed の場合、locations パラメータは必須です")
		}
	case OperationAccessSecretVersion:
		if err := requireSecret(); err != nil {
			return err
		}
		if strings.TrimSpace(cfg.Version) == "" {
			cfg.Version = defaultVersion
		}
	case OperationUpdateSecretLabels:
		if err := requireSecret(); err != nil {
			return err
		}
		if strings.TrimSpace(cfg.Labels) == "" {
			return fmt.Errorf("labels パラメータは必須です")
		}
	case OperationUpdateSecretVersionAliases:
		if err := requireSecret(); err != nil {
			return err
		}
		if err := validateAliasOption(cfg.AliasOption); err != nil {
			return err
		}
	default:
		return fmt.Errorf("未対応の操作です: %s", cfg.Operation)
	}

	return nil
}

func validateAliasOption(option string) error {
	trimmed := strings.TrimSpace(option)
	if trimmed == "" {
		return fmt.Errorf("alias-option パラメータは必須です")
	}

	if trimmed == "--clear-version-aliases" {
		return nil
	}
	if strings.HasPrefix(trimmed, "--remove-version-aliases=") {
		return nil
	}
	if strings.HasPrefix(trimmed, "--update-version-aliases=") {
		return nil
	}

	return fmt.Errorf("alias-option には --clear-version-aliases もしくは --remove-version-aliases= / --update-version-aliases= の形式を指定してください")
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud Secret Manager 用の gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n")
	fmt.Fprintf(os.Stderr, "  -secret-name string\n")
	fmt.Fprintf(os.Stderr, "        対象となるシークレット名\n")
	fmt.Fprintf(os.Stderr, "  -replication-policy string\n")
	fmt.Fprintf(os.Stderr, "        レプリケーションポリシー (automatic | user-managed)\n")
	fmt.Fprintf(os.Stderr, "  -locations string\n")
	fmt.Fprintf(os.Stderr, "        user-managed 時のレプリケーションロケーション (カンマ区切り)\n\n")

	fmt.Fprintf(os.Stderr, "create-secret / create-and-add-secret-version 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -secret-value string\n")
	fmt.Fprintf(os.Stderr, "        create-and-add-secret-version で使用するシークレット値\n\n")

	fmt.Fprintf(os.Stderr, "add-secret-version 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -secret-value string\n")
	fmt.Fprintf(os.Stderr, "        新しいバージョンに登録する値\n\n")

	fmt.Fprintf(os.Stderr, "access-secret-version 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -version string\n")
	fmt.Fprintf(os.Stderr, "        取得するバージョン番号 (デフォルト: %s)\n\n", defaultVersion)

	fmt.Fprintf(os.Stderr, "update-secret-labels 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -labels string\n")
	fmt.Fprintf(os.Stderr, "        更新するラベル (KEY=VALUE,...)\n\n")

	fmt.Fprintf(os.Stderr, "update-secret-version-aliases 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -alias-option string\n")
	fmt.Fprintf(os.Stderr, "        --clear-version-aliases などのオプションを指定\n\n")

	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=create-secret -secret-name=my-secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=add-secret-version -secret-name=my-secret -secret-value=\"value\"\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=access-secret-version -secret-name=my-secret -version=3\n", os.Args[0])
}
