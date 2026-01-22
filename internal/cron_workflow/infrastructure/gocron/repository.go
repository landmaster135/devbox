package gocron

import "context"

// Repository abstracts task scheduler creation backed by go-co-op/gocron.
type Repository interface {
	NewScheduler() (CronSchedulerRepository, error)
	NewTask(process ProcessFunc) *Task
}

// CronSchedulerRepository exposes operations required by the use case layer.
type CronSchedulerRepository interface {
	Start()
	RegisterJob(name, expression string, withSeconds bool, task *Task) (*Job, error)
	Shutdown() error
}

// ProcessFunc is the signature executed by a scheduled task.
type ProcessFunc func(ctx context.Context)

// NewRepository returns the default implementation backed by impl.go helpers.
func NewRepository() Repository {
	return &repository{}
}

type repository struct{}

func (repository) NewScheduler() (CronSchedulerRepository, error) {
	return NewCronScheduler()
}

func (repository) NewTask(process ProcessFunc) *Task {
	return NewTask(process)
}
