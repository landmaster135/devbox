package config

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
)

// Config はCLI設定を保持する構造体
type Config struct {
	Target    string // 対象パス（ファイルまたはディレクトリ）（必須）
	From      string // 置換元文字列（必須）
	To        string // 置換先文字列（必須）
	Encoding  string // 文字エンコーディング（オプション、デフォルト: utf-8）
	Recursive bool   // 再帰的処理（オプション）
	Backup    bool   // バックアップ作成（オプション）
	BackupDir string // バックアップディレクトリ（オプション）
	DryRun    bool   // ドライラン（オプション）
}

// NewConfig は新しいConfigを作成します
func NewConfig(target, from, to, encoding, backupDir string, recursive, backup, dryRun bool) (*Config, error) {
	c := &Config{
		Target:    target,
		From:      from,
		To:        to,
		Encoding:  encoding,
		Recursive: recursive,
		Backup:    backup,
		BackupDir: backupDir,
		DryRun:    dryRun,
	}

	var err error
	c, err = validateConfig(c)
	if err != nil {
		return nil, fmt.Errorf("設定の初期化に失敗しました: %v", err)
	}
	return c, nil
}

// ToReplacementConfig はConfigをdomain.ReplacementConfigに変換します
func (c *Config) ToReplacementConfig() *domain.ReplacementConfig {
	encoding := domain.EncodingType(c.Encoding)
	if encoding == "" {
		encoding = domain.EncodingUTF8
	}

	return &domain.ReplacementConfig{
		Target:    c.Target,
		From:      c.From,
		To:        c.To,
		Encoding:  encoding,
		Recursive: c.Recursive,
		Backup:    c.Backup,
		BackupDir: c.BackupDir,
		DryRun:    c.DryRun,
	}
}

// validateEncoding はエンコーディングの妥当性を検証します
func validateEncoding(encoding string) error {
	if encoding == "" {
		return nil // 空の場合はデフォルト値が使用される
	}

	validEncodings := []string{
		string(domain.EncodingUTF8),
		string(domain.EncodingShiftJIS),
		string(domain.EncodingEUCJP),
		string(domain.EncodingISO2022JP),
	}

	for _, valid := range validEncodings {
		if encoding == valid {
			return nil
		}
	}

	return fmt.Errorf("サポートされていないエンコーディングです: %s", encoding)
}

// validateConfig は設定の妥当性を検証します
func validateConfig(config *Config) (*Config, error) {
	// Targetは必須
	if config.Target == "" {
		return nil, fmt.Errorf("--target は必須パラメータです")
	}

	// Fromは必須
	if config.From == "" {
		return nil, fmt.Errorf("--from は必須パラメータです")
	}

	// Toは必須
	if config.To == "" {
		return nil, fmt.Errorf("--to は必須パラメータです")
	}

	// エンコーディングの検証
	err := validateEncoding(config.Encoding)
	if err != nil {
		return nil, err
	}

	// デフォルト値の設定
	if config.Encoding == "" {
		config.Encoding = string(domain.EncodingUTF8)
	}

	return config, nil
}

// ConfigParser はConfig解析を行う構造体
type ConfigParser struct {
	flagParser FlagParser
	osArgs     OSArgs
}

// NewConfigParser は新しいConfigParserを作成します
func NewConfigParser(flagParser FlagParser, osArgs OSArgs) *ConfigParser {
	return &ConfigParser{
		flagParser: flagParser,
		osArgs:     osArgs,
	}
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func (cp *ConfigParser) ParseFlags() (*Config, error) {
	var config Config

	cp.flagParser.StringVar(&config.Target, "target", "", "対象パス（ファイルまたはディレクトリ）（必須）")
	cp.flagParser.StringVar(&config.From, "from", "", "置換元文字列（必須）")
	cp.flagParser.StringVar(&config.To, "to", "", "置換先文字列（必須）")
	cp.flagParser.StringVar(&config.Encoding, "encoding", "utf-8", "文字エンコーディング（utf-8, shift_jis, euc-jp, iso-2022-jp）")
	cp.flagParser.StringVar(&config.BackupDir, "backup-dir", "", "バックアップディレクトリ（ディレクトリ処理時は必須）")
	cp.flagParser.BoolVar(&config.Recursive, "recursive", false, "ディレクトリの場合、再帰的に処理")
	cp.flagParser.BoolVar(&config.Backup, "backup", false, "バックアップファイルを作成")
	cp.flagParser.BoolVar(&config.DryRun, "dry-run", false, "実際の変更を行わず、変更予定を表示のみ")

	if err := cp.flagParser.Parse(); err != nil {
		return nil, err
	}

	return cp.validateConfig(&config)
}

// validateConfig は設定の妥当性を検証します
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
	fmt.Printf("  -target string\n")
	fmt.Printf("        対象パス（ファイルまたはディレクトリ）（必須）\n")
	fmt.Printf("  -from string\n")
	fmt.Printf("        置換元文字列（必須）\n")
	fmt.Printf("  -to string\n")
	fmt.Printf("        置換先文字列（必須）\n")
	fmt.Printf("  -encoding string\n")
	fmt.Printf("        文字エンコーディング（デフォルト: utf-8）\n")
	fmt.Printf("        対応エンコーディング: utf-8, shift_jis, euc-jp, iso-2022-jp\n")
	fmt.Printf("  -recursive\n")
	fmt.Printf("        ディレクトリの場合、再帰的に処理\n")
	fmt.Printf("  -backup\n")
	fmt.Printf("        バックアップファイルを作成\n")
	fmt.Printf("  -backup-dir string\n")
	fmt.Printf("        バックアップディレクトリ（ディレクトリ処理時は必須）\n")
	fmt.Printf("  -dry-run\n")
	fmt.Printf("        実際の変更を行わず、変更予定を表示のみ\n")
	fmt.Printf("\n使用例:\n")
	fmt.Printf("  %s -target=./test.txt -from=\"old\" -to=\"new\"\n", osArgs.Args()[0])
	fmt.Printf("  %s -target=./test.txt -from=\"old\" -to=\"new\" -backup -backup-dir=./backups\n", osArgs.Args()[0])
	fmt.Printf("  %s -target=./src -from=\"TODO\" -to=\"DONE\" -recursive -backup -backup-dir=./backups\n", osArgs.Args()[0])
	fmt.Printf("  %s -target=./data.txt -from=\"古い\" -to=\"新しい\" -encoding=shift_jis\n", osArgs.Args()[0])
	fmt.Printf("  %s -target=./project -from=\"debug\" -to=\"release\" -recursive -dry-run\n", osArgs.Args()[0])
}
