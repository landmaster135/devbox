package config

import "flag"

// FlagParser はフラグ解析のインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse()
}

// DefaultFlagParser は標準のflagパッケージを使用するFlagParserの実装
type DefaultFlagParser struct{}

// NewDefaultFlagParser は新しいDefaultFlagParserを作成する
func NewDefaultFlagParser() *DefaultFlagParser {
	return &DefaultFlagParser{}
}

// StringVar は文字列フラグを定義する
func (p *DefaultFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	flag.StringVar(ptr, name, value, usage)
}

// BoolVar はブールフラグを定義する
func (p *DefaultFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	flag.BoolVar(ptr, name, value, usage)
}

// Parse はフラグを解析する
func (p *DefaultFlagParser) Parse() {
	flag.Parse()
}
