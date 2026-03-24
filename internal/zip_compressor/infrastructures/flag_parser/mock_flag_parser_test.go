package flag_parser

import (
	"errors"
	"testing"
)

func TestMockFlagParser_StringVarAndBoolVar(t *testing.T) {
	t.Parallel()

	mock := NewMockFlagParser()
	mock.SetStringFlag("operation", "compress")
	mock.SetBoolFlag("help", true)

	operation := ""
	help := false

	mock.StringVar(&operation, "operation", "", "")
	mock.BoolVar(&help, "help", false, "")

	if operation != "compress" {
		t.Fatalf("operation = %q, want %q", operation, "compress")
	}
	if !help {
		t.Fatalf("help = false, want true")
	}
}

func TestMockFlagParser_SetterReflectsRegisteredPointers(t *testing.T) {
	t.Parallel()

	mock := NewMockFlagParser()
	operation := ""
	help := false

	mock.StringVar(&operation, "operation", "", "")
	mock.BoolVar(&help, "help", false, "")

	mock.SetStringFlag("operation", "decompress")
	mock.SetBoolFlag("help", true)

	if operation != "decompress" {
		t.Fatalf("operation = %q, want %q", operation, "decompress")
	}
	if !help {
		t.Fatalf("help = false, want true")
	}
}

func TestMockFlagParser_ArgsAndParse(t *testing.T) {
	t.Parallel()

	mock := NewMockFlagParser()
	mock.SetArgs([]string{"compress", "/tmp/file"})

	args := mock.Args()
	if len(args) != 2 {
		t.Fatalf("len(args) = %d, want 2", len(args))
	}

	if err := mock.Parse(); err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	wantErr := errors.New("parse error")
	mock.SetParseError(wantErr)
	if err := mock.Parse(); err != wantErr {
		t.Fatalf("Parse() error = %v, want %v", err, wantErr)
	}
}
