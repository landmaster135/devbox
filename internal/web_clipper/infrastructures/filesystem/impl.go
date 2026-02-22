package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type osRepository struct{}

func (r *osRepository) ReadFile(path string) ([]byte, error) {
	cleanPath, err := sanitizePath(path)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました (%s): %w", cleanPath, err)
	}

	return data, nil
}

func (r *osRepository) WriteFile(path string, data []byte) error {
	cleanPath, err := sanitizePath(path)
	if err != nil {
		return err
	}

	dirPath := filepath.Dir(cleanPath)
	if dirPath != "." {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("出力先ディレクトリの作成に失敗しました (%s): %w", dirPath, err)
		}
	}

	if err := os.WriteFile(cleanPath, data, 0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (%s): %w", cleanPath, err)
	}

	return nil
}

func sanitizePath(path string) (string, error) {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return "", fmt.Errorf("ファイルパスが空です")
	}

	return filepath.Clean(trimmedPath), nil
}
