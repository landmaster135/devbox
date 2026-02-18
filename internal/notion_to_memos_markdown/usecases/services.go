package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

const (
	defaultCategory      = "uncategorized"
	supportedPageType    = "content"
	defaultDirectoryPerm = os.FileMode(0755)
)

type Service struct {
	fileSystem filesystem.Repository
}

func NewService(fileSystem filesystem.Repository) *Service {
	repo := fileSystem
	if repo == nil {
		repo = filesystem.NewRepository()
	}

	return &Service{
		fileSystem: repo,
	}
}

func (s *Service) DistributeFiles(pageType, srcJSONFile, srcBodyDir, outDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType != supportedPageType {
		return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}

	rawJSON, err := s.fileSystem.ReadFile(srcJSONFile)
	if err != nil {
		return "", fmt.Errorf("src-json-file の読み込みに失敗しました: %w", err)
	}

	var contents []domain.Content
	if err := json.Unmarshal(rawJSON, &contents); err != nil {
		return "", fmt.Errorf("Content JSONの解析に失敗しました: %w", err)
	}

	srcBodyFiles, err := s.fileSystem.ListMarkdownFiles(srcBodyDir)
	if err != nil {
		return "", fmt.Errorf("src-body-dir の読み取りに失敗しました: %w", err)
	}
	jsonConIDSet := buildJSONConIDSet(contents)
	srcBodyTotal, srcBodyMapped, srcBodyUnmapped := countSrcBodyMetrics(srcBodyFiles, jsonConIDSet)

	total := len(contents)
	copied := 0
	missing := 0
	skipped := 0
	seenConID := map[string]struct{}{}

	for _, content := range contents {
		conID := strings.TrimSpace(content.ConID)
		if conID == "" {
			skipped++
			continue
		}

		if _, exists := seenConID[conID]; exists {
			skipped++
			continue
		}
		seenConID[conID] = struct{}{}

		srcPath := filepath.Join(srcBodyDir, conID+".md")
		exists, err := s.fileSystem.FileExists(srcPath)
		if err != nil {
			return "", fmt.Errorf("コピー元ファイルの確認に失敗しました (%s): %w", srcPath, err)
		}
		if !exists {
			missing++
			continue
		}

		categoryDir := sanitizeCategory(content.Category)
		destCategoryDir := filepath.Join(outDir, categoryDir)
		if err := s.fileSystem.MkdirAll(destCategoryDir, defaultDirectoryPerm); err != nil {
			return "", fmt.Errorf("カテゴリディレクトリの作成に失敗しました (%s): %w", destCategoryDir, err)
		}

		destPath := filepath.Join(destCategoryDir, conID+".md")
		if err := s.fileSystem.CopyFile(srcPath, destPath); err != nil {
			return "", fmt.Errorf("Markdownのコピーに失敗しました (%s -> %s): %w", srcPath, destPath, err)
		}

		copied++
	}

	return fmt.Sprintf(
		"処理完了\nJSON基準: 総件数=%d, コピー成功=%d, 未検出=%d, スキップ=%d\nsrc-body-dir基準: 総md件数=%d, JSON対応=%d, JSON未対応=%d",
		total, copied, missing, skipped,
		srcBodyTotal, srcBodyMapped, srcBodyUnmapped,
	), nil
}

func buildJSONConIDSet(contents []domain.Content) map[string]struct{} {
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

func countSrcBodyMetrics(srcBodyFiles []string, jsonConIDSet map[string]struct{}) (int, int, int) {
	total := len(srcBodyFiles)
	mapped := 0

	for _, srcBodyFile := range srcBodyFiles {
		conID := extractConIDFromPath(srcBodyFile)
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

func extractConIDFromPath(path string) string {
	baseName := filepath.Base(path)
	if !strings.EqualFold(filepath.Ext(baseName), ".md") {
		return ""
	}
	return strings.TrimSpace(strings.TrimSuffix(baseName, filepath.Ext(baseName)))
}

func sanitizeCategory(category string) string {
	trimmed := strings.TrimSpace(category)
	if trimmed == "" {
		return defaultCategory
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
		return defaultCategory
	}
	return strings.Join(safeSegments, "_")
}
