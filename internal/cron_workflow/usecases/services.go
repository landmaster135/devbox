package usecases

import (
	"context"
	"fmt"
	"strings"
)

// ProcessFunc defines the signature each workflow task must satisfy.
type ProcessFunc func(ctx context.Context) error

// Workflow describes a single scheduled job definition.
type Workflow struct {
	Description string
	Frequency   string
	Timezone    string
	Process     ProcessFunc
}

// GetCronDefinition returns the cron expression with timezone information and
// indicates whether the schedule uses a seconds field.
func (w Workflow) GetCronDefinition() (string, bool, error) {
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

	expression := fmt.Sprintf("CRON_TZ=%s %s", w.Timezone, base)

	withSeconds := false
	if len(fields) == 6 {
		withSeconds = true
	}

	return expression, withSeconds, nil
}
