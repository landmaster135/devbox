package config

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

type mockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string
	boolValues   map[string]bool
	parseErr     error
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
	}
}

func (m *mockFlagParser) StringVar(p *string, name, value, usage string) {
	if preset, ok := m.stringValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if preset, ok := m.boolValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

func (m *mockFlagParser) Parse() error   { return m.parseErr }
func (m *mockFlagParser) Args() []string { return nil }

func (m *mockFlagParser) setString(name, value string) {
	m.stringValues[name] = value
	if ptr, ok := m.stringVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setBool(name string, value bool) {
	m.boolValues[name] = value
	if ptr, ok := m.boolVars[name]; ok {
		*ptr = value
	}
}

func TestParseFlagsWithParser_Success(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationAuthLogin)
	parser.setString("project-id", " my-project ")
	parser.setString("additional-args", " --quiet ")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationAuthLogin {
		t.Fatalf("operation mismatch: %s", cfg.Operation)
	}
	if cfg.ProjectID != "my-project" {
		t.Fatalf("project-id mismatch: %s", cfg.ProjectID)
	}
	if cfg.AdditionalArgs != "--quiet" {
		t.Fatalf("additional args mismatch: %s", cfg.AdditionalArgs)
	}
}

func TestParseFlagsWithParser_SetProject(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationSetProjectConfig)
	parser.setString("project-id", "sample")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationSetProjectConfig || cfg.ProjectID != "sample" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("unknown operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unknown operation")
		}
	})

	t.Run("missing project", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAuthLogin)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when project-id is missing")
		}
	})

	t.Run("help bypass", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setBool("help", true)
		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !cfg.Help {
			t.Fatal("expected help flag to be true")
		}
	})
}

func TestParseFlags_StandardParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"cmd", "-operation=auth-login", "-project-id=my-project", "-additional-args=--quiet"}
	defer func() { os.Args = originalArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ProjectID != "my-project" || cfg.AdditionalArgs != "--quiet" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestStandardFlagParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"cmd", "-string=value", "-bool", "extra"}
	defer func() { os.Args = originalArgs }()

	parser := NewStandardFlagParser()
	var str string
	var flag bool

	parser.StringVar(&str, "string", "default", "")
	parser.BoolVar(&flag, "bool", false, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if str != "value" || !flag {
		t.Fatalf("unexpected parsed values: str=%s flag=%v", str, flag)
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Fatalf("unexpected remaining args: %v", args)
	}
}

func TestPrintUsage(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"gcloud-genset-init"}
	defer func() { os.Args = originalArgs }()

	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w

	PrintUsage()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read usage output: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "gcloud 初期設定向けのコマンド生成ツール") {
		t.Fatalf("usage output missing description: %s", output)
	}
	if !strings.Contains(output, "auth-login") || !strings.Contains(output, "set-project-config") {
		t.Fatalf("usage output missing operations: %s", output)
	}
}
