package filesystem

import "os"

// Repository はcoverage-badgeで必要なファイルI/Oを抽象化する。
type Repository interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
}

// NewRepository はOSファイルシステム実装を返す。
func NewRepository() Repository {
	return &osRepository{}
}
