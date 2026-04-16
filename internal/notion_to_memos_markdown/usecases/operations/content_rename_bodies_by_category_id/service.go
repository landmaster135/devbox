package renamebodiesbycategoryid

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	domain "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
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

func (s *Service) Execute(pageType string, conNumberStart, conNumberEnd int, srcJSONFile, srcResourceDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType != common.SupportedPageType {
		return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}
	if conNumberStart <= 0 || conNumberEnd <= 0 {
		return "", fmt.Errorf("con_number_start と con_number_end は1以上で指定してください")
	}
	if conNumberStart > conNumberEnd {
		return "", fmt.Errorf("con_number_start は con_number_end 以下である必要があります")
	}
	trimmedSrcJSONFile := strings.TrimSpace(srcJSONFile)
	if trimmedSrcJSONFile == "" {
		return "", fmt.Errorf("src-json-file パラメータは必須です")
	}
	trimmedSrcResourceDir := strings.TrimSpace(srcResourceDir)
	if trimmedSrcResourceDir == "" {
		return "", fmt.Errorf("src-resource-dir パラメータは必須です")
	}

	rawJSON, err := s.fileSystem.ReadFile(trimmedSrcJSONFile)
	if err != nil {
		return "", fmt.Errorf("src-json-file の読み込みに失敗しました: %w", err)
	}

	var contents []domain.Content
	if err := json.Unmarshal(rawJSON, &contents); err != nil {
		return "", fmt.Errorf("Content JSONの解析に失敗しました: %w", err)
	}

	categoryConIDMap, err := s.buildCategoryConIDMap(contents, conNumberStart, conNumberEnd)
	if err != nil {
		return "", err
	}

	filePaths, err := s.fileSystem.ListFilesRecursive(trimmedSrcResourceDir)
	if err != nil {
		return "", fmt.Errorf("src-resource-dir の読み取りに失敗しました: %w", err)
	}

	total := len(filePaths)
	renamed := 0
	skippedNoPrefix := 0
	skippedUnmapped := 0

	for _, srcPath := range filePaths {
		fileName := filepath.Base(srcPath)
		prefix, ok := extractPrefixToken(fileName)
		if !ok {
			skippedNoPrefix++
			continue
		}

		conID, exists := categoryConIDMap[prefix]
		if !exists {
			skippedUnmapped++
			continue
		}

		dstName := conID + fileName[len(prefix):]
		dstPath := filepath.Join(filepath.Dir(srcPath), dstName)

		exists, err := s.fileSystem.FileExists(dstPath)
		if err != nil {
			return "", fmt.Errorf("リネーム先ファイルの確認に失敗しました (%s): %w", dstPath, err)
		}
		if exists {
			return "", fmt.Errorf("リネーム先ファイルが既に存在します (%s)", dstPath)
		}

		if err := s.fileSystem.RenameFile(srcPath, dstPath); err != nil {
			return "", fmt.Errorf("ファイル名変更に失敗しました (%s -> %s): %w", srcPath, dstPath, err)
		}
		renamed++
	}

	return fmt.Sprintf(
		"処理完了\n対象Content件数=%d\n対象ファイル総数=%d\nリネーム成功=%d\nスキップ(プレフィックスなし)=%d\nスキップ(マップ未対応)=%d",
		len(categoryConIDMap), total, renamed, skippedNoPrefix, skippedUnmapped,
	), nil
}

func (s *Service) buildCategoryConIDMap(contents []domain.Content, conNumberStart, conNumberEnd int) (map[string]string, error) {
	result := map[string]string{}

	for _, content := range contents {
		conID := strings.TrimSpace(content.ConID)
		categoryID := strings.TrimSpace(content.CategoryID)
		if conID == "" || categoryID == "" {
			continue
		}

		conNumber, err := common.ParseConNumber(conID)
		if err != nil {
			return nil, fmt.Errorf("con_id の解析に失敗しました (%s): %w", conID, err)
		}
		if conNumber < conNumberStart || conNumber > conNumberEnd {
			continue
		}

		if existingConID, exists := result[categoryID]; exists && existingConID != conID {
			return nil, fmt.Errorf("category_id が重複しています (category_id=%s, con_id=%s, existing_con_id=%s)", categoryID, conID, existingConID)
		}
		result[categoryID] = conID
	}

	return result, nil
}

func extractPrefixToken(fileName string) (string, bool) {
	trimmed := strings.TrimSpace(fileName)
	if trimmed == "" {
		return "", false
	}

	var builder strings.Builder
	for _, r := range trimmed {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			break
		}
		builder.WriteRune(r)
	}

	prefix := builder.String()
	if prefix == "" {
		return "", false
	}

	return prefix, true
}
