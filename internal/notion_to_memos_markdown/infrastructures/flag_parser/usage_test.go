package flag_parser

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestPrintUsage(t *testing.T) {
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
	PrintUsage("cmd=%[1]s alias=%[1]s")
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close write pipe: %v", err)
	}

	outputBytes, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	output := string(outputBytes)
	if !strings.Contains(output, "cmd=custom-notion-cli alias=custom-notion-cli") {
		t.Fatalf("output = %q, want to contain command interpolation", output)
	}
}
