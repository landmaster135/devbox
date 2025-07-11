package config

import (
	"os"
	"testing"
)

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
// TestStandardOSArgs_Args_Normal はArgsの正常系テスト
func TestStandardOSArgs_Args_Normal(t *testing.T) {
	// Arrange
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// テスト用の引数を設定
	testArgs := []string{"test-program", "arg1", "arg2", "arg3"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()

	// Act
	result := osArgs.Args()

	// Assert
	if len(result) != len(testArgs) {
		t.Errorf("Expected %d args, got %d", len(testArgs), len(result))
	}

	for i, expected := range testArgs {
		if i >= len(result) || result[i] != expected {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expected, result[i])
		}
	}
}

// TestStandardOSArgs_Args_EmptyArgs は引数が空の場合のテスト
func TestStandardOSArgs_Args_EmptyArgs(t *testing.T) {
	// Arrange
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 空の引数を設定
	os.Args = []string{}

	osArgs := NewStandardOSArgs()

	// Act
	result := osArgs.Args()

	// Assert
	if len(result) != 0 {
		t.Errorf("Expected 0 args, got %d", len(result))
	}
}

// TestStandardOSArgs_Args_SingleArg は単一引数の場合のテスト
func TestStandardOSArgs_Args_SingleArg(t *testing.T) {
	// Arrange
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 単一引数を設定
	testArgs := []string{"test-program"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()

	// Act
	result := osArgs.Args()

	// Assert
	if len(result) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(result))
	}

	if result[0] != "test-program" {
		t.Errorf("Expected arg[0] to be 'test-program', got %s", result[0])
	}
}

// TestNewStandardOSArgs_Normal はNewStandardOSArgsの正常系テスト
func TestNewStandardOSArgs_Normal(t *testing.T) {
	// Act
	osArgs := NewStandardOSArgs()

	// Assert
	if osArgs == nil {
		t.Error("Expected osArgs to be non-nil")
	}

	// 型が正しいことを確認（インターフェースとして使用可能かチェック）
	var _ OSArgs = osArgs
}

// TestStandardOSArgs_Args_WithFlags はフラグ付き引数の場合のテスト
func TestStandardOSArgs_Args_WithFlags(t *testing.T) {
	// Arrange
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// フラグ付き引数を設定
	testArgs := []string{"test-program", "-flag1=value1", "-flag2", "arg1", "arg2"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()

	// Act
	result := osArgs.Args()

	// Assert
	if len(result) != len(testArgs) {
		t.Errorf("Expected %d args, got %d", len(testArgs), len(result))
	}

	for i, expected := range testArgs {
		if i >= len(result) || result[i] != expected {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expected, result[i])
		}
	}
}

// TestStandardOSArgs_Args_WithSpecialCharacters は特殊文字を含む引数の場合のテスト
func TestStandardOSArgs_Args_WithSpecialCharacters(t *testing.T) {
	// Arrange
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 特殊文字を含む引数を設定
	testArgs := []string{"test-program", "arg with spaces", "arg/with/slashes", "arg-with-dashes"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()

	// Act
	result := osArgs.Args()

	// Assert
	if len(result) != len(testArgs) {
		t.Errorf("Expected %d args, got %d", len(testArgs), len(result))
	}

	for i, expected := range testArgs {
		if i >= len(result) || result[i] != expected {
			t.Errorf("Expected arg[%d] to be %s, got %s", i, expected, result[i])
		}
	}
}

// TestStandardOSArgs_Args_Consistency は複数回呼び出しの一貫性テスト
func TestStandardOSArgs_Args_Consistency(t *testing.T) {
	// Arrange
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// テスト用の引数を設定
	testArgs := []string{"test-program", "arg1", "arg2"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()

	// Act
	result1 := osArgs.Args()
	result2 := osArgs.Args()

	// Assert
	if len(result1) != len(result2) {
		t.Errorf("Expected consistent results, got %d and %d", len(result1), len(result2))
	}

	for i := 0; i < len(result1) && i < len(result2); i++ {
		if result1[i] != result2[i] {
			t.Errorf("Expected consistent arg[%d], got %s and %s", i, result1[i], result2[i])
		}
	}
}
