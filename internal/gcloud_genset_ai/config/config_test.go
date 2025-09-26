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
	args         []string
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

func (m *mockFlagParser) StringVar(p *string, name string, value string, usage string) {
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
func (m *mockFlagParser) Args() []string { return m.args }

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

func (m *mockFlagParser) setArgs(args []string) { m.args = args }

func TestParseFlagsWithParser_Success(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationUndeployProcessorVersion)
	parser.setString("region", " us-central1 ")
	parser.setString("project-number", " 1234567890 ")
	parser.setString("processor-id", " processor-abc ")
	parser.setString("version-id", " version-001 ")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationUndeployProcessorVersion {
		t.Fatalf("operation mismatch: %s", cfg.Operation)
	}
	if cfg.Region != "us-central1" {
		t.Fatalf("region mismatch: %s", cfg.Region)
	}
	if cfg.ProjectNumber != "1234567890" {
		t.Fatalf("project number mismatch: %s", cfg.ProjectNumber)
	}
	if cfg.ProcessorID != "processor-abc" {
		t.Fatalf("processor id mismatch: %s", cfg.ProcessorID)
	}
	if cfg.VersionID != "version-001" {
		t.Fatalf("version id mismatch: %s", cfg.VersionID)
	}
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is unknown")
		}
	})

	t.Run("missing region", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUndeployProcessorVersion)
		parser.setString("project-number", "123")
		parser.setString("processor-id", "proc")
		parser.setString("version-id", "ver")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when region is missing")
		}
	})

	t.Run("missing project number", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUndeployProcessorVersion)
		parser.setString("region", "us")
		parser.setString("processor-id", "proc")
		parser.setString("version-id", "ver")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when project-number is missing")
		}
	})

	t.Run("missing processor id", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUndeployProcessorVersion)
		parser.setString("region", "us")
		parser.setString("project-number", "123")
		parser.setString("version-id", "ver")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when processor-id is missing")
		}
	})

	t.Run("missing version id", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUndeployProcessorVersion)
		parser.setString("region", "us")
		parser.setString("project-number", "123")
		parser.setString("processor-id", "proc")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when version-id is missing")
		}
	})

	t.Run("help flag bypasses validation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setBool("help", true)
		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !cfg.Help {
			t.Fatal("help flag should be true")
		}
	})

}

func TestParseFlags_StandardParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		"cmd",
		"-operation=undeploy-processor-version",
		"-region=us-central1",
		"-project-number=1234567890",
		"-processor-id=processor-abc",
		"-version-id=version-001",
	}
	defer func() { os.Args = originalArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationUndeployProcessorVersion || cfg.Region != "us-central1" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.ProjectNumber != "1234567890" || cfg.ProcessorID != "processor-abc" || cfg.VersionID != "version-001" {
		t.Fatalf("unexpected parsed values: %+v", cfg)
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
	os.Args = []string{"gcloud-genset-ai"}
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
	if !strings.Contains(output, "Document AI") {
		t.Fatalf("usage output missing description: %s", output)
	}
	if !strings.Contains(output, "undeploy-processor-version") {
		t.Fatalf("usage output missing operation description: %s", output)
	}
}
