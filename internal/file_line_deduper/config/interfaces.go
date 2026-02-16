package config

// FlagParser はコマンドラインフラグを解析するためのインターフェースです。
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	IntVar(p *int, name string, value int, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}
