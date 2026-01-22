package workflow

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	infraEnv "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/env"
	filesystem "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/filesystem"
	InfraTime "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/time"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
)

func workflowVolumeDir(fs filesystem.Repository) (string, error) {
	wd, err := fs.WorkingDir()
	if err != nil {
		return "", fmt.Errorf("determine workflow volume directory: %w", err)
	}
	return filepath.Join(wd, "volume"), nil
}

// List returns all configured workflows.
func List() ([]usecases.Workflow, error) {
	envRepo := infraEnv.NewRepository()
	filesystemRepo := filesystem.NewRepository()

	volumeDir, err := workflowVolumeDir(filesystemRepo)
	if err != nil {
		return nil, fmt.Errorf("resolve workflow volume directory: %w", err)
	}
	if err := filesystemRepo.EnsureDir(volumeDir); err != nil {
		return nil, fmt.Errorf("prepare workflow volume directory: %w", err)
	}

	heartOwner, err := envRepo.GetEnv("HEART_OWNER")
	if err != nil {
		return nil, fmt.Errorf("resolve heart owner from %s: %w", "HEART_OWNER", err)
	}

	const tz = "Asia/Tokyo"

	return []usecases.Workflow{
		heartbeatWorkflow(tz, heartOwner, filesystemRepo, volumeDir),
		hourlyHealthSnapshotWorkflow(tz),
	}, nil
}

func heartbeatWorkflow(timezone, owner string, fileRepo filesystem.Repository, volumeDir string) usecases.Workflow {
	return usecases.Workflow{
		Description: "Heartbeat monitor (every minute)",
		Frequency:   "*/1 * * * *",
		Timezone:    timezone,
		Process: func(ctx context.Context) error {
			timeRepo := InfraTime.NewRepository()

			now, err := timeRepo.Now(timezone)
			if err != nil {
				return fmt.Errorf("resolve current time: %w", err)
			}
			timestamp := now.Format("20060102150405")
			statusFile := filepath.Join(volumeDir, fmt.Sprintf("heartbeat-%s.status", timestamp))

			message := fmt.Sprintf("[heartbeat] alive: %s (owner=%s)", time.Now().Format(time.RFC3339), owner)
			log.Printf("%s", message)
			log.Printf("[heartbeat] writing status file: %s", statusFile)
			if err := fileRepo.Write(statusFile, true, message+"\n"); err != nil {
				return fmt.Errorf("write heartbeat status file: %w", err)
			}
			return nil
		},
	}
}

func hourlyHealthSnapshotWorkflow(timezone string) usecases.Workflow {
	return usecases.Workflow{
		Description: "Hourly state snapshot",
		Frequency:   "15 * * * *",
		Timezone:    timezone,
		Process: func(ctx context.Context) error {
			select {
			case <-time.After(2 * time.Second):
				log.Printf("[snapshot] captured at UTC=%s", time.Now().UTC().Format(time.RFC3339))
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}
}
