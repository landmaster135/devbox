package writer

import "os"

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#

// FileWriter はファイル書き込み操作のインターフェースです
type FileWriter interface {
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (*os.File, error)
	Write(file *os.File, content []byte) (int, error)
	Close(file *os.File) error
}


type TableDataWriter interface {
	WriteBatch(rows []map[string]any) error
	Close() error
	RowsWritten() int
}

// #==============================================================#
// ##          Default Implementations                           ##
// #==============================================================#

// DefaultFileWriter は標準のファイル操作を使用する実装
type DefaultFileWriter struct{}

func (w *DefaultFileWriter) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (w *DefaultFileWriter) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (w *DefaultFileWriter) Write(file *os.File, content []byte) (int, error) {
	return file.Write(content)
}

func (w *DefaultFileWriter) Close(file *os.File) error {
	return file.Close()
}
