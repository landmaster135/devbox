package flag_parser

import (
	"flag"
	"os"
)

type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
}

type StandardFlagParser struct {
	flagSet *flag.FlagSet
}

func NewStandardFlagParser() *StandardFlagParser {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flagSet.SetOutput(os.Stderr)
	return &StandardFlagParser{flagSet: flagSet}
}

func (p *StandardFlagParser) StringVar(value *string, name string, defaultValue string, usage string) {
	p.flagSet.StringVar(value, name, defaultValue, usage)
}

func (p *StandardFlagParser) BoolVar(value *bool, name string, defaultValue bool, usage string) {
	p.flagSet.BoolVar(value, name, defaultValue, usage)
}

func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(os.Args[1:])
}
