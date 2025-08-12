package config

import (
	"flag"
	"os"
	"testing"
)

// TestNewStandardFlagParser_Normal はNewStandardFlagParser()の正常系をテストします
func TestNewStandardFlagParser_Normal(t *testing.T) {
	parser := NewStandardFlagParser()

	if parser == nil {
		t.Error("NewStandardFlagParser() should not return nil")
		return
	}

	if parser.flagSet == nil {
		t.Error("flagSet should not be nil")
	}

	if parser.flagSet != flag.CommandLine {
		t.Error("flagSet should be flag.CommandLine")
	}
}

// TestStandardFlagParser_StringVar_Normal はStringVar()の正常系をテストします
func TestStandardFlagParser_StringVar_Normal(t *testing.T) {
	// テスト用のflagSetを作成
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: flagSet,
	}

	var testValue string
	parser.StringVar(&testValue, "test-string", "default", "test usage")

	// フラグが正しく定義されているかチェック
	testFlag := flagSet.Lookup("test-string")
	if testFlag == nil {
		t.Error("StringVar() should define the flag")
		return
	}

	if testFlag.DefValue != "default" {
		t.Errorf("Default value = %s, expected 'default'", testFlag.DefValue)
	}

	if testFlag.Usage != "test usage" {
		t.Errorf("Usage = %s, expected 'test usage'", testFlag.Usage)
	}
}

// TestStandardFlagParser_BoolVar_Normal はBoolVar()の正常系をテストします
func TestStandardFlagParser_BoolVar_Normal(t *testing.T) {
	// テスト用のflagSetを作成
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: flagSet,
	}

	var testValue bool
	parser.BoolVar(&testValue, "test-bool", true, "test bool usage")

	// フラグが正しく定義されているかチェック
	testFlag := flagSet.Lookup("test-bool")
	if testFlag == nil {
		t.Error("BoolVar() should define the flag")
		return
	}

	if testFlag.DefValue != "true" {
		t.Errorf("Default value = %s, expected 'true'", testFlag.DefValue)
	}

	if testFlag.Usage != "test bool usage" {
		t.Errorf("Usage = %s, expected 'test bool usage'", testFlag.Usage)
	}
}

// TestStandardFlagParser_Parse_Normal はParse()の正常系をテストします
func TestStandardFlagParser_Parse_Normal(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// テスト用の引数を設定
	os.Args = []string{"program", "-test-string=value", "-test-bool=true"}

	// テスト用のflagSetを作成
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: flagSet,
	}

	var stringValue string
	var boolValue bool
	parser.StringVar(&stringValue, "test-string", "default", "test usage")
	parser.BoolVar(&boolValue, "test-bool", false, "test bool usage")

	err := parser.Parse()
	if err != nil {
		t.Errorf("Parse() returned unexpected error: %v", err)
	}

	if stringValue != "value" {
		t.Errorf("stringValue = %s, expected 'value'", stringValue)
	}

	if !boolValue {
		t.Error("boolValue should be true")
	}
}

// TestStandardFlagParser_Parse_Error はParse()のエラーケースをテストします
func TestStandardFlagParser_Parse_Error(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 無効な引数を設定
	os.Args = []string{"program", "-undefined-flag=value"}

	// テスト用のflagSetを作成
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: flagSet,
	}

	err := parser.Parse()
	if err == nil {
		t.Error("Parse() should return error for undefined flag")
	}
}

// TestStandardFlagParser_Args_Normal はArgs()の正常系をテストします
func TestStandardFlagParser_Args_Normal(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// テスト用の引数を設定（フラグ以外の引数を含む）
	os.Args = []string{"program", "-test-string=value", "arg1", "arg2"}

	// テスト用のflagSetを作成
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: flagSet,
	}

	var stringValue string
	parser.StringVar(&stringValue, "test-string", "default", "test usage")

	// まずParseを実行
	err := parser.Parse()
	if err != nil {
		t.Errorf("Parse() returned unexpected error: %v", err)
	}

	// 残りの引数を取得
	args := parser.Args()
	expectedArgs := []string{"arg1", "arg2"}

	if len(args) != len(expectedArgs) {
		t.Errorf("Args() length = %d, expected %d", len(args), len(expectedArgs))
	}

	for i, arg := range args {
		if arg != expectedArgs[i] {
			t.Errorf("Args()[%d] = %s, expected %s", i, arg, expectedArgs[i])
		}
	}
}

// TestStandardFlagParser_Args_Empty は引数がない場合のArgs()をテストします
func TestStandardFlagParser_Args_Empty(t *testing.T) {
	// 元のos.Argsを保存
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// フラグのみの引数を設定
	os.Args = []string{"program", "-test-string=value"}

	// テスト用のflagSetを作成
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser := &StandardFlagParser{
		flagSet: flagSet,
	}

	var stringValue string
	parser.StringVar(&stringValue, "test-string", "default", "test usage")

	// まずParseを実行
	err := parser.Parse()
	if err != nil {
		t.Errorf("Parse() returned unexpected error: %v", err)
	}

	// 残りの引数を取得
	args := parser.Args()

	if len(args) != 0 {
		t.Errorf("Args() length = %d, expected 0", len(args))
	}
}
