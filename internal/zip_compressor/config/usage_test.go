package config

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

	os.Args = []string{"custom-zip-cli"}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stderr = w
	PrintUsage()
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close write pipe: %v", err)
	}

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	output := string(out)
	if !strings.Contains(output, "Zip圧縮CLIツール") {
		t.Fatalf("output missing title: %q", output)
	}
	if !strings.Contains(output, "custom-zip-cli -operation compress") {
		t.Fatalf("output missing command interpolation: %q", output)
	}
	if !strings.Contains(output, "-help, -h") {
		t.Fatalf("output missing options: %q", output)
	}
}
