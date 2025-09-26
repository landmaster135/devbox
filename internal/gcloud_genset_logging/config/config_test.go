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
	intVars      map[string]*int
	boolVars     map[string]*bool
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	parseErr     error
}

func newMockFlagParser() *mockFlagParser {
	return &mockFlagParser{
		stringVars:   make(map[string]*string),
		intVars:      make(map[string]*int),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		intValues:    make(map[string]int),
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

func (m *mockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if preset, ok := m.intValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.intVars[name] = p
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
func (m *mockFlagParser) Args() []string { return []string{} }

func (m *mockFlagParser) setString(name, value string) {
	m.stringValues[name] = value
	if ptr, ok := m.stringVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setInt(name string, value int) {
	m.intValues[name] = value
	if ptr, ok := m.intVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) setBool(name string, value bool) {
	m.boolValues[name] = value
	if ptr, ok := m.boolVars[name]; ok {
		*ptr = value
	}
}

func TestParseFlagsWithParser_LoggingRead(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationLoggingRead)
	parser.setString("severity", "ERROR")
	parser.setString("resource-type", "gce_instance")
	parser.setString("query", "textPayload:\"Database\"")
	parser.setInt("limit", 25)
	parser.setString("additional-args", "--format=json")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationLoggingRead {
		t.Errorf("operation mismatch: got %s", cfg.Operation)
	}
	if cfg.Severity != "ERROR" {
		t.Errorf("severity mismatch: got %s", cfg.Severity)
	}
	if cfg.ResourceType != "gce_instance" {
		t.Errorf("resource type mismatch: got %s", cfg.ResourceType)
	}
	if cfg.Query != "textPayload:\"Database\"" {
		t.Errorf("query mismatch: got %s", cfg.Query)
	}
	if cfg.Limit != 25 {
		t.Errorf("limit mismatch: got %d", cfg.Limit)
	}
	if cfg.AdditionalArgs != "--format=json" {
		t.Errorf("additional args mismatch: got %s", cfg.AdditionalArgs)
	}
}

func TestParseFlagsWithParser_CreateSink(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationCreateSink)
	parser.setString("sink-name", "my-sink")
	parser.setString("destination", "storage.googleapis.com/my-bucket")
	parser.setString("log-filter", "resource.type=gce_instance")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationCreateSink {
		t.Errorf("operation mismatch: got %s", cfg.Operation)
	}
	if cfg.SinkName != "my-sink" {
		t.Errorf("sink name mismatch: got %s", cfg.SinkName)
	}
	if cfg.Destination != "storage.googleapis.com/my-bucket" {
		t.Errorf("destination mismatch: got %s", cfg.Destination)
	}
	if cfg.LogFilter != "resource.type=gce_instance" {
		t.Errorf("log filter mismatch: got %s", cfg.LogFilter)
	}
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("logging-read without filter", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationLoggingRead)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when filter criteria are missing")
		}
	})

	t.Run("logging-read invalid limit", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationLoggingRead)
		parser.setString("filter", "severity>=ERROR")
		parser.setInt("limit", 0)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for limit <= 0")
		}
	})

	t.Run("create-sink missing sink name", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateSink)
		parser.setString("destination", "storage.googleapis.com/my-bucket")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when sink-name is missing")
		}
	})

	t.Run("create-sink missing destination", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationCreateSink)
		parser.setString("sink-name", "my-sink")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when destination is missing")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unsupported operation")
		}
	})
}

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := newMockFlagParser()
	parser.setBool("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatalf("expected help to be true")
	}
}

func TestParseFlags_StandardParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		"cmd",
		"-operation=logging-read",
		"-filter=severity>=ERROR",
		"-limit=10",
		"-additional-args=--format=json",
	}
	defer func() { os.Args = originalArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationLoggingRead || cfg.Filter != "severity>=ERROR" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.Limit != 10 || cfg.AdditionalArgs != "--format=json" {
		t.Fatalf("unexpected numeric/additional args: %+v", cfg)
	}
}

func TestStandardFlagParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"cmd", "-string=value", "-int=7", "-bool", "extra"}
	defer func() { os.Args = originalArgs }()

	parser := NewStandardFlagParser()
	var str string
	var num int
	var flag bool

	parser.StringVar(&str, "string", "default", "")
	parser.IntVar(&num, "int", 0, "")
	parser.BoolVar(&flag, "bool", false, "")

	if err := parser.Parse(); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if str != "value" || num != 7 || !flag {
		t.Fatalf("unexpected parsed values str=%s num=%d flag=%v", str, num, flag)
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Fatalf("unexpected remaining args: %v", args)
	}
}

func TestPrintUsage(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"gcloud-genset-logging"}
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
	if !strings.Contains(output, "Google Cloud Logging") {
		t.Fatalf("usage output missing description: %s", output)
	}
	if !strings.Contains(output, "logging-read") || !strings.Contains(output, "create-sink") {
		t.Fatalf("usage output missing operations: %s", output)
	}
}
