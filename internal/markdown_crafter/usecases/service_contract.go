package usecases

import "github.com/landmaster135/devbox/internal/markdown_crafter/domain"

type splitHeadingsOperation interface {
	Execute(filePath string, headingLevel int, outputDir, prefix string, sequencialDigits int) (string, error)
}

type addFrontMatterOperation interface {
	Execute(filePath string, kvPairs []string) (string, error)
}

type addTagsOperation interface {
	ExecuteByFile(filePath string, tagsCSV string) (string, error)
	ExecuteByDir(dirPath string, tagsCSV string) (string, error)
}

type deleteEmptyFilesOperation interface {
	Execute(dirPath string) (string, error)
}

type addHeading1Operation interface {
	Execute(filePath, headingText, headingPosition string) (string, error)
}

type replaceImagesOperation interface {
	Execute(filePath, replacementText string) (string, error)
}

type removeHeadingAnnotationsOperation interface {
	Execute(filePath string, headingLevel int) (string, error)
}

type removeTitleHashTagsOperation interface {
	Execute(dirPath string, startLine int, endLine int) (string, error)
}

func newServiceWithOperations(
	repository domain.Repository,
	splitHeadingsOp splitHeadingsOperation,
	addFrontMatterOp addFrontMatterOperation,
	addTagsOp addTagsOperation,
	deleteEmptyFilesOp deleteEmptyFilesOperation,
	addHeading1Op addHeading1Operation,
	replaceImagesOp replaceImagesOperation,
	removeHeadingAnnotationsOp removeHeadingAnnotationsOperation,
	removeTitleHashTagsOp removeTitleHashTagsOperation,
) *Service {
	return &Service{
		repository:                        repository,
		splitHeadingsOperation:            splitHeadingsOp,
		addFrontMatterOperation:           addFrontMatterOp,
		addTagsOperation:                  addTagsOp,
		deleteEmptyFilesOperation:         deleteEmptyFilesOp,
		addHeading1Operation:              addHeading1Op,
		replaceImagesOperation:            replaceImagesOp,
		removeHeadingAnnotationsOperation: removeHeadingAnnotationsOp,
		removeTitleHashTagsOperation:      removeTitleHashTagsOp,
	}
}
