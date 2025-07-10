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

// #==============================================================#
// ##          CalculatorService Handler Tests                   ##
// #==============================================================#

// TestCalculatorServiceHandleToCalculate は CalculatorService の HandleToCalculate メソッドをテストします
func TestCalculatorServiceHandleToCalculate(t *testing.T) {
	// テスト用の CalculatorService インスタンスを作成
	service := NewCalculatorService()

	// テストケースを定義
	testCases := []struct {
		name        string
		operation   string
		x           float64
		y           float64
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "加算の処理",
			operation:   "add",
			x:           5,
			y:           3,
			expected:    8,
			expectError: false,
		},
		{
			name:        "減算の処理",
			operation:   "subtract",
			x:           10,
			y:           4,
			expected:    6,
			expectError: false,
		},
		{
			name:        "乗算の処理",
			operation:   "multiply",
			x:           6,
			y:           7,
			expected:    42,
			expectError: false,
		},
		{
			name:        "除算の処理",
			operation:   "divide",
			x:           20,
			y:           5,
			expected:    4,
			expectError: false,
		},
		{
			name:        "ゼロ除算エラー",
			operation:   "divide",
			x:           5,
			y:           0,
			expected:    0,
			expectError: true,
			errorMsg:    "division by zero is not allowed",
		},
		{
			name:        "無効な演算子",
			operation:   "invalid",
			x:           5,
			y:           3,
			expected:    0,
			expectError: false,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleToCalculate(tc.operation, tc.x, tc.y)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Equal(t, tc.errorMsg, err.Error(), "エラーメッセージが期待値と一致しません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.Equal(t, tc.expected, result, "計算結果が期待値と一致しません")
			}
		})
	}
}

// TestCalculatorServiceHandleToCalculateWithArray は CalculatorService の HandleToCalculateWithArray メソッドをテストします
func TestCalculatorServiceHandleToCalculateWithArray(t *testing.T) {
	// テスト用の CalculatorService インスタンスを作成
	service := NewCalculatorService()

	// テストケースを定義
	testCases := []struct {
		name        string
		operation   string
		numbers     []float64
		expected    float64
		expectError bool
	}{
		{
			name:        "配列の合計",
			operation:   "sum",
			numbers:     []float64{1, 2, 3, 4, 5},
			expected:    15,
			expectError: false,
		},
		{
			name:        "空の配列の合計",
			operation:   "sum",
			numbers:     []float64{},
			expected:    0,
			expectError: false,
		},
		{
			name:        "負の数を含む配列の合計",
			operation:   "sum",
			numbers:     []float64{-5, -3, 0, 3, 5},
			expected:    0,
			expectError: false,
		},
		{
			name:        "小数点数の配列の合計",
			operation:   "sum",
			numbers:     []float64{1.5, 2.5, 3.5},
			expected:    7.5,
			expectError: false,
		},
		{
			name:        "無効な演算子",
			operation:   "invalid",
			numbers:     []float64{1, 2, 3},
			expected:    0,
			expectError: false,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleToCalculateWithArray(tc.operation, tc.numbers)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.Equal(t, tc.expected, result, "計算結果が期待値と一致しません")
			}
		})
	}
}

// #==============================================================#
// ##          FileEvaluatorService Handler Tests               ##
// #==============================================================#

// TestFileEvaluatorServiceHandleToEvaluateLineCount は FileEvaluatorService の HandleToEvaluateLineCount メソッドをテストします
func TestFileEvaluatorServiceHandleToEvaluateLineCount(t *testing.T) {
	// テスト用の FileEvaluatorService インスタンスを作成
	service := NewFileEvaluatorService()

	// テストケースを定義
	testCases := []struct {
		name        string
		lines       []string
		threshold   int
		expectError bool
		contains    []string // JSONレスポンスに含まれるべき文字列
	}{
		{
			name:        "閾値より少ない行数",
			lines:       []string{"1行目", "2行目", "3行目"},
			threshold:   5,
			expectError: false,
			contains:    []string{"\"is_greater\": false", "\"line_count\": 3", "\"threshold\": 5", "以下です"},
		},
		{
			name:        "閾値と同じ行数",
			lines:       []string{"1行目", "2行目", "3行目", "4行目", "5行目"},
			threshold:   5,
			expectError: false,
			contains:    []string{"\"is_greater\": false", "\"line_count\": 5", "\"threshold\": 5", "以下です"},
		},
		{
			name:        "閾値より多い行数",
			lines:       []string{"1行目", "2行目", "3行目", "4行目", "5行目", "6行目"},
			threshold:   5,
			expectError: false,
			contains:    []string{"\"is_greater\": true", "\"line_count\": 6", "\"threshold\": 5", "より大きいです"},
		},
		{
			name:        "空のファイル",
			lines:       []string{},
			threshold:   1,
			expectError: false,
			contains:    []string{"\"is_greater\": false", "\"line_count\": 0", "\"threshold\": 1", "以下です"},
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用の一時ファイルを作成
			filePath := createTempFileWithLines(t, tc.lines)

			// テスト対象の関数を実行
			result, err := service.HandleToEvaluateLineCount(filePath, tc.threshold)

			// 結果の検証
			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")

				// JSONレスポンスの内容を確認
				for _, expectedContent := range tc.contains {
					assert.Contains(t, result, expectedContent, "JSONレスポンスに期待される内容が含まれていません")
				}
			}
		})
	}

	// 存在しないファイルのテスト
	t.Run("存在しないファイル", func(t *testing.T) {
		_, err := service.HandleToEvaluateLineCount("/path/to/nonexistent/file.txt", 5)
		assert.Error(t, err, "存在しないファイルに対してエラーが発生すべきです")
	})
}

// TestFileEvaluatorServiceHandleToEvaluateLineCountWithMockError は JSONMarshaler のエラーをテストします
func TestFileEvaluatorServiceHandleToEvaluateLineCountWithMockError(t *testing.T) {
	// テスト用の一時ファイルを作成
	filePath := createTempFileWithLines(t, []string{"1行目", "2行目"})

	// モック用のJSONMarshalerを作成（エラーを返すように設定）
	mockJSONMarshaler := &MockJSONMarshaler{
		errToReturn: errors.New("JSON生成エラー"),
	}

	// テスト用のFileEvaluatorServiceを作成し、モックを注入
	service := NewFileEvaluatorServiceWithDependencies(
		&DefaultFileOpener{},
		&MockBufioScanner{},
		mockJSONMarshaler,
	)

	// HandleToEvaluateLineCountメソッドを呼び出す
	_, err := service.HandleToEvaluateLineCount(filePath, 5)

	// 結果の検証
	assert.Error(t, err, "JSON生成エラーが発生すべきです")
	assert.Equal(t, "JSON生成エラー", err.Error(), "エラーメッセージが期待値と一致しません")
}

// TestFileEvaluatorServiceFileOpenerError は FileOpener のエラーをテストします
func TestFileEvaluatorServiceFileOpenerError(t *testing.T) {
	// モック用のFileOpenerを作成（エラーを返すように設定）
	mockFileOpener := &MockFileOpener{
		errToReturn: errors.New("ファイルオープンエラー"),
	}

	// テスト用のFileEvaluatorServiceを作成し、モックを注入
	service := NewFileEvaluatorServiceWithDependencies(
		mockFileOpener,
		&MockBufioScanner{},
		&DefaultJSONMarshaler{},
	)

	// CountLinesメソッドを呼び出す
	_, err := service.CountLines("模擬的なファイルパス")

	// 結果の検証
	assert.Error(t, err, "ファイルオープンエラーが発生すべきです")
	assert.Contains(t, err.Error(), "ファイルを開けませんでした", "エラーメッセージが期待される内容を含んでいません")
}

// #==============================================================#
// ##          Helper Function Tests                            ##
// #==============================================================#

// TestIsGreaterDescription は isGreaterDescription 関数をテストします
func TestIsGreaterDescription(t *testing.T) {
	// テストケースを定義
	testCases := []struct {
		name      string
		isGreater bool
		expected  string
	}{
		{
			name:      "trueの場合",
			isGreater: true,
			expected:  "より大きいです。",
		},
		{
			name:      "falseの場合",
			isGreater: false,
			expected:  "以下です。",
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isGreaterDescription(tc.isGreater)
			assert.Equal(t, tc.expected, result, "説明文が期待値と一致しません")
		})
	}
}

// #==============================================================#
// ##          Constructor Tests                                ##
// #==============================================================#

// TestNewCalculatorService は NewCalculatorService 関数をテストします
func TestNewCalculatorService(t *testing.T) {
	service := NewCalculatorService()
	assert.NotNil(t, service, "CalculatorServiceが正しく作成されませんでした")
}

// TestNewFileEvaluatorService は NewFileEvaluatorService 関数をテストします
func TestNewFileEvaluatorService(t *testing.T) {
	service := NewFileEvaluatorService()
	assert.NotNil(t, service, "FileEvaluatorServiceが正しく作成されませんでした")
}

// TestNewFileEvaluatorServiceWithDependencies は NewFileEvaluatorServiceWithDependencies 関数をテストします
func TestNewFileEvaluatorServiceWithDependencies(t *testing.T) {
	mockFileOpener := &MockFileOpener{}
	mockScanner := &MockBufioScanner{}
	mockMarshaler := &MockJSONMarshaler{}

	service := NewFileEvaluatorServiceWithDependencies(mockFileOpener, mockScanner, mockMarshaler)
	assert.NotNil(t, service, "FileEvaluatorServiceが正しく作成されませんでした")
}
