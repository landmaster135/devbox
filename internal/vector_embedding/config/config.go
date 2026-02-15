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
	OperationOllama = "ollama"
	OperationOpenAI = "openai"

	defaultHost           = "127.0.0.1"
	defaultPort           = 11434
	defaultTimeoutSeconds = 60
)

var supportedOperations = []string{OperationOllama, OperationOpenAI}

// Config は CLI で受け取る設定値を保持する。
type Config struct {
	Operation      string
	Host           string
	Port           int
	Model          string
	APIKey         string
	Inputs         []string
	TimeoutSeconds int
	Help           bool
}

// ParseFlags は CLI フラグを解析する。
func ParseFlags() (*Config, error) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &Config{
		Host:           defaultHost,
		Port:           defaultPort,
		TimeoutSeconds: defaultTimeoutSeconds,
	}

	fs.StringVar(&cfg.Operation, "operation", cfg.Operation, fmt.Sprintf("実行する操作 (%s)", strings.Join(supportedOperations, ", ")))
	fs.StringVar(&cfg.Host, "host", cfg.Host, "Ollama API のホスト名")
	fs.IntVar(&cfg.Port, "port", cfg.Port, "Ollama API のポート番号")
	fs.StringVar(&cfg.Model, "model", cfg.Model, "使用するモデル名")
	fs.StringVar(&cfg.APIKey, "api-key", cfg.APIKey, "OpenAI API キー (operation=openai で必須)")
	inputValues := &stringList{}
	fs.Var(inputValues, "input", "ベクトル化するテキスト (複数指定可)")
	timeoutValue := newTrackedIntValue(defaultTimeoutSeconds)
	fs.Var(timeoutValue, "timeout", "HTTP リクエストのタイムアウト秒数")
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
	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	cfg.Inputs = inputValues.Items()
	cfg.TimeoutSeconds = timeoutValue.Value()

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
		return fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}
	if cfg.TimeoutSeconds <= 0 {
		return fmt.Errorf("timeout パラメータは 1 以上を指定してください")
	}

	switch cfg.Operation {
	case OperationOllama:
		if cfg.Host == "" {
			return fmt.Errorf("ollama 操作には host パラメータが必要です")
		}
		if cfg.Port <= 0 || cfg.Port > 65535 {
			return fmt.Errorf("ollama 操作の port パラメータが不正です: %d", cfg.Port)
		}
		if cfg.Model == "" {
			return fmt.Errorf("ollama 操作には model パラメータが必要です")
		}
		if len(cfg.Inputs) == 0 {
			return fmt.Errorf("ollama 操作には 1 件以上の input パラメータが必要です")
		}
	case OperationOpenAI:
		if cfg.APIKey == "" {
			return fmt.Errorf("openai 操作には api-key パラメータが必要です")
		}
		if cfg.Model == "" {
			return fmt.Errorf("openai 操作には model パラメータが必要です")
		}
		if len(cfg.Inputs) == 0 {
			return fmt.Errorf("openai 操作には 1 件以上の input パラメータが必要です")
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

// PrintUsage は CLI の利用方法を標準エラーに出力する。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Vector Embedding CLI\n\n")
	fmt.Fprintf(os.Stderr, "共通オプション:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n        実行する操作 (%s)\n", strings.Join(supportedOperations, ", "))
	fmt.Fprintf(os.Stderr, "  -host string\n        Ollama API のホスト名 (デフォルト: %s)\n", defaultHost)
	fmt.Fprintf(os.Stderr, "  -port int\n        Ollama API のポート番号 (デフォルト: %d)\n", defaultPort)
	fmt.Fprintf(os.Stderr, "  -api-key string\n        OpenAI API キー (operation=openai で必須)\n")
	fmt.Fprintf(os.Stderr, "  -timeout int\n        HTTP リクエストのタイムアウト秒数 (デフォルト: %d)\n", defaultTimeoutSeconds)
	fmt.Fprintf(os.Stderr, "  -help\n        このヘルプを表示\n\n")
	fmt.Fprintf(os.Stderr, "操作別パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  %s\n        -model (必須): モデル名\n        -input (複数可・必須): 埋め込み対象のテキスト\n        -host / -port: 呼び出す Ollama API エンドポイント\n", OperationOllama)
	fmt.Fprintf(os.Stderr, "  %s\n        -model (必須): モデル名\n        -input (複数可・必須): 埋め込み対象のテキスト\n        -api-key (必須): OpenAI API キー\n", OperationOpenAI)
}

// stringList は複数の -input フラグを収集する。
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
	out := make([]string, len(s.items))
	copy(out, s.items)
	return out
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
