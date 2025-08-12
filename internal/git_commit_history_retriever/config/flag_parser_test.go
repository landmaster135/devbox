package config

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
// TestStandardFlagParser_StringVar_Normal はStringVarの正常系テスト
func TestStandardFlagParser_StringVar_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue string

	// Act
	dv := "default-value"
	parser.StringVar(&testValue, "test-flag", dv, "test usage")

	// Assert
	// フラグが正常に定義されたことを確認
	// デフォルト値が設定されることを確認
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
	parser.BoolVar(&testValue, "test-bool-flag", false, "test bool usage")

	// Assert
	// フラグが正常に定義されたことを確認
	if testValue != false {
		t.Errorf("Expected initial value to be false, got %t", testValue)
	}
}

// TestStandardFlagParser_Parse_Normal はParseの正常系テスト
func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	// Arrange
	// 新しいFlagSetを作成してテスト用に使用
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

// TestStandardFlagParser_Parse_InvalidFlag は無効なフラグの場合のテスト
func TestStandardFlagParser_Parse_InvalidFlag(t *testing.T) {
	// Arrange
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)

	// 定義されていないフラグを含む引数
	testArgs := []string{"-undefined-flag=value"}

	// Act
	err := testFlagSet.Parse(testArgs)

	// Assert
	if err == nil {
		t.Error("Expected error for undefined flag, got nil")
	}
}

// TestStandardFlagParser_Args_Normal はArgsの正常系テスト
func TestStandardFlagParser_Args_Normal(t *testing.T) {
	// Arrange
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	// フラグと非フラグ引数を含むテストケース
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

// TestStandardFlagParser_Args_NoArgs は引数がない場合のテスト
func TestStandardFlagParser_Args_NoArgs(t *testing.T) {
	// Arrange
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	testArgs := []string{"-flag=value"}

	var flagValue string
	parser.StringVar(&flagValue, "flag", "", "test flag")

	// Act
	err := testFlagSet.Parse(testArgs)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	args := parser.Args()

	// Assert
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got %d", len(args))
	}
}

// TestNewStandardFlagParser_Normal はNewStandardFlagParserの正常系テスト
func TestNewStandardFlagParser_Normal(t *testing.T) {
	// Act
	parser := NewStandardFlagParser()

	// Assert
	if parser == nil {
		t.Error("Expected parser to be non-nil")
	}

	if parser.flagSet == nil {
		t.Error("Expected flagSet to be non-nil")
	}

	// flag.CommandLineが設定されていることを確認
	if parser.flagSet != flag.CommandLine {
		t.Error("Expected flagSet to be flag.CommandLine")
	}
}

// TestStandardFlagParser_Integration は統合テスト
func TestStandardFlagParser_Integration(t *testing.T) {
	// Arrange
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// テスト用の引数を設定
	os.Args = []string{"test-program", "-test-string=integration-test", "-test-bool"}

	// 新しいFlagSetを作成（グローバルなflag.CommandLineを汚染しないため）
	testFlagSet := flag.NewFlagSet("integration-test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	var testString string
	var testBool bool

	// Act
	parser.StringVar(&testString, "test-string", "default", "test string flag")
	parser.BoolVar(&testBool, "test-bool", false, "test bool flag")

	err := testFlagSet.Parse(os.Args[1:])

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if testString != "integration-test" {
		t.Errorf("Expected test-string to be 'integration-test', got %s", testString)
	}

	if testBool != true {
		t.Errorf("Expected test-bool to be true, got %t", testBool)
	}
}
