package workflow

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	workflowenv "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow/infrastructure/env"
)

const defaultTimezone = "Asia/Tokyo"

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
	heartOwner, err := envRepo.GetEnv("HEART_OWNER")
	if err != nil {
		return nil, fmt.Errorf("resolve heart owner from %s: %w", "HEART_OWNER", err)
	}

	return []Workflow{
		heartbeatWorkflow(heartOwner),
		hourlyHealthSnapshotWorkflow(),
	}, nil
}

func heartbeatWorkflow(owner string) Workflow {
	return Workflow{
		Description: "Heartbeat monitor (every minute)",
		Frequency:   "*/1 * * * *",
		Process: func(ctx context.Context) error {
			log.Printf("[heartbeat] alive: %s (owner=%s)", time.Now().Format(time.RFC3339), owner)
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
