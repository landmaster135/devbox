package config

import (
	"flag"
	"os"
)

// StandardFlagParser は標準のflagパッケージを使用したFlagParser実装。
type StandardFlagParser struct {
	flagSet *flag.FlagSet
}

// NewStandardFlagParser は新しいStandardFlagParserを作成する。
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{flagSet: flag.CommandLine}
}

// StringVar は文字列フラグを定義する。
func (s *StandardFlagParser) StringVar(p *string, name string, value string, usage string) {
	s.flagSet.StringVar(p, name, value, usage)
}

// IntVar は整数フラグを定義する。
func (s *StandardFlagParser) IntVar(p *int, name string, value int, usage string) {
	s.flagSet.IntVar(p, name, value, usage)
}

// Parse はフラグを解析する。
func (s *StandardFlagParser) Parse() error {
	return s.flagSet.Parse(os.Args[1:])
}
