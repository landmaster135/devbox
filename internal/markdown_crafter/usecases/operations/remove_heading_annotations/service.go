package removeheadingannotations

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/markdown_crafter/domain"
	"github.com/landmaster135/devbox/internal/markdown_crafter/usecases/common"
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Execute(filePath string, headingLevel int) (string, error) {
	if headingLevel < 1 || headingLevel > 6 {
		return "", fmt.Errorf("--heading-level は 1 から 6 の範囲で指定してください")
	}

	content, err := s.repository.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	updated, replacedCount := removeHeadingAnnotations(content, headingLevel)
	if err := s.repository.WriteFile(filePath, updated); err != nil {
		return "", err
	}

	return fmt.Sprintf("remove-heading-annotations: %s の見出し注釈 %d 件を除去しました\n", filePath, replacedCount), nil
}

func removeHeadingAnnotations(content string, headingLevel int) (string, int) {
	normalized := common.NormalizeNewlines(content)
	lines := strings.SplitAfter(normalized, "\n")
	targetPrefix := strings.Repeat("#", headingLevel)

	replacedCount := 0
	for i, line := range lines {
		updatedLine, replaced := removeHeadingAnnotationFromLine(line, targetPrefix)
		if replaced {
			lines[i] = updatedLine
			replacedCount++
		}
	}

	return strings.Join(lines, ""), replacedCount
}

func removeHeadingAnnotationFromLine(line, targetPrefix string) (string, bool) {
	lineEnding := ""
	if strings.HasSuffix(line, "\n") {
		lineEnding = "\n"
	}

	trimmed := strings.TrimSuffix(line, "\n")
	if trimmed == "" {
		return line, false
	}

	indentCount := leadingSpaceCount(trimmed)
	if indentCount > 3 {
		return line, false
	}

	indent := trimmed[:indentCount]
	body := trimmed[indentCount:]

	if !strings.HasPrefix(body, targetPrefix) {
		return line, false
	}
	if len(body) == len(targetPrefix) || body[len(targetPrefix)] != ' ' {
		return line, false
	}

	headingText := strings.TrimSpace(body[len(targetPrefix):])
	if !strings.HasPrefix(headingText, "**") || !strings.HasSuffix(headingText, "**") || len(headingText) < 4 {
		return line, false
	}

	unwrapped := strings.TrimSpace(headingText[2 : len(headingText)-2])
	if unwrapped == "" {
		return line, false
	}

	return indent + targetPrefix + " " + unwrapped + lineEnding, true
}

func leadingSpaceCount(value string) int {
	count := 0
	for count < len(value) && value[count] == ' ' {
		count++
	}
	return count
}
