package common

import (
	"fmt"
	"os"

	markdowndomain "github.com/landmaster135/devbox/internal/markdown_crafter/domain"
	markdownusecases "github.com/landmaster135/devbox/internal/markdown_crafter/usecases"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

type markdownCrafterRepositoryAdapter struct {
	fileSystem filesystem.Repository
}

func NewMarkdownService(fileSystem filesystem.Repository) *markdownusecases.Service {
	return markdownusecases.NewService(newMarkdownCrafterRepositoryAdapter(fileSystem))
}

func newMarkdownCrafterRepositoryAdapter(fileSystem filesystem.Repository) markdowndomain.Repository {
	return &markdownCrafterRepositoryAdapter{
		fileSystem: fileSystem,
	}
}

func (a *markdownCrafterRepositoryAdapter) ReadFile(filePath string) (string, error) {
	data, err := a.fileSystem.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *markdownCrafterRepositoryAdapter) WriteFile(filePath string, content string) error {
	return a.fileSystem.WriteFile(filePath, []byte(content))
}

func (a *markdownCrafterRepositoryAdapter) CreateDir(dirPath string) error {
	return a.fileSystem.MkdirAll(dirPath, DefaultDirectoryPerm)
}

func (a *markdownCrafterRepositoryAdapter) ListMarkdownFiles(dirPath string) ([]string, error) {
	return a.fileSystem.ListMarkdownFiles(dirPath)
}

func (a *markdownCrafterRepositoryAdapter) RemoveFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("ファイルの削除に失敗しました (%s): %w", filePath, err)
	}
	return nil
}
