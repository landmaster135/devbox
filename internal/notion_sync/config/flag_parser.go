package config

import (
	"flag"
	"os"
)

// FlagParser はフラグ解析のインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser は標準のflagパッケージを使用するFlagParserの実装
type StandardFlagParser struct {
	flagSet *flag.FlagSet
}

// NewStandardFlagParser は新しいStandardFlagParserを作成する
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{
		flagSet: flag.CommandLine,
	}
}

// StringVar は文字列フラグを定義する
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	p.flagSet.StringVar(ptr, name, value, usage)
}

// BoolVar はブールフラグを定義する
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	p.flagSet.BoolVar(ptr, name, value, usage)
}

// Parse はフラグを解析する
func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(os.Args[1:])
}

// Args は解析後の残りの引数を返す
func (p *StandardFlagParser) Args() []string {
	return p.flagSet.Args()
}
