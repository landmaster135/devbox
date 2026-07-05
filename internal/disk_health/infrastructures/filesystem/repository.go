package filesystem

type Repository interface {
	ReadFile(filePath string) ([]byte, error)
}
