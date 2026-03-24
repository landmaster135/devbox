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

	os.Args = []string{"custom-notion-cli"}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	os.Stderr = w
	PrintUsage()
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close write pipe: %v", err)
	}

	outputBytes, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	output := string(outputBytes)
	if !strings.Contains(output, "notion-to-memos-markdown CLI") {
		t.Fatalf("usage output missing title: %q", output)
	}
	if !strings.Contains(output, "custom-notion-cli --operation=distribute-files") {
		t.Fatalf("usage output missing command name interpolation: %q", output)
	}
	if !strings.Contains(output, "--src-resource-dir") {
		t.Fatalf("usage output missing option description: %q", output)
	}
}
