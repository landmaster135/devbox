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
	t.Run("auth login", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAuthLogin)
		parser.setString("project-id", " my-project ")
		parser.setString("additional-args", " --quiet ")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationAuthLogin {
			t.Fatalf("operation mismatch: got %s", cfg.Operation)
		}
		if cfg.ProjectID != "my-project" {
			t.Fatalf("project id mismatch: got %s", cfg.ProjectID)
		}
		if cfg.AdditionalArgs != "--quiet" {
			t.Fatalf("additional args mismatch: got %s", cfg.AdditionalArgs)
		}
	})

	t.Run("set project config", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationSetProjectConfig)
		parser.setString("project-id", "sample-project")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Operation != OperationSetProjectConfig {
			t.Fatalf("operation mismatch: got %s", cfg.Operation)
		}
		if cfg.ProjectID != "sample-project" {
			t.Fatalf("project id mismatch: got %s", cfg.ProjectID)
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

	t.Run("missing project id", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationAuthLogin)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when project-id is missing")
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
