package flag_parser

import "testing"

func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	parser := NewStandardFlagParserWithArgs([]string{
		"--src-dir", "src",
		"--out-dir", "out",
		"--archive-dir", "archive",
		"--move",
		"--ext", "webp",
		"--q", "99",
		"--workers", "2",
		"--recursive",
		"--lossless",
	})

	srcDir := ""
	outDir := ""
	archiveDir := ""
	move := false
	ext := ""
	quality := 0
	workers := 0
	recursive := false
	lossless := false

	parser.StringVar(&srcDir, "src-dir", ".", "")
	parser.StringVar(&outDir, "out-dir", "default-out", "")
	parser.StringVar(&archiveDir, "archive-dir", "", "")
	parser.BoolVar(&move, "move", false, "")
	parser.StringVar(&ext, "ext", "webp", "")
	parser.IntVar(&quality, "q", 80, "")
	parser.IntVar(&workers, "workers", 1, "")
	parser.BoolVar(&recursive, "recursive", false, "")
	parser.BoolVar(&lossless, "lossless", false, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if srcDir != "src" || outDir != "out" || archiveDir != "archive" {
		t.Fatalf("string flags = %q, %q, %q", srcDir, outDir, archiveDir)
	}
	if !move || !recursive || !lossless {
		t.Fatalf("bool flags = %v, %v, %v", move, recursive, lossless)
	}
	if ext != "webp" || quality != 99 || workers != 2 {
		t.Fatalf("other flags = %q, %d, %d", ext, quality, workers)
	}
}

func TestMockFlagParser_Parse_Normal(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetStringFlag("src-dir", "src")
	parser.SetBoolFlag("recursive", true)
	parser.SetIntFlag("q", 99)
	parser.SetArgs([]string{"extra"})

	srcDir := ""
	recursive := false
	quality := 0
	parser.StringVar(&srcDir, "src-dir", ".", "")
	parser.BoolVar(&recursive, "recursive", false, "")
	parser.IntVar(&quality, "q", 80, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if srcDir != "src" {
		t.Fatalf("srcDir = %q, want src", srcDir)
	}
	if !recursive {
		t.Fatalf("recursive = false, want true")
	}
	if quality != 99 {
		t.Fatalf("quality = %d, want 99", quality)
	}
	if len(parser.Args()) != 1 || parser.Args()[0] != "extra" {
		t.Fatalf("Args() = %#v", parser.Args())
	}
}
