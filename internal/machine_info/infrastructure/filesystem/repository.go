package filesystem

import "os"

// Repository はファイルシステム操作の抽象化を提供する
// Clean Architecture の層境界を守るため、ユースケース層はこのインターフェース経由で
// ファイル作成やディレクトリ作成を行う
type Repository interface {
	EnsureDir(path string, perm os.FileMode) error
	WriteFile(path string, data []byte, perm os.FileMode) error
	Join(elem ...string) string
}

// NewRepository はOSファイルシステムに基づく標準実装を返す
func NewRepository() Repository {
	return &osRepository{}
}
