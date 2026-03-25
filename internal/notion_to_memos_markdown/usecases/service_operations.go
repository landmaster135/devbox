package usecases

import (
	"fmt"
	"strings"

	common "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/common"
)

func (s *Service) DistributeFiles(pageType, srcJSONFile, srcBodyDir, outDir string) (string, error) {
	return s.distributeFilesOperation.Execute(pageType, srcJSONFile, srcBodyDir, outDir)
}

func (s *Service) CraftMarkdown(pageType, category string, skipsNoSrcBody bool, conNumberStart, conNumberEnd int, srcJSONFile, srcBodyDir, srcResourceDir, outDir, outResourceDir string) (string, error) {
	trimmedPageType := strings.TrimSpace(pageType)
	switch trimmedPageType {
	case common.SupportedPageTypeContent:
		return s.craftMarkdownOperation.Execute(trimmedPageType, category, skipsNoSrcBody, conNumberStart, conNumberEnd, srcJSONFile, srcBodyDir, srcResourceDir, outDir, outResourceDir)
	case common.SupportedPageTypeArtifact:
		if s.artifactCraftOperation == nil {
			return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
		}
		return s.artifactCraftOperation.Execute(trimmedPageType, category, skipsNoSrcBody, conNumberStart, conNumberEnd, srcJSONFile, srcBodyDir, srcResourceDir, outDir, outResourceDir)
	case common.SupportedPageTypeTask:
		if s.taskCraftOperation == nil {
			return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
		}
		return s.taskCraftOperation.Execute(trimmedPageType, category, skipsNoSrcBody, conNumberStart, conNumberEnd, srcJSONFile, srcBodyDir, srcResourceDir, outDir, outResourceDir)
	default:
		return "", fmt.Errorf("未対応のpage-typeです: %s", trimmedPageType)
	}
}

func (s *Service) CheckBodyLength(srcBodyDir string, threshold int) (string, error) {
	return s.checkBodyLengthOperation.Execute(srcBodyDir, threshold)
}

func (s *Service) GrepStr(srcBodyDir, targetStr string) (string, error) {
	return s.grepStrOperation.Execute(srcBodyDir, targetStr)
}

func (s *Service) RenameBodiesByCategoryID(pageType string, conNumberStart, conNumberEnd int, srcJSONFile, srcResourceDir string) (string, error) {
	return s.renameBodiesOperation.Execute(pageType, conNumberStart, conNumberEnd, srcJSONFile, srcResourceDir)
}

func (s *Service) MigrateToMemos(pageType, baseURL, apiToken, srcBodyDir, srcResourceDir string) (string, error) {
	return s.migrateToMemosOperation.Execute(pageType, baseURL, apiToken, srcBodyDir, srcResourceDir)
}
