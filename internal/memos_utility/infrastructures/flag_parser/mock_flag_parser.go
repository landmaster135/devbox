package flag_parser

// MockFlagParser はテスト用の FlagParser 実装。
type MockFlagParser struct {
	stringValues map[string]string
	intValues    map[string]int
	boolValues   map[string]bool
	args         []string
	parseError   error
}

// NewMockFlagParser は新しいモックを作成する。
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringValues: map[string]string{},
		intValues:    map[string]int{},
		boolValues:   map[string]bool{},
		args:         []string{},
	}
}

// StringVar は文字列フラグ値を適用する。
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

// IntVar は整数フラグ値を適用する。
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

// BoolVar は真偽値フラグ値を適用する。
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

// Parse は事前設定した parse エラーを返す。
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は事前設定した位置引数を返す。
func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetString は文字列フラグ値を設定する。
func (m *MockFlagParser) SetString(name, value string) {
	m.stringValues[name] = value
}

// SetInt は整数フラグ値を設定する。
func (m *MockFlagParser) SetInt(name string, value int) {
	m.intValues[name] = value
}

// SetBool は真偽値フラグ値を設定する。
func (m *MockFlagParser) SetBool(name string, value bool) {
	m.boolValues[name] = value
}

// SetArgs は位置引数を設定する。
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// SetParseError は Parse が返すエラーを設定する。
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}
