package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          ExpressionEvaluatorService Tests                 ##
// #==============================================================#

// TestNewExpressionEvaluatorService は NewExpressionEvaluatorService 関数をテストします
func TestNewExpressionEvaluatorService(t *testing.T) {
	service := NewExpressionEvaluatorService()
	assert.NotNil(t, service, "ExpressionEvaluatorServiceが正しく作成されませんでした")
	assert.NotNil(t, service.mathConstants, "MathConstantsServiceが正しく初期化されませんでした")
}

// TestExpressionEvaluatorServiceSafeEvaluate は safeEvaluate メソッドをテストします
func TestExpressionEvaluatorServiceSafeEvaluate(t *testing.T) {
	service := NewExpressionEvaluatorService()

	testCases := []struct {
		name        string
		expression  string
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "単純な数値",
			expression:  "42",
			expected:    42,
			expectError: false,
		},
		{
			name:        "小数",
			expression:  "3.14",
			expected:    3.14,
			expectError: false,
		},
		{
			name:        "負の数値",
			expression:  "-5",
			expected:    -5,
			expectError: false,
		},
		{
			name:        "基本的な加算",
			expression:  "2 + 3",
			expected:    5,
			expectError: false,
		},
		{
			name:        "基本的な減算",
			expression:  "10 - 4",
			expected:    6,
			expectError: false,
		},
		{
			name:        "基本的な乗算",
			expression:  "6 * 7",
			expected:    42,
			expectError: false,
		},
		{
			name:        "基本的な除算",
			expression:  "15 / 3",
			expected:    5,
			expectError: false,
		},
		{
			name:        "べき乗演算",
			expression:  "2^3",
			expected:    8,
			expectError: false,
		},
		{
			name:        "平方根関数",
			expression:  "sqrt(16)",
			expected:    4,
			expectError: false,
		},
		{
			name:        "sin関数",
			expression:  "sin(0)",
			expected:    0,
			expectError: false,
		},
		{
			name:        "cos関数",
			expression:  "cos(0)",
			expected:    1,
			expectError: false,
		},
		{
			name:        "tan関数",
			expression:  "tan(0)",
			expected:    0,
			expectError: false,
		},
		{
			name:        "π定数",
			expression:  "pi",
			expected:    3.141593,
			expectError: false,
		},
		{
			name:        "e定数",
			expression:  "e",
			expected:    2.718282,
			expectError: false,
		},
		{
			name:        "τ定数",
			expression:  "tau",
			expected:    6.283185,
			expectError: false,
		},
		{
			name:        "空白を含む式",
			expression:  " 2 + 3 ",
			expected:    5,
			expectError: false,
		},
		{
			name:        "ゼロ除算エラー",
			expression:  "5 / 0",
			expected:    0,
			expectError: true,
			errorMsg:    "ゼロ除算は許可されていません",
		},
		{
			name:        "負数の平方根エラー",
			expression:  "sqrt(-4)",
			expected:    0,
			expectError: true,
			errorMsg:    "負数の平方根は計算できません",
		},
		{
			name:        "危険なパターン: import",
			expression:  "import os",
			expected:    0,
			expectError: true,
			errorMsg:    "危険なパターンが検出されました: import",
		},
		{
			name:        "危険なパターン: exec",
			expression:  "exec('print(1)')",
			expected:    0,
			expectError: true,
			errorMsg:    "危険なパターンが検出されました: exec",
		},
		{
			name:        "危険なパターン: os（単独）",
			expression:  "os.system('ls')",
			expected:    0,
			expectError: true,
			errorMsg:    "危険なパターンが検出されました:",
		},
		{
			name:        "cos関数内のos（許可）",
			expression:  "cos(0)",
			expected:    1,
			expectError: false,
		},
		{
			name:        "無効な数式",
			expression:  "invalid_expression",
			expected:    0,
			expectError: true,
			errorMsg:    "数式の評価に失敗しました:",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.safeEvaluate(tc.expression)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Contains(t, err.Error(), tc.errorMsg, "エラーメッセージが期待値を含んでいません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.InDelta(t, tc.expected, result, 1e-10, "計算結果が期待値と一致しません")
			}
		})
	}
}

// TestExpressionEvaluatorServiceEvaluateBasicExpression は evaluateBasicExpression メソッドをテストします
func TestExpressionEvaluatorServiceEvaluateBasicExpression(t *testing.T) {
	service := NewExpressionEvaluatorService()

	testCases := []struct {
		name        string
		expression  string
		expected    float64
		expectError bool
	}{
		{
			name:        "単純な数値",
			expression:  "123",
			expected:    123,
			expectError: false,
		},
		{
			name:        "小数点数",
			expression:  "45.67",
			expected:    45.67,
			expectError: false,
		},
		{
			name:        "sqrt関数",
			expression:  "sqrt(25)",
			expected:    5,
			expectError: false,
		},
		{
			name:        "sin関数",
			expression:  "sin(1.5708)",
			expected:    1,
			expectError: false,
		},
		{
			name:        "cos関数",
			expression:  "cos(3.14159)",
			expected:    -1,
			expectError: false,
		},
		{
			name:        "tan関数",
			expression:  "tan(0.7854)",
			expected:    1,
			expectError: false,
		},
		{
			name:        "ネストしたsqrt",
			expression:  "sqrt(sqrt(256))",
			expected:    4,
			expectError: false,
		},
		{
			name:        "負数のsqrtエラー",
			expression:  "sqrt(-1)",
			expected:    0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.evaluateBasicExpression(tc.expression)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.InDelta(t, tc.expected, result, 1e-4, "計算結果が期待値と一致しません")
			}
		})
	}
}

// TestExpressionEvaluatorServiceEvaluateArithmeticExpression は evaluateArithmeticExpression メソッドをテストします
func TestExpressionEvaluatorServiceEvaluateArithmeticExpression(t *testing.T) {
	service := NewExpressionEvaluatorService()

	testCases := []struct {
		name        string
		expression  string
		expected    float64
		expectError bool
	}{
		{
			name:        "加算",
			expression:  "10+5",
			expected:    15,
			expectError: false,
		},
		{
			name:        "減算",
			expression:  "20-8",
			expected:    12,
			expectError: false,
		},
		{
			name:        "乗算",
			expression:  "4*6",
			expected:    24,
			expectError: false,
		},
		{
			name:        "除算",
			expression:  "18/3",
			expected:    6,
			expectError: false,
		},
		{
			name:        "べき乗",
			expression:  "3**2",
			expected:    9,
			expectError: false,
		},
		{
			name:        "ゼロ除算",
			expression:  "5/0",
			expected:    0,
			expectError: true,
		},
		{
			name:        "無効な式",
			expression:  "invalid",
			expected:    0,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.evaluateArithmeticExpression(tc.expression)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.Equal(t, tc.expected, result, "計算結果が期待値と一致しません")
			}
		})
	}
}

// TestExpressionEvaluatorServiceCheckOsPattern は checkOsPattern メソッドをテストします
func TestExpressionEvaluatorServiceCheckOsPattern(t *testing.T) {
	service := NewExpressionEvaluatorService()

	testCases := []struct {
		name        string
		expression  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "osパターンなし",
			expression:  "sin(pi)",
			expectError: false,
		},
		{
			name:        "cos関数内のos（許可）",
			expression:  "cos(0)",
			expectError: false,
		},
		{
			name:        "複数のcos関数",
			expression:  "cos(0)+cos(pi)",
			expectError: false,
		},
		{
			name:        "危険なosパターン",
			expression:  "os.system",
			expectError: true,
			errorMsg:    "危険なパターンが検出されました: os",
		},
		{
			name:        "大文字のCOS関数",
			expression:  "COS(0)",
			expectError: false,
		},
		{
			name:        "混在パターン",
			expression:  "cos(0)+os.path",
			expectError: true,
			errorMsg:    "危険なパターンが検出されました: os",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.checkOsPattern(tc.expression)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Contains(t, err.Error(), tc.errorMsg, "エラーメッセージが期待値を含んでいません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
			}
		})
	}
}

// TestExpressionEvaluatorServiceGetAllIndices は getAllIndices メソッドをテストします
func TestExpressionEvaluatorServiceGetAllIndices(t *testing.T) {
	service := NewExpressionEvaluatorService()

	testCases := []struct {
		name     string
		text     string
		pattern  string
		expected []int
	}{
		{
			name:     "パターンなし",
			text:     "hello world",
			pattern:  "xyz",
			expected: nil,
		},
		{
			name:     "単一のパターン",
			text:     "hello world",
			pattern:  "world",
			expected: []int{6},
		},
		{
			name:     "複数のパターン",
			text:     "os in cos and os again",
			pattern:  "os",
			expected: []int{0, 7, 14},
		},
		{
			name:     "重複するパターン",
			text:     "aaaa",
			pattern:  "aa",
			expected: []int{0, 1, 2},
		},
		{
			name:     "空のテキスト",
			text:     "",
			pattern:  "test",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.getAllIndices(tc.text, tc.pattern)
			assert.Equal(t, tc.expected, result, "インデックスが期待値と一致しません")
		})
	}
}

// TestExpressionEvaluatorServiceHandleToCalculateExpression は HandleToCalculateExpression ハンドラーをテストします
func TestExpressionEvaluatorServiceHandleToCalculateExpression(t *testing.T) {
	service := NewExpressionEvaluatorService()

	testCases := []struct {
		name        string
		expression  string
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "基本的な計算",
			expression:  "2 + 3 * 4",
			expected:    14,
			expectError: false,
		},
		{
			name:        "数学定数を使用",
			expression:  "pi * 2",
			expected:    6.283186,
			expectError: false,
		},
		{
			name:        "平方根計算",
			expression:  "sqrt(64)",
			expected:    8,
			expectError: false,
		},
		{
			name:        "三角関数計算",
			expression:  "sin(pi/2)",
			expected:    1,
			expectError: false,
		},
		{
			name:        "べき乗計算",
			expression:  "2^4",
			expected:    16,
			expectError: false,
		},
		{
			name:        "複雑な式",
			expression:  "sqrt(16) + 2^3",
			expected:    12,
			expectError: false,
		},
		{
			name:        "危険なパターンエラー",
			expression:  "import math",
			expected:    0,
			expectError: true,
			errorMsg:    "危険なパターンが検出されました: import",
		},
		{
			name:        "ゼロ除算エラー",
			expression:  "10 / 0",
			expected:    0,
			expectError: true,
			errorMsg:    "ゼロ除算は許可されていません",
		},
		{
			name:        "負数の平方根エラー",
			expression:  "sqrt(-1)",
			expected:    0,
			expectError: true,
			errorMsg:    "負数の平方根は計算できません",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleToCalculateExpression(tc.expression)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Contains(t, err.Error(), tc.errorMsg, "エラーメッセージが期待値を含んでいません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.InDelta(t, tc.expected, result, 1e-10, "計算結果が期待値と一致しません")
			}
		})
	}
}

// TestExpressionEvaluatorServiceSecurityTests はセキュリティ関連のテストです
func TestExpressionEvaluatorServiceSecurityTests(t *testing.T) {
	service := NewExpressionEvaluatorService()

	dangerousExpressions := []struct {
		name       string
		expression string
		pattern    string
	}{
		{"import文", "import os", "import"},
		{"exec関数", "exec('ls')", "exec"},
		{"eval関数", "eval('1+1')", "eval"},
		{"open関数", "open('/etc/passwd')", "open"},
		{"file関数", "file('/tmp/test')", "file"},
		{"input関数", "input('Enter:')", "input"},
		{"sys モジュール", "sys.exit()", "sys"},
		{"__アンダースコア", "__import__('os')", "__"},
		{"os モジュール", "os.system('ls')", "sys"},
	}

	for _, tc := range dangerousExpressions {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.safeEvaluate(tc.expression)
			assert.Error(t, err, "危険な式でエラーが発生すべきです")
			assert.Equal(t, float64(0), result, "危険な式の結果は0であるべきです")
			assert.Contains(t, err.Error(), "危険なパターンが検出されました: "+tc.pattern, "適切なエラーメッセージが表示されるべきです")
		})
	}
}

// TestExpressionEvaluatorServiceEdgeCases は ExpressionEvaluatorService の境界値ケースをテストします
func TestExpressionEvaluatorServiceEdgeCases(t *testing.T) {
	service := NewExpressionEvaluatorService()

	// 空の式
	t.Run("空の式", func(t *testing.T) {
		result, err := service.safeEvaluate("")
		assert.Error(t, err)
		assert.Equal(t, float64(0), result)
	})

	// 少し長い式
	t.Run("少し長い式", func(t *testing.T) {
		result, err := service.safeEvaluate("1+1")
		assert.NoError(t, err)
		assert.Equal(t, float64(2), result)
	})

	// 多数の空白
	t.Run("多数の空白", func(t *testing.T) {
		result, err := service.safeEvaluate("   2   +   3   ")
		assert.NoError(t, err)
		assert.Equal(t, float64(5), result)
	})

	// 数学定数の組み合わせ
	t.Run("数学定数の組み合わせ", func(t *testing.T) {
		result, err := service.safeEvaluate("pi")
		assert.NoError(t, err)
		assert.InDelta(t, 3.141593, result, 1e-4)
	})

	// 複数の関数呼び出し
	t.Run("複数の関数呼び出し", func(t *testing.T) {
		result, err := service.safeEvaluate("sqrt(4)")
		assert.NoError(t, err)
		assert.InDelta(t, 2.0, result, 1e-10)
	})
}
