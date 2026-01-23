package usecases

import (
	"math"
	"testing"
	"time"
)

func TestNewTime(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		hour     int
		minute   int
		second   int
		expected time.Time
	}{
		{
			name:     "基本的な日付",
			year:     2023,
			month:    5,
			day:      15,
			hour:     10,
			minute:   30,
			second:   45,
			expected: time.Date(2023, time.May, 15, 10, 30, 45, 0, time.Local),
		},
		{
			name:     "うるう年",
			year:     2024,
			month:    2,
			day:      29,
			hour:     0,
			minute:   0,
			second:   0,
			expected: time.Date(2024, time.February, 29, 0, 0, 0, 0, time.Local),
		},
		{
			name:     "年末",
			year:     2023,
			month:    12,
			day:      31,
			hour:     23,
			minute:   59,
			second:   59,
			expected: time.Date(2023, time.December, 31, 23, 59, 59, 0, time.Local),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.NewTime(tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second)
			if !result.Equal(tt.expected) {
				t.Errorf("NewTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStringToTime(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name     string
		dateStr  string
		expected time.Time
		wantErr  bool
	}{
		{
			name:     "有効な日付文字列",
			dateStr:  "2023-05-15 10:30:45",
			expected: time.Date(2023, time.May, 15, 10, 30, 45, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "無効な日付文字列",
			dateStr:  "2023/05/15 10:30:45",
			expected: time.Time{},
			wantErr:  true,
		},
		{
			name:     "空の文字列",
			dateStr:  "",
			expected: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.stringToTime(tt.dateStr)

			// エラーの有無をチェック
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// エラーがない場合は結果を検証
			if !tt.wantErr {
				// Parse関数はUTCで返すため、比較のためにUTCで検証
				if !result.UTC().Equal(tt.expected) {
					t.Errorf("StringToTime() = %v, want %v", result.UTC(), tt.expected)
				}
			}
		})
	}
}

func TestAddDuration(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		hour     int
		minute   int
		second   int
		duration time.Duration
		expected string
	}{
		{
			name:     "1日追加",
			year:     2023,
			month:    5,
			day:      15,
			hour:     10,
			minute:   30,
			second:   45,
			duration: 24 * time.Hour,
			expected: "2023-05-16 10:30:45",
		},
		{
			name:     "1時間追加",
			year:     2023,
			month:    5,
			day:      15,
			hour:     23,
			minute:   30,
			second:   45,
			duration: time.Hour,
			expected: "2023-05-16 00:30:45",
		},
		{
			name:     "複合的な時間追加",
			year:     2023,
			month:    12,
			day:      31,
			hour:     23,
			minute:   59,
			second:   59,
			duration: time.Hour + time.Minute + time.Second,
			expected: "2024-01-01 01:01:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.addDuration(tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second, tt.duration)
			if result != tt.expected {
				t.Errorf("AddDuration() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSubtractDuration(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name     string
		year     int
		month    int
		day      int
		hour     int
		minute   int
		second   int
		duration time.Duration
		expected string
	}{
		{
			name:     "1日減算",
			year:     2023,
			month:    5,
			day:      15,
			hour:     10,
			minute:   30,
			second:   45,
			duration: 24 * time.Hour,
			expected: "2023-05-14 10:30:45",
		},
		{
			name:     "1時間減算",
			year:     2023,
			month:    5,
			day:      15,
			hour:     0,
			minute:   30,
			second:   45,
			duration: time.Hour,
			expected: "2023-05-14 23:30:45",
		},
		{
			name:     "複合的な時間減算",
			year:     2024,
			month:    1,
			day:      1,
			hour:     0,
			minute:   0,
			second:   0,
			duration: time.Hour + time.Minute + time.Second,
			expected: "2023-12-31 22:58:59",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.subtractDuration(tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second, tt.duration)
			if result != tt.expected {
				t.Errorf("SubtractDuration() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAddDatetime(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name       string
		year       int
		month      int
		day        int
		hour       int
		minute     int
		second     int
		addYears   int
		addMonths  int
		addDays    int
		addHours   int
		addMinutes int
		addSeconds int
		expected   string
	}{
		{
			name:       "年を追加",
			year:       2023,
			month:      5,
			day:        15,
			hour:       10,
			minute:     30,
			second:     45,
			addYears:   1,
			addMonths:  0,
			addDays:    0,
			addHours:   0,
			addMinutes: 0,
			addSeconds: 0,
			expected:   "2024-05-15 10:30:45",
		},
		{
			name:       "月を追加（年が変わるケース）",
			year:       2023,
			month:      12,
			day:        15,
			hour:       10,
			minute:     30,
			second:     45,
			addYears:   0,
			addMonths:  1,
			addDays:    0,
			addHours:   0,
			addMinutes: 0,
			addSeconds: 0,
			expected:   "2024-01-15 10:30:45",
		},
		{
			name:       "複合的な時間追加",
			year:       2023,
			month:      12,
			day:        31,
			hour:       23,
			minute:     59,
			second:     59,
			addYears:   0,
			addMonths:  0,
			addDays:    0,
			addHours:   1,
			addMinutes: 1,
			addSeconds: 1,
			expected:   "2024-01-01 01:01:00",
		},
		{
			name:       "すべての単位を追加",
			year:       2023,
			month:      5,
			day:        15,
			hour:       10,
			minute:     30,
			second:     45,
			addYears:   1,
			addMonths:  2,
			addDays:    3,
			addHours:   4,
			addMinutes: 5,
			addSeconds: 6,
			expected:   "2024-07-18 14:35:51",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.addDatetime(tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second, tt.addYears, tt.addMonths, tt.addDays, tt.addHours, tt.addMinutes, tt.addSeconds)
			if result != tt.expected {
				t.Errorf("AddDatetime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAddDatetimeFloat(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name       string
		year       float64
		month      float64
		day        float64
		hour       float64
		minute     float64
		second     float64
		addYears   float64
		addMonths  float64
		addDays    float64
		addHours   float64
		addMinutes float64
		addSeconds float64
		expected   string
	}{
		{
			name:       "整数値と同等のfloat64値",
			year:       2023.0,
			month:      5.0,
			day:        15.0,
			hour:       10.0,
			minute:     30.0,
			second:     45.0,
			addYears:   1.0,
			addMonths:  0.0,
			addDays:    0.0,
			addHours:   0.0,
			addMinutes: 0.0,
			addSeconds: 0.0,
			expected:   "2024-05-15 10:30:45",
		},
		{
			name:       "小数点以下の値（切り捨て確認）",
			year:       2023.9,
			month:      12.9,
			day:        15.9,
			hour:       10.9,
			minute:     30.9,
			second:     45.9,
			addYears:   0.9,
			addMonths:  1.9,
			addDays:    0.9,
			addHours:   0.9,
			addMinutes: 0.9,
			addSeconds: 0.9,
			expected:   "2024-01-15 10:30:45",
		},
		{
			name:       "複合的な時間追加",
			year:       2023.0,
			month:      12.0,
			day:        31.0,
			hour:       23.0,
			minute:     59.0,
			second:     59.0,
			addYears:   0.0,
			addMonths:  0.0,
			addDays:    0.0,
			addHours:   1.0,
			addMinutes: 1.0,
			addSeconds: 1.0,
			expected:   "2024-01-01 01:01:00",
		},
		{
			name:       "すべての単位を追加",
			year:       2023.0,
			month:      5.0,
			day:        15.0,
			hour:       10.0,
			minute:     30.0,
			second:     45.0,
			addYears:   1.0,
			addMonths:  2.0,
			addDays:    3.0,
			addHours:   4.0,
			addMinutes: 5.0,
			addSeconds: 6.0,
			expected:   "2024-07-18 14:35:51",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.AddDatetimeFloat(tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second, tt.addYears, tt.addMonths, tt.addDays, tt.addHours, tt.addMinutes, tt.addSeconds)
			if result != tt.expected {
				t.Errorf("AddDatetimeFloat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSubtractDatetime(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name       string
		year       int
		month      int
		day        int
		hour       int
		minute     int
		second     int
		subYears   int
		subMonths  int
		subDays    int
		subHours   int
		subMinutes int
		subSeconds int
		expected   string
	}{
		{
			name:       "年を減算",
			year:       2023,
			month:      5,
			day:        15,
			hour:       10,
			minute:     30,
			second:     45,
			subYears:   1,
			subMonths:  0,
			subDays:    0,
			subHours:   0,
			subMinutes: 0,
			subSeconds: 0,
			expected:   "2022-05-15 10:30:45",
		},
		{
			name:       "月を減算（年が変わるケース）",
			year:       2023,
			month:      1,
			day:        15,
			hour:       10,
			minute:     30,
			second:     45,
			subYears:   0,
			subMonths:  1,
			subDays:    0,
			subHours:   0,
			subMinutes: 0,
			subSeconds: 0,
			expected:   "2022-12-15 10:30:45",
		},
		{
			name:       "複合的な時間減算",
			year:       2024,
			month:      1,
			day:        1,
			hour:       0,
			minute:     0,
			second:     0,
			subYears:   0,
			subMonths:  0,
			subDays:    0,
			subHours:   1,
			subMinutes: 1,
			subSeconds: 1,
			expected:   "2023-12-31 22:58:59",
		},
		{
			name:       "すべての単位を減算",
			year:       2024,
			month:      7,
			day:        18,
			hour:       14,
			minute:     35,
			second:     51,
			subYears:   1,
			subMonths:  2,
			subDays:    3,
			subHours:   4,
			subMinutes: 5,
			subSeconds: 6,
			expected:   "2023-05-15 10:30:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.subtractDatetime(tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second, tt.subYears, tt.subMonths, tt.subDays, tt.subHours, tt.subMinutes, tt.subSeconds)
			if result != tt.expected {
				t.Errorf("SubtractDatetime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSubtractDatetimeFloat(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name       string
		year       float64
		month      float64
		day        float64
		hour       float64
		minute     float64
		second     float64
		subYears   float64
		subMonths  float64
		subDays    float64
		subHours   float64
		subMinutes float64
		subSeconds float64
		expected   string
	}{
		{
			name:       "整数値と同等のfloat64値",
			year:       2023.0,
			month:      5.0,
			day:        15.0,
			hour:       10.0,
			minute:     30.0,
			second:     45.0,
			subYears:   1.0,
			subMonths:  0.0,
			subDays:    0.0,
			subHours:   0.0,
			subMinutes: 0.0,
			subSeconds: 0.0,
			expected:   "2022-05-15 10:30:45",
		},
		{
			name:       "小数点以下の値（切り捨て確認）",
			year:       2023.9,
			month:      1.9,
			day:        15.9,
			hour:       10.9,
			minute:     30.9,
			second:     45.9,
			subYears:   0.9,
			subMonths:  1.9,
			subDays:    0.9,
			subHours:   0.9,
			subMinutes: 0.9,
			subSeconds: 0.9,
			expected:   "2022-12-15 10:30:45",
		},
		{
			name:       "複合的な時間減算",
			year:       2024.0,
			month:      1.0,
			day:        1.0,
			hour:       0.0,
			minute:     0.0,
			second:     0.0,
			subYears:   0.0,
			subMonths:  0.0,
			subDays:    0.0,
			subHours:   1.0,
			subMinutes: 1.0,
			subSeconds: 1.0,
			expected:   "2023-12-31 22:58:59",
		},
		{
			name:       "すべての単位を減算",
			year:       2024.0,
			month:      7.0,
			day:        18.0,
			hour:       14.0,
			minute:     35.0,
			second:     51.0,
			subYears:   1.0,
			subMonths:  2.0,
			subDays:    3.0,
			subHours:   4.0,
			subMinutes: 5.0,
			subSeconds: 6.0,
			expected:   "2023-05-15 10:30:45",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.SubtractDatetimeFloat(tt.year, tt.month, tt.day, tt.hour, tt.minute, tt.second, tt.subYears, tt.subMonths, tt.subDays, tt.subHours, tt.subMinutes, tt.subSeconds)
			if result != tt.expected {
				t.Errorf("SubtractDatetimeFloat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDiffTime(t *testing.T) {
	calc := &DatetimeCalculator{}

	tests := []struct {
		name     string
		t1       time.Time
		t2       time.Time
		expected string
	}{
		{
			name:     "1日の差",
			t1:       time.Date(2023, time.May, 16, 10, 30, 45, 0, time.UTC),
			t2:       time.Date(2023, time.May, 15, 10, 30, 45, 0, time.UTC),
			expected: "24h0m0s",
		},
		{
			name:     "1時間の差",
			t1:       time.Date(2023, time.May, 15, 11, 30, 45, 0, time.UTC),
			t2:       time.Date(2023, time.May, 15, 10, 30, 45, 0, time.UTC),
			expected: "1h0m0s",
		},
		{
			name:     "複合的な時間差",
			t1:       time.Date(2023, time.May, 15, 11, 45, 30, 0, time.UTC),
			t2:       time.Date(2023, time.May, 15, 10, 30, 45, 0, time.UTC),
			expected: "1h14m45s",
		},
		{
			name:     "負の時間差",
			t1:       time.Date(2023, time.May, 15, 10, 30, 45, 0, time.UTC),
			t2:       time.Date(2023, time.May, 15, 11, 30, 45, 0, time.UTC),
			expected: "-1h0m0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.diffTime(tt.t1, tt.t2)
			if result != tt.expected {
				t.Errorf("DiffTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestDatetimeCalculator_convertTimeUnit_Normal は時間単位変換の正常系テスト
func TestDatetimeCalculator_convertTimeUnit_Normal(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name       string
		value      float64
		inputUnit  string
		outputUnit string
		expected   float64
	}{
		{"SecondToMinute", 60, "second", "minute", 1},
		{"MinuteToSecond", 1, "minute", "second", 60},
		{"HourToMinute", 1, "hour", "minute", 60},
		{"MinuteToHour", 60, "minute", "hour", 1},
		{"DayToHour", 1, "day", "hour", 24},
		{"HourToDay", 24, "hour", "day", 1},
		{"MonthToDay", 1, "month", "day", 30},
		{"DayToMonth", 30, "day", "month", 1},
		{"YearToMonth", 1, "year", "month", 12.166667}, // 365/30 = 12.166667
		{"MonthToYear", 12, "month", "year", 0.986301}, // 12*30/365 = 0.986301
		{"YearToDay", 1, "year", "day", 365},
		{"DayToYear", 365, "day", "year", 1},
		{"SameUnit", 100, "hour", "hour", 100},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := calculator.ConvertTimeUnit(tc.value, tc.inputUnit, tc.outputUnit)

			// Assert
			if err != nil {
				t.Errorf("convertTimeUnit returned error: %v", err)
			}
			// 浮動小数点の精度問題を考慮して、小さな誤差を許容
			tolerance := 0.000001
			if math.Abs(result-tc.expected) > tolerance {
				t.Errorf("Expected %f, got %f (difference: %f)", tc.expected, result, math.Abs(result-tc.expected))
			}
		})
	}
}

// TestDatetimeCalculator_convertTimeUnit_ComplexConversions は複雑な時間単位変換のテスト
func TestDatetimeCalculator_convertTimeUnit_ComplexConversions(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name       string
		value      float64
		inputUnit  string
		outputUnit string
		expected   float64
	}{
		{"SecondToHour", 3600, "second", "hour", 1},
		{"HourToSecond", 1, "hour", "second", 3600},
		{"SecondToDay", 86400, "second", "day", 1},
		{"DayToSecond", 1, "day", "second", 86400},
		{"MinuteToDay", 1440, "minute", "day", 1},
		{"DayToMinute", 1, "day", "minute", 1440},
		{"HourToMonth", 720, "hour", "month", 1}, // 720 / (30 * 24) = 1
		{"MonthToHour", 1, "month", "hour", 720},
		{"SecondToYear", 31536000, "second", "year", 1}, // 365 * 24 * 60 * 60
		{"YearToSecond", 1, "year", "second", 31536000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := calculator.ConvertTimeUnit(tc.value, tc.inputUnit, tc.outputUnit)

			// Assert
			if err != nil {
				t.Errorf("convertTimeUnit returned error: %v", err)
			}
			// 浮動小数点の精度問題を考慮して、小さな誤差を許容
			tolerance := 0.000001
			if math.Abs(result-tc.expected) > tolerance {
				t.Errorf("Expected %f, got %f (difference: %f)", tc.expected, result, math.Abs(result-tc.expected))
			}
		})
	}
}

// TestDatetimeCalculator_convertTimeUnit_DecimalValues は小数値の時間単位変換テスト
func TestDatetimeCalculator_convertTimeUnit_DecimalValues(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	testCases := []struct {
		name       string
		value      float64
		inputUnit  string
		outputUnit string
		expected   float64
	}{
		{"HalfHourToMinute", 0.5, "hour", "minute", 30},
		{"TwoAndHalfHoursToMinute", 2.5, "hour", "minute", 150},
		{"OneAndHalfDayToHour", 1.5, "day", "hour", 36},
		{"QuarterDayToMinute", 0.25, "day", "minute", 360}, // 0.25 * 24 * 60
		{"HalfMonthToDay", 0.5, "month", "day", 15},
		{"ThreeQuarterYearToMonth", 0.75, "year", "month", 9.125}, // 0.75 * 365 / 30 = 9.125
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := calculator.ConvertTimeUnit(tc.value, tc.inputUnit, tc.outputUnit)

			// Assert
			if err != nil {
				t.Errorf("convertTimeUnit returned error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %f, got %f", tc.expected, result)
			}
		})
	}
}

// TestDatetimeCalculator_convertTimeUnit_InvalidInputUnit は無効な入力単位のテスト
func TestDatetimeCalculator_convertTimeUnit_InvalidInputUnit(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	// Act
	result, err := calculator.ConvertTimeUnit(100, "invalid", "hour")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid input unit, got nil")
	}
	if result != 0 {
		t.Errorf("Expected result 0 for error case, got %f", result)
	}
}

// TestDatetimeCalculator_convertTimeUnit_InvalidOutputUnit は無効な出力単位のテスト
func TestDatetimeCalculator_convertTimeUnit_InvalidOutputUnit(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	// Act
	result, err := calculator.ConvertTimeUnit(100, "hour", "invalid")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid output unit, got nil")
	}
	if result != 0 {
		t.Errorf("Expected result 0 for error case, got %f", result)
	}
}

// TestDatetimeCalculator_convertTimeUnit_ZeroValue はゼロ値の変換テスト
func TestDatetimeCalculator_convertTimeUnit_ZeroValue(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	// Act
	result, err := calculator.ConvertTimeUnit(0, "hour", "minute")

	// Assert
	if err != nil {
		t.Errorf("convertTimeUnit returned error: %v", err)
	}
	if result != 0 {
		t.Errorf("Expected 0, got %f", result)
	}
}

// TestDatetimeCalculator_convertTimeUnit_NegativeValue は負の値の変換テスト
func TestDatetimeCalculator_convertTimeUnit_NegativeValue(t *testing.T) {
	// Arrange
	calculator := &DatetimeCalculator{}

	// Act
	result, err := calculator.ConvertTimeUnit(-2, "hour", "minute")

	// Assert
	if err != nil {
		t.Errorf("convertTimeUnit returned error: %v", err)
	}
	if result != -120 {
		t.Errorf("Expected -120, got %f", result)
	}
}
