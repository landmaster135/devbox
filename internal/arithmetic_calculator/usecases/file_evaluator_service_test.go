package usecases

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
// ##          Constructor Tests                                ##
// #==============================================================#

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
