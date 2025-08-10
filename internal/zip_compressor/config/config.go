package config

import (
	"fmt"
	"os"
)

// Config はZip圧縮CLIの設定を保持する構造体
type Config struct {
	Operation string // 操作タイプ (compress, decompress)
	Path      string // 対象ファイル/ディレクトリのパス
	Help      bool   // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation, path string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{"compress", "decompress"}
	isValid := false
	for _, op := range validOperations {
		if operation == op {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, fmt.Errorf("無効な操作タイプです: %s", operation)
	}

	if path == "" {
		return nil, fmt.Errorf("パスが指定されていません")
	}

	return &Config{
		Operation: operation,
		Path:      path,
	}, nil
}

// FlagParser インターフェースを定義（テスト用）
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
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
	return &StandardFlagParser{
		args: os.Args[1:],
	}
}

// StringVar は文字列フラグを定義する
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	// 簡易実装：実際のflagパッケージの代替
	for i, arg := range p.args {
		if arg == "-"+name || arg == "--"+name {
			if i+1 < len(p.args) {
				*ptr = p.args[i+1]
			}
		}
	}
}

// BoolVar はブールフラグを定義する
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	// 簡易実装：実際のflagパッケージの代替
	for _, arg := range p.args {
		if arg == "-"+name || arg == "--"+name {
			*ptr = true
		}
	}
}

// Parse はフラグを解析する
func (p *StandardFlagParser) Parse() error {
	return nil
}

// Args は残りの引数を返す
func (p *StandardFlagParser) Args() []string {
	var remainingArgs []string
	skipNext := false

	for _, arg := range p.args {
		if skipNext {
			skipNext = false
			continue
		}

		if arg == "-operation" || arg == "--operation" ||
			arg == "-o" || arg == "--o" ||
			arg == "-path" || arg == "--path" ||
			arg == "-p" || arg == "--p" {
			skipNext = true
			continue
		}

		if arg == "-help" || arg == "--help" || arg == "-h" || arg == "--h" {
			continue
		}

		if !skipNext && !startsWith(arg, "-") {
			remainingArgs = append(remainingArgs, arg)
		}
	}

	return remainingArgs
}

// startsWith は文字列が指定されたプレフィックスで始まるかチェックする
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation = ""
		path      = ""
		help      = false
	)

	parser.StringVar(&operation, "operation", operation, "操作タイプ (compress, decompress)")
	parser.StringVar(&operation, "o", operation, "操作タイプの短縮形")

	parser.StringVar(&path, "path", path, "対象ファイル/ディレクトリのパス")
	parser.StringVar(&path, "p", path, "パスの短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 残りの引数から operation, path を取得（位置引数として）
	args := parser.Args()
	if len(args) >= 1 && operation == "" {
		operation = args[0]
	}
	if len(args) >= 2 && path == "" {
		path = args[1]
	}

	return NewConfig(operation, path)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `Zip圧縮CLIツール

使用方法:
  ファイル/ディレクトリ圧縮:
    %s -operation compress -path /path/to/file_or_directory
    %s -o compress -p /path/to/file_or_directory
    %s compress /path/to/file_or_directory

  Zipファイル展開:
    %s -operation decompress -path /path/to/archive.zip
    %s -o decompress -p /path/to/archive.zip
    %s decompress /path/to/archive.zip

オプション:
  -operation, -o    操作タイプ (compress, decompress)
  -path, -p         対象ファイル/ディレクトリのパス
  -help, -h         このヘルプを表示

例:
  # ファイル圧縮
  %s compress /home/user/document.txt
  # → document.txt.zip が作成される

  # ディレクトリ圧縮
  %s compress /home/user/my_folder
  # → my_folder.zip が作成される

  # Zipファイル展開
  %s decompress /home/user/archive.zip
  # → archive_decompressed/ ディレクトリに展開される

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
