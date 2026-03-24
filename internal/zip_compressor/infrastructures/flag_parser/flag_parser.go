package flag_parser

import "os"

// FlagParser はフラグ解析を抽象化するインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// StandardFlagParser はコマンドライン引数を解析する実装
type StandardFlagParser struct {
	args []string
}

// NewStandardFlagParser は標準入力引数からパーサーを作成する
func NewStandardFlagParser() *StandardFlagParser {
	return &StandardFlagParser{args: os.Args[1:]}
}

// StringVar は文字列フラグを定義する
func (p *StandardFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	for i, arg := range p.args {
		if arg == "-"+name || arg == "--"+name {
			if i+1 < len(p.args) {
				*ptr = p.args[i+1]
			}
		}
	}
}

// BoolVar はブールフラグを定義する
func (p *StandardFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	for _, arg := range p.args {
		if arg == "-"+name || arg == "--"+name {
			*ptr = true
		}
	}
}

// Parse はフラグ解析を実行する
func (p *StandardFlagParser) Parse() error {
	return nil
}

// Args は未処理の位置引数を返す
func (p *StandardFlagParser) Args() []string {
	var remainingArgs []string
	skipNext := false

	for _, arg := range p.args {
		if skipNext {
			skipNext = false
			continue
		}

		if arg == "-operation" || arg == "--operation" ||
			arg == "-o" || arg == "--o" ||
			arg == "-path" || arg == "--path" ||
			arg == "-p" || arg == "--p" {
			skipNext = true
			continue
		}

		if arg == "-help" || arg == "--help" || arg == "-h" || arg == "--h" {
			continue
		}

		if !startsWith(arg, "-") {
			remainingArgs = append(remainingArgs, arg)
		}
	}

	return remainingArgs
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
