package usecases

import (
	"fmt"
	"io"
	"os"
)

// OutputWriterRepository は出力操作のインターフェース
type OutputWriterRepository interface {
	Printf(format string, args ...interface{})
	Println(args ...interface{})
}

// OutputWriter は実際の標準出力を使用する実装
type OutputWriter struct {
	writer io.Writer
}

// NewOutputWriter は新しいOutputWriterを作成
func NewOutputWriter() *OutputWriter {
	return &OutputWriter{
		writer: os.Stdout,
	}
}

// NewOutputWriterWithWriter は指定されたWriterを使用するOutputWriterを作成
func NewOutputWriterWithWriter(writer io.Writer) *OutputWriter {
	return &OutputWriter{
		writer: writer,
	}
}

// Printf はフォーマット付きで出力
func (o *OutputWriter) Printf(format string, args ...interface{}) {
	fmt.Fprintf(o.writer, format, args...)
}

// Println は改行付きで出力
func (o *OutputWriter) Println(args ...interface{}) {
	fmt.Fprintln(o.writer, args...)
}

// MockOutputWriter はテスト用のモック実装
type MockOutputWriter struct {
	PrintfFunc  func(format string, args ...interface{})
	PrintlnFunc func(args ...interface{})
	Output      []string // 出力内容を記録するためのスライス
}

// Printf はモック関数を実行
func (m *MockOutputWriter) Printf(format string, args ...interface{}) {
	if m.PrintfFunc != nil {
		m.PrintfFunc(format, args...)
	}
	if m.Output != nil {
		m.Output = append(m.Output, fmt.Sprintf(format, args...))
	}
}

// Println はモック関数を実行
func (m *MockOutputWriter) Println(args ...interface{}) {
	if m.PrintlnFunc != nil {
		m.PrintlnFunc(args...)
	}
	if m.Output != nil {
		m.Output = append(m.Output, fmt.Sprintln(args...))
	}
}
