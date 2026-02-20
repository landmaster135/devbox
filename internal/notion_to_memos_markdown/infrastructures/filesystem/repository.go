package filesystem

import "os"

type Repository interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	MkdirAll(path string, perm os.FileMode) error
	FileExists(path string) (bool, error)
	CopyFile(srcPath, dstPath string) error
	RenameFile(srcPath, dstPath string) error
	ListMarkdownFiles(dirPath string) ([]string, error)
	ListFilesRecursive(dirPath string) ([]string, error)
}

func NewRepository() Repository {
	return &osRepository{}
}
