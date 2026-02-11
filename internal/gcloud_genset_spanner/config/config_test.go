package config

import "testing"

type mockFlagParser struct {
	stringVars map[string]*string
	intVars    map[string]*int
	boolVars   map[string]*bool

	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool

	parseError error
	args       []string
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
	if v, ok := m.stringValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *mockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if v, ok := m.intValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
	m.intVars[name] = p
}

func (m *mockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if v, ok := m.boolValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

func (m *mockFlagParser) Parse() error {
	return m.parseError
}

func (m *mockFlagParser) Args() []string {
	return m.args
}

func (m *mockFlagParser) SetString(name, value string) {
	m.stringValues[name] = value
	if ptr, ok := m.stringVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) SetInt(name string, value int) {
	m.intValues[name] = value
	if ptr, ok := m.intVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) SetBool(name string, value bool) {
	m.boolValues[name] = value
	if ptr, ok := m.boolVars[name]; ok {
		*ptr = value
	}
}

func (m *mockFlagParser) SetArgs(args []string) {
	m.args = args
}

func (m *mockFlagParser) SetParseError(err error) {
	m.parseError = err
}

func TestParseFlagsWithParser_InstanceCreate(t *testing.T) {
	parser := newMockFlagParser()
	parser.SetString("operation", OperationInstanceCreate)
	parser.SetString("instance-id", "prod-instance")
	parser.SetString("config", "regional-asia-northeast1")
	parser.SetString("description", "Payments")
	parser.SetInt("nodes", 2)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Operation != OperationInstanceCreate {
		t.Fatalf("operation mismatch: %s", cfg.Operation)
	}
	if cfg.InstanceID != "prod-instance" {
		t.Fatalf("instance-id mismatch: %s", cfg.InstanceID)
	}
	if cfg.InstanceConfig != "regional-asia-northeast1" {
		t.Fatalf("config mismatch: %s", cfg.InstanceConfig)
	}
	if cfg.Description != "Payments" {
		t.Fatalf("description mismatch: %s", cfg.Description)
	}
	if cfg.Nodes != 2 {
		t.Fatalf("nodes mismatch: %d", cfg.Nodes)
	}
}

func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := newMockFlagParser()
	parser.SetBool("help", true)

	cfg, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Help {
		t.Fatalf("help flag should be true")
	}
}

func TestParseFlagsWithParser_ValidationErrors(t *testing.T) {
	cases := []struct {
		name  string
		setup func(*mockFlagParser)
	}{
		{
			name:  "missing operation",
			setup: func(m *mockFlagParser) {},
		},
		{
			name: "unsupported operation",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", "unknown")
			},
		},
		{
			name: "instance-create missing instance-id",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", OperationInstanceCreate)
				m.SetString("config", "regional-europe-west1")
				m.SetString("description", "desc")
				m.SetInt("nodes", 1)
			},
		},
		{
			name: "instance-create missing config",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", OperationInstanceCreate)
				m.SetString("instance-id", "prod")
				m.SetString("description", "desc")
				m.SetInt("nodes", 1)
			},
		},
		{
			name: "instance-create missing description",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", OperationInstanceCreate)
				m.SetString("instance-id", "prod")
				m.SetString("config", "regional")
				m.SetInt("nodes", 1)
			},
		},
		{
			name: "instance-create invalid nodes",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", OperationInstanceCreate)
				m.SetString("instance-id", "prod")
				m.SetString("config", "regional")
				m.SetString("description", "desc")
				m.SetInt("nodes", 0)
			},
		},
		{
			name: "db-create missing ddl",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", OperationDatabaseCreate)
				m.SetString("instance-id", "prod")
				m.SetString("db-id", "orders")
			},
		},
		{
			name: "db-list missing instance",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", OperationDatabaseList)
			},
		},
		{
			name: "db-describe missing db-id",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", OperationDatabaseDescribe)
				m.SetString("instance-id", "prod")
			},
		},
		{
			name: "unexpected args",
			setup: func(m *mockFlagParser) {
				m.SetString("operation", OperationInstanceList)
				m.SetArgs([]string{"extra"})
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parser := newMockFlagParser()
			tc.setup(parser)
			if _, err := ParseFlagsWithParser(parser); err == nil {
				t.Fatalf("expected error but got nil")
			}
		})
	}
}
