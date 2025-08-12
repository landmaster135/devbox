package config

import (
	"fmt"
	"os"
)

// Config はCLIパラメータの設定を保持する構造体
type Config struct {
	SrcDir  string
	DestDir string
	Help    bool
}

// ParseFlags はコマンドライン引数を解析してConfigを返す
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	cfg := &Config{}

	parser.StringVar(&cfg.SrcDir, "src-dir", "", "ソースディレクトリのパス (必須)")
	parser.StringVar(&cfg.DestDir, "dest-dir", "", "出力先ディレクトリのパス (必須)")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	if cfg.Help {
		return cfg, nil
	}

	// 必須パラメータの検証
	if cfg.SrcDir == "" {
		return nil, fmt.Errorf("src-dir パラメータは必須です")
	}
	if cfg.DestDir == "" {
		return nil, fmt.Errorf("dest-dir パラメータは必須です")
	}

	return cfg, nil
}

// PrintUsage はヘルプメッセージを表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "使用方法: %s [オプション]\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Markdownファイルからブログコンテンツを抽出するツール\n\n")
	fmt.Fprintf(os.Stderr, "必須パラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -src-dir string\n")
	fmt.Fprintf(os.Stderr, "        ソースディレクトリのパス（Markdownファイルが格納されているディレクトリ）\n")
	fmt.Fprintf(os.Stderr, "  -dest-dir string\n")
	fmt.Fprintf(os.Stderr, "        出力先ディレクトリのパス（抽出したコンテンツを保存するディレクトリ）\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "オプションパラメータ:\n")
	fmt.Fprintf(os.Stderr, "  -help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプメッセージを表示\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "説明:\n")
	fmt.Fprintf(os.Stderr, "  このツールは、指定されたディレクトリ内のMarkdownファイルから、\n")
	fmt.Fprintf(os.Stderr, "  特定のマーカー（# Content → 空行 → ## はじまり）以降の内容を抽出し、\n")
	fmt.Fprintf(os.Stderr, "  新しいファイルとして出力先ディレクトリに保存します。\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  %s -src-dir=./blog-drafts -dest-dir=./extracted-content\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -src-dir=/path/to/markdown/files -dest-dir=/path/to/output\n", os.Args[0])
}
