package assessdmesg

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/disk_health/infrastructures/filesystem"
)

type Service struct {
	fileSystem filesystem.Repository
}

func NewService(fileSystem filesystem.Repository) *Service {
	return &Service{
		fileSystem: fileSystem,
	}
}

func (s *Service) Execute(srcFile string, outputJSON bool, verbose bool) (string, error) {
	content, err := s.fileSystem.ReadFile(srcFile)
	if err != nil {
		return "", fmt.Errorf("dmesgログファイルの読み込みに失敗しました: %w", err)
	}

	events, err := s.ParseDmesgLog(string(content))
	if err != nil {
		return "", err
	}

	assessment := s.AssessDmesgEvents(events)
	if outputJSON {
		return s.FormatJSON(assessment, verbose)
	}
	return s.FormatText(assessment, verbose), nil
}
