package config

import "testing"

func TestStandardFlagParser_ParseAndArgs_Normal(t *testing.T) {
	parser := newStandardFlagParser([]string{
		"-file", "sample.txt",
		"-start", "2",
		"-end", "8",
		"extra",
	})

	var (
		filePath string
		startPos int
		endPos   int
	)

	parser.StringVar(&filePath, "file", "", "")
	parser.IntVar(&startPos, "start", 0, "")
	parser.IntVar(&endPos, "end", 0, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	if filePath != "sample.txt" {
		t.Errorf("filePath = %q, want %q", filePath, "sample.txt")
	}

	if startPos != 2 {
		t.Errorf("startPos = %d, want %d", startPos, 2)
	}

	if endPos != 8 {
		t.Errorf("endPos = %d, want %d", endPos, 8)
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Errorf("Args() = %v, want [extra]", args)
	}
}

func TestStandardFlagParser_Parse_Error(t *testing.T) {
	parser := newStandardFlagParser([]string{"-start", "invalid"})

	var startPos int
	parser.IntVar(&startPos, "start", 0, "")

	if err := parser.Parse(); err == nil {
		t.Fatal("Parse() error = nil, want error")
	}
}
