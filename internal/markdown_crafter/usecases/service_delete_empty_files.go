package usecases

import (
	"fmt"
	"strings"
)

const miscellaneousNotesTemplate = "# Miscellaneous notes\n- "

func (s *Service) DeleteEmptyFiles(directoryPath string) (string, error) {
	markdownFiles, err := s.repository.ListMarkdownFiles(directoryPath)
	if err != nil {
		return "", err
	}

	deletedFiles := make([]string, 0)
	for _, markdownFile := range markdownFiles {
		content, err := s.repository.ReadFile(markdownFile)
		if err != nil {
			return "", err
		}

		if !isDeleteTarget(content) {
			continue
		}

		if err := s.repository.RemoveFile(markdownFile); err != nil {
			return "", err
		}
		deletedFiles = append(deletedFiles, markdownFile)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("delete-empty-files: %d ファイルを削除しました\n", len(deletedFiles)))
	for _, deletedFile := range deletedFiles {
		builder.WriteString(fmt.Sprintf("- %s\n", deletedFile))
	}

	return builder.String(), nil
}

func isDeleteTarget(content string) bool {
	normalized := normalizeNewlines(content)
	if normalized == "" {
		return true
	}

	return normalized == miscellaneousNotesTemplate || normalized == miscellaneousNotesTemplate+"\n"
}
