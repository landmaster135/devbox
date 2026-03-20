package config

import (
	"fmt"
	"os"
	"strings"
)

// Config はカラーコード変換CLIの設定を保持する構造体
type Config struct {
	SrcFormat  string // 変換元のカラーコード形式 (hex, rgb, hsl, hsv, dec)
	DestFormat string // 変換先のカラーコード形式 (hex, rgb, hsl, hsv, dec)
	Value      string // 変換するカラーコード値
	Help       bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(srcFormat, destFormat, value string) (*Config, error) {
	if srcFormat == "" {
		return nil, fmt.Errorf("変換元形式が指定されていません")
	}
	if destFormat == "" {
		return nil, fmt.Errorf("変換先形式が指定されていません")
	}
	if value == "" {
		return nil, fmt.Errorf("変換するカラーコード値が指定されていません")
	}

	// サポートされている形式の検証
	validFormats := []string{"hex", "rgb", "hsl", "hsv", "dec"}
	if !isValidFormat(srcFormat, validFormats) {
		return nil, fmt.Errorf("無効な変換元形式です: %s (サポート形式: %s)", srcFormat, strings.Join(validFormats, ", "))
	}
	if !isValidFormat(destFormat, validFormats) {
		return nil, fmt.Errorf("無効な変換先形式です: %s (サポート形式: %s)", destFormat, strings.Join(validFormats, ", "))
	}

	return &Config{
		SrcFormat:  strings.ToLower(srcFormat),
		DestFormat: strings.ToLower(destFormat),
		Value:      strings.TrimSpace(value),
	}, nil
}

// isValidFormat は指定された形式が有効かどうかを確認する
func isValidFormat(format string, validFormats []string) bool {
	format = strings.ToLower(format)
	for _, valid := range validFormats {
		if format == valid {
			return true
		}
	}
	return false
}

// FlagParser はフラグ解析のインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser は標準のflagパッケージを使用するFlagParser
type StandardFlagParser struct {
	args []string
}

// NewStandardFlagParser は新しいStandardFlagParserを作成する
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{}
}

func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	// 実際の実装では flag.StringVar を使用
	for i, arg := range os.Args {
		if arg == "-"+name || arg == "--"+name {
			if i+1 < len(os.Args) {
				*ptr = os.Args[i+1]
			}
		}
	}
}

func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	// 実際の実装では flag.BoolVar を使用
	for _, arg := range os.Args {
		if arg == "-"+name || arg == "--"+name {
			*ptr = true
		}
	}
}

func (p *StandardFlagParser) Parse() error {
	// 実際の実装では flag.Parse を使用
	return nil
}

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
		srcFormat  = ""
		destFormat = ""
		value      = ""
		help       = false
	)

	parser.StringVar(&srcFormat, "src-format", srcFormat, "変換元のカラーコード形式 (hex, rgb, hsl, hsv, dec)")
	parser.StringVar(&srcFormat, "s", srcFormat, "変換元形式の短縮形")

	parser.StringVar(&destFormat, "dest-format", destFormat, "変換先のカラーコード形式 (hex, rgb, hsl, hsv, dec)")
	parser.StringVar(&destFormat, "d", destFormat, "変換先形式の短縮形")

	parser.StringVar(&value, "value", value, "変換するカラーコード値")
	parser.StringVar(&value, "v", value, "カラーコード値の短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 残りの引数から値を取得（位置引数として）
	args := parser.Args()
	if len(args) >= 3 {
		srcFormat = args[0]
		destFormat = args[1]
		value = args[2]
	}

	return NewConfig(srcFormat, destFormat, value)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `カラーコード変換CLIツール

使用方法:
  フラグ形式:
    %s -src-format hex -dest-format rgb -value "#FF0000"
    %s -s hex -d rgb -v "#FF0000"

  位置引数形式:
    %s hex rgb "#FF0000"

サポートされているカラーコード形式:
  hex  - 16進数形式 (例: #FF0000, #ff0000)
  rgb  - RGB形式 (例: rgb(255,0,0))
  hsl  - HSL形式 (例: hsl(0,100%%,50%%))
  hsv  - HSV形式 (例: hsv(0,100%%,100%%))
  dec  - 10進数形式 (例: 16711680)

オプション:
  -src-format, -s   変換元のカラーコード形式
  -dest-format, -d  変換先のカラーコード形式
  -value, -v        変換するカラーコード値
  -help, -h         このヘルプを表示

例:
  %s -s hex -d rgb -v "#FF0000"     # HEXからRGBに変換
  %s -s rgb -d hsl -v "rgb(255,0,0)" # RGBからHSLに変換
  %s hex hsv "#00FF00"              # HEXからHSVに変換（位置引数）
  %s -s dec -d hex -v "16711680"    # 10進数からHEXに変換

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
