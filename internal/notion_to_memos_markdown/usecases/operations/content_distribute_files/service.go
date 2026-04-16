package distributefiles

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
	domain "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
	common "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/common"
)

type Service struct {
	fileSystem filesystem.Repository
}

func NewService(fileSystem filesystem.Repository) *Service {
	return &Service{
		fileSystem: fileSystem,
	}
}

func (s *Service) Execute(pageType, srcJSONFile, srcBodyDir, outDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType != common.SupportedPageType {
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
	jsonConIDSet := common.BuildJSONConIDSet(contents)
	srcBodyTotal, srcBodyMapped, srcBodyUnmapped := common.CountSrcBodyMetrics(srcBodyFiles, jsonConIDSet)

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

		categoryDir := common.SanitizeCategory(content.Category)
		destCategoryDir := filepath.Join(outDir, categoryDir)
		if err := s.fileSystem.MkdirAll(destCategoryDir, common.DefaultDirectoryPerm); err != nil {
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
