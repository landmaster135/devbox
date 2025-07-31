package usecases

import (
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

// モック用のFileReaderインターフェース実装
type MockFileReader struct {
	dataToReturn []byte
	errToReturn  error
}

func (m *MockFileReader) ReadFile(filename string) ([]byte, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	return m.dataToReturn, nil
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
