package config

import (
	"fmt"
	"strings"
)

const (
	// OperationCreateBadge はカバレッジバッジを生成するoperation。
	OperationCreateBadge = "create-badge"
	// OperationPatchBadge はMarkdownへカバレッジバッジを追加/更新するoperation。
	OperationPatchBadge = "patch-badge"
)

const (
	defaultBadgeTitle      = "Coverage"
	defaultCoverageFile    = "coverage.out"
	defaultGreenThreshold  = 70
	defaultYellowThreshold = 30
	defaultTargetFile      = "README.md"
)

var allowedOperations = []string{
	OperationCreateBadge,
	OperationPatchBadge,
}

var allowedColors = map[string]struct{}{
	"green":  {},
	"yellow": {},
	"red":    {},
}

// Config はcoverage-badge CLIの設定を保持する。
type Config struct {
	Operation       string
	BadgeTitle      string
	CoverageFile    string
	GreenThreshold  int
	YellowThreshold int
	ForceColor      string
	BadgeLink       string
	BadgeValue      string
	TargetFile      string
	DryRun          bool
	Help            bool
}

// NewConfig は新しいConfigを作成し、検証する。
func NewConfig(
	operation string,
	badgeTitle string,
	coverageFile string,
	greenThreshold int,
	yellowThreshold int,
	forceColor string,
	badgeLink string,
	badgeValue string,
	targetFile string,
	dryRun bool,
	help bool,
) (*Config, error) {
	cfg := &Config{
		Operation:       strings.TrimSpace(operation),
		BadgeTitle:      strings.TrimSpace(badgeTitle),
		CoverageFile:    strings.TrimSpace(coverageFile),
		GreenThreshold:  greenThreshold,
		YellowThreshold: yellowThreshold,
		ForceColor:      strings.ToLower(strings.TrimSpace(forceColor)),
		BadgeLink:       strings.TrimSpace(badgeLink),
		BadgeValue:      strings.TrimSpace(badgeValue),
		TargetFile:      strings.TrimSpace(targetFile),
		DryRun:          dryRun,
		Help:            help,
	}

	if cfg.Help {
		return cfg, nil
	}

	applyDefaults(cfg)

	normalized, err := validateConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("設定の初期化に失敗しました: %w", err)
	}

	return normalized, nil
}

func applyDefaults(cfg *Config) {
	if cfg.BadgeTitle == "" {
		cfg.BadgeTitle = defaultBadgeTitle
	}
	if cfg.CoverageFile == "" {
		cfg.CoverageFile = defaultCoverageFile
	}
	if cfg.TargetFile == "" {
		cfg.TargetFile = defaultTargetFile
	}
}

func validateConfig(cfg *Config) (*Config, error) {
	if cfg.Operation == "" {
		return nil, fmt.Errorf("--operation は必須です")
	}
	if !isValidOperation(cfg.Operation) {
		return nil, fmt.Errorf("--operation は次のいずれかを指定してください: %s", strings.Join(allowedOperations, ", "))
	}

	if cfg.BadgeTitle == "" {
		return nil, fmt.Errorf("--badge-title は空文字を指定できません")
	}
	if cfg.CoverageFile == "" && cfg.BadgeValue == "" {
		return nil, fmt.Errorf("--coverage-file または --badge-value のいずれかは必須です")
	}

	if cfg.GreenThreshold < 0 || cfg.GreenThreshold > 100 {
		return nil, fmt.Errorf("--green-threshold は 0 から 100 の範囲で指定してください")
	}
	if cfg.YellowThreshold < 0 || cfg.YellowThreshold > 100 {
		return nil, fmt.Errorf("--yellow-threshold は 0 から 100 の範囲で指定してください")
	}
	if cfg.GreenThreshold <= cfg.YellowThreshold {
		return nil, fmt.Errorf("--green-threshold は --yellow-threshold より大きい値を指定してください")
	}

	if cfg.ForceColor != "" {
		if _, ok := allowedColors[cfg.ForceColor]; !ok {
			return nil, fmt.Errorf("--force-color は green, yellow, red のいずれかを指定してください")
		}
	}

	if cfg.Operation == OperationPatchBadge && cfg.TargetFile == "" {
		return nil, fmt.Errorf("--target-file は operation=patch-badge の場合に必須です")
	}

	return cfg, nil
}

func isValidOperation(operation string) bool {
	for _, allowed := range allowedOperations {
		if operation == allowed {
			return true
		}
	}
	return false
}

// ConfigParser はConfigを生成する責務を持つ。
type ConfigParser struct {
	flagParser FlagParser
	osArgs     OSArgs
}

// NewConfigParser は新しいConfigParserを作成する。
func NewConfigParser(flagParser FlagParser, osArgs OSArgs) *ConfigParser {
	return &ConfigParser{
		flagParser: flagParser,
		osArgs:     osArgs,
	}
}

// ParseFlags はCLIフラグを解析してConfigを返す。
func (cp *ConfigParser) ParseFlags() (*Config, error) {
	var (
		operation       string
		badgeTitle      string
		coverageFile    string
		greenThreshold  int
		yellowThreshold int
		forceColor      string
		badgeLink       string
		badgeValue      string
		targetFile      string
		dryRun          bool
		help            bool
	)

	cp.flagParser.StringVar(&operation, "operation", "", fmt.Sprintf("実行する操作（必須: %s）", strings.Join(allowedOperations, ", ")))
	cp.flagParser.StringVar(&badgeTitle, "badge-title", defaultBadgeTitle, "バッジ左側のタイトル")
	cp.flagParser.StringVar(&coverageFile, "coverage-file", defaultCoverageFile, "go tool cover -func 出力ファイル")
	cp.flagParser.IntVar(&greenThreshold, "green-threshold", defaultGreenThreshold, "緑色へ切り替える閾値（0-100）")
	cp.flagParser.IntVar(&yellowThreshold, "yellow-threshold", defaultYellowThreshold, "黄色へ切り替える閾値（0-100）")
	cp.flagParser.StringVar(&forceColor, "force-color", "", "色を強制指定（green|yellow|red）")
	cp.flagParser.StringVar(&badgeLink, "badge-link", "", "バッジクリック時のリンク先URL")
	cp.flagParser.StringVar(&badgeValue, "badge-value", "", "カバレッジ値を手動指定（例: 58.6, 58.6%）")
	cp.flagParser.StringVar(&targetFile, "target-file", defaultTargetFile, "operation=patch-badge の更新対象ファイル")
	cp.flagParser.BoolVar(&dryRun, "dry-run", false, "operation=patch-badge の書き込みを行わず内容のみ出力")
	cp.flagParser.BoolVar(&help, "help", false, "使用方法を表示")
	cp.flagParser.BoolVar(&help, "h", false, "使用方法を表示（短縮形）")

	if err := cp.flagParser.Parse(); err != nil {
		return nil, err
	}

	return NewConfig(
		operation,
		badgeTitle,
		coverageFile,
		greenThreshold,
		yellowThreshold,
		forceColor,
		badgeLink,
		badgeValue,
		targetFile,
		dryRun,
		help,
	)
}

// ParseFlags は標準FlagParserとOSArgsを使ってConfigを返す。
func ParseFlags() (*Config, error) {
	parser := NewConfigParser(NewStandardFlagParser(), NewStandardOSArgs())
	return parser.ParseFlags()
}

// PrintUsage はCLIの利用方法を表示する。
func PrintUsage() {
	osArgs := NewStandardOSArgs()
	progName := "coverage-badge"
	if len(osArgs.Args()) > 0 && strings.TrimSpace(osArgs.Args()[0]) != "" {
		progName = osArgs.Args()[0]
	}

	fmt.Printf("使用方法: %s [オプション]\n", progName)
	fmt.Printf("\nオプション:\n")
	fmt.Printf("  -operation string\n")
	fmt.Printf("        実行する操作（必須: %s）\n", strings.Join(allowedOperations, ", "))
	fmt.Printf("  -badge-title string\n")
	fmt.Printf("        バッジ左側タイトル（デフォルト: %s）\n", defaultBadgeTitle)
	fmt.Printf("  -coverage-file string\n")
	fmt.Printf("        go tool cover -func 出力ファイル（デフォルト: %s）\n", defaultCoverageFile)
	fmt.Printf("  -green-threshold int\n")
	fmt.Printf("        緑色へ切り替える閾値（デフォルト: %d）\n", defaultGreenThreshold)
	fmt.Printf("  -yellow-threshold int\n")
	fmt.Printf("        黄色へ切り替える閾値（デフォルト: %d）\n", defaultYellowThreshold)
	fmt.Printf("  -force-color string\n")
	fmt.Printf("        強制色指定（green|yellow|red）\n")
	fmt.Printf("  -badge-link string\n")
	fmt.Printf("        バッジクリック時のリンク先URL\n")
	fmt.Printf("  -badge-value string\n")
	fmt.Printf("        カバレッジ値を直接指定（例: 58.6, 58.6%%）\n")
	fmt.Printf("  -target-file string\n")
	fmt.Printf("        operation=patch-badge の更新対象（デフォルト: %s）\n", defaultTargetFile)
	fmt.Printf("  -dry-run\n")
	fmt.Printf("        operation=patch-badge で書き込みせず結果のみ出力\n")
	fmt.Printf("  -help, -h\n")
	fmt.Printf("        使用方法を表示\n")
	fmt.Printf("\n使用例:\n")
	fmt.Printf("  %s -operation=create-badge -coverage-file=coverage.out\n", progName)
	fmt.Printf("  %s -operation=create-badge -badge-value=58.6 -badge-link=https://example.com/report\n", progName)
	fmt.Printf("  %s -operation=patch-badge -target-file=README.md -coverage-file=coverage.out\n", progName)
	fmt.Printf("  %s -operation=patch-badge -target-file=README.md -badge-value=72.1 -dry-run\n", progName)
}
