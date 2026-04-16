package flag_parser

// MockFlagParser is a lightweight test double for FlagParser.
// It is shared by tests that verify ParseFlagsWithParser behavior.
type MockFlagParser struct {
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	parseError   error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: map[string]string{},
		intValues:    map[string]int{},
		boolValues:   map[string]bool{},
	}
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if v, ok := m.stringValues[name]; ok {
		*p = v
		return
	}
	if *p != "" {
		return
	}
	*p = value
}

func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if v, ok := m.intValues[name]; ok {
		*p = v
		return
	}
	if *p != 0 {
		return
	}
	*p = value
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if v, ok := m.boolValues[name]; ok {
		*p = v
		return
	}
	if *p {
		return
	}
	*p = value
}

func (m *MockFlagParser) Parse() error {
	return m.parseError
}

func (m *MockFlagParser) SetString(name, value string) {
	m.stringValues[name] = value
}

func (m *MockFlagParser) SetInt(name string, value int) {
	m.intValues[name] = value
}

func (m *MockFlagParser) SetBool(name string, value bool) {
	m.boolValues[name] = value
}

func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}
