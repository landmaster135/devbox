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

func (s *Service) Execute(dirPath string, startLine int, endLine int) (string, error) {
	if startLine < 1 || endLine < 1 {
		return "", fmt.Errorf("startLine と endLine は 1 以上で指定してください")
	}
	if startLine > endLine {
		return "", fmt.Errorf("startLine は endLine 以下で指定してください")
	}

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

		updated := removeTitleHashTagsFromContent(content, startLine, endLine)
		if content == updated {
			continue
		}

		if err := s.repository.WriteFile(filePath, updated); err != nil {
			return "", fmt.Errorf("%s の更新に失敗しました: %w", filePath, err)
		}
		updatedFiles = append(updatedFiles, filePath)
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("remove-title-hash-tags: %d ファイルの %d 行目から %d 行目までのハッシュタグを除去しました\n", len(updatedFiles), startLine, endLine))
	for _, filePath := range updatedFiles {
		out.WriteString(fmt.Sprintf("- %s\n", filePath))
	}
	return out.String(), nil
}

func removeTitleHashTagsFromContent(content string, startLine int, endLine int) string {
	normalized := common.NormalizeNewlines(content)
	lines := strings.SplitAfter(normalized, "\n")

	startIndex := startLine - 1
	if startIndex >= len(lines) {
		return strings.Join(lines, "")
	}
	endIndex := endLine - 1
	if endIndex >= len(lines) {
		endIndex = len(lines) - 1
	}

	for i := startIndex; i <= endIndex; i++ {
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
