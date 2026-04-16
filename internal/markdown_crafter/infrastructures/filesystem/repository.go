package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) ReadFile(filePath string) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("ファイルの読み込みに失敗しました: %s: %w", filePath, err)
	}
	return string(b), nil
}

func (r *Repository) WriteFile(filePath string, content string) error {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %s: %w", filePath, err)
	}
	return nil
}

func (r *Repository) CreateDir(dirPath string) error {
	cleaned := filepath.Clean(dirPath)
	if err := os.MkdirAll(cleaned, 0755); err != nil {
		return fmt.Errorf("ディレクトリ作成に失敗しました: %s: %w", cleaned, err)
	}
	return nil
}

func (r *Repository) ListMarkdownFiles(dirPath string) ([]string, error) {
	cleaned := filepath.Clean(dirPath)
	entries, err := os.ReadDir(cleaned)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み込みに失敗しました: %s: %w", cleaned, err)
	}

	markdownFiles := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			markdownFiles = append(markdownFiles, filepath.Join(cleaned, entry.Name()))
		}
	}

	sort.Strings(markdownFiles)
	return markdownFiles, nil
}

func (r *Repository) RemoveFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("ファイルの削除に失敗しました: %s: %w", filePath, err)
	}
	return nil
}
