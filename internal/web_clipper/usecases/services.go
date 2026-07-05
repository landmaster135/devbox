package usecases

import (
	"fmt"
	"time"

	"github.com/landmaster135/devbox/internal/web_clipper/config"
	"github.com/landmaster135/devbox/internal/web_clipper/infrastructures/filesystem"
	patchmarkdown "github.com/landmaster135/devbox/internal/web_clipper/usecases/operations/patch_markdown"
	renameattachments "github.com/landmaster135/devbox/internal/web_clipper/usecases/operations/rename_attachments"
)

type Service struct {
	repository                 filesystem.Repository
	now                        func() time.Time
	patchMarkdownOperation     patchMarkdownOperation
	renameAttachmentsOperation renameAttachmentsOperation
}

func NewService(repository filesystem.Repository) *Service {
	repo := repository
	if repo == nil {
		repo = filesystem.NewRepository()
	}

	return newServiceWithOperations(
		repo,
		patchmarkdown.NewService(repo),
		renameattachments.NewService(repo),
	)
}

func (s *Service) PatchMarkdown(targetTitle, targetURL, srcMarkdownContent, srcMarkdownFile, outFilePath string, topHeadingLevel int) (string, error) {
	return s.patchMarkdownOperation.Execute(targetTitle, targetURL, srcMarkdownContent, srcMarkdownFile, outFilePath, topHeadingLevel)
}

type RenameAttachmentsOptions = renameattachments.Options

func (s *Service) RenameAttachments(opts RenameAttachmentsOptions) (string, error) {
	return s.renameAttachmentsOperation.Execute(opts, s.now())
}

func (s *Service) ExecuteByConfig(cfg *config.Config) (string, error) {
	switch cfg.Operation {
	case config.OperationPatchMarkdown:
		return s.PatchMarkdown(
			cfg.TargetTitle,
			cfg.TargetURL,
			cfg.SrcMarkdownContent,
			cfg.SrcMarkdownFile,
			cfg.OutFilePath,
			cfg.TopHeadingLevel,
		)
	case config.OperationRenameAttachments:
		return s.RenameAttachments(RenameAttachmentsOptions{
			SrcDir:     cfg.SrcDir,
			Slug:       cfg.Slug,
			Start:      cfg.Start,
			Digits:     cfg.Digits,
			SortByTime: cfg.SortByTime,
			SortByName: cfg.SortByName,
			JSON:       cfg.JSON,
			Verbose:    cfg.Verbose,
		})
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}
}
