package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	workflow "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow"
	schedulerService "github.com/landmaster135/devbox/internal/cron_workflow/usecases/scheduler_service"
	logging "github.com/landmaster135/devbox/internal/logging"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "cron-workflow schedules predefined jobs using gocron.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s\n", os.Args[0])
	}
	flag.Parse()

	logger := logging.New(logging.WithWriter(os.Stdout)).WithTags("CRON workflow")
	if err := run(logger); err != nil {
		logger.WithTags("fatal").Errorf("cron-workflow: %v", err)
		os.Exit(1)
	}
}

func run(logger *logging.StructuredLogger) error {
	workflows, err := workflow.List(logger)
	if err != nil {
		return err
	}
	if len(workflows) == 0 {
		return errors.New("no workflows configured")
	}

	if err := schedulerService.Schedule(logger, workflows); err != nil {
		return err
	}

	logger.WithTags("scheduler").Infof("scheduler stopped cleanly")
	return nil
}
