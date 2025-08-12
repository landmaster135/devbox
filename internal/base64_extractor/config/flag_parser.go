package config

import (
	"flag"
	"os"
)

// RealFlagParser は実際のflagパッケージを使用するフラグパーサー
type RealFlagParser struct {
	flagSet *flag.FlagSet
}

// NewRealFlagParser は新しいRealFlagParserを作成する
func NewRealFlagParser() *RealFlagParser {
	return &RealFlagParser{
		flagSet: flag.CommandLine,
	}
}

// StringVar は文字列フラグを設定する
func (fp *RealFlagParser) StringVar(p *string, name string, value string, usage string) {
	fp.flagSet.StringVar(p, name, value, usage)
}

// BoolVar はブールフラグを設定する
func (fp *RealFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	fp.flagSet.BoolVar(p, name, value, usage)
}

// Parse はフラグを解析する
func (fp *RealFlagParser) Parse() error {
	return fp.flagSet.Parse(os.Args[1:])
}
