package usecases

import (
	"fmt"
	"math"
	"testing"
)

// floatEquals は浮動小数点数を許容誤差で比較する
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}



// TestDatetimeCalculator_sumTimeFloat_Normal はsumTimeFloat関数の正常系テスト
func TestDatetimeCalculator_sumTimeFloat_Normal(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name     string
		figures  []float64
		expected float64
	}{
		{
			name:     "正の数値の合計",
			figures:  []float64{1.5, 2.5, 3.0},
			expected: 7.0,
		},
		{
			name:     "整数値の合計",
			figures:  []float64{10, 20, 30},
			expected: 60,
		},
		{
			name:     "小数値の合計",
			figures:  []float64{0.1, 0.2, 0.3},
			expected: 0.6000000000000001, // 浮動小数点の精度問題
		},
		{
			name:     "単一の値",
			figures:  []float64{42.5},
			expected: 42.5,
		},
		{
			name:     "ゼロを含む合計",
			figures:  []float64{0, 5, 0, 10},
			expected: 15,
		},
		{
			name:     "負の値を含む合計",
			figures:  []float64{-5, 10, -2},
			expected: 3,
		},
		{
			name:     "すべて負の値",
			figures:  []float64{-1, -2, -3},
			expected: -6,
		},
		{
			name:     "大きな値の合計",
			figures:  []float64{1000000, 2000000, 3000000},
			expected: 6000000,
		},
		{
			name:     "非常に小さな値の合計",
			figures:  []float64{0.001, 0.002, 0.003},
			expected: 0.006,
		},
		{
			name:     "多数の値",
			figures:  []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			expected: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := calculator.SumTimeFloat(tc.figures)

			// Assert
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

// TestDatetimeCalculator_GetUnitName_Normal はGetUnitName関数の正常系テスト
func TestDatetimeCalculator_GetUnitName_Normal(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name     string
		unit     string
		expected string
	}{
		{
			name:     "年単位",
			unit:     "year",
			expected: "年",
		},
		{
			name:     "月単位",
			unit:     "month",
			expected: "月",
		},
		{
			name:     "日単位",
			unit:     "day",
			expected: "日",
		},
		{
			name:     "時間単位",
			unit:     "hour",
			expected: "時間",
		},
		{
			name:     "分単位",
			unit:     "minute",
			expected: "分",
		},
		{
			name:     "秒単位",
			unit:     "second",
			expected: "秒",
		},
		{
			name:     "無効な単位_そのまま返す",
			unit:     "invalid_unit",
			expected: "invalid_unit",
		},
		{
			name:     "空文字列",
			unit:     "",
			expected: "",
		},
		{
			name:     "大文字混在_そのまま返す",
			unit:     "Year",
			expected: "Year",
		},
		{
			name:     "複数形_そのまま返す",
			unit:     "years",
			expected: "years",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := calculator.GetUnitName(tc.unit)

			// Assert
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestDatetimeCalculator_GetUnitName_AllValidUnits は全ての有効な単位のテスト
func TestDatetimeCalculator_GetUnitName_AllValidUnits(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	validUnits := map[string]string{
		"year":   "年",
		"month":  "月",
		"day":    "日",
		"hour":   "時間",
		"minute": "分",
		"second": "秒",
	}

	for unit, expectedJapanese := range validUnits {
		t.Run(fmt.Sprintf("単位_%s", unit), func(t *testing.T) {
			// Act
			result := calculator.GetUnitName(unit)

			// Assert
			if result != expectedJapanese {
				t.Errorf("Unit %s: expected %s, got %s", unit, expectedJapanese, result)
			}
		})
	}
}

// TestDatetimeCalculator_GetUnitName_EdgeCases はエッジケースのテスト
func TestDatetimeCalculator_GetUnitName_EdgeCases(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name     string
		unit     string
		expected string
	}{
		{
			name:     "スペースを含む文字列",
			unit:     "hour minute",
			expected: "hour minute",
		},
		{
			name:     "数字を含む文字列",
			unit:     "hour1",
			expected: "hour1",
		},
		{
			name:     "特殊文字を含む文字列",
			unit:     "hour@#$",
			expected: "hour@#$",
		},
		{
			name:     "非常に長い文字列",
			unit:     "this_is_a_very_long_unit_name_that_should_be_returned_as_is",
			expected: "this_is_a_very_long_unit_name_that_should_be_returned_as_is",
		},
		{
			name:     "日本語文字列",
			unit:     "時間単位",
			expected: "時間単位",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := calculator.GetUnitName(tc.unit)

			// Assert
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestDatetimeCalculator_sumTimeFloat_EmptySlice は空のスライスのテスト
func TestDatetimeCalculator_sumTimeFloat_EmptySlice(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}
	figures := []float64{}

	// Act
	result := calculator.SumTimeFloat(figures)

	// Assert
	if result != 0 {
		t.Errorf("Expected 0 for empty slice, got %f", result)
	}
}

// TestDatetimeCalculator_extractTimeFromText_Normal はextractTimeFromText関数の正常系テスト
func TestDatetimeCalculator_extractTimeFromText_Normal(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name       string
		text       string
		outputUnit string
		expected   float64
		wantErr    bool
	}{
		{
			name:       "単一のマッチ_分単位",
			text:       "作業は合計30分掛かった。",
			outputUnit: "minute",
			expected:   30,
			wantErr:    false,
		},
		{
			name:       "複数のマッチ_分単位",
			text:       "作業は合計30分掛かった。別の作業は合計45分掛かった。",
			outputUnit: "minute",
			expected:   75,
			wantErr:    false,
		},
		{
			name:       "分から時間への変換",
			text:       "合計120分掛かった。",
			outputUnit: "hour",
			expected:   2,
			wantErr:    false,
		},
		{
			name:       "分から秒への変換",
			text:       "合計5分掛かった。",
			outputUnit: "second",
			expected:   300,
			wantErr:    false,
		},
		{
			name:       "分から日への変換",
			text:       "合計1440分掛かった。",
			outputUnit: "day",
			expected:   1,
			wantErr:    false,
		},
		{
			name:       "分から月への変換",
			text:       "合計43200分掛かった。",
			outputUnit: "month",
			expected:   1.0,
			wantErr:    false,
		},
		{
			name:       "分から年への変換",
			text:       "合計525600分掛かった。",
			outputUnit: "year",
			expected:   1,
			wantErr:    false,
		},
		{
			name:       "複数のマッチ_時間変換",
			text:       "会議は合計90分掛かった。資料作成は合計60分掛かった。",
			outputUnit: "hour",
			expected:   2.5,
			wantErr:    false,
		},
		{
			name:       "大きな数値",
			text:       "プロジェクトは合計1000分掛かった。",
			outputUnit: "hour",
			expected:   16.666666666666668,
			wantErr:    false,
		},
		{
			name:       "小さな数値",
			text:       "確認は合計1分掛かった。",
			outputUnit: "second",
			expected:   60,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := calculator.ExtractTimeFromText(tc.text, tc.outputUnit)

			// Assert
			if (err != nil) != tc.wantErr {
				t.Errorf("extractTimeFromText() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if !floatEquals(result, tc.expected, 1e-10) {
					t.Errorf("Expected %f, got %f", tc.expected, result)
				}
			}
		})
	}
}

// TestDatetimeCalculator_extractTimeFromText_NoMatches はマッチしない場合のテスト
func TestDatetimeCalculator_extractTimeFromText_NoMatches(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name       string
		text       string
		outputUnit string
		expected   float64
	}{
		{
			name:       "マッチしないテキスト",
			text:       "今日は良い天気です。",
			outputUnit: "minute",
			expected:   0,
		},
		{
			name:       "パターンが異なるテキスト",
			text:       "作業に30分かかりました。",
			outputUnit: "minute",
			expected:   0,
		},
		{
			name:       "数値がないテキスト",
			text:       "作業は合計分掛かった。",
			outputUnit: "minute",
			expected:   0,
		},
		{
			name:       "空のテキスト",
			text:       "",
			outputUnit: "minute",
			expected:   0,
		},
		{
			name:       "部分的にマッチするテキスト",
			text:       "合計30分で終了した。",
			outputUnit: "minute",
			expected:   0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := calculator.ExtractTimeFromText(tc.text, tc.outputUnit)

			// Assert
			if err != nil {
				t.Errorf("extractTimeFromText() returned unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

// TestDatetimeCalculator_extractTimeFromText_InvalidUnit は無効な単位のテスト
func TestDatetimeCalculator_extractTimeFromText_InvalidUnit(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}
	text := "作業は合計30分掛かった。"
	outputUnit := "invalid_unit"

	// Act
	result, err := calculator.ExtractTimeFromText(text, outputUnit)

	// Assert
	if err == nil {
		t.Error("Expected error for invalid output unit, got nil")
	}
	if result != 0 {
		t.Errorf("Expected result 0 for error case, got %f", result)
	}
}

// TestDatetimeCalculator_extractTimeFromText_ComplexText は複雑なテキストのテスト
func TestDatetimeCalculator_extractTimeFromText_ComplexText(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name       string
		text       string
		outputUnit string
		expected   float64
	}{
		{
			name: "複数の文章に分散",
			text: `今日の作業報告です。
			午前中の会議は合計90分掛かった。
			午後の資料作成は合計120分掛かった。
			最後のレビューは合計30分掛かった。`,
			outputUnit: "hour",
			expected:   4,
		},
		{
			name: "数値が混在するテキスト",
			text: `プロジェクトは2024年に開始され、合計180分掛かった。
			参加者は15名で、合計60分掛かった。`,
			outputUnit: "minute",
			expected:   240,
		},
		{
			name:       "改行を含むテキスト",
			text:       "第1フェーズは合計45分掛かった。\n第2フェーズは合計75分掛かった。\n第3フェーズは合計30分掛かった。",
			outputUnit: "hour",
			expected:   2.5,
		},
		{
			name:       "特殊文字を含むテキスト",
			text:       "タスク①は合計20分掛かった。タスク②は合計40分掛かった。タスク③は合計60分掛かった。",
			outputUnit: "minute",
			expected:   120,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := calculator.ExtractTimeFromText(tc.text, tc.outputUnit)

			// Assert
			if err != nil {
				t.Errorf("extractTimeFromText() returned unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

// TestDatetimeCalculator_extractTimeFromText_EdgeCases はエッジケースのテスト
func TestDatetimeCalculator_extractTimeFromText_EdgeCases(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name       string
		text       string
		outputUnit string
		expected   float64
	}{
		{
			name:       "ゼロ分",
			text:       "作業は合計0分掛かった。",
			outputUnit: "minute",
			expected:   0,
		},
		{
			name:       "大きな数値",
			text:       "長期プロジェクトは合計99999分掛かった。",
			outputUnit: "hour",
			expected:   1666.65,
		},
		{
			name:       "同じパターンの重複",
			text:       "作業は合計30分掛かった。作業は合計30分掛かった。作業は合計30分掛かった。",
			outputUnit: "minute",
			expected:   90,
		},
		{
			name:       "数値のみ異なる複数パターン",
			text:       "Aタスクは合計1分掛かった。Bタスクは合計999分掛かった。",
			outputUnit: "minute",
			expected:   1000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := calculator.ExtractTimeFromText(tc.text, tc.outputUnit)

			// Assert
			if err != nil {
				t.Errorf("extractTimeFromText() returned unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}
