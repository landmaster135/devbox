package usecases

import (
	"fmt"
	"os"
	"path/filepath"
)

// #==============================================================#
// ##          FileWriter                                        ##
// #==============================================================#
// FileWriter はファイル書き込み操作のインターフェース
type FileWriter interface {
	WriteToFile(filePath, content string) error
	EnsureDirectory(dirPath string) error
}

// FileWriterImpl はFileWriterの実装
type FileWriterImpl struct{}

// NewFileWriter は新しいFileWriterインスタンスを作成する
func NewFileWriter() FileWriter {
	return &FileWriterImpl{}
}

// WriteToFile はファイルに内容を書き込む
func (fw *FileWriterImpl) WriteToFile(filePath, content string) error {
	// ディレクトリを確保
	dir := filepath.Dir(filePath)
	if err := fw.EnsureDirectory(dir); err != nil {
		return err
	}

	// ファイルに書き込み
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %v", err)
	}

	return nil
}

// EnsureDirectory はディレクトリが存在することを確認し、必要に応じて作成する
func (fw *FileWriterImpl) EnsureDirectory(dirPath string) error {
	if dirPath == "." || dirPath == "" {
		return nil
	}

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("ディレクトリの作成に失敗しました: %v", err)
	}

	return nil
}
