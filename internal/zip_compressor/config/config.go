package config

import (
	"fmt"
	"slices"

	flagParser "github.com/landmaster135/devbox/internal/zip_compressor/infrastructures/flag_parser"
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

	validOperations := []string{"compress", "decompress"}
	if !slices.Contains(validOperations, operation) {
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

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser flagParser.FlagParser) (*Config, error) {
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

	if help {
		return &Config{Help: true}, nil
	}

	args := parser.Args()
	if len(args) >= 1 && operation == "" {
		operation = args[0]
	}
	if len(args) >= 2 && path == "" {
		path = args[1]
	}

	return NewConfig(operation, path)
}
