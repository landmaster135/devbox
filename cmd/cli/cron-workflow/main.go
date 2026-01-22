package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "cron-workflow schedules predefined jobs using gocron.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s\n", os.Args[0])
	}
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	if err := run(); err != nil {
		log.Fatalf("cron-workflow: %v", err)
	}
}

func run() error {
	workflows, err := workflow.List()
	if err != nil {
		return err
	}
	if len(workflows) == 0 {
		return errors.New("no workflows configured")
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("failed to initialize scheduler: %w", err)
	}

	if err := registerWorkflows(scheduler, workflows); err != nil {
		return err
	}

	scheduler.Start()
	log.Printf("scheduler started (%d workflow(s)). waiting for termination signal...", len(workflows))

	if err := waitForShutdownSignal(scheduler); err != nil {
		return err
	}

	log.Printf("scheduler stopped cleanly")
	return nil
}

func registerWorkflows(s gocron.Scheduler, workflows []usecases.Workflow) error {
	for i := range workflows {
		wf := workflows[i]
		expression, withSeconds, err := wf.GetCronDefinition()
		if err != nil {
			return err
		}

		task := gocron.NewTask(func(ctx context.Context) {
			if err := wf.Process(ctx); err != nil {
				log.Printf("workflow %q failed: %v", wf.Description, err)
				return
			}
			log.Printf("workflow %q completed", wf.Description)
		})

		if _, err := s.NewJob(
			gocron.CronJob(expression, withSeconds),
			task,
			gocron.WithName(wf.Description),
		); err != nil {
			return fmt.Errorf("failed to register workflow %q: %w", wf.Description, err)
		}

		log.Printf("registered workflow %q (cron=%s)", wf.Description, expression)
	}
	return nil
}

func waitForShutdownSignal(s gocron.Scheduler) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	sig := <-sigCh
	log.Printf("signal received: %s. shutting down...", sig)

	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- s.Shutdown()
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
