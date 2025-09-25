package config

import "testing"

type MockFlagParser struct {
	stringVars   map[string]*string
	intVars      map[string]*int
	boolVars     map[string]*bool
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	parseError   error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		intVars:      make(map[string]*int),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		intValues:    make(map[string]int),
		boolValues:   make(map[string]bool),
	}
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if preset, ok := m.stringValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if preset, ok := m.intValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.intVars[name] = p
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if preset, ok := m.boolValues[name]; ok {
		*p = preset
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

func (m *MockFlagParser) Parse() error {
	return m.parseError
}

func (m *MockFlagParser) Args() []string {
	return []string{}
}

func (m *MockFlagParser) SetStringFlag(name, value string) {
	m.stringValues[name] = value
	if ptr, ok := m.stringVars[name]; ok {
		*ptr = value
	}
}

func (m *MockFlagParser) SetIntFlag(name string, value int) {
	m.intValues[name] = value
	if ptr, ok := m.intVars[name]; ok {
		*ptr = value
	}
}

func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	m.boolValues[name] = value
	if ptr, ok := m.boolVars[name]; ok {
		*ptr = value
	}
}

func TestParseFlagsWithParser_LoggingRead(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetStringFlag("operation", OperationLoggingRead)
	parser.SetStringFlag("severity", "ERROR")
	parser.SetStringFlag("resource-type", "gce_instance")
	parser.SetStringFlag("query", "textPayload:\"Database\"")
	parser.SetIntFlag("limit", 25)
	parser.SetStringFlag("additional-args", "--format=json")

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
	parser := NewMockFlagParser()
	parser.SetStringFlag("operation", OperationCreateSink)
	parser.SetStringFlag("sink-name", "my-sink")
	parser.SetStringFlag("destination", "storage.googleapis.com/my-bucket")
	parser.SetStringFlag("log-filter", "resource.type=gce_instance")

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
		parser := NewMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("logging-read without filter", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationLoggingRead)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when filter criteria are missing")
		}
	})

	t.Run("logging-read invalid limit", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationLoggingRead)
		parser.SetStringFlag("filter", "severity>=ERROR")
		parser.SetIntFlag("limit", 0)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for limit <= 0")
		}
	})

	t.Run("create-sink missing sink name", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationCreateSink)
		parser.SetStringFlag("destination", "storage.googleapis.com/my-bucket")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when sink-name is missing")
		}
	})

	t.Run("create-sink missing destination", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationCreateSink)
		parser.SetStringFlag("sink-name", "my-sink")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when destination is missing")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unsupported operation")
		}
	})
}

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetBoolFlag("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatalf("expected help to be true")
	}
}
