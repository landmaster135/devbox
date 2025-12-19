package config

import (
	"fmt"
	"strings"
)

// Config はCLI設定を保持する構造体
type Config struct {
	RootDir    string   // ルートディレクトリ（必須）
	TargetDirs []string // 対象ディレクトリ（必須、カンマ区切り）
	Operation  string   // 実行する操作（必須）
}

var allowedOperations = []string{"output"}

// NewConfig は新しいConfigを作成する
func NewConfig(rootDir, targetDirs, operation string) (*Config, error) {
	c := &Config{
		RootDir:    rootDir,
		TargetDirs: parseTargetDirs(targetDirs),
		Operation:  operation,
	}

	var err error
	c, err = validateConfig(c)
	if err != nil {
		return nil, fmt.Errorf("設定の初期化に失敗しました: %v", err)
	}
	return c, nil
}

// parseTargetDirs はカンマ区切りの文字列をスライスに変換する
func parseTargetDirs(targetDirs string) []string {
	if targetDirs == "" {
		return []string{}
	}

	dirs := strings.Split(targetDirs, ",")
	result := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		trimmed := strings.TrimSpace(dir)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// validateConfig は設定の妥当性を検証する
func validateConfig(config *Config) (*Config, error) {
	// RootDirは必須
	if config.RootDir == "" {
		return nil, fmt.Errorf("--root-dir は必須パラメータです")
	}

	// TargetDirsは必須
	if len(config.TargetDirs) == 0 {
		return nil, fmt.Errorf("--target-dirs は必須パラメータです")
	}

	// Operationは必須
	if config.Operation == "" {
		return nil, fmt.Errorf("--operation は必須です")
	}

	if !isValidOperation(config.Operation) {
		return nil, fmt.Errorf("--operation は次のいずれかを指定してください: %s", strings.Join(allowedOperations, ", "))
	}

	return config, nil
}

func isValidOperation(operation string) bool {
	for _, allowed := range allowedOperations {
		if allowed == operation {
			return true
		}
	}
	return false
}

// ConfigParser はConfig解析を行う構造体
type ConfigParser struct {
	flagParser FlagParser
	osArgs     OSArgs
}

// NewConfigParser は新しいConfigParserを作成する
func NewConfigParser(flagParser FlagParser, osArgs OSArgs) *ConfigParser {
	return &ConfigParser{
		flagParser: flagParser,
		osArgs:     osArgs,
	}
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func (cp *ConfigParser) ParseFlags() (*Config, error) {
	var rootDir, targetDirs, operation string

	cp.flagParser.StringVar(&rootDir, "root-dir", "", "ルートディレクトリ（必須）")
	cp.flagParser.StringVar(&targetDirs, "target-dirs", "", "対象ディレクトリ（必須、カンマ区切り）")
	cp.flagParser.StringVar(&operation, "operation", "", fmt.Sprintf("実行する操作（必須: %s）", strings.Join(allowedOperations, ", ")))

	if err := cp.flagParser.Parse(); err != nil {
		return nil, err
	}

	return NewConfig(rootDir, targetDirs, operation)
}

// ParseFlags は後方互換性のための関数
func ParseFlags() (*Config, error) {
	flagParser := NewStandardFlagParser()
	osArgs := NewStandardOSArgs()
	configParser := NewConfigParser(flagParser, osArgs)
	return configParser.ParseFlags()
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	osArgs := NewStandardOSArgs()

	fmt.Printf("使用方法: %s [オプション]\n", osArgs.Args()[0])
	fmt.Printf("\nオプション:\n")
	fmt.Printf("  -root-dir string\n")
	fmt.Printf("        ルートディレクトリ（必須）\n")
	fmt.Printf("  -target-dirs string\n")
	fmt.Printf("        対象ディレクトリ（必須、カンマ区切り）\n")
	fmt.Printf("  -operation string\n")
	fmt.Printf("        実行する操作（必須: %s）\n", strings.Join(allowedOperations, ", "))
	fmt.Printf("\n使用例:\n")
	fmt.Printf("  %s -root-dir=$HOME/devbox -target-dirs=cli,mcp -operation=output\n", osArgs.Args()[0])
	fmt.Printf("  %s -root-dir=/path/to/project -target-dirs=\"cli,mcp,powershell\" -operation=output\n", osArgs.Args()[0])
}
