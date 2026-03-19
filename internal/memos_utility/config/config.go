package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	OperationCreateWebClip   = "create-web-clip"
	OperationCreateMovieClip = "create-movie-clip"
	OperationCreateClips     = "create-clips"

	defaultTimeoutSeconds = 30
)

var supportedOperations = []string{
	OperationCreateWebClip,
	OperationCreateMovieClip,
	OperationCreateClips,
}

var webClipFilePattern = regexp.MustCompile(`^web-summary-(\d{8})-(\d{6})-([A-Za-z0-9][A-Za-z0-9_-]*)\.md$`)
var movieClipFilePattern = regexp.MustCompile(`^movie-summary-(\d{8})-(\d{6})-([A-Za-z0-9][A-Za-z0-9_-]*)\.md$`)

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
	return ParseFlagsFromArgs(os.Args[1:])
}

// ParseFlagsFromArgs はテスト容易性のため引数を受け取って解析する。
func ParseFlagsFromArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &Config{TimeoutSeconds: defaultTimeoutSeconds}

	fs.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	fs.StringVar(&cfg.BaseURL, "base-url", "", "Memos のベースURL (例: https://memos.example.com)")
	fs.StringVar(&cfg.APIToken, "api-token", "", "Memos API の Bearer トークン")
	fs.IntVar(&cfg.TimeoutSeconds, "timeout", cfg.TimeoutSeconds, "HTTP リクエストのタイムアウト秒数")
	fs.StringVar(&cfg.ContentFile, "content-file", "", "メモ本文を読み込むファイルのパス")
	fs.StringVar(&cfg.Attachments, "attachments", "", "添付するファイルパス（カンマ区切り）")
	fs.StringVar(&cfg.ContentDir, "content-dir", "", "一括作成対象の Markdown ファイルを格納したディレクトリ")
	fs.StringVar(&cfg.AttachmentDir, "attachment-dir", "", "一括作成時に添付ファイルを格納したディレクトリ")
	fs.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")
	fs.BoolVar(&cfg.Help, "h", false, "ヘルプを表示（短縮）")

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return &Config{Help: true}, nil
		}
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if len(fs.Args()) > 0 {
		return nil, fmt.Errorf("未処理の位置引数があります: %v", fs.Args())
	}

	cfg.Operation = strings.ToLower(strings.TrimSpace(cfg.Operation))
	cfg.BaseURL = strings.TrimSpace(cfg.BaseURL)
	cfg.APIToken = strings.TrimSpace(cfg.APIToken)
	cfg.ContentFile = strings.TrimSpace(cfg.ContentFile)
	cfg.Attachments = strings.TrimSpace(cfg.Attachments)
	cfg.ContentDir = strings.TrimSpace(cfg.ContentDir)
	cfg.AttachmentDir = strings.TrimSpace(cfg.AttachmentDir)

	if cfg.Help {
		return cfg, nil
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
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

	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Memos Utility CLI（Web/Movie サマリー用メモ作成）\n\n")

	fmt.Fprintf(os.Stderr, "共通オプション:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(ops, ", "))
	fmt.Fprintf(os.Stderr, "  -base-url string\n        Memos のベースURL (必須)\n")
	fmt.Fprintf(os.Stderr, "  -api-token string\n        Memos API の Bearer トークン (必須)\n")
	fmt.Fprintf(os.Stderr, "  -timeout int\n        HTTP リクエストのタイムアウト秒数 (デフォルト: %d)\n", defaultTimeoutSeconds)
	fmt.Fprintf(os.Stderr, "  -content-file string\n        メモ本文を読み込むファイルパス (create-web-clip/create-movie-clip で必須)\n")
	fmt.Fprintf(os.Stderr, "  -attachments string\n        添付するファイルパス（任意、カンマ区切り。create-web-clip/create-movie-clip で利用）\n")
	fmt.Fprintf(os.Stderr, "  -content-dir string\n        一括作成対象の Markdown ファイルを格納したディレクトリ (create-clips で必須)\n")
	fmt.Fprintf(os.Stderr, "  -attachment-dir string\n        一括作成時に添付ファイルを格納したディレクトリ (create-clips で任意)\n")
	fmt.Fprintf(os.Stderr, "  -help, -h\n        ヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "operation 別 content-file 制約:\n")
	fmt.Fprintf(os.Stderr, "  create-web-clip\n")
	fmt.Fprintf(os.Stderr, "        web-summary-YYYYMMDD-hhmmss-<slug>.md 形式のみ\n")
	fmt.Fprintf(os.Stderr, "  create-movie-clip\n")
	fmt.Fprintf(os.Stderr, "        movie-summary-YYYYMMDD-hhmmss-<slug>.md 形式のみ\n\n")
	fmt.Fprintf(os.Stderr, "  create-clips\n")
	fmt.Fprintf(os.Stderr, "        content-dir 配下の全ファイルが web-summary-... または movie-summary-... 形式のみ\n")
	fmt.Fprintf(os.Stderr, "        attachment-dir 指定時は *_<number>.<extension> 形式のみ\n\n")

	fmt.Fprintf(os.Stderr, "例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=create-web-clip -base-url=https://memos.example.com -api-token=token -content-file=./web-summary-20240719-231059-sample.md\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=create-movie-clip -base-url=https://memos.example.com -api-token=token -content-file=./movie-summary-20260319-055716-sample.md -attachments=./a.png,./b.txt\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=create-clips -base-url=https://memos.example.com -api-token=token -content-dir=./clips -attachment-dir=./attachments\n", os.Args[0])
}
