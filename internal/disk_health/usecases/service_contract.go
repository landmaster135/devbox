package usecases

import "github.com/landmaster135/devbox/internal/disk_health/infrastructures/filesystem"

type assessSmartOperation interface {
	Execute(srcFile string, outputJSON bool, verbose bool) (string, error)
}

type assessDmesgOperation interface {
	Execute(srcFile string, outputJSON bool, verbose bool) (string, error)
}

func newServiceWithOperations(
	fileSystem filesystem.Repository,
	assessSmartOp assessSmartOperation,
	assessDmesgOp assessDmesgOperation,
) *Service {
	return &Service{
		fileSystem:           fileSystem,
		assessSmartOperation: assessSmartOp,
		assessDmesgOperation: assessDmesgOp,
	}
}
