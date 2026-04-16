package addheading1

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/markdown_crafter/config"
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

func (s *Service) Execute(filePath, headingText, headingPosition string) (string, error) {
	content, err := s.repository.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	trimmedHeading := strings.TrimSpace(headingText)
	trimmedHeading = strings.TrimLeft(trimmedHeading, "#")
	trimmedHeading = strings.TrimSpace(trimmedHeading)
	if trimmedHeading == "" {
		return "", fmt.Errorf("見出しテキストが空です")
	}
	headingLine := "# " + trimmedHeading

	hasFrontMatter, block, body, err := common.SplitFrontMatterBlock(content)
	if err != nil {
		return "", err
	}

	updatedBody, err := insertHeadingByPosition(body, headingLine, headingPosition)
	if err != nil {
		return "", err
	}

	var updated strings.Builder
	if hasFrontMatter {
		updated.WriteString(block)
		updated.WriteString("\n")
		updated.WriteString(updatedBody)
	} else {
		updated.WriteString(updatedBody)
	}

	if err := s.repository.WriteFile(filePath, updated.String()); err != nil {
		return "", err
	}

	return fmt.Sprintf("add-heading1: %s に見出しを追加しました (%s)\n", filePath, headingPosition), nil
}

func insertHeadingByPosition(body, headingLine, headingPosition string) (string, error) {
	normalizedBody := strings.TrimPrefix(common.NormalizeNewlines(body), "\n")

	switch headingPosition {
	case config.HeadingPositionHead:
		if normalizedBody == "" {
			return headingLine + "\n", nil
		}
		return headingLine + "\n\n" + normalizedBody, nil
	case config.HeadingPositionTail:
		trimmedBody := strings.TrimRight(normalizedBody, "\n")
		if trimmedBody == "" {
			return headingLine + "\n", nil
		}
		return trimmedBody + "\n\n" + headingLine + "\n", nil
	default:
		return "", fmt.Errorf("未対応の heading-position です: %s", headingPosition)
	}
}
