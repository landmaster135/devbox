package usecases

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/disk_health/infrastructures/filesystem"
)

type Service struct {
	fileSystem filesystem.Repository
}

type ServiceOptions struct {
	FileSystem filesystem.Repository
}

func NewService(options ServiceOptions) *Service {
	fileSystem := options.FileSystem
	if fileSystem == nil {
		fileSystem = filesystem.NewOSRepository()
	}
	return &Service{fileSystem: fileSystem}
}

func (s *Service) AssessSmart(srcFile string, outputJSON bool, verbose bool) (string, error) {
	content, err := s.fileSystem.ReadFile(srcFile)
	if err != nil {
		return "", fmt.Errorf("SMART情報ファイルの読み込みに失敗しました: %w", err)
	}

	report, err := s.ParseSmartReport(string(content))
	if err != nil {
		return "", err
	}

	assessment := s.AssessReport(report)
	if outputJSON {
		return s.FormatJSON(assessment, verbose)
	}
	return s.FormatText(assessment, verbose), nil
}
