package config

import (
	"testing"
)

// TestConfigStruct はConfig構造体のテストクラス
type TestConfigStruct struct{}

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
