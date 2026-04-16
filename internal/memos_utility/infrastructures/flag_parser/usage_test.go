package flag_parser

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintUsage_Normal(t *testing.T) {
	oldStderr := os.Stderr
	oldArgs := os.Args
	defer func() {
		os.Stderr = oldStderr
		os.Args = oldArgs
	}()

	os.Args = []string{"memos-utility"}
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error = %v", err)
	}

	os.Stderr = writer
	PrintUsage("使用方法: %[1]s [オプション]\n例: %[1]s -help\n")

	if err := writer.Close(); err != nil {
		t.Fatalf("writer.Close() error = %v", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}

	out := string(data)
	if !strings.Contains(out, "使用方法: memos-utility [オプション]") {
		t.Fatalf("usage output does not include usage header: %s", out)
	}
	if !strings.Contains(out, "例: memos-utility -help") {
		t.Fatalf("usage output does not include example: %s", out)
	}
}
