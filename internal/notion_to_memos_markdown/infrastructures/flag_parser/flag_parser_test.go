package flag_parser

import (
	"os"
	"testing"
)

func TestStandardFlagParser_Parse(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{
		"notion-to-memos-markdown",
		"--operation=grep-str",
		"--src-body-dir=/tmp/body",
		"--target-str=TODO",
		"--help",
		"--threshold=10",
	}

	parser := NewStandardFlagParser()

	var operation string
	var srcBodyDir string
	var targetStr string
	var help bool
	var threshold int

	parser.StringVar(&operation, "operation", "", "")
	parser.StringVar(&srcBodyDir, "src-body-dir", "", "")
	parser.StringVar(&targetStr, "target-str", "", "")
	parser.BoolVar(&help, "help", false, "")
	parser.IntVar(&threshold, "threshold", -1, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if operation != "grep-str" {
		t.Fatalf("operation = %q, want %q", operation, "grep-str")
	}
	if srcBodyDir != "/tmp/body" {
		t.Fatalf("srcBodyDir = %q, want %q", srcBodyDir, "/tmp/body")
	}
	if targetStr != "TODO" {
		t.Fatalf("targetStr = %q, want %q", targetStr, "TODO")
	}
	if !help {
		t.Fatalf("help = false, want true")
	}
	if threshold != 10 {
		t.Fatalf("threshold = %d, want 10", threshold)
	}
}

func TestStandardFlagParser_ParseError(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{"notion-to-memos-markdown", "--unknown-flag"}

	parser := NewStandardFlagParser()
	if err := parser.Parse(); err == nil {
		t.Fatalf("Parse() error = nil, want non-nil")
	}
}
