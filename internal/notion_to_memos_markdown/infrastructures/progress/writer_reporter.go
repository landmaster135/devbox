package progress

import (
	"fmt"
	"io"
)

type WriterReporter struct {
	writer io.Writer
}

func NewWriterReporter(writer io.Writer) *WriterReporter {
	return &WriterReporter{
		writer: writer,
	}
}

func (r *WriterReporter) Report(message string) {
	if r == nil || r.writer == nil {
		return
	}
	fmt.Fprintln(r.writer, message)
}
