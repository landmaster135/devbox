package config

import (
	"os"
	"testing"
)

func TestNewStandardFlagParser(t *testing.T) {
	parser := NewStandardFlagParser()

	if parser == nil {
		t.Fatal("Expected parser to be non-nil")
	}

	if parser.flagSet == nil {
		t.Error("Expected flagSet to be initialized")
	}

	if parser.flagSet.Name() != os.Args[0] {
		t.Errorf("Expected flagSet name to be %s, got %s", os.Args[0], parser.flagSet.Name())
	}
}

func TestStandardFlagParser_StringVar(t *testing.T) {
	parser := NewStandardFlagParser()
	var testString string

	parser.StringVar(&testString, "test", "default", "test usage")

	// デフォルト値が設定されているか確認
	if testString != "default" {
		t.Errorf("Expected default value 'default', got '%s'", testString)
	}

	// フラグが定義されているか確認
	flag := parser.flagSet.Lookup("test")
	if flag == nil {
		t.Error("Expected flag 'test' to be defined")
		return
	}

	if flag.Usage != "test usage" {
		t.Errorf("Expected usage 'test usage', got '%s'", flag.Usage)
	}
}

func TestStandardFlagParser_BoolVar(t *testing.T) {
	parser := NewStandardFlagParser()
	var testBool bool

	parser.BoolVar(&testBool, "test-bool", true, "test bool usage")

	// デフォルト値が設定されているか確認
	if !testBool {
		t.Error("Expected default value true, got false")
	}

	// フラグが定義されているか確認
	flag := parser.flagSet.Lookup("test-bool")
	if flag == nil {
		t.Error("Expected flag 'test-bool' to be defined")
		return
	}

	if flag.Usage != "test bool usage" {
		t.Errorf("Expected usage 'test bool usage', got '%s'", flag.Usage)
	}
}

func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// テスト用の引数を設定
	os.Args = []string{"test-program", "-test=value", "-test-bool=false"}

	parser := NewStandardFlagParser()
	var testString string
	var testBool bool

	parser.StringVar(&testString, "test", "default", "test usage")
	parser.BoolVar(&testBool, "test-bool", true, "test bool usage")

	err := parser.Parse()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if testString != "value" {
		t.Errorf("Expected 'value', got '%s'", testString)
	}

	if testBool {
		t.Error("Expected false, got true")
	}
}

func TestStandardFlagParser_Parse_WithInvalidFlag(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// 無効なフラグを含むテスト用の引数を設定
	os.Args = []string{"test-program", "-invalid-flag=value"}

	parser := NewStandardFlagParser()
	var testString string

	parser.StringVar(&testString, "test", "default", "test usage")

	err := parser.Parse()

	if err == nil {
		t.Error("Expected error for invalid flag, but got none")
	}
}

func TestStandardFlagParser_Args(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// フラグと残りの引数を含むテスト用の引数を設定
	os.Args = []string{"test-program", "-test=value", "arg1", "arg2"}

	parser := NewStandardFlagParser()
	var testString string

	parser.StringVar(&testString, "test", "default", "test usage")

	err := parser.Parse()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	args := parser.Args()
	expectedArgs := []string{"arg1", "arg2"}

	if len(args) != len(expectedArgs) {
		t.Errorf("Expected %d args, got %d", len(expectedArgs), len(args))
		return
	}

	for i, arg := range args {
		if arg != expectedArgs[i] {
			t.Errorf("Expected arg[%d] to be '%s', got '%s'", i, expectedArgs[i], arg)
		}
	}
}

func TestStandardFlagParser_Parse_EmptyArgs(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// 引数なしのテスト用の設定
	os.Args = []string{"test-program"}

	parser := NewStandardFlagParser()
	var testString string

	parser.StringVar(&testString, "test", "default", "test usage")

	err := parser.Parse()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// デフォルト値が保持されているか確認
	if testString != "default" {
		t.Errorf("Expected default value 'default', got '%s'", testString)
	}

	args := parser.Args()
	if len(args) != 0 {
		t.Errorf("Expected 0 args, got %d", len(args))
	}
}

func TestStandardFlagParser_Parse_ContinueOnError(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// 無効なフラグを含むテスト用の引数を設定
	os.Args = []string{"test-program", "-unknown"}

	parser := NewStandardFlagParser()

	// ContinueOnErrorが設定されているため、エラーが返される
	err := parser.Parse()

	if err == nil {
		t.Error("Expected error for unknown flag, but got none")
	}

	// エラーメッセージを確認
	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}
