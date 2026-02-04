package writer

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
	WriteFunc     func(file *os.File, content []byte) (int, error)
	CloseFunc     func(file *os.File) error
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

func (m *MockFileWriter) Write(file *os.File, content []byte) (int, error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(file, content)
	}
	if file == nil {
		return 0, nil
	}
	return file.Write(content)
}

func (m *MockFileWriter) Close(file *os.File) error {
	if m.CloseFunc != nil {
		return m.CloseFunc(file)
	}
	if file == nil {
		return nil
	}
	return file.Close()
}
