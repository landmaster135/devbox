package config

import "flag"

// FlagParser はフラグ解析のインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	IntVar(p *int, name string, value int, usage string)
	Parse() error
}

// StandardFlagParser は標準のflagパッケージを使用するFlagParser
type StandardFlagParser struct{}

// NewStandardFlagParser は新しいStandardFlagParserを作成する
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{}
}

// StandardFlagParser の実装
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	flag.StringVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	flag.BoolVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) IntVar(ptr *int, name string, value int, usage string) {
	flag.IntVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) Parse() error {
	flag.Parse()
	return nil
}
