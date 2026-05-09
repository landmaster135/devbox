package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/forgejo/infrastructures/flag_parser"
)

const (
	dotenvFilename       = ".env"
	operationRepoList    = "repo list"
	operationProjectList = "project list"
	envKeyHost           = "forgejo-host"
	envKeyUsername       = "forgejo-username"
	envKeyToken          = "forgejo-token"
	defaultWorkers       = 4
)

// Config はForgejo CLIの設定です。
type Config struct {
	Operation string // 実行する操作 (repo list, project list)
	Host      string // Forgejoホスト
	Username  string // Forgejoユーザー名
	Token     string // APIトークン
	JSON      bool   // JSON形式で出力するか
	// PullsPageWorkers はrepo list取得時の同時実行ワーカー数です。
	PullsPageWorkers int
	Help             bool // ヘルプ表示フラグ
}

var supportedOperations = []string{operationRepoList, operationProjectList}

// NewConfig は新しいConfigを作成します。
func NewConfig(operation, host, username, token string, jsonOutput bool) (*Config, error) {
	operation = strings.TrimSpace(operation)
	if operation == "" {
		return nil, fmt.Errorf("operationが指定されていません")
	}

	isSupported := false
	for _, op := range supportedOperations {
		if operation == op {
			isSupported = true
			break
		}
	}
	if !isSupported {
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
		Operation:        operation,
		Host:             host,
		Username:         username,
		Token:            token,
		JSON:             jsonOutput,
		PullsPageWorkers: defaultWorkers,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成します。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(flag_parser.NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析します。
func ParseFlagsWithParser(parser flag_parser.FlagParser) (*Config, error) {
	return ParseFlagsWithParserWithEnvFile(parser, dotenvFilename)
}

// ParseFlagsWithParserWithEnvFile は指定した .env ファイルを環境変数のデフォルト値として使用します。
func ParseFlagsWithParserWithEnvFile(parser flag_parser.FlagParser, envFilePath string) (*Config, error) {
	var (
		operation        = ""
		host             = ""
		username         = ""
		token            = ""
		pullsPageWorkers = ""
		jsonOutput       = false
		help             = false
	)

	parser.StringVar(&operation, "operation", operation, "実行する操作 (repo list, project list)")
	parser.BoolVar(&jsonOutput, "json", jsonOutput, "JSON形式で出力")
	parser.StringVar(&host, "forgejo-host", host, "Forgejoホスト（https://example.com）")
	parser.StringVar(&username, "forgejo-username", username, "Forgejoユーザー名")
	parser.StringVar(&token, "forgejo-token", token, "Forgejo APIトークン")
	parser.StringVar(&pullsPageWorkers, "forgejo-pulls-page-workers", pullsPageWorkers, "repo list 取得時の同時実行数")
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

	envValues, err := loadDotEnv(envFilePath)
	if err != nil {
		return nil, err
	}

	if host == "" {
		host = lookupEnvWithDotEnv(envKeyHost, host, envValues)
	}
	if username == "" {
		username = lookupEnvWithDotEnv(envKeyUsername, username, envValues)
	}
	if token == "" {
		token = lookupEnvWithDotEnv(envKeyToken, token, envValues)
	}

	pullsPageWorkersInt, err := parsePullsPageWorkers(pullsPageWorkers)
	if err != nil {
		return nil, err
	}

	cfg, err := NewConfig(operation, host, username, token, jsonOutput)
	if err != nil {
		return nil, err
	}
	cfg.PullsPageWorkers = pullsPageWorkersInt
	return cfg, nil
}

func parsePullsPageWorkers(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultWorkers, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("forgejo-pulls-page-workers が不正です: %v", err)
	}
	if value <= 0 {
		return 0, fmt.Errorf("forgejo-pulls-page-workers は1以上を指定してください")
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
	if len(args) >= 2 && strings.EqualFold(strings.TrimSpace(args[0]), "project") && strings.EqualFold(strings.TrimSpace(args[1]), "list") {
		return operationProjectList
	}
	if len(args) >= 1 {
		return strings.TrimSpace(args[0])
	}
	return ""
}

func lookupEnvWithDotEnv(key, flagValue string, envValues map[string]string) string {
	flagValue = strings.TrimSpace(flagValue)
	if flagValue != "" {
		return flagValue
	}
	if osValue := strings.TrimSpace(os.Getenv(key)); osValue != "" {
		return osValue
	}
	if fileValue, ok := envValues[key]; ok {
		return fileValue
	}
	return ""
}

func loadDotEnv(path string) (map[string]string, error) {
	if path == "" {
		return map[string]string{}, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}

	values := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(string(raw)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := splitKeyValue(line)
		if !ok {
			continue
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

func splitKeyValue(line string) (key string, value string, ok bool) {
	if !strings.Contains(line, "=") {
		return "", "", false
	}
	parts := strings.SplitN(line, "=", 2)
	key = strings.TrimSpace(parts[0])
	value = strings.TrimSpace(parts[1])
	if key == "" {
		return "", "", false
	}
	value = strings.TrimSpace(strings.Trim(value, "\""))
	return key, value, true
}
