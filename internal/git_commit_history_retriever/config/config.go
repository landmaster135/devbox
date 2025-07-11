package config

import (
	"fmt"
	"regexp"
)

// 見出し定数
const (
	HeaderCommitHistory = "=== Commit History ==="
	HeaderCommitDetails = "=== Commit Details ==="
)

// Config はCLI設定を保持する構造体
type Config struct {
	GitDir  string // 対象Gitディレクトリ（必須）
	Keyword string // 検索キーワード（オプション）
	Since   string // 開始年月日（オプション、YYYY-MM-DD形式）
	Until   string // 終了年月日（オプション、YYYY-MM-DD形式）
}

func NewConfig(gitDir, keyword, since, until string)(*Config, error){
	c := &Config{
		GitDir:  gitDir,
		Keyword: keyword,
		Since:   since,
		Until:   until,
	}
	var err error
	c, err = validateConfig(c); if err != nil{
		return nil, fmt.Errorf("設定の初期化に失敗しました: %v", err)
	}
	return c, nil
}

// validateDateFormat は日付フォーマットを検証する
func validateDateFormat(date string) error {
	if date == "" {
		return nil
	}
	datePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !datePattern.MatchString(date) {
		return fmt.Errorf("日付フォーマットが正しくありません。YYYY-MM-DD形式で入力してください: %s", date)
	}
	return nil
}

func validateConfig(config *Config) (*Config, error) {
	// GitDirは必須
	if config.GitDir == "" {
		return nil, fmt.Errorf("--git-dir は必須パラメータです")
	}

	err := validateDateFormat(config.Since); if err != nil {
		return nil, fmt.Errorf("--since の日付フォーマットが正しくありません。YYYY-MM-DD形式で入力してください")
	}

	err = validateDateFormat(config.Until); if err != nil {
		return nil, fmt.Errorf("--until の日付フォーマットが正しくありません。YYYY-MM-DD形式で入力してください")
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

	cp.flagParser.StringVar(&config.GitDir, "git-dir", "", "対象Gitディレクトリ（必須）")
	cp.flagParser.StringVar(&config.Keyword, "keyword", "", "検索キーワード（オプション）")
	cp.flagParser.StringVar(&config.Since, "since", "", "開始年月日（オプション、YYYY-MM-DD形式）")
	cp.flagParser.StringVar(&config.Until, "until", "", "終了年月日（オプション、YYYY-MM-DD形式）")

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
	fmt.Printf("  -git-dir string\n")
	fmt.Printf("        対象Gitディレクトリ（必須）\n")
	fmt.Printf("  -keyword string\n")
	fmt.Printf("        検索キーワード（オプション）\n")
	fmt.Printf("  -since string\n")
	fmt.Printf("        開始年月日（オプション、YYYY-MM-DD形式）\n")
	fmt.Printf("  -until string\n")
	fmt.Printf("        終了年月日（オプション、YYYY-MM-DD形式）\n")
	fmt.Printf("\n使用例:\n")
	fmt.Printf("  %s -git-dir=/path/to/repo\n", osArgs.Args()[0])
	fmt.Printf("  %s -git-dir=/path/to/repo -keyword=\"feat:\"\n", osArgs.Args()[0])
	fmt.Printf("  %s -git-dir=/path/to/repo -since=\"2025-01-01\" -until=\"2025-01-31\"\n", osArgs.Args()[0])
}
