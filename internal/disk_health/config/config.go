package config

import (
	"fmt"

	flagParser "github.com/landmaster135/devbox/internal/disk_health/infrastructures/flag_parser"
)

const (
	OperationAssessSmart = "assess-smart"
)

const usageTemplate = `使用方法: %[1]s -operation=assess-smart -src-file=<SMARTログファイル> [オプション]

オプション:
  -operation string  実行する操作。対応値: assess-smart
  -src-file string   smartctl -a の出力を保存した .log / .txt ファイル (必須)
  -json              JSON形式で出力
  -verbose           判定根拠に使ったSMART属性を詳細表示
  -help              ヘルプを表示
`

type Config struct {
	Operation string
	SrcFile   string
	JSON      bool
	Verbose   bool
	Help      bool
}

func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(flagParser.NewStandardFlagParser())
}

func ParseFlagsWithParser(parser flagParser.FlagParser) (*Config, error) {
	cfg := &Config{}
	parser.StringVar(&cfg.Operation, "operation", "", "実行する操作")
	parser.StringVar(&cfg.SrcFile, "src-file", "", "SMARTログファイル")
	parser.BoolVar(&cfg.JSON, "json", false, "JSON形式で出力")
	parser.BoolVar(&cfg.Verbose, "verbose", false, "詳細表示")
	parser.BoolVar(&cfg.Help, "help", false, "ヘルプを表示")

	if err := parser.Parse(); err != nil {
		return nil, err
	}
	if cfg.Help {
		return cfg, nil
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Operation == "" {
		return fmt.Errorf("--operation は必須パラメータです")
	}
	if c.Operation != OperationAssessSmart {
		return fmt.Errorf("未対応のoperationです: %s", c.Operation)
	}
	if c.SrcFile == "" {
		return fmt.Errorf("--src-file は必須パラメータです")
	}
	return nil
}

func PrintUsage() {
	flagParser.PrintUsage(usageTemplate)
}
