package config

import "testing"

type MockFlagParser struct {
	stringVars   map[string]*string
	intVars      map[string]*int
	boolVars     map[string]*bool
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	args         []string
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

func TestParseFlagsWithParser_ListDashboards(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetStringFlag("operation", OperationListDashboards)
	parser.SetStringFlag("project", "my-project")
	parser.SetStringFlag("filter", "displayName:test")
	parser.SetStringFlag("format", "json")
	parser.SetIntFlag("page-size", 50)
	parser.SetStringFlag("sort-by", "name")
	parser.SetIntFlag("limit", 100)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationListDashboards {
		t.Errorf("operation mismatch: %s", cfg.Operation)
	}
	if cfg.Project != "my-project" {
		t.Errorf("project mismatch: %s", cfg.Project)
	}
	if cfg.Filter != "displayName:test" {
		t.Errorf("filter mismatch: %s", cfg.Filter)
	}
	if cfg.Format != "json" {
		t.Errorf("format mismatch: %s", cfg.Format)
	}
	if cfg.PageSize != 50 {
		t.Errorf("page-size mismatch: %d", cfg.PageSize)
	}
	if cfg.SortBy != "name" {
		t.Errorf("sort-by mismatch: %s", cfg.SortBy)
	}
	if cfg.Limit != 100 {
		t.Errorf("limit mismatch: %d", cfg.Limit)
	}
}

func TestParseFlagsWithParser_DescribeDashboard(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetStringFlag("operation", OperationDescribeDashboard)
	parser.SetStringFlag("project", "my-project")
	parser.SetStringFlag("dashboard-id", "dash123")
	parser.SetStringFlag("format", "yaml")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DashboardID != "dash123" {
		t.Errorf("dashboard-id mismatch: %s", cfg.DashboardID)
	}
	if cfg.Format != "yaml" {
		t.Errorf("format mismatch: %s", cfg.Format)
	}
}

func TestParseFlagsWithParser_ListSnoozes(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetStringFlag("operation", OperationListSnoozes)
	parser.SetStringFlag("project", "my-project")
	parser.SetStringFlag("filter", "displayName:maintenance")
	parser.SetIntFlag("page-size", 10)
	parser.SetIntFlag("limit", 20)
	parser.SetBoolFlag("uri", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.URI {
		t.Fatal("expected uri flag to be true")
	}
	if cfg.PageSize != 10 {
		t.Errorf("page-size mismatch: %d", cfg.PageSize)
	}
	if cfg.Limit != 20 {
		t.Errorf("limit mismatch: %d", cfg.Limit)
	}
}

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetBoolFlag("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatal("expected help to be true")
	}
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := NewMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unsupported operation")
		}
	})

	t.Run("describe missing dashboard id", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationDescribeDashboard)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when dashboard-id is missing")
		}
	})

	t.Run("list dashboards invalid limit", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationListDashboards)
		parser.SetIntFlag("limit", 0)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when limit is zero")
		}
	})

	t.Run("describe with filter", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationDescribeDashboard)
		parser.SetStringFlag("dashboard-id", "dash123")
		parser.SetStringFlag("filter", "name=test")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when filter is provided for describe")
		}
	})

	t.Run("unexpected args", func(t *testing.T) {
		parser := NewMockFlagParser()
		parser.SetStringFlag("operation", OperationListDashboards)
		parser.SetArgs([]string{"extra"})
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when extra args exist")
		}
	})
}
