package config

import "os"

// FlagParser はコマンドラインフラグの解析を行うインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
}

// OSArgs はOS引数を取得するインターフェース
type OSArgs interface {
	Args() []string
}

// StandardOSArgs は標準のOS引数実装
type StandardOSArgs struct{}

// NewStandardOSArgs は新しいStandardOSArgsを作成する
func NewStandardOSArgs() *StandardOSArgs {
	return &StandardOSArgs{}
}

// Args はOS引数を返す
func (s *StandardOSArgs) Args() []string {
	return os.Args
}
