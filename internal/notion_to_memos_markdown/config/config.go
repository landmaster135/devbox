package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	OperationDistributeFiles = "distribute-files"
	OperationCraftMarkdown   = "craft-markdown"
	OperationCheckBodyLength = "check-body-length"
	OperationGrepStr         = "grep-str"
	PageTypeContent          = "content"
)

type Config struct {
	Operation      string
	PageType       string
	Category       string
	SkipsNoSrcBody bool
	SrcJSONFile    string
	SrcBodyDir     string
	OutDir         string
	TargetStr      string
	ConNumberStart int
	ConNumberEnd   int
	Threshold      int
	Help           bool
}

func NewConfig(operation, pageType, category string, skipsNoSrcBody bool, srcJSONFile, srcBodyDir, outDir, targetStr string, conNumberStart, conNumberEnd, threshold int) (*Config, error) {
	trimmedOperation := strings.TrimSpace(operation)
	if trimmedOperation == "" {
		return nil, fmt.Errorf("operation パラメータは必須です")
	}
	if trimmedOperation != OperationDistributeFiles &&
		trimmedOperation != OperationCraftMarkdown &&
		trimmedOperation != OperationCheckBodyLength &&
		trimmedOperation != OperationGrepStr {
		return nil, fmt.Errorf("未対応のoperationです: %s", trimmedOperation)
	}

	trimmedPageType := strings.TrimSpace(pageType)
	trimmedSrcJSONFile := strings.TrimSpace(srcJSONFile)
	trimmedSrcBodyDir := strings.TrimSpace(srcBodyDir)
	trimmedOutDir := strings.TrimSpace(outDir)
	trimmedTargetStr := strings.TrimSpace(targetStr)

	switch trimmedOperation {
	case OperationDistributeFiles, OperationCraftMarkdown:
		if trimmedPageType == "" {
			return nil, fmt.Errorf("page-type パラメータは必須です")
		}
		if trimmedPageType != PageTypeContent {
			return nil, fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
		}
		if trimmedSrcJSONFile == "" {
			return nil, fmt.Errorf("src-json-file パラメータは必須です")
		}
		if trimmedSrcBodyDir == "" {
			return nil, fmt.Errorf("src-body-dir パラメータは必須です")
		}
		if trimmedOutDir == "" {
			return nil, fmt.Errorf("out-dir パラメータは必須です")
		}
	case OperationCheckBodyLength:
		if trimmedSrcBodyDir == "" {
			return nil, fmt.Errorf("src-body-dir パラメータは必須です")
		}
		if threshold < 0 {
			return nil, fmt.Errorf("threshold パラメータは0以上で必須です")
		}
	case OperationGrepStr:
		if trimmedSrcBodyDir == "" {
			return nil, fmt.Errorf("src-body-dir パラメータは必須です")
		}
		if trimmedTargetStr == "" {
			return nil, fmt.Errorf("target-str パラメータは必須です")
		}
	}

	if trimmedOperation == OperationCraftMarkdown {
		if conNumberStart <= 0 {
			return nil, fmt.Errorf("con_number_start パラメータは1以上で必須です")
		}
		if conNumberEnd <= 0 {
			return nil, fmt.Errorf("con_number_end パラメータは1以上で必須です")
		}
		if conNumberStart > conNumberEnd {
			return nil, fmt.Errorf("con_number_start は con_number_end 以下である必要があります")
		}
	}

	return &Config{
		Operation:      trimmedOperation,
		PageType:       trimmedPageType,
		Category:       strings.TrimSpace(category),
		SkipsNoSrcBody: skipsNoSrcBody,
		SrcJSONFile:    trimmedSrcJSONFile,
		SrcBodyDir:     trimmedSrcBodyDir,
		OutDir:         trimmedOutDir,
		TargetStr:      trimmedTargetStr,
		ConNumberStart: conNumberStart,
		ConNumberEnd:   conNumberEnd,
		Threshold:      threshold,
	}, nil
}

func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation      string
		pageType       string
		category       string
		skipsNoSrcBody bool
		srcJSONFile    string
		srcBodyDir     string
		outDir         string
		targetStr      string
		conNumberStart int
		conNumberEnd   int
		threshold      int
		help           bool
	)

	parser.StringVar(&operation, "operation", "", "操作タイプ（必須）")
	parser.StringVar(&pageType, "page-type", "", "ページタイプ（必須）")
	parser.StringVar(&category, "category", "", "対象category（craft-markdownで任意）")
	parser.BoolVar(&skipsNoSrcBody, "skips-no-src-body", false, "コピー元MarkdownがないContentをスキップする（craft-markdownで任意）")
	parser.StringVar(&srcJSONFile, "src-json-file", "", "Content JSONファイルのパス（必須）")
	parser.StringVar(&srcJSONFile, "src-json-path", "", "Content JSONファイルのパス（後方互換）")
	parser.StringVar(&srcBodyDir, "src-body-dir", "", "Markdown本文ファイル群のディレクトリ（必須）")
	parser.StringVar(&targetStr, "target-str", "", "検索文字列（grep-strで必須）")
	parser.StringVar(&outDir, "out-dir", "", "カテゴリ別出力先ディレクトリ（必須）")
	parser.IntVar(&conNumberStart, "con_number_start", 0, "con_id範囲の開始番号（craft-markdownで必須）")
	parser.IntVar(&conNumberEnd, "con_number_end", 0, "con_id範囲の終了番号（craft-markdownで必須）")
	parser.IntVar(&threshold, "threshold", -1, "文字数の閾値（check-body-lengthで必須、0以上）")
	parser.BoolVar(&help, "help", false, "ヘルプを表示")
	parser.BoolVar(&help, "h", false, "ヘルプを表示")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	if help {
		return &Config{Help: true}, nil
	}

	return NewConfig(
		operation,
		pageType,
		category,
		skipsNoSrcBody,
		srcJSONFile,
		srcBodyDir,
		outDir,
		targetStr,
		conNumberStart,
		conNumberEnd,
		threshold,
	)
}

func PrintUsage() {
	fmt.Fprintf(os.Stderr, `notion-to-memos-markdown CLI

使用方法:
  %s --operation=distribute-files --page-type=content --src-json-file=./tmp/contents.json --src-body-dir=./tmp/body --out-dir=./tmp/out
  %s --operation=craft-markdown --page-type=content --category=software --skips-no-src-body=false --con_number_start=1 --con_number_end=9999 --src-json-file=./tmp/contents.json --src-body-dir=./tmp/body --out-dir=./tmp/out
  %s --operation=check-body-length --src-body-dir=./tmp/body --threshold=1000
  %s --operation=grep-str --src-body-dir=./tmp/body --target-str=TODO

オプション:
  --operation        操作タイプ（必須: distribute-files, craft-markdown, check-body-length, grep-str）
  --page-type        ページタイプ（distribute-files/craft-markdownで必須: content）
  --category         対象category（craft-markdownで任意。指定時は一致するContentのみ処理）
  --skips-no-src-body コピー元Markdownなしをスキップするか（craft-markdownで任意。デフォルト:false）
  --con_number_start craft-markdown時の開始con番号（必須）
  --con_number_end   craft-markdown時の終了con番号（必須）
  --threshold        文字数の閾値（check-body-lengthで必須、0以上）
  --target-str       検索文字列（grep-strで必須）
  --src-json-file    Content JSONファイルのパス（distribute-files/craft-markdownで必須）
  --src-body-dir     入力ディレクトリ（全operationで必須）
  --out-dir          出力先ルートディレクトリ（distribute-files/craft-markdownで必須）
  -help, -h          このヘルプを表示
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
