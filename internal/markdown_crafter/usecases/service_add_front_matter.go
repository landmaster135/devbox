package usecases

import (
	"fmt"
	"strings"
)

func (s *Service) AddFrontMatter(filePath string, kvPairs []string) (string, error) {
	content, err := s.repository.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hasFrontMatter, block, body, err := splitFrontMatterBlock(content)
	if err != nil {
		return "", err
	}

	existingKeys := make([]string, 0)
	existingValues := map[string]string{}
	if hasFrontMatter {
		existingKeys, existingValues, err = parseFrontMatterMap(block)
		if err != nil {
			return "", err
		}
	}

	newKeys, newValues, err := parseKVPairs(kvPairs)
	if err != nil {
		return "", err
	}

	for _, key := range newKeys {
		if _, exists := existingValues[key]; !exists {
			existingKeys = append(existingKeys, key)
		}
		existingValues[key] = newValues[key]
	}

	frontMatter := buildFrontMatter(existingKeys, existingValues)
	body = strings.TrimPrefix(body, "\n")

	var updated strings.Builder
	updated.WriteString(frontMatter)
	if body != "" {
		updated.WriteString("\n")
		updated.WriteString(body)
	}

	if err := s.repository.WriteFile(filePath, updated.String()); err != nil {
		return "", err
	}

	return fmt.Sprintf("add-front-matter: %s を更新しました (%d キー)\n", filePath, len(existingKeys)), nil
}
