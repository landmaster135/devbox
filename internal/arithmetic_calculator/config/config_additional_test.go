package config

import (
	"testing"
)

// TestNewConfig_PowerOperation はpower操作の正常系テスト
func TestNewConfig_PowerOperation(t *testing.T) {
	// Arrange
	operation := "power"
	base := 2.0
	exponent := 3.0

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, base, exponent, 0.0, 0.0, 0, "", "", "", "")

	// Assert
	if err != nil {
		t.Errorf("NewConfig returned error for power operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for power operation")
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.Base != base {
		t.Errorf("Expected Base %f, got %f", base, config.Base)
	}
	if config.Exponent != exponent {
		t.Errorf("Expected Exponent %f, got %f", exponent, config.Exponent)
	}
}

// TestNewConfig_SquareRootOperation はsquare_root操作の正常系テスト
func TestNewConfig_SquareRootOperation(t *testing.T) {
	// Arrange
	operation := "square_root"
	number := 16.0

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, number, 0.0, 0, "", "", "", "")

	// Assert
	if err != nil {
		t.Errorf("NewConfig returned error for square_root operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for square_root operation")
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.Number != number {
		t.Errorf("Expected Number %f, got %f", number, config.Number)
	}
}

// TestNewConfig_FactorialOperation はfactorial操作の正常系テスト
func TestNewConfig_FactorialOperation(t *testing.T) {
	// Arrange
	operation := "factorial"
	n := 5

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, 0.0, 0.0, n, "", "", "", "")

	// Assert
	if err != nil {
		t.Errorf("NewConfig returned error for factorial operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for factorial operation")
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.N != n {
		t.Errorf("Expected N %d, got %d", n, config.N)
	}
}

// TestNewConfig_FactorialWithNegativeNumber はfactorial操作で負数の場合のテスト
func TestNewConfig_FactorialWithNegativeNumber(t *testing.T) {
	// Arrange
	operation := "factorial"
	n := -1

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, 0.0, 0.0, n, "", "", "", "")

	// Assert
	if err == nil {
		t.Error("Expected error for factorial operation with negative number, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for factorial operation with negative number")
	}
}

// TestNewConfig_TrigonometryOperation はtrigonometry操作の正常系テスト
func TestNewConfig_TrigonometryOperation(t *testing.T) {
	// Arrange
	operation := "trigonometry"
	function := "sin"
	angle := 90.0
	unit := "degrees"

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, 0.0, angle, 0, function, unit, "", "")

	// Assert
	if err != nil {
		t.Errorf("NewConfig returned error for trigonometry operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for trigonometry operation")
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.Function != function {
		t.Errorf("Expected Function %s, got %s", function, config.Function)
	}
	if config.Angle != angle {
		t.Errorf("Expected Angle %f, got %f", angle, config.Angle)
	}
	if config.Unit != unit {
		t.Errorf("Expected Unit %s, got %s", unit, config.Unit)
	}
}

// TestNewConfig_TrigonometryWithEmptyFunction はtrigonometry操作で関数名が空の場合のテスト
func TestNewConfig_TrigonometryWithEmptyFunction(t *testing.T) {
	// Arrange
	operation := "trigonometry"
	function := ""
	angle := 90.0
	unit := "degrees"

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, 0.0, angle, 0, function, unit, "", "")

	// Assert
	if err == nil {
		t.Error("Expected error for trigonometry operation with empty function, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for trigonometry operation with empty function")
	}
}

// TestNewConfig_TrigonometryWithDefaultUnit はtrigonometry操作でデフォルト単位のテスト
func TestNewConfig_TrigonometryWithDefaultUnit(t *testing.T) {
	// Arrange
	operation := "trigonometry"
	function := "cos"
	angle := 1.57
	unit := ""

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, 0.0, angle, 0, function, unit, "", "")

	// Assert
	if err != nil {
		t.Errorf("NewConfig returned error for trigonometry operation with default unit: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for trigonometry operation with default unit")
	}
	if config.Unit != "radians" {
		t.Errorf("Expected Unit 'radians' (default), got %s", config.Unit)
	}
}

// TestNewConfig_CalculateOperation はcalculate操作の正常系テスト
func TestNewConfig_CalculateOperation(t *testing.T) {
	// Arrange
	operation := "calculate"
	expression := "2 + 3 * 4"

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, 0.0, 0.0, 0, "", "", expression, "")

	// Assert
	if err != nil {
		t.Errorf("NewConfig returned error for calculate operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for calculate operation")
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.Expression != expression {
		t.Errorf("Expected Expression %s, got %s", expression, config.Expression)
	}
}

// TestNewConfig_CalculateWithEmptyExpression はcalculate操作で数式が空の場合のテスト
func TestNewConfig_CalculateWithEmptyExpression(t *testing.T) {
	// Arrange
	operation := "calculate"
	expression := ""

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, 0.0, 0.0, 0, "", "", expression, "")

	// Assert
	if err == nil {
		t.Error("Expected error for calculate operation with empty expression, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for calculate operation with empty expression")
	}
}

// TestNewConfig_GetConstantsOperation はget_constants操作の正常系テスト
func TestNewConfig_GetConstantsOperation(t *testing.T) {
	// Arrange
	operation := "get_constants"

	// Act
	config, err := NewConfig(operation, 0.0, 0.0, []float64{}, "", 0, 0.0, 0.0, 0.0, 0.0, 0, "", "", "", "")

	// Assert
	if err != nil {
		t.Errorf("NewConfig returned error for get_constants operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for get_constants operation")
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
}

// TestParseFlagsWithParser_PowerOperation はpower操作のパースのテスト
func TestParseFlagsWithParser_PowerOperation(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "power")
	mockParser.SetStringFlag("base", "2")
	mockParser.SetStringFlag("exponent", "8")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for power operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for power operation")
	}
	if config.Operation != "power" {
		t.Errorf("Expected operation 'power', got %s", config.Operation)
	}
	if config.Base != 2.0 {
		t.Errorf("Expected Base 2.0, got %f", config.Base)
	}
	if config.Exponent != 8.0 {
		t.Errorf("Expected Exponent 8.0, got %f", config.Exponent)
	}
}

// TestParseFlagsWithParser_SquareRootOperation はsquare_root操作のパースのテスト
func TestParseFlagsWithParser_SquareRootOperation(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "square_root")
	mockParser.SetStringFlag("number", "25")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for square_root operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for square_root operation")
	}
	if config.Operation != "square_root" {
		t.Errorf("Expected operation 'square_root', got %s", config.Operation)
	}
	if config.Number != 25.0 {
		t.Errorf("Expected Number 25.0, got %f", config.Number)
	}
}

// TestParseFlagsWithParser_FactorialOperation はfactorial操作のパースのテスト
func TestParseFlagsWithParser_FactorialOperation(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "factorial")
	mockParser.SetStringFlag("n", "6")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for factorial operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for factorial operation")
	}
	if config.Operation != "factorial" {
		t.Errorf("Expected operation 'factorial', got %s", config.Operation)
	}
	if config.N != 6 {
		t.Errorf("Expected N 6, got %d", config.N)
	}
}

// TestParseFlagsWithParser_TrigonometryOperation はtrigonometry操作のパースのテスト
func TestParseFlagsWithParser_TrigonometryOperation(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "trigonometry")
	mockParser.SetStringFlag("function", "tan")
	mockParser.SetStringFlag("angle", "45")
	mockParser.SetStringFlag("unit", "degrees")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for trigonometry operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for trigonometry operation")
	}
	if config.Operation != "trigonometry" {
		t.Errorf("Expected operation 'trigonometry', got %s", config.Operation)
	}
	if config.Function != "tan" {
		t.Errorf("Expected Function 'tan', got %s", config.Function)
	}
	if config.Angle != 45.0 {
		t.Errorf("Expected Angle 45.0, got %f", config.Angle)
	}
	if config.Unit != "degrees" {
		t.Errorf("Expected Unit 'degrees', got %s", config.Unit)
	}
}

// TestParseFlagsWithParser_CalculateOperation はcalculate操作のパースのテスト
func TestParseFlagsWithParser_CalculateOperation(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "calculate")
	mockParser.SetStringFlag("expression", "sqrt(16) + 2^3")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for calculate operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for calculate operation")
	}
	if config.Operation != "calculate" {
		t.Errorf("Expected operation 'calculate', got %s", config.Operation)
	}
	if config.Expression != "sqrt(16) + 2^3" {
		t.Errorf("Expected Expression 'sqrt(16) + 2^3', got %s", config.Expression)
	}
}

// TestParseFlagsWithParser_GetConstantsOperation はget_constants操作のパースのテスト
func TestParseFlagsWithParser_GetConstantsOperation(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "get_constants")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for get_constants operation: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for get_constants operation")
	}
	if config.Operation != "get_constants" {
		t.Errorf("Expected operation 'get_constants', got %s", config.Operation)
	}
}

// TestParseFlagsWithParser_InvalidBaseValue は無効なbase値のテスト
func TestParseFlagsWithParser_InvalidBaseValue(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "power")
	mockParser.SetStringFlag("base", "invalid")
	mockParser.SetStringFlag("exponent", "2")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid base value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid base value")
	}
}

// TestParseFlagsWithParser_InvalidExponentValue は無効なexponent値のテスト
func TestParseFlagsWithParser_InvalidExponentValue(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "power")
	mockParser.SetStringFlag("base", "3")
	mockParser.SetStringFlag("exponent", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid exponent value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid exponent value")
	}
}

// TestParseFlagsWithParser_InvalidNumberValue は無効なnumber値のテスト
func TestParseFlagsWithParser_InvalidNumberValue(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "square_root")
	mockParser.SetStringFlag("number", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid number value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid number value")
	}
}

// TestParseFlagsWithParser_InvalidNValue は無効なn値のテスト
func TestParseFlagsWithParser_InvalidNValue(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "factorial")
	mockParser.SetStringFlag("n", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid n value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid n value")
	}
}

// TestParseFlagsWithParser_InvalidAngleValue は無効なangle値のテスト
func TestParseFlagsWithParser_InvalidAngleValue(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "trigonometry")
	mockParser.SetStringFlag("function", "sin")
	mockParser.SetStringFlag("angle", "invalid")
	mockParser.SetStringFlag("unit", "radians")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid angle value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid angle value")
	}
}

// TestParseFlagsWithParser_ShortFormNewParameters は新しいパラメータの短縮形フラグのテスト
func TestParseFlagsWithParser_ShortFormNewParameters(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("o", "power")
	mockParser.SetStringFlag("b", "4")
	mockParser.SetStringFlag("exp", "3")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for short form new parameters: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for short form new parameters")
	}
	if config.Operation != "power" {
		t.Errorf("Expected operation 'power', got %s", config.Operation)
	}
	if config.Base != 4.0 {
		t.Errorf("Expected Base 4.0, got %f", config.Base)
	}
	if config.Exponent != 3.0 {
		t.Errorf("Expected Exponent 3.0, got %f", config.Exponent)
	}
}

// TestParseFlagsWithParser_TrigonometryShortForm はtrigonometry操作の短縮形フラグのテスト
func TestParseFlagsWithParser_TrigonometryShortForm(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("o", "trigonometry")
	mockParser.SetStringFlag("func", "cos")
	mockParser.SetStringFlag("a", "60")
	mockParser.SetStringFlag("u", "degrees")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for trigonometry short form: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for trigonometry short form")
	}
	if config.Operation != "trigonometry" {
		t.Errorf("Expected operation 'trigonometry', got %s", config.Operation)
	}
	if config.Function != "cos" {
		t.Errorf("Expected Function 'cos', got %s", config.Function)
	}
	if config.Angle != 60.0 {
		t.Errorf("Expected Angle 60.0, got %f", config.Angle)
	}
	if config.Unit != "degrees" {
		t.Errorf("Expected Unit 'degrees', got %s", config.Unit)
	}
}

// TestParseFlagsWithParser_CalculateShortForm はcalculate操作の短縮形フラグのテスト
func TestParseFlagsWithParser_CalculateShortForm(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("o", "calculate")
	mockParser.SetStringFlag("expr", "log(100) + sin(pi/2)")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error for calculate short form: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for calculate short form")
	}
	if config.Operation != "calculate" {
		t.Errorf("Expected operation 'calculate', got %s", config.Operation)
	}
	if config.Expression != "log(100) + sin(pi/2)" {
		t.Errorf("Expected Expression 'log(100) + sin(pi/2)', got %s", config.Expression)
	}
}
