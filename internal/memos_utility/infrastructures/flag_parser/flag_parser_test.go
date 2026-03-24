package flag_parser

import (
	"os"
	"testing"
)

func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	parser := NewStandardFlagParser("memos-utility", []string{
		"-operation=create-web-clip",
		"-timeout=45",
		"-help",
		"positional-arg",
	})

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

	if operation != "create-web-clip" {
		t.Fatalf("operation = %q, want %q", operation, "create-web-clip")
	}
	if timeout != 45 {
		t.Fatalf("timeout = %d, want 45", timeout)
	}
	if !help {
		t.Fatal("help = false, want true")
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "positional-arg" {
		t.Fatalf("args = %#v, want [positional-arg]", args)
	}
}

func TestStandardFlagParser_ParseError_Error(t *testing.T) {
	parser := NewStandardFlagParser("memos-utility", []string{"-unknown-flag"})

	if err := parser.Parse(); err == nil {
		t.Fatal("Parse() error = nil, want non-nil")
	}
}

func TestStandardFlagParser_NewStandardFlagParserFromOSArgs_Normal(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{
		"memos-utility",
		"-operation=create-common-memos",
		"leftover",
	}

	parser := NewStandardFlagParserFromOSArgs()

	var operation string
	parser.StringVar(&operation, "operation", "", "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if operation != "create-common-memos" {
		t.Fatalf("operation = %q, want %q", operation, "create-common-memos")
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "leftover" {
		t.Fatalf("args = %#v, want [leftover]", args)
	}
}
