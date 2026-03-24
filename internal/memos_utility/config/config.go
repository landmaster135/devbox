package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	flagParser "github.com/landmaster135/devbox/internal/memos_utility/infrastructures/flag_parser"
)

const (
	OperationCreateWebClip     = "create-web-clip"
	OperationCreateMovieClip   = "create-movie-clip"
	OperationCreateClips       = "create-clips"
	OperationCreateCommonMemos = "create-common-memos"

	defaultTimeoutSeconds = 30
)

var supportedOperations = []string{
	OperationCreateWebClip,
	OperationCreateMovieClip,
	OperationCreateClips,
	OperationCreateCommonMemos,
}

var webClipFilePattern = regexp.MustCompile(`^web-summary-(\d{8})-(\d{6})-([A-Za-z0-9][A-Za-z0-9_-]*)\.md$`)
var movieClipFilePattern = regexp.MustCompile(`^movie-summary-(\d{8})-(\d{6})-([A-Za-z0-9][A-Za-z0-9_-]*)\.md$`)

const usageTemplate = `使用方法: %[1]s [オプション]

Memos Utility CLI（Web/Movie サマリー用メモ作成）

共通オプション:
  -operation string
        実行する操作 (__SUPPORTED_OPERATIONS__)
  -base-url string
        Memos のベースURL (必須)
  -api-token string
        Memos API の Bearer トークン (必須)
  -timeout int
        HTTP リクエストのタイムアウト秒数 (デフォルト: __DEFAULT_TIMEOUT_SECONDS__)
  -content-file string
        メモ本文を読み込むファイルパス (create-web-clip/create-movie-clip で必須)
  -attachments string
        添付するファイルパス（任意、カンマ区切り。create-web-clip/create-movie-clip で利用）
  -content-dir string
        一括作成対象の Markdown ファイルを格納したディレクトリ (create-clips/create-common-memos で必須)
  -attachment-dir string
        一括作成時に添付ファイルを格納したディレクトリ (create-clips/create-common-memos で任意)
  -help, -h
        ヘルプを表示

operation 別 content-file 制約:
  create-web-clip
        web-summary-YYYYMMDD-hhmmss-<slug>.md 形式のみ
  create-movie-clip
        movie-summary-YYYYMMDD-hhmmss-<slug>.md 形式のみ

  create-clips
        content-dir 配下の全ファイルが web-summary-... または movie-summary-... 形式のみ
        attachment-dir 指定時は *_<number>.<extension> 形式のみ

  create-common-memos
        content-dir 配下の全ファイルが YYYYMMDDhhmmss_<number>.md 形式のみ
        attachment-dir 指定時は YYYYMMDDhhmmss_<number>_<index>.<extension> 形式のみ

例:
  %[1]s -operation=create-web-clip -base-url=https://memos.example.com -api-token=token -content-file=./web-summary-20240719-231059-sample.md
  %[1]s -operation=create-movie-clip -base-url=https://memos.example.com -api-token=token -content-file=./movie-summary-20260319-055716-sample.md -attachments=./a.png,./b.txt
  %[1]s -operation=create-clips -base-url=https://memos.example.com -api-token=token -content-dir=./clips -attachment-dir=./attachments
  %[1]s -operation=create-common-memos -base-url=https://memos.example.com -api-token=token -content-dir=./common-memos -attachment-dir=./attachments
`

// Config は memos-utility CLI の設定を保持する。
type Config struct {
	Operation      string
	BaseURL        string
	APIToken       string
	TimeoutSeconds int
	Help           bool

	ContentFile   string
	Attachments   string
	ContentDir    string
	AttachmentDir string
}

// ParseFlags は CLI フラグを解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParserFromOSArgs())
}

// ParseFlagsFromArgs はテスト容易性のため引数を受け取って解析する。
func ParseFlagsFromArgs(args []string) (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParser(os.Args[0], args))
}

// ParseFlagsWithParser は parser を受け取って解析する。
func ParseFlagsWithParser(parser flagParser.FlagParser) (*Config, error) {
	cfg := &Config{TimeoutSeconds: defaultTimeoutSeconds}

	parser.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	parser.StringVar(&cfg.BaseURL, "base-url", "", "Memos のベースURL (例: https://memos.example.com)")
	parser.StringVar(&cfg.APIToken, "api-token", "", "Memos API の Bearer トークン")
	parser.IntVar(&cfg.TimeoutSeconds, "timeout", cfg.TimeoutSeconds, "HTTP リクエストのタイムアウト秒数")
	parser.StringVar(&cfg.ContentFile, "content-file", "", "メモ本文を読み込むファイルのパス")
	parser.StringVar(&cfg.Attachments, "attachments", "", "添付するファイルパス（カンマ区切り）")
	parser.StringVar(&cfg.ContentDir, "content-dir", "", "一括作成対象の Markdown ファイルを格納したディレクトリ")
	parser.StringVar(&cfg.AttachmentDir, "attachment-dir", "", "一括作成時に添付ファイルを格納したディレクトリ")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")
	parser.BoolVar(&cfg.Help, "h", false, "ヘルプを表示（短縮）")

	if err := parser.Parse(); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return &Config{Help: true}, nil
		}
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if len(parser.Args()) > 0 {
		return nil, fmt.Errorf("未処理の位置引数があります: %v", parser.Args())
	}

	return NewConfig(*cfg)
}

// NewConfig は Config を正規化し、妥当性を検証して返す。
func NewConfig(raw Config) (*Config, error) {
	cfg := raw

	cfg.Operation = strings.ToLower(strings.TrimSpace(cfg.Operation))
	cfg.BaseURL = strings.TrimSpace(cfg.BaseURL)
	cfg.APIToken = strings.TrimSpace(cfg.APIToken)
	cfg.ContentFile = strings.TrimSpace(cfg.ContentFile)
	cfg.Attachments = strings.TrimSpace(cfg.Attachments)
	cfg.ContentDir = strings.TrimSpace(cfg.ContentDir)
	cfg.AttachmentDir = strings.TrimSpace(cfg.AttachmentDir)

	if cfg.Help {
		return &cfg, nil
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}
	if !isSupportedOperation(cfg.Operation) {
		return fmt.Errorf("未対応の operation です: %s", cfg.Operation)
	}
	if cfg.BaseURL == "" {
		return fmt.Errorf("base-url パラメータは必須です")
	}
	if cfg.APIToken == "" {
		return fmt.Errorf("api-token パラメータは必須です")
	}
	if cfg.TimeoutSeconds <= 0 {
		return fmt.Errorf("timeout パラメータは 1 以上で指定してください")
	}
	switch cfg.Operation {
	case OperationCreateWebClip, OperationCreateMovieClip:
		if cfg.ContentFile == "" {
			return fmt.Errorf("content-file パラメータは必須です")
		}
		if err := validateContentFileByOperation(cfg.Operation, cfg.ContentFile); err != nil {
			return err
		}
	case OperationCreateClips:
		if cfg.ContentDir == "" {
			return fmt.Errorf("content-dir パラメータは必須です")
		}
		if cfg.ContentFile != "" {
			return fmt.Errorf("create-clips では content-file は指定できません")
		}
		if cfg.Attachments != "" {
			return fmt.Errorf("create-clips では attachments は指定できません。attachment-dir を使用してください")
		}
	case OperationCreateCommonMemos:
		if cfg.ContentDir == "" {
			return fmt.Errorf("content-dir パラメータは必須です")
		}
		if cfg.ContentFile != "" {
			return fmt.Errorf("create-common-memos では content-file は指定できません")
		}
		if cfg.Attachments != "" {
			return fmt.Errorf("create-common-memos では attachments は指定できません。attachment-dir を使用してください")
		}
	default:
		return fmt.Errorf("未対応の operation です: %s", cfg.Operation)
	}

	return nil
}

func isSupportedOperation(operation string) bool {
	for _, op := range supportedOperations {
		if op == operation {
			return true
		}
	}
	return false
}

func validateContentFileByOperation(operation, contentFile string) error {
	baseName := filepath.Base(contentFile)

	switch operation {
	case OperationCreateWebClip:
		if !webClipFilePattern.MatchString(baseName) {
			return fmt.Errorf("create-web-clip の content-file は web-summary-YYYYMMDD-hhmmss-<slug>.md 形式のみ指定できます: %s", baseName)
		}
	case OperationCreateMovieClip:
		if !movieClipFilePattern.MatchString(baseName) {
			return fmt.Errorf("create-movie-clip の content-file は movie-summary-YYYYMMDD-hhmmss-<slug>.md 形式のみ指定できます: %s", baseName)
		}
	default:
		return fmt.Errorf("未対応の operation です: %s", operation)
	}

	return nil
}

// PrintUsage は CLI の利用方法を表示する。
func PrintUsage() {
	ops := make([]string, len(supportedOperations))
	copy(ops, supportedOperations)
	sort.Strings(ops)

	usage := strings.ReplaceAll(usageTemplate, "__SUPPORTED_OPERATIONS__", strings.Join(ops, ", "))
	usage = strings.ReplaceAll(usage, "__DEFAULT_TIMEOUT_SECONDS__", strconv.Itoa(defaultTimeoutSeconds))

	flagParser.PrintUsage(usage)
}
