package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	// OperationInstanceList は Spanner インスタンス一覧取得コマンドを表す。
	OperationInstanceList = "instance-list"
	// OperationInstanceCreate は Spanner インスタンス作成コマンドを表す。
	OperationInstanceCreate = "instance-create"
	// OperationDatabaseCreate は Spanner データベース新規作成コマンドを表す。
	OperationDatabaseCreate = "db-create"
	// OperationDatabaseList は Spanner データベース一覧取得コマンドを表す。
	OperationDatabaseList = "db-list"
	// OperationDatabaseDescribe は Spanner データベース詳細表示コマンドを表す。
	OperationDatabaseDescribe = "db-describe"
)

var supportedOperations = []string{
	OperationDatabaseCreate,
	OperationDatabaseDescribe,
	OperationDatabaseList,
	OperationInstanceCreate,
	OperationInstanceList,
}

// Config は CLI で渡されるパラメータを保持する。
type Config struct {
	Operation      string
	Help           bool
	InstanceID     string
	InstanceConfig string
	Description    string
	Nodes          int
	DatabaseID     string
	DDLFilePath    string
}

// ParseFlags は標準パーサーで CLI 引数を解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は注入したパーサーで CLI 引数を解析する。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")
	parser.StringVar(&cfg.InstanceID, "instance-id", "", "Spanner インスタンス ID")
	parser.StringVar(&cfg.InstanceConfig, "config", "", "Spanner インスタンス構成 ID (例: regional-asia-northeast1)")
	parser.StringVar(&cfg.Description, "description", "", "インスタンスの説明")
	parser.IntVar(&cfg.Nodes, "nodes", 0, "作成するノード数")
	parser.StringVar(&cfg.DatabaseID, "db-id", "", "Spanner データベース ID")
	parser.StringVar(&cfg.DDLFilePath, "ddl-file-path", "", "DDL ファイルへのパス")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if len(parser.Args()) > 0 {
		return nil, fmt.Errorf("未処理の位置引数があります: %v", parser.Args())
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
	cfg.Operation = strings.ToLower(strings.TrimSpace(cfg.Operation))
	cfg.InstanceID = strings.TrimSpace(cfg.InstanceID)
	cfg.InstanceConfig = strings.TrimSpace(cfg.InstanceConfig)
	cfg.Description = strings.TrimSpace(cfg.Description)
	cfg.DatabaseID = strings.TrimSpace(cfg.DatabaseID)
	cfg.DDLFilePath = strings.TrimSpace(cfg.DDLFilePath)
}

func validateConfig(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}

	if !isSupportedOperation(cfg.Operation) {
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	switch cfg.Operation {
	case OperationInstanceList:
		// 追加パラメータなし
	case OperationInstanceCreate:
		if cfg.InstanceID == "" {
			return fmt.Errorf("instance-id パラメータは必須です")
		}
		if cfg.InstanceConfig == "" {
			return fmt.Errorf("config パラメータは必須です")
		}
		if cfg.Description == "" {
			return fmt.Errorf("description パラメータは必須です")
		}
		if cfg.Nodes <= 0 {
			return fmt.Errorf("nodes パラメータは1以上で指定してください")
		}
	case OperationDatabaseCreate:
		if cfg.InstanceID == "" {
			return fmt.Errorf("instance-id パラメータは必須です")
		}
		if cfg.DatabaseID == "" {
			return fmt.Errorf("db-id パラメータは必須です")
		}
		if cfg.DDLFilePath == "" {
			return fmt.Errorf("ddl-file-path パラメータは必須です")
		}
	case OperationDatabaseList:
		if cfg.InstanceID == "" {
			return fmt.Errorf("instance-id パラメータは必須です")
		}
	case OperationDatabaseDescribe:
		if cfg.InstanceID == "" {
			return fmt.Errorf("instance-id パラメータは必須です")
		}
		if cfg.DatabaseID == "" {
			return fmt.Errorf("db-id パラメータは必須です")
		}
	}

	return nil
}

func isSupportedOperation(operation string) bool {
	for _, op := range supportedOperations {
		if op == operation {
			return true
		}
	}
	return false
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Cloud Spanner 向け gcloud コマンド生成ツール\n\n")

	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(supportedOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -help\n        このヘルプを表示\n")
	fmt.Fprintf(os.Stderr, "  -instance-id string\n        対象の Spanner インスタンス ID (operation によって必須)\n")
	fmt.Fprintf(os.Stderr, "  -db-id string\n        対象の Spanner データベース ID (db-create / db-describe で必須)\n")
	fmt.Fprintf(os.Stderr, "  -ddl-file-path string\n        適用する DDL ファイルへのパス (db-create で必須)\n")
	fmt.Fprintf(os.Stderr, "  -config string\n        インスタンス構成 ID (instance-create で必須)\n")
	fmt.Fprintf(os.Stderr, "  -description string\n        インスタンスの説明 (instance-create で必須)\n")
	fmt.Fprintf(os.Stderr, "  -nodes int\n        作成するノード数 (instance-create で必須)\n\n")

	fmt.Fprintf(os.Stderr, "operation 別の使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=%s\n", os.Args[0], OperationInstanceList)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -instance-id=my-instance -config=regional-asia-northeast1 -description=\"Project API\" -nodes=1\n", os.Args[0], OperationInstanceCreate)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -instance-id=my-instance -db-id=orders -ddl-file-path=schema.sql\n", os.Args[0], OperationDatabaseCreate)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -instance-id=my-instance\n", os.Args[0], OperationDatabaseList)
	fmt.Fprintf(os.Stderr, "  %s -operation=%s -instance-id=my-instance -db-id=orders\n", os.Args[0], OperationDatabaseDescribe)
}

func init() {
	sort.Strings(supportedOperations)
}
