package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	OperationDistributeFiles = "distribute-files"
	PageTypeContent          = "content"
)

type Config struct {
	Operation   string
	PageType    string
	SrcJSONPath string
	SrcBodyDir  string
	OutDir      string
	Help        bool
}

func NewConfig(operation, pageType, srcJSONPath, srcBodyDir, outDir string) (*Config, error) {
	trimmedOperation := strings.TrimSpace(operation)
	if trimmedOperation == "" {
		return nil, fmt.Errorf("operation パラメータは必須です")
	}
	if trimmedOperation != OperationDistributeFiles {
		return nil, fmt.Errorf("未対応のoperationです: %s", trimmedOperation)
	}

	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType == "" {
		return nil, fmt.Errorf("page-type パラメータは必須です")
	}
	if trimmedPageType != PageTypeContent {
		return nil, fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}

	trimmedSrcJSONPath := strings.TrimSpace(srcJSONPath)
	if trimmedSrcJSONPath == "" {
		return nil, fmt.Errorf("src-json-path パラメータは必須です")
	}

	trimmedSrcBodyDir := strings.TrimSpace(srcBodyDir)
	if trimmedSrcBodyDir == "" {
		return nil, fmt.Errorf("src-body-dir パラメータは必須です")
	}

	trimmedOutDir := strings.TrimSpace(outDir)
	if trimmedOutDir == "" {
		return nil, fmt.Errorf("out-dir パラメータは必須です")
	}

	return &Config{
		Operation:   trimmedOperation,
		PageType:    trimmedPageType,
		SrcJSONPath: trimmedSrcJSONPath,
		SrcBodyDir:  trimmedSrcBodyDir,
		OutDir:      trimmedOutDir,
	}, nil
}

func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation   string
		pageType    string
		srcJSONPath string
		srcBodyDir  string
		outDir      string
		help        bool
	)

	parser.StringVar(&operation, "operation", "", "操作タイプ（必須）")
	parser.StringVar(&pageType, "page-type", "", "ページタイプ（必須）")
	parser.StringVar(&srcJSONPath, "src-json-path", "", "Content JSONファイルのパス（必須）")
	parser.StringVar(&srcBodyDir, "src-body-dir", "", "Markdown本文ファイル群のディレクトリ（必須）")
	parser.StringVar(&outDir, "out-dir", "", "カテゴリ別出力先ディレクトリ（必須）")
	parser.BoolVar(&help, "help", false, "ヘルプを表示")
	parser.BoolVar(&help, "h", false, "ヘルプを表示")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(operation, pageType, srcJSONPath, srcBodyDir, outDir)
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, `notion-to-memos-markdown CLI

使用方法:
  %s --operation=distribute-files --page-type=content --src-json-path=./tmp/contents.json --src-body-dir=./tmp/body --out-dir=./tmp/out

オプション:
  --operation      操作タイプ（必須: distribute-files）
  --page-type      ページタイプ（必須: content）
  --src-json-path  Content JSONファイルのパス（必須）
  --src-body-dir   con_id.md が配置されているディレクトリ（必須）
  --out-dir        カテゴリごとの出力先ルートディレクトリ（必須）
  -help, -h        このヘルプを表示
`, os.Args[0])
}
