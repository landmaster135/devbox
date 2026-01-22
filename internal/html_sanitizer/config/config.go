package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const flagSetName = "html-sanitizer"

// Config はHTMLサニタイズCLIの設定を保持します。
type Config struct {
	InputPath     string
	OutputPath    string
	OmitsFullBody bool
	ShowHelp      bool
}

// ParseFlags は引数を解析し、Configと使用したFlagSetを返します。
func ParseFlags(args []string) (*Config, *flag.FlagSet, error) {
	cfg := &Config{}
	var htmlFileAlias string

	fs := flag.NewFlagSet(flagSetName, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&cfg.InputPath, "input-file", "", "サニタイズ対象のHTMLファイルパス")
	fs.StringVar(&htmlFileAlias, "html-file", "", "--input-fileのエイリアス")
	fs.StringVar(&cfg.OutputPath, "output-file", "", "サニタイズ結果を書き出すファイルパス")
	fs.BoolVar(&cfg.OmitsFullBody, "omits-full-body", false, "trueでエラー時に入力HTML全文を出力しません")
	fs.BoolVar(&cfg.ShowHelp, "help", false, "ヘルプを表示")
	fs.Usage = func() {
		PrintUsage(fs)
	}

	if err := fs.Parse(args); err != nil {
		return nil, fs, err
	}

	if cfg.InputPath == "" {
		cfg.InputPath = htmlFileAlias
	}

	if cfg.ShowHelp {
		return cfg, fs, nil
	}

	if cfg.InputPath == "" {
		return nil, fs, errors.New("--input-fileまたは--html-fileは必須です")
	}
	if cfg.OutputPath == "" {
		return nil, fs, errors.New("--output-fileは必須です")
	}

	return cfg, fs, nil
}

// PrintUsage はCLIの使用方法を表示します。
func PrintUsage(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.NewFlagSet(flagSetName, flag.ContinueOnError)
		fs.SetOutput(os.Stderr)
	}

	fmt.Fprintf(fs.Output(), "Usage: %s (--input-file|--html-file) <path> --output-file <path> [--omits-full-body]\n", fs.Name())
	fmt.Fprintln(fs.Output(), "Options:")
	fs.PrintDefaults()
}
