package flag_parser

// MockFlagParser は config テスト用の FlagParser 実装です。
type MockFlagParser struct {
	stringValues map[string]string
	boolValues   map[string]bool
	intValues    map[string]int
	args         []string
	parseError   error
}

func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		intValues:    make(map[string]int),
		args:         []string{},
	}
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue
		return
	}
	*p = value
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue
		return
	}
	*p = value
}

func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if presetValue, exists := m.intValues[name]; exists {
		*p = presetValue
		return
	}
	*p = value
}

func (m *MockFlagParser) Parse() error {
	return m.parseError
}

func (m *MockFlagParser) Args() []string {
	return m.args
}

func (m *MockFlagParser) SetStringFlag(name string, value string) {
	m.stringValues[name] = value
}

func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	m.boolValues[name] = value
}

func (m *MockFlagParser) SetIntFlag(name string, value int) {
	m.intValues[name] = value
}

func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}
