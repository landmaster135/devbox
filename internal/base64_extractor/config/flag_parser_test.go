package config

import (
	"flag"
	"os"
	"testing"
)

// TestRealFlagParser はRealFlagParserのテストクラス
type TestRealFlagParser struct{}

// TestNewRealFlagParser_Normal はNewRealFlagParserの正常系テスト
func TestNewRealFlagParser_Normal(t *testing.T) {
	// Act
	parser := NewRealFlagParser()

	// Assert
	if parser == nil {
		t.Fatal("parserがnilです")
		return
	}
	if parser.flagSet == nil {
		t.Fatal("flagSetがnilです")
	}
	if parser.flagSet != flag.CommandLine {
		t.Error("flagSetがflag.CommandLineと異なります")
	}
}

// TestRealFlagParser_StringVar_Normal はStringVarの正常系テスト
func TestRealFlagParser_StringVar_Normal(t *testing.T) {
	// Arrange
	parser := NewRealFlagParser()
	var testValue string
	name := "test-string"
	defaultValue := "default"
	usage := "test usage"

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// Act
	parser.StringVar(&testValue, name, defaultValue, usage)

	// Assert
	if testValue != defaultValue {
		t.Errorf("testValueが期待値と異なります。期待値: %s, 実際: %s", defaultValue, testValue)
	}

	// フラグが登録されているかを確認
	flag := testFlagSet.Lookup(name)
	if flag == nil {
		t.Fatalf("フラグ '%s' が登録されていません", name)
	}
	if flag.Usage != usage {
		t.Errorf("フラグのUsageが期待値と異なります。期待値: %s, 実際: %s", usage, flag.Usage)
	}
}

// TestRealFlagParser_BoolVar_Normal はBoolVarの正常系テスト
func TestRealFlagParser_BoolVar_Normal(t *testing.T) {
	// Arrange
	parser := NewRealFlagParser()
	var testValue bool
	name := "test-bool"
	defaultValue := true
	usage := "test bool usage"

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// Act
	parser.BoolVar(&testValue, name, defaultValue, usage)

	// Assert
	if testValue != defaultValue {
		t.Errorf("testValueが期待値と異なります。期待値: %t, 実際: %t", defaultValue, testValue)
	}

	// フラグが登録されているかを確認
	flag := testFlagSet.Lookup(name)
	if flag == nil {
		t.Fatalf("フラグ '%s' が登録されていません", name)
		return
	}
	if flag.Usage != usage {
		t.Errorf("フラグのUsageが期待値と異なります。期待値: %s, 実際: %s", usage, flag.Usage)
	}
}

// TestRealFlagParser_Parse_Normal はParseの正常系テスト
func TestRealFlagParser_Parse_Normal(t *testing.T) {
	// Arrange
	parser := NewRealFlagParser()
	var testString string
	var testBool bool

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	parser.StringVar(&testString, "string", "default", "string flag")
	parser.BoolVar(&testBool, "bool", false, "bool flag")

	// テスト用の引数を設定
	testArgs := []string{"-string=test", "-bool=true"}

	// os.Argsを一時的に変更
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = append([]string{"test"}, testArgs...)

	// Act
	err := parser.Parse()

	// Assert
	if err != nil {
		t.Fatalf("Parseでエラーが発生しました: %v", err)
	}
}

// TestRealFlagParser_Parse_WithInvalidFlag は無効なフラグでのテスト
func TestRealFlagParser_Parse_WithInvalidFlag(t *testing.T) {
	// Arrange
	parser := NewRealFlagParser()

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

	// テスト用の無効な引数を設定
	testArgs := []string{"-invalid-flag=value"}

	// os.Argsを一時的に変更
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()
	os.Args = append([]string{"test"}, testArgs...)

	// Act
	err := parser.Parse()

	// Assert
	if err == nil {
		t.Fatal("エラーが発生するべきです")
	}
}

// TestRealFlagParser_StringVar_MultipleFlags は複数のStringVarのテスト
func TestRealFlagParser_StringVar_MultipleFlags(t *testing.T) {
	// Arrange
	parser := NewRealFlagParser()
	var value1, value2 string

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

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

	// 両方のフラグが登録されているかを確認
	flag1 := testFlagSet.Lookup("flag1")
	flag2 := testFlagSet.Lookup("flag2")
	if flag1 == nil {
		t.Error("flag1が登録されていません")
	}
	if flag2 == nil {
		t.Error("flag2が登録されていません")
	}
}

// TestRealFlagParser_BoolVar_MultipleFlags は複数のBoolVarのテスト
func TestRealFlagParser_BoolVar_MultipleFlags(t *testing.T) {
	// Arrange
	parser := NewRealFlagParser()
	var value1, value2 bool

	// 新しいFlagSetを作成してテスト用に使用
	testFlagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	parser.flagSet = testFlagSet

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

	// 両方のフラグが登録されているかを確認
	flag1 := testFlagSet.Lookup("bool1")
	flag2 := testFlagSet.Lookup("bool2")
	if flag1 == nil {
		t.Error("bool1が登録されていません")
	}
	if flag2 == nil {
		t.Error("bool2が登録されていません")
	}
}
