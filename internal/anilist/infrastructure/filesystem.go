package infrastructure

import (
	"os"
	"path/filepath"
)

// FileSystem はファイルシステム操作のインターフェース
type FileSystem interface {
	MkdirAll(path string, perm os.FileMode) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
	Join(elem ...string) string
}

// OSFileSystem は実際のファイルシステム操作を行う実装
type OSFileSystem struct{}

// NewOSFileSystem は新しいOSFileSystemを作成する
func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

// MkdirAll はディレクトリを作成する
func (fs *OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// WriteFile はファイルに書き込む
func (fs *OSFileSystem) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

// Join はパスを結合する
func (fs *OSFileSystem) Join(elem ...string) string {
	return filepath.Join(elem...)
}
