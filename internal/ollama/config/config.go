package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	OperationVersion    = "version"
	OperationListModels = "list-models"
	OperationEmbed      = "embed"
	OperationGenerate   = "generate"
	OperationPull       = "pull"
	OperationDescribe   = "describe"
	OperationDelete     = "delete"

	defaultHost                = "127.0.0.1"
	defaultPort                = 11434
	defaultShortTimeoutSeconds = 30
	defaultLongTimeoutSeconds  = 300
)

var supportedOperations = []string{
	OperationEmbed,
	OperationDelete,
	OperationDescribe,
	OperationGenerate,
	OperationListModels,
	OperationPull,
	OperationVersion,
}

// Config は CLI フラグから構成情報を保持する。
type Config struct {
	Operation      string
	Host           string
	Port           int
	TimeoutSeconds int
	RunningOnly    bool
	Model          string
	Inputs         []string
	Prompt         string
	Help           bool
}

// ParseFlags は CLI フラグを解析する。
func ParseFlags() (*Config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &Config{}
	cfg.Operation = ""
	cfg.Host = defaultHost
	cfg.Port = defaultPort
	cfg.TimeoutSeconds = defaultShortTimeoutSeconds

	fs.StringVar(&cfg.Operation, "operation", cfg.Operation, fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	fs.StringVar(&cfg.Host, "host", cfg.Host, "Ollama API のホスト名")
	fs.IntVar(&cfg.Port, "port", cfg.Port, "Ollama API のポート番号")
	fs.BoolVar(&cfg.RunningOnly, "running-only", false, "list-models 操作時に稼働中モデルのみを表示")
	timeoutValue := newTrackedIntValue(defaultShortTimeoutSeconds)
	fs.Var(timeoutValue, "timeout", "HTTP リクエストのタイムアウト秒数 (デフォルト: 30, embed/generate/pull では 300)")
	fs.StringVar(&cfg.Model, "model", "", "embed/generate/pull/describe/delete で使用するモデル名")
	fs.StringVar(&cfg.Prompt, "prompt", "", "generate で送信するプロンプト")

	inputValues := &stringList{}
	fs.Var(inputValues, "input", "embed で使用する入力文字列 (複数指定可)")

	fs.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return &Config{Help: true}, nil
		}
		return nil, fmt.Errorf("フラグの解析に失敗しました: %w", err)
	}

	if len(fs.Args()) > 0 {
		return nil, fmt.Errorf("未処理の位置引数があります: %v", fs.Args())
	}

	cfg.Operation = strings.ToLower(strings.TrimSpace(cfg.Operation))
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Model = strings.TrimSpace(cfg.Model)
	cfg.Prompt = strings.TrimSpace(cfg.Prompt)
	cfg.Inputs = inputValues.Items()
	cfg.TimeoutSeconds = timeoutValue.Value()

	if cfg.Help {
		return cfg, nil
	}

	timeoutChanged := timeoutValue.Changed()

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	if !timeoutChanged {
		if needsLongTimeout(cfg.Operation) {
			cfg.TimeoutSeconds = defaultLongTimeoutSeconds
		} else {
			cfg.TimeoutSeconds = defaultShortTimeoutSeconds
		}
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Operation == "" {
		return fmt.Errorf("operation パラメータは必須です")
	}

	if !isSupportedOperation(cfg.Operation) {
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}

	if cfg.Host == "" {
		return fmt.Errorf("host パラメータは必須です")
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("port パラメータは 1 から 65535 の範囲で指定してください")
	}

	if cfg.TimeoutSeconds <= 0 {
		return fmt.Errorf("timeout パラメータは 1 以上で指定してください")
	}

	switch cfg.Operation {
	case OperationVersion:
		// 追加パラメータ不要
	case OperationListModels:
		// running-only のみ利用
	case OperationEmbed:
		if cfg.Model == "" {
			return fmt.Errorf("embed 操作には model パラメータが必要です")
		}
		if len(cfg.Inputs) == 0 {
			return fmt.Errorf("embed 操作には 1 件以上の input パラメータが必要です")
		}
	case OperationGenerate:
		if cfg.Model == "" {
			return fmt.Errorf("generate 操作には model パラメータが必要です")
		}
		if cfg.Prompt == "" {
			return fmt.Errorf("generate 操作には prompt パラメータが必要です")
		}
	case OperationPull:
		if cfg.Model == "" {
			return fmt.Errorf("pull 操作には model パラメータが必要です")
		}
	case OperationDescribe, OperationDelete:
		if cfg.Model == "" {
			return fmt.Errorf("%s 操作には model パラメータが必要です", cfg.Operation)
		}
	}

	return nil
}

func isSupportedOperation(operation string) bool {
	for _, op := range supportedOperations {
		if operation == op {
			return true
		}
	}
	return false
}

func needsLongTimeout(operation string) bool {
	switch operation {
	case OperationEmbed, OperationGenerate, OperationPull:
		return true
	default:
		return false
	}
}

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Ollama API クライアント\n\n")

	fmt.Fprintf(os.Stderr, "共通オプション:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(supportedOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -host string\n        Ollama API のホスト名 (デフォルト: %s)\n", defaultHost)
	fmt.Fprintf(os.Stderr, "  -port int\n        Ollama API のポート番号 (デフォルト: %d)\n", defaultPort)
	fmt.Fprintf(os.Stderr, "  -timeout int\n        HTTP リクエストのタイムアウト秒数 (デフォルト: %d、embed/generate/pull では %d)\n", defaultShortTimeoutSeconds, defaultLongTimeoutSeconds)
	fmt.Fprintf(os.Stderr, "  -help\n        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "operation 別パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  %s\n        その他の追加パラメータ不要\n", OperationVersion)
	fmt.Fprintf(os.Stderr, "  %s\n        -running-only (任意): 稼働中モデルのみを表示\n", OperationListModels)
	fmt.Fprintf(os.Stderr, "  %s\n        -model (必須): モデル名\n        -input (複数可・必須): 埋め込み対象テキスト\n", OperationEmbed)
	fmt.Fprintf(os.Stderr, "  %s\n        -model (必須): モデル名\n        -prompt (必須): プロンプト文字列\n", OperationGenerate)
	fmt.Fprintf(os.Stderr, "  %s\n        -model (必須): モデル名\n", OperationPull)
	fmt.Fprintf(os.Stderr, "  %s\n        -model (必須): 詳細を取得するモデル名\n", OperationDescribe)
	fmt.Fprintf(os.Stderr, "  %s\n        -model (必須): 削除するモデル名\n", OperationDelete)
}

// stringList は複数回指定可能な input フラグを表す。
type stringList struct {
	items []string
}

func (s *stringList) String() string {
	return strings.Join(s.items, ",")
}

func (s *stringList) Set(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fmt.Errorf("input パラメータに空文字は指定できません")
	}
	s.items = append(s.items, trimmed)
	return nil
}

func (s *stringList) Items() []string {
	if len(s.items) == 0 {
		return nil
	}
	result := make([]string, len(s.items))
	copy(result, s.items)
	return result
}

type trackedIntValue struct {
	value   int
	changed bool
}

func newTrackedIntValue(defaultValue int) *trackedIntValue {
	return &trackedIntValue{value: defaultValue}
}

func (t *trackedIntValue) String() string {
	return strconv.Itoa(t.value)
}

func (t *trackedIntValue) Set(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	t.value = v
	t.changed = true
	return nil
}

func (t *trackedIntValue) Value() int {
	return t.value
}

func (t *trackedIntValue) Changed() bool {
	return t.changed
}

func init() {
	sort.Strings(supportedOperations)
}
