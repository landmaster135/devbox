package filesystem

type Repository interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
}

func NewRepository() Repository {
	return &osRepository{}
}
