package config

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

const (
	// OperationListDeployments は Deployment Manager のデプロイメント一覧取得を表す。
	OperationListDeployments = "list-deployments"
)

// Config は CLI の入力値を保持する。
// operation ごとに利用するフィールドが異なる。
type Config struct {
	Operation       string
	ListDeployments ListDeploymentsConfig
}

// ListDeploymentsConfig は list-deployments 操作向けの設定値。
type ListDeploymentsConfig struct {
	Project     string
	Filter      string
	Format      string
	Limit       string
	Simple      bool
	ShowCommand bool
}

// ParseArgs は CLI 引数を解析し Config を返す。
func ParseArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet("gcloud-genset-deployment", flag.ContinueOnError)
	// flag パッケージのエラーメッセージを抑制する
	fs.SetOutput(new(strings.Builder))

	operation := fs.String("operation", "", "実行するオペレーション (必須)")
	project := fs.String("project", "", "GCP プロジェクト ID")
	filter := fs.String("filter", "", "結果をフィルタリングする式")
	format := fs.String("format", "table(name,insertTime,operation.operationType,operation.status,description)", "出力フォーマット")
	limit := fs.String("limit", "", "取得件数の上限")
	simple := fs.Bool("simple", false, "シンプルな出力フォーマットを使用する")
	showCommand := fs.Bool("show-command", false, "実行コマンドを事前に表示する")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if strings.TrimSpace(*operation) == "" {
		return nil, errors.New("--operation は必須です")
	}

	if !isSupportedOperation(*operation) {
		return nil, fmt.Errorf("未対応の operation です: %s", *operation)
	}

	cfg := &Config{
		Operation: *operation,
		ListDeployments: ListDeploymentsConfig{
			Project:     *project,
			Filter:      *filter,
			Format:      *format,
			Limit:       *limit,
			Simple:      *simple,
			ShowCommand: *showCommand,
		},
	}

	if *simple {
		cfg.ListDeployments.Format = "table(name,insertTime)"
	}

	return cfg, nil
}

// Usage は CLI の使用方法を返す。
func Usage() string {
	builder := &strings.Builder{}
	builder.WriteString("使用方法:\n")
	builder.WriteString("  gcloud-genset-deployment --operation=<operation> [オプション]\n\n")
	builder.WriteString("必須:\n")
	builder.WriteString("  --operation=list-deployments  デプロイメント一覧を取得\n\n")
	builder.WriteString("list-deployments オプション:\n")
	builder.WriteString("  --project=<PROJECT_ID>        対象プロジェクトを指定\n")
	builder.WriteString("  --filter=<FILTER>             フィルタ式を指定\n")
	builder.WriteString("  --format=<FORMAT>             出力フォーマットを指定\n")
	builder.WriteString("  --limit=<NUM>                 取得件数の上限を指定\n")
	builder.WriteString("  --simple                      シンプルなフォーマットを使用\n")
	builder.WriteString("  --show-command                実行前にコマンドを表示\n")
	return builder.String()
}

func isSupportedOperation(operation string) bool {
	switch operation {
	case OperationListDeployments:
		return true
	default:
		return false
	}
}
