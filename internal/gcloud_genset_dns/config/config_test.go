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
	args         []string
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
	return m.args
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

func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

func TestParseFlagsWithParser_ManagedZonesList(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetStringFlag("operation", OperationManagedZonesList)
	parser.SetStringFlag("project", "my-project")
	parser.SetStringFlag("format", "json")
	parser.SetStringFlag("filter", "name=example")
	parser.SetIntFlag("limit", 20)
	parser.SetIntFlag("page-size", 10)
	parser.SetStringFlag("sort-by", "name")
	parser.SetStringFlag("verbosity", "debug")
	parser.SetBoolFlag("uri", true)
	parser.SetStringFlag("additional-args", "--account=my-account")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationManagedZonesList {
		t.Errorf("operation mismatch: got %s", cfg.Operation)
	}
	if cfg.Project != "my-project" {
		t.Errorf("project mismatch: got %s", cfg.Project)
	}
	if cfg.Format != "json" {
		t.Errorf("format mismatch: got %s", cfg.Format)
	}
	if cfg.Filter != "name=example" {
		t.Errorf("filter mismatch: got %s", cfg.Filter)
	}
	if cfg.Limit != 20 {
		t.Errorf("limit mismatch: got %d", cfg.Limit)
	}
	if cfg.PageSize != 10 {
		t.Errorf("page size mismatch: got %d", cfg.PageSize)
	}
	if cfg.SortBy != "name" {
		t.Errorf("sort-by mismatch: got %s", cfg.SortBy)
	}
	if cfg.Verbosity != "debug" {
		t.Errorf("verbosity mismatch: got %s", cfg.Verbosity)
	}
	if !cfg.URI {
		t.Errorf("expected uri flag to be true")
	}
	if cfg.AdditionalArgs != "--account=my-account" {
		t.Errorf("additional args mismatch: got %s", cfg.AdditionalArgs)
	}
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := NewMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("negative limit", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationManagedZonesList)
		parser.SetIntFlag("limit", -1)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for negative limit")
		}
	})

	t.Run("negative page size", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationManagedZonesList)
		parser.SetIntFlag("page-size", -5)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for negative page size")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unsupported operation")
		}
	})

	t.Run("unexpected args", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationManagedZonesList)
		parser.SetArgs([]string{"extra"})
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unexpected positional arguments")
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
