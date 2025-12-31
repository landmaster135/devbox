package writers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// FileWriter はDOMをファイルへ保存する実装です。
type FileWriter struct{}

// NewFileWriter はFileWriterを生成します。
func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

// Write は指定されたパスに新規ファイルとしてcontentを書き込みます。
func (w *FileWriter) Write(path string, content string) error {
	cleanPath := filepath.Clean(path)
	if cleanPath == "" {
		return errors.New("出力パスが空です")
	}

	dir := filepath.Dir(cleanPath)
	if dir != "." && dir != "" {
		info, err := os.Stat(dir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("出力先ディレクトリが存在しません: %s", dir)
			}
			return fmt.Errorf("出力先ディレクトリの確認に失敗しました: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("出力先がディレクトリではありません: %s", dir)
		}
	}

	if _, err := os.Stat(cleanPath); err == nil {
		return fmt.Errorf("出力ファイルが既に存在します: %s", cleanPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("出力ファイルの確認に失敗しました: %w", err)
	}

	if err := os.WriteFile(cleanPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("出力ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}
