package workflow

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	workflowenv "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow/infrastructure/env"
	filesystem "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow/infrastructure/filesystem"
	workflowtime "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow/infrastructure/time"
)

const defaultTimezone = "Asia/Tokyo"

var workflowVolumePath string

func workflowVolumeDir() (string, error) {
	if workflowVolumePath != "" {
		return workflowVolumePath, nil
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("determine workflow volume directory: runtime.Caller failed")
	}
	workflowVolumePath = filepath.Join(filepath.Dir(file), "volume")
	return workflowVolumePath, nil
}

// ProcessFunc defines the signature each workflow task must satisfy.
type ProcessFunc func(ctx context.Context) error

// Workflow describes a single scheduled job definition.
type Workflow struct {
	Description string
	Frequency   string
	Timezone    string
	Process     ProcessFunc
}

// CronDefinition returns the cron expression with timezone information and
// indicates whether the schedule uses a seconds field.
func (w Workflow) CronDefinition() (string, bool, error) {
	if strings.TrimSpace(w.Frequency) == "" {
		return "", false, fmt.Errorf("frequency is not set: %s", w.Description)
	}
	if w.Process == nil {
		return "", false, fmt.Errorf("process is not set: %s", w.Description)
	}

	base := strings.TrimSpace(w.Frequency)
	fields := strings.Fields(base)
	if len(fields) == 0 {
		return "", false, fmt.Errorf("invalid frequency: %s", w.Description)
	}

	tz := w.timezoneOrDefault()
	expression := fmt.Sprintf("CRON_TZ=%s %s", tz, base)

	withSeconds := false
	if len(fields) == 6 {
		withSeconds = true
	}

	return expression, withSeconds, nil
}

func (w Workflow) timezoneOrDefault() string {
	if strings.TrimSpace(w.Timezone) == "" {
		return defaultTimezone
	}
	return w.Timezone
}

// List returns all configured workflows.
func List() ([]Workflow, error) {
	envRepo := workflowenv.NewRepository()
	filesystemRepo := filesystem.NewRepository()

	volumeDir, err := workflowVolumeDir()
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

	return []Workflow{
		heartbeatWorkflow(heartOwner, filesystemRepo, volumeDir),
		hourlyHealthSnapshotWorkflow(),
	}, nil
}

func heartbeatWorkflow(owner string, fileRepo filesystem.Repository, volumeDir string) Workflow {
	return Workflow{
		Description: "Heartbeat monitor (every minute)",
		Frequency:   "*/1 * * * *",
		Process: func(ctx context.Context) error {
			timeRepo := workflowtime.NewRepository()

			now, err := timeRepo.Now(defaultTimezone)
			if err != nil {
				return fmt.Errorf("resolve current time: %w", err)
			}
			timestamp := now.Format("20060102150405")
			statusFile := filepath.Join(volumeDir, fmt.Sprintf("heartbeat-%s.status", timestamp))

			message := fmt.Sprintf("[heartbeat] alive: %s (owner=%s)", time.Now().Format(time.RFC3339), owner)
			log.Printf("%s", message)
			if err := fileRepo.Write(statusFile, true, message+"\n"); err != nil {
				return fmt.Errorf("write heartbeat status file: %w", err)
			}
			return nil
		},
	}
}

func hourlyHealthSnapshotWorkflow() Workflow {
	return Workflow{
		Description: "Hourly state snapshot",
		Frequency:   "15 * * * *",
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
