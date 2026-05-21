package flag_parser

import (
	"errors"
	"testing"
)

func TestMockFlagParser_SetPresetValues(t *testing.T) {
	mock := NewMockFlagParser()
	mock.SetStringValue("src-file", "input.mp4")
	mock.SetIntValue("fps", 8)
	mock.SetBoolValue("help", true)

	var (
		srcFile string
		fps     int
		help    bool
	)

	mock.StringVar(&srcFile, "src-file", "", "src")
	mock.IntVar(&fps, "fps", 1, "fps")
	mock.BoolVar(&help, "help", false, "help")

	if srcFile != "input.mp4" {
		t.Fatalf("src-file mismatch: got=%s", srcFile)
	}
	if fps != 8 {
		t.Fatalf("fps mismatch: got=%d", fps)
	}
	if !help {
		t.Fatal("help should be true")
	}
}

func TestMockFlagParser_Parse(t *testing.T) {
	mock := NewMockFlagParser()
	expectedErr := errors.New("parse failed")
	mock.SetParseError(expectedErr)

	err := mock.Parse()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if !mock.ParseCalled() {
		t.Fatal("parse should be marked as called")
	}
}
