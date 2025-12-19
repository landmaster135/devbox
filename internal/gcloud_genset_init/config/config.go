package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	// OperationAuthLogin は gcloud auth login コマンドを生成する操作を表す。
	OperationAuthLogin = "auth-login"
	// OperationSetProjectConfig は gcloud config set project コマンドを生成する操作を表す。
	OperationSetProjectConfig = "set-project-config"
)

// Config は CLI で指定されたパラメータを保持する。
type Config struct {
	Operation      string
	Help           bool
	ProjectID      string
	AdditionalArgs string
}

// ParseFlags は標準のフラグパーサーを使用して引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたフラグパーサーで引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: auth-login, set-project-config)")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")
	parser.StringVar(&cfg.ProjectID, "project-id", "", "対象のプロジェクトID")
	parser.StringVar(&cfg.AdditionalArgs, "additional-args", "", "gcloud コマンドに追加する引数")

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
	cfg.Operation = strings.TrimSpace(cfg.Operation)
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}

	cfg.ProjectID = strings.TrimSpace(cfg.ProjectID)
	cfg.AdditionalArgs = strings.TrimSpace(cfg.AdditionalArgs)

	switch cfg.Operation {
	case OperationAuthLogin, OperationSetProjectConfig:
		if cfg.ProjectID == "" {
			return fmt.Errorf("project-id パラメータは必須です")
		}
	default:
		return fmt.Errorf("未対応の操作です: %s", cfg.Operation)
	}

	return nil
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "gcloud 初期設定向けのコマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (auth-login | set-project-config)\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "auth-login 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -project-id string\n")
	fmt.Fprintf(os.Stderr, "        認証対象のプロジェクトID (必須)\n")
	fmt.Fprintf(os.Stderr, "  -additional-args string\n")
	fmt.Fprintf(os.Stderr, "        gcloud auth login に渡す追加引数\n\n")

	fmt.Fprintf(os.Stderr, "set-project-config 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -project-id string\n")
	fmt.Fprintf(os.Stderr, "        設定するプロジェクトID (必須)\n")
	fmt.Fprintf(os.Stderr, "  -additional-args string\n")
	fmt.Fprintf(os.Stderr, "        gcloud config set project に渡す追加引数\n\n")

	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=auth-login -project-id=my-project-123\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=set-project-config -project-id=my-project-123 --additional-args='--quiet'\n", os.Args[0])
}
