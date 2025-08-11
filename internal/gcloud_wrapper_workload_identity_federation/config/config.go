package config

import (
	"fmt"
	"os"
)

// Config はCLIの設定を保持する構造体
type Config struct {
	ProjectID        string
	PoolID           string
	ProviderID       string
	ServiceAccountID string
	Location         string
	PoolDescription  string
	RepoOwner        string
	RepoName         string
	WebhookURL       string
	Help             bool
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	parser := NewDefaultFlagParser()
	return ParseFlagsWithParser(parser)
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	config := &Config{}

	// フラグの定義
	parser.StringVar(&config.ProjectID, "project-id", "", "Google CloudプロジェクトID (必須)")
	parser.StringVar(&config.PoolID, "pool-id", "", "Workload Identity Pool ID (必須)")
	parser.StringVar(&config.ProviderID, "provider-id", "", "OIDC Provider ID (必須)")
	parser.StringVar(&config.ServiceAccountID, "service-account-id", "", "サービスアカウントID (必須)")
	parser.StringVar(&config.Location, "location", "global", "リソースのロケーション (デフォルト: global)")
	parser.StringVar(&config.PoolDescription, "pool-description", "", "Workload Identity Poolの説明 (任意)")
	parser.StringVar(&config.RepoOwner, "repo-owner", "", "GitHubリポジトリのオーナー (必須)")
	parser.StringVar(&config.RepoName, "repo-name", "", "GitHubリポジトリ名 (必須)")
	parser.StringVar(&config.WebhookURL, "webhook-url", "", "Discord WebhookのURL (必須)")
	parser.BoolVar(&config.Help, "help", false, "ヘルプを表示")

	// フラグの解析
	parser.Parse()

	// ヘルプが要求された場合は早期リターン
	if config.Help {
		return config, nil
	}

	// 必須パラメータの検証
	if err := validateRequiredParams(config); err != nil {
		return nil, err
	}

	return config, nil
}

// validateRequiredParams は必須パラメータの検証を行う
func validateRequiredParams(config *Config) error {
	if config.ProjectID == "" {
		return fmt.Errorf("project-idは必須です")
	}
	if config.PoolID == "" {
		return fmt.Errorf("pool-idは必須です")
	}
	if config.ProviderID == "" {
		return fmt.Errorf("provider-idは必須です")
	}
	if config.ServiceAccountID == "" {
		return fmt.Errorf("service-account-idは必須です")
	}
	if config.RepoOwner == "" {
		return fmt.Errorf("repo-ownerは必須です")
	}
	if config.RepoName == "" {
		return fmt.Errorf("repo-nameは必須です")
	}
	if config.WebhookURL == "" {
		return fmt.Errorf("webhook-urlは必須です")
	}
	return nil
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud Workload Identity FederationとGitHub Actions認証の設定を自動化するツールです。\n\n")
	fmt.Fprintf(os.Stderr, "必須オプション:\n")
	fmt.Fprintf(os.Stderr, "  -project-id string\n")
	fmt.Fprintf(os.Stderr, "        Google CloudプロジェクトID\n")
	fmt.Fprintf(os.Stderr, "  -pool-id string\n")
	fmt.Fprintf(os.Stderr, "        Workload Identity Pool ID\n")
	fmt.Fprintf(os.Stderr, "  -provider-id string\n")
	fmt.Fprintf(os.Stderr, "        OIDC Provider ID\n")
	fmt.Fprintf(os.Stderr, "  -service-account-id string\n")
	fmt.Fprintf(os.Stderr, "        サービスアカウントID\n")
	fmt.Fprintf(os.Stderr, "  -repo-owner string\n")
	fmt.Fprintf(os.Stderr, "        GitHubリポジトリのオーナー\n")
	fmt.Fprintf(os.Stderr, "  -repo-name string\n")
	fmt.Fprintf(os.Stderr, "        GitHubリポジトリ名\n")
	fmt.Fprintf(os.Stderr, "  -webhook-url string\n")
	fmt.Fprintf(os.Stderr, "        Discord WebhookのURL\n")
	fmt.Fprintf(os.Stderr, "\n任意オプション:\n")
	fmt.Fprintf(os.Stderr, "  -location string\n")
	fmt.Fprintf(os.Stderr, "        リソースのロケーション (デフォルト: global)\n")
	fmt.Fprintf(os.Stderr, "  -pool-description string\n")
	fmt.Fprintf(os.Stderr, "        Workload Identity Poolの説明\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n")
	fmt.Fprintf(os.Stderr, "\n使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s \\\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "    -project-id my-project \\\n")
	fmt.Fprintf(os.Stderr, "    -pool-id github-pool \\\n")
	fmt.Fprintf(os.Stderr, "    -provider-id github-provider \\\n")
	fmt.Fprintf(os.Stderr, "    -service-account-id monitoring-sa \\\n")
	fmt.Fprintf(os.Stderr, "    -repo-owner myorg \\\n")
	fmt.Fprintf(os.Stderr, "    -repo-name myrepo \\\n")
	fmt.Fprintf(os.Stderr, "    -webhook-url \"https://discord.com/api/webhooks/...\"\n")
}
