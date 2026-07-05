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

func (s *Service) Execute(filePath string, headingLevel int, outputDir, prefix string, sequencialDigits int) (string, error) {
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

	effectivePrefix := prefix
	if strings.TrimSpace(effectivePrefix) == "" {
		effectivePrefix = defaultPrefix(filePath)
	}

	outputs := make([]string, 0, len(sections))
	for i, section := range sections {
		fileName := buildOutputFileName(effectivePrefix, sequencialDigits, i+1)
		outputPath := filepath.Join(outputDir, fileName)
		if err := s.repository.WriteFile(outputPath, withSingleTrailingNewline(section)); err != nil {
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

func defaultPrefix(filePath string) string {
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if strings.TrimSpace(name) == "" {
		return "section"
	}
	return name
}

func withSingleTrailingNewline(content string) string {
	trimmed := strings.TrimRight(content, "\r\n")
	return trimmed + "\n"
}

func buildOutputFileName(prefix string, sequencialDigits int, index int) string {
	normalizedPrefix := prefix
	if !strings.HasSuffix(normalizedPrefix, "_") {
		normalizedPrefix += "_"
	}
	return fmt.Sprintf("%s%0*d.md", normalizedPrefix, sequencialDigits, index)
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
