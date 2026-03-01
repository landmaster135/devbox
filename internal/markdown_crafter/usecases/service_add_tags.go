package usecases

import (
	"fmt"
	"strings"
)

func (s *Service) AddTags(filePath string, tagsCSV string) (string, error) {
	tags := uniqueTrimmedTags(tagsCSV)
	tagLine, err := buildTagLine(tags)
	if err != nil {
		return "", err
	}

	if err := s.addTagsToFile(filePath, tagLine); err != nil {
		return "", err
	}

	return fmt.Sprintf("add-tags: %s にタグを追加しました (%s)\n", filePath, tagLine), nil
}

func (s *Service) AddTagsByDir(dirPath string, tagsCSV string) (string, error) {
	tags := uniqueTrimmedTags(tagsCSV)
	tagLine, err := buildTagLine(tags)
	if err != nil {
		return "", err
	}

	markdownFiles, err := s.repository.ListMarkdownFiles(dirPath)
	if err != nil {
		return "", err
	}
	if len(markdownFiles) == 0 {
		return fmt.Sprintf("add-tags: %s 内のMarkdownファイルは0件でした\n", dirPath), nil
	}

	updatedFiles := make([]string, 0, len(markdownFiles))
	for _, filePath := range markdownFiles {
		if err := s.addTagsToFile(filePath, tagLine); err != nil {
			return "", fmt.Errorf("%s へのタグ追加に失敗しました: %w", filePath, err)
		}
		updatedFiles = append(updatedFiles, filePath)
	}

	var out strings.Builder
	out.WriteString(fmt.Sprintf("add-tags: %d ファイルにタグを追加しました (%s)\n", len(updatedFiles), tagLine))
	for _, filePath := range updatedFiles {
		out.WriteString(fmt.Sprintf("- %s\n", filePath))
	}
	return out.String(), nil
}

func (s *Service) addTagsToFile(filePath string, tagLine string) error {
	content, err := s.repository.ReadFile(filePath)
	if err != nil {
		return err
	}

	hasFrontMatter, block, body, err := splitFrontMatterBlock(content)
	if err != nil {
		return err
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
		return err
	}

	return nil
}
