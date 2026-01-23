package dump

import (
	"os"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockFileWriter はテスト用のFileWriterモック
type MockFileWriter struct {
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
	MkdirAllFunc  func(path string, perm os.FileMode) error
	CreateFunc    func(name string) (*os.File, error)
}

func (m *MockFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(filename, data, perm)
	}
	return nil
}

func (m *MockFileWriter) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

func (m *MockFileWriter) Create(name string) (*os.File, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(name)
	}
	return nil, nil
}
