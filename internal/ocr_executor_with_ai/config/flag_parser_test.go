package config

import (
	"testing"
)

// TestStandardFlagParser はStandardFlagParser構造体のテストクラス
type TestStandardFlagParser struct{}

func (t *TestStandardFlagParser) TestStandardFlagParser_AllMethods_Normal(test *testing.T) {
	parser := NewStandardFlagParser()

	if parser == nil {
		test.Fatal("NewStandardFlagParser() が nil を返しました")
	}

	// 各メソッドが正常に呼び出せることを確認
	var (
		testString  string
		testBool    bool
		testFloat64 float64
		testInt     int
	)

	// StringVar のテスト
	parser.StringVar(&testString, "test-string", "default", "test string usage")
	if testString != "default" {
		test.Errorf("StringVar() でデフォルト値が設定されませんでした: got %v, want %v", testString, "default")
	}

	// BoolVar のテスト
	parser.BoolVar(&testBool, "test-bool", true, "test bool usage")
	if testBool != true {
		test.Errorf("BoolVar() でデフォルト値が設定されませんでした: got %v, want %v", testBool, true)
	}

	// Float64Var のテスト
	parser.Float64Var(&testFloat64, "test-float64", 1.5, "test float64 usage")
	if testFloat64 != 1.5 {
		test.Errorf("Float64Var() でデフォルト値が設定されませんでした: got %v, want %v", testFloat64, 1.5)
	}

	// IntVar のテスト
	parser.IntVar(&testInt, "test-int", 42, "test int usage")
	if testInt != 42 {
		test.Errorf("IntVar() でデフォルト値が設定されませんでした: got %v, want %v", testInt, 42)
	}

	// Parse のテスト（エラーが発生しないことを確認）
	err := parser.Parse()
	if err != nil {
		test.Errorf("Parse() でエラーが発生しました: %v", err)
	}
}

// TestMockFlagParser はMockFlagParser構造体のテストクラス
type TestMockFlagParser struct{}

func (t *TestMockFlagParser) TestMockFlagParser_StringVar_Normal(test *testing.T) {
	const (
		testName         = "test-string"
		testDefaultValue = "default"
		testPresetValue  = "preset"
		testUsage        = "test usage"
	)

	tests := []struct {
		name           string
		presetValue    string
		hasPresetValue bool
		expectedValue  string
	}{
		{
			name:           "WithPresetValue",
			presetValue:    testPresetValue,
			hasPresetValue: true,
			expectedValue:  testPresetValue,
		},
		{
			name:           "WithoutPresetValue",
			presetValue:    "",
			hasPresetValue: false,
			expectedValue:  testDefaultValue,
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			mockParser := &MockFlagParser{
				stringValues: make(map[string]string),
				stringVars:   make(map[string]*string),
			}

			if tt.hasPresetValue {
				mockParser.stringValues[testName] = tt.presetValue
			}

			var testVar string
			mockParser.StringVar(&testVar, testName, testDefaultValue, testUsage)

			if testVar != tt.expectedValue {
				t.Errorf("StringVar() = %v, want %v", testVar, tt.expectedValue)
			}

			// stringVarsに正しく保存されているかを確認
			if storedPtr, exists := mockParser.stringVars[testName]; !exists {
				t.Error("stringVarsに変数が保存されていません")
			} else if storedPtr != &testVar {
				t.Error("stringVarsに正しいポインタが保存されていません")
			}
		})
	}
}

func (t *TestMockFlagParser) TestMockFlagParser_BoolVar_Normal(test *testing.T) {
	const (
		testName         = "test-bool"
		testDefaultValue = false
		testPresetValue  = true
		testUsage        = "test usage"
	)

	tests := []struct {
		name           string
		presetValue    bool
		hasPresetValue bool
		expectedValue  bool
	}{
		{
			name:           "WithPresetValue",
			presetValue:    testPresetValue,
			hasPresetValue: true,
			expectedValue:  testPresetValue,
		},
		{
			name:           "WithoutPresetValue",
			presetValue:    false,
			hasPresetValue: false,
			expectedValue:  testDefaultValue,
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			mockParser := &MockFlagParser{
				boolValues: make(map[string]bool),
				boolVars:   make(map[string]*bool),
			}

			if tt.hasPresetValue {
				mockParser.boolValues[testName] = tt.presetValue
			}

			var testVar bool
			mockParser.BoolVar(&testVar, testName, testDefaultValue, testUsage)

			if testVar != tt.expectedValue {
				t.Errorf("BoolVar() = %v, want %v", testVar, tt.expectedValue)
			}

			// boolVarsに正しく保存されているかを確認
			if storedPtr, exists := mockParser.boolVars[testName]; !exists {
				t.Error("boolVarsに変数が保存されていません")
			} else if storedPtr != &testVar {
				t.Error("boolVarsに正しいポインタが保存されていません")
			}
		})
	}
}

func (t *TestMockFlagParser) TestMockFlagParser_Float64Var_Normal(test *testing.T) {
	const (
		testName         = "test-float64"
		testDefaultValue = 1.0
		testPresetValue  = 2.5
		testUsage        = "test usage"
	)

	tests := []struct {
		name           string
		presetValue    float64
		hasPresetValue bool
		expectedValue  float64
	}{
		{
			name:           "WithPresetValue",
			presetValue:    testPresetValue,
			hasPresetValue: true,
			expectedValue:  testPresetValue,
		},
		{
			name:           "WithoutPresetValue",
			presetValue:    0.0,
			hasPresetValue: false,
			expectedValue:  testDefaultValue,
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			mockParser := &MockFlagParser{
				float64Values: make(map[string]float64),
				float64Vars:   make(map[string]*float64),
			}

			if tt.hasPresetValue {
				mockParser.float64Values[testName] = tt.presetValue
			}

			var testVar float64
			mockParser.Float64Var(&testVar, testName, testDefaultValue, testUsage)

			if testVar != tt.expectedValue {
				t.Errorf("Float64Var() = %v, want %v", testVar, tt.expectedValue)
			}

			// float64Varsに正しく保存されているかを確認
			if storedPtr, exists := mockParser.float64Vars[testName]; !exists {
				t.Error("float64Varsに変数が保存されていません")
			} else if storedPtr != &testVar {
				t.Error("float64Varsに正しいポインタが保存されていません")
			}
		})
	}
}

func (t *TestMockFlagParser) TestMockFlagParser_IntVar_Normal(test *testing.T) {
	const (
		testName         = "test-int"
		testDefaultValue = 10
		testPresetValue  = 20
		testUsage        = "test usage"
	)

	tests := []struct {
		name           string
		presetValue    int
		hasPresetValue bool
		expectedValue  int
	}{
		{
			name:           "WithPresetValue",
			presetValue:    testPresetValue,
			hasPresetValue: true,
			expectedValue:  testPresetValue,
		},
		{
			name:           "WithoutPresetValue",
			presetValue:    0,
			hasPresetValue: false,
			expectedValue:  testDefaultValue,
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			mockParser := &MockFlagParser{
				intValues: make(map[string]int),
				intVars:   make(map[string]*int),
			}

			if tt.hasPresetValue {
				mockParser.intValues[testName] = tt.presetValue
			}

			var testVar int
			mockParser.IntVar(&testVar, testName, testDefaultValue, testUsage)

			if testVar != tt.expectedValue {
				t.Errorf("IntVar() = %v, want %v", testVar, tt.expectedValue)
			}

			// intVarsに正しく保存されているかを確認
			if storedPtr, exists := mockParser.intVars[testName]; !exists {
				t.Error("intVarsに変数が保存されていません")
			} else if storedPtr != &testVar {
				t.Error("intVarsに正しいポインタが保存されていません")
			}
		})
	}
}

func (t *TestMockFlagParser) TestMockFlagParser_Parse_Normal(test *testing.T) {
	tests := []struct {
		name        string
		parseError  error
		expectError bool
	}{
		{
			name:        "ParseSuccess",
			parseError:  nil,
			expectError: false,
		},
		{
			name:        "ParseError",
			parseError:  &MockError{message: "parse error"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		test.Run(tt.name, func(t *testing.T) {
			mockParser := &MockFlagParser{
				parseError: tt.parseError,
			}

			err := mockParser.Parse()

			if tt.expectError {
				if err == nil {
					t.Error("Parse() でエラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("Parse() で予期しないエラーが発生しました: %v", err)
				}
			}
		})
	}
}

// テスト実行用の関数
func TestStandardFlagParser_AllMethods_Normal(t *testing.T) {
	testInstance := &TestStandardFlagParser{}
	testInstance.TestStandardFlagParser_AllMethods_Normal(t)
}

func TestMockFlagParser_StringVar_Normal(t *testing.T) {
	testInstance := &TestMockFlagParser{}
	testInstance.TestMockFlagParser_StringVar_Normal(t)
}

func TestMockFlagParser_BoolVar_Normal(t *testing.T) {
	testInstance := &TestMockFlagParser{}
	testInstance.TestMockFlagParser_BoolVar_Normal(t)
}

func TestMockFlagParser_Float64Var_Normal(t *testing.T) {
	testInstance := &TestMockFlagParser{}
	testInstance.TestMockFlagParser_Float64Var_Normal(t)
}

func TestMockFlagParser_IntVar_Normal(t *testing.T) {
	testInstance := &TestMockFlagParser{}
	testInstance.TestMockFlagParser_IntVar_Normal(t)
}

func TestMockFlagParser_Parse_Normal(t *testing.T) {
	testInstance := &TestMockFlagParser{}
	testInstance.TestMockFlagParser_Parse_Normal(t)
}

// MockError はテスト用のエラー実装
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}
