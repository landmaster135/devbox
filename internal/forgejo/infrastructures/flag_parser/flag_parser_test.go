package flag_parser

import (
	"os"
	"testing"
)

func TestStandardFlagParser_Parse_WithInterspersedFlags(t *testing.T) {
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()
	os.Args = []string{
		"forgejo.exe",
		"repo",
		"list",
		"-forgejo-host",
		"https://codeberg.org",
		"-forgejo-username",
		"myuser",
		"-json",
	}

	parser := NewStandardFlagParser()
	var host string
	var user string
	var jsonFlag bool
	parser.StringVar(&host, "forgejo-host", "", "")
	parser.StringVar(&user, "forgejo-username", "", "")
	parser.BoolVar(&jsonFlag, "json", false, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if got, want := host, "https://codeberg.org"; got != want {
		t.Fatalf("host = %q, want %q", got, want)
	}
	if got, want := user, "myuser"; got != want {
		t.Fatalf("user = %q, want %q", got, want)
	}
	if !jsonFlag {
		t.Fatalf("json flag should be true")
	}
	args := parser.Args()
	if len(args) != 2 || args[0] != "repo" || args[1] != "list" {
		t.Fatalf("Args() = %q, want [repo list]", args)
	}
}
