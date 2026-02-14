package config

import "testing"

func TestParseFlagsKaTeXTable_Normal(t *testing.T) {
	cfg, err := ParseFlags([]string{
		"-input-file-path", "./input.data",
		"-output-file-path", "./output.data",
		"-input-format", "md-table",
		"-output-format", "katex-table",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.InputFormat != FormatMDTable {
		t.Fatalf("input format mismatch: got=%s want=%s", cfg.InputFormat, FormatMDTable)
	}
	if cfg.OutputFormat != FormatKaTeXTable {
		t.Fatalf("output format mismatch: got=%s want=%s", cfg.OutputFormat, FormatKaTeXTable)
	}
}
