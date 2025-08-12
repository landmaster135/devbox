package config

import (
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars map[string]*string
	intVars    map[string]*int
	boolVars   map[string]*bool
	args       []string
	parseError error
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars: make(map[string]*string),
		intVars:    make(map[string]*int),
		boolVars:   make(map[string]*bool),
		args:       []string{},
	}
}

// StringVar は文字列フラグを定義する
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	*p = value
	m.stringVars[name] = p
}

// IntVar は整数フラグを定義する
func (m *MockFlagParser) IntVar(p *int, name string, value int, usage string) {
	*p = value
	m.intVars[name] = p
}

// BoolVar はブールフラグを定義する
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	*p = value
	m.boolVars[name] = p
}

// Parse はコマンドライン引数を解析する
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は解析後の残りの引数を返す
func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringValue はテスト用に文字列値を設定する
func (m *MockFlagParser) SetStringValue(name, value string) {
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

// SetArgs はテスト用に引数を設定する
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// SetParseError はテスト用にパースエラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// TestParseFlagsWithParser_PositionalArgs は位置引数のテストを行う
func TestParseFlagsWithParser_PositionalArgs(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetArgs([]string{"test-api-key", "Tokyo,JP", "3"})

	config, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Errorf("ParseFlagsWithParser() error = %v, want nil", err)
		return
	}

	if config.APIKey != "test-api-key" {
		t.Errorf("APIKey = %v, want test-api-key", config.APIKey)
	}
	if config.City != "Tokyo,JP" {
		t.Errorf("City = %v, want Tokyo,JP", config.City)
	}
	if config.MaxDays != 3 {
		t.Errorf("MaxDays = %v, want 3", config.MaxDays)
	}
}

// TestStandardFlagParser_Normal は標準フラグパーサーの基本機能をテストする
func TestStandardFlagParser_Normal(t *testing.T) {
	parser := NewStandardFlagParser()

	var testString string
	var testInt int
	var testBool bool

	parser.StringVar(&testString, "test-string", "default", "test string")
	parser.IntVar(&testInt, "test-int", 42, "test int")
	parser.BoolVar(&testBool, "test-bool", false, "test bool")

	// デフォルト値の確認
	if testString != "default" {
		t.Errorf("testString = %v, want default", testString)
	}
	if testInt != 42 {
		t.Errorf("testInt = %v, want 42", testInt)
	}
	if testBool != false {
		t.Errorf("testBool = %v, want false", testBool)
	}
}
