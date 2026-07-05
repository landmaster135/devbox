package flag_parser

import (
	"flag"
	"os"
)

// FlagParser はフラグ解析を抽象化します。
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	IntVar(p *int, name string, value int, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser は Go 標準 flag パッケージによる実装です。
type StandardFlagParser struct {
	flagSet *flag.FlagSet
	args    []string
}

// NewStandardFlagParser は os.Args から parser を作成します。
func NewStandardFlagParser() *StandardFlagParser {
	return NewStandardFlagParserWithArgs(os.Args[1:])
}

// NewStandardFlagParserWithArgs は指定引数から parser を作成します。
func NewStandardFlagParserWithArgs(args []string) *StandardFlagParser {
	flagSet := flag.NewFlagSet("image-converter-by-libwebp", flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	return &StandardFlagParser{
		flagSet: flagSet,
		args:    args,
	}
}

func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	p.flagSet.StringVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	p.flagSet.BoolVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) IntVar(ptr *int, name string, value int, usage string) {
	p.flagSet.IntVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(p.args)
}

func (p *StandardFlagParser) Args() []string {
	return p.flagSet.Args()
}
