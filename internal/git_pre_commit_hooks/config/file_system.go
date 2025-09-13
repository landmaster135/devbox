package config

import (
	"os"
)

// FileSystemRepository はファイルシステム操作のインターフェース
type FileSystemRepository interface {
	Stat(filename string) (os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

// FileSystem は実際のファイルシステム操作を行う実装
type FileSystem struct{}

// NewFileSystem は新しいFileSystemを作成
func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

// Stat はファイル情報を取得
func (f *FileSystem) Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

// ReadFile はファイルを読み込み
func (f *FileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// MockFileSystem はテスト用のモック実装
type MockFileSystem struct {
	StatFunc     func(filename string) (os.FileInfo, error)
	ReadFileFunc func(filename string) ([]byte, error)
}

// Stat はモック関数を実行
func (m *MockFileSystem) Stat(filename string) (os.FileInfo, error) {
	return m.StatFunc(filename)
}

// ReadFile はモック関数を実行
func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	return m.ReadFileFunc(filename)
}
