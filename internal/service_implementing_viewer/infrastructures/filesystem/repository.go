package filesystem

import "os"

// Repository は service_implementing_viewer で利用するファイルシステム操作を抽象化する。
type Repository interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	ListDirectories(path string) ([]string, error)
	Join(elem ...string) string
}

// NewRepository はOSファイルシステム実装を返す。
func NewRepository() Repository {
	return &osRepository{}
}
