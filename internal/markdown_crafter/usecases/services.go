package usecases

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/markdown_crafter/domain"
	"github.com/landmaster135/devbox/internal/markdown_crafter/infrastructures/filesystem"
)

type Service struct {
	repository domain.Repository
}

func NewService(repository domain.Repository) *Service {
	repo := repository
	if repo == nil {
		repo = filesystem.NewRepository()
	}

	return &Service{
		repository: repo,
	}
}

func normalizeNewlines(content string) string {
	return strings.ReplaceAll(content, "\r\n", "\n")
}

func splitFrontMatterBlock(content string) (bool, string, string, error) {
	normalized := normalizeNewlines(content)
	if !strings.HasPrefix(normalized, "---\n") {
		return false, "", normalized, nil
	}

	rest := strings.TrimPrefix(normalized, "---\n")
	endIdx := strings.Index(rest, "\n---\n")
	if endIdx >= 0 {
		block := "---\n" + rest[:endIdx] + "\n---\n"
		body := rest[endIdx+len("\n---\n"):]
		return true, block, body, nil
	}

	if strings.HasSuffix(rest, "\n---") {
		block := "---\n" + strings.TrimSuffix(rest, "\n---") + "\n---\n"
		return true, block, "", nil
	}

	return false, "", "", fmt.Errorf("front matter の終端 '---' が見つかりません")
}

func parseFrontMatterMap(block string) ([]string, map[string]string, error) {
	values := map[string]string{}
	keys := make([]string, 0)

	trimmed := strings.TrimPrefix(block, "---\n")
	trimmed = strings.TrimSuffix(trimmed, "\n---\n")
	if strings.TrimSpace(trimmed) == "" {
		return keys, values, nil
	}

	lines := strings.Split(trimmed, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("front matter の書式が不正です: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, nil, fmt.Errorf("front matter のキーが空です: %s", line)
		}
		value := strings.TrimSpace(parts[1])

		if _, exists := values[key]; !exists {
			keys = append(keys, key)
		}
		values[key] = value
	}

	return keys, values, nil
}

func parseKVPairs(kvPairs []string) ([]string, map[string]string, error) {
	values := map[string]string{}
	keys := make([]string, 0, len(kvPairs))

	for _, kv := range kvPairs {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return nil, nil, fmt.Errorf("--kv の形式が不正です: %s (key=value 形式で指定してください)", kv)
		}

		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, nil, fmt.Errorf("--kv のキーが空です: %s", kv)
		}
		value := strings.TrimSpace(parts[1])

		if _, exists := values[key]; !exists {
			keys = append(keys, key)
		}
		values[key] = value
	}

	return keys, values, nil
}

func buildFrontMatter(keys []string, values map[string]string) string {
	var builder strings.Builder
	builder.WriteString("---\n")
	for _, key := range keys {
		builder.WriteString(fmt.Sprintf("%s: %s\n", key, values[key]))
	}
	builder.WriteString("---\n")
	return builder.String()
}

func uniqueTrimmedTags(tagsCSV string) []string {
	rawTags := strings.Split(tagsCSV, ",")
	seen := map[string]struct{}{}
	tags := make([]string, 0, len(rawTags))

	for _, tag := range rawTags {
		trimmed := strings.TrimSpace(tag)
		trimmed = strings.TrimPrefix(trimmed, "#")
		if trimmed == "" {
			continue
		}
		if _, exists := seen[trimmed]; exists {
			continue
		}
		seen[trimmed] = struct{}{}
		tags = append(tags, trimmed)
	}
	return tags
}

func buildTagLine(tags []string) (string, error) {
	if len(tags) == 0 {
		return "", fmt.Errorf("有効なタグが見つかりません")
	}

	prefixed := make([]string, len(tags))
	for i, tag := range tags {
		prefixed[i] = "#" + tag
	}

	return strings.Join(prefixed, " "), nil
}
