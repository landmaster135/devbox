package flag_parser

import (
	"errors"
	"testing"
)

func TestNewMockFlagParser(t *testing.T) {
	t.Parallel()

	parser := NewMockFlagParser()
	if parser == nil {
		t.Fatalf("NewMockFlagParser() returned nil")
	}
	if parser.stringValues == nil {
		t.Fatalf("stringValues was nil")
	}
	if parser.intValues == nil {
		t.Fatalf("intValues was nil")
	}
	if parser.boolValues == nil {
		t.Fatalf("boolValues was nil")
	}
}

func TestMockFlagParser_StringVar(t *testing.T) {
	t.Parallel()

	t.Run("set value has priority", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetString("operation", "grep-str")

		var got string
		parser.StringVar(&got, "operation", "default", "")

		if got != "grep-str" {
			t.Fatalf("got = %q, want %q", got, "grep-str")
		}
	})

	t.Run("default value is used", func(t *testing.T) {
		parser := NewMockFlagParser()

		var got string
		parser.StringVar(&got, "operation", "default", "")

		if got != "default" {
			t.Fatalf("got = %q, want %q", got, "default")
		}
	})

	t.Run("existing non-empty value is kept", func(t *testing.T) {
		parser := NewMockFlagParser()

		got := "already-set"
		parser.StringVar(&got, "operation", "default", "")

		if got != "already-set" {
			t.Fatalf("got = %q, want %q", got, "already-set")
		}
	})
}

func TestMockFlagParser_IntVar(t *testing.T) {
	t.Parallel()

	t.Run("set value has priority", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetInt("threshold", 120)

		var got int
		parser.IntVar(&got, "threshold", -1, "")

		if got != 120 {
			t.Fatalf("got = %d, want %d", got, 120)
		}
	})

	t.Run("default value is used", func(t *testing.T) {
		parser := NewMockFlagParser()

		var got int
		parser.IntVar(&got, "threshold", -1, "")

		if got != -1 {
			t.Fatalf("got = %d, want %d", got, -1)
		}
	})

	t.Run("existing non-zero value is kept", func(t *testing.T) {
		parser := NewMockFlagParser()

		got := 99
		parser.IntVar(&got, "threshold", -1, "")

		if got != 99 {
			t.Fatalf("got = %d, want %d", got, 99)
		}
	})
}

func TestMockFlagParser_BoolVar(t *testing.T) {
	t.Parallel()

	t.Run("set value has priority", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetBool("help", true)

		got := false
		parser.BoolVar(&got, "help", false, "")

		if !got {
			t.Fatalf("got = false, want true")
		}
	})

	t.Run("default value is used", func(t *testing.T) {
		parser := NewMockFlagParser()

		got := false
		parser.BoolVar(&got, "help", true, "")

		if !got {
			t.Fatalf("got = false, want true")
		}
	})

	t.Run("existing true value is kept", func(t *testing.T) {
		parser := NewMockFlagParser()

		got := true
		parser.BoolVar(&got, "help", false, "")

		if !got {
			t.Fatalf("got = false, want true")
		}
	})
}

func TestMockFlagParser_Parse(t *testing.T) {
	t.Parallel()

	parser := NewMockFlagParser()
	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	wantErr := errors.New("parse failed")
	parser.SetParseError(wantErr)
	if err := parser.Parse(); err != wantErr {
		t.Fatalf("Parse() error = %v, want %v", err, wantErr)
	}
}
