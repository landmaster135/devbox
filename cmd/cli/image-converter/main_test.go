package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_MissingOutputDirFlag_ReturnsError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-src-dir", ".", "-ext", "png"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("expected exitCodeError, got %d", code)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("expected empty stdout, got %q", got)
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "エラー: -output-dir は必須です") {
		t.Fatalf("expected required -output-dir error, got %q", errOutput)
	}
	if !strings.Contains(errOutput, "-output-dir string") {
		t.Fatalf("expected usage output to include -output-dir flag, got %q", errOutput)
	}
}

func TestRun_EmptyOutputDirFlag_ReturnsError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-src-dir", ".", "-output-dir", " ", "-ext", "png"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("expected exitCodeError, got %d", code)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("expected empty stdout, got %q", got)
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "エラー: -output-dir は必須です") {
		t.Fatalf("expected required -output-dir error, got %q", errOutput)
	}
}
