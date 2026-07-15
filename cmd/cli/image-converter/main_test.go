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

func TestRun_Help_IncludesArchiveDirFlag(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-h"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("expected exitCodeError, got %d", code)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("expected empty stdout, got %q", got)
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "-archive-dir string") {
		t.Fatalf("expected usage output to include -archive-dir flag, got %q", errOutput)
	}
	if strings.Contains(errOutput, "-archive string") {
		t.Fatalf("expected usage output not to include legacy -archive flag, got %q", errOutput)
	}
}

func TestRun_LegacyArchiveFlag_ReturnsError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"-src-dir", ".", "-output-dir", "./out", "-archive", "./archive"}, &stdout, &stderr)

	if code != exitCodeError {
		t.Fatalf("expected exitCodeError, got %d", code)
	}
	if got := stdout.String(); got != "" {
		t.Fatalf("expected empty stdout, got %q", got)
	}

	errOutput := stderr.String()
	if !strings.Contains(errOutput, "flag provided but not defined: -archive") {
		t.Fatalf("expected undefined legacy -archive error, got %q", errOutput)
	}
}
