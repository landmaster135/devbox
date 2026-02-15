package usecases

import (
	"fmt"
	"strings"
)

// Service は Cloud Spanner 向けの gcloud コマンドを生成する。
type Service struct{}

// NewService は Service を生成する。
func NewService() *Service {
	return &Service{}
}

// InstanceCreateParams はインスタンス作成に必要な入力値。
type InstanceCreateParams struct {
	InstanceID     string
	InstanceConfig string
	Description    string
	Nodes          int
}

// DatabaseCreateParams はデータベース作成に必要な入力値。
type DatabaseCreateParams struct {
	InstanceID  string
	DatabaseID  string
	DDLFilePath string
}

// DatabaseListParams はデータベース一覧表示に必要な入力値。
type DatabaseListParams struct {
	InstanceID string
}

// DatabaseDescribeParams はデータベース詳細表示に必要な入力値。
type DatabaseDescribeParams struct {
	InstanceID string
	DatabaseID string
}

// BuildInstanceListCommand はインスタンス一覧コマンドを生成する。
func (s *Service) BuildInstanceListCommand() (string, error) {
	return "gcloud spanner instances list", nil
}

// BuildInstanceCreateCommand はインスタンス作成コマンドを生成する。
func (s *Service) BuildInstanceCreateCommand(params InstanceCreateParams) (string, error) {
	instanceID, err := requireValue(params.InstanceID, "instance-id")
	if err != nil {
		return "", err
	}
	config, err := requireValue(params.InstanceConfig, "config")
	if err != nil {
		return "", err
	}
	description, err := requireValue(params.Description, "description")
	if err != nil {
		return "", err
	}
	if params.Nodes <= 0 {
		return "", fmt.Errorf("nodes は1以上で指定してください")
	}

	command := fmt.Sprintf("gcloud spanner instances create %s \\\n    --config=%s \\\n    --description=%s \\\n    --nodes=%d",
		shellQuote(instanceID),
		shellQuote(config),
		shellQuote(description),
		params.Nodes,
	)

	return command, nil
}

// BuildDatabaseCreateCommand はデータベース作成コマンドを生成する。
func (s *Service) BuildDatabaseCreateCommand(params DatabaseCreateParams) (string, error) {
	instanceID, err := requireValue(params.InstanceID, "instance-id")
	if err != nil {
		return "", err
	}
	databaseID, err := requireValue(params.DatabaseID, "db-id")
	if err != nil {
		return "", err
	}
	ddlPath, err := requireValue(params.DDLFilePath, "ddl-file-path")
	if err != nil {
		return "", err
	}

	command := fmt.Sprintf("gcloud spanner databases create %s \\\n    --instance=%s \\\n    --ddl-file=%s",
		shellQuote(databaseID),
		shellQuote(instanceID),
		shellQuote(ddlPath),
	)

	return command, nil
}

// BuildDatabaseListCommand はデータベース一覧コマンドを生成する。
func (s *Service) BuildDatabaseListCommand(params DatabaseListParams) (string, error) {
	instanceID, err := requireValue(params.InstanceID, "instance-id")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("gcloud spanner databases list --instance=%s", shellQuote(instanceID)), nil
}

// BuildDatabaseDescribeCommand はデータベース詳細コマンドを生成する。
func (s *Service) BuildDatabaseDescribeCommand(params DatabaseDescribeParams) (string, error) {
	instanceID, err := requireValue(params.InstanceID, "instance-id")
	if err != nil {
		return "", err
	}
	databaseID, err := requireValue(params.DatabaseID, "db-id")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("gcloud spanner databases describe %s --instance=%s", shellQuote(databaseID), shellQuote(instanceID)), nil
}

// PrintHighlightedCommand は生成したコマンドを読みやすく出力する。
func (s *Service) PrintHighlightedCommand(command string) {
	fmt.Println()
	fmt.Println("==============================")
	fmt.Println("生成された gcloud コマンド")
	fmt.Println("==============================")
	fmt.Println(command)
	fmt.Println("==============================")
}

func requireValue(value string, field string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", fmt.Errorf("%s は必須です", field)
	}
	return trimmed, nil
}

func shellQuote(value string) string {
	escaped := strings.ReplaceAll(value, "'", "'\"'\"'")
	return "'" + escaped + "'"
}
