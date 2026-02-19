package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	domain "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

var conNumberRegexp = regexp.MustCompile(`\d+`)

func BuildFrontMatterPairs(content domain.Content) []string {
	conID := strings.TrimSpace(content.ConID)
	url := NormalizeFrontMatterURL(content.URL)
	return []string{
		fmt.Sprintf("bought_at=%s", strings.TrimSpace(content.BoughtAt)),
		fmt.Sprintf("score_of_100=%d", content.Score),
		fmt.Sprintf("price_yen=%d", content.Price),
		fmt.Sprintf("con_id=%s", conID),
		fmt.Sprintf("url=%s", url),
	}
}

func NormalizeFrontMatterURL(url string) string {
	trimmed := strings.TrimSpace(url)
	if trimmed == "" {
		return `""`
	}
	return trimmed
}

func BuildTagsForContent(content domain.Content, frequentTags []string) ([]string, error) {
	tags := make([]string, 0, len(content.Tags)+4)

	categoryTag, err := MapCategoryTag(content.Category)
	if err != nil {
		return nil, err
	}
	owningStatusTag, err := MapOwningStatusTag(content.OwningStatus)
	if err != nil {
		return nil, err
	}

	tags = append(tags, RequiredBackupTag, categoryTag, owningStatusTag)

	trimmedColor := strings.TrimSpace(content.Color)
	if trimmedColor != "" {
		colorTag, err := MapColorTag(trimmedColor)
		if err != nil {
			return nil, err
		}
		tags = append(tags, colorTag)
	}
	for _, contentTag := range content.Tags {
		name := strings.ToLower(strings.TrimSpace(contentTag.PageTitle))
		if name == "" {
			continue
		}
		tags = append(tags, ResolveFrequentTag(name, frequentTags))
	}

	return uniqueTags(tags), nil
}

func MapCategoryTag(category string) (string, error) {
	switch NormalizeKey(category) {
	case "webclip":
		return "0a-content/web-clip", nil
	case "device":
		return "0a-content/device", nil
	case "software":
		return "0a-content/software", nil
	case "disc":
		return "0a-content/disc", nil
	case "subscriptionyear":
		return "0a-content/subscription-year", nil
	case "subscriptionmonth":
		return "0a-content/subscription-month", nil
	case "book":
		return "0a-content/book", nil
	default:
		return "", fmt.Errorf("未対応のcategoryです: %s", strings.TrimSpace(category))
	}
}

func MapOwningStatusTag(owningStatus string) (string, error) {
	switch NormalizeKey(owningStatus) {
	case "yet":
		return "01-p/own-status/1-yet", nil
	case "already":
		return "01-p/own-status/2-already", nil
	case "gone":
		return "01-p/own-status/3-gone", nil
	default:
		return "", fmt.Errorf("未対応のowning_statusです: %s", strings.TrimSpace(owningStatus))
	}
}

func MapColorTag(color string) (string, error) {
	switch NormalizeKey(color) {
	case "gray":
		return "01-p/color/gray", nil
	case "white":
		return "01-p/color/white", nil
	case "black":
		return "01-p/color/black", nil
	case "red":
		return "01-p/color/red", nil
	case "blue":
		return "01-p/color/blue", nil
	case "green":
		return "01-p/color/green", nil
	case "yellow":
		return "01-p/color/yellow", nil
	case "orange":
		return "01-p/color/orange", nil
	case "brown":
		return "01-p/color/brown", nil
	case "pink":
		return "01-p/color/pink", nil
	case "purple":
		return "01-p/color/purple", nil
	default:
		return "", fmt.Errorf("未対応のcolorです: %s", strings.TrimSpace(color))
	}
}

func NormalizeKey(value string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	return strings.NewReplacer("-", "", "_", "", " ", "").Replace(trimmed)
}

func ResolveFrequentTag(tagName string, frequentTags []string) string {
	normalizedTagName := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(tagName), "#"))
	if normalizedTagName == "" {
		return normalizedTagName
	}

	for _, frequentTag := range frequentTags {
		normalizedFrequentTag := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(frequentTag), "#"))
		if normalizedFrequentTag == "" {
			continue
		}

		if !strings.Contains(normalizedFrequentTag, "/") {
			if normalizedFrequentTag == normalizedTagName {
				return strings.TrimPrefix(strings.TrimSpace(frequentTag), "#")
			}
			continue
		}

		lastSlash := strings.LastIndex(normalizedFrequentTag, "/")
		if lastSlash < 0 || lastSlash+1 >= len(normalizedFrequentTag) {
			continue
		}
		if normalizedFrequentTag[lastSlash+1:] == normalizedTagName {
			return strings.TrimPrefix(strings.TrimSpace(frequentTag), "#")
		}
	}
	return normalizedTagName
}

func ParseConNumber(conID string) (int, error) {
	digits := conNumberRegexp.FindAllString(conID, -1)
	if len(digits) == 0 {
		return 0, fmt.Errorf("数値部分が見つかりません")
	}

	number, err := strconv.Atoi(digits[len(digits)-1])
	if err != nil {
		return 0, fmt.Errorf("数値変換に失敗しました: %w", err)
	}
	return number, nil
}

func LoadFrequentTags(fileSystem filesystem.Repository, tagsPath string) ([]string, error) {
	data, err := fileSystem.ReadFile(tagsPath)
	if err != nil {
		return nil, fmt.Errorf("tags.md の読み込みに失敗しました (%s): %w", tagsPath, err)
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	inFrequentSection := false
	seen := map[string]struct{}{}
	result := make([]string, 0)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "## ") {
			if trimmed == "## Frequent Tags" {
				inFrequentSection = true
				continue
			}
			if inFrequentSection {
				break
			}
			continue
		}
		if !inFrequentSection || strings.HasPrefix(trimmed, "### ") || trimmed == "" {
			continue
		}

		parts := strings.Fields(trimmed)
		for _, part := range parts {
			if !strings.HasPrefix(part, "#") {
				continue
			}
			normalized := strings.TrimPrefix(strings.TrimSpace(part), "#")
			if normalized == "" {
				continue
			}
			if _, exists := seen[normalized]; exists {
				continue
			}
			seen[normalized] = struct{}{}
			result = append(result, normalized)
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("tags.md の ## Frequent Tags セクションにタグが見つかりません (%s)", tagsPath)
	}
	return result, nil
}

func uniqueTags(tags []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		normalized := strings.TrimPrefix(strings.TrimSpace(tag), "#")
		if normalized == "" {
			continue
		}
		if _, exists := seen[normalized]; exists {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}
	return result
}
