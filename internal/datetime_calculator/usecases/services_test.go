package usecases

import (
	"testing"
)

// TestNewDatetimeCalculatorService は NewDatetimeCalculatorService のテストです
func TestNewDatetimeCalculatorService(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "正常なサービス作成",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewDatetimeCalculatorService()

			// サービスがnilでないことを確認
			if service == nil {
				t.Error("NewDatetimeCalculatorService() returned nil")
				return
			}

			// calculatorフィールドがnilでないことを確認
			if service.calculator == nil {
				t.Error("NewDatetimeCalculatorService() calculator field is nil")
			}
		})
	}
}

// TestDatetimeCalculatorService_HandleTimeUnitSum は HandleTimeUnitSum のテストです
func TestDatetimeCalculatorService_HandleTimeUnitSum(t *testing.T) {
	service := NewDatetimeCalculatorService()

	tests := []struct {
		name       string
		figures    []float64
		inputUnit  string
		outputUnit string
		expected   float64
		wantErr    bool
	}{
		{
			name:       "秒から時間への変換_正常ケース",
			figures:    []float64{3600, 1800, 7200},
			inputUnit:  "second",
			outputUnit: "hour",
			expected:   3.5,
			wantErr:    false,
		},
		{
			name:       "時間から分への変換_正常ケース",
			figures:    []float64{2.5},
			inputUnit:  "hour",
			outputUnit: "minute",
			expected:   150,
			wantErr:    false,
		},
		{
			name:       "日から月への変換_正常ケース",
			figures:    []float64{30, 15, 45},
			inputUnit:  "day",
			outputUnit: "month",
			expected:   3,
			wantErr:    false,
		},
		{
			name:       "分から秒への変換_正常ケース",
			figures:    []float64{1, 2, 3},
			inputUnit:  "minute",
			outputUnit: "second",
			expected:   360,
			wantErr:    false,
		},
		{
			name:       "年から日への変換_正常ケース",
			figures:    []float64{1},
			inputUnit:  "year",
			outputUnit: "day",
			expected:   365,
			wantErr:    false,
		},
		{
			name:       "同じ単位での合計_正常ケース",
			figures:    []float64{10, 20, 30},
			inputUnit:  "hour",
			outputUnit: "hour",
			expected:   60,
			wantErr:    false,
		},
		{
			name:       "小数値の処理_正常ケース",
			figures:    []float64{0.5, 1.5, 2.0},
			inputUnit:  "hour",
			outputUnit: "minute",
			expected:   240,
			wantErr:    false,
		},
		{
			name:       "複雑な変換_秒から年",
			figures:    []float64{31536000},
			inputUnit:  "second",
			outputUnit: "year",
			expected:   1,
			wantErr:    false,
		},
		{
			name:       "空の配列_エラーケース",
			figures:    []float64{},
			inputUnit:  "second",
			outputUnit: "hour",
			expected:   0,
			wantErr:    true,
		},
		{
			name:       "無効な入力単位_エラーケース",
			figures:    []float64{3600},
			inputUnit:  "invalid",
			outputUnit: "hour",
			expected:   0,
			wantErr:    true,
		},
		{
			name:       "無効な出力単位_エラーケース",
			figures:    []float64{3600},
			inputUnit:  "second",
			outputUnit: "invalid",
			expected:   0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.HandleTimeUnitSum(tt.figures, tt.inputUnit, tt.outputUnit)

			// エラーの有無をチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleTimeUnitSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーがない場合は結果を検証
			if !tt.wantErr {
				if result != tt.expected {
					t.Errorf("HandleTimeUnitSum() = %v, want %v", result, tt.expected)
				}
			}

			// エラーがある場合は結果が0であることを確認
			if tt.wantErr {
				if result != 0 {
					t.Errorf("HandleTimeUnitSum() error case result = %v, want 0", result)
				}
			}
		})
	}
}

// TestDatetimeCalculatorService_HandleTimeUnitSum_EdgeCases はエッジケースのテストです
func TestDatetimeCalculatorService_HandleTimeUnitSum_EdgeCases(t *testing.T) {
	service := NewDatetimeCalculatorService()

	tests := []struct {
		name       string
		figures    []float64
		inputUnit  string
		outputUnit string
		expected   float64
		wantErr    bool
	}{
		{
			name:       "ゼロ値を含む配列",
			figures:    []float64{0, 3600, 0},
			inputUnit:  "second",
			outputUnit: "hour",
			expected:   1,
			wantErr:    false,
		},
		{
			name:       "負の値を含む配列",
			figures:    []float64{-3600, 7200},
			inputUnit:  "second",
			outputUnit: "hour",
			expected:   1,
			wantErr:    false,
		},
		{
			name:       "すべて負の値",
			figures:    []float64{-1, -2, -3},
			inputUnit:  "hour",
			outputUnit: "minute",
			expected:   -360,
			wantErr:    false,
		},
		{
			name:       "大きな値の処理",
			figures:    []float64{1000000, 2000000},
			inputUnit:  "second",
			outputUnit: "hour",
			expected:   833.3333333333334,
			wantErr:    false,
		},
		{
			name:       "非常に小さな値の処理",
			figures:    []float64{0.001, 0.002, 0.003},
			inputUnit:  "second",
			outputUnit: "minute",
			expected:   0.0001,
			wantErr:    false,
		},
		{
			name:       "単一の値",
			figures:    []float64{7200},
			inputUnit:  "second",
			outputUnit: "hour",
			expected:   2,
			wantErr:    false,
		},
		{
			name:       "多数の値",
			figures:    []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			inputUnit:  "hour",
			outputUnit: "minute",
			expected:   600,
			wantErr:    false,
		},
		{
			name:       "全ての有効な時間単位の組み合わせ_年から秒",
			figures:    []float64{1},
			inputUnit:  "year",
			outputUnit: "second",
			expected:   31536000,
			wantErr:    false,
		},
		{
			name:       "全ての有効な時間単位の組み合わせ_月から分",
			figures:    []float64{1},
			inputUnit:  "month",
			outputUnit: "minute",
			expected:   43200,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.HandleTimeUnitSum(tt.figures, tt.inputUnit, tt.outputUnit)

			// エラーの有無をチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleTimeUnitSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーがない場合は結果を検証
			if !tt.wantErr {
				if result != tt.expected {
					t.Errorf("HandleTimeUnitSum() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

// TestDatetimeCalculatorService_HandleDatetimeCalc は HandleDatetimeCalc のテストです
func TestDatetimeCalculatorService_HandleDatetimeCalc(t *testing.T) {
	service := NewDatetimeCalculatorService()

	tests := []struct {
		name           string
		op             string
		year1          float64
		month1         float64
		day1           float64
		hour1          float64
		minute1        float64
		second1        float64
		durationYear   float64
		durationMonth  float64
		durationDay    float64
		durationHour   float64
		durationMinute float64
		durationSecond float64
		expected       string
		wantErr        bool
	}{
		{
			name:           "add操作_正常ケース",
			op:             "add",
			year1:          2023.0,
			month1:         5.0,
			day1:           15.0,
			hour1:          10.0,
			minute1:        30.0,
			second1:        45.0,
			durationYear:   1.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "2024-05-15 10:30:45",
			wantErr:        false,
		},
		{
			name:           "subtract操作_正常ケース",
			op:             "subtract",
			year1:          2023.0,
			month1:         5.0,
			day1:           15.0,
			hour1:          10.0,
			minute1:        30.0,
			second1:        45.0,
			durationYear:   1.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "2022-05-15 10:30:45",
			wantErr:        false,
		},
		{
			name:           "add操作_複合的な時間追加",
			op:             "add",
			year1:          2023.0,
			month1:         12.0,
			day1:           31.0,
			hour1:          23.0,
			minute1:        59.0,
			second1:        59.0,
			durationYear:   0.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   1.0,
			durationMinute: 1.0,
			durationSecond: 1.0,
			expected:       "2024-01-01 01:01:00",
			wantErr:        false,
		},
		{
			name:           "subtract操作_複合的な時間減算",
			op:             "subtract",
			year1:          2024.0,
			month1:         1.0,
			day1:           1.0,
			hour1:          0.0,
			minute1:        0.0,
			second1:        0.0,
			durationYear:   0.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   1.0,
			durationMinute: 1.0,
			durationSecond: 1.0,
			expected:       "2023-12-31 22:58:59",
			wantErr:        false,
		},
		{
			name:           "add操作_すべての単位を追加",
			op:             "add",
			year1:          2023.0,
			month1:         5.0,
			day1:           15.0,
			hour1:          10.0,
			minute1:        30.0,
			second1:        45.0,
			durationYear:   1.0,
			durationMonth:  2.0,
			durationDay:    3.0,
			durationHour:   4.0,
			durationMinute: 5.0,
			durationSecond: 6.0,
			expected:       "2024-07-18 14:35:51",
			wantErr:        false,
		},
		{
			name:           "subtract操作_すべての単位を減算",
			op:             "subtract",
			year1:          2024.0,
			month1:         7.0,
			day1:           18.0,
			hour1:          14.0,
			minute1:        35.0,
			second1:        51.0,
			durationYear:   1.0,
			durationMonth:  2.0,
			durationDay:    3.0,
			durationHour:   4.0,
			durationMinute: 5.0,
			durationSecond: 6.0,
			expected:       "2023-05-15 10:30:45",
			wantErr:        false,
		},
		{
			name:           "無効な操作_エラーケース",
			op:             "invalid",
			year1:          2023.0,
			month1:         5.0,
			day1:           15.0,
			hour1:          10.0,
			minute1:        30.0,
			second1:        45.0,
			durationYear:   1.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "",
			wantErr:        true,
		},
		{
			name:           "空の操作_エラーケース",
			op:             "",
			year1:          2023.0,
			month1:         5.0,
			day1:           15.0,
			hour1:          10.0,
			minute1:        30.0,
			second1:        45.0,
			durationYear:   1.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "",
			wantErr:        true,
		},
		{
			name:           "multiply操作_エラーケース",
			op:             "multiply",
			year1:          2023.0,
			month1:         5.0,
			day1:           15.0,
			hour1:          10.0,
			minute1:        30.0,
			second1:        45.0,
			durationYear:   2.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "",
			wantErr:        true,
		},
		{
			name:           "divide操作_エラーケース",
			op:             "divide",
			year1:          2023.0,
			month1:         5.0,
			day1:           15.0,
			hour1:          10.0,
			minute1:        30.0,
			second1:        45.0,
			durationYear:   2.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.HandleDatetimeCalc(tt.op, tt.year1, tt.month1, tt.day1, tt.hour1, tt.minute1, tt.second1, tt.durationYear, tt.durationMonth, tt.durationDay, tt.durationHour, tt.durationMinute, tt.durationSecond)

			// エラーの有無をチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleDatetimeCalc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーがない場合は結果を検証
			if !tt.wantErr {
				if result != tt.expected {
					t.Errorf("HandleDatetimeCalc() = %v, want %v", result, tt.expected)
				}
			}

			// エラーがある場合は結果が空文字列であることを確認
			if tt.wantErr {
				if result != tt.expected {
					t.Errorf("HandleDatetimeCalc() error case result = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

// TestDatetimeCalculatorService_HandleDatetimeCalc_EdgeCases はエッジケースのテストです
func TestDatetimeCalculatorService_HandleDatetimeCalc_EdgeCases(t *testing.T) {
	service := NewDatetimeCalculatorService()

	tests := []struct {
		name           string
		op             string
		year1          float64
		month1         float64
		day1           float64
		hour1          float64
		minute1        float64
		second1        float64
		durationYear   float64
		durationMonth  float64
		durationDay    float64
		durationHour   float64
		durationMinute float64
		durationSecond float64
		expected       string
		wantErr        bool
	}{
		{
			name:           "うるう年の処理_add",
			op:             "add",
			year1:          2024.0,
			month1:         2.0,
			day1:           29.0,
			hour1:          12.0,
			minute1:        0.0,
			second1:        0.0,
			durationYear:   1.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "2025-03-01 12:00:00",
			wantErr:        false,
		},
		{
			name:           "月末日の処理_add",
			op:             "add",
			year1:          2023.0,
			month1:         1.0,
			day1:           31.0,
			hour1:          12.0,
			minute1:        0.0,
			second1:        0.0,
			durationYear:   0.0,
			durationMonth:  1.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "2023-03-03 12:00:00",
			wantErr:        false,
		},
		{
			name:           "年末年始の処理_subtract",
			op:             "subtract",
			year1:          2024.0,
			month1:         1.0,
			day1:           1.0,
			hour1:          0.0,
			minute1:        0.0,
			second1:        1.0,
			durationYear:   0.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 1.0,
			expected:       "2024-01-01 00:00:00",
			wantErr:        false,
		},
		{
			name:           "大きな値の処理_add",
			op:             "add",
			year1:          2023.0,
			month1:         6.0,
			day1:           15.0,
			hour1:          12.0,
			minute1:        30.0,
			second1:        45.0,
			durationYear:   100.0,
			durationMonth:  0.0,
			durationDay:    0.0,
			durationHour:   0.0,
			durationMinute: 0.0,
			durationSecond: 0.0,
			expected:       "2123-06-15 12:30:45",
			wantErr:        false,
		},
		{
			name:           "小数点以下の切り捨て確認_add",
			op:             "add",
			year1:          2023.9,
			month1:         5.9,
			day1:           15.9,
			hour1:          10.9,
			minute1:        30.9,
			second1:        45.9,
			durationYear:   1.9,
			durationMonth:  2.9,
			durationDay:    3.9,
			durationHour:   4.9,
			durationMinute: 5.9,
			durationSecond: 6.9,
			expected:       "2024-07-18 14:35:51",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.HandleDatetimeCalc(tt.op, tt.year1, tt.month1, tt.day1, tt.hour1, tt.minute1, tt.second1, tt.durationYear, tt.durationMonth, tt.durationDay, tt.durationHour, tt.durationMinute, tt.durationSecond)

			// エラーの有無をチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleDatetimeCalc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーがない場合は結果を検証
			if !tt.wantErr {
				if result != tt.expected {
					t.Errorf("HandleDatetimeCalc() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}
