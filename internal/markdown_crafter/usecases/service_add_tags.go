package usecases

import (
	"fmt"
	"strings"
)

func (s *Service) AddTags(filePath string, tagsCSV string) (string, error) {
	content, err := s.repository.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	tags := uniqueTrimmedTags(tagsCSV)
	tagLine, err := buildTagLine(tags)
	if err != nil {
		return "", err
	}

	hasFrontMatter, block, body, err := splitFrontMatterBlock(content)
	if err != nil {
		return "", err
	}

	body = strings.TrimPrefix(body, "\n")
	normalized := normalizeNewlines(content)

	var updated strings.Builder
	if hasFrontMatter {
		updated.WriteString(block)
		updated.WriteString("\n")
		updated.WriteString(tagLine)
		if body != "" {
			updated.WriteString("\n\n")
			updated.WriteString(body)
		} else {
			updated.WriteString("\n")
		}
	} else {
		plainBody := strings.TrimPrefix(normalized, "\n")
		updated.WriteString(tagLine)
		if plainBody != "" {
			updated.WriteString("\n\n")
			updated.WriteString(plainBody)
		} else {
			updated.WriteString("\n")
		}
	}

	if err := s.repository.WriteFile(filePath, updated.String()); err != nil {
		return "", err
	}

	return fmt.Sprintf("add-tags: %s にタグを追加しました (%s)\n", filePath, tagLine), nil
}
