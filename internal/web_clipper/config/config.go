package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	OperationPatchMarkdown = "patch-markdown"
)

var allowedOperations = []string{
	OperationPatchMarkdown,
}

type Config struct {
	Operation          string
	TargetTitle        string
	TargetURL          string
	SrcMarkdownContent string
	SrcMarkdownFile    string
	OutFilePath        string
	TopHeadingLevel    int
	Help               bool
}

func NewConfig(operation, targetTitle, targetURL, srcMarkdownContent, srcMarkdownFile, outFilePath string, topHeadingLevel int, help bool) (*Config, error) {
	trimmedOperation := strings.TrimSpace(operation)
	trimmedTargetTitle := strings.TrimSpace(targetTitle)
	trimmedTargetURL := strings.TrimSpace(targetURL)
	trimmedSrcMarkdownFile := strings.TrimSpace(srcMarkdownFile)
	trimmedOutFilePath := strings.TrimSpace(outFilePath)

	cfg := &Config{
		Operation:          trimmedOperation,
		TargetTitle:        trimmedTargetTitle,
		TargetURL:          trimmedTargetURL,
		SrcMarkdownContent: srcMarkdownContent,
		SrcMarkdownFile:    trimmedSrcMarkdownFile,
		OutFilePath:        trimmedOutFilePath,
		TopHeadingLevel:    topHeadingLevel,
		Help:               help,
	}

	if cfg.Help {
		return cfg, nil
	}

	if cfg.Operation == "" {
		return nil, fmt.Errorf("--operation は必須です")
	}
	if !isAllowed(cfg.Operation, allowedOperations) {
		return nil, fmt.Errorf("--operation には %s を指定してください", strings.Join(allowedOperations, ", "))
	}

	if cfg.Operation == OperationPatchMarkdown {
		if cfg.TargetTitle == "" {
			return nil, fmt.Errorf("--target-title は必須です (--operation=patch-markdown)")
		}
		if cfg.TargetURL == "" {
			return nil, fmt.Errorf("--target-url は必須です (--operation=patch-markdown)")
		}
		if strings.TrimSpace(cfg.SrcMarkdownContent) == "" && cfg.SrcMarkdownFile == "" {
			return nil, fmt.Errorf("--src-markdown-content または --src-markdown-file のいずれかは必須です (--operation=patch-markdown)")
		}
		if strings.TrimSpace(cfg.SrcMarkdownContent) != "" && cfg.SrcMarkdownFile != "" {
			return nil, fmt.Errorf("--src-markdown-content と --src-markdown-file は同時に指定できません")
		}
		if cfg.OutFilePath == "" {
			return nil, fmt.Errorf("--out-file-path は必須です (--operation=patch-markdown)")
		}
		if cfg.TopHeadingLevel < 1 {
			return nil, fmt.Errorf("--top-heading-level は 1 以上で指定してください")
		}
	}

	return cfg, nil
}

func ParseFlags() (*Config, error) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	flagSet.Usage = PrintUsage

	var (
		operation          string
		targetTitle        string
		targetURL          string
		srcMarkdownContent string
		srcMarkdownFile    string
		outFilePath        string
		topHeadingLevel    int
		help               bool
	)

	flagSet.StringVar(&operation, "operation", "", fmt.Sprintf("実行する操作（%s）", strings.Join(allowedOperations, ", ")))
	flagSet.StringVar(&targetTitle, "target-title", "", "リンク表示用の記事タイトル（必須）")
	flagSet.StringVar(&targetURL, "target-url", "", "リンク先URL（必須）")
	flagSet.StringVar(&srcMarkdownContent, "src-markdown-content", "", "処理対象Markdown本文")
	flagSet.StringVar(&srcMarkdownFile, "src-markdown-file", "", "処理対象Markdownファイルパス")
	flagSet.StringVar(&outFilePath, "out-file-path", "", "出力先ファイルパス（必須）")
	flagSet.IntVar(&topHeadingLevel, "top-heading-level", 0, "リンク挿入対象の見出しレベル（1以上, 必須）")
	flagSet.BoolVar(&help, "help", false, "ヘルプを表示")
	flagSet.BoolVar(&help, "h", false, "ヘルプを表示（短縮形）")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return NewConfig(
		operation,
		targetTitle,
		targetURL,
		srcMarkdownContent,
		srcMarkdownFile,
		outFilePath,
		topHeadingLevel,
		help,
	)
}

func PrintUsage() {
	exeName := os.Args[0]
	fmt.Fprintf(os.Stderr, `web-clipper は記事要約MarkdownへWeb記事リンクを挿入するCLIツールです。

使用方法:
  %[1]s --operation=patch-markdown --target-title="OpenAI" --target-url="https://openai.com" --src-markdown-content=$'## 記事タイトル 要約\n\n### 見出し\n本文\n' --out-file-path=./out.md --top-heading-level=2
  %[1]s --operation=patch-markdown --target-title="OpenAI" --target-url="https://openai.com" --src-markdown-file=./in.md --out-file-path=./out.md --top-heading-level=2
  %[1]s --help

オプション:
  --operation         実行する操作（必須: patch-markdown）
  --target-title      リンク表示用の記事タイトル（必須）
  --target-url        リンク先URL（必須）
  --src-markdown-content 処理対象Markdown本文（--src-markdown-file と同時指定不可）
  --src-markdown-file 処理対象Markdownファイルパス（--src-markdown-content と同時指定不可）
  --out-file-path     出力先ファイルパス（必須）
  --top-heading-level リンク挿入対象の見出しレベル（1以上, 必須）
  -help, -h           ヘルプを表示
`, exeName)
}

func isAllowed(value string, candidates []string) bool {
	for _, candidate := range candidates {
		if candidate == value {
			return true
		}
	}

	return false
}
