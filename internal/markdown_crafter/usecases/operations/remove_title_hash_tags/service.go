package removetitlehashtags

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/landmaster135/devbox/internal/markdown_crafter/domain"
	"github.com/landmaster135/devbox/internal/markdown_crafter/usecases/common"
)

var (
	lineStartHashTagPattern = regexp.MustCompile(`^#[\p{L}\p{N}_.+\-]+\s*`)
	inlineHashTagPattern    = regexp.MustCompile(`[ \t　]+#[\p{L}\p{N}_.+\-]+`)
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Execute(dirPath string) (string, error) {
	markdownFiles, err := s.repository.ListMarkdownFiles(dirPath)
	if err != nil {
		return "", err
	}
	if len(markdownFiles) == 0 {
		return fmt.Sprintf("remove-title-hash-tags: %s 内のMarkdownファイルは0件でした\n", dirPath), nil
	}

	updatedFiles := make([]string, 0, len(markdownFiles))
	for _, filePath := range markdownFiles {
		content, err := s.repository.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("%s の読み込みに失敗しました: %w", filePath, err)
		}

		updated := removeTitleHashTagsFromContent(content)
		if content == updated {
			continue
		}

		if err := s.repository.WriteFile(filePath, updated); err != nil {
			return "", fmt.Errorf("%s の更新に失敗しました: %w", filePath, err)
		}
		updatedFiles = append(updatedFiles, filePath)
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("remove-title-hash-tags: %d ファイルの先頭2行からハッシュタグを除去しました\n", len(updatedFiles)))
	for _, filePath := range updatedFiles {
		out.WriteString(fmt.Sprintf("- %s\n", filePath))
	}
	return out.String(), nil
}

func removeTitleHashTagsFromContent(content string) string {
	normalized := common.NormalizeNewlines(content)
	lines := strings.SplitAfter(normalized, "\n")

	targetLineCount := 2
	if len(lines) < targetLineCount {
		targetLineCount = len(lines)
	}

	for i := 0; i < targetLineCount; i++ {
		lines[i] = removeHashTagsFromLine(lines[i])
	}

	return strings.Join(lines, "")
}

func removeHashTagsFromLine(line string) string {
	lineEnding := ""
	if strings.HasSuffix(line, "\n") {
		lineEnding = "\n"
	}

	trimmed := strings.TrimSuffix(line, "\n")
	if trimmed == "" {
		return line
	}

	updated := lineStartHashTagPattern.ReplaceAllString(trimmed, "")
	updated = inlineHashTagPattern.ReplaceAllString(updated, "")
	return updated + lineEnding
}
