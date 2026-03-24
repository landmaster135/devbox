package config

import (
	"fmt"
	"strings"

	flagParser "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/flag_parser"
)

type Operation string

type PageType string

const (
	OperationDistributeFiles          Operation = "distribute-files"
	OperationCraftMarkdown            Operation = "craft-markdown"
	OperationCheckBodyLength          Operation = "check-body-length"
	OperationGrepStr                  Operation = "grep-str"
	OperationRenameBodiesByCategoryID Operation = "rename-bodies-by-category-id"
	OperationMigrateToMemos           Operation = "migrate-to-memos"
	PageTypeContent                   PageType  = "content"
	PageTypeArtifact                  PageType  = "artifact"
	PageTypeTask                      PageType  = "task"
)

type Config struct {
	Operation      Operation
	PageType       PageType
	BaseURL        string
	APIToken       string
	Category       string
	SkipsNoSrcBody bool
	SrcJSONFile    string
	SrcBodyDir     string
	SrcResourceDir string
	OutDir         string
	TargetStr      string
	ConNumberStart int
	ConNumberEnd   int
	Threshold      int
	Help           bool
}

func NewConfig(operation, pageType, baseURL, apiToken, category string, skipsNoSrcBody bool, srcJSONFile, srcBodyDir, srcResourceDir, outDir, targetStr string, conNumberStart, conNumberEnd, threshold int) (*Config, error) {
	trimmedOperation := Operation(strings.TrimSpace(operation))
	if trimmedOperation == "" {
		return nil, fmt.Errorf("operation パラメータは必須です")
	}
	if trimmedOperation != OperationDistributeFiles &&
		trimmedOperation != OperationCraftMarkdown &&
		trimmedOperation != OperationCheckBodyLength &&
		trimmedOperation != OperationGrepStr &&
		trimmedOperation != OperationRenameBodiesByCategoryID &&
		trimmedOperation != OperationMigrateToMemos {
		return nil, fmt.Errorf("未対応のoperationです: %s", trimmedOperation)
	}

	trimmedPageType := PageType(strings.TrimSpace(pageType))
	trimmedBaseURL := strings.TrimSpace(baseURL)
	trimmedAPIToken := strings.TrimSpace(apiToken)
	trimmedSrcJSONFile := strings.TrimSpace(srcJSONFile)
	trimmedSrcBodyDir := strings.TrimSpace(srcBodyDir)
	trimmedSrcResourceDir := strings.TrimSpace(srcResourceDir)
	trimmedOutDir := strings.TrimSpace(outDir)
	trimmedTargetStr := strings.TrimSpace(targetStr)

	switch trimmedOperation {
	case OperationDistributeFiles:
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
	case OperationCraftMarkdown:
		if trimmedPageType == "" {
			return nil, fmt.Errorf("page-type パラメータは必須です")
		}
		if trimmedPageType != PageTypeContent && trimmedPageType != PageTypeArtifact && trimmedPageType != PageTypeTask {
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
	case OperationRenameBodiesByCategoryID:
		if trimmedPageType == "" {
			return nil, fmt.Errorf("page-type パラメータは必須です")
		}
		if trimmedPageType != PageTypeContent {
			return nil, fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
		}
		if trimmedSrcJSONFile == "" {
			return nil, fmt.Errorf("src-json-file パラメータは必須です")
		}
		if trimmedSrcResourceDir == "" {
			return nil, fmt.Errorf("src-resource-dir パラメータは必須です")
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
	case OperationMigrateToMemos:
		if trimmedPageType == "" {
			return nil, fmt.Errorf("page-type パラメータは必須です")
		}
		if trimmedPageType != PageTypeContent {
			return nil, fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
		}
		if trimmedBaseURL == "" {
			return nil, fmt.Errorf("base-url パラメータは必須です")
		}
		if trimmedAPIToken == "" {
			return nil, fmt.Errorf("api-token パラメータは必須です")
		}
		if trimmedSrcBodyDir == "" {
			return nil, fmt.Errorf("src-body-dir パラメータは必須です")
		}
		if trimmedSrcResourceDir == "" {
			return nil, fmt.Errorf("src-resource-dir パラメータは必須です")
		}
	}

	if trimmedOperation == OperationCraftMarkdown || trimmedOperation == OperationRenameBodiesByCategoryID {
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
		BaseURL:        trimmedBaseURL,
		APIToken:       trimmedAPIToken,
		Category:       strings.TrimSpace(category),
		SkipsNoSrcBody: skipsNoSrcBody,
		SrcJSONFile:    trimmedSrcJSONFile,
		SrcBodyDir:     trimmedSrcBodyDir,
		SrcResourceDir: trimmedSrcResourceDir,
		OutDir:         trimmedOutDir,
		TargetStr:      trimmedTargetStr,
		ConNumberStart: conNumberStart,
		ConNumberEnd:   conNumberEnd,
		Threshold:      threshold,
	}, nil
}

func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParser())
}

func ParseFlagsWithParser(parser flagParser.FlagParser) (*Config, error) {
	var (
		operation      string
		pageType       string
		baseURL        string
		apiToken       string
		category       string
		skipsNoSrcBody bool
		srcJSONFile    string
		srcBodyDir     string
		srcResourceDir string
		outDir         string
		targetStr      string
		conNumberStart int
		conNumberEnd   int
		threshold      int
		help           bool
	)

	parser.StringVar(&operation, "operation", "", "操作タイプ（必須）")
	parser.StringVar(&pageType, "page-type", "", "ページタイプ（必須）")
	parser.StringVar(&baseURL, "base-url", "", "Memos API のベースURL（migrate-to-memosで必須）")
	parser.StringVar(&apiToken, "api-token", "", "Memos API のトークン（migrate-to-memosで必須）")
	parser.StringVar(&category, "category", "", "対象category（craft-markdownで任意）")
	parser.BoolVar(&skipsNoSrcBody, "skips-no-src-body", false, "コピー元MarkdownがないContentをスキップする（craft-markdownで任意）")
	parser.StringVar(&srcJSONFile, "src-json-file", "", "入力JSONファイルのパス（必須）")
	parser.StringVar(&srcJSONFile, "src-json-path", "", "入力JSONファイルのパス（後方互換）")
	parser.StringVar(&srcBodyDir, "src-body-dir", "", "Markdown本文ファイル群のディレクトリ（distribute-files/craft-markdown/check-body-length/grep-str/migrate-to-memosで必須）")
	parser.StringVar(&srcResourceDir, "src-resource-dir", "", "リソースファイル群のディレクトリ（rename-bodies-by-category-id/migrate-to-memosで必須）")
	parser.StringVar(&targetStr, "target-str", "", "検索文字列（grep-strで必須）")
	parser.StringVar(&outDir, "out-dir", "", "カテゴリ別出力先ディレクトリ（必須）")
	parser.IntVar(&conNumberStart, "con_number_start", 0, "con_id範囲の開始番号（craft-markdown/rename-bodies-by-category-idで必須）")
	parser.IntVar(&conNumberEnd, "con_number_end", 0, "con_id範囲の終了番号（craft-markdown/rename-bodies-by-category-idで必須）")
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
		baseURL,
		apiToken,
		category,
		skipsNoSrcBody,
		srcJSONFile,
		srcBodyDir,
		srcResourceDir,
		outDir,
		targetStr,
		conNumberStart,
		conNumberEnd,
		threshold,
	)
}
