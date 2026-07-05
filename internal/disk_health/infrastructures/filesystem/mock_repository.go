package filesystem

type MockRepository struct {
	ReadFileFunc func(filePath string) ([]byte, error)
}

func (r *MockRepository) ReadFile(filePath string) ([]byte, error) {
	return r.ReadFileFunc(filePath)
}
