package flag_parser

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintUsage(t *testing.T) {
	t.Parallel()

	oldArgs := os.Args
	oldStderr := os.Stderr
	defer func() {
		os.Args = oldArgs
		os.Stderr = oldStderr
	}()

	os.Args = []string{"zip-compressor"}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stderr = w
	PrintUsage("tool=%[1]s")
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close write pipe: %v", err)
	}

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read pipe: %v", err)
	}

	if !strings.Contains(string(out), "tool=zip-compressor") {
		t.Fatalf("output = %q, want to contain tool=zip-compressor", string(out))
	}
}
