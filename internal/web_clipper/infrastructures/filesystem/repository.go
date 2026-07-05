package filesystem

import "time"

type FileInfo struct {
	Name    string
	Path    string
	ModTime time.Time
	IsDir   bool
}

type Repository interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	ListDirectory(path string) ([]FileInfo, error)
	Exists(path string) (bool, error)
	Rename(oldPath, newPath string) error
}

func NewRepository() Repository {
	return &osRepository{}
}
