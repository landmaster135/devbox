package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_MissingOutFlag_ReturnsError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-src", ".", "-ext", "png"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("expected exitCodeError, got %d", code)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("expected empty stdout, got %q", got)
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "エラー: -out は必須です") {
		t.Fatalf("expected required -out error, got %q", errOutput)
	}
	if !strings.Contains(errOutput, "-out string") {
		t.Fatalf("expected usage output to include -out flag, got %q", errOutput)
	}
}

func TestRun_EmptyOutFlag_ReturnsError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-src", ".", "-out", " ", "-ext", "png"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("expected exitCodeError, got %d", code)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("expected empty stdout, got %q", got)
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "エラー: -out は必須です") {
		t.Fatalf("expected required -out error, got %q", errOutput)
	}
}
