package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	dateLayout            = "20060102"
	defaultLimit          = 200
	retrieveSessionOpName = "retrieve-session"
	codexAgentType        = "codex"
)

// Config はCLI設定を保持する構造体。
type Config struct {
	Operation    string
	AgentType    string
	Limit        int
	StartDate    string
	EndDate      string
	AgentHomeDir string

	StartDateValue *time.Time
	EndDateValue   *time.Time
}

// NewConfig は新しいConfigを作成する。
func NewConfig(operation, agentType string, limit int, startDate, endDate, agentHomeDir string) (*Config, error) {
	cfg := &Config{
		Operation:    operation,
		AgentType:    agentType,
		Limit:        limit,
		StartDate:    startDate,
		EndDate:      endDate,
		AgentHomeDir: agentHomeDir,
	}

	validatedConfig, err := validateConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("設定の初期化に失敗しました: %w", err)
	}

	return validatedConfig, nil
}

// ConfigParser はConfig解析を行う構造体。
type ConfigParser struct {
	flagParser FlagParser
	osArgs     OSArgs
}

// NewConfigParser は新しいConfigParserを作成する。
func NewConfigParser(flagParser FlagParser, osArgs OSArgs) *ConfigParser {
	return &ConfigParser{flagParser: flagParser, osArgs: osArgs}
}

// ParseFlags はコマンドライン引数を解析してConfigを返す。
func (cp *ConfigParser) ParseFlags() (*Config, error) {
	defaultAgentHomeDir := resolveDefaultAgentHomeDir()

	var operation string
	var agentType string
	limit := defaultLimit
	var startDate string
	var endDate string
	agentHomeDir := defaultAgentHomeDir

	cp.flagParser.StringVar(&operation, "operation", "", "実行する操作（必須: retrieve-session）")
	cp.flagParser.StringVar(&agentType, "agent-type", "", "エージェントタイプ（必須: codex）")
	cp.flagParser.IntVar(&limit, "limit", defaultLimit, "取得件数の上限（任意。デフォルト: 200）")
	cp.flagParser.StringVar(&startDate, "start-date", "", "開始日（任意: yyyyMMdd）")
	cp.flagParser.StringVar(&endDate, "end-date", "", "終了日（任意: yyyyMMdd）")
	cp.flagParser.StringVar(&agentHomeDir, "agent-home-dir", defaultAgentHomeDir, "エージェントホームディレクトリ（任意。デフォルト: ~/.codex）")

	if err := cp.flagParser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	return NewConfig(operation, agentType, limit, startDate, endDate, agentHomeDir)
}

// ParseFlags は後方互換性のための関数。
func ParseFlags() (*Config, error) {
	parser := NewStandardFlagParser()
	osArgs := NewStandardOSArgs()
	configParser := NewConfigParser(parser, osArgs)
	return configParser.ParseFlags()
}

// PrintUsage は使用方法を表示する。
func PrintUsage() {
	osArgs := NewStandardOSArgs()
	program := "agent-session"
	if len(osArgs.Args()) > 0 {
		program = filepath.Base(osArgs.Args()[0])
	}

	fmt.Printf("使用方法: %s [オプション]\n", program)
	fmt.Printf("\nオプション:\n")
	fmt.Printf("  -operation string\n")
	fmt.Printf("        実行する操作（必須: retrieve-session）\n")
	fmt.Printf("  -agent-type string\n")
	fmt.Printf("        エージェントタイプ（必須: codex）\n")
	fmt.Printf("  -limit int\n")
	fmt.Printf("        取得件数の上限（任意、デフォルト: %d）\n", defaultLimit)
	fmt.Printf("  -start-date string\n")
	fmt.Printf("        開始日（任意: yyyyMMdd）\n")
	fmt.Printf("  -end-date string\n")
	fmt.Printf("        終了日（任意: yyyyMMdd）\n")
	fmt.Printf("  -agent-home-dir string\n")
	fmt.Printf("        エージェントホームディレクトリ（任意、デフォルト: ~/.codex）\n")

	fmt.Printf("\n使用例:\n")
	fmt.Printf("  %s -operation=retrieve-session -agent-type=codex\n", program)
	fmt.Printf("  %s -operation=retrieve-session -agent-type=codex -limit=50 -start-date=20260301 -end-date=20260331\n", program)
	fmt.Printf("  %s -operation=retrieve-session -agent-type=codex -agent-home-dir=$HOME/.codex\n", program)
}

func validateConfig(config *Config) (*Config, error) {
	if strings.TrimSpace(config.Operation) == "" {
		return nil, fmt.Errorf("--operation は必須です")
	}
	if config.Operation != retrieveSessionOpName {
		return nil, fmt.Errorf("--operation は %s のみ対応しています", retrieveSessionOpName)
	}
	if strings.TrimSpace(config.AgentType) == "" {
		return nil, fmt.Errorf("--agent-type は必須です")
	}
	if config.AgentType != codexAgentType {
		return nil, fmt.Errorf("--agent-type は %s のみ対応しています", codexAgentType)
	}
	if config.Limit <= 0 {
		return nil, fmt.Errorf("--limit は1以上を指定してください")
	}
	if strings.TrimSpace(config.AgentHomeDir) == "" {
		return nil, fmt.Errorf("--agent-home-dir は空にできません")
	}

	startDateValue, err := parseDate(config.StartDate)
	if err != nil {
		return nil, fmt.Errorf("--start-date の形式が不正です: %w", err)
	}
	endDateValue, err := parseDate(config.EndDate)
	if err != nil {
		return nil, fmt.Errorf("--end-date の形式が不正です: %w", err)
	}

	if startDateValue != nil && endDateValue != nil && startDateValue.After(*endDateValue) {
		return nil, fmt.Errorf("--start-date は --end-date 以下を指定してください")
	}

	config.StartDateValue = startDateValue
	config.EndDateValue = endDateValue

	return config, nil
}

func parseDate(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	parsed, err := time.ParseInLocation(dateLayout, trimmed, time.Local)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func resolveDefaultAgentHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ".codex"
	}
	return filepath.Join(home, ".codex")
}
