package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

const (
	OperationCreateMemo = "create-memo"
	OperationGetMemo    = "get-memo"
	OperationListMemos  = "list-memos"
	OperationUpdateMemo = "update-memo"

	defaultTimeoutSeconds = 30
	defaultPageSize       = 20
)

var supportedOperations = []string{
	OperationCreateMemo,
	OperationGetMemo,
	OperationListMemos,
	OperationUpdateMemo,
}

var supportedVisibility = map[string]struct{}{
	"VISIBILITY_UNSPECIFIED": {},
	"PRIVATE":                {},
	"PROTECTED":              {},
	"PUBLIC":                 {},
}

var supportedState = map[string]struct{}{
	"STATE_UNSPECIFIED": {},
	"NORMAL":            {},
	"ARCHIVED":          {},
}

// Config は memos CLI の設定を保持する。
type Config struct {
	Operation      string
	BaseURL        string
	APIToken       string
	TimeoutSeconds int
	Help           bool

	MemoID string
	Memo   string

	Content     string
	Visibility  string
	State       string
	Pinned      bool
	PinnedSet   bool
	DisplayTime string

	PageSize  int
	PageToken string
	OrderBy   string

	UpdateMask string
}

// ParseFlags は CLI フラグを解析する。
func ParseFlags() (*Config, error) {
	return ParseFlagsFromArgs(os.Args[1:])
}

// ParseFlagsFromArgs はテスト容易性のため引数を受け取って解析する。
func ParseFlagsFromArgs(args []string) (*Config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &Config{
		TimeoutSeconds: defaultTimeoutSeconds,
		PageSize:       defaultPageSize,
	}

	pinnedValue := newTrackedBoolValue(false)

	fs.StringVar(&cfg.Operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	fs.StringVar(&cfg.BaseURL, "base-url", "", "Memos のベースURL (例: https://memos.example.com)")
	fs.StringVar(&cfg.APIToken, "api-token", "", "Memos API の Bearer トークン")
	fs.IntVar(&cfg.TimeoutSeconds, "timeout", cfg.TimeoutSeconds, "HTTP リクエストのタイムアウト秒数")
	fs.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")
	fs.BoolVar(&cfg.Help, "h", false, "ヘルプを表示（短縮）")

	fs.StringVar(&cfg.MemoID, "memo-id", "", "create-memo で作成する memoId（任意）")
	fs.StringVar(&cfg.Memo, "memo", "", "get-memo/update-memo で対象にする memo 識別子")
	fs.StringVar(&cfg.Content, "content", "", "create-memo/update-memo で設定する本文")
	fs.StringVar(&cfg.Visibility, "visibility", "", "visibility（PRIVATE/PROTECTED/PUBLIC）")
	fs.StringVar(&cfg.State, "state", "", "state（NORMAL/ARCHIVED）")
	fs.Var(pinnedValue, "pinned", "pinned の設定値（true/false）")
	fs.StringVar(&cfg.DisplayTime, "display-time", "", "create-memo で設定する表示日時（RFC3339）")

	fs.IntVar(&cfg.PageSize, "page-size", cfg.PageSize, "list-memos の取得件数")
	fs.StringVar(&cfg.PageToken, "page-token", "", "list-memos のページトークン")
	fs.StringVar(&cfg.OrderBy, "order-by", "", "list-memos のソート指定（例: update_time desc）")
	fs.StringVar(&cfg.UpdateMask, "update-mask", "", "update-memo の updateMask（例: content,visibility）")

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
	cfg.MemoID = strings.TrimSpace(cfg.MemoID)
	cfg.Memo = strings.TrimSpace(cfg.Memo)
	cfg.Content = strings.TrimSpace(cfg.Content)
	cfg.Visibility = strings.ToUpper(strings.TrimSpace(cfg.Visibility))
	cfg.State = strings.ToUpper(strings.TrimSpace(cfg.State))
	cfg.DisplayTime = strings.TrimSpace(cfg.DisplayTime)
	cfg.PageToken = strings.TrimSpace(cfg.PageToken)
	cfg.OrderBy = strings.TrimSpace(cfg.OrderBy)
	cfg.UpdateMask = strings.TrimSpace(cfg.UpdateMask)
	cfg.Pinned = pinnedValue.Value()
	cfg.PinnedSet = pinnedValue.Changed()

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
	if cfg.PageSize < 0 {
		return fmt.Errorf("page-size パラメータは 0 以上で指定してください")
	}
	if cfg.Visibility != "" && !isSupportedValue(cfg.Visibility, supportedVisibility) {
		return fmt.Errorf("visibility の値が不正です: %s", cfg.Visibility)
	}
	if cfg.State != "" && !isSupportedValue(cfg.State, supportedState) {
		return fmt.Errorf("state の値が不正です: %s", cfg.State)
	}

	switch cfg.Operation {
	case OperationCreateMemo:
		if cfg.Content == "" {
			return fmt.Errorf("create-memo 操作には content パラメータが必要です")
		}
	case OperationGetMemo:
		if cfg.Memo == "" {
			return fmt.Errorf("get-memo 操作には memo パラメータが必要です")
		}
	case OperationListMemos:
		if cfg.PageSize == 0 {
			return fmt.Errorf("list-memos 操作では page-size に 1 以上を指定してください")
		}
	case OperationUpdateMemo:
		if cfg.Memo == "" {
			return fmt.Errorf("update-memo 操作には memo パラメータが必要です")
		}
		if cfg.Content == "" {
			return fmt.Errorf("update-memo 操作には content パラメータが必要です")
		}
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

func isSupportedValue(value string, candidates map[string]struct{}) bool {
	_, ok := candidates[value]
	return ok
}

// PrintUsage は CLI の利用方法を表示する。
func PrintUsage() {
	ops := make([]string, len(supportedOperations))
	copy(ops, supportedOperations)
	sort.Strings(ops)

	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Memos API クライアント\n\n")

	fmt.Fprintf(os.Stderr, "共通オプション:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(ops, ", "))
	fmt.Fprintf(os.Stderr, "  -base-url string\n        Memos のベースURL (必須)\n")
	fmt.Fprintf(os.Stderr, "  -api-token string\n        Memos API の Bearer トークン (必須)\n")
	fmt.Fprintf(os.Stderr, "  -timeout int\n        HTTP リクエストのタイムアウト秒数 (デフォルト: %d)\n", defaultTimeoutSeconds)
	fmt.Fprintf(os.Stderr, "  -help, -h\n        ヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "operation 別オプション:\n")
	fmt.Fprintf(os.Stderr, "  create-memo\n")
	fmt.Fprintf(os.Stderr, "        -content (必須), -memo-id (任意), -visibility (任意), -state (任意), -pinned (任意), -display-time (任意)\n")
	fmt.Fprintf(os.Stderr, "  get-memo\n")
	fmt.Fprintf(os.Stderr, "        -memo (必須)\n")
	fmt.Fprintf(os.Stderr, "  list-memos\n")
	fmt.Fprintf(os.Stderr, "        -page-size (任意), -page-token (任意), -state (任意), -order-by (任意)\n")
	fmt.Fprintf(os.Stderr, "  update-memo\n")
	fmt.Fprintf(os.Stderr, "        -memo (必須), -content (必須), -visibility (任意), -state (任意), -pinned (任意), -update-mask (任意)\n\n")

	fmt.Fprintf(os.Stderr, "例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=create-memo -base-url=https://memos.example.com -api-token=token -content='hello'\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=get-memo -base-url=https://memos.example.com -api-token=token -memo=abc123\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=list-memos -base-url=https://memos.example.com -api-token=token -page-size=20 -state=NORMAL\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=update-memo -base-url=https://memos.example.com -api-token=token -memo=abc123 -content='updated'\n", os.Args[0])
}

type trackedBoolValue struct {
	value   bool
	changed bool
}

func newTrackedBoolValue(defaultValue bool) *trackedBoolValue {
	return &trackedBoolValue{value: defaultValue}
}

func (t *trackedBoolValue) String() string {
	if t.value {
		return "true"
	}
	return "false"
}

func (t *trackedBoolValue) Set(value string) error {
	parsed, err := parseBool(value)
	if err != nil {
		return err
	}
	t.value = parsed
	t.changed = true
	return nil
}

func (t *trackedBoolValue) IsBoolFlag() bool {
	return true
}

func (t *trackedBoolValue) Value() bool {
	return t.value
}

func (t *trackedBoolValue) Changed() bool {
	return t.changed
}

func parseBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "1", "t", "true", "y", "yes":
		return true, nil
	case "0", "f", "false", "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("bool 値が不正です: %s", value)
	}
}
