package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/markdown_crafter/domain"
)

var allowedOperations = []string{
	domain.OperationSplitHeadings,
	domain.OperationAddFrontMatter,
	domain.OperationAddTags,
}

type Config struct {
	Operation    string
	FilePath     string
	HeadingLevel int
	OutputDir    string
	KVPairs      []string
	Tags         string
	Help         bool
}

func NewConfig(operation, filePath string, headingLevel int, outputDir string, kvPairs []string, tags string, help bool) (*Config, error) {
	cfg := &Config{
		Operation:    operation,
		FilePath:     filePath,
		HeadingLevel: headingLevel,
		OutputDir:    outputDir,
		KVPairs:      kvPairs,
		Tags:         tags,
		Help:         help,
	}

	if cfg.Help {
		return cfg, nil
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if strings.TrimSpace(c.Operation) == "" {
		return fmt.Errorf("--operation は必須です")
	}
	if !isAllowed(c.Operation, allowedOperations) {
		return fmt.Errorf("--operation には %s のいずれかを指定してください", strings.Join(allowedOperations, ", "))
	}
	if strings.TrimSpace(c.FilePath) == "" {
		return fmt.Errorf("--file-path は必須です")
	}

	switch c.Operation {
	case domain.OperationSplitHeadings:
		if c.HeadingLevel < 1 || c.HeadingLevel > 6 {
			return fmt.Errorf("--heading-level は 1 から 6 の範囲で指定してください")
		}
		if strings.TrimSpace(c.OutputDir) == "" {
			return fmt.Errorf("--output-dir は必須です (--operation=split-headings)")
		}
	case domain.OperationAddFrontMatter:
		if len(c.KVPairs) == 0 {
			return fmt.Errorf("--kv は1件以上指定してください (--operation=add-front-matter)")
		}
	case domain.OperationAddTags:
		if strings.TrimSpace(c.Tags) == "" {
			return fmt.Errorf("--tags は必須です (--operation=add-tags)")
		}
	}

	return nil
}

func isAllowed(value string, candidates []string) bool {
	for _, candidate := range candidates {
		if candidate == value {
			return true
		}
	}
	return false
}

type stringSliceValue []string

func (s *stringSliceValue) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceValue) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func ParseFlags() (*Config, error) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	flagSet.Usage = func() {
		PrintUsage()
	}

	var (
		operation    string
		filePath     string
		headingLevel int
		outputDir    string
		tags         string
		help         bool
		kvPairs      stringSliceValue
	)

	flagSet.StringVar(&operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(allowedOperations, ", ")))
	flagSet.StringVar(&filePath, "file-path", "", "対象のMarkdownファイルパス")
	flagSet.IntVar(&headingLevel, "heading-level", 0, "分割対象の見出しレベル (1-6, split-headingsで必須)")
	flagSet.StringVar(&outputDir, "output-dir", "", "分割後ファイルの出力先ディレクトリ (split-headingsで必須)")
	flagSet.Var(&kvPairs, "kv", "front matter に追加する key=value (add-front-matter で複数指定可)")
	flagSet.StringVar(&tags, "tags", "", "追加するタグ（カンマ区切り, 例: go,markdown）")
	flagSet.BoolVar(&help, "help", false, "ヘルプを表示")
	flagSet.BoolVar(&help, "h", false, "ヘルプを表示 (短縮形)")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return NewConfig(operation, filePath, headingLevel, outputDir, []string(kvPairs), tags, help)
}

func PrintUsage() {
	exeName := os.Args[0]
	fmt.Fprintf(os.Stderr, "Markdown Crafter CLI ツール\n\n")
	fmt.Fprintf(os.Stderr, "使用方法:\n")
	fmt.Fprintf(os.Stderr, "  %s --operation split-headings --file-path ./note.md --heading-level 2 --output-dir ./out\n", exeName)
	fmt.Fprintf(os.Stderr, "  %s --operation add-front-matter --file-path ./note.md --kv title=記事 --kv author=nov\n", exeName)
	fmt.Fprintf(os.Stderr, "  %s --operation add-tags --file-path ./note.md --tags go,markdown\n\n", exeName)

	fmt.Fprintf(os.Stderr, "オプション:\n")
	fmt.Fprintf(os.Stderr, "  --operation string\n")
	fmt.Fprintf(os.Stderr, "        実行する操作 (%s)\n", strings.Join(allowedOperations, ", "))
	fmt.Fprintf(os.Stderr, "  --file-path string\n")
	fmt.Fprintf(os.Stderr, "        対象のMarkdownファイルパス\n")
	fmt.Fprintf(os.Stderr, "  --heading-level int\n")
	fmt.Fprintf(os.Stderr, "        分割対象の見出しレベル (1-6, split-headingsで必須)\n")
	fmt.Fprintf(os.Stderr, "  --output-dir string\n")
	fmt.Fprintf(os.Stderr, "        分割後ファイルの出力先ディレクトリ (split-headingsで必須)\n")
	fmt.Fprintf(os.Stderr, "  --kv key=value\n")
	fmt.Fprintf(os.Stderr, "        front matter に追加する key=value (add-front-matter で複数指定可)\n")
	fmt.Fprintf(os.Stderr, "  --tags string\n")
	fmt.Fprintf(os.Stderr, "        追加するタグ（カンマ区切り, 例: go,markdown）\n")
	fmt.Fprintf(os.Stderr, "  --help\n")
	fmt.Fprintf(os.Stderr, "        このヘルプを表示\n")
}
