package config

import (
	"flag"
	"os"
)

// FlagParser は標準flagパッケージの抽象化インターフェース
// テスト時に任意の実装へ差し替えられるようにする
// (Goのflagパッケージはグローバル状態を持つため)
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	IntVar(p *int, name string, value int, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser はflagパッケージに基づく FlagParser 実装
// ContinueOnError を利用して、エラー時に即座に終了しないようにする
type StandardFlagParser struct {
	flagSet *flag.FlagSet
}

// NewStandardFlagParser は標準実装を返す
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}
}

func (s *StandardFlagParser) StringVar(p *string, name string, value string, usage string) {
	s.flagSet.StringVar(p, name, value, usage)
}

func (s *StandardFlagParser) IntVar(p *int, name string, value int, usage string) {
	s.flagSet.IntVar(p, name, value, usage)
}

func (s *StandardFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	s.flagSet.BoolVar(p, name, value, usage)
}

func (s *StandardFlagParser) Parse() error {
	return s.flagSet.Parse(os.Args[1:])
}

func (s *StandardFlagParser) Args() []string {
	return s.flagSet.Args()
}
