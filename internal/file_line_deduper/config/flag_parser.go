package config

import (
	"flag"
	"io"
	"os"
)

// StandardFlagParser は標準のflagパッケージを使用する実装です。
type StandardFlagParser struct {
	flagSet *flag.FlagSet
	args    []string
}

// NewStandardFlagParser は os.Args を使ってフラグパーサーを作成します。
func NewStandardFlagParser() *StandardFlagParser {
	return newStandardFlagParser(os.Args[1:])
}

func newStandardFlagParser(args []string) *StandardFlagParser {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	return &StandardFlagParser{
		flagSet: flagSet,
		args:    args,
	}
}

// StringVar は文字列フラグを定義します。
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	p.flagSet.StringVar(ptr, name, value, usage)
}

// IntVar は整数フラグを定義します。
func (p *StandardFlagParser) IntVar(ptr *int, name string, value int, usage string) {
	p.flagSet.IntVar(ptr, name, value, usage)
}

// BoolVar はブールフラグを定義します。
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	p.flagSet.BoolVar(ptr, name, value, usage)
}

// Parse はフラグを解析します。
func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(p.args)
}

// Args は解析後の残り引数を返します。
func (p *StandardFlagParser) Args() []string {
	return p.flagSet.Args()
}
