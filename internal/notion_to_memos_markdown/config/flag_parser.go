package config

import (
	"flag"
	"os"
)

type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	IntVar(p *int, name string, value int, usage string)
	Parse() error
}

type StandardFlagParser struct {
	flagSet *flag.FlagSet
}

func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}
}

func (s *StandardFlagParser) StringVar(p *string, name string, value string, usage string) {
	s.flagSet.StringVar(p, name, value, usage)
}

func (s *StandardFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	s.flagSet.BoolVar(p, name, value, usage)
}

func (s *StandardFlagParser) IntVar(p *int, name string, value int, usage string) {
	s.flagSet.IntVar(p, name, value, usage)
}

func (s *StandardFlagParser) Parse() error {
	return s.flagSet.Parse(os.Args[1:])
}
