package config

import "testing"

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

func (m *mockFlagParser) Parse() error {
	return m.parseErr
}

func (m *mockFlagParser) Args() []string {
	return nil
}

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

func TestParseFlagsWithParser(t *testing.T) {
	t.Run("undeploy processor version", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUndeployProcessorVersion)
		parser.setString("region", " us-central1 ")
		parser.setString("project-number", " 1234567890 ")
		parser.setString("processor-id", " proc-abc123 ")
		parser.setString("version-id", " version-001 ")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationUndeployProcessorVersion {
			t.Fatalf("operation mismatch: got %s", cfg.Operation)
		}
		if cfg.Region != "us-central1" {
			t.Fatalf("region mismatch: got %s", cfg.Region)
		}
		if cfg.ProjectNumber != "1234567890" {
			t.Fatalf("project number mismatch: got %s", cfg.ProjectNumber)
		}
		if cfg.ProcessorID != "proc-abc123" {
			t.Fatalf("processor id mismatch: got %s", cfg.ProcessorID)
		}
		if cfg.VersionID != "version-001" {
			t.Fatalf("version id mismatch: got %s", cfg.VersionID)
		}
	})
}

func TestParseFlagsWithParserErrors(t *testing.T) {
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
		parser.setString("region", "us-central1")
		parser.setString("processor-id", "proc")
		parser.setString("version-id", "ver")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when project-number is missing")
		}
	})

	t.Run("missing processor id", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUndeployProcessorVersion)
		parser.setString("region", "us-central1")
		parser.setString("project-number", "123")
		parser.setString("version-id", "ver")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when processor-id is missing")
		}
	})

	t.Run("missing version id", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationUndeployProcessorVersion)
		parser.setString("region", "us-central1")
		parser.setString("project-number", "123")
		parser.setString("processor-id", "proc")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when version-id is missing")
		}
	})

	t.Run("help flag skips validation", func(t *testing.T) {
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
