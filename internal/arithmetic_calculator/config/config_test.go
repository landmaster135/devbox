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
