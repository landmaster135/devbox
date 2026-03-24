package flag_parser

import "testing"

func TestStandardFlagParser_StringVarAndBoolVar(t *testing.T) {
	t.Parallel()

	parser := &StandardFlagParser{
		args: []string{"--operation", "compress", "-p", "/tmp/input", "--help"},
	}

	operation := ""
	path := ""
	help := false

	parser.StringVar(&operation, "operation", "", "")
	parser.StringVar(&path, "p", "", "")
	parser.BoolVar(&help, "help", false, "")

	if operation != "compress" {
		t.Fatalf("operation = %q, want %q", operation, "compress")
	}
	if path != "/tmp/input" {
		t.Fatalf("path = %q, want %q", path, "/tmp/input")
	}
	if !help {
		t.Fatalf("help = false, want true")
	}
}

func TestStandardFlagParser_Parse(t *testing.T) {
	t.Parallel()

	parser := &StandardFlagParser{args: []string{"--operation", "compress"}}
	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}
}

func TestStandardFlagParser_Args(t *testing.T) {
	t.Parallel()

	parser := &StandardFlagParser{
		args: []string{
			"--operation", "compress",
			"--path", "/tmp/input",
			"compress", "/tmp/from-positional",
			"-h",
		},
	}

	got := parser.Args()
	if len(got) != 2 {
		t.Fatalf("len(Args()) = %d, want 2", len(got))
	}
	if got[0] != "compress" {
		t.Fatalf("Args()[0] = %q, want %q", got[0], "compress")
	}
	if got[1] != "/tmp/from-positional" {
		t.Fatalf("Args()[1] = %q, want %q", got[1], "/tmp/from-positional")
	}
}

func TestStartsWith(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		s      string
		prefix string
		want   bool
	}{
		{name: "true", s: "hello", prefix: "he", want: true},
		{name: "false", s: "hello", prefix: "lo", want: false},
		{name: "empty prefix", s: "hello", prefix: "", want: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := startsWith(tt.s, tt.prefix)
			if got != tt.want {
				t.Fatalf("startsWith(%q, %q) = %v, want %v", tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}
