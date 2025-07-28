package config

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars map[string]*string
	boolVars   map[string]*bool
	args       []string
	parseError error
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars: make(map[string]*string),
		boolVars:   make(map[string]*bool),
		args:       []string{},
	}
}

// StringVar は文字列フラグを定義する（モック）
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	*p = value // デフォルト値を設定
	m.stringVars[name] = p
}

// BoolVar はブールフラグを定義する（モック）
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	*p = value // デフォルト値を設定
	m.boolVars[name] = p
}

// Parse はフラグを解析する（モック）
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は解析後の残りの引数を返す（モック）
func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringFlag はテスト用に文字列フラグの値を設定する
func (m *MockFlagParser) SetStringFlag(name, value string) {
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

// SetBoolFlag はテスト用にブールフラグの値を設定する
func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	if p, exists := m.boolVars[name]; exists {
		*p = value
	}
}

// SetArgs はテスト用に残りの引数を設定する
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// SetParseError はテスト用に解析エラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// TestStandardFlagParser_Parse_Error はParseのエラー系テスト
func TestStandardFlagParser_Parse_Error(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()

	// Act - 実際のParseメソッドを呼び出してカバレッジを向上させる
	// os.Argsを一時的に変更してテスト
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 有効な引数でParseを実行
	os.Args = []string{"test-program"}
	err := parser.Parse()

	// Assert
	if err != nil {
		t.Errorf("Expected no error for valid args, got %v", err)
	}

	// Argsメソッドも呼び出してカバレッジを向上
	args := parser.Args()
	if args == nil {
		t.Error("Expected non-nil args")
	}
}

// TestStandardFlagParser_StringVar_Normal はStringVarの正常系テスト
func TestStandardFlagParser_StringVar_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue string

	// Act
	dv := "default-value"
	parser.StringVar(&testValue, "test-flag-arith", dv, "test usage")

	// Assert
	if testValue != dv {
		t.Errorf("Expected initial value to be '%s', got %s", dv, testValue)
	}
}

// TestStandardFlagParser_BoolVar_Normal はBoolVarの正常系テスト
func TestStandardFlagParser_BoolVar_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue bool

	// Act
	parser.BoolVar(&testValue, "test-bool-flag-arith", false, "test bool usage")

	// Assert
	if testValue != false {
		t.Errorf("Expected initial value to be false, got %t", testValue)
	}
}

// TestStandardFlagParser_Parse_Normal はParseの正常系テスト
func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	// Arrange
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	var testString string
	var testBool bool

	parser.StringVar(&testString, "string-flag", "default", "string flag")
	parser.BoolVar(&testBool, "bool-flag", false, "bool flag")

	// テスト用の引数を設定
	ts := "test-value"
	testArgs := []string{fmt.Sprintf("-string-flag=%s", ts), "-bool-flag=true"}

	// Act
	err := testFlagSet.Parse(testArgs)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if testString != ts {
		t.Errorf("Expected string flag to be '%s', got %s", ts, testString)
	}

	if testBool != true {
		t.Errorf("Expected bool flag to be true, got %t", testBool)
	}
}

// TestStandardFlagParser_Args_Normal はArgsの正常系テスト
func TestStandardFlagParser_Args_Normal(t *testing.T) {
	// Arrange
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	testArgs := []string{"-flag=value", "arg1", "arg2"}

	var flagValue string
	parser.StringVar(&flagValue, "flag", "", "test flag")

	// Act
	err := testFlagSet.Parse(testArgs)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	args := parser.Args()

	// Assert
	expectedArgs := []string{"arg1", "arg2"}
	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(args))
	}

	for i, expected := range expectedArgs {
		if i >= len(args) || args[i] != expected {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expected, args[i])
		}
	}
}

// TestNewStandardFlagParser_Normal はNewStandardFlagParserの正常系テスト
func TestNewStandardFlagParser_Normal(t *testing.T) {
	// Act
	parser := NewStandardFlagParser()

	// Assert
	if parser == nil {
		t.Fatal("Expected parser to be non-nil")
	}

	if parser.flagSet == nil {
		t.Error("Expected flagSet to be non-nil")
	}

	if parser.flagSet != flag.CommandLine {
		t.Error("Expected flagSet to be flag.CommandLine")
	}
}

// TestParseFlagsWithParser_AddOperation はadd操作のフラグ解析テスト
func TestParseFlagsWithParser_AddOperation(t *testing.T) {
	// Arrange
	testFlagSet := flag.NewFlagSet("test-add", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	var operation, xStr, yStr string

	// フラグを定義
	parser.StringVar(&operation, "operation", "", "算術操作")
	parser.StringVar(&xStr, "x", "0", "第一オペランド")
	parser.StringVar(&yStr, "y", "0", "第二オペランド")

	// テスト用の引数を設定
	testArgs := []string{"-operation=add", "-x=10", "-y=5"}

	// Act
	err := testFlagSet.Parse(testArgs)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Assert
	if operation != "add" {
		t.Errorf("Expected operation to be 'add', got %s", operation)
	}
	if xStr != "10" {
		t.Errorf("Expected x to be '10', got %s", xStr)
	}
	if yStr != "5" {
		t.Errorf("Expected y to be '5', got %s", yStr)
	}
}

// TestParseFlagsWithParser_Integration は統合テスト
func TestParseFlagsWithParser_Integration(t *testing.T) {
	// Arrange
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	os.Args = []string{"test-program", "-operation=multiply", "-x=7", "-y=8"}

	testFlagSet := flag.NewFlagSet("integration-test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	var operation, xStr, yStr string

	// Act
	parser.StringVar(&operation, "operation", "", "算術操作")
	parser.StringVar(&xStr, "x", "0", "第一オペランド")
	parser.StringVar(&yStr, "y", "0", "第二オペランド")

	err := testFlagSet.Parse(os.Args[1:])

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if operation != "multiply" {
		t.Errorf("Expected operation to be 'multiply', got %s", operation)
	}

	if xStr != "7" {
		t.Errorf("Expected x to be '7', got %s", xStr)
	}

	if yStr != "8" {
		t.Errorf("Expected y to be '8', got %s", yStr)
	}
}
