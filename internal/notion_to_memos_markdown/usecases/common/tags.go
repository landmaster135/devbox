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

var artifactSystemTagMatchers = []struct {
	keyword string
	tag     string
}{
	{keyword: "devbox", tag: "06-af/system/devbox"},
	{keyword: "dotfiles", tag: "06-af/system/dotfiles"},
	{keyword: "db-server-brewery", tag: "06-af/system/db-server-brewery"},
	{keyword: "notion-synchronizer", tag: "06-af/system/notion-synchronizer"},
	{keyword: "chrome-forge", tag: "06-af/system/chrome-forge"},
	{keyword: "others", tag: "06-af/system/others"},
}

var artifactArticleTagByPageTitle = map[string]string{
	"ブログ活動（投稿部）：エンドルフィン風呂に浸かる": "06-af/article/blog",
	"note活動（投稿部）": "06-af/article/note",
	"INTPなワイは100話で3年勤めたサラリーマンを辞める（マンガ）": "06-af/article/comic",
	"Palworldプレイ日記":                        "06-af/diary/game",
	"Azur Laneプレイ日記かつ作業記録":                 "06-af/diary/game",
	"World of Warshipsプレイ日記":               "06-af/diary/game",
	"Epistory_プレイ日記":                       "06-af/diary/game",
	"Monster Hunter World: Ice Borneプレイ日記": "06-af/diary/game",
	"原神（Genshin Impact）_ゲームプレイ日記":          "06-af/diary/game",
	"Monster Hunter Wildsプレイ日記":            "06-af/diary/game",
	"ワインを飲んだり勉強するムーヴメント":                   "06-af/diary/hobby",
	"料理のレシピのまとめ集":                          "06-af/diary/hobby",
	"ペンギンを愛でたり勉強するムーヴメント":                  "06-af/diary/hobby",
	"サバ缶を食ったり鯖缶を勉強するムーヴメント":                "06-af/diary/hobby",
	"旅行日程プラン計画まとめ集":                        "06-af/diary/hobby",
	"画像生成AIによるイラストで同人活動":                   "06-af/diary/hobby",
	"生け花を学ぶムーヴメント":                         "06-af/diary/hobby",
	"過去のメールのやり取りまとめ集":                      "mail",
}

var ignoredArtifactTagTitles = map[string]struct{}{
	"programming": {},
	"google":      {},
}

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

func BuildTagsForArtifact(artifact domain.Artifact, frequentTags, artifactTags []string) []string {
	tags := make([]string, 0, len(artifact.Tags)+2)
	tags = append(tags, RequiredBackupTag)

	for _, artifactTag := range artifact.Tags {
		name := strings.ToLower(strings.TrimSpace(artifactTag.PageTitle))
		if name == "" {
			continue
		}
		if _, shouldIgnore := ignoredArtifactTagTitles[name]; shouldIgnore {
			continue
		}
		tags = append(tags, ResolveFrequentTagByPartialMatch(name, frequentTags))
	}

	normalizedCategory := NormalizeKey(artifact.Category)
	switch normalizedCategory {
	case "system":
		if systemTag := MapArtifactSystemTagByPageTitle(artifact.PageTitle); systemTag != "" {
			tags = append(tags, systemTag)
		}
	case "article":
		if articleTag := MapArtifactArticleTagByPageTitle(artifact.PageTitle); articleTag != "" {
			tags = append(tags, articleTag)
		}
	default:
		if categoryTag := ResolveArtifactCategoryTag(artifact.Category, artifactTags); categoryTag != "" {
			tags = append(tags, categoryTag)
		}
	}

	return uniqueTags(tags)
}

func MapArtifactSystemTagByPageTitle(pageTitle string) string {
	normalizedPageTitle := NormalizeKey(pageTitle)
	if normalizedPageTitle == "" {
		return ""
	}

	for _, matcher := range artifactSystemTagMatchers {
		if strings.Contains(normalizedPageTitle, NormalizeKey(matcher.keyword)) {
			return matcher.tag
		}
	}
	return ""
}

func MapArtifactArticleTagByPageTitle(pageTitle string) string {
	return artifactArticleTagByPageTitle[strings.TrimSpace(pageTitle)]
}

func ResolveArtifactCategoryTag(category string, artifactTags []string) string {
	normalizedCategory := NormalizeKey(category)
	if normalizedCategory == "" {
		return ""
	}

	for _, artifactTag := range artifactTags {
		trimmedTag := strings.TrimPrefix(strings.TrimSpace(artifactTag), "#")
		if trimmedTag == "" {
			continue
		}
		normalizedTag := NormalizeKey(trimmedTag)
		if strings.Contains(normalizedTag, normalizedCategory) || strings.Contains(normalizedCategory, normalizedTag) {
			return trimmedTag
		}
	}
	return ""
}

func ResolveFrequentTagByPartialMatch(tagName string, frequentTags []string) string {
	resolved := ResolveFrequentTag(tagName, frequentTags)
	if resolved != strings.TrimSpace(strings.TrimPrefix(strings.ToLower(tagName), "#")) {
		return resolved
	}

	normalizedTagName := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(tagName), "#"))
	if normalizedTagName == "" {
		return normalizedTagName
	}
	if len(normalizedTagName) < 5 {
		return normalizedTagName
	}

	for _, frequentTag := range frequentTags {
		trimmedFrequentTag := strings.TrimPrefix(strings.TrimSpace(frequentTag), "#")
		if trimmedFrequentTag == "" {
			continue
		}
		normalizedFrequentTag := strings.ToLower(trimmedFrequentTag)
		if strings.Contains(normalizedFrequentTag, normalizedTagName) || strings.Contains(normalizedTagName, normalizedFrequentTag) {
			return trimmedFrequentTag
		}
	}
	return normalizedTagName
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
	return loadTagsBySection(fileSystem, tagsPath, "## Frequent Tags")
}

func LoadArtifactTags(fileSystem filesystem.Repository, tagsPath string) ([]string, error) {
	return loadTagsBySection(fileSystem, tagsPath, "## Artifact")
}

func loadTagsBySection(fileSystem filesystem.Repository, tagsPath, targetSection string) ([]string, error) {
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
			if trimmed == targetSection {
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
		return nil, fmt.Errorf("tags.md の %s セクションにタグが見つかりません (%s)", targetSection, tagsPath)
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
