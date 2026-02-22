package config

import "os"

// FlagParser はCLIフラグの解析を抽象化する。
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	IntVar(p *int, name string, value int, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// OSArgs はOS引数取得を抽象化する。
type OSArgs interface {
	Args() []string
}

// StandardOSArgs は標準のOS引数実装。
type StandardOSArgs struct{}

// NewStandardOSArgs はStandardOSArgsを返す。
func NewStandardOSArgs() *StandardOSArgs {
	return &StandardOSArgs{}
}

// Args はプロセス引数を返す。
func (s *StandardOSArgs) Args() []string {
	return os.Args
}
