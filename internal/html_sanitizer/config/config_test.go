package config

import (
	"strings"
	"testing"
)

func TestParseFlags_Success(t *testing.T) {
	t.Parallel()

	cfg, _, err := ParseFlags([]string{"--input-file", "input.html", "--output-file", "out.html"})
	if err != nil {
		t.Fatalf("ParseFlags returned error: %v", err)
	}

	if cfg.InputPath != "input.html" {
		t.Fatalf("expected InputPath=input.html, got %q", cfg.InputPath)
	}
	if cfg.OutputPath != "out.html" {
		t.Fatalf("expected OutputPath=out.html, got %q", cfg.OutputPath)
	}
}

func TestParseFlags_HtmlFileAlias(t *testing.T) {
	t.Parallel()

	cfg, _, err := ParseFlags([]string{"--html-file", "src.html", "--output-file", "dst.html"})
	if err != nil {
		t.Fatalf("ParseFlags returned error: %v", err)
	}

	if cfg.InputPath != "src.html" {
		t.Fatalf("expected InputPath=src.html, got %q", cfg.InputPath)
	}
}

func TestParseFlags_HelpSkipsValidation(t *testing.T) {
	t.Parallel()

	cfg, _, err := ParseFlags([]string{"--help"})
	if err != nil {
		t.Fatalf("ParseFlags returned error: %v", err)
	}

	if !cfg.ShowHelp {
		t.Fatalf("expected ShowHelp=true when --help is provided")
	}
}

func TestParseFlags_MissingInput(t *testing.T) {
	t.Parallel()

	_, _, err := ParseFlags([]string{"--output-file", "out.html"})
	if err == nil {
		t.Fatalf("expected error when input flag missing")
	}

	if !strings.Contains(err.Error(), "--input-file") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseFlags_MissingOutput(t *testing.T) {
	t.Parallel()

	_, _, err := ParseFlags([]string{"--input-file", "input.html"})
	if err == nil {
		t.Fatalf("expected error when output flag missing")
	}

	if !strings.Contains(err.Error(), "--output-file") {
		t.Fatalf("unexpected error: %v", err)
	}
}
