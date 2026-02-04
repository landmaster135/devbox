package gocron

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
)

type Task struct {
	Task gocron.Task
}

func NewTask(process ProcessFunc) *Task {
	task := gocron.NewTask(process)
	return &Task{
		Task: task,
	}
}

type Job struct {
	Job gocron.Job
}

type CronScheduler struct {
	Scheduler gocron.Scheduler
}

func NewCronScheduler() (*CronScheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize scheduler: %w", err)
	}
	return &CronScheduler{
		Scheduler: scheduler,
	}, nil
}

func (cs *CronScheduler) Start() {
	cs.Scheduler.Start()
}

func (cs *CronScheduler) RegisterJob(name, expression string, withSeconds bool, task *Task) (*Job, error) {
	j, err := cs.Scheduler.NewJob(
		gocron.CronJob(expression, withSeconds),
		task.Task,
		gocron.WithName(name),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register workflow %q: %w", name, err)
	}

	return &Job{
		Job: j,
	}, nil
}

func (cs *CronScheduler) Shutdown() error {
	return cs.Scheduler.Shutdown()
}
