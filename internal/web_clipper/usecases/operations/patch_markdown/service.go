package patchmarkdown

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/web_clipper/infrastructures/filesystem"
	"github.com/landmaster135/devbox/internal/web_clipper/usecases/common"
)

type Service struct {
	repository filesystem.Repository
}

func NewService(repository filesystem.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Execute(targetTitle, targetURL, srcMarkdownContent, srcMarkdownFile, outFilePath string, topHeadingLevel int) (string, error) {
	trimmedTargetTitle := strings.TrimSpace(targetTitle)
	trimmedTargetURL := strings.TrimSpace(targetURL)
	trimmedSrcMarkdownFile := strings.TrimSpace(srcMarkdownFile)
	trimmedOutFilePath := strings.TrimSpace(outFilePath)

	if trimmedTargetTitle == "" {
		return "", fmt.Errorf("--target-title は必須です")
	}
	if trimmedTargetURL == "" {
		return "", fmt.Errorf("--target-url は必須です")
	}
	if strings.TrimSpace(srcMarkdownContent) == "" && trimmedSrcMarkdownFile == "" {
		return "", fmt.Errorf("--src-markdown-content または --src-markdown-file のいずれかは必須です")
	}
	if strings.TrimSpace(srcMarkdownContent) != "" && trimmedSrcMarkdownFile != "" {
		return "", fmt.Errorf("--src-markdown-content と --src-markdown-file は同時に指定できません")
	}
	if trimmedOutFilePath == "" {
		return "", fmt.Errorf("--out-file-path は必須です")
	}
	if strings.Contains(trimmedOutFilePath, ",") {
		return "", fmt.Errorf("--out-file-path にカンマは使用できません")
	}
	if topHeadingLevel < 1 {
		return "", fmt.Errorf("--top-heading-level は 1 以上で指定してください")
	}

	markdownContent, err := s.resolveMarkdownContent(srcMarkdownContent, trimmedSrcMarkdownFile)
	if err != nil {
		return "", err
	}

	normalizedContent := common.NormalizeNewlines(markdownContent)
	if common.ContainsHeadingLevel4OrMore(normalizedContent) {
		return "", fmt.Errorf("見出しレベル4以上（#### 以降）は使用できません")
	}

	patchedContent, err := addWebArticleInfo(normalizedContent, trimmedTargetTitle, trimmedTargetURL, topHeadingLevel)
	if err != nil {
		return "", err
	}

	if err := s.repository.WriteFile(trimmedOutFilePath, []byte(patchedContent)); err != nil {
		return "", fmt.Errorf("出力ファイルへの書き込みに失敗しました: %w", err)
	}

	return fmt.Sprintf("出力しました: %s", trimmedOutFilePath), nil
}

func (s *Service) resolveMarkdownContent(srcMarkdownContent, srcMarkdownFile string) (string, error) {
	if srcMarkdownFile != "" {
		data, err := s.repository.ReadFile(srcMarkdownFile)
		if err != nil {
			return "", fmt.Errorf("入力ファイルの読み込みに失敗しました: %w", err)
		}
		return string(data), nil
	}

	return srcMarkdownContent, nil
}

func addWebArticleInfo(markdownContent, targetTitle, targetURL string, topHeadingLevel int) (string, error) {
	lines := strings.Split(markdownContent, "\n")
	targetHeadingIndex := -1

	for idx, line := range lines {
		if common.HeadingLevel(line) == topHeadingLevel {
			targetHeadingIndex = idx
			break
		}
	}

	if targetHeadingIndex == -1 {
		return "", fmt.Errorf("見出しレベル%d が見つかりませんでした", topHeadingLevel)
	}

	linkLine := fmt.Sprintf("- [%s](%s)", targetTitle, targetURL)

	outputLines := make([]string, 0, len(lines)+1)
	outputLines = append(outputLines, lines[:targetHeadingIndex+1]...)
	outputLines = append(outputLines, linkLine)
	outputLines = append(outputLines, lines[targetHeadingIndex+1:]...)

	return strings.Join(outputLines, "\n"), nil
}
