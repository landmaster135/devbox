package config

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars map[string]*string
	intVars    map[string]*int
	boolVars   map[string]*bool
	args       []string
	parseError error
	parseHook  func()
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
	if m.parseHook != nil {
		m.parseHook()
	}
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

// SetBoolValue はテスト用にブール値を設定する
func (m *MockFlagParser) SetBoolValue(name string, value bool) {
	if p, exists := m.boolVars[name]; exists {
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

// SetParseHook はパース時に実行するフックを設定する
func (m *MockFlagParser) SetParseHook(hook func()) {
	m.parseHook = hook
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

// TestParseFlagsWithParser_HelpFlag はヘルプフラグ指定時の挙動をテストする
func TestParseFlagsWithParser_HelpFlag(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetParseHook(func() {
		parser.SetBoolValue("help", true)
	})

	config, err := ParseFlagsWithParser(parser)
	if err != nil {
		t.Fatalf("ParseFlagsWithParser() error = %v, want nil", err)
	}

	if !config.Help {
		t.Error("Help flag was not propagated to config")
	}
}

// TestParseFlagsWithParser_ParseError はパーサーでエラー発生時の挙動をテストする
func TestParseFlagsWithParser_ParseError(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetParseError(errors.New("parse failure"))

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatal("ParseFlagsWithParser() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "フラグの解析に失敗しました") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestParseFlagsWithParser_InvalidMaxDays は最大日数が無効な場合をテストする
func TestParseFlagsWithParser_InvalidMaxDays(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetParseHook(func() {
		parser.SetStringValue("max-days", "invalid")
	})

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatal("ParseFlagsWithParser() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "無効な最大日数です") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestParseFlagsWithParser_NewConfigError はConfig生成でエラーとなる場合をテストする
func TestParseFlagsWithParser_NewConfigError(t *testing.T) {
	parser := NewMockFlagParser()
	parser.SetParseHook(func() {
		parser.SetStringValue("max-days", "10")
		parser.SetStringValue("api-key", "test-api-key")
		parser.SetStringValue("city", "Tokyo,JP")
	})

	_, err := ParseFlagsWithParser(parser)
	if err == nil {
		t.Fatal("ParseFlagsWithParser() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "最大日数は1-5の範囲で指定してください") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestStandardFlagParser_ParseAndArgs はParseおよびArgsの挙動をテストする
func TestStandardFlagParser_ParseAndArgs(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{
		"cmd",
		"-test-string", "updated",
		"-test-int", "99",
		"-test-bool",
		"pos1", "pos2",
	}

	parser := NewStandardFlagParser()

	var testString string
	var testInt int
	var testBool bool

	parser.StringVar(&testString, "test-string", "default", "test string")
	parser.IntVar(&testInt, "test-int", 42, "test int")
	parser.BoolVar(&testBool, "test-bool", false, "test bool")

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}

	if testString != "updated" {
		t.Errorf("testString = %v, want updated", testString)
	}
	if testInt != 99 {
		t.Errorf("testInt = %v, want 99", testInt)
	}
	if testBool != true {
		t.Errorf("testBool = %v, want true", testBool)
	}

	args := parser.Args()
	if len(args) != 2 || args[0] != "pos1" || args[1] != "pos2" {
		t.Errorf("Args() = %v, want [pos1 pos2]", args)
	}
}
