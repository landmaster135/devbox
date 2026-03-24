package flag_parser

import (
	"errors"
	"testing"
)

func TestMockFlagParser_ValuesAndParse_Normal(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetString("operation", "create-clips")
	parser.SetInt("timeout", 60)
	parser.SetBool("help", true)
	parser.SetArgs([]string{"left1", "left2"})

	var (
		operation string
		timeout   int
		help      bool
	)

	parser.StringVar(&operation, "operation", "", "")
	parser.IntVar(&timeout, "timeout", 30, "")
	parser.BoolVar(&help, "help", false, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if operation != "create-clips" {
		t.Fatalf("operation = %q, want %q", operation, "create-clips")
	}
	if timeout != 60 {
		t.Fatalf("timeout = %d, want 60", timeout)
	}
	if !help {
		t.Fatal("help = false, want true")
	}

	args := parser.Args()
	if len(args) != 2 || args[0] != "left1" || args[1] != "left2" {
		t.Fatalf("args = %#v, want [left1 left2]", args)
	}
}

func TestMockFlagParser_ParseError_Error(t *testing.T) {
	parser := NewMockFlagParser()
	wantErr := errors.New("parse failed")
	parser.SetParseError(wantErr)

	if err := parser.Parse(); !errors.Is(err, wantErr) {
		t.Fatalf("Parse() error = %v, want %v", err, wantErr)
	}
}

func TestMockFlagParser_DefaultAndKeepExistingValues_Normal(t *testing.T) {
	parser := NewMockFlagParser()

	emptyString := ""
	parser.StringVar(&emptyString, "operation", "default-operation", "")
	if emptyString != "default-operation" {
		t.Fatalf("emptyString = %q, want %q", emptyString, "default-operation")
	}

	existingString := "keep"
	parser.StringVar(&existingString, "operation", "overwrite", "")
	if existingString != "keep" {
		t.Fatalf("existingString = %q, want %q", existingString, "keep")
	}

	emptyInt := 0
	parser.IntVar(&emptyInt, "timeout", 30, "")
	if emptyInt != 30 {
		t.Fatalf("emptyInt = %d, want 30", emptyInt)
	}

	existingInt := 99
	parser.IntVar(&existingInt, "timeout", 30, "")
	if existingInt != 99 {
		t.Fatalf("existingInt = %d, want 99", existingInt)
	}

	emptyBool := false
	parser.BoolVar(&emptyBool, "help", true, "")
	if !emptyBool {
		t.Fatal("emptyBool = false, want true")
	}

	existingBool := true
	parser.BoolVar(&existingBool, "help", false, "")
	if !existingBool {
		t.Fatal("existingBool = false, want true")
	}
}
