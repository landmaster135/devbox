package artifactcraftmarkdown

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	markdownconfig "github.com/landmaster135/devbox/internal/markdown_crafter/config"

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

func (s *Service) Execute(pageType, category string, skipsNoSrcBody bool, conNumberStart, conNumberEnd int, srcJSONFile, srcBodyDir, outDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	if trimmedPageType != common.SupportedPageTypeArtifact {
		return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}
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

	var artifacts []domain.Artifact
	if err := json.Unmarshal(rawJSON, &artifacts); err != nil {
		return "", fmt.Errorf("Artifact JSONの解析に失敗しました: %w", err)
	}

	tagsPath := filepath.Join(filepath.Dir(srcJSONFile), "tags.md")
	frequentTags, err := common.LoadFrequentTags(s.fileSystem, tagsPath)
	if err != nil {
		return "", err
	}
	artifactTags, err := common.LoadArtifactTags(s.fileSystem, tagsPath)
	if err != nil {
		return "", err
	}

	if err := s.fileSystem.MkdirAll(outDir, common.DefaultDirectoryPerm); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました (%s): %w", outDir, err)
	}

	normalizedCategoryFilter := common.NormalizeKey(strings.TrimSpace(category))
	markdownService := common.NewMarkdownService(s.fileSystem)

	totalTarget := 0
	processed := 0
	seenConID := map[string]struct{}{}
	for _, artifact := range artifacts {
		conID := strings.TrimSpace(artifact.ConID)
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
		if normalizedCategoryFilter != "" && common.NormalizeKey(artifact.Category) != normalizedCategoryFilter {
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

		if err := s.prependOutputURL(outPath, artifact.OutputURL); err != nil {
			return "", fmt.Errorf("output_url の追加に失敗しました (con_id=%s): %w", conID, err)
		}

		tags := common.BuildTagsForArtifact(artifact, frequentTags, artifactTags)
		if _, err := markdownService.AddTags(outPath, strings.Join(tags, ",")); err != nil {
			return "", fmt.Errorf("タグ追加に失敗しました (con_id=%s): %w", conID, err)
		}

		headingText := strings.TrimSpace(artifact.PageTitle)
		if headingText == "" {
			return "", fmt.Errorf("page_title が空です (con_id=%s)", conID)
		}
		if _, err := markdownService.AddHeading1(outPath, headingText, markdownconfig.HeadingPositionHead); err != nil {
			return "", fmt.Errorf("見出し追加に失敗しました (con_id=%s): %w", conID, err)
		}

		processed++
	}

	return fmt.Sprintf("処理完了\n対象件数=%d, 加工成功=%d", totalTarget, processed), nil
}

func (s *Service) prependOutputURL(filePath, outputURL string) error {
	trimmedOutputURL := strings.TrimSpace(outputURL)
	if trimmedOutputURL == "" {
		return nil
	}

	data, err := s.fileSystem.ReadFile(filePath)
	if err != nil {
		return err
	}

	normalizedBody := strings.TrimPrefix(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	var updated strings.Builder
	updated.WriteString(trimmedOutputURL)
	if normalizedBody != "" {
		updated.WriteString("\n\n")
		updated.WriteString(normalizedBody)
	} else {
		updated.WriteString("\n")
	}

	return s.fileSystem.WriteFile(filePath, []byte(updated.String()))
}
