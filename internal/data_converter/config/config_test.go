package config

import "testing"

func TestParseFlags_FormatsInferredFromExtensions(t *testing.T) {
	cfg, err := ParseFlags([]string{
		"-input-file-path", "./input.tsv",
		"-output-file-path", "./output.html",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.InputFormat != FormatTSV {
		t.Fatalf("input format mismatch: got=%s want=%s", cfg.InputFormat, FormatTSV)
	}
	if cfg.OutputFormat != FormatHTML {
		t.Fatalf("output format mismatch: got=%s want=%s", cfg.OutputFormat, FormatHTML)
	}
}

func TestParseFlags_ExplicitFormats(t *testing.T) {
	cfg, err := ParseFlags([]string{
		"-input-file-path", "./input.data",
		"-output-file-path", "./output.data",
		"-input-format", "JSON",
		"-output-format", "tsv",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.InputFormat != FormatJSON {
		t.Fatalf("input format mismatch: got=%s want=%s", cfg.InputFormat, FormatJSON)
	}
	if cfg.OutputFormat != FormatTSV {
		t.Fatalf("output format mismatch: got=%s want=%s", cfg.OutputFormat, FormatTSV)
	}
}

func TestParseFlags_InvalidInputFormat(t *testing.T) {
	_, err := ParseFlags([]string{
		"-input-file-path", "./in.txt",
		"-output-file-path", "./out.csv",
		"-input-format", "txt",
	})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestParseFlags_MissingRequired(t *testing.T) {
	_, err := ParseFlags([]string{"-output-file-path", "./out.csv"})
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestParseFlags_HelpBypassesValidation(t *testing.T) {
	cfg, err := ParseFlags([]string{"-help"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatal("expected help=true")
	}
}
