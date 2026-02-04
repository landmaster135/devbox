package gocron

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestNewCronSchedulerStartAndShutdown(t *testing.T) {
	cs, cleanup := newTestScheduler(t)
	defer cleanup()

	if cs.Scheduler == nil {
		t.Fatalf("NewCronScheduler() returned nil Scheduler")
	}

	cs.Start()

	if jobs := cs.Scheduler.Jobs(); jobs == nil {
		t.Fatalf("Scheduler.Jobs() returned nil slice")
	}
}

func TestCronSchedulerRegisterJobRunsTask(t *testing.T) {
	cs, cleanup := newTestScheduler(t)
	defer cleanup()

	triggered := make(chan context.Context, 1)
	task := NewTask(func(ctx context.Context) {
		select {
		case triggered <- ctx:
		default:
		}
	})

	if _, err := cs.RegisterJob("workflow", "*/1 * * * * *", true, task); err != nil {
		t.Fatalf("RegisterJob() error = %v", err)
	}

	cs.Start()

	select {
	case ctx := <-triggered:
		if ctx == nil {
			t.Fatalf("RegisterJob() delivered nil context")
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("RegisterJob() did not execute task within timeout")
	}
}

func TestCronSchedulerRegisterJobRejectsInvalidExpression(t *testing.T) {
	cs, cleanup := newTestScheduler(t)
	defer cleanup()

	task := NewTask(func(ctx context.Context) {})

	_, err := cs.RegisterJob("invalid", "* * *", false, task)
	if err == nil {
		t.Fatalf("RegisterJob() error = nil, want non-nil for invalid expression")
	}

	if !strings.Contains(err.Error(), "failed to register workflow \"invalid\"") {
		t.Fatalf("RegisterJob() error = %v, want wrapped message", err)
	}
}

func TestRepositoryCreatesFunctionalSchedulerAndTask(t *testing.T) {
	repo := NewRepository()

	scheduler, err := repo.NewScheduler()
	if err != nil {
		t.Fatalf("NewScheduler() error = %v", err)
	}

	defer func() {
		if err := scheduler.Shutdown(); err != nil {
			t.Fatalf("Shutdown() error = %v", err)
		}
	}()

	executed := make(chan struct{}, 1)
	task := repo.NewTask(func(ctx context.Context) {
		select {
		case <-ctx.Done():
		default:
		}
		select {
		case executed <- struct{}{}:
		default:
		}
	})

	if _, err := scheduler.RegisterJob("repo-workflow", "*/1 * * * * *", true, task); err != nil {
		t.Fatalf("RegisterJob() error = %v", err)
	}

	scheduler.Start()

	select {
	case <-executed:
	case <-time.After(3 * time.Second):
		t.Fatalf("repository-backed scheduler did not execute task")
	}
}

func newTestScheduler(t *testing.T) (*CronScheduler, func()) {
	t.Helper()

	cs, err := NewCronScheduler()
	if err != nil {
		t.Fatalf("NewCronScheduler() error = %v", err)
	}

	cleanup := func() {
		if err := cs.Shutdown(); err != nil {
			t.Fatalf("Shutdown() error = %v", err)
		}
	}

	return cs, cleanup
}
