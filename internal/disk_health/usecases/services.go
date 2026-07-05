package usecases

import (
	"fmt"

	config "github.com/landmaster135/devbox/internal/disk_health/config"
	"github.com/landmaster135/devbox/internal/disk_health/infrastructures/filesystem"
	assessdmesg "github.com/landmaster135/devbox/internal/disk_health/usecases/operations/assess_dmesg"
	assesssmart "github.com/landmaster135/devbox/internal/disk_health/usecases/operations/assess_smart"
)

type Service struct {
	fileSystem           filesystem.Repository
	assessSmartOperation assessSmartOperation
	assessDmesgOperation assessDmesgOperation
}

type ServiceOptions struct {
	FileSystem filesystem.Repository
}

func NewService(options ServiceOptions) *Service {
	fileSystem := options.FileSystem
	if fileSystem == nil {
		fileSystem = filesystem.NewOSRepository()
	}

	service := &Service{fileSystem: fileSystem}
	service.assessSmartOperation = assesssmart.NewService(fileSystem)
	service.assessDmesgOperation = assessdmesg.NewService(fileSystem)
	return service
}

func (s *Service) AssessSmart(srcFile string, outputJSON bool, verbose bool) (string, error) {
	return s.assessSmartOperation.Execute(srcFile, outputJSON, verbose)
}

func (s *Service) AssessDmesg(srcFile string, outputJSON bool, verbose bool) (string, error) {
	return s.assessDmesgOperation.Execute(srcFile, outputJSON, verbose)
}

func (s *Service) ExecuteByConfig(cfg *config.Config) (string, error) {
	switch cfg.Operation {
	case config.OperationAssessSmart:
		return s.AssessSmart(cfg.SrcFile, cfg.JSON, cfg.Verbose)
	case config.OperationAssessDmesg:
		return s.AssessDmesg(cfg.SrcFile, cfg.JSON, cfg.Verbose)
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", cfg.Operation)
	}
}
