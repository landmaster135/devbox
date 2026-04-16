package flag_parser

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string
	boolValues   map[string]bool
	args         []string
	parseError   error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		args:         []string{},
	}
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue
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
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	m.boolValues[name] = value
	if p, exists := m.boolVars[name]; exists {
		*p = value
	}
}

func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}
