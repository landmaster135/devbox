package config

import (
	"errors"
	"testing"
)

// TestFlagParserInterface_MockImplementation はFlagParserインターフェースのモック実装テスト
func TestFlagParserInterface_MockImplementation(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()

	// Act & Assert - インターフェースが正しく実装されているかテスト
	var parser FlagParser = mockParser

	// StringVarメソッドのテスト
	var testString string
	parser.StringVar(&testString, "test", "default", "test usage")
	if testString != "default" {
		t.Errorf("Expected default value 'default', got '%s'", testString)
	}

	// BoolVarメソッドのテスト
	var testBool bool
	parser.BoolVar(&testBool, "test", true, "test usage")
	if !testBool {
		t.Error("Expected default value true, got false")
	}

	// Parseメソッドのテスト
	err := parser.Parse()
	if err != nil {
		t.Errorf("Expected no error from Parse, got %v", err)
	}

	// Argsメソッドのテスト
	args := parser.Args()
	if args == nil {
		t.Error("Expected non-nil args from Args method")
	}
}

// TestFlagParserInterface_StandardImplementation はFlagParserインターフェースの標準実装テスト
func TestFlagParserInterface_StandardImplementation(t *testing.T) {
	// Arrange
	standardParser := NewStandardFlagParser()

	// Act & Assert - インターフェースが正しく実装されているかテスト
	var parser FlagParser = standardParser

	// StringVarメソッドのテスト（ユニークな名前を使用）
	var testString string
	parser.StringVar(&testString, "test-string-std", "default", "test usage")
	if testString != "default" {
		t.Errorf("Expected default value 'default', got '%s'", testString)
	}

	// BoolVarメソッドのテスト（ユニークな名前を使用）
	var testBool bool
	parser.BoolVar(&testBool, "test-bool-std", true, "test usage")
	if !testBool {
		t.Error("Expected default value true, got false")
	}

	// Argsメソッドのテスト（Parseを呼ぶ前でも動作するはず）
	args := parser.Args()
	if args == nil {
		t.Error("Expected non-nil args from Args method")
	}
}

// TestMockFlagParser_SetMethods はMockFlagParserのセッターメソッドのテスト
func TestMockFlagParser_SetMethods(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	var testString string
	var testBool bool

	// StringVarを設定
	mockParser.StringVar(&testString, "testString", "default", "test usage")
	mockParser.BoolVar(&testBool, "testBool", false, "test usage")

	// Act
	mockParser.SetStringFlag("testString", "modified")
	mockParser.SetBoolFlag("testBool", true)
	mockParser.SetArgs([]string{"arg1", "arg2"})
	mockParser.SetParseError(errors.New("test error"))

	// Assert
	if testString != "modified" {
		t.Errorf("Expected modified string 'modified', got '%s'", testString)
	}
	if !testBool {
		t.Error("Expected modified bool true, got false")
	}

	args := mockParser.Args()
	expectedArgs := []string{"arg1", "arg2"}
	if len(args) != len(expectedArgs) {
		t.Errorf("Expected args length %d, got %d", len(expectedArgs), len(args))
	}
	for i, expected := range expectedArgs {
		if args[i] != expected {
			t.Errorf("Expected args[%d] '%s', got '%s'", i, expected, args[i])
		}
	}

	err := mockParser.Parse()
	if err == nil {
		t.Error("Expected error from Parse, got nil")
	}
	if err.Error() != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", err.Error())
	}
}

// TestMockFlagParser_NonExistentFlags は存在しないフラグの設定テスト
func TestMockFlagParser_NonExistentFlags(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()

	// Act - 存在しないフラグを設定しようとする
	mockParser.SetStringFlag("nonexistent", "value")
	mockParser.SetBoolFlag("nonexistent", true)

	// Assert - エラーが発生しないことを確認（単に無視される）
	// この動作は正常で、実際のフラグが定義されていない場合は無視される
	args := mockParser.Args()
	if args == nil {
		t.Error("Expected non-nil args even with nonexistent flags")
	}
}

// TestMockFlagParser_MultipleStringVars は複数の文字列フラグのテスト
func TestMockFlagParser_MultipleStringVars(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	var flag1, flag2, flag3 string

	// Act
	mockParser.StringVar(&flag1, "flag1", "default1", "usage1")
	mockParser.StringVar(&flag2, "flag2", "default2", "usage2")
	mockParser.StringVar(&flag3, "flag3", "default3", "usage3")

	mockParser.SetStringFlag("flag1", "value1")
	mockParser.SetStringFlag("flag2", "value2")
	mockParser.SetStringFlag("flag3", "value3")

	// Assert
	if flag1 != "value1" {
		t.Errorf("Expected flag1 'value1', got '%s'", flag1)
	}
	if flag2 != "value2" {
		t.Errorf("Expected flag2 'value2', got '%s'", flag2)
	}
	if flag3 != "value3" {
		t.Errorf("Expected flag3 'value3', got '%s'", flag3)
	}
}

// TestMockFlagParser_MultipleBoolVars は複数のブールフラグのテスト
func TestMockFlagParser_MultipleBoolVars(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	var flag1, flag2, flag3 bool

	// Act
	mockParser.BoolVar(&flag1, "flag1", false, "usage1")
	mockParser.BoolVar(&flag2, "flag2", true, "usage2")
	mockParser.BoolVar(&flag3, "flag3", false, "usage3")

	mockParser.SetBoolFlag("flag1", true)
	mockParser.SetBoolFlag("flag2", false)
	mockParser.SetBoolFlag("flag3", true)

	// Assert
	if !flag1 {
		t.Error("Expected flag1 true, got false")
	}
	if flag2 {
		t.Error("Expected flag2 false, got true")
	}
	if !flag3 {
		t.Error("Expected flag3 true, got false")
	}
}
