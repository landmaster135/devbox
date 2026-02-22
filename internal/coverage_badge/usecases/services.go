package usecases

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	filesystem "github.com/landmaster135/devbox/internal/coverage_badge/infrastructures/filesystem"
)

const (
	defaultWritePermission = 0o644
)

var totalCoveragePattern = regexp.MustCompile(`(?i)^total:\s*\(statements\)\s*([0-9]+(?:\.[0-9]+)?)%`)

// CreateBadgeInput はcreate-badge operationの入力。
type CreateBadgeInput struct {
	BadgeTitle      string
	CoverageFile    string
	GreenThreshold  int
	YellowThreshold int
	ForceColor      string
	BadgeLink       string
	BadgeValue      string
}

// PatchBadgeInput はpatch-badge operationの入力。
type PatchBadgeInput struct {
	CreateBadgeInput
	TargetFile string
	DryRun     bool
}

// PatchBadgeResult はpatch-badge operationの結果。
type PatchBadgeResult struct {
	TargetFile      string
	BadgeMarkdown   string
	PatchedContent  string
	ContentModified bool
	FileWritten     bool
}

// CoverageBadgeService はcoverage badge生成と更新処理を提供する。
type CoverageBadgeService struct {
	fileSystem filesystem.Repository
}

// NewCoverageBadgeService は標準Repositoryでサービスを作成する。
func NewCoverageBadgeService() *CoverageBadgeService {
	return NewCoverageBadgeServiceWithRepository(filesystem.NewRepository())
}

// NewCoverageBadgeServiceWithRepository は依存注入付きでサービスを作成する。
func NewCoverageBadgeServiceWithRepository(repo filesystem.Repository) *CoverageBadgeService {
	if repo == nil {
		repo = filesystem.NewRepository()
	}
	return &CoverageBadgeService{
		fileSystem: repo,
	}
}

// CreateBadge はバッジMarkdownを生成する。
func (s *CoverageBadgeService) CreateBadge(input CreateBadgeInput) (string, error) {
	coverageValue, err := s.resolveCoverageValue(input.BadgeValue, input.CoverageFile)
	if err != nil {
		return "", err
	}

	color := determineBadgeColor(coverageValue, input.GreenThreshold, input.YellowThreshold, input.ForceColor)
	coverageText := formatCoverageValue(coverageValue)

	return buildBadgeMarkdown(input.BadgeTitle, coverageText, color, input.BadgeLink), nil
}

// PatchBadge は対象Markdownのバッジ行を追加または更新する。
func (s *CoverageBadgeService) PatchBadge(input PatchBadgeInput) (*PatchBadgeResult, error) {
	badgeMarkdown, err := s.CreateBadge(input.CreateBadgeInput)
	if err != nil {
		return nil, err
	}

	rawContent, err := s.fileSystem.ReadFile(input.TargetFile)
	if err != nil {
		return nil, err
	}

	patchedContent, changed := patchCoverageBadge(
		string(rawContent),
		badgeMarkdown,
		input.BadgeTitle,
	)

	result := &PatchBadgeResult{
		TargetFile:      input.TargetFile,
		BadgeMarkdown:   badgeMarkdown,
		PatchedContent:  patchedContent,
		ContentModified: changed,
		FileWritten:     false,
	}

	if input.DryRun || !changed {
		return result, nil
	}

	if err := s.fileSystem.WriteFile(input.TargetFile, []byte(patchedContent), defaultWritePermission); err != nil {
		return nil, err
	}
	result.FileWritten = true

	return result, nil
}

func (s *CoverageBadgeService) resolveCoverageValue(badgeValue string, coverageFile string) (float64, error) {
	trimmedValue := strings.TrimSpace(badgeValue)
	if trimmedValue != "" {
		return parseCoverageValue(trimmedValue)
	}

	coverageReport, err := s.fileSystem.ReadFile(coverageFile)
	if err != nil {
		return 0, err
	}

	return parseCoverageFromReport(string(coverageReport))
}

func parseCoverageFromReport(content string) (float64, error) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		matches := totalCoveragePattern.FindStringSubmatch(trimmed)
		if len(matches) != 2 {
			continue
		}

		value, err := parseCoverageValue(matches[1])
		if err != nil {
			return 0, err
		}
		return value, nil
	}

	return 0, fmt.Errorf("カバレッジ値を coverage report から抽出できませんでした")
}

func parseCoverageValue(valueText string) (float64, error) {
	trimmed := strings.TrimSpace(valueText)
	trimmed = strings.TrimSuffix(trimmed, "%")
	if trimmed == "" {
		return 0, fmt.Errorf("カバレッジ値が空です")
	}

	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, fmt.Errorf("カバレッジ値の解析に失敗しました: %w", err)
	}
	if value < 0 || value > 100 {
		return 0, fmt.Errorf("カバレッジ値は 0 から 100 の範囲で指定してください")
	}

	return value, nil
}

func determineBadgeColor(coverageValue float64, greenThreshold int, yellowThreshold int, forceColor string) string {
	trimmedForceColor := strings.ToLower(strings.TrimSpace(forceColor))
	if trimmedForceColor != "" {
		return trimmedForceColor
	}

	if coverageValue >= float64(greenThreshold) {
		return "green"
	}
	if coverageValue >= float64(yellowThreshold) {
		return "yellow"
	}
	return "red"
}

func formatCoverageValue(value float64) string {
	return fmt.Sprintf("%.1f", value)
}

func buildBadgeMarkdown(badgeTitle string, coverageValue string, color string, badgeLink string) string {
	badgeURL := fmt.Sprintf(
		"https://img.shields.io/badge/%s-%s-%s",
		url.PathEscape(badgeTitle),
		url.PathEscape(coverageValue+"%"),
		url.PathEscape(color),
	)

	trimmedLink := strings.TrimSpace(badgeLink)
	if trimmedLink == "" {
		return fmt.Sprintf("![%s](%s)", badgeTitle, badgeURL)
	}

	return fmt.Sprintf("[![%s](%s)](%s)", badgeTitle, badgeURL, trimmedLink)
}

func patchCoverageBadge(content string, badgeMarkdown string, badgeTitle string) (string, bool) {
	trailingNewlineCount := len(content) - len(strings.TrimRight(content, "\n"))
	contentWithoutTrailingNewline := strings.TrimRight(content, "\n")
	lines := strings.Split(contentWithoutTrailingNewline, "\n")
	if contentWithoutTrailingNewline == "" {
		lines = []string{}
	}

	updatedLines := make([]string, 0, len(lines)+1)
	replaced := false
	changed := false
	titleLower := strings.ToLower(strings.TrimSpace(badgeTitle))

	for _, line := range lines {
		if isCoverageBadgeLine(line, titleLower) {
			if !replaced {
				updatedLines = append(updatedLines, badgeMarkdown)
				if strings.TrimSpace(line) != badgeMarkdown {
					changed = true
				}
				replaced = true
			} else {
				changed = true
			}
			continue
		}
		updatedLines = append(updatedLines, line)
	}

	if !replaced {
		insertIndex := findInsertIndex(updatedLines)
		updatedLines = insertAt(updatedLines, insertIndex, badgeMarkdown)
		changed = true
	}

	patched := strings.Join(updatedLines, "\n")
	if trailingNewlineCount > 0 {
		patched += strings.Repeat("\n", trailingNewlineCount)
	}

	return patched, changed
}

func isCoverageBadgeLine(line string, titleLower string) bool {
	trimmed := strings.TrimSpace(line)
	lower := strings.ToLower(trimmed)

	if !strings.Contains(lower, "img.shields.io/badge/") {
		return false
	}

	if titleLower != "" && strings.Contains(lower, titleLower) {
		return true
	}

	return strings.Contains(lower, "coverage")
}

func findInsertIndex(lines []string) int {
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "# ") {
			return i + 1
		}
	}
	return 0
}

func insertAt(lines []string, index int, value string) []string {
	if index <= 0 {
		return append([]string{value}, lines...)
	}
	if index >= len(lines) {
		return append(lines, value)
	}

	result := make([]string, 0, len(lines)+1)
	result = append(result, lines[:index]...)
	result = append(result, value)
	result = append(result, lines[index:]...)
	return result
}
