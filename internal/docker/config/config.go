package config

import (
	"flag"
	"fmt"
	"os"
)

const OperationEnvIntoCompose = "env-into-compose"

// Config はCLIフラグの値を保持する
type Config struct {
	Operation   string
	ComposePath string
	EnvYAMLPath string
}

// ParseFlags はCLIフラグを解析してConfigを返す
func ParseFlags() (*Config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var cfg Config
	fs.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作（必須: %s）", OperationEnvIntoCompose))
	fs.StringVar(&cfg.ComposePath, "compose-path", "docker-compose.yml", "環境変数を書き込む docker-compose.yml へのパス")
	fs.StringVar(&cfg.EnvYAMLPath, "env-yaml-path", "env.yml", "環境変数を読み出す YAML ファイルのパス")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	if err := validate(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("--operation は必須です")
	}
	if cfg.Operation != OperationEnvIntoCompose {
		return fmt.Errorf("--operation は %s のみサポートされています", OperationEnvIntoCompose)
	}
	if cfg.ComposePath == "" {
		return fmt.Errorf("--compose-path は必須です")
	}
	if cfg.EnvYAMLPath == "" {
		return fmt.Errorf("--env-yaml-path は必須です")
	}
	return nil
}

// PrintUsage はCLIの使用方法を表示する
func PrintUsage() {
	executable := os.Args[0]
	fmt.Fprintf(os.Stderr, "使用方法: %s --operation=%s [オプション]\n", executable, OperationEnvIntoCompose)
	fmt.Fprintln(os.Stderr, "\nオプション:")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作（必須: %s）\n", OperationEnvIntoCompose)
	fmt.Fprintln(os.Stderr, "  -compose-path string\n        docker-compose.yml のパス（デフォルト: docker-compose.yml）")
	fmt.Fprintln(os.Stderr, "  -env-yaml-path string\n        env.yml のパス（デフォルト: env.yml）")
}
