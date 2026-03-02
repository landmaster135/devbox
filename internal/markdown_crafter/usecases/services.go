package usecases

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/markdown_crafter/config"
	"github.com/landmaster135/devbox/internal/markdown_crafter/domain"
	"github.com/landmaster135/devbox/internal/markdown_crafter/infrastructures/filesystem"
	"github.com/landmaster135/devbox/internal/markdown_crafter/usecases/common"
	addfrontmatter "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/add_front_matter"
	addheading1 "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/add_heading1"
	addtags "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/add_tags"
	deleteemptyfiles "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/delete_empty_files"
	replaceimages "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/replace_images"
	splitheadings "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/split_headings"
)

type Service struct {
	repository                domain.Repository
	splitHeadingsOperation    splitHeadingsOperation
	addFrontMatterOperation   addFrontMatterOperation
	addTagsOperation          addTagsOperation
	deleteEmptyFilesOperation deleteEmptyFilesOperation
	addHeading1Operation      addHeading1Operation
	replaceImagesOperation    replaceImagesOperation
}

func NewService(repository domain.Repository) *Service {
	repo := repository
	if repo == nil {
		repo = filesystem.NewRepository()
	}

	return newServiceWithOperations(
		repo,
		splitheadings.NewService(repo),
		addfrontmatter.NewService(repo),
		addtags.NewService(repo),
		deleteemptyfiles.NewService(repo),
		addheading1.NewService(repo),
		replaceimages.NewService(repo),
	)
}

func (s *Service) SplitHeadings(filePath string, headingLevel int, outputDir string) (string, error) {
	return s.splitHeadingsOperation.Execute(filePath, headingLevel, outputDir)
}

func (s *Service) AddFrontMatter(filePath string, kvPairs []string) (string, error) {
	return s.addFrontMatterOperation.Execute(filePath, kvPairs)
}

func (s *Service) AddTags(filePath string, tagsCSV string) (string, error) {
	return s.addTagsOperation.ExecuteByFile(filePath, tagsCSV)
}

func (s *Service) AddTagsByDir(dirPath string, tagsCSV string) (string, error) {
	return s.addTagsOperation.ExecuteByDir(dirPath, tagsCSV)
}

func (s *Service) DeleteEmptyFiles(directoryPath string) (string, error) {
	return s.deleteEmptyFilesOperation.Execute(directoryPath)
}

func (s *Service) AddHeading1(filePath, headingText, headingPosition string) (string, error) {
	return s.addHeading1Operation.Execute(filePath, headingText, headingPosition)
}

func (s *Service) ReplaceImages(filePath, replacementText string) (string, error) {
	return s.replaceImagesOperation.Execute(filePath, replacementText)
}

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

func normalizeNewlines(content string) string {
	return common.NormalizeNewlines(content)
}

func splitFrontMatterBlock(content string) (bool, string, string, error) {
	return common.SplitFrontMatterBlock(content)
}

func parseFrontMatterMap(block string) ([]string, map[string]string, error) {
	return common.ParseFrontMatterMap(block)
}

func parseKVPairs(kvPairs []string) ([]string, map[string]string, error) {
	return common.ParseKVPairs(kvPairs)
}

func buildFrontMatter(keys []string, values map[string]string) string {
	return common.BuildFrontMatter(keys, values)
}

func uniqueTrimmedTags(tagsCSV string) []string {
	return common.UniqueTrimmedTags(tagsCSV)
}

func buildTagLine(tags []string) (string, error) {
	return common.BuildTagLine(tags)
}
