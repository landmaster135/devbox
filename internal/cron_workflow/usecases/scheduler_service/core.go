package schedulerService

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	gocronInfra "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/gocron"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
	logging "github.com/landmaster135/devbox/internal/logging"
)

var schedulerShutdownTimeout = 15 * time.Second

func Schedule(logger *logging.StructuredLogger, workflows []usecases.Workflow) error {
	repo := gocronInfra.NewRepository()
	scheduleLogger := logging.Ensure(logger)

	scheduler, err := repo.NewScheduler()
	if err != nil {
		return fmt.Errorf("failed to initialize CronScheduler: %w", err)
	}

	if err := registerWorkflows(scheduleLogger, repo, scheduler, workflows); err != nil {
		return err
	}

	scheduler.Start()
	scheduleLogger.WithTags("scheduler").Infof("scheduler started (%d workflow(s)). waiting for termination signal...", len(workflows))

	if err := waitForShutdownSignal(scheduleLogger, scheduler); err != nil {
		return err
	}

	return nil
}

func newTask(logger *logging.StructuredLogger, repo gocronInfra.Repository, wf *usecases.Workflow) *gocronInfra.Task {
	taskLogger := logging.Ensure(logger).WithTags("workflow", wf.Description)
	task := repo.NewTask(func(ctx context.Context) {
		if err := wf.Process(ctx); err != nil {
			taskLogger.WithTags("error").Errorf("failed: %v", err)
			return
		}
		taskLogger.Infof("completed")
	})
	return task
}

func registerWorkflows(logger *logging.StructuredLogger, repo gocronInfra.Repository, cs gocronInfra.CronSchedulerRepository, workflows []usecases.Workflow) error {
	for _, wf := range workflows {
		task := newTask(logger, repo, &wf)

		expression, withSeconds, err := wf.GetCronDefinition()
		if err != nil {
			return fmt.Errorf("failed to get cron definition: %w", err)
		}

		if _, err := cs.RegisterJob(wf.Description, expression, withSeconds, task); err != nil {
			return fmt.Errorf("failed to register job: %w", err)
		}

		logger.WithTags("registry", wf.Description).Infof("registered workflow (cron=%s)", expression)
	}
	return nil
}

func waitForShutdownSignal(logger *logging.StructuredLogger, cs gocronInfra.CronSchedulerRepository) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	sig := <-sigCh
	logging.Ensure(logger).WithTags("scheduler").Infof("signal received: %s. shutting down...", sig)

	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- cs.Shutdown()
	}()

	select {
	case err := <-shutdownDone:
		if err != nil {
			return fmt.Errorf("scheduler shutdown failed: %w", err)
		}
		return nil
	case <-time.After(schedulerShutdownTimeout):
		return errors.New("scheduler shutdown timed out")
	}
}
