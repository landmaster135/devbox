package config

import (
	"flag"
	"os"
	"testing"
)

// TestNewStandardFlagParser_Normal はNewStandardFlagParser関数の正常系テスト
func TestNewStandardFlagParser_Normal(t *testing.T) {
	// Act
	parser := NewStandardFlagParser()

	// Assert
	if parser == nil {
		t.Fatal("Expected non-nil StandardFlagParser")
	}
	if parser.flagSet == nil {
		t.Error("Expected non-nil flagSet")
	}
	if parser.flagSet != flag.CommandLine {
		t.Error("Expected flagSet to be flag.CommandLine")
	}
}

// TestStandardFlagParser_StringVar はStringVarメソッドのテスト
func TestStandardFlagParser_StringVar(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue string
	expectedDefault := "default_value"
	flagName := "test-string-flag"
	usage := "Test string flag"

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// Act
	parser.StringVar(&testValue, flagName, expectedDefault, usage)

	// Assert
	if testValue != expectedDefault {
		t.Errorf("Expected testValue to be %s, got %s", expectedDefault, testValue)
	}

	// フラグが正しく登録されているかチェック
	testFlag := testFlagSet.Lookup(flagName)
	if testFlag == nil {
		t.Fatalf("Expected flag %s to be registered", flagName)
	}
	if testFlag.Usage != usage {
		t.Errorf("Expected flag usage to be %s, got %s", usage, testFlag.Usage)
	}
}

// TestStandardFlagParser_BoolVar はBoolVarメソッドのテスト
func TestStandardFlagParser_BoolVar(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var testValue bool
	expectedDefault := true
	flagName := "test-bool-flag"
	usage := "Test bool flag"

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// Act
	parser.BoolVar(&testValue, flagName, expectedDefault, usage)

	// Assert
	if testValue != expectedDefault {
		t.Errorf("Expected testValue to be %t, got %t", expectedDefault, testValue)
	}

	// フラグが正しく登録されているかチェック
	testFlag := testFlagSet.Lookup(flagName)
	if testFlag == nil {
		t.Fatalf("Expected flag %s to be registered", flagName)
	}
	if testFlag.Usage != usage {
		t.Errorf("Expected flag usage to be %s, got %s", usage, testFlag.Usage)
	}
}

// TestStandardFlagParser_Parse_Normal はParseメソッドの正常系テスト
func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var stringValue string
	var boolValue bool

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	parser.StringVar(&stringValue, "string-flag", "default", "Test string flag")
	parser.BoolVar(&boolValue, "bool-flag", false, "Test bool flag")

	// os.Argsを一時的に変更
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()
	os.Args = []string{"test-program", "-string-flag=test-value", "-bool-flag=true"}

	// Act
	err := parser.Parse()

	// Assert
	if err != nil {
		t.Errorf("Parse returned error: %v", err)
	}
	if stringValue != "test-value" {
		t.Errorf("Expected stringValue to be 'test-value', got %s", stringValue)
	}
	if !boolValue {
		t.Errorf("Expected boolValue to be true, got %t", boolValue)
	}
}

// TestStandardFlagParser_Parse_WithCustomArgs はカスタム引数でのParseメソッドのテスト
func TestStandardFlagParser_Parse_WithCustomArgs(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var stringValue string

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	parser.StringVar(&stringValue, "test-flag", "default", "Test flag")

	// カスタム引数を直接パース
	args := []string{"-test-flag=custom-value"}

	// Act
	err := testFlagSet.Parse(args)

	// Assert
	if err != nil {
		t.Errorf("Parse returned error: %v", err)
	}
	if stringValue != "custom-value" {
		t.Errorf("Expected stringValue to be 'custom-value', got %s", stringValue)
	}
}

// TestStandardFlagParser_Parse_InvalidFlag は無効なフラグでのParseメソッドのテスト
func TestStandardFlagParser_Parse_InvalidFlag(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// 無効な引数を直接パース
	args := []string{"-invalid-flag=value"}

	// Act
	err := testFlagSet.Parse(args)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid flag, got nil")
	}
}

// TestStandardFlagParser_Args はArgsメソッドのテスト
func TestStandardFlagParser_Args_Normal(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// フラグと非フラグ引数を含む引数をパース
	args := []string{"-flag=value", "arg1", "arg2"}
	var flagValue string
	parser.StringVar(&flagValue, "flag", "default", "Test flag")

	err := testFlagSet.Parse(args)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Act
	remainingArgs := parser.Args()

	// Assert
	expectedArgs := []string{"arg1", "arg2"}
	if len(remainingArgs) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(remainingArgs))
	}
	for i, expected := range expectedArgs {
		if i >= len(remainingArgs) || remainingArgs[i] != expected {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expected, remainingArgs[i])
		}
	}
}

// TestStandardFlagParser_Args_NoArgs は引数がない場合のArgsメソッドのテスト
func TestStandardFlagParser_Args_NoArgs(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// フラグのみの引数をパース
	args := []string{"-flag=value"}
	var flagValue string
	parser.StringVar(&flagValue, "flag", "default", "Test flag")

	err := testFlagSet.Parse(args)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Act
	remainingArgs := parser.Args()

	// Assert
	if len(remainingArgs) != 0 {
		t.Errorf("Expected 0 args, got %d", len(remainingArgs))
	}
}

// TestStandardFlagParser_Integration は統合テスト
func TestStandardFlagParser_Integration(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var operation string
	var year string
	var help bool

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	parser.StringVar(&operation, "operation", "", "Operation type")
	parser.StringVar(&operation, "o", "", "Operation short form")
	parser.StringVar(&year, "year", "2025", "Year value")
	parser.StringVar(&year, "y", "2025", "Year short form")
	parser.BoolVar(&help, "help", false, "Show help")
	parser.BoolVar(&help, "h", false, "Help short form")

	// テスト引数
	args := []string{"-o=add", "-y=2024", "-h=true", "extra", "args"}

	// Act
	err := testFlagSet.Parse(args)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	remainingArgs := parser.Args()

	// Assert
	if operation != "add" {
		t.Errorf("Expected operation to be 'add', got %s", operation)
	}
	if year != "2024" {
		t.Errorf("Expected year to be '2024', got %s", year)
	}
	if !help {
		t.Errorf("Expected help to be true, got %t", help)
	}

	expectedArgs := []string{"extra", "args"}
	if len(remainingArgs) != len(expectedArgs) {
		t.Errorf("Expected %d remaining args, got %d", len(expectedArgs), len(remainingArgs))
	}
	for i, expected := range expectedArgs {
		if i >= len(remainingArgs) || remainingArgs[i] != expected {
			t.Errorf("Expected remaining arg[%d] to be %s, got %s", i, expected, remainingArgs[i])
		}
	}
}

// TestStandardFlagParser_MultipleStringVarCalls は同じ変数に対する複数のStringVar呼び出しのテスト
func TestStandardFlagParser_MultipleStringVarCalls(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var value string

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// 同じ変数に対して複数のフラグを設定
	parser.StringVar(&value, "long-flag", "default", "Long flag name")
	parser.StringVar(&value, "s", "default", "Short flag name")

	// 短縮形フラグを使用
	args := []string{"-s=short-value"}

	// Act
	err := testFlagSet.Parse(args)

	// Assert
	if err != nil {
		t.Errorf("Parse returned error: %v", err)
	}
	if value != "short-value" {
		t.Errorf("Expected value to be 'short-value', got %s", value)
	}
}

// TestStandardFlagParser_MultipleBoolVarCalls は同じ変数に対する複数のBoolVar呼び出しのテスト
func TestStandardFlagParser_MultipleBoolVarCalls(t *testing.T) {
	// Arrange
	parser := NewStandardFlagParser()
	var value bool

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// 同じ変数に対して複数のフラグを設定
	parser.BoolVar(&value, "verbose", false, "Verbose output")
	parser.BoolVar(&value, "v", false, "Verbose short form")

	// 短縮形フラグを使用
	args := []string{"-v=true"}

	// Act
	err := testFlagSet.Parse(args)

	// Assert
	if err != nil {
		t.Errorf("Parse returned error: %v", err)
	}
	if !value {
		t.Errorf("Expected value to be true, got %t", value)
	}
}
