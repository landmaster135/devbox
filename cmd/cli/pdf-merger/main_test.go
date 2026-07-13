package main

import (
	"bytes"
	"testing"
)

func TestParseFlags_Recursive(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseFlags([]string{"-dir", "images", "-output-dir", "output", "-recursive"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags()でエラーが発生: %v", err)
	}

	if opts.Dir != "images" {
		t.Errorf("Dir: 期待 %s, 実際 %s", "images", opts.Dir)
	}
	if !opts.Recursive {
		t.Errorf("Recursive: 期待 true, 実際 false")
	}
	if opts.OutputDir != "output" {
		t.Errorf("OutputDir: 期待 %s, 実際 %s", "output", opts.OutputDir)
	}
}

func TestParseFlags_RecursiveDefault(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseFlags([]string{"-dir", "images", "-output-dir", "output"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags()でエラーが発生: %v", err)
	}

	if opts.Recursive {
		t.Errorf("Recursive: 期待 false, 実際 true")
	}
}

func TestParseFlags_OutFlagRemoved(t *testing.T) {
	var stderr bytes.Buffer

	_, err := parseFlags([]string{"-dir", "images", "-out", "output.pdf"}, &stderr)
	if err == nil {
		t.Fatal("廃止された-outフラグではエラーが期待されます")
	}
}

func TestParseFlags_OutputDirRequired(t *testing.T) {
	var stderr bytes.Buffer

	_, err := parseFlags([]string{"-dir", "images"}, &stderr)
	if err == nil {
		t.Fatal("-output-dir未指定ではエラーが期待されます")
	}
}
