package flag_parser

import (
	"flag"
	"io"
	"os"
)

// FlagParser はフラグ解析を抽象化するインターフェースです。
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	IntVar(p *int, name string, value int, usage string)
	Float64Var(p *float64, name string, value float64, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
}

// StandardFlagParser は標準の flag パッケージを利用する実装です。
type StandardFlagParser struct {
	flagSet *flag.FlagSet
	args    []string
}

// NewStandardFlagParser は新しい StandardFlagParser を生成します。
func NewStandardFlagParser() *StandardFlagParser {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	return &StandardFlagParser{
		flagSet: flagSet,
		args:    os.Args[1:],
	}
}

func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	p.flagSet.StringVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) IntVar(ptr *int, name string, value int, usage string) {
	p.flagSet.IntVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) Float64Var(ptr *float64, name string, value float64, usage string) {
	p.flagSet.Float64Var(ptr, name, value, usage)
}

func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	p.flagSet.BoolVar(ptr, name, value, usage)
}

func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(p.args)
}
