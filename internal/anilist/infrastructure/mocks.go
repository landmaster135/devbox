package infrastructure

import (
	"os"
)

// MockFileSystem はテスト用のファイルシステムモック
type MockFileSystem struct {
	MkdirAllFunc  func(path string, perm os.FileMode) error
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
	JoinFunc      func(elem ...string) string
}

// MkdirAll はディレクトリを作成する（モック）
func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

// WriteFile はファイルに書き込む（モック）
func (m *MockFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(filename, data, perm)
	}
	return nil
}

// Join はパスを結合する（モック）
func (m *MockFileSystem) Join(elem ...string) string {
	if m.JoinFunc != nil {
		return m.JoinFunc(elem...)
	}
	// デフォルトの実装
	result := ""
	for i, e := range elem {
		if i > 0 {
			result += "/"
		}
		result += e
	}
	return result
}

// MockJSONProcessor はテスト用のJSON処理モック
type MockJSONProcessor struct {
	MarshalIndentFunc func(v interface{}, prefix, indent string) ([]byte, error)
}

// MarshalIndent はJSONにマーシャルする（モック）
func (m *MockJSONProcessor) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	if m.MarshalIndentFunc != nil {
		return m.MarshalIndentFunc(v, prefix, indent)
	}
	return []byte("{}"), nil
}
