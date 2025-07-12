package config

import (
	"os"
	"testing"
)

// TestNewStandardOSArgs_Normal はNewStandardOSArgs()の正常系をテストします
func TestNewStandardOSArgs_Normal(t *testing.T) {
	osArgs := NewStandardOSArgs()

	if osArgs == nil {
		t.Error("NewStandardOSArgs() should not return nil")
		return
	}

	// 型が正しいかチェック
	if _, ok := interface{}(osArgs).(*StandardOSArgs); !ok {
		t.Error("NewStandardOSArgs() should return *StandardOSArgs")
	}
}

// TestStandardOSArgs_Args_Normal はArgs()の正常系をテストします
func TestStandardOSArgs_Args_Normal(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// テスト用の引数を設定
	testArgs := []string{"program", "arg1", "arg2", "arg3"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()
	args := osArgs.Args()

	if len(args) != len(testArgs) {
		t.Errorf("Args() length = %d, expected %d", len(args), len(testArgs))
		return
	}

	for i, arg := range args {
		if arg != testArgs[i] {
			t.Errorf("Args()[%d] = %s, expected %s", i, arg, testArgs[i])
		}
	}
}

// TestStandardOSArgs_Args_Empty は引数が空の場合のArgs()をテストします
func TestStandardOSArgs_Args_Empty(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 空の引数を設定
	os.Args = []string{}

	osArgs := NewStandardOSArgs()
	args := osArgs.Args()

	if len(args) != 0 {
		t.Errorf("Args() length = %d, expected 0", len(args))
	}
}

// TestStandardOSArgs_Args_SingleProgram はプログラム名のみの場合のArgs()をテストします
func TestStandardOSArgs_Args_SingleProgram(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// プログラム名のみを設定
	testArgs := []string{"program"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()
	args := osArgs.Args()

	if len(args) != 1 {
		t.Errorf("Args() length = %d, expected 1", len(args))
		return
	}

	if args[0] != "program" {
		t.Errorf("Args()[0] = %s, expected 'program'", args[0])
	}
}

// TestStandardOSArgs_Args_WithFlags はフラグを含む引数の場合のArgs()をテストします
func TestStandardOSArgs_Args_WithFlags(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// フラグを含む引数を設定
	testArgs := []string{"program", "-flag1", "value1", "--flag2=value2", "arg1"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()
	args := osArgs.Args()

	if len(args) != len(testArgs) {
		t.Errorf("Args() length = %d, expected %d", len(args), len(testArgs))
		return
	}

	for i, arg := range args {
		if arg != testArgs[i] {
			t.Errorf("Args()[%d] = %s, expected %s", i, arg, testArgs[i])
		}
	}
}

// TestStandardOSArgs_Args_WithSpecialCharacters は特殊文字を含む引数の場合のArgs()をテストします
func TestStandardOSArgs_Args_WithSpecialCharacters(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 特殊文字を含む引数を設定
	testArgs := []string{"program", "arg with spaces", "arg/with/slashes", "arg=with=equals"}
	os.Args = testArgs

	osArgs := NewStandardOSArgs()
	args := osArgs.Args()

	if len(args) != len(testArgs) {
		t.Errorf("Args() length = %d, expected %d", len(args), len(testArgs))
		return
	}

	for i, arg := range args {
		if arg != testArgs[i] {
			t.Errorf("Args()[%d] = %s, expected %s", i, arg, testArgs[i])
		}
	}
}

// TestStandardOSArgs_Interface はStandardOSArgsがOSArgsインターフェースを実装しているかテストします
func TestStandardOSArgs_Interface(t *testing.T) {
	var _ OSArgs = &StandardOSArgs{}
	var _ OSArgs = NewStandardOSArgs()

	// インターフェースメソッドが呼び出せるかテスト
	osArgs := NewStandardOSArgs()
	args := osArgs.Args()

	// argsがnilでないことを確認
	if args == nil {
		t.Error("Args() should not return nil")
	}
}
