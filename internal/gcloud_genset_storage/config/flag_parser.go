package config

import "flag"

// FlagParser is an interface to abstract flag parsing for testability.
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser implements FlagParser using the standard flag package.
type StandardFlagParser struct{}

// NewStandardFlagParser creates a new StandardFlagParser.
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{}
}

// StringVar delegates to flag.StringVar.
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	flag.StringVar(ptr, name, value, usage)
}

// BoolVar delegates to flag.BoolVar.
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	flag.BoolVar(ptr, name, value, usage)
}

// Parse delegates to flag.Parse.
func (p *StandardFlagParser) Parse() error {
	flag.Parse()
	return nil
}

// Args delegates to flag.Args.
func (p *StandardFlagParser) Args() []string {
	return flag.Args()
}
