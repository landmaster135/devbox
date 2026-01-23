package dump

import "os"

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#

// FileWriter はファイル書き込み操作のインターフェースです
type FileWriter interface {
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (*os.File, error)
}

// #==============================================================#
// ##          Default Implementations                           ##
// #==============================================================#

// DefaultFileWriter は標準のファイル操作を使用する実装
type DefaultFileWriter struct{}

func (w *DefaultFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (w *DefaultFileWriter) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (w *DefaultFileWriter) Create(name string) (*os.File, error) {
	return os.Create(name)
}
