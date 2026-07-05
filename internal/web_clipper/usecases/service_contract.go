package usecases

import (
	"time"

	"github.com/landmaster135/devbox/internal/web_clipper/infrastructures/filesystem"
)

type patchMarkdownOperation interface {
	Execute(targetTitle, targetURL, srcMarkdownContent, srcMarkdownFile, outFilePath string, topHeadingLevel int) (string, error)
}

type renameAttachmentsOperation interface {
	Execute(opts RenameAttachmentsOptions, now time.Time) (string, error)
}

func newServiceWithOperations(
	repository filesystem.Repository,
	patchMarkdownOp patchMarkdownOperation,
	renameAttachmentsOp renameAttachmentsOperation,
) *Service {
	return &Service{
		repository:                 repository,
		now:                        time.Now,
		patchMarkdownOperation:     patchMarkdownOp,
		renameAttachmentsOperation: renameAttachmentsOp,
	}
}
