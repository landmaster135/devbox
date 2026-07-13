package main

import (
	"bytes"
	"testing"
)

func TestParseFlags_Recursive(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseFlags([]string{"-dir", "images", "-recursive"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags()でエラーが発生: %v", err)
	}

	if opts.Dir != "images" {
		t.Errorf("Dir: 期待 %s, 実際 %s", "images", opts.Dir)
	}
	if !opts.Recursive {
		t.Errorf("Recursive: 期待 true, 実際 false")
	}
}

func TestParseFlags_RecursiveDefault(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseFlags([]string{"-dir", "images"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags()でエラーが発生: %v", err)
	}

	if opts.Recursive {
		t.Errorf("Recursive: 期待 false, 実際 true")
	}
}
