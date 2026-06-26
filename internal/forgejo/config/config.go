package config

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	envLoader "github.com/landmaster135/devbox/internal/env_loader"
	flagParser "github.com/landmaster135/devbox/internal/forgejo/infrastructures/flag_parser"
)

const (
	operationRepoList  = "repo list"
	operationIssueList = "issue list"
	envKeyHost         = "FORGEJO_HOST"
	envKeyUsername     = "FORGEJO_USERNAME"
	envKeyToken        = "FORGEJO_TOKEN"
	defaultWorkers     = 4
)

// Config はForgejo CLIの設定です。
type Config struct {
	Operation string // 実行する操作 (repo list, issue list)
	Host      string // Forgejoホスト
	Username  string // Forgejoユーザー名
	Token     string // APIトークン
	JSON      bool   // JSON形式で出力するか
	// ReposWorkers はrepo list取得時の同時実行ワーカー数です。
	ReposWorkers int
	Help         bool // ヘルプ表示フラグ
}

var supportedOperations = []string{operationRepoList, operationIssueList}

// NewConfig は新しいConfigを作成します。
func NewConfig(operation, host, username, token string, jsonOutput bool) (*Config, error) {
	operation = strings.TrimSpace(operation)
	if operation == "" {
		return nil, fmt.Errorf("operationが指定されていません")
	}

	if !slices.Contains(supportedOperations, operation) {
		return nil, fmt.Errorf("未対応のoperationです: %s", operation)
	}

	host = strings.TrimSpace(host)
	if host == "" {
		return nil, fmt.Errorf("%s が指定されていません", envKeyHost)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, fmt.Errorf("%s が指定されていません", envKeyUsername)
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, fmt.Errorf("%s が指定されていません", envKeyToken)
	}

	return &Config{
		Operation:    operation,
		Host:         host,
		Username:     username,
		Token:        token,
		JSON:         jsonOutput,
		ReposWorkers: defaultWorkers,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成します。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析します。
func ParseFlagsWithParser(parser flagParser.FlagParser) (*Config, error) {
	var (
		operation    = ""
		reposWorkers = ""
		jsonOutput   = false
		help         = false
	)

	parser.StringVar(&operation, "operation", operation, "実行する操作 (repo list, issue list)")
	parser.BoolVar(&jsonOutput, "json", jsonOutput, "JSON形式で出力")
	parser.StringVar(&reposWorkers, "repos-workers", reposWorkers, "repo list 取得時の同時実行数")
	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプを表示（短縮）")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグ解析に失敗しました: %w", err)
	}
	if help {
		return &Config{Help: true}, nil
	}

	args := parser.Args()
	operation = resolveOperation(operation, args)

	envValues, err := envLoader.Load([]string{envKeyHost, envKeyUsername, envKeyToken})
	if err != nil {
		return nil, err
	}

	reposWorkersInt, err := parseReposWorkers(reposWorkers)
	if err != nil {
		return nil, err
	}

	cfg, err := NewConfig(operation, envValues[envKeyHost], envValues[envKeyUsername], envValues[envKeyToken], jsonOutput)
	if err != nil {
		return nil, err
	}
	cfg.ReposWorkers = reposWorkersInt
	return cfg, nil
}

func parseReposWorkers(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultWorkers, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("repos-workers が不正です: %v", err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("repos-workers は1以上を指定してください")
	}
	return value, nil
}

func resolveOperation(operation string, args []string) string {
	operation = strings.TrimSpace(operation)
	if operation != "" {
		return operation
	}
	if len(args) >= 2 && strings.EqualFold(strings.TrimSpace(args[0]), "repo") && strings.EqualFold(strings.TrimSpace(args[1]), "list") {
		return operationRepoList
	}
	if len(args) >= 2 && strings.EqualFold(strings.TrimSpace(args[0]), "issue") && strings.EqualFold(strings.TrimSpace(args[1]), "list") {
		return operationIssueList
	}
	if len(args) >= 1 {
		return strings.TrimSpace(args[0])
	}
	return ""
}
