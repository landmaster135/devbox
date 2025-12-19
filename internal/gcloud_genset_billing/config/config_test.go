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
	parseError   error
	args         []string
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
	if len(m.args) == 0 {
		return []string{}
	}
	return append([]string(nil), m.args...)
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

func (m *mockFlagParser) setArgs(args []string) {
	m.args = append([]string(nil), args...)
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

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := newMockFlagParser()
	parser.setBool("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatalf("expected help flag to be true")
	}
}

func TestParseFlags_StandardParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		"cmd",
		"-operation=list-projects",
		"-limit=7",
		"-filter=project_id:test",
		"-billing-account=0000-1111-2222",
	}
	defer func() { os.Args = originalArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationListProjects {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
	if cfg.Limit != 7 {
		t.Fatalf("unexpected limit: %d", cfg.Limit)
	}
	if cfg.Filter != "project_id:test" {
		t.Fatalf("unexpected filter: %s", cfg.Filter)
	}
	if cfg.BillingAccount != "0000-1111-2222" {
		t.Fatalf("unexpected billing account: %s", cfg.BillingAccount)
	}
}

func TestStandardFlagParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"cmd", "-string=value", "-int=9", "-bool", "extra"}
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

	if str != "value" {
		t.Fatalf("unexpected string value: %s", str)
	}
	if num != 9 {
		t.Fatalf("unexpected int value: %d", num)
	}
	if !flag {
		t.Fatalf("expected bool flag to be true")
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Fatalf("unexpected remaining args: %v", args)
	}
}

func TestPrintUsage(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"gcloud-genset-billing"}
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
	if !strings.Contains(output, "Google Cloud Billing") {
		t.Fatalf("usage output missing expected description: %s", output)
	}
	if !strings.Contains(output, "list-budgets / list-projects 操作用") {
		t.Fatalf("usage output missing expected section: %s", output)
	}
}
