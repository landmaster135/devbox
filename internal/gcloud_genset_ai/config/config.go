package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	// OperationUndeployProcessorVersion は Document AI のプロセッサバージョンをアンデプロイする操作を表す。
	OperationUndeployProcessorVersion = "undeploy-processor-version"
)

// Config は CLI で指定されたパラメータを保持する。
type Config struct {
	Operation     string
	Help          bool
	Region        string
	ProjectNumber string
	ProcessorID   string
	VersionID     string
}

// ParseFlags は標準のフラグパーサーを使用して引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたフラグパーサーで引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: undeploy-processor-version)")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")
	parser.StringVar(&cfg.Region, "region", "", "Document AI のリージョン")
	parser.StringVar(&cfg.ProjectNumber, "project-number", "", "対象となる Google Cloud プロジェクト番号")
	parser.StringVar(&cfg.ProcessorID, "processor-id", "", "Document AI プロセッサ ID")
	parser.StringVar(&cfg.VersionID, "version-id", "", "アンデプロイ対象のプロセッサバージョン ID")

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

	cfg.Region = strings.TrimSpace(cfg.Region)
	cfg.ProjectNumber = strings.TrimSpace(cfg.ProjectNumber)
	cfg.ProcessorID = strings.TrimSpace(cfg.ProcessorID)
	cfg.VersionID = strings.TrimSpace(cfg.VersionID)

	switch cfg.Operation {
	case OperationUndeployProcessorVersion:
		if cfg.Region == "" {
			return fmt.Errorf("region パラメータは必須です")
		}
		if cfg.ProjectNumber == "" {
			return fmt.Errorf("project-number パラメータは必須です")
		}
		if cfg.ProcessorID == "" {
			return fmt.Errorf("processor-id パラメータは必須です")
		}
		if cfg.VersionID == "" {
			return fmt.Errorf("version-id パラメータは必須です")
		}
	default:
		return fmt.Errorf("未対応の操作です: %s", cfg.Operation)
	}

	return nil
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Document AI 関連の gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (undeploy-processor-version)\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "undeploy-processor-version 操作用:\n")
	fmt.Fprintf(os.Stderr, "  -region string\n")
	fmt.Fprintf(os.Stderr, "        Document AI のリージョン (例: us-central1)\n")
	fmt.Fprintf(os.Stderr, "  -project-number string\n")
	fmt.Fprintf(os.Stderr, "        Google Cloud プロジェクト番号 (必須)\n")
	fmt.Fprintf(os.Stderr, "  -processor-id string\n")
	fmt.Fprintf(os.Stderr, "        Document AI プロセッサ ID (必須)\n")
	fmt.Fprintf(os.Stderr, "  -version-id string\n")
	fmt.Fprintf(os.Stderr, "        アンデプロイ対象のプロセッサバージョン ID (必須)\n\n")

	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=undeploy-processor-version -region=us-central1 -project-number=123456789 -processor-id=abc123def456 -version-id=20240901\n", os.Args[0])
}
