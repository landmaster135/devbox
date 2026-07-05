package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_Help_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"--help"}, &stdout, &stderr)
	if code != exitCodeOK {
		t.Fatalf("run() code = %d, want %d", code, exitCodeOK)
	}
}

func TestRun_ParseError(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := run([]string{"--ext", "jpg"}, &stdout, &stderr)
	if code != exitCodeError {
		t.Fatalf("run() code = %d, want %d", code, exitCodeError)
	}
	if !strings.Contains(stderr.String(), "サポートされていない出力フォーマットです: jpg") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
