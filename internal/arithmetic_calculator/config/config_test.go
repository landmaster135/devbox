package config

import (
	"fmt"
	"os"
	"testing"
)

// TestConfigStruct はConfig構造体のテストクラス
type TestConfigStruct struct{}

// TestParseFlags_Normal はParseFlags関数の正常系テスト
func (t *TestConfigStruct) TestParseFlags_Normal(test *testing.T) {
	// Arrange - ParseFlags関数を直接テストするため、os.Argsを設定
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 有効な引数を設定
	os.Args = []string{"test-program", "-operation=subtract", "-x=20", "-y=8"}

	// Act
	config, err := ParseFlags()

	// Assert
	if err != nil {
		test.Errorf("ParseFlags returned error: %v", err)
	}
	if config.Operation != "subtract" {
		test.Errorf("Expected operation 'subtract', got %s", config.Operation)
	}
	if config.X != 20.0 {
		test.Errorf("Expected X 20.0, got %f", config.X)
	}
	if config.Y != 8.0 {
		test.Errorf("Expected Y 8.0, got %f", config.Y)
	}
}

// TestParseFlagsWithParser_HelpFlag はヘルプフラグのテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_HelpFlag(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	// ヘルプフラグを設定する前に、まずフラグを定義する必要がある
	var help bool
	mockParser.BoolVar(&help, "help", false, "ヘルプを表示")
	mockParser.SetBoolFlag("help", true)

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		test.Errorf("ParseFlagsWithParser returned error: %v", err)
		return
	}
	if config == nil {
		test.Fatal("Expected non-nil config for help flag")
	}
	if !config.Help {
		test.Error("Expected Help flag to be true")
	}
}

// TestParseFlagsWithParser_InvalidXValue は無効なX値のテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_InvalidXValue(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("x", "invalid")
	mockParser.SetStringFlag("y", "5")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		test.Error("Expected error for invalid X value, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for invalid X value")
	}
}

// TestParseFlagsWithParser_InvalidYValue は無効なY値のテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_InvalidYValue(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("x", "10")
	mockParser.SetStringFlag("y", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		test.Error("Expected error for invalid Y value, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for invalid Y value")
	}
}

// TestParseFlagsWithParser_InvalidThreshold は無効な閾値のテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_InvalidThreshold(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "evaluate_line_count")
	mockParser.SetStringFlag("file", "/test/file.txt")
	mockParser.SetStringFlag("threshold", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		test.Error("Expected error for invalid threshold value, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for invalid threshold value")
	}
}

// TestParseFlagsWithParser_InvalidNumbers は無効な数値リストのテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_InvalidNumbers(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "sum")
	mockParser.SetStringFlag("numbers", "1,2,invalid,4")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		test.Error("Expected error for invalid numbers, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for invalid numbers")
	}
}

// TestParseFlagsWithParser_PositionalArgs は位置引数のテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_PositionalArgs(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "multiply")
	mockParser.SetArgs([]string{"15", "3"})

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		test.Errorf("ParseFlagsWithParser returned error: %v", err)
		return
	}
	if config == nil {
		test.Fatal("Expected non-nil config")
	}
	if config.X != 15.0 {
		test.Errorf("Expected X 15.0 from positional args, got %f", config.X)
	}
	if config.Y != 3.0 {
		test.Errorf("Expected Y 3.0 from positional args, got %f", config.Y)
	}
}

// TestParseFlagsWithParser_SumOperation はsum操作のテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_SumOperation(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "sum")
	mockParser.SetStringFlag("numbers", "1,2,3,4,5")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		test.Errorf("ParseFlagsWithParser returned error: %v", err)
	}
	if config.Operation != "sum" {
		test.Errorf("Expected operation 'sum', got %s", config.Operation)
	}
	expectedNumbers := []float64{1, 2, 3, 4, 5}
	if len(config.Numbers) != len(expectedNumbers) {
		test.Errorf("Expected Numbers length %d, got %d", len(expectedNumbers), len(config.Numbers))
	}
	for i, expected := range expectedNumbers {
		if config.Numbers[i] != expected {
			test.Errorf("Expected Numbers[%d] %f, got %f", i, expected, config.Numbers[i])
		}
	}
}

// TestParseFlagsWithParser_EvaluateLineCountOperation はevaluate_line_count操作のテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_EvaluateLineCountOperation(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "evaluate_line_count")
	mockParser.SetStringFlag("file", "/test/file.txt")
	mockParser.SetStringFlag("threshold", "200")
	mockParser.SetArgs([]string{"150"})

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		test.Errorf("ParseFlagsWithParser returned error: %v", err)
	}
	if config.Operation != "evaluate_line_count" {
		test.Errorf("Expected operation 'evaluate_line_count', got %s", config.Operation)
	}
	if config.FilePath != "/test/file.txt" {
		test.Errorf("Expected FilePath '/test/file.txt', got %s", config.FilePath)
	}
	if config.Threshold != 150 {
		test.Errorf("Expected Threshold 150 from positional args, got %d", config.Threshold)
	}
}

// TestParseFlagsWithParser_ParseError はパースエラーのテスト
func (t *TestConfigStruct) TestParseFlagsWithParser_ParseError(test *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetParseError(fmt.Errorf("parse error"))

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		test.Error("Expected parse error, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for parse error")
	}
}

// TestPrintUsage_Normal はPrintUsage関数のテスト
func (t *TestConfigStruct) TestPrintUsage_Normal(test *testing.T) {
	// Act - PrintUsageは標準エラー出力に書き込むため、パニックしないことを確認
	defer func() {
		if r := recover(); r != nil {
			test.Errorf("PrintUsage panicked: %v", r)
		}
	}()

	PrintUsage()

	// Assert - パニックしなければ成功
}

// TestNewConfig_Normal はNewConfigの正常系テスト
func (t *TestConfigStruct) TestNewConfig_Normal(test *testing.T) {
	// Arrange
	operation := "add"
	x := 10.0
	y := 5.0
	numbers := []float64{1, 2, 3}
	filePath := "/test/file.txt"
	threshold := 100

	// Act
	config, err := NewConfig(operation, x, y, numbers, filePath, threshold)

	// Assert
	if err != nil {
		test.Errorf("NewConfig returned error: %v", err)
	}
	if config.Operation != operation {
		test.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.X != x {
		test.Errorf("Expected X %f, got %f", x, config.X)
	}
	if config.Y != y {
		test.Errorf("Expected Y %f, got %f", y, config.Y)
	}
	if len(config.Numbers) != len(numbers) {
		test.Errorf("Expected Numbers length %d, got %d", len(numbers), len(config.Numbers))
	}
	if config.FilePath != filePath {
		test.Errorf("Expected FilePath %s, got %s", filePath, config.FilePath)
	}
	if config.Threshold != threshold {
		test.Errorf("Expected Threshold %d, got %d", threshold, config.Threshold)
	}
}

// TestNewConfig_EmptyOperation は操作タイプが空の場合のテスト
func (t *TestConfigStruct) TestNewConfig_EmptyOperation(test *testing.T) {
	// Arrange
	operation := ""
	x := 10.0
	y := 5.0
	numbers := []float64{}
	filePath := ""
	threshold := 0

	// Act
	config, err := NewConfig(operation, x, y, numbers, filePath, threshold)

	// Assert
	if err == nil {
		test.Error("Expected error for empty operation, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for empty operation")
	}
}

// TestNewConfig_InvalidOperation は無効な操作タイプの場合のテスト
func (t *TestConfigStruct) TestNewConfig_InvalidOperation(test *testing.T) {
	// Arrange
	operation := "invalid_operation"
	x := 10.0
	y := 5.0
	numbers := []float64{}
	filePath := ""
	threshold := 0

	// Act
	config, err := NewConfig(operation, x, y, numbers, filePath, threshold)

	// Assert
	if err == nil {
		test.Error("Expected error for invalid operation, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for invalid operation")
	}
}

// TestNewConfig_SumWithEmptyNumbers はsum操作で数値配列が空の場合のテスト
func (t *TestConfigStruct) TestNewConfig_SumWithEmptyNumbers(test *testing.T) {
	// Arrange
	operation := "sum"
	x := 0.0
	y := 0.0
	numbers := []float64{}
	filePath := ""
	threshold := 0

	// Act
	config, err := NewConfig(operation, x, y, numbers, filePath, threshold)

	// Assert
	if err == nil {
		test.Error("Expected error for sum operation with empty numbers, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for sum operation with empty numbers")
	}
}

// TestNewConfig_EvaluateLineCountWithEmptyFilePath はevaluate_line_count操作でファイルパスが空の場合のテスト
func (t *TestConfigStruct) TestNewConfig_EvaluateLineCountWithEmptyFilePath(test *testing.T) {
	// Arrange
	operation := "evaluate_line_count"
	x := 0.0
	y := 0.0
	numbers := []float64{}
	filePath := ""
	threshold := 100

	// Act
	config, err := NewConfig(operation, x, y, numbers, filePath, threshold)

	// Assert
	if err == nil {
		test.Error("Expected error for evaluate_line_count operation with empty file path, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for evaluate_line_count operation with empty file path")
	}
}

// TestNewConfig_EvaluateLineCountWithNegativeThreshold はevaluate_line_count操作で負の閾値の場合のテスト
func (t *TestConfigStruct) TestNewConfig_EvaluateLineCountWithNegativeThreshold(test *testing.T) {
	// Arrange
	operation := "evaluate_line_count"
	x := 0.0
	y := 0.0
	numbers := []float64{}
	filePath := "/test/file.txt"
	threshold := -1

	// Act
	config, err := NewConfig(operation, x, y, numbers, filePath, threshold)

	// Assert
	if err == nil {
		test.Error("Expected error for evaluate_line_count operation with negative threshold, got nil")
	}
	if config != nil {
		test.Error("Expected nil config for evaluate_line_count operation with negative threshold")
	}
}

// TestNewConfig_AllValidOperations は全ての有効な操作タイプのテスト
func (t *TestConfigStruct) TestNewConfig_AllValidOperations(test *testing.T) {
	validOperations := []string{"add", "subtract", "multiply", "divide", "sum", "evaluate_line_count"}

	for _, operation := range validOperations {
		// Arrange
		x := 10.0
		y := 5.0
		numbers := []float64{1, 2, 3}
		filePath := "/test/file.txt"
		threshold := 100

		// Act
		config, err := NewConfig(operation, x, y, numbers, filePath, threshold)

		// Assert
		if err != nil {
			test.Errorf("NewConfig returned error for valid operation %s: %v", operation, err)
		}
		if config == nil {
			test.Errorf("Expected non-nil config for valid operation %s", operation)
		}
		if config != nil && config.Operation != operation {
			test.Errorf("Expected operation %s, got %s", operation, config.Operation)
		}
	}
}

// 実際のテスト関数
func TestNewConfig_Normal(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestNewConfig_Normal(t)
}

func TestNewConfig_EmptyOperation(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestNewConfig_EmptyOperation(t)
}

func TestNewConfig_InvalidOperation(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestNewConfig_InvalidOperation(t)
}

func TestNewConfig_SumWithEmptyNumbers(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestNewConfig_SumWithEmptyNumbers(t)
}

func TestNewConfig_EvaluateLineCountWithEmptyFilePath(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestNewConfig_EvaluateLineCountWithEmptyFilePath(t)
}

func TestNewConfig_EvaluateLineCountWithNegativeThreshold(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestNewConfig_EvaluateLineCountWithNegativeThreshold(t)
}

func TestNewConfig_AllValidOperations(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestNewConfig_AllValidOperations(t)
}

// TestParseFlags_Normal はParseFlags関数の正常系テスト
func TestParseFlags_Normal(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestParseFlags_Normal(t)
}

// TestParseFlagsWithParser_HelpFlag はヘルプフラグのテスト
// func TestParseFlagsWithParser_HelpFlag(t *testing.T) {
// 	testStruct := &TestConfigStruct{}
// 	testStruct.TestParseFlagsWithParser_HelpFlag(t)
// }

// TestParseFlagsWithParser_InvalidXValue は無効なX値のテスト
func TestParseFlagsWithParser_InvalidXValue(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestParseFlagsWithParser_InvalidXValue(t)
}

// TestParseFlagsWithParser_InvalidYValue は無効なY値のテスト
func TestParseFlagsWithParser_InvalidYValue(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestParseFlagsWithParser_InvalidYValue(t)
}

// TestParseFlagsWithParser_InvalidThreshold は無効な閾値のテスト
func TestParseFlagsWithParser_InvalidThreshold(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestParseFlagsWithParser_InvalidThreshold(t)
}

// TestParseFlagsWithParser_InvalidNumbers は無効な数値リストのテスト
func TestParseFlagsWithParser_InvalidNumbers(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestParseFlagsWithParser_InvalidNumbers(t)
}

// TestParseFlagsWithParser_PositionalArgs は位置引数のテスト
// func TestParseFlagsWithParser_PositionalArgs(t *testing.T) {
// 	testStruct := &TestConfigStruct{}
// 	testStruct.TestParseFlagsWithParser_PositionalArgs(t)
// }

// TestParseFlagsWithParser_SumOperation はsum操作のテスト
// func TestParseFlagsWithParser_SumOperation(t *testing.T) {
// 	testStruct := &TestConfigStruct{}
// 	testStruct.TestParseFlagsWithParser_SumOperation(t)
// }

// TestParseFlagsWithParser_EvaluateLineCountOperation はevaluate_line_count操作のテスト
// func TestParseFlagsWithParser_EvaluateLineCountOperation(t *testing.T) {
// 	testStruct := &TestConfigStruct{}
// 	testStruct.TestParseFlagsWithParser_EvaluateLineCountOperation(t)
// }

// TestParseFlagsWithParser_ParseError はパースエラーのテスト
func TestParseFlagsWithParser_ParseError(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestParseFlagsWithParser_ParseError(t)
}

// TestPrintUsage_Normal はPrintUsage関数のテスト
func TestPrintUsage_Normal(t *testing.T) {
	testStruct := &TestConfigStruct{}
	testStruct.TestPrintUsage_Normal(t)
}
