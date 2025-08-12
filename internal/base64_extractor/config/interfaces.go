package config

import "os"

// FlagParser はフラグ解析のインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
}

// OSArgs はOS引数のインターフェース
type OSArgs interface {
	Args() []string
}

// StandardFlagParser は標準のフラグパーサー実装
type StandardFlagParser struct {
	flagSet map[string]interface{}
	parsed  bool
}

// NewStandardFlagParser は新しいStandardFlagParserを作成する
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{
		flagSet: make(map[string]interface{}),
	}
}

// StringVar は文字列フラグを設定する
func (fp *StandardFlagParser) StringVar(p *string, name string, value string, usage string) {
	*p = value
	fp.flagSet[name] = p
}

// BoolVar はブールフラグを設定する
func (fp *StandardFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	*p = value
	fp.flagSet[name] = p
}

// Parse はフラグを解析する
func (fp *StandardFlagParser) Parse() error {
	// 実際のフラグ解析はflagパッケージを使用
	// ここでは簡略化
	fp.parsed = true
	return nil
}

// StandardOSArgs は標準のOS引数実装
type StandardOSArgs struct{}

// NewStandardOSArgs は新しいStandardOSArgsを作成する
func NewStandardOSArgs() *StandardOSArgs {
	return &StandardOSArgs{}
}

// Args はOS引数を返す
func (oa *StandardOSArgs) Args() []string {
	return os.Args
}
