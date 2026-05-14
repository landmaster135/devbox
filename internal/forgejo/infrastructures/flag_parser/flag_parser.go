package flag_parser

import (
	"flag"
	"os"
	"strings"
)

// FlagParser はフラグ解析を抽象化するインターフェースです。
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

type flagKind int

const (
	flagKindUnknown flagKind = iota
	flagKindBool
	flagKindString
)

// StandardFlagParser は標準入力引数からフラグを解析する実装です。
type StandardFlagParser struct {
	flagSet *flag.FlagSet
	args    []string
	kinds   map[string]flagKind
}

// NewStandardFlagParser は標準入力引数からパーサーを作成します。
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{
		flagSet: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
		args:    os.Args[1:],
		kinds:   map[string]flagKind{},
	}
}

// StringVar は文字列フラグを定義します。
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	p.flagSet.StringVar(ptr, name, value, usage)
	if p.kinds == nil {
		p.kinds = map[string]flagKind{}
	}
	p.kinds[name] = flagKindString
}

// BoolVar はブールフラグを定義します。
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	p.flagSet.BoolVar(ptr, name, value, usage)
	if p.kinds == nil {
		p.kinds = map[string]flagKind{}
	}
	p.kinds[name] = flagKindBool
}

// Parse はフラグ解析を実行します。
func (p *StandardFlagParser) Parse() error {
	return p.flagSet.Parse(p.reorderArgs(p.args))
}

// Args は未処理の位置引数を返します。
func (p *StandardFlagParser) Args() []string {
	return p.flagSet.Args()
}

func (p *StandardFlagParser) reorderArgs(rawArgs []string) []string {
	flags := make([]string, 0, len(rawArgs))
	positionals := make([]string, 0, len(rawArgs))
	for i := 0; i < len(rawArgs); i++ {
		arg := rawArgs[i]
		if arg == "--" {
			positionals = append(positionals, rawArgs[i:]...)
			break
		}

		if strings.HasPrefix(arg, "-") {
			name := strings.TrimLeft(arg, "-")
			valueIndex := strings.Index(name, "=")
			if valueIndex >= 0 {
				name = name[:valueIndex]
			}

			flags = append(flags, arg)
			if p.kinds[name] == flagKindString && i+1 < len(rawArgs) && !strings.HasPrefix(rawArgs[i+1], "-") {
				i++
				flags = append(flags, rawArgs[i])
			}
			continue
		}

		positionals = append(positionals, arg)
	}

	return append(flags, positionals...)
}
