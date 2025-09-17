package usecases

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          AdvancedMathService Tests                         ##
// #==============================================================#

// TestNewAdvancedMathService は NewAdvancedMathService 関数をテストします
func TestNewAdvancedMathService(t *testing.T) {
	service := NewAdvancedMathService()
	assert.NotNil(t, service, "AdvancedMathServiceが正しく作成されませんでした")
}

// TestAdvancedMathServicePower は power メソッドをテストします
func TestAdvancedMathServicePower(t *testing.T) {
	service := NewAdvancedMathService()

	testCases := []struct {
		name     string
		base     float64
		exponent float64
		expected float64
	}{
		{
			name:     "正の数のべき乗",
			base:     2,
			exponent: 3,
			expected: 8,
		},
		{
			name:     "負の数のべき乗",
			base:     -2,
			exponent: 3,
			expected: -8,
		},
		{
			name:     "小数のべき乗",
			base:     2.5,
			exponent: 2,
			expected: 6.25,
		},
		{
			name:     "0のべき乗",
			base:     0,
			exponent: 5,
			expected: 0,
		},
		{
			name:     "1のべき乗",
			base:     1,
			exponent: 100,
			expected: 1,
		},
		{
			name:     "負の指数",
			base:     2,
			exponent: -2,
			expected: 0.25,
		},
		{
			name:     "0乗",
			base:     5,
			exponent: 0,
			expected: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.power(tc.base, tc.exponent)
			assert.InDelta(t, tc.expected, result, 1e-10, "べき乗計算が期待値と一致しません")
		})
	}
}

// TestAdvancedMathServiceSquareRoot は squareRoot メソッドをテストします
func TestAdvancedMathServiceSquareRoot(t *testing.T) {
	service := NewAdvancedMathService()

	testCases := []struct {
		name        string
		number      float64
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "正の数の平方根",
			number:      16,
			expected:    4,
			expectError: false,
		},
		{
			name:        "0の平方根",
			number:      0,
			expected:    0,
			expectError: false,
		},
		{
			name:        "小数の平方根",
			number:      2.25,
			expected:    1.5,
			expectError: false,
		},
		{
			name:        "1の平方根",
			number:      1,
			expected:    1,
			expectError: false,
		},
		{
			name:        "負の数の平方根",
			number:      -4,
			expected:    0,
			expectError: true,
			errorMsg:    "負数の平方根は計算できません",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.squareRoot(tc.number)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Equal(t, tc.errorMsg, err.Error(), "エラーメッセージが期待値と一致しません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.InDelta(t, tc.expected, result, 1e-10, "平方根計算が期待値と一致しません")
			}
		})
	}
}

// TestAdvancedMathServiceFactorial は factorial メソッドをテストします
func TestAdvancedMathServiceFactorial(t *testing.T) {
	service := NewAdvancedMathService()

	testCases := []struct {
		name        string
		n           int
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "0の階乗",
			n:           0,
			expected:    1,
			expectError: false,
		},
		{
			name:        "1の階乗",
			n:           1,
			expected:    1,
			expectError: false,
		},
		{
			name:        "5の階乗",
			n:           5,
			expected:    120,
			expectError: false,
		},
		{
			name:        "10の階乗",
			n:           10,
			expected:    3628800,
			expectError: false,
		},
		{
			name:        "負の数の階乗",
			n:           -1,
			expected:    0,
			expectError: true,
			errorMsg:    "負数の階乗は定義されていません",
		},
		{
			name:        "大きすぎる数の階乗",
			n:           171,
			expected:    0,
			expectError: true,
			errorMsg:    "数値が大きすぎて階乗計算でオーバーフローします",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.factorial(tc.n)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Equal(t, tc.errorMsg, err.Error(), "エラーメッセージが期待値と一致しません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.Equal(t, tc.expected, result, "階乗計算が期待値と一致しません")
			}
		})
	}
}

// TestAdvancedMathServiceHandleToPower は HandleToPower ハンドラーをテストします
func TestAdvancedMathServiceHandleToPower(t *testing.T) {
	service := NewAdvancedMathService()

	testCases := []struct {
		name     string
		base     float64
		exponent float64
		expected float64
	}{
		{
			name:     "基本的なべき乗",
			base:     3,
			exponent: 4,
			expected: 81,
		},
		{
			name:     "小数のべき乗",
			base:     1.5,
			exponent: 3,
			expected: 3.375,
		},
		{
			name:     "負の指数",
			base:     4,
			exponent: -1,
			expected: 0.25,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleToPower(tc.base, tc.exponent)
			assert.NoError(t, err, "エラーが発生すべきではありません")
			assert.InDelta(t, tc.expected, result, 1e-10, "べき乗計算が期待値と一致しません")
		})
	}
}

// TestAdvancedMathServiceHandleToSquareRoot は HandleToSquareRoot ハンドラーをテストします
func TestAdvancedMathServiceHandleToSquareRoot(t *testing.T) {
	service := NewAdvancedMathService()

	testCases := []struct {
		name        string
		number      float64
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "正の数の平方根",
			number:      25,
			expected:    5,
			expectError: false,
		},
		{
			name:        "小数の平方根",
			number:      6.25,
			expected:    2.5,
			expectError: false,
		},
		{
			name:        "負の数の平方根エラー",
			number:      -9,
			expected:    0,
			expectError: true,
			errorMsg:    "負数の平方根は計算できません",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleToSquareRoot(tc.number)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Equal(t, tc.errorMsg, err.Error(), "エラーメッセージが期待値と一致しません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.InDelta(t, tc.expected, result, 1e-10, "平方根計算が期待値と一致しません")
			}
		})
	}
}

// TestAdvancedMathServiceHandleToFactorial は HandleToFactorial ハンドラーをテストします
func TestAdvancedMathServiceHandleToFactorial(t *testing.T) {
	service := NewAdvancedMathService()

	testCases := []struct {
		name        string
		n           int
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "基本的な階乗",
			n:           6,
			expected:    720,
			expectError: false,
		},
		{
			name:        "0の階乗",
			n:           0,
			expected:    1,
			expectError: false,
		},
		{
			name:        "負の数の階乗エラー",
			n:           -5,
			expected:    0,
			expectError: true,
			errorMsg:    "負数の階乗は定義されていません",
		},
		{
			name:        "オーバーフローエラー",
			n:           200,
			expected:    0,
			expectError: true,
			errorMsg:    "数値が大きすぎて階乗計算でオーバーフローします",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleToFactorial(tc.n)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Equal(t, tc.errorMsg, err.Error(), "エラーメッセージが期待値と一致しません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.Equal(t, tc.expected, result, "階乗計算が期待値と一致しません")
			}
		})
	}
}

// TestAdvancedMathServiceEdgeCases は AdvancedMathService の境界値ケースをテストします
func TestAdvancedMathServiceEdgeCases(t *testing.T) {
	service := NewAdvancedMathService()

	// 非常に大きな数のべき乗
	t.Run("非常に大きな数のべき乗", func(t *testing.T) {
		result := service.power(10, 10)
		assert.Equal(t, 1e10, result)
	})

	// 非常に小さな数の平方根
	t.Run("非常に小さな数の平方根", func(t *testing.T) {
		result, err := service.squareRoot(1e-10)
		assert.NoError(t, err)
		assert.InDelta(t, 1e-5, result, 1e-15)
	})

	// 境界値の階乗（170は計算可能な最大値）
	t.Run("境界値の階乗", func(t *testing.T) {
		result, err := service.factorial(170)
		assert.NoError(t, err)
		assert.True(t, result > 0, "170の階乗は正の値であるべきです")
		assert.False(t, math.IsInf(result, 1), "170の階乗は無限大であってはいけません")
	})
}
