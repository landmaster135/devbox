package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type osRepository struct{}

func (r *osRepository) EnsureDir(path string, perm os.FileMode) error {
	cleanPath := sanitizePath(path)
	if cleanPath == "" {
		cleanPath = "."
	}
	if err := os.MkdirAll(cleanPath, perm); err != nil {
		return fmt.Errorf("ディレクトリ作成に失敗しました (%s): %w", cleanPath, err)
	}
	return nil
}

func (r *osRepository) WriteFile(path string, data []byte, perm os.FileMode) error {
	cleanPath := sanitizePath(path)
	if cleanPath == "" {
		return fmt.Errorf("ファイルパスが空です")
	}
	if err := os.WriteFile(cleanPath, data, perm); err != nil {
		return fmt.Errorf("ファイル書き込みに失敗しました (%s): %w", cleanPath, err)
	}
	return nil
}

func (r *osRepository) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func sanitizePath(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}
	return filepath.Clean(trimmed)
}
