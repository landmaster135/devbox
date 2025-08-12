package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config は天気通知CLIの設定を保持する構造体
type Config struct {
	APIKey     string // OpenWeather API キー（必須）
	City       string // 都市名（必須）
	MaxDays    int    // 最大日数（必須）
	WebhookURL string // Discord Webhook URL（必須）
	Help       bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(apiKey, city string, maxDays int, webhookURL string) (*Config, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API キーが指定されていません")
	}

	if city == "" {
		return nil, fmt.Errorf("都市名が指定されていません")
	}

	if maxDays <= 0 {
		return nil, fmt.Errorf("最大日数は1以上である必要があります")
	}

	if maxDays > 5 {
		return nil, fmt.Errorf("最大日数は5日以下である必要があります（OpenWeather API制限）")
	}

	if webhookURL == "" {
		return nil, fmt.Errorf("Discord Webhook URLが指定されていません")
	}

	return &Config{
		APIKey:     apiKey,
		City:       city,
		MaxDays:    maxDays,
		WebhookURL: webhookURL,
	}, nil
}

// FlagParser はフラグ解析のインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	IntVar(p *int, name string, value int, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser は標準のflagパッケージを使用する実装
type StandardFlagParser struct {
	args []string
}

// NewStandardFlagParser は新しいStandardFlagParserを作成する
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{}
}

// StringVar は文字列フラグを定義する
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	for i, arg := range os.Args {
		if arg == "-"+name || arg == "--"+name {
			if i+1 < len(os.Args) {
				*ptr = os.Args[i+1]
			}
		}
	}
}

// IntVar は整数フラグを定義する
func (p *StandardFlagParser) IntVar(ptr *int, name string, value int, usage string) {
	for i, arg := range os.Args {
		if arg == "-"+name || arg == "--"+name {
			if i+1 < len(os.Args) {
				if val, err := strconv.Atoi(os.Args[i+1]); err == nil {
					*ptr = val
				}
			}
		}
	}
}

// BoolVar はブールフラグを定義する
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	for _, arg := range os.Args {
		if arg == "-"+name || arg == "--"+name {
			*ptr = true
		}
	}
}

// Parse はフラグを解析する
func (p *StandardFlagParser) Parse() error {
	p.args = os.Args[1:]
	return nil
}

// Args は残りの引数を返す
func (p *StandardFlagParser) Args() []string {
	return p.args
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		apiKey     = ""
		city       = ""
		maxDays    = 1
		webhookURL = ""
		help       = false
	)

	parser.StringVar(&apiKey, "api-key", apiKey, "OpenWeather API キー（必須）")
	parser.StringVar(&city, "city", city, "都市名（必須）")
	parser.IntVar(&maxDays, "max-days", maxDays, "最大日数（必須、1-5日）")
	parser.StringVar(&webhookURL, "webhook-url", webhookURL, "Discord Webhook URL（必須）")
	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(apiKey, city, maxDays, webhookURL)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `天気通知CLIツール

使用方法:
  %s -api-key YOUR_API_KEY -city Tokyo -max-days 3 -webhook-url YOUR_WEBHOOK_URL

オプション:
  -api-key      OpenWeather API キー（必須）
  -city         都市名（必須）例: Tokyo, Osaka, "New York"
  -max-days     最大日数（必須、1-5日）
  -webhook-url  Discord Webhook URL（必須）
  -help, -h     このヘルプを表示

例:
  %s -api-key abc123 -city Tokyo -max-days 3 -webhook-url https://discord.com/api/webhooks/...
  %s -api-key abc123 -city "New York" -max-days 5 -webhook-url https://discord.com/api/webhooks/...

`, os.Args[0], os.Args[0], os.Args[0])
}
