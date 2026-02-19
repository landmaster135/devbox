package common

import (
	"os"
	"path/filepath"
	"strings"

	domain "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
)

func BuildJSONConIDSet(contents []domain.Content) map[string]struct{} {
	result := map[string]struct{}{}
	for _, content := range contents {
		conID := strings.TrimSpace(content.ConID)
		if conID == "" {
			continue
		}
		result[conID] = struct{}{}
	}
	return result
}

func CountSrcBodyMetrics(srcBodyFiles []string, jsonConIDSet map[string]struct{}) (int, int, int) {
	total := len(srcBodyFiles)
	mapped := 0

	for _, srcBodyFile := range srcBodyFiles {
		conID := ExtractConIDFromPath(srcBodyFile)
		if conID == "" {
			continue
		}
		if _, exists := jsonConIDSet[conID]; exists {
			mapped++
		}
	}

	unmapped := total - mapped
	return total, mapped, unmapped
}

func ExtractConIDFromPath(path string) string {
	baseName := filepath.Base(path)
	if !strings.EqualFold(filepath.Ext(baseName), ".md") {
		return ""
	}
	return strings.TrimSpace(strings.TrimSuffix(baseName, filepath.Ext(baseName)))
}

func SanitizeCategory(category string) string {
	trimmed := strings.TrimSpace(category)
	if trimmed == "" {
		return DefaultCategory
	}

	normalized := strings.ReplaceAll(trimmed, "\\", "/")
	segments := strings.Split(normalized, "/")
	safeSegments := make([]string, 0, len(segments))

	for _, segment := range segments {
		cleaned := strings.TrimSpace(segment)
		if cleaned == "" || cleaned == "." || cleaned == ".." {
			continue
		}

		cleaned = strings.Map(func(r rune) rune {
			if r == '/' || r == '\\' || r == os.PathSeparator || r == 0 {
				return '_'
			}
			return r
		}, cleaned)
		if cleaned == "" {
			continue
		}
		safeSegments = append(safeSegments, cleaned)
	}

	if len(safeSegments) == 0 {
		return DefaultCategory
	}
	return strings.Join(safeSegments, "_")
}
