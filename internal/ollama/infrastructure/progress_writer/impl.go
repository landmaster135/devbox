package infrastructure

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type pullProgressWriter struct {
	out        io.Writer
	isTTY      bool
	buffer     bytes.Buffer
	lastDigest string
}

func NewPullProgressWriter(out io.Writer) *pullProgressWriter {
	return &pullProgressWriter{
		out:   out,
		isTTY: isTerminalWriter(out),
	}
}

func (w *pullProgressWriter) Write(p []byte) (int, error) {
	total := len(p)
	for len(p) > 0 {
		idx := bytes.IndexByte(p, '\n')
		if idx == -1 {
			w.buffer.Write(p)
			break
		}
		w.buffer.Write(p[:idx])
		if err := w.flushLine(); err != nil {
			return 0, err
		}
		p = p[idx+1:]
	}
	return total, nil
}

func (w *pullProgressWriter) flushLine() error {
	line := w.buffer.String()
	w.buffer.Reset()
	if line == "" {
		return nil
	}
	digest := extractDigestFromLine(line)
	if !w.isTTY || digest == "" || w.lastDigest == "" || digest != w.lastDigest {
		if _, err := fmt.Fprintln(w.out, line); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(w.out, "\033[1A\033[2K%s\n", line); err != nil {
			return err
		}
	}
	if digest != "" {
		w.lastDigest = digest
	} else {
		w.lastDigest = ""
	}
	return nil
}

func extractDigestFromLine(line string) string {
	fields := strings.Fields(line)
	if len(fields) >= 2 && fields[0] == "pulling" {
		return fields[1]
	}
	return ""
}

func isTerminalWriter(out io.Writer) bool {
	file, ok := out.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
