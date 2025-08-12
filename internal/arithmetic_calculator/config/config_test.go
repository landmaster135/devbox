package config

import (
	"fmt"
	"os"
	"testing"
)

// TestParseFlags_Normal はParseFlags関数の正常系テスト
func TestParseFlags_Normal(t *testing.T) {
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
		t.Errorf("ParseFlags returned error: %v", err)
	}
	if config.Operation != "subtract" {
		t.Errorf("Expected operation 'subtract', got %s", config.Operation)
	}
	if config.X != 20.0 {
		t.Errorf("Expected X 20.0, got %f", config.X)
	}
	if config.Y != 8.0 {
		t.Errorf("Expected Y 8.0, got %f", config.Y)
	}
}

// TestParseFlagsWithParser_InvalidXValue は無効なX値のテスト
func TestParseFlagsWithParser_InvalidXValue(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("x", "invalid")
	mockParser.SetStringFlag("y", "5")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid X value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid X value")
	}
}

// TestParseFlagsWithParser_InvalidYValue は無効なY値のテスト
func TestParseFlagsWithParser_InvalidYValue(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("x", "10")
	mockParser.SetStringFlag("y", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid Y value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid Y value")
	}
}

// TestParseFlagsWithParser_InvalidThreshold は無効な閾値のテスト
func TestParseFlagsWithParser_InvalidThreshold(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "evaluate_line_count")
	mockParser.SetStringFlag("file", "/test/file.txt")
	mockParser.SetStringFlag("threshold", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid threshold value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid threshold value")
	}
}

// TestParseFlagsWithParser_InvalidNumbers は無効な数値リストのテスト
func TestParseFlagsWithParser_InvalidNumbers(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "sum")
	mockParser.SetStringFlag("numbers", "1,2,invalid,4")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid numbers, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid numbers")
	}
}

// TestParseFlagsWithParser_ParseError はパースエラーのテスト
func TestParseFlagsWithParser_ParseError(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetParseError(fmt.Errorf("parse error"))

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected parse error, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for parse error")
	}
}

// TestPrintUsage_Normal はPrintUsage関数のテスト
func TestPrintUsage_Normal(t *testing.T) {
	// Act - PrintUsageは標準エラー出力に書き込むため、パニックしないことを確認
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("PrintUsage panicked: %v", r)
		}
	}()

	PrintUsage()

	// Assert - パニックしなければ成功
}

// TestNewConfig_Normal はNewConfigの正常系テスト
func TestNewConfig_Normal(t *testing.T) {
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
		t.Errorf("NewConfig returned error: %v", err)
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.X != x {
		t.Errorf("Expected X %f, got %f", x, config.X)
	}
	if config.Y != y {
		t.Errorf("Expected Y %f, got %f", y, config.Y)
	}
	if len(config.Numbers) != len(numbers) {
		t.Errorf("Expected Numbers length %d, got %d", len(numbers), len(config.Numbers))
	}
	if config.FilePath != filePath {
		t.Errorf("Expected FilePath %s, got %s", filePath, config.FilePath)
	}
	if config.Threshold != threshold {
		t.Errorf("Expected Threshold %d, got %d", threshold, config.Threshold)
	}
}

// TestNewConfig_EmptyOperation は操作タイプが空の場合のテスト
func TestNewConfig_EmptyOperation(t *testing.T) {
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
		t.Error("Expected error for empty operation, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for empty operation")
	}
}

// TestNewConfig_InvalidOperation は無効な操作タイプの場合のテスト
func TestNewConfig_InvalidOperation(t *testing.T) {
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
		t.Error("Expected error for invalid operation, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid operation")
	}
}

// TestNewConfig_SumWithEmptyNumbers はsum操作で数値配列が空の場合のテスト
func TestNewConfig_SumWithEmptyNumbers(t *testing.T) {
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
		t.Error("Expected error for sum operation with empty numbers, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for sum operation with empty numbers")
	}
}

// TestNewConfig_EvaluateLineCountWithEmptyFilePath はevaluate_line_count操作でファイルパスが空の場合のテスト
func TestNewConfig_EvaluateLineCountWithEmptyFilePath(t *testing.T) {
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
		t.Error("Expected error for evaluate_line_count operation with empty file path, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for evaluate_line_count operation with empty file path")
	}
}

// TestNewConfig_EvaluateLineCountWithNegativeThreshold はevaluate_line_count操作で負の閾値の場合のテスト
func TestNewConfig_EvaluateLineCountWithNegativeThreshold(t *testing.T) {
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
		t.Error("Expected error for evaluate_line_count operation with negative threshold, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for evaluate_line_count operation with negative threshold")
	}
}

// TestNewConfig_AllValidOperations は全ての有効な操作タイプのテスト
func TestNewConfig_AllValidOperations(t *testing.T) {
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
			t.Errorf("NewConfig returned error for valid operation %s: %v", operation, err)
		}
		if config == nil {
			t.Errorf("Expected non-nil config for valid operation %s", operation)
		}
		if config != nil && config.Operation != operation {
			t.Errorf("Expected operation %s, got %s", operation, config.Operation)
		}
	}
}

// TestNewConfigForParseApiCost_Normal はNewConfigForParseApiCostの正常系テスト
func TestNewConfigForParseApiCost_Normal(t *testing.T) {
	// Arrange
	operation := "parse-api-cost"
	filePath := "/test/file.txt"
	textInput := ""

	// Act
	config, err := NewConfigForParseApiCost(operation, filePath, textInput)

	// Assert
	if err != nil {
		t.Errorf("NewConfigForParseApiCost returned error: %v", err)
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.FilePath != filePath {
		t.Errorf("Expected FilePath %s, got %s", filePath, config.FilePath)
	}
	if config.TextInput != textInput {
		t.Errorf("Expected TextInput %s, got %s", textInput, config.TextInput)
	}
}

// TestNewConfigForParseApiCost_WithTextInput はテキスト入力を使用したテスト
func TestNewConfigForParseApiCost_WithTextInput(t *testing.T) {
	// Arrange
	operation := "parse-api-cost"
	filePath := ""
	textInput := "API料金が100円掛かった"

	// Act
	config, err := NewConfigForParseApiCost(operation, filePath, textInput)

	// Assert
	if err != nil {
		t.Errorf("NewConfigForParseApiCost returned error: %v", err)
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.FilePath != filePath {
		t.Errorf("Expected FilePath %s, got %s", filePath, config.FilePath)
	}
	if config.TextInput != textInput {
		t.Errorf("Expected TextInput %s, got %s", textInput, config.TextInput)
	}
}

// TestNewConfigForParseApiCost_EmptyOperation は操作タイプが空の場合のテスト
func TestNewConfigForParseApiCost_EmptyOperation(t *testing.T) {
	// Arrange
	operation := ""
	filePath := "/test/file.txt"
	textInput := ""

	// Act
	config, err := NewConfigForParseApiCost(operation, filePath, textInput)

	// Assert
	if err == nil {
		t.Error("Expected error for empty operation, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for empty operation")
	}
}

// TestNewConfigForParseApiCost_InvalidOperation は無効な操作タイプの場合のテスト
func TestNewConfigForParseApiCost_InvalidOperation(t *testing.T) {
	// Arrange
	operation := "invalid-operation"
	filePath := "/test/file.txt"
	textInput := ""

	// Act
	config, err := NewConfigForParseApiCost(operation, filePath, textInput)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid operation, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid operation")
	}
}

// TestNewConfigForParseApiCost_BothFilePathAndTextInput はファイルパスとテキスト入力の両方が指定された場合のテスト
func TestNewConfigForParseApiCost_BothFilePathAndTextInput(t *testing.T) {
	// Arrange
	operation := "parse-api-cost"
	filePath := "/test/file.txt"
	textInput := "API料金が100円掛かった"

	// Act
	config, err := NewConfigForParseApiCost(operation, filePath, textInput)

	// Assert
	if err == nil {
		t.Error("Expected error for both file path and text input, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for both file path and text input")
	}
}

// TestNewConfigForParseApiCost_NeitherFilePathNorTextInput はファイルパスもテキスト入力も指定されない場合のテスト
func TestNewConfigForParseApiCost_NeitherFilePathNorTextInput(t *testing.T) {
	// Arrange
	operation := "parse-api-cost"
	filePath := ""
	textInput := ""

	// Act
	config, err := NewConfigForParseApiCost(operation, filePath, textInput)

	// Assert
	if err == nil {
		t.Error("Expected error for neither file path nor text input, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for neither file path nor text input")
	}
}

// TestNewConfigForParseApiCost_InvalidFileExtension は無効なファイル拡張子の場合のテスト
func TestNewConfigForParseApiCost_InvalidFileExtension(t *testing.T) {
	// Arrange
	operation := "parse-api-cost"
	filePath := "/test/file.json"
	textInput := ""

	// Act
	config, err := NewConfigForParseApiCost(operation, filePath, textInput)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid file extension, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid file extension")
	}
}

// TestNewConfigForParseApiCost_ValidMdExtension は.md拡張子の場合のテスト
func TestNewConfigForParseApiCost_ValidMdExtension(t *testing.T) {
	// Arrange
	operation := "parse-api-cost"
	filePath := "/test/file.md"
	textInput := ""

	// Act
	config, err := NewConfigForParseApiCost(operation, filePath, textInput)

	// Assert
	if err != nil {
		t.Errorf("NewConfigForParseApiCost returned error for .md file: %v", err)
	}
	if config == nil {
		t.Error("Expected non-nil config for .md file")
	}
}

// TestParseFlagsWithParser_HelpFlag はヘルプフラグのテスト
func TestParseFlagsWithParser_HelpFlag(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()

	// ParseFlagsWithParserが内部でフラグを定義するので、先にダミーの値を設定
	var help bool
	mockParser.BoolVar(&help, "help", false, "ヘルプ")
	mockParser.SetBoolFlag("help", true)

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for help flag: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for help flag")
	}
	if !config.Help {
		t.Error("Expected Help flag to be true")
	}
}

// TestParseFlagsWithParser_ShortFormFlags は短縮形フラグのテスト
func TestParseFlagsWithParser_ShortFormFlags(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("o", "add")
	mockParser.SetStringFlag("x", "15")
	mockParser.SetStringFlag("y", "25")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for short form flags: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for short form flags")
	}
	if config.Operation != "add" {
		t.Errorf("Expected operation 'add', got %s", config.Operation)
	}
	if config.X != 15.0 {
		t.Errorf("Expected X 15.0, got %f", config.X)
	}
	if config.Y != 25.0 {
		t.Errorf("Expected Y 25.0, got %f", config.Y)
	}
}

// TestParseFlagsWithParser_SumOperationWithShortForm はsum操作の短縮形フラグのテスト
func TestParseFlagsWithParser_SumOperationWithShortForm(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("o", "sum")
	mockParser.SetStringFlag("n", "1,2,3,4,5")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for sum operation with short form: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for sum operation")
	}
	if config.Operation != "sum" {
		t.Errorf("Expected operation 'sum', got %s", config.Operation)
	}
	expectedNumbers := []float64{1, 2, 3, 4, 5}
	if len(config.Numbers) != len(expectedNumbers) {
		t.Errorf("Expected Numbers length %d, got %d", len(expectedNumbers), len(config.Numbers))
	}
	for i, expected := range expectedNumbers {
		if config.Numbers[i] != expected {
			t.Errorf("Expected Numbers[%d] %f, got %f", i, expected, config.Numbers[i])
		}
	}
}

// TestParseFlagsWithParser_EvaluateLineCountWithShortForm はevaluate_line_count操作の短縮形フラグのテスト
func TestParseFlagsWithParser_EvaluateLineCountWithShortForm(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("o", "evaluate_line_count")
	mockParser.SetStringFlag("f", "/test/file.txt")
	mockParser.SetStringFlag("t", "200")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for evaluate_line_count with short form: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for evaluate_line_count operation")
	}
	if config.Operation != "evaluate_line_count" {
		t.Errorf("Expected operation 'evaluate_line_count', got %s", config.Operation)
	}
	if config.FilePath != "/test/file.txt" {
		t.Errorf("Expected FilePath '/test/file.txt', got %s", config.FilePath)
	}
	if config.Threshold != 200 {
		t.Errorf("Expected Threshold 200, got %d", config.Threshold)
	}
}

// TestParseFlagsWithParser_ParseApiCostWithShortForm はparse-api-cost操作の短縮形フラグのテスト
func TestParseFlagsWithParser_ParseApiCostWithShortForm(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("o", "parse-api-cost")
	mockParser.SetStringFlag("ti", "API料金が300円掛かった")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for parse-api-cost with short form: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for parse-api-cost operation")
	}
	if config.Operation != "parse-api-cost" {
		t.Errorf("Expected operation 'parse-api-cost', got %s", config.Operation)
	}
	if config.TextInput != "API料金が300円掛かった" {
		t.Errorf("Expected TextInput 'API料金が300円掛かった', got %s", config.TextInput)
	}
}

// TestParseFlagsWithParser_PositionalArguments は位置引数のテスト
func TestParseFlagsWithParser_PositionalArguments(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "multiply")
	mockParser.SetArgs([]string{"7", "8"})

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for positional arguments: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for positional arguments")
	}
	if config.Operation != "multiply" {
		t.Errorf("Expected operation 'multiply', got %s", config.Operation)
	}
	if config.X != 7.0 {
		t.Errorf("Expected X 7.0, got %f", config.X)
	}
	if config.Y != 8.0 {
		t.Errorf("Expected Y 8.0, got %f", config.Y)
	}
}

// TestParseFlagsWithParser_EvaluateLineCountPositionalArgument はevaluate_line_count操作の位置引数テスト
func TestParseFlagsWithParser_EvaluateLineCountPositionalArgument(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "evaluate_line_count")
	mockParser.SetStringFlag("file", "/test/file.txt")
	mockParser.SetArgs([]string{"150"})

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for evaluate_line_count positional argument: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for evaluate_line_count positional argument")
	}
	if config.Operation != "evaluate_line_count" {
		t.Errorf("Expected operation 'evaluate_line_count', got %s", config.Operation)
	}
	if config.FilePath != "/test/file.txt" {
		t.Errorf("Expected FilePath '/test/file.txt', got %s", config.FilePath)
	}
	if config.Threshold != 150 {
		t.Errorf("Expected Threshold 150, got %d", config.Threshold)
	}
}

// TestParseFlagsWithParser_InvalidPositionalArgument は無効な位置引数のテスト
func TestParseFlagsWithParser_InvalidPositionalArgument(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetArgs([]string{"invalid", "5"})

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for invalid positional argument: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config even with invalid positional argument")
	}
	// 無効な位置引数は無視され、デフォルト値が使用される
	if config.X != 0.0 {
		t.Errorf("Expected X 0.0 (default), got %f", config.X)
	}
	if config.Y != 5.0 {
		t.Errorf("Expected Y 5.0, got %f", config.Y)
	}
}
