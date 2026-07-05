package flag_parser

// MockFlagParser は FlagParser のモック実装です。
type MockFlagParser struct {
	stringValues map[string]string
	intValues    map[string]int
	floatValues  map[string]float64
	boolValues   map[string]bool
	parseError   error
	parseCalled  bool
}

// NewMockFlagParser は新しいモックを生成します。
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: make(map[string]string),
		intValues:    make(map[string]int),
		floatValues:  make(map[string]float64),
		boolValues:   make(map[string]bool),
	}
}

func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if v, ok := m.stringValues[name]; ok {
		*p = v
		return
	}
	*p = value
}

func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	if v, ok := m.intValues[name]; ok {
		*p = v
		return
	}
	*p = value
}

func (m *MockFlagParser) Float64Var(p *float64, name string, value float64, usage string) {
	if v, ok := m.floatValues[name]; ok {
		*p = v
		return
	}
	*p = value
}

func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if v, ok := m.boolValues[name]; ok {
		*p = v
		return
	}
	*p = value
}

func (m *MockFlagParser) Parse() error {
	m.parseCalled = true
	return m.parseError
}

func (m *MockFlagParser) SetStringValue(name, value string) {
	m.stringValues[name] = value
}

func (m *MockFlagParser) SetIntValue(name string, value int) {
	m.intValues[name] = value
}

func (m *MockFlagParser) SetFloat64Value(name string, value float64) {
	m.floatValues[name] = value
}

func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	m.boolValues[name] = value
}

func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

func (m *MockFlagParser) ParseCalled() bool {
	return m.parseCalled
}
