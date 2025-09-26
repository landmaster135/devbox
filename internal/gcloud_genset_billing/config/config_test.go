package config

import "testing"

type mockFlagParser struct {
	stringVars   map[string]*string
	intVars      map[string]*int
	boolVars     map[string]*bool
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	parseError   error
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

func (m *mockFlagParser) Parse() error {
	return m.parseError
}

func (m *mockFlagParser) Args() []string {
	return []string{}
}

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

func TestParseFlagsWithParser_Success(t *testing.T) {
	t.Run("list budgets", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListBudgets)
		parser.setInt("limit", 25)
		parser.setString("billing-account", "0000-9999-8888")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Limit != 25 {
			t.Fatalf("unexpected limit: %d", cfg.Limit)
		}
		if cfg.BillingAccount != "0000-9999-8888" {
			t.Fatalf("unexpected billing account: %s", cfg.BillingAccount)
		}
	})

	t.Run("list projects", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListProjects)
		parser.setString("filter", "project_id:test")
		parser.setInt("limit", 5)

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.Filter != "project_id:test" {
			t.Fatalf("unexpected filter: %s", cfg.Filter)
		}
	})

	t.Run("describe project", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDescribeProject)
		parser.setString("project-id", "my-project")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.ProjectID != "my-project" {
			t.Fatalf("unexpected project id: %s", cfg.ProjectID)
		}
	})

	t.Run("describe budget", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDescribeBudget)
		parser.setString("budget-id", "00AA00-123456-AAAA")

		cfg, err := ParseFlagsWithParser(parser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if cfg.BudgetID != "00AA00-123456-AAAA" {
			t.Fatalf("unexpected budget id: %s", cfg.BudgetID)
		}
	})
}

func TestParseFlagsWithParser_Errors(t *testing.T) {
	t.Run("missing operation", func(t *testing.T) {
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("invalid operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is invalid")
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListBudgets)
		parser.setInt("limit", 0)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when limit <= 0")
		}
	})

	t.Run("missing budget id", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDescribeBudget)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when budget-id is missing")
		}
	})
}
