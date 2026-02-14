package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

const (
	FormatJSON            = "json"
	FormatYAML            = "yaml"
	FormatCSV             = "csv"
	FormatTSV             = "tsv"
	FormatHTML            = "html"
	FormatMDOrderedList   = "md-ordered-list"
	FormatMDUnorderedList = "md-unordered-list"
	FormatMDTable         = "md-table"
)

// Config は data-converter のCLI設定です。
type Config struct {
	InputFilePath  string
	OutputFilePath string
	InputFormat    string
	OutputFormat   string
	Help           bool
}

var supportedFormats = map[string]struct{}{
	FormatJSON:            {},
	FormatYAML:            {},
	FormatCSV:             {},
	FormatTSV:             {},
	FormatHTML:            {},
	FormatMDOrderedList:   {},
	FormatMDUnorderedList: {},
	FormatMDTable:         {},
}

var extensionToFormat = map[string]string{
	".json": FormatJSON,
	".yaml": FormatYAML,
	".yml":  FormatYAML,
	".csv":  FormatCSV,
	".tsv":  FormatTSV,
	".html": FormatHTML,
	".htm":  FormatHTML,
}

// ParseFlags はCLI引数を解析し、設定を返します。
func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("data-converter", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cfg := &Config{}
	fs.StringVar(&cfg.InputFilePath, "input-file-path", "", "入力ファイルパス")
	fs.StringVar(&cfg.OutputFilePath, "output-file-path", "", "出力ファイルパス")
	fs.StringVar(&cfg.InputFormat, "input-format", "", "入力形式 (json|yaml|csv|tsv|html|md-ordered-list|md-unordered-list|md-table)")
	fs.StringVar(&cfg.OutputFormat, "output-format", "", "出力形式 (json|yaml|csv|tsv|html|md-ordered-list|md-unordered-list|md-table)")
	fs.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if cfg.Help {
		return cfg, nil
	}

	cfg.InputFilePath = strings.TrimSpace(cfg.InputFilePath)
	cfg.OutputFilePath = strings.TrimSpace(cfg.OutputFilePath)
	cfg.InputFormat = normalizeFormat(cfg.InputFormat)
	cfg.OutputFormat = normalizeFormat(cfg.OutputFormat)

	if cfg.InputFilePath == "" {
		return nil, errors.New("-input-file-path を指定してください")
	}
	if cfg.OutputFilePath == "" {
		return nil, errors.New("-output-file-path を指定してください")
	}

	if cfg.InputFormat == "" {
		inferred, err := inferFormatFromPath(cfg.InputFilePath)
		if err != nil {
			return nil, fmt.Errorf("入力形式を判定できません: %w", err)
		}
		cfg.InputFormat = inferred
	}

	if cfg.OutputFormat == "" {
		inferred, err := inferFormatFromPath(cfg.OutputFilePath)
		if err != nil {
			return nil, fmt.Errorf("出力形式を判定できません: %w", err)
		}
		cfg.OutputFormat = inferred
	}

	if !isSupportedFormat(cfg.InputFormat) {
		return nil, fmt.Errorf("未対応の入力形式です: %s", cfg.InputFormat)
	}
	if !isSupportedFormat(cfg.OutputFormat) {
		return nil, fmt.Errorf("未対応の出力形式です: %s", cfg.OutputFormat)
	}

	return cfg, nil
}

// PrintUsage はCLIの使い方を表示します。
func PrintUsage(w io.Writer) {
	fmt.Fprintln(w, "Data Converter")
	fmt.Fprintln(w, "任意の入力ファイルを key-value リストへ正規化してから、別形式ファイルへ変換します。")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "使用例:")
	fmt.Fprintln(w, "  go run ./cmd/cli/data-converter \\")
	fmt.Fprintln(w, "    -input-file-path ./data.json \\")
	fmt.Fprintln(w, "    -output-file-path ./data.csv")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "フラグ:")
	fmt.Fprintln(w, "  -input-file-path   入力ファイルパス (必須)")
	fmt.Fprintln(w, "  -output-file-path  出力ファイルパス (必須)")
	fmt.Fprintln(w, "  -input-format      入力形式 (json|yaml|csv|tsv|html|md-ordered-list|md-unordered-list|md-table) 省略時は拡張子推定")
	fmt.Fprintln(w, "  -output-format     出力形式 (json|yaml|csv|tsv|html|md-ordered-list|md-unordered-list|md-table) 省略時は拡張子推定")
	fmt.Fprintln(w, "  -help              ヘルプ表示")
}

func normalizeFormat(format string) string {
	return strings.ToLower(strings.TrimSpace(format))
}

func isSupportedFormat(format string) bool {
	_, ok := supportedFormats[format]
	return ok
}

func inferFormatFromPath(path string) (string, error) {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(path)))
	format, ok := extensionToFormat[ext]
	if !ok {
		return "", fmt.Errorf("拡張子 %s は未対応です", ext)
	}
	return format, nil
}
