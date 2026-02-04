package reversePolishNotation

import (
	"math"
	"testing"
)

func TestEvaluateSuccess(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		expression string
		expected   float64
	}{
		{
			name:       "単純な加算と乗算",
			expression: "2+3*4",
			expected:   14,
		},
		{
			name:       "括弧と演算子優先順位",
			expression: "(2+3)*4",
			expected:   20,
		},
		{
			name:       "先頭の単項マイナス",
			expression: "-5+3",
			expected:   -2,
		},
		{
			name:       "括弧にかかる単項マイナス",
			expression: "-(2+3)",
			expected:   -5,
		},
		{
			name:       "小数の計算",
			expression: "3.5+0.5",
			expected:   4,
		},
		{
			name:       "平方根と加算",
			expression: "sqrt(16)+2",
			expected:   6,
		},
		{
			name:       "三角関数 sin",
			expression: "sin(0)",
			expected:   0,
		},
		{
			name:       "三角関数 cos",
			expression: "cos(0)",
			expected:   1,
		},
		{
			name:       "三角関数 tan",
			expression: "tan(0)",
			expected:   0,
		},
		{
			name:       "累乗演算",
			expression: "2^3",
			expected:   8,
		},
		{
			name:       "複合式",
			expression: "sqrt(9) + cos(0) + 1",
			expected:   5,
		},
	}

	const tolerance = 1e-9

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := Evaluate(tc.expression)
			if err != nil {
				t.Fatalf("期待しないエラー: %v", err)
			}

			if math.Abs(result-tc.expected) > tolerance {
				t.Fatalf("計算結果が一致しません: got=%v want=%v", result, tc.expected)
			}
		})
	}
}

func TestEvaluateError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		expression string
	}{
		{
			name:       "空文字列",
			expression: "",
		},
		{
			name:       "不正な文字",
			expression: "2@3",
		},
		{
			name:       "括弧の不一致",
			expression: "2+(3*4",
		},
		{
			name:       "ゼロ除算",
			expression: "10/0",
		},
		{
			name:       "未対応の識別子",
			expression: "foo(1)",
		},
		{
			name:       "負数の平方根",
			expression: "sqrt(-1)",
		},
		{
			name:       "演算子の並びが不正",
			expression: "2+*3",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if _, err := Evaluate(tc.expression); err == nil {
				t.Fatalf("エラーを期待しましたが成功しました: expression=%s", tc.expression)
			}
		})
	}
}
