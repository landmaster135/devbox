package workflowCreator

import (
	"fmt"
	"path/filepath"

	infraEnv "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/env"
	infraFilesystem "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/filesystem"
	infraTime "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/time"
)

type WorkflowCreator struct {
	Timezone  string
	EnvRepo   infraEnv.Repository
	FileRepo  infraFilesystem.Repository
	TimeRepo  infraTime.Repository
	VolumeDir string
}

func workflowVolumeDir(fs infraFilesystem.Repository) (string, error) {
	wd, err := fs.WorkingDir()
	if err != nil {
		return "", fmt.Errorf("determine workflow volume directory: %w", err)
	}
	return filepath.Join(wd, "volume"), nil
}

func NewWorkflowCreator(timezone string, envRepo infraEnv.Repository, fileRepo infraFilesystem.Repository, timeRepo infraTime.Repository) (*WorkflowCreator, error) {
	volumeDir, err := workflowVolumeDir(fileRepo)
	if err != nil {
		return nil, fmt.Errorf("resolve workflow volume directory: %w", err)
	}
	if err := fileRepo.EnsureDir(volumeDir); err != nil {
		return nil, fmt.Errorf("prepare workflow volume directory: %w", err)
	}

	return &WorkflowCreator{
		Timezone:  timezone,
		EnvRepo:   envRepo,
		FileRepo:  fileRepo,
		TimeRepo:  timeRepo,
		VolumeDir: volumeDir,
	}, nil
}
