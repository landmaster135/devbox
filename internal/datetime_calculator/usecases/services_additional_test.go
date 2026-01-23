package usecases

import (
	"testing"
)

// TestDatetimeCalculatorService_GetCalculator はGetCalculator関数のテスト
func TestDatetimeCalculatorService_GetCalculator(t *testing.T) {
	// Arrange
	service := NewDatetimeCalculatorService()

	// Act
	calculator := service.GetCalculator()

	// Assert
	if calculator == nil {
		t.Error("GetCalculator() returned nil")
		return
	}

	// serviceのcalculatorフィールドと同じインスタンスであることを確認
	if calculator != service.calculator {
		t.Error("GetCalculator() returned different instance than service.calculator")
	}
}

// TestDatetimeCalculatorService_GetCalculator_WithFileReader はFileReaderを注入したサービスのGetCalculatorテスト
func TestDatetimeCalculatorService_GetCalculator_WithFileReader(t *testing.T) {
	// Arrange
	mockFileReader := NewMockFileReader()
	service := NewDatetimeCalculatorServiceWithFileReader(mockFileReader)

	// Act
	calculator := service.GetCalculator()

	// Assert
	if calculator == nil {
		t.Error("GetCalculator() returned nil")
		return
	}

	// serviceのcalculatorフィールドと同じインスタンスであることを確認
	if calculator != service.calculator {
		t.Error("GetCalculator() returned different instance than service.calculator")
	}
}

// TestDatetimeCalculatorService_GetCalculator_Consistency は複数回呼び出しの一貫性テスト
func TestDatetimeCalculatorService_GetCalculator_Consistency(t *testing.T) {
	// Arrange
	service := NewDatetimeCalculatorService()

	// Act
	calculator1 := service.GetCalculator()
	calculator2 := service.GetCalculator()

	// Assert
	if calculator1 != calculator2 {
		t.Error("GetCalculator() returned different instances on multiple calls")
	}

	if calculator1 == nil || calculator2 == nil {
		t.Error("GetCalculator() returned nil on one or more calls")
	}
}

// TestDatetimeCalculatorService_GetCalculator_FunctionalTest はGetCalculatorで取得したcalculatorの機能テスト
func TestDatetimeCalculatorService_GetCalculator_FunctionalTest(t *testing.T) {
	// Arrange
	service := NewDatetimeCalculatorService()
	calculator := service.GetCalculator()

	if calculator == nil {
		t.Fatal("GetCalculator() returned nil")
	}

	// Act & Assert - calculatorのメソッドが正常に動作することを確認
	testCases := []struct {
		name     string
		unit     string
		expected string
	}{
		{
			name:     "年単位の日本語名取得",
			unit:     "year",
			expected: "年",
		},
		{
			name:     "時間単位の日本語名取得",
			unit:     "hour",
			expected: "時間",
		},
		{
			name:     "分単位の日本語名取得",
			unit:     "minute",
			expected: "分",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculator.GetUnitName(tc.unit)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}

	// sumTimeFloat機能のテスト
	figures := []float64{1.5, 2.5, 3.0}
	expectedSum := 7.0
	actualSum := calculator.SumTimeFloat(figures)
	if actualSum != expectedSum {
		t.Errorf("sumTimeFloat: expected %f, got %f", expectedSum, actualSum)
	}
}
