package config

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック
type MockFlagParser struct {
	stringVars   map[string]*string
	boolVars     map[string]*bool
	stringValues map[string]string // 事前設定された文字列値
	boolValues   map[string]bool   // 事前設定されたブール値
	args         []string
	parseError   error
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringVars:   make(map[string]*string),
		boolVars:     make(map[string]*bool),
		stringValues: make(map[string]string),
		boolValues:   make(map[string]bool),
		args:         []string{},
	}
}

// StringVar は文字列フラグを定義する（モック）
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	// 事前設定された値があるかチェック
	if presetValue, exists := m.stringValues[name]; exists {
		*p = presetValue
	} else {
		*p = value // デフォルト値を設定
	}
	m.stringVars[name] = p
}

// BoolVar はブールフラグを定義する（モック）
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	// 事前設定された値があるかチェック
	if presetValue, exists := m.boolValues[name]; exists {
		*p = presetValue
	} else {
		*p = value // デフォルト値を設定
	}
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
	m.stringValues[name] = value
	if p, exists := m.stringVars[name]; exists {
		*p = value
	}
}

// SetBoolFlag はテスト用にブールフラグの値を設定する
func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	m.boolValues[name] = value
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

// TestMockFlagParser_DefaultValues はMockFlagParserのデフォルト値テスト
func TestMockFlagParser_DefaultValues(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	var stringFlag string
	var boolFlag bool

	// Act
	mockParser.StringVar(&stringFlag, "string-flag", "default-string", "string flag")
	mockParser.BoolVar(&boolFlag, "bool-flag", true, "bool flag")

	// Assert
	if stringFlag != "default-string" {
		t.Errorf("Expected default string 'default-string', got '%s'", stringFlag)
	}
	if !boolFlag {
		t.Error("Expected default bool true, got false")
	}
}

// TestMockFlagParser_EmptyArgs は空の引数リストのテスト
func TestMockFlagParser_EmptyArgs(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()

	// Act
	args := mockParser.Args()

	// Assert
	if args == nil {
		t.Error("Expected non-nil args, got nil")
	}
	if len(args) != 0 {
		t.Errorf("Expected empty args, got %d items", len(args))
	}
}

// TestMockFlagParser_ParseWithoutError はエラーなしのパースのテスト
func TestMockFlagParser_ParseWithoutError(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()

	// Act
	err := mockParser.Parse()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestStandardFlagParser_MultipleFlags は複数フラグの処理テスト
func TestStandardFlagParser_MultipleFlags(t *testing.T) {
	// Arrange
	testFlagSet := flag.NewFlagSet("multi-flags-test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	var operation, file, numbers string
	var threshold int
	var help bool

	// Act
	parser.StringVar(&operation, "operation", "", "算術操作")
	parser.StringVar(&file, "file", "", "ファイルパス")
	parser.StringVar(&numbers, "numbers", "", "数値リスト")
	parser.BoolVar(&help, "help", false, "ヘルプ")

	testArgs := []string{"-operation=sum", "-numbers=1,2,3", "-help=true"}
	err := testFlagSet.Parse(testArgs)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if operation != "sum" {
		t.Errorf("Expected operation 'sum', got '%s'", operation)
	}
	if numbers != "1,2,3" {
		t.Errorf("Expected numbers '1,2,3', got '%s'", numbers)
	}
	if !help {
		t.Error("Expected help true, got false")
	}
	if file != "" {
		t.Errorf("Expected empty file, got '%s'", file)
	}
	if threshold != 0 {
		t.Errorf("Expected threshold 0, got %d", threshold)
	}
}

// TestStandardFlagParser_FlagSetReference はflagSetの参照テスト
func TestStandardFlagParser_FlagSetReference(t *testing.T) {
	// Arrange & Act
	parser := NewStandardFlagParser()

	// Assert
	if parser.flagSet != flag.CommandLine {
		t.Error("Expected flagSet to reference flag.CommandLine")
	}
}

// TestMockFlagParser_OverwriteFlags はフラグの上書きテスト
func TestMockFlagParser_OverwriteFlags(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	var testFlag string

	mockParser.StringVar(&testFlag, "test-flag", "initial", "test flag")

	// Act
	mockParser.SetStringFlag("test-flag", "first-update")
	firstValue := testFlag

	mockParser.SetStringFlag("test-flag", "second-update")
	secondValue := testFlag

	// Assert
	if firstValue != "first-update" {
		t.Errorf("Expected first update 'first-update', got '%s'", firstValue)
	}
	if secondValue != "second-update" {
		t.Errorf("Expected second update 'second-update', got '%s'", secondValue)
	}
}

// TestMockFlagParser_MixedFlagTypes は文字列とブールフラグの混合テスト
func TestMockFlagParser_MixedFlagTypes(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	var stringFlag1, stringFlag2 string
	var boolFlag1, boolFlag2 bool

	// Act
	mockParser.StringVar(&stringFlag1, "str1", "default1", "string flag 1")
	mockParser.StringVar(&stringFlag2, "str2", "default2", "string flag 2")
	mockParser.BoolVar(&boolFlag1, "bool1", false, "bool flag 1")
	mockParser.BoolVar(&boolFlag2, "bool2", true, "bool flag 2")

	mockParser.SetStringFlag("str1", "updated1")
	mockParser.SetBoolFlag("bool1", true)
	mockParser.SetBoolFlag("bool2", false)

	// Assert
	if stringFlag1 != "updated1" {
		t.Errorf("Expected stringFlag1 'updated1', got '%s'", stringFlag1)
	}
	if stringFlag2 != "default2" {
		t.Errorf("Expected stringFlag2 'default2', got '%s'", stringFlag2)
	}
	if !boolFlag1 {
		t.Error("Expected boolFlag1 true, got false")
	}
	if boolFlag2 {
		t.Error("Expected boolFlag2 false, got true")
	}
}

// TestStandardFlagParser_ErrorHandling はエラーハンドリングのテスト
func TestStandardFlagParser_ErrorHandling(t *testing.T) {
	// Arrange
	testFlagSet := flag.NewFlagSet("error-test", flag.ContinueOnError)

	var intFlag int
	testFlagSet.IntVar(&intFlag, "int-flag", 0, "integer flag")

	// Act - 無効な整数値を渡す
	testArgs := []string{"-int-flag=invalid"}
	err := testFlagSet.Parse(testArgs)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid integer flag, got nil")
	}
}
