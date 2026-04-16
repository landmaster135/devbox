package craftmarkdown

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	markdownconfig "github.com/landmaster135/devbox/internal/markdown_crafter/config"
	markdowncommon "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/common"

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

func (s *Service) Execute(pageType, category string, skipsNoSrcBody bool, conNumberStart, conNumberEnd int, srcJSONFile, srcBodyDir, srcResourceDir, outDir, outResourceDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType != common.SupportedPageType {
		return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}
	_ = srcResourceDir
	_ = outResourceDir
	trimmedCategory := strings.TrimSpace(category)
	normalizedCategoryFilter := common.NormalizeKey(trimmedCategory)
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
	frequentTags, err := common.LoadFrequentTags(s.fileSystem, tagsPath)
	if err != nil {
		return "", err
	}

	if err := s.fileSystem.MkdirAll(outDir, common.DefaultDirectoryPerm); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました (%s): %w", outDir, err)
	}

	markdownService := common.NewMarkdownService(s.fileSystem)
	srcBodyFiles, err := s.fileSystem.ListMarkdownFiles(srcBodyDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("src-body-dir の読み込みに失敗しました (%s): %w", srcBodyDir, err)
		}
		srcBodyFiles = []string{}
	}
	splitIndex := common.BuildSplitMarkdownIndex(srcBodyFiles)

	processOutput := func(content domain.Content, conID, outPath string) error {
		headingText := strings.TrimSpace(content.PageTitle)
		if headingText == "" {
			return fmt.Errorf("page_title が空です (con_id=%s)", conID)
		}

		hasFrontMatter, err := s.hasFrontMatter(outPath)
		if err != nil {
			return fmt.Errorf("front matter の判定に失敗しました (con_id=%s): %w", conID, err)
		}
		if !hasFrontMatter {
			if _, err := markdownService.AddHeading1(outPath, headingText, markdownconfig.HeadingPositionHead); err != nil {
				return fmt.Errorf("見出し追加に失敗しました (con_id=%s): %w", conID, err)
			}
		}

		frontMatterPairs := common.BuildFrontMatterPairs(content)
		if _, err := markdownService.AddFrontMatter(outPath, frontMatterPairs); err != nil {
			return fmt.Errorf("front matter 追加に失敗しました (con_id=%s): %w", conID, err)
		}

		if hasFrontMatter {
			return nil
		}

		tags, err := common.BuildTagsForContent(content, frequentTags)
		if err != nil {
			return fmt.Errorf("タグ生成に失敗しました (con_id=%s): %w", conID, err)
		}
		if _, err := markdownService.AddTags(outPath, strings.Join(tags, ",")); err != nil {
			return fmt.Errorf("タグ追加に失敗しました (con_id=%s): %w", conID, err)
		}
		return nil
	}

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

		conNumber, err := common.ParseConNumber(conID)
		if err != nil {
			return "", fmt.Errorf("con_id の解析に失敗しました (%s): %w", conID, err)
		}
		if conNumber < conNumberStart || conNumber > conNumberEnd {
			continue
		}
		if normalizedCategoryFilter != "" && common.NormalizeKey(content.Category) != normalizedCategoryFilter {
			continue
		}

		exactSrcPath := filepath.Join(srcBodyDir, conID+".md")
		exactExists, err := s.fileSystem.FileExists(exactSrcPath)
		if err != nil {
			return "", fmt.Errorf("コピー元ファイルの確認に失敗しました (%s): %w", exactSrcPath, err)
		}

		sourcePaths := splitIndex.Resolve(conID)
		if exactExists {
			sourcePaths = []string{exactSrcPath}
		}

		if len(sourcePaths) == 0 {
			totalTarget++
			outPath := filepath.Join(outDir, conID+".md")
			if skipsNoSrcBody {
				continue
			}
			if err := s.fileSystem.WriteFile(outPath, []byte("")); err != nil {
				return "", fmt.Errorf("空Markdownの作成に失敗しました (%s): %w", outPath, err)
			}
			if err := processOutput(content, conID, outPath); err != nil {
				return "", err
			}
			processed++
			continue
		}

		totalTarget += len(sourcePaths)
		for _, srcPath := range sourcePaths {
			outPath := filepath.Join(outDir, filepath.Base(srcPath))
			if err := s.fileSystem.CopyFile(srcPath, outPath); err != nil {
				return "", fmt.Errorf("Markdownのコピーに失敗しました (%s -> %s): %w", srcPath, outPath, err)
			}
			if err := processOutput(content, conID, outPath); err != nil {
				return "", err
			}
			processed++
		}
	}

	return fmt.Sprintf("処理完了\n対象件数=%d, 加工成功=%d", totalTarget, processed), nil
}

func (s *Service) hasFrontMatter(filePath string) (bool, error) {
	content, err := s.fileSystem.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	hasFrontMatter, _, _, err := markdowncommon.SplitFrontMatterBlock(string(content))
	if err != nil {
		return false, err
	}
	return hasFrontMatter, nil
}
