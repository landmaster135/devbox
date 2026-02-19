package usecases

func (s *Service) DistributeFiles(pageType, srcJSONFile, srcBodyDir, outDir string) (string, error) {
	return s.distributeFilesOperation.Execute(pageType, srcJSONFile, srcBodyDir, outDir)
}

func (s *Service) CraftMarkdown(pageType, category string, skipsNoSrcBody bool, conNumberStart, conNumberEnd int, srcJSONFile, srcBodyDir, outDir string) (string, error) {
	return s.craftMarkdownOperation.Execute(pageType, category, skipsNoSrcBody, conNumberStart, conNumberEnd, srcJSONFile, srcBodyDir, outDir)
}

func (s *Service) CheckBodyLength(srcBodyDir string, threshold int) (string, error) {
	return s.checkBodyLengthOperation.Execute(srcBodyDir, threshold)
}
