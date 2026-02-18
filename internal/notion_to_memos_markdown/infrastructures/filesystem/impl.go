package filesystem

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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

	if err := os.WriteFile(cleanPath, data, 0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (%s): %w", cleanPath, err)
	}
	return nil
}

func (r *osRepository) MkdirAll(path string, perm os.FileMode) error {
	cleanPath, err := sanitizePath(path)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(cleanPath, perm); err != nil {
		return fmt.Errorf("ディレクトリの作成に失敗しました (%s): %w", cleanPath, err)
	}
	return nil
}

func (r *osRepository) FileExists(path string) (bool, error) {
	cleanPath, err := sanitizePath(path)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(cleanPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, fmt.Errorf("ファイル状態の確認に失敗しました (%s): %w", cleanPath, err)
	}

	return !info.IsDir(), nil
}

func (r *osRepository) CopyFile(srcPath, dstPath string) error {
	cleanSrcPath, err := sanitizePath(srcPath)
	if err != nil {
		return err
	}
	cleanDstPath, err := sanitizePath(dstPath)
	if err != nil {
		return err
	}

	srcFile, err := os.Open(cleanSrcPath)
	if err != nil {
		return fmt.Errorf("コピー元ファイルのオープンに失敗しました (%s): %w", cleanSrcPath, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(cleanDstPath)
	if err != nil {
		return fmt.Errorf("コピー先ファイルの作成に失敗しました (%s): %w", cleanDstPath, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("ファイルコピーに失敗しました (%s -> %s): %w", cleanSrcPath, cleanDstPath, err)
	}
	return nil
}

func (r *osRepository) ListMarkdownFiles(dirPath string) ([]string, error) {
	cleanDirPath, err := sanitizePath(dirPath)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(cleanDirPath)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み取りに失敗しました (%s): %w", cleanDirPath, err)
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			files = append(files, filepath.Join(cleanDirPath, entry.Name()))
		}
	}
	sort.Strings(files)
	return files, nil
}

func sanitizePath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("ファイルパスが空です")
	}
	return filepath.Clean(trimmed), nil
}
