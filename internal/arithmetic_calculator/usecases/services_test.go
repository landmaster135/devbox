package usecases

import (
	"errors"
	"math"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// モック用のFileOpenerインターフェース実装
type MockFileOpener struct {
	errToReturn error
	file        *os.File
}

func (m *MockFileOpener) Open(name string) (*os.File, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	if m.file != nil {
		return m.file, nil
	}
	// 実際のファイルを開く（通常のテストケースで使用）
	return os.Open(name)
}

// モック用のBufioScannerインターフェース実装
type MockBufioScanner struct {
	boolToReturn bool
	errToReturn  error
}

func (m *MockBufioScanner) Scan() bool {
	return m.boolToReturn
}

func (m *MockBufioScanner) Err() error {
	return m.errToReturn
}

// モック用のJSONMarshalerインターフェース実装
type MockJSONMarshaler struct {
	errToReturn error
}

func (m *MockJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	return []byte(`{"mocked":"json"}`), nil
}

// テスト用の一時ファイルを作成する関数
func createTempFileWithLines(t *testing.T, lines []string) string {
	// 一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "test-line-count-*.txt")
	assert.NoError(t, err, "一時ファイルの作成に失敗しました")

	// テスト終了時に一時ファイルを削除
	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	// ファイルに行を書き込む
	for _, line := range lines {
		_, err := tmpFile.WriteString(line + "\n")
		assert.NoError(t, err, "ファイルへの書き込みに失敗しました")
	}

	// ファイルを閉じる
	err = tmpFile.Close()
	assert.NoError(t, err, "ファイルのクローズに失敗しました")

	return tmpFile.Name()
}

// #==============================================================#
// ##          CalculatorService Tests                           ##
// #==============================================================#

// TestCalculatorServiceAdd は CalculatorService の Add メソッドをテストします
func TestCalculatorServiceAdd(t *testing.T) {
	// テスト用の CalculatorService インスタンスを作成
	service := NewCalculatorService()

	// テストケースを定義
	testCases := []struct {
		name     string
		x        float64
		y        float64
		expected float64
	}{
		{
			name:     "正の数の加算",
			x:        5,
			y:        3,
			expected: 8,
		},
		{
			name:     "負の数の加算",
			x:        -2,
			y:        -3,
			expected: -5,
		},
		{
			name:     "正と負の数の加算",
			x:        5,
			y:        -3,
			expected: 2,
		},
		{
			name:     "ゼロとの加算",
			x:        5,
			y:        0,
			expected: 5,
		},
		{
			name:     "小数点数の加算",
			x:        2.5,
			y:        3.5,
			expected: 6.0,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.Add(tc.x, tc.y)
			assert.Equal(t, tc.expected, result, "%f + %f should equal %f", tc.x, tc.y, tc.expected)
		})
	}
}

// TestCalculatorServiceSubtract は CalculatorService の Subtract メソッドをテストします
func TestCalculatorServiceSubtract(t *testing.T) {
	// テスト用の CalculatorService インスタンスを作成
	service := NewCalculatorService()

	// テストケースを定義
	testCases := []struct {
		name     string
		x        float64
		y        float64
		expected float64
	}{
		{
			name:     "正の数の減算",
			x:        10,
			y:        4,
			expected: 6,
		},
		{
			name:     "負の数の減算",
			x:        -2,
			y:        -3,
			expected: 1,
		},
		{
			name:     "正と負の数の減算",
			x:        5,
			y:        -3,
			expected: 8,
		},
		{
			name:     "ゼロとの減算",
			x:        5,
			y:        0,
			expected: 5,
		},
		{
			name:     "自身からの減算",
			x:        5,
			y:        5,
			expected: 0,
		},
		{
			name:     "小数点数の減算",
			x:        5.5,
			y:        2.2,
			expected: 3.3,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.Subtract(tc.x, tc.y)
			// 小数点の計算誤差を考慮して、ほぼ等しいかをチェック
			if tc.name == "小数点数の減算" {
				assert.InDelta(t, tc.expected, result, 0.0001, "%f - %f should approximately equal %f", tc.x, tc.y, tc.expected)
			} else {
				assert.Equal(t, tc.expected, result, "%f - %f should equal %f", tc.x, tc.y, tc.expected)
			}
		})
	}
}

// TestCalculatorServiceMultiply は CalculatorService の Multiply メソッドをテストします
func TestCalculatorServiceMultiply(t *testing.T) {
	// テスト用の CalculatorService インスタンスを作成
	service := NewCalculatorService()

	// テストケースを定義
	testCases := []struct {
		name     string
		x        float64
		y        float64
		expected float64
	}{
		{
			name:     "正の数の乗算",
			x:        6,
			y:        7,
			expected: 42,
		},
		{
			name:     "負の数の乗算",
			x:        -3,
			y:        -4,
			expected: 12,
		},
		{
			name:     "正と負の数の乗算",
			x:        5,
			y:        -3,
			expected: -15,
		},
		{
			name:     "ゼロとの乗算",
			x:        5,
			y:        0,
			expected: 0,
		},
		{
			name:     "小数点数の乗算",
			x:        2.5,
			y:        4,
			expected: 10,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.Multiply(tc.x, tc.y)
			assert.Equal(t, tc.expected, result, "%f * %f should equal %f", tc.x, tc.y, tc.expected)
		})
	}
}

// TestCalculatorServiceDivide は CalculatorService の Divide メソッドをテストします
func TestCalculatorServiceDivide(t *testing.T) {
	// テスト用の CalculatorService インスタンスを作成
	service := NewCalculatorService()

	// テストケースを定義
	testCases := []struct {
		name     string
		x        float64
		y        float64
		expected float64
		isInf    bool
	}{
		{
			name:     "正の数の除算",
			x:        20,
			y:        5,
			expected: 4,
			isInf:    false,
		},
		{
			name:     "負の数の除算",
			x:        -12,
			y:        -4,
			expected: 3,
			isInf:    false,
		},
		{
			name:     "正と負の数の除算",
			x:        15,
			y:        -3,
			expected: -5,
			isInf:    false,
		},
		{
			name:     "ゼロの除算",
			x:        0,
			y:        5,
			expected: 0,
			isInf:    false,
		},
		{
			name:     "小数点数の除算",
			x:        10,
			y:        4,
			expected: 2.5,
			isInf:    false,
		},
		{
			name:     "ゼロによる除算",
			x:        5,
			y:        0,
			expected: 0, // 実際には無限大が返されるが、テストのために0を設定
			isInf:    true,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.Divide(tc.x, tc.y)

			if tc.isInf {
				// ゼロ除算の場合は無限大が返されることを確認
				assert.True(t, math.IsInf(result, 1), "%f / %f should return +Inf", tc.x, tc.y)
			} else {
				assert.Equal(t, tc.expected, result, "%f / %f should equal %f", tc.x, tc.y, tc.expected)
			}
		})
	}
}

// TestCalculatorServiceEdgeCases は CalculatorService の境界値ケースをテストします
func TestCalculatorServiceEdgeCases(t *testing.T) {
	// テスト用の CalculatorService インスタンスを作成
	service := NewCalculatorService()

	// 大きな数値の演算
	t.Run("大きな数値の加算", func(t *testing.T) {
		result := service.Add(1e15, 1e15)
		assert.Equal(t, 2e15, result)
	})

	// 非常に小さな数値の演算
	t.Run("非常に小さな数値の乗算", func(t *testing.T) {
		result := service.Multiply(1e-15, 1e-15)
		assert.Equal(t, 1e-30, result)
	})

	// 精度の問題を確認
	t.Run("精度の問題", func(t *testing.T) {
		// 0.1 + 0.2 は浮動小数点の精度の問題で厳密には 0.3 にならない
		result := service.Add(0.1, 0.2)
		assert.InDelta(t, 0.3, result, 1e-10)
	})
}

// TestCalculatorServiceSum は CalculatorService の Sum メソッドをテストします
func TestCalculatorServiceSum(t *testing.T) {
	// テスト用の CalculatorService インスタンスを作成
	service := NewCalculatorService()

	// テストケースを定義
	testCases := []struct {
		name     string
		arr      []float64
		expected float64
	}{
		{
			name:     "正の数の合計",
			arr:      []float64{1, 2, 3, 4, 5},
			expected: 15,
		},
		{
			name:     "負の数の合計",
			arr:      []float64{-1, -2, -3, -4, -5},
			expected: -15,
		},
		{
			name:     "正と負の数の合計",
			arr:      []float64{-5, -3, 0, 3, 5},
			expected: 0,
		},
		{
			name:     "小数点数の合計",
			arr:      []float64{1.5, 2.5, 3.5},
			expected: 7.5,
		},
		{
			name:     "空の配列",
			arr:      []float64{},
			expected: 0,
		},
		{
			name:     "単一要素の配列",
			arr:      []float64{42},
			expected: 42,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.Sum(tc.arr)
			assert.Equal(t, tc.expected, result, "Sum of %v should equal %f", tc.arr, tc.expected)
		})
	}
}

// #==============================================================#
// ##          FileEvaluatorService Tests                        ##
// #==============================================================#

// TestFileEvaluatorServiceCountLines は FileEvaluatorService の CountLines メソッドをテストします
func TestFileEvaluatorServiceCountLines(t *testing.T) {
	// テスト用の FileEvaluatorService インスタンスを作成
	service := NewFileEvaluatorService()

	// テストケースを定義
	testCases := []struct {
		name          string
		lines         []string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "空のファイル",
			lines:         []string{},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name:          "1行のファイル",
			lines:         []string{"これは1行のファイルです"},
			expectedCount: 1,
			expectError:   false,
		},
		{
			name:          "複数行のファイル",
			lines:         []string{"1行目", "2行目", "3行目", "4行目", "5行目"},
			expectedCount: 5,
			expectError:   false,
		},
		{
			name:          "空行を含むファイル",
			lines:         []string{"1行目", "", "3行目", "", "5行目"},
			expectedCount: 5,
			expectError:   false,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用の一時ファイルを作成
			filePath := createTempFileWithLines(t, tc.lines)

			// テスト対象の関数を実行
			count, err := service.CountLines(filePath)

			// 結果の検証
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, count, "行数が期待値と一致しません")
			}
		})
	}

	// 存在しないファイルのテスト
	t.Run("存在しないファイル", func(t *testing.T) {
		_, err := service.CountLines("/path/to/nonexistent/file.txt")
		assert.Error(t, err, "存在しないファイルに対してエラーが発生すべきです")
	})
}

// TestFileEvaluatorServiceIsLineCountGreaterThan は FileEvaluatorService の IsLineCountGreaterThan メソッドをテストします
func TestFileEvaluatorServiceIsLineCountGreaterThan(t *testing.T) {
	// テスト用の FileEvaluatorService インスタンスを作成
	service := NewFileEvaluatorService()

	// テストケースを定義
	testCases := []struct {
		name           string
		lines          []string
		threshold      int
		expectedResult bool
		expectedCount  int
		expectError    bool
	}{
		{
			name:           "閾値より少ない行数",
			lines:          []string{"1行目", "2行目", "3行目"},
			threshold:      5,
			expectedResult: false,
			expectedCount:  3,
			expectError:    false,
		},
		{
			name:           "閾値と同じ行数",
			lines:          []string{"1行目", "2行目", "3行目", "4行目", "5行目"},
			threshold:      5,
			expectedResult: false,
			expectedCount:  5,
			expectError:    false,
		},
		{
			name:           "閾値より多い行数",
			lines:          []string{"1行目", "2行目", "3行目", "4行目", "5行目", "6行目"},
			threshold:      5,
			expectedResult: true,
			expectedCount:  6,
			expectError:    false,
		},
		{
			name:           "閾値が0の場合",
			lines:          []string{"1行目"},
			threshold:      0,
			expectedResult: true,
			expectedCount:  1,
			expectError:    false,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用の一時ファイルを作成
			filePath := createTempFileWithLines(t, tc.lines)

			// テスト対象の関数を実行
			result, count, err := service.IsLineCountGreaterThan(filePath, tc.threshold)

			// 結果の検証
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result, "比較結果が期待値と一致しません")
				assert.Equal(t, tc.expectedCount, count, "行数が期待値と一致しません")
			}
		})
	}

	// 存在しないファイルのテスト
	t.Run("存在しないファイル", func(t *testing.T) {
		_, _, err := service.IsLineCountGreaterThan("/path/to/nonexistent/file.txt", 5)
		assert.Error(t, err, "存在しないファイルに対してエラーが発生すべきです")
	})
}

// TestFileEvaluatorServiceScannerError は scanner.Err() がエラーを返す場合をテストします
func TestFileEvaluatorServiceScannerError(t *testing.T) {
	// テスト用の一時ファイルを作成（実際のファイルが必要）
	tmpFile, err := os.CreateTemp("", "mock-file-*.txt")
	assert.NoError(t, err, "一時ファイルの作成に失敗しました")
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// モック用のBufioScannerを作成
	mockScanner := &MockBufioScanner{
		boolToReturn: false,
		errToReturn:  errors.New("模擬的なスキャナーエラー"),
	}

	// モック用のFileOpenerを作成（一時ファイルを返すように設定）
	mockFileOpener := &MockFileOpener{
		errToReturn: nil,
		file:        tmpFile,
	}

	// テスト用のFileEvaluatorServiceを作成し、モックを注入
	service := NewFileEvaluatorServiceWithDependencies(
		mockFileOpener,
		mockScanner,
		&DefaultJSONMarshaler{},
	)

	// CountLinesメソッドを呼び出す
	lc, err := service.CountLines("模擬的なファイルパス")

	// 結果の検証
	assert.Equal(t, 0, lc, "行数が期待値と一致しません")
	assert.Error(t, err, "scanner.Err()がエラーを返すべきです")
	assert.Equal(t, "模擬的なスキャナーエラー", err.Error(), "エラーメッセージが期待値と一致しません")
}
