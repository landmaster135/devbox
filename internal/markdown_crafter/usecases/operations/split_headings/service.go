package splitheadings

import (
	"fmt"
	"path/filepath"
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

func (s *Service) Execute(filePath string, headingLevel int, outputDir string) (string, error) {
	content, err := s.repository.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	sections := ExtractSectionsByHeadingLevel(content, headingLevel)
	if len(sections) == 0 {
		return "", fmt.Errorf("見出しレベル %d のセクションが見つかりませんでした", headingLevel)
	}

	if err := s.repository.CreateDir(outputDir); err != nil {
		return "", err
	}

	outputs := make([]string, 0, len(sections))
	for i, section := range sections {
		fileName := fmt.Sprintf("%03d.md", i+1)
		outputPath := filepath.Join(outputDir, fileName)
		if err := s.repository.WriteFile(outputPath, section); err != nil {
			return "", err
		}
		outputs = append(outputs, outputPath)
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("split-headings: %d ファイルを出力しました\n", len(outputs)))
	for _, output := range outputs {
		builder.WriteString(fmt.Sprintf("- %s\n", output))
	}
	return builder.String(), nil
}

func ExtractSectionsByHeadingLevel(content string, headingLevel int) []string {
	normalized := common.NormalizeNewlines(content)
	lines := strings.SplitAfter(normalized, "\n")

	sections := make([]string, 0)
	var current strings.Builder
	inSection := false

	for _, line := range lines {
		if isHeadingOfLevel(line, headingLevel) {
			if inSection {
				sections = append(sections, current.String())
				current.Reset()
			}
			inSection = true
			current.WriteString(line)
			continue
		}

		if inSection {
			current.WriteString(line)
		}
	}

	if inSection {
		sections = append(sections, current.String())
	}

	return sections
}

func isHeadingOfLevel(line string, targetLevel int) bool {
	trimmed := strings.TrimRight(line, "\n")
	trimmed = strings.TrimRight(trimmed, "\r")
	if trimmed == "" {
		return false
	}

	level := 0
	for level < len(trimmed) && trimmed[level] == '#' {
		level++
	}
	if level != targetLevel {
		return false
	}

	if level >= len(trimmed) {
		return false
	}

	return trimmed[level] == ' '
}
