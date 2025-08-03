package config

import "os"

// FlagParser はコマンドラインフラグを解析するためのインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	IntVar(p *int, name string, value int, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// FileReader はファイル読み取りのためのインターフェース
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

// StandardFileReader は標準のos.ReadFileを使用する実装
type StandardFileReader struct{}

func (r *StandardFileReader) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
