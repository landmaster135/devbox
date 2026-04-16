package flag_parser

import (
	"flag"
	"io"
	"os"
)

// FlagParser はフラグ解析を抽象化するインターフェース。
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	IntVar(p *int, name string, value int, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser は flag.FlagSet を使った標準実装。
type StandardFlagParser struct {
	flagSet *flag.FlagSet
	args    []string
}

// NewStandardFlagParser は引数を受け取って parser を作成する。
func NewStandardFlagParser(programName string, args []string) *StandardFlagParser {
	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	return &StandardFlagParser{
		flagSet: flagSet,
		args:    args,
	}
}

// NewStandardFlagParserFromOSArgs は os.Args から parser を作成する。
func NewStandardFlagParserFromOSArgs() *StandardFlagParser {
	return NewStandardFlagParser(os.Args[0], os.Args[1:])
}

// StringVar は文字列フラグを定義する。
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	p.flagSet.StringVar(ptr, name, value, usage)
}

// IntVar は整数フラグを定義する。
func (p *StandardFlagParser) IntVar(ptr *int, name string, value int, usage string) {
	p.flagSet.IntVar(ptr, name, value, usage)
}

// BoolVar は真偽値フラグを定義する。
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	p.flagSet.BoolVar(ptr, name, value, usage)
}

// Parse はフラグ解析を実行する。
func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(p.args)
}

// Args は解析後に残った位置引数を返す。
func (p *StandardFlagParser) Args() []string {
	return p.flagSet.Args()
}
