package config

import (
	"flag"
	"os"
)

// FlagParser はフラグ解析のインターフェースです（テスト用）
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser は標準のflagパッケージを使用する実装です
type StandardFlagParser struct {
	flagSet *flag.FlagSet
}

// NewStandardFlagParser は新しいStandardFlagParserを作成します
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}
}

// StringVar は文字列フラグを定義します
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	p.flagSet.StringVar(ptr, name, value, usage)
}

// BoolVar はブールフラグを定義します
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	p.flagSet.BoolVar(ptr, name, value, usage)
}

// Parse はフラグを解析します
func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(os.Args[1:])
}

// Args は解析後の残りの引数を返します
func (p *StandardFlagParser) Args() []string {
	return p.flagSet.Args()
}
