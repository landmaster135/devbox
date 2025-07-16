package config

import (
	"fmt"
	"os"
)

// 見出し定数
const (
	HeaderBase64Results = "=== Base64 Extraction Results ==="
)

// Config はCLI設定を保持する構造体
type Config struct {
	Path         string // 対象パス（必須）
	Recursive    bool   // 再帰検索（オプション、デフォルト: true）
	OutputFormat string // 出力形式（オプション、text/json、デフォルト: text）
}

// NewConfig は新しいConfigを作成する
func NewConfig(path string, recursive bool, outputFormat string) (*Config, error) {
	c := &Config{
		Path:         path,
		Recursive:    recursive,
		OutputFormat: outputFormat,
	}

	var err error
	c, err = validateConfig(c)
	if err != nil {
		return nil, fmt.Errorf("設定の初期化に失敗しました: %v", err)
	}
	return c, nil
}

// validateConfig は設定の妥当性を検証する
func validateConfig(config *Config) (*Config, error) {
	// Pathは必須
	if config.Path == "" {
		return nil, fmt.Errorf("--path は必須パラメータです")
	}

	// パスの存在確認
	if _, err := os.Stat(config.Path); err != nil {
		return nil, fmt.Errorf("指定されたパスが存在しません: %s", config.Path)
	}

	// 出力形式の検証
	if config.OutputFormat != "text" && config.OutputFormat != "json" {
		return nil, fmt.Errorf("--output-format は 'text' または 'json' を指定してください")
	}

	return config, nil
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
	var config Config

	cp.flagParser.StringVar(&config.Path, "path", "", "対象パス（必須）")
	cp.flagParser.BoolVar(&config.Recursive, "recursive", true, "再帰検索（デフォルト: true）")
	cp.flagParser.StringVar(&config.OutputFormat, "output-format", "text", "出力形式（text/json、デフォルト: text）")

	if err := cp.flagParser.Parse(); err != nil {
		return nil, err
	}

	return cp.validateConfig(&config)
}

// validateConfig は設定の妥当性を検証する
func (cp *ConfigParser) validateConfig(config *Config) (*Config, error) {
	config, err := validateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// ParseFlags は後方互換性のための関数
func ParseFlags() (*Config, error) {
	flagParser := NewRealFlagParser()
	osArgs := NewStandardOSArgs()
	configParser := NewConfigParser(flagParser, osArgs)
	return configParser.ParseFlags()
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	osArgs := NewStandardOSArgs()

	fmt.Printf("使用方法: %s [オプション]\n", osArgs.Args()[0])
	fmt.Printf("\nオプション:\n")
	fmt.Printf("  -path string\n")
	fmt.Printf("        対象パス（必須）\n")
	fmt.Printf("  -recursive\n")
	fmt.Printf("        再帰検索（デフォルト: true）\n")
	fmt.Printf("  -output-format string\n")
	fmt.Printf("        出力形式（text/json、デフォルト: text）\n")
	fmt.Printf("\n使用例:\n")
	fmt.Printf("  %s -path=/path/to/image.jpg\n", osArgs.Args()[0])
	fmt.Printf("  %s -path=/path/to/directory\n", osArgs.Args()[0])
	fmt.Printf("  %s -path=/path/to/directory -recursive=false\n", osArgs.Args()[0])
	fmt.Printf("  %s -path=/path/to/directory -output-format=json\n", osArgs.Args()[0])
}
