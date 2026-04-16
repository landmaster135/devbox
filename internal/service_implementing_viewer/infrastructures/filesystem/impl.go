package filesystem

import (
	"errors"
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

func (r *osRepository) ListDirectories(path string) ([]string, error) {
	cleanPath, err := sanitizePath(path)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(cleanPath)
	if err != nil {
		// 既存仕様に合わせ、存在しないディレクトリは空配列を返す。
		if errors.Is(err, os.ErrNotExist) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("ディレクトリの読み取りに失敗しました (%s): %w", cleanPath, err)
	}

	directories := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			directories = append(directories, entry.Name())
		}
	}
	return directories, nil
}

func (r *osRepository) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func sanitizePath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("ファイルパスが空です")
	}
	return filepath.Clean(trimmed), nil
}
