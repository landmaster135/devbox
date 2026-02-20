package filesystem

import "os"

type MkdirAllCall struct {
	Path string
	Perm os.FileMode
}

type CopyFileCall struct {
	SrcPath string
	DstPath string
}

type MockRepository struct {
	ReadFileFunc           func(path string) ([]byte, error)
	WriteFileFunc          func(path string, data []byte) error
	MkdirAllFunc           func(path string, perm os.FileMode) error
	FileExistsFunc         func(path string) (bool, error)
	CopyFileFunc           func(srcPath, dstPath string) error
	RenameFileFunc         func(srcPath, dstPath string) error
	ListMarkdownFilesFunc  func(dirPath string) ([]string, error)
	ListFilesRecursiveFunc func(dirPath string) ([]string, error)

	ReadFileCalls           []string
	WriteFileCalls          []WriteFileCall
	MkdirAllCalls           []MkdirAllCall
	FileExistsCalls         []string
	CopyFileCalls           []CopyFileCall
	RenameFileCalls         []CopyFileCall
	ListMarkdownFilesCalls  []string
	ListFilesRecursiveCalls []string
}

type WriteFileCall struct {
	Path string
	Data []byte
}

func (m *MockRepository) ReadFile(path string) ([]byte, error) {
	m.ReadFileCalls = append(m.ReadFileCalls, path)
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(path)
	}
	return nil, nil
}

func (m *MockRepository) WriteFile(path string, data []byte) error {
	cloned := append([]byte(nil), data...)
	m.WriteFileCalls = append(m.WriteFileCalls, WriteFileCall{
		Path: path,
		Data: cloned,
	})
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(path, cloned)
	}
	return nil
}

func (m *MockRepository) MkdirAll(path string, perm os.FileMode) error {
	m.MkdirAllCalls = append(m.MkdirAllCalls, MkdirAllCall{
		Path: path,
		Perm: perm,
	})
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

func (m *MockRepository) FileExists(path string) (bool, error) {
	m.FileExistsCalls = append(m.FileExistsCalls, path)
	if m.FileExistsFunc != nil {
		return m.FileExistsFunc(path)
	}
	return false, nil
}

func (m *MockRepository) CopyFile(srcPath, dstPath string) error {
	m.CopyFileCalls = append(m.CopyFileCalls, CopyFileCall{
		SrcPath: srcPath,
		DstPath: dstPath,
	})
	if m.CopyFileFunc != nil {
		return m.CopyFileFunc(srcPath, dstPath)
	}
	return nil
}

func (m *MockRepository) RenameFile(srcPath, dstPath string) error {
	m.RenameFileCalls = append(m.RenameFileCalls, CopyFileCall{
		SrcPath: srcPath,
		DstPath: dstPath,
	})
	if m.RenameFileFunc != nil {
		return m.RenameFileFunc(srcPath, dstPath)
	}
	return nil
}

func (m *MockRepository) ListMarkdownFiles(dirPath string) ([]string, error) {
	m.ListMarkdownFilesCalls = append(m.ListMarkdownFilesCalls, dirPath)
	if m.ListMarkdownFilesFunc != nil {
		return m.ListMarkdownFilesFunc(dirPath)
	}
	return []string{}, nil
}

func (m *MockRepository) ListFilesRecursive(dirPath string) ([]string, error) {
	m.ListFilesRecursiveCalls = append(m.ListFilesRecursiveCalls, dirPath)
	if m.ListFilesRecursiveFunc != nil {
		return m.ListFilesRecursiveFunc(dirPath)
	}
	return []string{}, nil
}
