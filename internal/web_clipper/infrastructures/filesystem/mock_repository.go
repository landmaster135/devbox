package filesystem

type WriteFileCall struct {
	Path string
	Data []byte
}

type MockRepository struct {
	ReadFileFunc   func(path string) ([]byte, error)
	WriteFileFunc  func(path string, data []byte) error
	ReadFileCalls  []string
	WriteFileCalls []WriteFileCall
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
