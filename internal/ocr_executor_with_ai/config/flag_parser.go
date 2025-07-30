package config

// FlagParser はフラグ解析のインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Float64Var(p *float64, name string, value float64, usage string)
	IntVar(p *int, name string, value int, usage string)
	Parse() error
}

// StandardFlagParser は標準のflagパッケージを使用するFlagParser
type StandardFlagParser struct{}

// NewStandardFlagParser は新しいStandardFlagParserを作成する
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{}
}
