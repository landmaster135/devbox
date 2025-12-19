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
	args         []string
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
func (m *mockFlagParser) Args() []string { return m.args }

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

func (m *mockFlagParser) setArgs(args []string) { m.args = args }

func TestParseFlagsWithParser_Success(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationListDashboards)
	parser.setString("project", "my-project")
	parser.setString("filter", "displayName:test")
	parser.setString("format", "json")
	parser.setInt("page-size", 50)
	parser.setString("sort-by", "name")
	parser.setInt("limit", 100)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationListDashboards || cfg.Project != "my-project" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.Filter != "displayName:test" || cfg.Format != "json" {
		t.Fatalf("unexpected filter/format: %+v", cfg)
	}
	if cfg.PageSize != 50 || cfg.SortBy != "name" || cfg.Limit != 100 {
		t.Fatalf("unexpected paging values: %+v", cfg)
	}
}

func TestParseFlagsWithParser_DescribeDashboard(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationDescribeDashboard)
	parser.setString("project", "my-project")
	parser.setString("dashboard-id", "dash123")
	parser.setString("format", "yaml")

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DashboardID != "dash123" || cfg.Format != "yaml" {
		t.Fatalf("unexpected describe values: %+v", cfg)
	}
}

func TestParseFlagsWithParser_ListSnoozes(t *testing.T) {
	parser := newMockFlagParser()
	parser.setString("operation", OperationListSnoozes)
	parser.setString("project", "my-project")
	parser.setInt("page-size", 10)
	parser.setInt("limit", 20)
	parser.setBool("uri", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.URI || cfg.PageSize != 10 || cfg.Limit != 20 {
		t.Fatalf("unexpected snoozes config: %+v", cfg)
	}
}

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := newMockFlagParser()
	parser.setBool("help", true)

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
		parser := newMockFlagParser()
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when operation is missing")
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", "unknown")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error for unsupported operation")
		}
	})

	t.Run("describe missing dashboard id", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDescribeDashboard)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when dashboard-id is missing")
		}
	})

	t.Run("list dashboards invalid limit", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListDashboards)
		parser.setInt("limit", 0)
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when limit is zero")
		}
	})

	t.Run("describe with filter", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationDescribeDashboard)
		parser.setString("dashboard-id", "dash123")
		parser.setString("filter", "name=test")
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when filter is provided for describe")
		}
	})

	t.Run("unexpected args", func(t *testing.T) {
		parser := newMockFlagParser()
		parser.setString("operation", OperationListDashboards)
		parser.setArgs([]string{"extra"})
		if _, err := ParseFlagsWithParser(parser); err == nil {
			t.Fatal("expected error when extra args exist")
		}
	})
}

func TestParseFlags_StandardParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{
		"cmd",
		"-operation=list-snoozes",
		"-project=my-project",
		"-filter=displayName:test",
		"-page-size=5",
		"-limit=10",
		"-format=json",
		"-uri",
	}
	defer func() { os.Args = originalArgs }()

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Operation != OperationListSnoozes || cfg.Project != "my-project" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
	if cfg.Format != "json" || cfg.Filter != "displayName:test" || cfg.PageSize != 5 || cfg.Limit != 10 || !cfg.URI {
		t.Fatalf("unexpected flag values: %+v", cfg)
	}
}

func TestStandardFlagParser(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"cmd", "-string=value", "-int=3", "-bool", "extra"}
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
	if str != "value" || num != 3 || !flag {
		t.Fatalf("unexpected parsed values: str=%s num=%d flag=%v", str, num, flag)
	}

	args := parser.Args()
	if len(args) != 1 || args[0] != "extra" {
		t.Fatalf("unexpected remaining args: %v", args)
	}
}

func TestPrintUsage(t *testing.T) {
	originalArgs := os.Args
	os.Args = []string{"gcloud-genset-monitoring"}
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
	if !strings.Contains(output, "Google Cloud Monitoring") {
		t.Fatalf("usage output missing description: %s", output)
	}
	if !strings.Contains(output, "list-dashboards") || !strings.Contains(output, "list-snoozes") {
		t.Fatalf("usage output missing operations: %s", output)
	}
}
