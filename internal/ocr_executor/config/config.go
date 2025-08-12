package config

import (
	"fmt"
	"os"
	"strings"
)

// 見出し定数
const (
	HeaderOCRResults = "=== OCR Results ==="
)

// Config はCLI設定を保持する構造体
type Config struct {
	Path         string // 対象パス（必須）
	Recursive    bool   // 再帰検索（オプション、デフォルト: false）
	Language     string // OCR言語（カンマ区切り、デフォルト: "jpn,eng"）
	OutputFormat string // 出力形式（オプション、text/json、デフォルト: text）
	OutputDir    string // 出力ディレクトリ（オプション）
}

// NewConfig は新しいConfigを作成する
func NewConfig(path string, recursive bool, language string, outputFormat string, outputDir string) (*Config, error) {
	c := &Config{
		Path:         path,
		Recursive:    recursive,
		Language:     language,
		OutputFormat: outputFormat,
		OutputDir:    outputDir,
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

	// 出力ディレクトリの検証（指定されている場合）
	if config.OutputDir != "" {
		if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
			return nil, fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
		}
	}

	// 言語設定の検証
	if err := validateLanguages(config.Language); err != nil {
		return nil, fmt.Errorf("言語設定が無効です: %v", err)
	}

	return config, nil
}

// validateLanguages は言語設定の妥当性を検証する
func validateLanguages(languages string) error {
	if languages == "" {
		return fmt.Errorf("言語設定が空です")
	}

	// カンマ区切りで分割
	langList := strings.Split(languages, ",")
	for _, lang := range langList {
		lang = strings.TrimSpace(lang)
		if lang == "" {
			return fmt.Errorf("空の言語コードが含まれています")
		}
		// 基本的な言語コードの検証（3文字以下）
		if len(lang) > 3 {
			return fmt.Errorf("無効な言語コード: %s", lang)
		}
	}

	return nil
}

// GetTesseractLanguages はTesseract用の言語設定を返す
func (c *Config) GetTesseractLanguages() string {
	// カンマ区切りをプラス区切りに変換
	langList := strings.Split(c.Language, ",")
	for i, lang := range langList {
		langList[i] = strings.TrimSpace(lang)
	}
	return strings.Join(langList, "+")
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
	cp.flagParser.BoolVar(&config.Recursive, "recursive", false, "再帰検索（デフォルト: false）")
	cp.flagParser.StringVar(&config.Language, "language", "jpn,eng", "OCR言語（カンマ区切り、デフォルト: jpn,eng）")
	cp.flagParser.StringVar(&config.OutputFormat, "output-format", "text", "出力形式（text/json、デフォルト: text）")
	cp.flagParser.StringVar(&config.OutputDir, "output-dir", "", "出力ディレクトリ（オプション）")

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
	fmt.Printf("        再帰検索（デフォルト: false）\n")
	fmt.Printf("  -language string\n")
	fmt.Printf("        OCR言語（カンマ区切り、デフォルト: jpn,eng）\n")
	fmt.Printf("  -output-format string\n")
	fmt.Printf("        出力形式（text/json、デフォルト: text）\n")
	fmt.Printf("  -output-dir string\n")
	fmt.Printf("        出力ディレクトリ（オプション、未指定時は標準出力のみ）\n")
	fmt.Printf("\n使用例:\n")
	fmt.Printf("  %s -path=/path/to/image.jpg\n", osArgs.Args()[0])
	fmt.Printf("  %s -path=/path/to/directory\n", osArgs.Args()[0])
	fmt.Printf("  %s -path=/path/to/directory -recursive=true\n", osArgs.Args()[0])
	fmt.Printf("  %s -path=/path/to/directory -language=jpn,eng,fra\n", osArgs.Args()[0])
	fmt.Printf("  %s -path=/path/to/directory -output-format=json -output-dir=/path/to/output\n", osArgs.Args()[0])
}
