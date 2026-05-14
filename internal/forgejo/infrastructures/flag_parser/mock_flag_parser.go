package flag_parser

// MockFlagParser はテスト用のFlagParserモックです。
type MockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string
	boolValues   map[string]bool
	args         []string
	parseError   error
}

// NewMockFlagParser は新しいMockFlagParserを作成します。
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		args:         []string{},
	}
}

// StringVar は文字列フラグを登録します。
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	_ = usage
	if v, ok := m.stringValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
	m.stringVars[name] = p
}

// BoolVar はブールフラグを登録します。
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	_ = usage
	if v, ok := m.boolValues[name]; ok {
		*p = v
	} else {
		*p = value
	}
	m.boolVars[name] = p
}

// Parse は解析を模擬します。
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は登録済みの位置引数を返します。
func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringFlag はテスト用に文字列フラグ値を事前設定します。
func (m *MockFlagParser) SetStringFlag(name, value string) {
	m.stringValues[name] = value
	if p, ok := m.stringVars[name]; ok {
		*p = value
	}
}

// SetBoolFlag はテスト用にブールフラグ値を事前設定します。
func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	m.boolValues[name] = value
	if p, ok := m.boolVars[name]; ok {
		*p = value
	}
}

// SetArgs はテスト用に位置引数を事前設定します。
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// SetParseError はParse時のエラーを設定します。
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}
