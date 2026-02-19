package usecases

type distributeFilesOperation interface {
	Execute(pageType, srcJSONFile, srcBodyDir, outDir string) (string, error)
}

type craftMarkdownOperation interface {
	Execute(pageType, category string, skipsNoSrcBody bool, conNumberStart, conNumberEnd int, srcJSONFile, srcBodyDir, outDir string) (string, error)
}

type checkBodyLengthOperation interface {
	Execute(srcBodyDir string, threshold int) (string, error)
}

type grepStrOperation interface {
	Execute(srcBodyDir, targetStr string) (string, error)
}

func newServiceWithOperations(distributeFilesOp distributeFilesOperation, craftMarkdownOp craftMarkdownOperation, checkBodyLengthOp checkBodyLengthOperation, grepStrOp grepStrOperation) *Service {
	return &Service{
		distributeFilesOperation: distributeFilesOp,
		craftMarkdownOperation:   craftMarkdownOp,
		checkBodyLengthOperation: checkBodyLengthOp,
		grepStrOperation:         grepStrOp,
	}
}
