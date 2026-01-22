package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow"
	schedulerService "github.com/landmaster135/devbox/internal/cron_workflow/usecases/scheduler_service"
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

	if err := schedulerService.Schedule(workflows); err != nil {
		return err
	}

	log.Printf("scheduler stopped cleanly")
	return nil
}
