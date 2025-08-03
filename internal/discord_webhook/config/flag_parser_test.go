package config

import (
	"flag"
	"os"
	"testing"
)

// TestNewStandardFlagParser_Normal はStandardFlagParserの正常作成をテストする
func TestNewStandardFlagParser_Normal(t *testing.T) {
	parser := NewStandardFlagParser()

	if parser == nil {
		t.Errorf("NewStandardFlagParser() returned nil")
		return
	}

	if parser.flagSet == nil {
		t.Errorf("flagSet is nil")
	}

	if parser.flagSet != flag.CommandLine {
		t.Errorf("flagSet should be flag.CommandLine")
	}
}

// TestStandardFlagParser_StringVar はStringVarメソッドをテストする
func TestStandardFlagParser_StringVar(t *testing.T) {
	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	var testValue string
	parser.StringVar(&testValue, "test-flag", "default", "test usage")

	// フラグが正しく定義されているかチェック
	testFlag := testFlagSet.Lookup("test-flag")
	if testFlag == nil {
		t.Errorf("Flag 'test-flag' was not defined")
		return
	}

	if testFlag.DefValue != "default" {
		t.Errorf("Default value = %v, want %v", testFlag.DefValue, "default")
	}

	if testFlag.Usage != "test usage" {
		t.Errorf("Usage = %v, want %v", testFlag.Usage, "test usage")
	}
}

// TestStandardFlagParser_BoolVar はBoolVarメソッドをテストする
func TestStandardFlagParser_BoolVar(t *testing.T) {
	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	var testValue bool
	parser.BoolVar(&testValue, "test-bool", true, "test bool usage")

	// フラグが正しく定義されているかチェック
	testFlag := testFlagSet.Lookup("test-bool")
	if testFlag == nil {
		t.Errorf("Flag 'test-bool' was not defined")
		return
	}

	if testFlag.DefValue != "true" {
		t.Errorf("Default value = %v, want %v", testFlag.DefValue, "true")
	}

	if testFlag.Usage != "test bool usage" {
		t.Errorf("Usage = %v, want %v", testFlag.Usage, "test bool usage")
	}
}

// TestStandardFlagParser_Parse はParseメソッドをテストする
func TestStandardFlagParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "正常な引数",
			args:        []string{"-test-flag", "value"},
			expectError: false,
		},
		{
			name:        "引数なし",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "無効なフラグ",
			args:        []string{"-invalid-flag", "value"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 新しいFlagSetを作成してテスト用に使用
			testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			parser := &StandardFlagParser{
				flagSet: testFlagSet,
			}

			// テスト用フラグを定義
			var testValue string
			parser.StringVar(&testValue, "test-flag", "default", "test usage")

			// 元のos.Argsを保存
			originalArgs := os.Args
			defer func() {
				os.Args = originalArgs
			}()

			// テスト用の引数を設定
			os.Args = append([]string{"test"}, tt.args...)

			// Parseメソッドをテスト
			err := parser.Parse()

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
			}
		})
	}
}

// TestStandardFlagParser_Args はArgsメソッドをテストする
func TestStandardFlagParser_Args(t *testing.T) {
	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	// テスト用フラグを定義
	var testValue string
	parser.StringVar(&testValue, "test-flag", "default", "test usage")

	// テスト用の引数を設定してパース
	testArgs := []string{"-test-flag", "value", "arg1", "arg2"}
	err := testFlagSet.Parse(testArgs)
	if err != nil {
		t.Errorf("Parse failed: %v", err)
		return
	}

	// Argsメソッドをテスト
	args := parser.Args()
	expectedArgs := []string{"arg1", "arg2"}

	if len(args) != len(expectedArgs) {
		t.Errorf("Args length = %v, want %v", len(args), len(expectedArgs))
		return
	}

	for i, arg := range args {
		if arg != expectedArgs[i] {
			t.Errorf("Args[%d] = %v, want %v", i, arg, expectedArgs[i])
		}
	}
}

// TestStandardFlagParser_Integration は統合テストを行う
func TestStandardFlagParser_Integration(t *testing.T) {
	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	// 複数のフラグを定義
	var stringValue string
	var boolValue bool
	parser.StringVar(&stringValue, "string-flag", "default-string", "string flag usage")
	parser.BoolVar(&boolValue, "bool-flag", false, "bool flag usage")

	// テスト用の引数を設定してパース
	testArgs := []string{"-string-flag", "test-value", "-bool-flag", "remaining", "args"}
	err := testFlagSet.Parse(testArgs)
	if err != nil {
		t.Errorf("Parse failed: %v", err)
		return
	}

	// 値が正しく設定されているかチェック
	if stringValue != "test-value" {
		t.Errorf("stringValue = %v, want %v", stringValue, "test-value")
	}

	if !boolValue {
		t.Errorf("boolValue = %v, want %v", boolValue, true)
	}

	// 残りの引数をチェック
	args := parser.Args()
	expectedArgs := []string{"remaining", "args"}

	if len(args) != len(expectedArgs) {
		t.Errorf("Args length = %v, want %v", len(args), len(expectedArgs))
		return
	}

	for i, arg := range args {
		if arg != expectedArgs[i] {
			t.Errorf("Args[%d] = %v, want %v", i, arg, expectedArgs[i])
		}
	}
}

// TestStandardFlagParser_EmptyArgs は引数なしの場合をテストする
func TestStandardFlagParser_EmptyArgs(t *testing.T) {
	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: testFlagSet,
	}

	// テスト用フラグを定義
	var testValue string
	parser.StringVar(&testValue, "test-flag", "default", "test usage")

	// 引数なしでパース
	err := testFlagSet.Parse([]string{})
	if err != nil {
		t.Errorf("Parse failed: %v", err)
		return
	}

	// デフォルト値が設定されているかチェック
	if testValue != "default" {
		t.Errorf("testValue = %v, want %v", testValue, "default")
	}

	// 引数が空であることをチェック
	args := parser.Args()
	if len(args) != 0 {
		t.Errorf("Args should be empty, got %v", args)
	}
}
