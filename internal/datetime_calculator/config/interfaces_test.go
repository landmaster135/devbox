package config

import (
	"testing"
)

// TestFlagParserInterface_StandardFlagParserImplementation はStandardFlagParserがFlagParserインターフェースを実装していることを確認するテスト
func TestFlagParserInterface_StandardFlagParserImplementation(t *testing.T) {
	// Arrange & Act
	var parser FlagParser = NewStandardFlagParser()

	// インターフェースのメソッドが呼び出せることを確認
	var testString string
	var testBool bool

	// StringVarメソッドの呼び出し確認
	parser.StringVar(&testString, "test", "default", "test usage")
	if testString != "default" {
		t.Errorf("Expected testString to be 'default', got %s", testString)
	}

	// BoolVarメソッドの呼び出し確認
	parser.BoolVar(&testBool, "test-bool", true, "test bool usage")
	if !testBool {
		t.Errorf("Expected testBool to be true, got %t", testBool)
	}

	// Argsメソッドの呼び出し確認
	args := parser.Args()
	if args == nil {
		t.Error("Expected Args() to return non-nil slice")
	}
}

// TestFlagParserInterface_MockImplementation はMockFlagParserがFlagParserインターフェースを実装していることを確認するテスト
func TestFlagParserInterface_MockImplementation(t *testing.T) {
	// Arrange & Act
	var parser FlagParser = NewMockFlagParser()

	// インターフェースのメソッドが呼び出せることを確認
	var testString string
	var testBool bool

	// StringVarメソッドの呼び出し確認
	parser.StringVar(&testString, "test", "default", "test usage")
	if testString != "default" {
		t.Errorf("Expected testString to be 'default', got %s", testString)
	}

	// BoolVarメソッドの呼び出し確認
	parser.BoolVar(&testBool, "test-bool", true, "test bool usage")
	if !testBool {
		t.Errorf("Expected testBool to be true, got %t", testBool)
	}

	// Parseメソッドの呼び出し確認
	err := parser.Parse()
	if err != nil {
		t.Errorf("Expected Parse() to return nil error, got %v", err)
	}

	// Argsメソッドの呼び出し確認
	args := parser.Args()
	if args == nil {
		t.Error("Expected Args() to return non-nil slice")
	}
}

// TestFlagParserInterface_PolymorphicUsage はポリモーフィックな使用のテスト
func TestFlagParserInterface_PolymorphicUsage(t *testing.T) {
	// テスト対象のパーサー実装
	parsers := []struct {
		name   string
		parser FlagParser
	}{
		{"StandardFlagParser", NewStandardFlagParser()},
		{"MockFlagParser", NewMockFlagParser()},
	}

	for _, tc := range parsers {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			parser := tc.parser
			var stringValue string
			var boolValue bool

			// Act - インターフェースを通じてメソッドを呼び出し
			parser.StringVar(&stringValue, "string-flag", "default-string", "String flag usage")
			parser.BoolVar(&boolValue, "bool-flag", false, "Bool flag usage")

			// Parse呼び出し（MockFlagParserの場合はエラーなし、StandardFlagParserの場合は実際の引数に依存）
			err := parser.Parse()

			// Args呼び出し
			args := parser.Args()

			// Assert
			if stringValue != "default-string" {
				t.Errorf("Expected stringValue to be 'default-string', got %s", stringValue)
			}
			if boolValue != false {
				t.Errorf("Expected boolValue to be false, got %t", boolValue)
			}

			// MockFlagParserの場合はエラーなし、StandardFlagParserの場合は実装に依存
			if tc.name == "MockFlagParser" && err != nil {
				t.Errorf("Expected MockFlagParser Parse() to return nil error, got %v", err)
			}

			if args == nil {
				t.Error("Expected Args() to return non-nil slice")
			}
		})
	}
}

// TestFlagParserInterface_MethodSignatures はインターフェースのメソッドシグネチャのテスト
func TestFlagParserInterface_MethodSignatures(t *testing.T) {
	// Arrange
	parser := NewMockFlagParser() // StandardFlagParserではなくMockFlagParserを使用してフラグ重複を避ける

	// Act & Assert - コンパイル時にメソッドシグネチャが正しいことを確認
	var stringValue string
	var boolValue bool

	// StringVarメソッドのシグネチャ確認
	parser.StringVar(&stringValue, "string-name", "value", "usage")

	// BoolVarメソッドのシグネチャ確認
	parser.BoolVar(&boolValue, "bool-name", true, "usage")

	// Parseメソッドのシグネチャ確認
	var err error = parser.Parse()
	_ = err

	// Argsメソッドのシグネチャ確認
	var args []string = parser.Args()
	_ = args
}

// TestFlagParserInterface_NilPointerHandling はnilポインタの処理テスト
func TestFlagParserInterface_NilPointerHandling(t *testing.T) {
	// Arrange
	parser := NewMockFlagParser()

	// Act & Assert - nilポインタを渡すとパニックが発生することを確認
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when passing nil pointer to StringVar, but no panic occurred")
		}
	}()

	// nilポインタでStringVarを呼び出し（パニックが発生することを期待）
	parser.StringVar(nil, "test", "default", "usage")
}

// TestFlagParserInterface_EmptyValues は空の値での処理テスト
func TestFlagParserInterface_EmptyValues(t *testing.T) {
	// Arrange
	parser := NewMockFlagParser()
	var stringValue string
	var boolValue bool

	// Act
	parser.StringVar(&stringValue, "", "", "")
	parser.BoolVar(&boolValue, "", false, "")

	// Assert - 空の値でも正常に動作することを確認
	if stringValue != "" {
		t.Errorf("Expected empty string, got %s", stringValue)
	}
	if boolValue != false {
		t.Errorf("Expected false, got %t", boolValue)
	}
}

// TestFlagParserInterface_MultipleCallsWithSameVariable は同じ変数への複数回呼び出しテスト
func TestFlagParserInterface_MultipleCallsWithSameVariable(t *testing.T) {
	// Arrange
	parser := NewMockFlagParser()
	var value string

	// Act - 同じ変数に対して複数回StringVarを呼び出し
	parser.StringVar(&value, "flag1", "default1", "usage1")
	parser.StringVar(&value, "flag2", "default2", "usage2")

	// Assert - MockFlagParserの実装では、非空の値が既に設定されている場合は保持される
	// フラグが設定されていない場合、最初に設定された非空のデフォルト値が保持される
	if value != "default1" {
		t.Errorf("Expected value to be 'default1', got %s", value)
	}
}

// TestFlagParserInterface_ConcurrentAccess は並行アクセスの基本テスト
func TestFlagParserInterface_ConcurrentAccess(t *testing.T) {
	// Arrange
	parser := NewMockFlagParser()

	// Act & Assert - 並行してメソッドを呼び出してもパニックしないことを確認
	done := make(chan bool, 2)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Concurrent StringVar caused panic: %v", r)
			}
			done <- true
		}()

		var value string
		parser.StringVar(&value, "concurrent1", "default1", "usage1")
	}()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Concurrent BoolVar caused panic: %v", r)
			}
			done <- true
		}()

		var value bool
		parser.BoolVar(&value, "concurrent2", false, "usage2")
	}()

	// 両方のgoroutineの完了を待つ
	<-done
	<-done
}
