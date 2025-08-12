package config

import (
	"os"
	"testing"
)

// TestStandardFlagParser はStandardFlagParserのテストクラス
type TestStandardFlagParser struct{}

// TestNewStandardFlagParser_Normal はNewStandardFlagParserの正常系テスト
func TestNewStandardFlagParser_Normal(t *testing.T) {
	// Act
	parser := NewStandardFlagParser()

	// Assert
	if parser == nil {
		t.Fatal("parserがnilです")
	}
	if parser.flagSet == nil {
		t.Fatal("flagSetがnilです")
	}
	if parser.parsed != false {
		t.Errorf("parsedが期待値と異なります。期待値: false, 実際: %t", parser.parsed)
	}
}

// TestStandardFlagParser_StringVar_Normal はStringVarの正常系テスト
func TestStandardFlagParser_StringVar_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue string
	name := "test-string"
	defaultValue := "default"
	usage := "test usage"

	// Act
	parser.StringVar(&testValue, name, defaultValue, usage)

	// Assert
	if testValue != defaultValue {
		t.Errorf("testValueが期待値と異なります。期待値: %s, 実際: %s", defaultValue, testValue)
	}

	// flagSetに保存されているかを確認
	if storedPtr, exists := parser.flagSet[name]; !exists {
		t.Errorf("フラグ '%s' がflagSetに保存されていません", name)
	} else if storedPtr != &testValue {
		t.Errorf("フラグ '%s' のポインタが期待値と異なります", name)
	}
}

// TestStandardFlagParser_BoolVar_Normal はBoolVarの正常系テスト
func TestStandardFlagParser_BoolVar_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue bool
	name := "test-bool"
	defaultValue := true
	usage := "test bool usage"

	// Act
	parser.BoolVar(&testValue, name, defaultValue, usage)

	// Assert
	if testValue != defaultValue {
		t.Errorf("testValueが期待値と異なります。期待値: %t, 実際: %t", defaultValue, testValue)
	}

	// flagSetに保存されているかを確認
	if storedPtr, exists := parser.flagSet[name]; !exists {
		t.Errorf("フラグ '%s' がflagSetに保存されていません", name)
	} else if storedPtr != &testValue {
		t.Errorf("フラグ '%s' のポインタが期待値と異なります", name)
	}
}

// TestStandardFlagParser_Parse_Normal はParseの正常系テスト
func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()

	// Act
	err := parser.Parse()

	// Assert
	if err != nil {
		t.Fatalf("Parseでエラーが発生しました: %v", err)
	}
	if parser.parsed != true {
		t.Errorf("parsedが期待値と異なります。期待値: true, 実際: %t", parser.parsed)
	}
}

// TestStandardFlagParser_StringVar_MultipleFlags は複数のStringVarのテスト
func TestStandardFlagParser_StringVar_MultipleFlags(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var value1, value2 string

	// Act
	parser.StringVar(&value1, "flag1", "default1", "usage1")
	parser.StringVar(&value2, "flag2", "default2", "usage2")

	// Assert
	if value1 != "default1" {
		t.Errorf("value1が期待値と異なります。期待値: default1, 実際: %s", value1)
	}
	if value2 != "default2" {
		t.Errorf("value2が期待値と異なります。期待値: default2, 実際: %s", value2)
	}

	// 両方のフラグがflagSetに保存されているかを確認
	if _, exists := parser.flagSet["flag1"]; !exists {
		t.Error("flag1がflagSetに保存されていません")
	}
	if _, exists := parser.flagSet["flag2"]; !exists {
		t.Error("flag2がflagSetに保存されていません")
	}
}

// TestStandardFlagParser_BoolVar_MultipleFlags は複数のBoolVarのテスト
func TestStandardFlagParser_BoolVar_MultipleFlags(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var value1, value2 bool

	// Act
	parser.BoolVar(&value1, "bool1", true, "usage1")
	parser.BoolVar(&value2, "bool2", false, "usage2")

	// Assert
	if value1 != true {
		t.Errorf("value1が期待値と異なります。期待値: true, 実際: %t", value1)
	}
	if value2 != false {
		t.Errorf("value2が期待値と異なります。期待値: false, 実際: %t", value2)
	}

	// 両方のフラグがflagSetに保存されているかを確認
	if _, exists := parser.flagSet["bool1"]; !exists {
		t.Error("bool1がflagSetに保存されていません")
	}
	if _, exists := parser.flagSet["bool2"]; !exists {
		t.Error("bool2がflagSetに保存されていません")
	}
}

// TestStandardFlagParser_MixedFlags は文字列とブールフラグの混合テスト
func TestStandardFlagParser_MixedFlags(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var stringValue string
	var boolValue bool

	// Act
	parser.StringVar(&stringValue, "string-flag", "string-default", "string usage")
	parser.BoolVar(&boolValue, "bool-flag", true, "bool usage")

	// Assert
	if stringValue != "string-default" {
		t.Errorf("stringValueが期待値と異なります。期待値: string-default, 実際: %s", stringValue)
	}
	if boolValue != true {
		t.Errorf("boolValueが期待値と異なります。期待値: true, 実際: %t", boolValue)
	}

	// 両方のフラグがflagSetに保存されているかを確認
	if _, exists := parser.flagSet["string-flag"]; !exists {
		t.Error("string-flagがflagSetに保存されていません")
	}
	if _, exists := parser.flagSet["bool-flag"]; !exists {
		t.Error("bool-flagがflagSetに保存されていません")
	}
}

// TestStandardOSArgs はStandardOSArgsのテストクラス
type TestStandardOSArgs struct{}

// TestNewStandardOSArgs_Normal はNewStandardOSArgsの正常系テスト
func TestNewStandardOSArgs_Normal(t *testing.T) {
	// Act
	osArgs := NewStandardOSArgs()

	// Assert
	if osArgs == nil {
		t.Fatal("osArgsがnilです")
	}
}

// TestStandardOSArgs_Args_Normal はArgsの正常系テスト
func TestStandardOSArgs_Args_Normal(t *testing.T) {
	// Arrange
	osArgs := NewStandardOSArgs()

	// os.Argsを一時的に変更
	originalArgs := os.Args
	testArgs := []string{"test-program", "arg1", "arg2"}
	os.Args = testArgs
	defer func() { os.Args = originalArgs }()

	// Act
	result := osArgs.Args()

	// Assert
	if len(result) != len(testArgs) {
		t.Errorf("結果の長さが期待値と異なります。期待値: %d, 実際: %d", len(testArgs), len(result))
	}
	for i, expected := range testArgs {
		if i >= len(result) || result[i] != expected {
			t.Errorf("引数[%d]が期待値と異なります。期待値: %s, 実際: %s", i, expected, result[i])
		}
	}
}

// TestStandardOSArgs_Args_Empty は空の引数でのテスト
func TestStandardOSArgs_Args_Empty(t *testing.T) {
	// Arrange
	osArgs := NewStandardOSArgs()

	// os.Argsを一時的に変更
	originalArgs := os.Args
	os.Args = []string{}
	defer func() { os.Args = originalArgs }()

	// Act
	result := osArgs.Args()

	// Assert
	if len(result) != 0 {
		t.Errorf("結果の長さが期待値と異なります。期待値: 0, 実際: %d", len(result))
	}
}

// TestStandardOSArgs_Args_SingleArg は単一引数でのテスト
func TestStandardOSArgs_Args_SingleArg(t *testing.T) {
	// Arrange
	osArgs := NewStandardOSArgs()

	// os.Argsを一時的に変更
	originalArgs := os.Args
	testArgs := []string{"test-program"}
	os.Args = testArgs
	defer func() { os.Args = originalArgs }()

	// Act
	result := osArgs.Args()

	// Assert
	if len(result) != 1 {
		t.Errorf("結果の長さが期待値と異なります。期待値: 1, 実際: %d", len(result))
	}
	if result[0] != "test-program" {
		t.Errorf("引数[0]が期待値と異なります。期待値: test-program, 実際: %s", result[0])
	}
}
