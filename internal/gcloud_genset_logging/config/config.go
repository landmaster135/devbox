package config

import (
	"fmt"
	"os"
)

const (
	OperationLoggingRead = "logging-read"
	OperationCreateSink  = "create-sink"

	defaultLimit = 10
)

// Config はCLIで指定されたパラメータを保持する
type Config struct {
	Operation      string
	Help           bool
	Severity       string
	Limit          int
	Query          string
	ResourceType   string
	Filter         string
	AdditionalArgs string
	SinkName       string
	Destination    string
	LogFilter      string
}

// ParseFlags は標準のフラグパーサーを使用して引数を解析する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたフラグパーサーで引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{
		Limit: defaultLimit,
	}

	parser.StringVar(&cfg.Operation, "operation", "", "実行する操作 (必須: logging-read, create-sink)")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	// logging-read 向けパラメータ
	parser.StringVar(&cfg.Severity, "severity", "", "ログの重要度 (例: ERROR, WARNING)")
	parser.IntVar(&cfg.Limit, "limit", defaultLimit, "取得するログの最大数 (logging-read 操作用)")
	parser.StringVar(&cfg.Query, "query", "", "追加のクエリフィルター")
	parser.StringVar(&cfg.ResourceType, "resource-type", "", "リソースタイプ (例: gce_instance)")
	parser.StringVar(&cfg.Filter, "filter", "", "完全なフィルター文字列を指定")

	// 共通で利用する追加引数
	parser.StringVar(&cfg.AdditionalArgs, "additional-args", "", "gcloud コマンドに追加する引数を指定")

	// create-sink 向けパラメータ
	parser.StringVar(&cfg.SinkName, "sink-name", "", "作成するシンク名 (create-sink 操作用)")
	parser.StringVar(&cfg.Destination, "destination", "", "シンクの送信先 (例: storage.googleapis.com/my-bucket)")
	parser.StringVar(&cfg.LogFilter, "log-filter", "", "シンクに適用するログフィルター")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	if cfg.Operation == "" {
		return nil, fmt.Errorf("operation パラメータは必須です")
	}

	switch cfg.Operation {
	case OperationLoggingRead:
		if cfg.Filter == "" && cfg.Severity == "" && cfg.Query == "" && cfg.ResourceType == "" {
			return nil, fmt.Errorf("logging-read 操作にはフィルター条件を少なくとも1つ指定してください")
		}
		if cfg.Limit <= 0 {
			return nil, fmt.Errorf("limit パラメータは1以上で指定してください")
		}
	case OperationCreateSink:
		if cfg.SinkName == "" {
			return nil, fmt.Errorf("sink-name パラメータは必須です")
		}
		if cfg.Destination == "" {
			return nil, fmt.Errorf("destination パラメータは必須です")
		}
	default:
		return nil, fmt.Errorf("未対応の操作です: %s", cfg.Operation)
	}

	return cfg, nil
}

// PrintUsage はCLIの利用方法を標準エラーに出力する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Google Cloud Logging 用の gcloud コマンド生成ツール\n\n")
	fmt.Fprintf(os.Stderr, "共通パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (logging-read | create-sink)\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n\n")

	fmt.Fprintf(os.Stderr, "logging-read 操作用のパラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -severity string\n")
	fmt.Fprintf(os.Stderr, "        ログの重要度 (例: ERROR)\n")
	fmt.Fprintf(os.Stderr, "  -limit int\n")
	fmt.Fprintf(os.Stderr, "        取得するログの最大数 (デフォルト: %d)\n", defaultLimit)
	fmt.Fprintf(os.Stderr, "  -query string\n")
	fmt.Fprintf(os.Stderr, "        追加のクエリフィルター\n")
	fmt.Fprintf(os.Stderr, "  -resource-type string\n")
	fmt.Fprintf(os.Stderr, "        ログのリソースタイプ (例: gce_instance)\n")
	fmt.Fprintf(os.Stderr, "  -filter string\n")
	fmt.Fprintf(os.Stderr, "        完全なフィルター文字列を直接指定\n")
	fmt.Fprintf(os.Stderr, "  -additional-args string\n")
	fmt.Fprintf(os.Stderr, "        gcloud logging read に渡す追加引数\n\n")

	fmt.Fprintf(os.Stderr, "create-sink 操作用のパラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -sink-name string\n")
	fmt.Fprintf(os.Stderr, "        作成するシンク名 (必須)\n")
	fmt.Fprintf(os.Stderr, "  -destination string\n")
	fmt.Fprintf(os.Stderr, "        シンクの送信先 (必須)\n")
	fmt.Fprintf(os.Stderr, "  -log-filter string\n")
	fmt.Fprintf(os.Stderr, "        シンクに適用するログフィルター\n")
	fmt.Fprintf(os.Stderr, "  -additional-args string\n")
	fmt.Fprintf(os.Stderr, "        gcloud logging sinks create に渡す追加引数\n\n")

	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -operation=logging-read -severity=ERROR -resource-type=gce_instance\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=logging-read -filter='resource.type=gce_instance' -limit=20\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -operation=create-sink -sink-name=my-sink -destination=storage.googleapis.com/my-bucket\n", os.Args[0])
}
