package schedulerService

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	gocronInfra "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/gocron"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
)

func Schedule(workflows []usecases.Workflow) error {
	scheduler, err := gocronInfra.NewCronScheduler()
	if err != nil {
		return fmt.Errorf("failed to initialize CronScheduler: %w", err)
	}

	if err := registerWorkflows(scheduler, workflows); err != nil {
		return err
	}

	scheduler.Start()
	log.Printf("scheduler started (%d workflow(s)). waiting for termination signal...", len(workflows))

	if err := waitForShutdownSignal(scheduler); err != nil {
		return err
	}

	return nil
}

func newTask(wf *usecases.Workflow) *gocronInfra.Task {
	task := gocronInfra.NewTask(func(ctx context.Context) {
		if err := wf.Process(ctx); err != nil {
			log.Printf("workflow %q failed: %v", wf.Description, err)
			return
		}
		log.Printf("workflow %q completed", wf.Description)
	})
	return task
}

func registerWorkflows(cs *gocronInfra.CronScheduler, workflows []usecases.Workflow) error {
	for _, wf := range workflows {
		task := newTask(&wf)

		expression, withSeconds, err := wf.GetCronDefinition()
		if err != nil {
			return fmt.Errorf("failed to get cron definition: %w", err)
		}

		if _, err := cs.RegisterJob(wf.Description, expression, withSeconds, task); err != nil {
			return fmt.Errorf("failed to register job: %w", err)
		}
	}
	return nil
}

func waitForShutdownSignal(cs *gocronInfra.CronScheduler) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	sig := <-sigCh
	log.Printf("signal received: %s. shutting down...", sig)

	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- cs.Scheduler.Shutdown()
	}()

	select {
	case err := <-shutdownDone:
		if err != nil {
			return fmt.Errorf("scheduler shutdown failed: %w", err)
		}
		return nil
	case <-time.After(15 * time.Second):
		return errors.New("scheduler shutdown timed out")
	}
}
