package filesystem

import "os"

type OSRepository struct{}

func NewOSRepository() *OSRepository {
	return &OSRepository{}
}

func (r *OSRepository) ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}
