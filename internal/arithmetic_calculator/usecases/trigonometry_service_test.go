package usecases

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          TrigonometryService Tests                         ##
// #==============================================================#

// TestNewTrigonometryService は NewTrigonometryService 関数をテストします
func TestNewTrigonometryService(t *testing.T) {
	service := NewTrigonometryService()
	assert.NotNil(t, service, "TrigonometryServiceが正しく作成されませんでした")
}

// TestTrigonometryServiceTrigonometry は trigonometry メソッドをテストします
func TestTrigonometryServiceTrigonometry(t *testing.T) {
	service := NewTrigonometryService()

	testCases := []struct {
		name        string
		function    string
		angle       float64
		unit        string
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "sin(30度)",
			function:    "sin",
			angle:       30,
			unit:        "degrees",
			expected:    0.5,
			expectError: false,
		},
		{
			name:        "cos(60度)",
			function:    "cos",
			angle:       60,
			unit:        "degrees",
			expected:    0.5,
			expectError: false,
		},
		{
			name:        "tan(45度)",
			function:    "tan",
			angle:       45,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "sin(π/2ラジアン)",
			function:    "sin",
			angle:       math.Pi / 2,
			unit:        "radians",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "cos(πラジアン)",
			function:    "cos",
			angle:       math.Pi,
			unit:        "radians",
			expected:    -1.0,
			expectError: false,
		},
		{
			name:        "tan(π/4ラジアン)",
			function:    "tan",
			angle:       math.Pi / 4,
			unit:        "radians",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "sin(0度)",
			function:    "sin",
			angle:       0,
			unit:        "degrees",
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "cos(0度)",
			function:    "cos",
			angle:       0,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "tan(0度)",
			function:    "tan",
			angle:       0,
			unit:        "degrees",
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "SIN(大文字)",
			function:    "SIN",
			angle:       90,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "COS(大文字)",
			function:    "COS",
			angle:       0,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "TAN(大文字)",
			function:    "TAN",
			angle:       45,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "DEGREES(大文字単位)",
			function:    "sin",
			angle:       90,
			unit:        "DEGREES",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "無効な関数名",
			function:    "invalid",
			angle:       45,
			unit:        "degrees",
			expected:    0,
			expectError: true,
			errorMsg:    "未知の三角関数です: invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.trigonometry(tc.function, tc.angle, tc.unit)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Equal(t, tc.errorMsg, err.Error(), "エラーメッセージが期待値と一致しません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.InDelta(t, tc.expected, result, 1e-10, "三角関数計算が期待値と一致しません")
			}
		})
	}
}

// TestTrigonometryServiceHandleToTrigonometry は HandleToTrigonometry ハンドラーをテストします
func TestTrigonometryServiceHandleToTrigonometry(t *testing.T) {
	service := NewTrigonometryService()

	testCases := []struct {
		name        string
		function    string
		angle       float64
		unit        string
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "基本的なsin計算",
			function:    "sin",
			angle:       90,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "基本的なcos計算",
			function:    "cos",
			angle:       0,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "基本的なtan計算",
			function:    "tan",
			angle:       45,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "ラジアン単位でのsin計算",
			function:    "sin",
			angle:       math.Pi / 6,
			unit:        "radians",
			expected:    0.5,
			expectError: false,
		},
		{
			name:        "ラジアン単位でのcos計算",
			function:    "cos",
			angle:       math.Pi / 3,
			unit:        "radians",
			expected:    0.5,
			expectError: false,
		},
		{
			name:        "ラジアン単位でのtan計算",
			function:    "tan",
			angle:       math.Pi / 4,
			unit:        "radians",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "負の角度",
			function:    "sin",
			angle:       -30,
			unit:        "degrees",
			expected:    -0.5,
			expectError: false,
		},
		{
			name:        "大きな角度",
			function:    "sin",
			angle:       450,
			unit:        "degrees",
			expected:    1.0,
			expectError: false,
		},
		{
			name:        "無効な関数名エラー",
			function:    "log",
			angle:       45,
			unit:        "degrees",
			expected:    0,
			expectError: true,
			errorMsg:    "未知の三角関数です: log",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleToTrigonometry(tc.function, tc.angle, tc.unit)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Equal(t, tc.errorMsg, err.Error(), "エラーメッセージが期待値と一致しません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.InDelta(t, tc.expected, result, 1e-10, "三角関数計算が期待値と一致しません")
			}
		})
	}
}

// TestTrigonometryServiceSpecialAngles は特別な角度での三角関数をテストします
func TestTrigonometryServiceSpecialAngles(t *testing.T) {
	service := NewTrigonometryService()

	// 特別な角度のテストケース
	specialCases := []struct {
		name     string
		function string
		angle    float64
		unit     string
		expected float64
	}{
		// 0度/0ラジアン
		{"sin(0°)", "sin", 0, "degrees", 0},
		{"cos(0°)", "cos", 0, "degrees", 1},
		{"tan(0°)", "tan", 0, "degrees", 0},

		// 30度/π/6ラジアン
		{"sin(30°)", "sin", 30, "degrees", 0.5},
		{"cos(30°)", "cos", 30, "degrees", math.Sqrt(3) / 2},
		{"tan(30°)", "tan", 30, "degrees", 1 / math.Sqrt(3)},

		// 45度/π/4ラジアン
		{"sin(45°)", "sin", 45, "degrees", math.Sqrt(2) / 2},
		{"cos(45°)", "cos", 45, "degrees", math.Sqrt(2) / 2},
		{"tan(45°)", "tan", 45, "degrees", 1},

		// 60度/π/3ラジアン
		{"sin(60°)", "sin", 60, "degrees", math.Sqrt(3) / 2},
		{"cos(60°)", "cos", 60, "degrees", 0.5},
		{"tan(60°)", "tan", 60, "degrees", math.Sqrt(3)},

		// 90度/π/2ラジアン
		{"sin(90°)", "sin", 90, "degrees", 1},
		{"cos(90°)", "cos", 90, "degrees", 0},

		// 180度/πラジアン
		{"sin(180°)", "sin", 180, "degrees", 0},
		{"cos(180°)", "cos", 180, "degrees", -1},
		{"tan(180°)", "tan", 180, "degrees", 0},

		// 270度/3π/2ラジアン
		{"sin(270°)", "sin", 270, "degrees", -1},
		{"cos(270°)", "cos", 270, "degrees", 0},

		// 360度/2πラジアン
		{"sin(360°)", "sin", 360, "degrees", 0},
		{"cos(360°)", "cos", 360, "degrees", 1},
		{"tan(360°)", "tan", 360, "degrees", 0},
	}

	for _, tc := range specialCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.trigonometry(tc.function, tc.angle, tc.unit)
			assert.NoError(t, err, "エラーが発生すべきではありません")
			assert.InDelta(t, tc.expected, result, 1e-10, "特別な角度の三角関数計算が期待値と一致しません")
		})
	}
}

// TestTrigonometryServiceUnitConversion は単位変換をテストします
func TestTrigonometryServiceUnitConversion(t *testing.T) {
	service := NewTrigonometryService()

	// 同じ角度を度数とラジアンで計算して結果が同じことを確認
	conversionCases := []struct {
		name         string
		function     string
		angleDegrees float64
		angleRadians float64
	}{
		{"30度とπ/6ラジアン", "sin", 30, math.Pi / 6},
		{"45度とπ/4ラジアン", "cos", 45, math.Pi / 4},
		{"60度とπ/3ラジアン", "tan", 60, math.Pi / 3},
		{"90度とπ/2ラジアン", "sin", 90, math.Pi / 2},
		{"180度とπラジアン", "cos", 180, math.Pi},
	}

	for _, tc := range conversionCases {
		t.Run(tc.name, func(t *testing.T) {
			resultDegrees, err1 := service.trigonometry(tc.function, tc.angleDegrees, "degrees")
			resultRadians, err2 := service.trigonometry(tc.function, tc.angleRadians, "radians")

			assert.NoError(t, err1, "度数での計算でエラーが発生すべきではありません")
			assert.NoError(t, err2, "ラジアンでの計算でエラーが発生すべきではありません")
			assert.InDelta(t, resultDegrees, resultRadians, 1e-10, "度数とラジアンでの計算結果が一致しません")
		})
	}
}

// TestTrigonometryServiceEdgeCases は TrigonometryService の境界値ケースをテストします
func TestTrigonometryServiceEdgeCases(t *testing.T) {
	service := NewTrigonometryService()

	// 非常に大きな角度
	t.Run("非常に大きな角度", func(t *testing.T) {
		result, err := service.trigonometry("sin", 3600, "degrees") // 10回転
		assert.NoError(t, err)
		assert.InDelta(t, 0, result, 1e-10, "3600度のsinは0に近いはずです")
	})

	// 非常に小さな角度
	t.Run("非常に小さな角度", func(t *testing.T) {
		result, err := service.trigonometry("sin", 0.001, "degrees")
		assert.NoError(t, err)
		assert.True(t, result > 0 && result < 0.1, "非常に小さな角度のsinは小さな正の値であるべきです")
	})

	// 負の大きな角度
	t.Run("負の大きな角度", func(t *testing.T) {
		result, err := service.trigonometry("cos", -720, "degrees") // -2回転
		assert.NoError(t, err)
		assert.InDelta(t, 1, result, 1e-10, "-720度のcosは1に近いはずです")
	})

	// ラジアンでの大きな値
	t.Run("ラジアンでの大きな値", func(t *testing.T) {
		result, err := service.trigonometry("sin", 4*math.Pi, "radians") // 2回転
		assert.NoError(t, err)
		assert.InDelta(t, 0, result, 1e-10, "4πラジアンのsinは0に近いはずです")
	})
}

// TestTrigonometryServiceCaseInsensitive は大文字小文字を区別しないことをテストします
func TestTrigonometryServiceCaseInsensitive(t *testing.T) {
	service := NewTrigonometryService()

	// 関数名の大文字小文字
	functionCases := []string{"sin", "SIN", "Sin", "sIn"}
	for _, funcName := range functionCases {
		t.Run("関数名: "+funcName, func(t *testing.T) {
			result, err := service.trigonometry(funcName, 30, "degrees")
			assert.NoError(t, err)
			assert.InDelta(t, 0.5, result, 1e-10)
		})
	}

	// 単位の大文字小文字
	unitCases := []string{"degrees", "DEGREES", "Degrees", "dEgReEs"}
	for _, unit := range unitCases {
		t.Run("単位: "+unit, func(t *testing.T) {
			result, err := service.trigonometry("sin", 30, unit)
			assert.NoError(t, err)
			assert.InDelta(t, 0.5, result, 1e-10)
		})
	}
}
