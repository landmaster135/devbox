package config

import (
	"flag"
	"os"
)

// StandardFlagParser は標準のflagパッケージを使用したFlagParserの実装
type StandardFlagParser struct {
	flagSet *flag.FlagSet
}

// NewStandardFlagParser は新しいStandardFlagParserを作成する
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
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

// IntVar は整数フラグを定義する
func (p *StandardFlagParser) IntVar(ptr *int, name string, value int, usage string) {
	p.flagSet.IntVar(ptr, name, value, usage)
}

// Parse はコマンドライン引数を解析する
func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(os.Args[1:])
}

// Args は解析後の残りの引数を返す
func (p *StandardFlagParser) Args() []string {
	return p.flagSet.Args()
}
