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

func (r *osRepository) WriteFile(path string, data []byte, perm os.FileMode) error {
	cleanPath, err := sanitizePath(path)
	if err != nil {
		return err
	}

	if err := os.WriteFile(cleanPath, data, perm); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (%s): %w", cleanPath, err)
	}

	return nil
}

func sanitizePath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("ファイルパスが空です")
	}
	return filepath.Clean(trimmed), nil
}
