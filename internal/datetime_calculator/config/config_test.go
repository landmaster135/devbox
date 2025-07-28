package config

import (
	"fmt"
	"os"
	"testing"
)

// MockFlagParser はテスト用のFlagParserモック実装
type MockFlagParser struct {
	stringFlags map[string]string
	boolFlags   map[string]bool
	parseError  error
	args        []string
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser() *MockFlagParser {
	return &MockFlagParser{
		stringFlags: make(map[string]string),
		boolFlags:   make(map[string]bool),
		args:        []string{},
	}
}

// StringVar は文字列フラグを設定する（モック用）
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
	if val, exists := m.stringFlags[name]; exists {
		*p = val
	} else {
		*p = value
	}
}

// BoolVar はブールフラグを設定する（モック用）
func (m *MockFlagParser) BoolVar(p *bool, name string, value bool, usage string) {
	if val, exists := m.boolFlags[name]; exists {
		*p = val
	} else {
		*p = value
	}
}

// Parse はフラグを解析する（モック用）
func (m *MockFlagParser) Parse() error {
	return m.parseError
}

// Args は解析後の残りの引数を返す（モック用）
func (m *MockFlagParser) Args() []string {
	return m.args
}

// SetStringFlag はモック用の文字列フラグを設定する
func (m *MockFlagParser) SetStringFlag(name, value string) {
	m.stringFlags[name] = value
}

// SetBoolFlag はモック用のブールフラグを設定する
func (m *MockFlagParser) SetBoolFlag(name string, value bool) {
	m.boolFlags[name] = value
}

// SetParseError はモック用のパースエラーを設定する
func (m *MockFlagParser) SetParseError(err error) {
	m.parseError = err
}

// SetArgs はモック用の引数を設定する
func (m *MockFlagParser) SetArgs(args []string) {
	m.args = args
}

// TestNewConfig_Normal はNewConfig関数の正常系テスト
func TestNewConfig_Normal(t *testing.T) {
	// Arrange
	operation := "add"
	year1 := 2025.0
	month1 := 1.0
	day1 := 15.0
	hour1 := 12.0
	minute1 := 30.0
	second1 := 45.0
	durationYear := 1.0
	durationMonth := 2.0
	durationDay := 10.0
	durationHour := 5.0
	durationMinute := 30.0
	durationSecond := 15.0

	// Act
	config, err := NewConfig(operation, year1, month1, day1, hour1, minute1, second1, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond)

	// Assert
	if err != nil {
		t.Errorf("NewConfig returned error: %v", err)
	}
	if config.Operation != operation {
		t.Errorf("Expected operation %s, got %s", operation, config.Operation)
	}
	if config.Year1 != year1 {
		t.Errorf("Expected Year1 %f, got %f", year1, config.Year1)
	}
	if config.Month1 != month1 {
		t.Errorf("Expected Month1 %f, got %f", month1, config.Month1)
	}
	if config.Day1 != day1 {
		t.Errorf("Expected Day1 %f, got %f", day1, config.Day1)
	}
	if config.Hour1 != hour1 {
		t.Errorf("Expected Hour1 %f, got %f", hour1, config.Hour1)
	}
	if config.Minute1 != minute1 {
		t.Errorf("Expected Minute1 %f, got %f", minute1, config.Minute1)
	}
	if config.Second1 != second1 {
		t.Errorf("Expected Second1 %f, got %f", second1, config.Second1)
	}
	if config.DurationYear != durationYear {
		t.Errorf("Expected DurationYear %f, got %f", durationYear, config.DurationYear)
	}
	if config.DurationMonth != durationMonth {
		t.Errorf("Expected DurationMonth %f, got %f", durationMonth, config.DurationMonth)
	}
	if config.DurationDay != durationDay {
		t.Errorf("Expected DurationDay %f, got %f", durationDay, config.DurationDay)
	}
	if config.DurationHour != durationHour {
		t.Errorf("Expected DurationHour %f, got %f", durationHour, config.DurationHour)
	}
	if config.DurationMinute != durationMinute {
		t.Errorf("Expected DurationMinute %f, got %f", durationMinute, config.DurationMinute)
	}
	if config.DurationSecond != durationSecond {
		t.Errorf("Expected DurationSecond %f, got %f", durationSecond, config.DurationSecond)
	}
}

// TestNewConfig_EmptyOperation は操作タイプが空の場合のテスト
func TestNewConfig_EmptyOperation(t *testing.T) {
	// Arrange
	operation := ""

	// Act
	config, err := NewConfig(operation, 2025, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0)

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

	// Act
	config, err := NewConfig(operation, 2025, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid operation, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid operation")
	}
}

// TestNewConfig_AllValidOperations は全ての有効な操作タイプのテスト
func TestNewConfig_AllValidOperations(t *testing.T) {
	validOperations := []string{"add", "subtract"}

	for _, operation := range validOperations {
		// Arrange & Act
		config, err := NewConfig(operation, 2025, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0)

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

// TestParseFlags_Normal はParseFlags関数の正常系テスト
func TestParseFlags_Normal(t *testing.T) {
	// Arrange - ParseFlags関数を直接テストするため、os.Argsを設定
	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	// 有効な引数を設定
	os.Args = []string{"test-program", "-operation=add", "-year=2025", "-month=3", "-day=15", "-hour=12", "-minute=30", "-second=45", "-duration-year=1", "-duration-month=2", "-duration-day=10", "-duration-hour=5", "-duration-minute=30", "-duration-second=15"}

	// Act
	config, err := ParseFlags()

	// Assert
	if err != nil {
		t.Errorf("ParseFlags returned error: %v", err)
	}
	if config.Operation != "add" {
		t.Errorf("Expected operation 'add', got %s", config.Operation)
	}
	if config.Year1 != 2025.0 {
		t.Errorf("Expected Year1 2025.0, got %f", config.Year1)
	}
	if config.Month1 != 3.0 {
		t.Errorf("Expected Month1 3.0, got %f", config.Month1)
	}
	if config.Day1 != 15.0 {
		t.Errorf("Expected Day1 15.0, got %f", config.Day1)
	}
}

// TestParseFlagsWithParser_HelpFlag はヘルプフラグのテスト
func TestParseFlagsWithParser_HelpFlag(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	// ヘルプフラグの短縮形も設定
	mockParser.SetBoolFlag("help", true)
	mockParser.SetBoolFlag("h", true)

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error: %v", err)
	}
	if config == nil {
		t.Fatal("Expected non-nil config for help flag")
	}
	if !config.Help {
		t.Error("Expected Help flag to be true")
	}
}

// TestParseFlagsWithParser_ShortFlags は短縮形フラグのテスト
func TestParseFlagsWithParser_ShortFlags(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("o", "subtract")
	mockParser.SetStringFlag("y", "2024")
	mockParser.SetStringFlag("m", "12")
	mockParser.SetStringFlag("d", "25")
	mockParser.SetStringFlag("hr", "18")
	mockParser.SetStringFlag("min", "45")
	mockParser.SetStringFlag("s", "30")
	mockParser.SetStringFlag("dy", "0")
	mockParser.SetStringFlag("dm", "1")
	mockParser.SetStringFlag("dd", "5")
	mockParser.SetStringFlag("dh", "2")
	mockParser.SetStringFlag("dmin", "15")
	mockParser.SetStringFlag("ds", "45")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err != nil {
		t.Errorf("ParseFlagsWithParser returned error: %v", err)
	}
	if config.Operation != "subtract" {
		t.Errorf("Expected operation 'subtract', got %s", config.Operation)
	}
	if config.Year1 != 2024.0 {
		t.Errorf("Expected Year1 2024.0, got %f", config.Year1)
	}
	if config.Month1 != 12.0 {
		t.Errorf("Expected Month1 12.0, got %f", config.Month1)
	}
}

// TestParseFlagsWithParser_InvalidYear は無効な年の値のテスト
func TestParseFlagsWithParser_InvalidYear(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("year", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid year value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid year value")
	}
}

// TestParseFlagsWithParser_InvalidMonth は無効な月の値のテスト
func TestParseFlagsWithParser_InvalidMonth(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("month", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid month value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid month value")
	}
}

// TestParseFlagsWithParser_InvalidDay は無効な日の値のテスト
func TestParseFlagsWithParser_InvalidDay(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("day", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid day value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid day value")
	}
}

// TestParseFlagsWithParser_InvalidHour は無効な時の値のテスト
func TestParseFlagsWithParser_InvalidHour(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("hour", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid hour value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid hour value")
	}
}

// TestParseFlagsWithParser_InvalidMinute は無効な分の値のテスト
func TestParseFlagsWithParser_InvalidMinute(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("minute", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid minute value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid minute value")
	}
}

// TestParseFlagsWithParser_InvalidSecond は無効な秒の値のテスト
func TestParseFlagsWithParser_InvalidSecond(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("second", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid second value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid second value")
	}
}

// TestParseFlagsWithParser_InvalidDurationYear は無効な期間年の値のテスト
func TestParseFlagsWithParser_InvalidDurationYear(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("duration-year", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid duration year value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid duration year value")
	}
}

// TestParseFlagsWithParser_InvalidDurationMonth は無効な期間月の値のテスト
func TestParseFlagsWithParser_InvalidDurationMonth(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("duration-month", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid duration month value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid duration month value")
	}
}

// TestParseFlagsWithParser_InvalidDurationDay は無効な期間日の値のテスト
func TestParseFlagsWithParser_InvalidDurationDay(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("duration-day", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid duration day value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid duration day value")
	}
}

// TestParseFlagsWithParser_InvalidDurationHour は無効な期間時の値のテスト
func TestParseFlagsWithParser_InvalidDurationHour(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("duration-hour", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid duration hour value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid duration hour value")
	}
}

// TestParseFlagsWithParser_InvalidDurationMinute は無効な期間分の値のテスト
func TestParseFlagsWithParser_InvalidDurationMinute(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("duration-minute", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid duration minute value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid duration minute value")
	}
}

// TestParseFlagsWithParser_InvalidDurationSecond は無効な期間秒の値のテスト
func TestParseFlagsWithParser_InvalidDurationSecond(t *testing.T) {
	// Arrange
	mockParser := NewMockFlagParser()
	mockParser.SetStringFlag("operation", "add")
	mockParser.SetStringFlag("duration-second", "invalid")

	// Act
	config, err := ParseFlagsWithParser(mockParser)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid duration second value, got nil")
	}
	if config != nil {
		t.Error("Expected nil config for invalid duration second value")
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
