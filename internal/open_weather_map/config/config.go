package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config はOpenWeather API CLIの設定を保持する構造体
type Config struct {
	APIKey  string // APIキー（必須）
	City    string // 都市名（必須）
	MaxDays int    // 最大日数（必須、1-5の範囲）
	Help    bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(apiKey, city string, maxDays int) (*Config, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("APIキーが指定されていません")
	}

	if city == "" {
		return nil, fmt.Errorf("都市名が指定されていません")
	}

	if maxDays < 1 || maxDays > 5 {
		return nil, fmt.Errorf("最大日数は1-5の範囲で指定してください: %d", maxDays)
	}

	return &Config{
		APIKey:  apiKey,
		City:    city,
		MaxDays: maxDays,
	}, nil
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
		maxDaysStr = "3"
		help       = false
	)

	parser.StringVar(&apiKey, "api-key", apiKey, "OpenWeather API キー（必須）")
	parser.StringVar(&apiKey, "k", apiKey, "APIキーの短縮形")

	parser.StringVar(&city, "city", city, "都市名（例: Tokyo,JP）（必須）")
	parser.StringVar(&city, "c", city, "都市名の短縮形")

	parser.StringVar(&maxDaysStr, "max-days", maxDaysStr, "最大日数（1-5）（必須）")
	parser.StringVar(&maxDaysStr, "d", maxDaysStr, "最大日数の短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 文字列から数値に変換
	maxDays, err := strconv.Atoi(maxDaysStr)
	if err != nil {
		return nil, fmt.Errorf("無効な最大日数です: %s", maxDaysStr)
	}

	// 残りの引数から位置引数を取得
	args := parser.Args()
	if len(args) >= 1 && apiKey == "" {
		apiKey = args[0]
	}
	if len(args) >= 2 && city == "" {
		city = args[1]
	}
	if len(args) >= 3 {
		if val, err := strconv.Atoi(args[2]); err == nil {
			maxDays = val
		}
	}

	return NewConfig(apiKey, city, maxDays)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `OpenWeather API CLIツール

使用方法:
  基本的な使用方法:
    %s -api-key YOUR_API_KEY -city "Tokyo,JP" -max-days 3
    %s -k YOUR_API_KEY -c "Tokyo,JP" -d 5

  位置引数での指定:
    %s YOUR_API_KEY "Tokyo,JP" 3

オプション:
  -api-key, -k     OpenWeather APIキー（必須）
  -city, -c        都市名（例: Tokyo,JP, London,UK）（必須）
  -max-days, -d    取得する最大日数（1-5）（必須、デフォルト: 3）
  -help, -h        このヘルプを表示

例:
  %s -k abc123 -c "Tokyo,JP" -d 3
  %s -api-key abc123 -city "London,UK" -max-days 5
  %s abc123 "New York,US" 2

注意:
  - APIキーはOpenWeatherMapから取得してください
  - 都市名は "都市名,国コード" の形式で指定してください
  - 無料プランでは最大5日間の予報が取得可能です

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
