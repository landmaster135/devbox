package filesystem

import "os"

// MockRepository はRepositoryのテスト用実装。
type MockRepository struct {
	ReadFileFunc  func(path string) ([]byte, error)
	WriteFileFunc func(path string, data []byte, perm os.FileMode) error

	LastReadPath string

	LastWritePath string
	LastWriteData []byte
	LastWritePerm os.FileMode

	WriteCallCount int
}

func (m *MockRepository) ReadFile(path string) ([]byte, error) {
	m.LastReadPath = path
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(path)
	}
	return nil, nil
}

func (m *MockRepository) WriteFile(path string, data []byte, perm os.FileMode) error {
	m.LastWritePath = path
	m.LastWriteData = append([]byte(nil), data...)
	m.LastWritePerm = perm
	m.WriteCallCount++
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(path, data, perm)
	}
	return nil
}
