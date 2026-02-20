package grepstr

import (
	"fmt"
	"strings"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

type Service struct {
	fileSystem filesystem.Repository
}

func NewService(fileSystem filesystem.Repository) *Service {
	return &Service{
		fileSystem: fileSystem,
	}
}

func (s *Service) Execute(srcBodyDir, targetStr string) (string, error) {
	trimmedSrcBodyDir := strings.TrimSpace(srcBodyDir)
	if trimmedSrcBodyDir == "" {
		return "", fmt.Errorf("src-body-dir パラメータは必須です")
	}

	trimmedTargetStr := strings.TrimSpace(targetStr)
	if trimmedTargetStr == "" {
		return "", fmt.Errorf("target-str パラメータは必須です")
	}

	filePaths, err := s.fileSystem.ListFilesRecursive(trimmedSrcBodyDir)
	if err != nil {
		return "", fmt.Errorf("src-body-dir の読み取りに失敗しました: %w", err)
	}

	matchedFiles := make([]string, 0)
	for _, filePath := range filePaths {
		data, err := s.fileSystem.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("ファイルの読み込みに失敗しました (%s): %w", filePath, err)
		}
		if strings.Contains(string(data), trimmedTargetStr) {
			matchedFiles = append(matchedFiles, filePath)
		}
	}

	var output strings.Builder
	output.WriteString("処理完了\n")
	output.WriteString(fmt.Sprintf("対象ファイル総数=%d\n", len(filePaths)))
	output.WriteString(fmt.Sprintf("該当ファイル総数=%d\n", len(matchedFiles)))
	output.WriteString("該当ファイル一覧:\n")
	if len(matchedFiles) == 0 {
		output.WriteString("(なし)")
		return output.String(), nil
	}

	for i, filePath := range matchedFiles {
		if i > 0 {
			output.WriteString("\n")
		}
		output.WriteString(filePath)
	}

	return output.String(), nil
}
