package config

import (
	"errors"
	"strings"
	"testing"

	flagParser "github.com/landmaster135/devbox/internal/movie_extractor/infrastructures/flag_parser"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		srcFile     string
		fps         int
		quality     int
		startPos    string
		outDir      string
		expectError string
	}{
		{
			name:      "正常系",
			operation: "extract-frames",
			srcFile:   "input.mp4",
			fps:       3,
			quality:   2,
			startPos:  "00:00:10.5",
			outDir:    "frames",
		},
		{
			name:        "未対応operation",
			operation:   "unknown",
			srcFile:     "input.mp4",
			fps:         1,
			quality:     2,
			outDir:      "frames",
			expectError: "未対応のoperationです",
		},
		{
			name:        "src-file未指定",
			operation:   "extract-frames",
			fps:         1,
			quality:     2,
			outDir:      "frames",
			expectError: "src-file は必須です",
		},
		{
			name:        "fps不正",
			operation:   "extract-frames",
			srcFile:     "input.mp4",
			fps:         0,
			quality:     2,
			outDir:      "frames",
			expectError: "fps は1以上",
		},
		{
			name:        "quality不正",
			operation:   "extract-frames",
			srcFile:     "input.mp4",
			fps:         1,
			quality:     32,
			outDir:      "frames",
			expectError: "quality は1から31",
		},
		{
			name:        "start-position不正",
			operation:   "extract-frames",
			srcFile:     "input.mp4",
			fps:         1,
			quality:     2,
			startPos:    "99:99:99",
			outDir:      "frames",
			expectError: "start-position の形式が不正",
		},
		{
			name:        "out-dir未指定",
			operation:   "extract-frames",
			srcFile:     "input.mp4",
			fps:         1,
			quality:     2,
			expectError: "out-dir は必須です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConfig(tt.operation, tt.srcFile, tt.fps, tt.quality, tt.startPos, tt.outDir)

			if tt.expectError == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if cfg.Operation != tt.operation {
					t.Fatalf("operation mismatch: got=%s", cfg.Operation)
				}
				return
			}

			if err == nil {
				t.Fatalf("expected error: %s", tt.expectError)
			}
			if !strings.Contains(err.Error(), tt.expectError) {
				t.Fatalf("error mismatch: got=%v wantContains=%s", err, tt.expectError)
			}
		})
	}
}

func TestParseFlagsWithParser(t *testing.T) {
	t.Run("help指定", func(t *testing.T) {
		mock := flagParser.NewMockFlagParser()
		mock.SetBoolValue("help", true)

		cfg, err := ParseFlagsWithParser(mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !cfg.Help {
			t.Fatal("help should be true")
		}
	})

	t.Run("正常系", func(t *testing.T) {
		mock := flagParser.NewMockFlagParser()
		mock.SetStringValue("operation", "extract-frames")
		mock.SetStringValue("src-file", "input.mp4")
		mock.SetIntValue("fps", 4)
		mock.SetIntValue("quality", 3)
		mock.SetStringValue("start-position", "12.5")
		mock.SetStringValue("out-dir", "frames")

		cfg, err := ParseFlagsWithParser(mock)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != "extract-frames" || cfg.SrcFile != "input.mp4" {
			t.Fatalf("unexpected config: %+v", cfg)
		}
	})

	t.Run("parse失敗", func(t *testing.T) {
		mock := flagParser.NewMockFlagParser()
		mock.SetParseError(errors.New("invalid flag"))

		_, err := ParseFlagsWithParser(mock)
		if err == nil {
			t.Fatal("expected parse error")
		}
		if !strings.Contains(err.Error(), "フラグ解析に失敗しました") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestIsValidStartPosition(t *testing.T) {
	valid := []string{"", "0", "12.5", "00:00", "00:00:00", "01:02:03.9", "10:59"}
	for _, v := range valid {
		if !isValidStartPosition(v) {
			t.Fatalf("expected valid start-position: %s", v)
		}
	}

	invalid := []string{"abc", "1:2", "00:61", "00:00:60", "00:00:00.x"}
	for _, v := range invalid {
		if isValidStartPosition(v) {
			t.Fatalf("expected invalid start-position: %s", v)
		}
	}
}
