package usecases

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeBaseURL(raw string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(raw), "/")
	if baseURL == "" {
		return "/api/v1"
	}
	if strings.HasSuffix(baseURL, "/api/v1") {
		return baseURL
	}
	return baseURL + "/api/v1"
}

func normalizeMemoIdentifier(memo string) string {
	trimmed := strings.TrimSpace(memo)
	if trimmed == "" {
		return ""
	}

	if strings.Contains(trimmed, "://") {
		if parsed, err := url.Parse(trimmed); err == nil {
			trimmed = strings.Trim(parsed.Path, "/")
		}
	}

	trimmed = strings.Trim(trimmed, "/")
	trimmed = strings.TrimPrefix(trimmed, "api/v1/")
	trimmed = strings.TrimPrefix(trimmed, "memos/")
	return strings.TrimSpace(trimmed)
}

func buildMemoResourceName(memo string) string {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return ""
	}
	return "memos/" + memoID
}

func (s *Service) resolveContent(content string, contentFile string) (string, error) {
	hasContent := strings.TrimSpace(content) != ""
	hasContentFile := strings.TrimSpace(contentFile) != ""

	if hasContent && hasContentFile {
		return "", fmt.Errorf("content と content-file は同時に指定できません")
	}
	if !hasContent && !hasContentFile {
		return "", fmt.Errorf("content または content-file の指定が必要です")
	}

	if hasContent {
		return content, nil
	}

	data, err := s.fileSystem.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("content-file の読み込みに失敗しました: %w", err)
	}
	fileContent := string(data)
	if strings.TrimSpace(fileContent) == "" {
		return "", fmt.Errorf("content-file が空です: %s", contentFile)
	}
	return fileContent, nil
}
