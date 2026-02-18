package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	markdownconfig "github.com/landmaster135/devbox/internal/markdown_crafter/config"
	markdowndomain "github.com/landmaster135/devbox/internal/markdown_crafter/domain"
	markdownusecases "github.com/landmaster135/devbox/internal/markdown_crafter/usecases"
	"github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

const requiredBackupTag = "91-backup/tool-migration/202602-notion"

var conNumberRegexp = regexp.MustCompile(`\d+`)

func (s *Service) CraftMarkdown(pageType, category string, skipsNoSrcBody bool, conNumberStart, conNumberEnd int, srcJSONFile, srcBodyDir, outDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType != supportedPageType {
		return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}
	trimmedCategory := strings.TrimSpace(category)
	normalizedCategoryFilter := normalizeKey(trimmedCategory)
	if conNumberStart <= 0 || conNumberEnd <= 0 {
		return "", fmt.Errorf("con_number_start と con_number_end は1以上で指定してください")
	}
	if conNumberStart > conNumberEnd {
		return "", fmt.Errorf("con_number_start は con_number_end 以下である必要があります")
	}

	rawJSON, err := s.fileSystem.ReadFile(srcJSONFile)
	if err != nil {
		return "", fmt.Errorf("src-json-file の読み込みに失敗しました: %w", err)
	}

	var contents []domain.Content
	if err := json.Unmarshal(rawJSON, &contents); err != nil {
		return "", fmt.Errorf("Content JSONの解析に失敗しました: %w", err)
	}

	tagsPath := filepath.Join(filepath.Dir(srcJSONFile), "tags.md")
	frequentTags, err := s.loadFrequentTags(tagsPath)
	if err != nil {
		return "", err
	}

	if err := s.fileSystem.MkdirAll(outDir, defaultDirectoryPerm); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました (%s): %w", outDir, err)
	}

	markdownService := markdownusecases.NewService(newMarkdownCrafterRepositoryAdapter(s.fileSystem))

	totalTarget := 0
	processed := 0
	seenConID := map[string]struct{}{}
	for _, content := range contents {
		conID := strings.TrimSpace(content.ConID)
		if conID == "" {
			continue
		}
		if _, exists := seenConID[conID]; exists {
			continue
		}
		seenConID[conID] = struct{}{}

		conNumber, err := parseConNumber(conID)
		if err != nil {
			return "", fmt.Errorf("con_id の解析に失敗しました (%s): %w", conID, err)
		}
		if conNumber < conNumberStart || conNumber > conNumberEnd {
			continue
		}
		if normalizedCategoryFilter != "" && normalizeKey(content.Category) != normalizedCategoryFilter {
			continue
		}
		totalTarget++

		srcPath := filepath.Join(srcBodyDir, conID+".md")
		exists, err := s.fileSystem.FileExists(srcPath)
		if err != nil {
			return "", fmt.Errorf("コピー元ファイルの確認に失敗しました (%s): %w", srcPath, err)
		}

		outPath := filepath.Join(outDir, conID+".md")
		if exists {
			if err := s.fileSystem.CopyFile(srcPath, outPath); err != nil {
				return "", fmt.Errorf("Markdownのコピーに失敗しました (%s -> %s): %w", srcPath, outPath, err)
			}
		} else {
			if skipsNoSrcBody {
				continue
			}
			if err := s.fileSystem.WriteFile(outPath, []byte("")); err != nil {
				return "", fmt.Errorf("空Markdownの作成に失敗しました (%s): %w", outPath, err)
			}
		}

		headingText := strings.TrimSpace(content.PageTitle)
		if headingText == "" {
			return "", fmt.Errorf("page_title が空です (con_id=%s)", conID)
		}
		if _, err := markdownService.AddHeading1(outPath, headingText, markdownconfig.HeadingPositionHead); err != nil {
			return "", fmt.Errorf("見出し追加に失敗しました (con_id=%s): %w", conID, err)
		}

		frontMatterPairs := buildFrontMatterPairs(content)
		if _, err := markdownService.AddFrontMatter(outPath, frontMatterPairs); err != nil {
			return "", fmt.Errorf("front matter 追加に失敗しました (con_id=%s): %w", conID, err)
		}

		tags, err := buildTagsForContent(content, frequentTags)
		if err != nil {
			return "", fmt.Errorf("タグ生成に失敗しました (con_id=%s): %w", conID, err)
		}
		if _, err := markdownService.AddTags(outPath, strings.Join(tags, ",")); err != nil {
			return "", fmt.Errorf("タグ追加に失敗しました (con_id=%s): %w", conID, err)
		}

		processed++
	}

	return fmt.Sprintf("処理完了\n対象件数=%d, 加工成功=%d", totalTarget, processed), nil
}

func buildFrontMatterPairs(content domain.Content) []string {
	conID := strings.TrimSpace(content.ConID)
	url := normalizeFrontMatterURL(content.URL)
	return []string{
		fmt.Sprintf("bought_at=%s", strings.TrimSpace(content.BoughtAt)),
		fmt.Sprintf("score_of_100=%d", content.Score),
		fmt.Sprintf("price_yen=%d", content.Price),
		fmt.Sprintf("con_id=%s", conID),
		fmt.Sprintf("url=%s", url),
	}
}

func normalizeFrontMatterURL(url string) string {
	trimmed := strings.TrimSpace(url)
	if trimmed == "" {
		return `""`
	}
	return trimmed
}

func buildTagsForContent(content domain.Content, frequentTags []string) ([]string, error) {
	tags := make([]string, 0, len(content.Tags)+4)

	categoryTag, err := mapCategoryTag(content.Category)
	if err != nil {
		return nil, err
	}
	owningStatusTag, err := mapOwningStatusTag(content.OwningStatus)
	if err != nil {
		return nil, err
	}

	tags = append(tags, requiredBackupTag, categoryTag, owningStatusTag)

	trimmedColor := strings.TrimSpace(content.Color)
	if trimmedColor != "" {
		colorTag, err := mapColorTag(trimmedColor)
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
		tags = append(tags, resolveFrequentTag(name, frequentTags))
	}

	return uniqueTags(tags), nil
}

func mapCategoryTag(category string) (string, error) {
	switch normalizeKey(category) {
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

func mapOwningStatusTag(owningStatus string) (string, error) {
	switch normalizeKey(owningStatus) {
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

func mapColorTag(color string) (string, error) {
	switch normalizeKey(color) {
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

func normalizeKey(value string) string {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	return strings.NewReplacer("-", "", "_", "", " ", "").Replace(trimmed)
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

func resolveFrequentTag(tagName string, frequentTags []string) string {
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

func parseConNumber(conID string) (int, error) {
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

func (s *Service) loadFrequentTags(tagsPath string) ([]string, error) {
	data, err := s.fileSystem.ReadFile(tagsPath)
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

type markdownCrafterRepositoryAdapter struct {
	fileSystem filesystem.Repository
}

func newMarkdownCrafterRepositoryAdapter(fileSystem filesystem.Repository) markdowndomain.Repository {
	return &markdownCrafterRepositoryAdapter{
		fileSystem: fileSystem,
	}
}

func (a *markdownCrafterRepositoryAdapter) ReadFile(filePath string) (string, error) {
	data, err := a.fileSystem.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (a *markdownCrafterRepositoryAdapter) WriteFile(filePath string, content string) error {
	return a.fileSystem.WriteFile(filePath, []byte(content))
}

func (a *markdownCrafterRepositoryAdapter) CreateDir(dirPath string) error {
	return a.fileSystem.MkdirAll(dirPath, defaultDirectoryPerm)
}

func (a *markdownCrafterRepositoryAdapter) ListMarkdownFiles(dirPath string) ([]string, error) {
	return a.fileSystem.ListMarkdownFiles(dirPath)
}

func (a *markdownCrafterRepositoryAdapter) RemoveFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("ファイルの削除に失敗しました (%s): %w", filePath, err)
	}
	return nil
}
