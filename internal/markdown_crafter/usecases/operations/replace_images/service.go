package replaceimages

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/landmaster135/devbox/internal/markdown_crafter/domain"
)

var markdownImagePattern = regexp.MustCompile(`!\[[^\]]*]\([^)]+\)`)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Execute(filePath, replacementText string) (string, error) {
	trimmedReplacement := strings.TrimSpace(replacementText)
	if trimmedReplacement == "" {
		return "", fmt.Errorf("置換文字列が空です")
	}

	content, err := s.repository.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	updated, replacedCount := replaceMarkdownImages(content, trimmedReplacement)
	if err := s.repository.WriteFile(filePath, updated); err != nil {
		return "", err
	}

	return fmt.Sprintf("replace-images: %s の画像記法 %d 件を置換しました\n", filePath, replacedCount), nil
}

func replaceMarkdownImages(content, replacementText string) (string, int) {
	replaceCount := 0
	replaced := markdownImagePattern.ReplaceAllStringFunc(content, func(_ string) string {
		replaceCount++
		return replacementText
	})

	return replaced, replaceCount
}
