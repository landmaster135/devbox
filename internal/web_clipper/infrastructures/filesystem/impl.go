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

func (r *osRepository) ListDirectory(path string) ([]FileInfo, error) {
	cleanPath, err := sanitizePath(path)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリ情報の取得に失敗しました (%s): %w", cleanPath, err)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("指定されたパスはディレクトリではありません: %s", cleanPath)
	}

	entries, err := os.ReadDir(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み込みに失敗しました (%s): %w", cleanPath, err)
	}

	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("ファイル情報の取得に失敗しました (%s): %w", entry.Name(), err)
		}

		files = append(files, FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(cleanPath, entry.Name()),
			ModTime: info.ModTime(),
			IsDir:   entry.IsDir(),
		})
	}

	return files, nil
}

func (r *osRepository) Exists(path string) (bool, error) {
	cleanPath, err := sanitizePath(path)
	if err != nil {
		return false, err
	}

	_, err = os.Stat(cleanPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("パスの存在確認に失敗しました (%s): %w", cleanPath, err)
}

func (r *osRepository) Rename(oldPath, newPath string) error {
	cleanOldPath, err := sanitizePath(oldPath)
	if err != nil {
		return err
	}
	cleanNewPath, err := sanitizePath(newPath)
	if err != nil {
		return err
	}

	if err := os.Rename(cleanOldPath, cleanNewPath); err != nil {
		return fmt.Errorf("ファイルのリネームに失敗しました (%s -> %s): %w", cleanOldPath, cleanNewPath, err)
	}

	return nil
}

func sanitizePath(path string) (string, error) {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return "", fmt.Errorf("ファイルパスが空です")
	}

	cleanPath := filepath.Clean(trimmedPath)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("ファイルパスの絶対パス変換に失敗しました (%s): %w", cleanPath, err)
	}

	return absPath, nil
}
