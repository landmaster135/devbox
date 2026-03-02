package usecases

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/markdown_crafter/config"
)

func (s *Service) ExecuteByConfig(cfg *config.Config) (string, error) {
	switch cfg.Operation {
	case config.OperationSplitHeadings:
		return s.SplitHeadings(cfg.FilePath, cfg.HeadingLevel, cfg.OutputDir)
	case config.OperationAddFrontMatter:
		return s.AddFrontMatter(cfg.FilePath, cfg.KVPairs)
	case config.OperationAddTags:
		if strings.TrimSpace(cfg.DirPath) != "" {
			return s.AddTagsByDir(cfg.DirPath, cfg.Tags)
		}
		return s.AddTags(cfg.FilePath, cfg.Tags)
	case config.OperationDeleteEmptyFiles:
		return s.DeleteEmptyFiles(cfg.DirectoryPath)
	case config.OperationAddHeading1:
		return s.AddHeading1(cfg.FilePath, cfg.HeadingText, cfg.HeadingPosition)
	case config.OperationReplaceImages:
		return s.ReplaceImages(cfg.FilePath, cfg.ReplacementText)
	default:
		return "", fmt.Errorf("未サポートのoperationです: %s", cfg.Operation)
	}
}
