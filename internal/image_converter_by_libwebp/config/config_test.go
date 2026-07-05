package config

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	flagParser "github.com/landmaster135/devbox/internal/image_converter_by_libwebp/infrastructures/flag_parser"
)

func TestConfig_NewConfig_Normal(t *testing.T) {
	cfg, err := NewConfig("src", "out", "archive", ".WEBP", true, 99, 2, true, true)
	if err != nil {
		t.Fatalf("NewConfig() error = %v", err)
	}

	if cfg.SrcDir != "src" {
		t.Fatalf("SrcDir = %q, want src", cfg.SrcDir)
	}
	if cfg.OutDir != "out" {
		t.Fatalf("OutDir = %q, want out", cfg.OutDir)
	}
	if cfg.ArchiveDir != "archive" {
		t.Fatalf("ArchiveDir = %q, want archive", cfg.ArchiveDir)
	}
	if !cfg.Move {
		t.Fatalf("Move = false, want true")
	}
	if cfg.OutExt != "webp" {
		t.Fatalf("OutExt = %q, want webp", cfg.OutExt)
	}
	if cfg.Quality != 99 {
		t.Fatalf("Quality = %d, want 99", cfg.Quality)
	}
	if cfg.Workers != 2 {
		t.Fatalf("Workers = %d, want 2", cfg.Workers)
	}
	if !cfg.Recursive {
		t.Fatalf("Recursive = false, want true")
	}
	if !cfg.Lossless {
		t.Fatalf("Lossless = false, want true")
	}
}

func TestConfig_NewConfig_Error(t *testing.T) {
	tests := []struct {
		name    string
		srcDir  string
		outDir  string
		outExt  string
		quality int
		workers int
		want    string
	}{
		{name: "入力ディレクトリが空", srcDir: "", outDir: "out", outExt: "webp", quality: 99, workers: 1, want: "入力ディレクトリが指定されていません"},
		{name: "出力ディレクトリが空", srcDir: "src", outDir: "", outExt: "webp", quality: 99, workers: 1, want: "出力ディレクトリが指定されていません"},
		{name: "出力形式が空", srcDir: "src", outDir: "out", outExt: "", quality: 99, workers: 1, want: "出力形式が指定されていません"},
		{name: "未対応形式", srcDir: "src", outDir: "out", outExt: "jpg", quality: 99, workers: 1, want: "サポートされていない出力フォーマットです: jpg"},
		{name: "品質が小さい", srcDir: "src", outDir: "out", outExt: "webp", quality: 0, workers: 1, want: "品質は1から100の範囲で指定してください: 0"},
		{name: "品質が大きい", srcDir: "src", outDir: "out", outExt: "webp", quality: 101, workers: 1, want: "品質は1から100の範囲で指定してください: 101"},
		{name: "workersが小さい", srcDir: "src", outDir: "out", outExt: "webp", quality: 99, workers: 0, want: "workers は1以上で指定してください: 0"},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewConfig(tc.srcDir, tc.outDir, "", tc.outExt, false, tc.quality, tc.workers, false, false)
			if err == nil {
				t.Fatalf("NewConfig() error = nil, want %q", tc.want)
			}
			if err.Error() != tc.want {
				t.Fatalf("NewConfig() error = %q, want %q", err.Error(), tc.want)
			}
		})
	}
}

func TestConfig_ParseFlagsWithParser_Normal(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetStringFlag("src-dir", "src")
	parser.SetStringFlag("out-dir", "out")
	parser.SetStringFlag("archive-dir", "archive")
	parser.SetBoolFlag("move", true)
	parser.SetStringFlag("ext", "webp")
	parser.SetIntFlag("q", 88)
	parser.SetIntFlag("workers", 3)
	parser.SetBoolFlag("recursive", true)
	parser.SetBoolFlag("lossless", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.SrcDir != "src" || cfg.OutDir != "out" || cfg.ArchiveDir != "archive" {
		t.Fatalf("Config directories = %#v", cfg)
	}
	if !cfg.Move || !cfg.Recursive || !cfg.Lossless {
		t.Fatalf("Config bools = %#v", cfg)
	}
	if cfg.Quality != 88 || cfg.Workers != 3 {
		t.Fatalf("Config numeric values = %#v", cfg)
	}
}

func TestConfig_ParseFlagsWithParser_Default_Normal(t *testing.T) {
	cfg, err := ParseFlagsWithParser(flagParser.NewMockFlagParser())
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if cfg.SrcDir != defaultSrcDir {
		t.Fatalf("SrcDir = %q, want %q", cfg.SrcDir, defaultSrcDir)
	}
	if cfg.OutDir != defaultOutDir {
		t.Fatalf("OutDir = %q, want %q", cfg.OutDir, defaultOutDir)
	}
	if cfg.OutExt != defaultOutExt {
		t.Fatalf("OutExt = %q, want %q", cfg.OutExt, defaultOutExt)
	}
	if cfg.Quality != 99 {
		t.Fatalf("Quality = %d, want 99", cfg.Quality)
	}
	if cfg.Workers != runtime.NumCPU() {
		t.Fatalf("Workers = %d, want %d", cfg.Workers, runtime.NumCPU())
	}
}

func TestConfig_ParseFlagsWithParser_Help_Normal(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetBoolFlag("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v", err)
	}
	if !cfg.Help {
		t.Fatalf("Help = false, want true")
	}
}

func TestConfig_ParseFlagsWithParser_ParseError(t *testing.T) {
	parser := flagParser.NewMockFlagParser()
	parser.SetParseError(errors.New("parse failed"))

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatalf("ParseFlagsWithParser() error = nil")
	}
	if !strings.Contains(err.Error(), "フラグの解析に失敗しました") {
		t.Fatalf("ParseFlagsWithParser() error = %q", err.Error())
	}
}

func TestConfig_ParseFlagsWithArgs_Error(t *testing.T) {
	_, err := ParseFlagsWithArgs([]string{"--src", "photos"})
	if err == nil {
		t.Fatalf("ParseFlagsWithArgs() error = nil")
	}
	if !strings.Contains(err.Error(), "flag provided but not defined") {
		t.Fatalf("ParseFlagsWithArgs() error = %q", err.Error())
	}
}
