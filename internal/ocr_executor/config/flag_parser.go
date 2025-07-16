package config

import (
	"flag"
	"os"
)

// RealFlagParser は実際のflagパッケージを使用する実装
type RealFlagParser struct {
	flagSet *flag.FlagSet
}

// NewRealFlagParser は新しいRealFlagParserを作成する
func NewRealFlagParser() *RealFlagParser {
	return &RealFlagParser{
		flagSet: flag.CommandLine,
	}
}

// StringVar は文字列フラグを定義する
func (r *RealFlagParser) StringVar(p *string, name string, value string, usage string) {
	r.flagSet.StringVar(p, name, value, usage)
}

// BoolVar はブールフラグを定義する
func (r *RealFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	r.flagSet.BoolVar(p, name, value, usage)
}

// Parse はフラグを解析する
func (r *RealFlagParser) Parse() error {
	r.flagSet.Parse(os.Args[1:])
	return nil
}
