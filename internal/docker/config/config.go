package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	OperationEnvIntoCompose     = "env-into-compose"
	OperationPortsIntoCompose   = "ports-into-compose"
	OperationVolumesIntoCompose = "volumes-into-compose"
	OperationUserIntoCompose    = "user-into-compose"
)

// Config はCLIフラグの値を保持する
type Config struct {
	Operation   string
	ComposePath string
	EnvYAMLPath string
	PortKey     string
	VolumeKey   string
	UserKey     string
	Service     string
}

// ParseFlags はCLIフラグを解析してConfigを返す
func ParseFlags() (*Config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var cfg Config
	fs.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作（必須: %s）", OperationEnvIntoCompose))
	fs.StringVar(&cfg.ComposePath, "compose-path", "docker-compose.yml", "環境変数を書き込む docker-compose.yml へのパス")
	fs.StringVar(&cfg.EnvYAMLPath, "env-yaml-path", "env.yml", "環境変数を読み出す YAML ファイルのパス")
	fs.StringVar(&cfg.PortKey, "port-key", "", "operation=ports-into-compose のとき参照する env.yml のキー")
	fs.StringVar(&cfg.VolumeKey, "volume-key", "", "operation=volumes-into-compose のとき参照する env.yml のキー")
	fs.StringVar(&cfg.UserKey, "user-key", "", "operation=user-into-compose のとき参照する env.yml のキー")
	fs.StringVar(&cfg.Service, "service", "", "operation=ports-into-compose / volumes-into-compose / user-into-compose のとき更新対象となるサービス名")

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
	switch cfg.Operation {
	case OperationEnvIntoCompose:
		// 追加の必須パラメータなし
	case OperationPortsIntoCompose:
		if cfg.PortKey == "" {
			return fmt.Errorf("--port-key は operation=%s の場合に必須です", OperationPortsIntoCompose)
		}
		if cfg.Service == "" {
			return fmt.Errorf("--service は operation=%s の場合に必須です", OperationPortsIntoCompose)
		}
	case OperationVolumesIntoCompose:
		if cfg.VolumeKey == "" {
			return fmt.Errorf("--volume-key は operation=%s の場合に必須です", OperationVolumesIntoCompose)
		}
		if cfg.Service == "" {
			return fmt.Errorf("--service は operation=%s の場合に必須です", OperationVolumesIntoCompose)
		}
	case OperationUserIntoCompose:
		if cfg.UserKey == "" {
			return fmt.Errorf("--user-key は operation=%s の場合に必須です", OperationUserIntoCompose)
		}
		if cfg.Service == "" {
			return fmt.Errorf("--service は operation=%s の場合に必須です", OperationUserIntoCompose)
		}
	default:
		return fmt.Errorf("--operation は %s, %s, %s, %s のいずれかを指定してください", OperationEnvIntoCompose, OperationPortsIntoCompose, OperationVolumesIntoCompose, OperationUserIntoCompose)
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
	fmt.Fprintf(os.Stderr, "使用方法: %s --operation=<%s|%s|%s|%s> [オプション]\n", executable, OperationEnvIntoCompose, OperationPortsIntoCompose, OperationVolumesIntoCompose, OperationUserIntoCompose)
	fmt.Fprintln(os.Stderr, "\nオプション:")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作（%s, %s, %s, %s）\n", OperationEnvIntoCompose, OperationPortsIntoCompose, OperationVolumesIntoCompose, OperationUserIntoCompose)
	fmt.Fprintln(os.Stderr, "  -compose-path string\n        docker-compose.yml のパス（デフォルト: docker-compose.yml）")
	fmt.Fprintln(os.Stderr, "  -env-yaml-path string\n        env.yml のパス（デフォルト: env.yml）")
	fmt.Fprintln(os.Stderr, "  -port-key string\n        operation=ports-into-compose のとき参照する env.yml のキー")
	fmt.Fprintln(os.Stderr, "  -volume-key string\n        operation=volumes-into-compose のとき参照する env.yml のキー")
	fmt.Fprintln(os.Stderr, "  -user-key string\n        operation=user-into-compose のとき参照する env.yml のキー")
	fmt.Fprintln(os.Stderr, "  -service string\n        operation=ports-into-compose / volumes-into-compose / user-into-compose のとき更新対象となるサービス名")
}
