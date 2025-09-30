package config

import (
	"flag"
	"os"
)

// FlagParser はフラグ解析を抽象化するインターフェース。
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser は flag パッケージを利用した実装。
type StandardFlagParser struct {
	flagSet *flag.FlagSet
}

// NewStandardFlagParser は StandardFlagParser を生成する。
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}
}

// StringVar は文字列フラグを登録する。
func (s *StandardFlagParser) StringVar(p *string, name string, value string, usage string) {
	s.flagSet.StringVar(p, name, value, usage)
}

// BoolVar は真偽値フラグを登録する。
func (s *StandardFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	s.flagSet.BoolVar(p, name, value, usage)
}

// Parse はフラグを解析する。
func (s *StandardFlagParser) Parse() error {
	return s.flagSet.Parse(os.Args[1:])
}

// Args は解析後の引数を返す。
func (s *StandardFlagParser) Args() []string {
	return s.flagSet.Args()
}
