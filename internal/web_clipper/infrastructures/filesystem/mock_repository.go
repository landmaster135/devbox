package filesystem

type WriteFileCall struct {
	Path string
	Data []byte
}

type MockRepository struct {
	ReadFileFunc       func(path string) ([]byte, error)
	WriteFileFunc      func(path string, data []byte) error
	ListDirectoryFunc  func(path string) ([]FileInfo, error)
	ExistsFunc         func(path string) (bool, error)
	RenameFunc         func(oldPath, newPath string) error
	ReadFileCalls      []string
	WriteFileCalls     []WriteFileCall
	ListDirectoryCalls []string
	ExistsCalls        []string
	RenameCalls        []RenameCall
}

type RenameCall struct {
	OldPath string
	NewPath string
}

func (m *MockRepository) ReadFile(path string) ([]byte, error) {
	m.ReadFileCalls = append(m.ReadFileCalls, path)
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(path)
	}

	return nil, nil
}

func (m *MockRepository) WriteFile(path string, data []byte) error {
	clonedData := append([]byte(nil), data...)
	m.WriteFileCalls = append(m.WriteFileCalls, WriteFileCall{
		Path: path,
		Data: clonedData,
	})

	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(path, clonedData)
	}

	return nil
}

func (m *MockRepository) ListDirectory(path string) ([]FileInfo, error) {
	m.ListDirectoryCalls = append(m.ListDirectoryCalls, path)
	if m.ListDirectoryFunc != nil {
		return m.ListDirectoryFunc(path)
	}

	return nil, nil
}

func (m *MockRepository) Exists(path string) (bool, error) {
	m.ExistsCalls = append(m.ExistsCalls, path)
	if m.ExistsFunc != nil {
		return m.ExistsFunc(path)
	}

	return false, nil
}

func (m *MockRepository) Rename(oldPath, newPath string) error {
	m.RenameCalls = append(m.RenameCalls, RenameCall{
		OldPath: oldPath,
		NewPath: newPath,
	})
	if m.RenameFunc != nil {
		return m.RenameFunc(oldPath, newPath)
	}

	return nil
}
