package config

import "os"

// FlagParser はflagパッケージのインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
}

// OSArgs はos.Argsのインターフェース
type OSArgs interface {
	Args() []string
}

// StandardOSArgs は標準のos.Argsを使用する実装
type StandardOSArgs struct{}

// NewStandardOSArgs は新しいStandardOSArgsを作成する
func NewStandardOSArgs() *StandardOSArgs {
	return &StandardOSArgs{}
}

// Args はos.Argsを返す
func (s *StandardOSArgs) Args() []string {
	return os.Args
}
