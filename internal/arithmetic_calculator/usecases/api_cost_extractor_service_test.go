package usecases

import (
	"errors"
	"os"
	"testing"

	config "github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          ApiCostExtractorService Tests                     ##
// #==============================================================#

// TestApiCostExtractorServiceExtractApiCostFromText は ApiCostExtractorService の extractApiCostFromText メソッドをテストします
func TestApiCostExtractorServiceExtractApiCostFromText(t *testing.T) {
	// テスト用の ApiCostExtractorService インスタンスを作成
	service := NewApiCostExtractorService()

	// テストケースを定義
	testCases := []struct {
		name     string
		text     string
		expected float64
	}{
		{
			name:     "単一のAPI料金パターン",
			text:     "今日はAPI料金が100円掛かった処理を実行しました。",
			expected: 100.0,
		},
		{
			name:     "複数のAPI料金パターン",
			text:     "API料金が50円掛かった処理とAPI料金が75円掛かった処理を実行しました。",
			expected: 125.0,
		},
		{
			name:     "API料金パターンが存在しない場合",
			text:     "今日は通常の処理を実行しました。",
			expected: 0.0,
		},
		{
			name:     "大きな金額のAPI料金",
			text:     "API料金が9999円掛かった重い処理を実行しました。",
			expected: 9999.0,
		},
		{
			name:     "複数の同じ金額のAPI料金",
			text:     "API料金が200円掛かった処理を3回実行しました。API料金が200円掛かった、API料金が200円掛かった。",
			expected: 600.0,
		},
		{
			name:     "0円のAPI料金",
			text:     "API料金が0円掛かった無料処理を実行しました。",
			expected: 0.0,
		},
		{
			name:     "空の文字列",
			text:     "",
			expected: 0.0,
		},
		{
			name:     "パターンに似ているが異なる文字列",
			text:     "API料金は100円でした。料金が100円掛かりました。",
			expected: 0.0,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.extractApiCostFromText(tc.text)
			assert.NoError(t, err, "エラーが発生すべきではありません")
			assert.Equal(t, tc.expected, result, "抽出された金額が期待値と一致しません")
		})
	}
}

// TestApiCostExtractorServiceHandleApiCostExtraction は ApiCostExtractorService の HandleApiCostExtraction メソッドをテストします
func TestApiCostExtractorServiceHandleApiCostExtraction(t *testing.T) {
	// テスト用の ApiCostExtractorService インスタンスを作成
	service := NewApiCostExtractorService()

	// テストケースを定義
	testCases := []struct {
		name        string
		filePath    string
		textInput   string
		fileContent string
		expected    float64
		expectError bool
		errorMsg    string
	}{
		{
			name:        "テキスト入力からの抽出",
			filePath:    "",
			textInput:   "API料金が150円掛かった処理を実行しました。",
			expected:    150.0,
			expectError: false,
		},
		{
			name:        "複数のAPI料金を含むテキスト",
			filePath:    "",
			textInput:   "API料金が100円掛かった処理とAPI料金が200円掛かった処理を実行しました。",
			expected:    300.0,
			expectError: false,
		},
		{
			name:        "API料金が含まれないテキスト",
			filePath:    "",
			textInput:   "通常の処理を実行しました。",
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "ファイルパスとテキスト両方指定エラー",
			filePath:    "test.md",
			textInput:   "API料金が100円掛かった",
			expected:    0.0,
			expectError: true,
			errorMsg:    "ファイルパスとテキスト入力は同時に指定できません",
		},
		{
			name:        "どちらも指定なしエラー",
			filePath:    "",
			textInput:   "",
			expected:    0.0,
			expectError: true,
			errorMsg:    "ファイルパスまたはテキスト入力のいずれかを指定してください",
		},
		{
			name:        "不正なファイル拡張子エラー",
			filePath:    "test.json",
			textInput:   "",
			expected:    0.0,
			expectError: true,
			errorMsg:    "ファイルは.mdまたは.txt形式である必要があります",
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.HandleApiCostExtraction(tc.filePath, tc.textInput)

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				assert.Equal(t, tc.errorMsg, err.Error(), "エラーメッセージが期待値と一致しません")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.Equal(t, tc.expected, result, "抽出された金額が期待値と一致しません")
			}
		})
	}
}

// TestApiCostExtractorServiceHandleApiCostExtractionWithFile は ファイルからのAPI料金抽出をテストします
func TestApiCostExtractorServiceHandleApiCostExtractionWithFile(t *testing.T) {
	// テスト用の ApiCostExtractorService インスタンスを作成
	service := NewApiCostExtractorService()

	// テストケースを定義
	testCases := []struct {
		name        string
		fileContent string
		extension   string
		expected    float64
		expectError bool
	}{
		{
			name:        "MDファイルからの抽出",
			fileContent: "# レポート\nAPI料金が250円掛かった処理を実行しました。\nAPI料金が150円掛かった別の処理も実行しました。",
			extension:   ".md",
			expected:    400.0,
			expectError: false,
		},
		{
			name:        "TXTファイルからの抽出",
			fileContent: "ログファイル\nAPI料金が500円掛かった重い処理を実行しました。",
			extension:   ".txt",
			expected:    500.0,
			expectError: false,
		},
		{
			name:        "API料金が含まれないファイル",
			fileContent: "通常のログファイルです。\n特に料金は発生していません。",
			extension:   ".md",
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "空のファイル",
			fileContent: "",
			extension:   ".txt",
			expected:    0.0,
			expectError: false,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用の一時ファイルを作成
			tmpFile, err := os.CreateTemp("", "test-api-cost-*"+tc.extension)
			assert.NoError(t, err, "一時ファイルの作成に失敗しました")

			// テスト終了時に一時ファイルを削除
			t.Cleanup(func() {
				os.Remove(tmpFile.Name())
			})

			// ファイルに内容を書き込む
			_, err = tmpFile.WriteString(tc.fileContent)
			assert.NoError(t, err, "ファイルへの書き込みに失敗しました")

			// ファイルを閉じる
			err = tmpFile.Close()
			assert.NoError(t, err, "ファイルのクローズに失敗しました")

			// テスト対象の関数を実行
			result, err := service.HandleApiCostExtraction(tmpFile.Name(), "")

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.Equal(t, tc.expected, result, "抽出された金額が期待値と一致しません")
			}
		})
	}

	// 存在しないファイルのテスト
	t.Run("存在しないファイル", func(t *testing.T) {
		_, err := service.HandleApiCostExtraction("/path/to/nonexistent/file.md", "")
		assert.Error(t, err, "存在しないファイルに対してエラーが発生すべきです")
		assert.Contains(t, err.Error(), "ファイル読み込みエラー", "エラーメッセージが期待される内容を含んでいません")
	})
}

// TestApiCostExtractorServiceWithMockFileReader は モックFileReaderを使用したテストです
func TestApiCostExtractorServiceWithMockFileReader(t *testing.T) {
	// テストケースを定義
	testCases := []struct {
		name          string
		mockData      []byte
		mockError     error
		filePath      string
		expected      float64
		expectError   bool
		errorContains string
	}{
		{
			name:        "モックFileReaderからの正常な読み込み",
			mockData:    []byte("API料金が300円掛かった処理を実行しました。"),
			mockError:   nil,
			filePath:    "test.md",
			expected:    300.0,
			expectError: false,
		},
		{
			name:          "モックFileReaderからのエラー",
			mockData:      nil,
			mockError:     errors.New("ファイル読み込み失敗"),
			filePath:      "test.txt",
			expected:      0.0,
			expectError:   true,
			errorContains: "ファイル読み込みエラー",
		},
		{
			name:        "モックFileReaderから空データ",
			mockData:    []byte(""),
			mockError:   nil,
			filePath:    "empty.md",
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "モックFileReaderから複数のAPI料金データ",
			mockData:    []byte("API料金が100円掛かった処理とAPI料金が200円掛かった処理とAPI料金が300円掛かった処理を実行しました。"),
			mockError:   nil,
			filePath:    "multiple.txt",
			expected:    600.0,
			expectError: false,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// モックFileReaderを作成
			mockFileReader := &MockFileReader{
				dataToReturn: tc.mockData,
				errToReturn:  tc.mockError,
			}

			// テスト用のApiCostExtractorServiceを作成し、モックを注入
			service := NewApiCostExtractorServiceWithFileReader(mockFileReader)

			// テスト対象の関数を実行
			result, err := service.HandleApiCostExtraction(tc.filePath, "")

			if tc.expectError {
				assert.Error(t, err, "エラーが発生すべきです")
				if tc.errorContains != "" {
					assert.Contains(t, err.Error(), tc.errorContains, "エラーメッセージが期待される内容を含んでいません")
				}
			} else {
				assert.NoError(t, err, "エラーが発生すべきではありません")
				assert.Equal(t, tc.expected, result, "抽出された金額が期待値と一致しません")
			}
		})
	}
}

// #==============================================================#
// ##          Constructor Tests                                ##
// #==============================================================#

// TestNewApiCostExtractorService は NewApiCostExtractorService 関数をテストします
func TestNewApiCostExtractorService(t *testing.T) {
	service := NewApiCostExtractorService()
	assert.NotNil(t, service, "ApiCostExtractorServiceが正しく作成されませんでした")
}

// TestNewApiCostExtractorServiceWithFileReader は NewApiCostExtractorServiceWithFileReader 関数をテストします
func TestNewApiCostExtractorServiceWithFileReader(t *testing.T) {
	mockFileReader := &config.StandardFileReader{}
	service := NewApiCostExtractorServiceWithFileReader(mockFileReader)
	assert.NotNil(t, service, "ApiCostExtractorServiceが正しく作成されませんでした")
}

// TestNewApiCostExtractorServiceWithMockFileReader は モックFileReaderを使用したコンストラクタをテストします
func TestNewApiCostExtractorServiceWithMockFileReader(t *testing.T) {
	mockFileReader := &MockFileReader{}
	service := NewApiCostExtractorServiceWithFileReader(mockFileReader)
	assert.NotNil(t, service, "ApiCostExtractorServiceが正しく作成されませんでした")
}

// #==============================================================#
// ##          Edge Case Tests                                  ##
// #==============================================================#

// TestApiCostExtractorServiceEdgeCases は ApiCostExtractorService のエッジケースをテストします
func TestApiCostExtractorServiceEdgeCases(t *testing.T) {
	// テスト用の ApiCostExtractorService インスタンスを作成
	service := NewApiCostExtractorService()

	// エッジケースのテストケースを定義
	testCases := []struct {
		name     string
		text     string
		expected float64
	}{
		{
			name:     "非常に大きな金額",
			text:     "API料金が999999999円掛かった処理を実行しました。",
			expected: 999999999.0,
		},
		{
			name:     "改行を含むテキスト",
			text:     "処理1:\nAPI料金が100円掛かった\n処理2:\nAPI料金が200円掛かった",
			expected: 300.0,
		},
		{
			name:     "特殊文字を含むテキスト",
			text:     "【重要】API料金が500円掛かった処理を実行しました！",
			expected: 500.0,
		},
		{
			name:     "数値が文字列の最後にある場合",
			text:     "今回の処理でAPI料金が777円掛かった",
			expected: 777.0,
		},
		{
			name:     "複数の数値パターンが混在",
			text:     "API料金が123円掛かった処理とAPI料金が456円掛かった処理とAPI料金が789円掛かった処理",
			expected: 1368.0,
		},
	}

	// 各テストケースを実行
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.extractApiCostFromText(tc.text)
			assert.NoError(t, err, "エラーが発生すべきではありません")
			assert.Equal(t, tc.expected, result, "抽出された金額が期待値と一致しません")
		})
	}
}
