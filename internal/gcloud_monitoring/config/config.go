package config

import (
	"fmt"
	"os"
)

// Config はCLIパラメータの設定を保持する構造体
type Config struct {
	Operation        string
	Project          string
	Location         string
	Service          string
	ServiceAccountID string
	Help             bool
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: create-dashboard-for-cloud-run)")
	parser.StringVar(&cfg.Project, "project", "", "Google Cloud プロジェクトID (必須)")
	parser.StringVar(&cfg.Location, "location", "", "Cloud Run サービスのロケーション (必須)")
	parser.StringVar(&cfg.Service, "service", "", "Cloud Run サービス名 (必須)")
	parser.StringVar(&cfg.ServiceAccountID, "service-account-id", "", "サービスアカウントID (オプション)")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	// 必須パラメータの検証
	if cfg.Operation == "" {
		return nil, fmt.Errorf("operation パラメータは必須です")
	}
	if cfg.Project == "" {
		return nil, fmt.Errorf("project パラメータは必須です")
	}
	if cfg.Location == "" {
		return nil, fmt.Errorf("location パラメータは必須です")
	}
	if cfg.Service == "" {
		return nil, fmt.Errorf("service パラメータは必須です")
	}

	// サポートされている操作の検証
	if cfg.Operation != "create-dashboard-for-cloud-run" {
		return nil, fmt.Errorf("未対応の操作です: %s", cfg.Operation)
	}

	return cfg, nil
}

// PrintUsage はヘルプメッセージを表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud Run サービスのモニタリング用ダッシュボードを作成するツール\n\n")
	fmt.Fprintf(os.Stderr, "必須パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (create-dashboard-for-cloud-run)\n")
	fmt.Fprintf(os.Stderr, "  -project string\n")
	fmt.Fprintf(os.Stderr, "        Google Cloud プロジェクトID\n")
	fmt.Fprintf(os.Stderr, "  -location string\n")
	fmt.Fprintf(os.Stderr, "        Cloud Run サービスのロケーション (例: us-central1)\n")
	fmt.Fprintf(os.Stderr, "  -service string\n")
	fmt.Fprintf(os.Stderr, "        Cloud Run サービス名\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "オプションパラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -service-account-id string\n")
	fmt.Fprintf(os.Stderr, "        サービスアカウントID (省略時は現在の認証情報を使用)\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプメッセージを表示\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=create-dashboard-for-cloud-run -project=my-project -location=us-central1 -service=my-service\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=create-dashboard-for-cloud-run -project=my-project -location=us-central1 -service=my-service -service-account-id=monitoring-sa\n", os.Args[0])
}

// GetServiceAccountEmail はサービスアカウントの完全なメールアドレスを返す
func (c *Config) GetServiceAccountEmail() string {
	if c.ServiceAccountID == "" {
		return ""
	}
	return fmt.Sprintf("%s@%s.iam.gserviceaccount.com", c.ServiceAccountID, c.Project)
}
