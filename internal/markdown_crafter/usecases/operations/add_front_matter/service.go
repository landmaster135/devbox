package addfrontmatter

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

func (s *Service) Execute(filePath string, kvPairs []string) (string, error) {
	content, err := s.repository.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hasFrontMatter, block, body, err := common.SplitFrontMatterBlock(content)
	if err != nil {
		return "", err
	}

	existingKeys := make([]string, 0)
	existingValues := map[string]string{}
	if hasFrontMatter {
		existingKeys, existingValues, err = common.ParseFrontMatterMap(block)
		if err != nil {
			return "", err
		}
	}

	newKeys, newValues, err := common.ParseKVPairs(kvPairs)
	if err != nil {
		return "", err
	}

	for _, key := range newKeys {
		if _, exists := existingValues[key]; !exists {
			existingKeys = append(existingKeys, key)
		}
		existingValues[key] = newValues[key]
	}

	frontMatter := common.BuildFrontMatter(existingKeys, existingValues)
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
