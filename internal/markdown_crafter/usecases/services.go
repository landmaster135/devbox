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
	removeheadingannotations "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/remove_heading_annotations"
	removetitlehashtags "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/remove_title_hash_tags"
	replaceimages "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/replace_images"
	splitheadings "github.com/landmaster135/devbox/internal/markdown_crafter/usecases/operations/split_headings"
)

type Service struct {
	repository                        domain.Repository
	splitHeadingsOperation            splitHeadingsOperation
	addFrontMatterOperation           addFrontMatterOperation
	addTagsOperation                  addTagsOperation
	deleteEmptyFilesOperation         deleteEmptyFilesOperation
	addHeading1Operation              addHeading1Operation
	replaceImagesOperation            replaceImagesOperation
	removeHeadingAnnotationsOperation removeHeadingAnnotationsOperation
	removeTitleHashTagsOperation      removeTitleHashTagsOperation
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
		removeheadingannotations.NewService(repo),
		removetitlehashtags.NewService(repo),
	)
}

func (s *Service) SplitHeadings(filePath string, headingLevel int, outputDir, prefix string, sequencialDigits int) (string, error) {
	return s.splitHeadingsOperation.Execute(filePath, headingLevel, outputDir, prefix, sequencialDigits)
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

func (s *Service) DeleteEmptyFiles(dirPath string) (string, error) {
	return s.deleteEmptyFilesOperation.Execute(dirPath)
}

func (s *Service) AddHeading1(filePath, headingText, headingPosition string) (string, error) {
	return s.addHeading1Operation.Execute(filePath, headingText, headingPosition)
}

func (s *Service) ReplaceImages(filePath, replacementText string) (string, error) {
	return s.replaceImagesOperation.Execute(filePath, replacementText)
}

func (s *Service) RemoveHeadingAnnotations(filePath string, headingLevel int) (string, error) {
	return s.removeHeadingAnnotationsOperation.Execute(filePath, headingLevel)
}

func (s *Service) RemoveTitleHashTags(dirPath string, startLine int, endLine int) (string, error) {
	return s.removeTitleHashTagsOperation.Execute(dirPath, startLine, endLine)
}

func (s *Service) ExecuteByConfig(cfg *config.Config) (string, error) {
	switch cfg.Operation {
	case config.OperationSplitHeadings:
		return s.SplitHeadings(cfg.FilePath, cfg.HeadingLevel, cfg.OutputDir, cfg.Prefix, cfg.SequencialDigits)
	case config.OperationAddFrontMatter:
		return s.AddFrontMatter(cfg.FilePath, cfg.KVPairs)
	case config.OperationAddTags:
		if strings.TrimSpace(cfg.DirPath) != "" {
			return s.AddTagsByDir(cfg.DirPath, cfg.Tags)
		}
		return s.AddTags(cfg.FilePath, cfg.Tags)
	case config.OperationDeleteEmptyFiles:
		return s.DeleteEmptyFiles(cfg.DirPath)
	case config.OperationAddHeading1:
		return s.AddHeading1(cfg.FilePath, cfg.HeadingText, cfg.HeadingPosition)
	case config.OperationReplaceImages:
		return s.ReplaceImages(cfg.FilePath, cfg.ReplacementText)
	case config.OperationRemoveHeadingAnnotations:
		return s.RemoveHeadingAnnotations(cfg.FilePath, cfg.HeadingLevel)
	case config.OperationRemoveTitleHashTags:
		return s.RemoveTitleHashTags(cfg.DirPath, cfg.StartLine, cfg.EndLine)
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
