package filesystem

import (
	"os"
	"path/filepath"
)

// MockRepository は Repository のテスト用実装。
// usecases テストから利用することを想定している。
type MockRepository struct {
	ReadFileFunc        func(path string) ([]byte, error)
	WriteFileFunc       func(path string, data []byte, perm os.FileMode) error
	ListDirectoriesFunc func(path string) ([]string, error)
	JoinFunc            func(elem ...string) string

	LastWritePath       string
	LastWriteContent    []byte
	LastWritePermission os.FileMode
}

func (m *MockRepository) ReadFile(path string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(path)
	}
	return nil, nil
}

func (m *MockRepository) WriteFile(path string, data []byte, perm os.FileMode) error {
	m.LastWritePath = path
	m.LastWriteContent = append([]byte(nil), data...)
	m.LastWritePermission = perm
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(path, data, perm)
	}
	return nil
}

func (m *MockRepository) ListDirectories(path string) ([]string, error) {
	if m.ListDirectoriesFunc != nil {
		return m.ListDirectoriesFunc(path)
	}
	return []string{}, nil
}

func (m *MockRepository) Join(elem ...string) string {
	if m.JoinFunc != nil {
		return m.JoinFunc(elem...)
	}
	return filepath.Join(elem...)
}
