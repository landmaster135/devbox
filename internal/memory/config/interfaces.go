package config

import "os"

// FlagParser はフラグ解析のインターフェース
type FlagParser interface {
	StringVar(p *string, name string, value string, usage string)
	BoolVar(p *bool, name string, value bool, usage string)
	Parse() error
	Args() []string
}

// FileReader はファイル読み込み機能を提供するインターフェース
type FileReader interface {
	ReadFile(filename string) ([]byte, error)
}

// FileWriter はファイル書き込み機能を提供するインターフェース
type FileWriter interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// StandardFileReader は標準のファイル読み込み実装
type StandardFileReader struct{}

// ReadFile はファイルを読み込んでバイト配列を返す
func (r *StandardFileReader) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// StandardFileWriter は標準のファイル書き込み実装
type StandardFileWriter struct{}

// WriteFile はファイルにバイト配列を書き込む
func (w *StandardFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}
