package common

import (
	"fmt"
	"net/url"
	"strings"
)

// ContentReader は content-file 解決で利用する最小契約。
type ContentReader interface {
	ReadFile(filePath string) ([]byte, error)
}

func NormalizeBaseURL(raw string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(raw), "/")
	if baseURL == "" {
		return "/api/v1"
	}
	if strings.HasSuffix(baseURL, "/api/v1") {
		return baseURL
	}
	return baseURL + "/api/v1"
}

func NormalizeMemoIdentifier(memo string) string {
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

func BuildMemoResourceName(memo string) string {
	memoID := NormalizeMemoIdentifier(memo)
	if memoID == "" {
		return ""
	}
	return "memos/" + memoID
}

func ResolveContent(content string, contentFile string, reader ContentReader) (string, error) {
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

	data, err := reader.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("content-file の読み込みに失敗しました: %w", err)
	}
	fileContent := string(data)
	if strings.TrimSpace(fileContent) == "" {
		return "", fmt.Errorf("content-file が空です: %s", contentFile)
	}
	return fileContent, nil
}
