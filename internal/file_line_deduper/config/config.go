package config

import (
	"fmt"
	"os"
)

// Config は file-line-deduper CLI の設定です。
type Config struct {
	FilePath string
	StartPos int
	EndPos   int
	Help     bool
}

// NewConfig は検証済みConfigを生成します。
func NewConfig(filePath string, startPos, endPos int) (*Config, error) {
	if filePath == "" {
		return nil, fmt.Errorf("ファイルパスを指定してください（-file オプション）")
	}

	if startPos == 0 && endPos == 0 {
		return nil, fmt.Errorf("開始位置と終了位置を指定してください（-start と -end オプション）")
	}

	return &Config{
		FilePath: filePath,
		StartPos: startPos,
		EndPos:   endPos,
	}, nil
}

// ParseFlags は os.Args からフラグを解析します。
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithArgs は指定引数からフラグを解析します。
func ParseFlagsWithArgs(args []string) (*Config, error) {
	return ParseFlagsWithParser(newStandardFlagParser(args))
}

// ParseFlagsWithParser は指定されたパーサーでフラグを解析します。
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		filePath = ""
		startPos = 0
		endPos   = 0
		help     = false
	)

	parser.StringVar(&filePath, "file", filePath, "処理するファイルのパス")
	parser.IntVar(&startPos, "start", startPos, "各行の文字列を取得する開始位置（0ベース）")
	parser.IntVar(&endPos, "end", endPos, "各行の文字列を取得する終了位置")
	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプを表示（短縮形）")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(filePath, startPos, endPos)
}

// PrintUsage は使い方を標準エラーに表示します。
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `file-line-deduper は指定範囲の文字列をキーに重複行を削除するCLIツールです。

使用方法:
  %s -file <ファイルパス> -start <開始位置> -end <終了位置>
  %s -help

オプション:
  -file    処理対象ファイルのパス（必須）
  -start   文字列抽出の開始位置（0ベース、必須）
  -end     文字列抽出の終了位置（必須）
  -help    ヘルプを表示

例:
  %s -file data.txt -start 5 -end 10
`, os.Args[0], os.Args[0], os.Args[0])
}
