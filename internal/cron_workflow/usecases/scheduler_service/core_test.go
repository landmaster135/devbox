package schedulerService

import (
	"context"
	"errors"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	gocronInfra "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/gocron"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
)

func TestRegisterWorkflowsSuccess(t *testing.T) {

	repo := &stubRepository{}
	cs := &recordingCronScheduler{}
	workflows := []usecases.Workflow{
		{
			Description: "wf-no-seconds",
			Frequency:   "*/5 * * * *",
			Timezone:    "UTC",
			Process: func(ctx context.Context) error {
				return nil
			},
		},
		{
			Description: "wf-with-seconds",
			Frequency:   "*/2 * * * * *",
			Timezone:    "Asia/Tokyo",
			Process: func(ctx context.Context) error {
				return nil
			},
		},
	}

	if err := registerWorkflows(repo, cs, workflows); err != nil {
		t.Fatalf("registerWorkflows() error = %v", err)
	}

	if got, want := cs.count(), len(workflows); got != want {
		t.Fatalf("RegisterJob() called %d times, want %d", got, want)
	}

	jobs := cs.snapshot()
	if jobs[0].expression != "CRON_TZ=UTC */5 * * * *" || jobs[0].withSeconds {
		t.Fatalf("first workflow registered incorrectly: %+v", jobs[0])
	}
	if jobs[1].expression != "CRON_TZ=Asia/Tokyo */2 * * * * *" || !jobs[1].withSeconds {
		t.Fatalf("second workflow registered incorrectly: %+v", jobs[1])
	}
}

func TestRegisterWorkflowsCronDefinitionError(t *testing.T) {

	repo := &stubRepository{}
	cs := &recordingCronScheduler{}

	workflows := []usecases.Workflow{
		{
			Description: "invalid",
			Frequency:   " ",
			Timezone:    "UTC",
			Process: func(ctx context.Context) error {
				return nil
			},
		},
	}

	err := registerWorkflows(repo, cs, workflows)
	if err == nil {
		t.Fatalf("registerWorkflows() error = nil, want error")
	}

	if cs.count() != 0 {
		t.Fatalf("RegisterJob() called despite cron definition error")
	}
}

func TestRegisterWorkflowsRegisterJobError(t *testing.T) {

	repo := &stubRepository{}
	cs := &recordingCronScheduler{registerErr: errors.New("boom")}

	workflows := []usecases.Workflow{
		{
			Description: "wf",
			Frequency:   "*/5 * * * *",
			Timezone:    "UTC",
			Process: func(ctx context.Context) error {
				return nil
			},
		},
	}

	err := registerWorkflows(repo, cs, workflows)
	if err == nil {
		t.Fatalf("registerWorkflows() error = nil, want error")
	}
}

func TestNewTaskExecutesWorkflowProcess(t *testing.T) {

	repo := &stubRepository{}
	executed := make(chan struct{}, 1)

	wf := &usecases.Workflow{
		Description: "wf",
		Process: func(ctx context.Context) error {
			select {
			case executed <- struct{}{}:
			default:
			}
			if ctx == nil {
				t.Fatalf("context should not be nil")
			}
			return nil
		},
	}

	task := newTask(repo, wf)
	if task == nil {
		t.Fatalf("newTask() returned nil task")
	}

	process := repo.latestProcess()
	if process == nil {
		t.Fatalf("repo.latestProcess() = nil, want process")
	}

	process(context.Background())

	select {
	case <-executed:
	case <-time.After(time.Second):
		t.Fatalf("workflow process was not executed")
	}
}

func TestNewTaskSwallowsWorkflowErrors(t *testing.T) {

	repo := &stubRepository{}
	called := make(chan struct{}, 1)

	wf := &usecases.Workflow{
		Description: "wf",
		Process: func(ctx context.Context) error {
			called <- struct{}{}
			return errors.New("process failed")
		},
	}

	newTask(repo, wf)
	process := repo.latestProcess()
	if process == nil {
		t.Fatalf("repo.latestProcess() = nil, want process")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		process(context.Background())
	}()

	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatalf("workflow process not invoked")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("process context did not return after error")
	}
}

func TestWaitForShutdownSignalSuccess(t *testing.T) {
	scheduler := newSignalTestScheduler(nil)
	errCh := make(chan error, 1)

	go func() {
		errCh <- waitForShutdownSignal(scheduler)
	}()

	// Give goroutine time to register signal handler.
	time.Sleep(50 * time.Millisecond)
	sendSignal(t, syscall.SIGINT)

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("waitForShutdownSignal() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("waitForShutdownSignal() did not return")
	}

	expectShutdownCalled(t, scheduler.shutdownCalled)
}

func TestWaitForShutdownSignalShutdownError(t *testing.T) {
	scheduler := newSignalTestScheduler(errors.New("shutdown failed"))
	errCh := make(chan error, 1)

	go func() {
		errCh <- waitForShutdownSignal(scheduler)
	}()

	time.Sleep(50 * time.Millisecond)
	sendSignal(t, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatalf("waitForShutdownSignal() error = nil, want error")
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("waitForShutdownSignal() did not return")
	}

	expectShutdownCalled(t, scheduler.shutdownCalled)
}

func TestWaitForShutdownSignalTimeout(t *testing.T) {
	scheduler := newBlockingSignalTestScheduler()
	errCh := make(chan error, 1)
	setSchedulerShutdownTimeout(t, 50*time.Millisecond)

	go func() {
		errCh <- waitForShutdownSignal(scheduler)
	}()

	time.Sleep(50 * time.Millisecond)
	sendSignal(t, syscall.SIGINT)

	select {
	case err := <-errCh:
		if err == nil || err.Error() != "scheduler shutdown timed out" {
			t.Fatalf("waitForShutdownSignal() error = %v, want timeout", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("waitForShutdownSignal() did not return on timeout")
	}

	// Unblock scheduler shutdown goroutine to avoid leaks.
	scheduler.release()
	expectShutdownCalled(t, scheduler.shutdownCalled)
}

func TestScheduleRunsWorkflowAndShutsDown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping scheduler integration test in short mode")
	}

	executed := make(chan struct{}, 1)

	workflow := usecases.Workflow{
		Description: "integration",
		Frequency:   "*/1 * * * * *",
		Timezone:    "UTC",
		Process: func(ctx context.Context) error {
			select {
			case executed <- struct{}{}:
			default:
			}
			return nil
		},
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- Schedule([]usecases.Workflow{workflow})
	}()

	select {
	case <-executed:
	case <-time.After(5 * time.Second):
		t.Fatalf("workflow did not execute")
	}

	sendSignal(t, syscall.SIGINT)

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Schedule() error = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("Schedule() did not exit after signal")
	}
}

// Helpers

type stubRepository struct {
	mu        sync.Mutex
	processes []gocronInfra.ProcessFunc
}

func (s *stubRepository) NewScheduler() (gocronInfra.CronSchedulerRepository, error) {
	return nil, errors.New("not implemented")
}

func (s *stubRepository) NewTask(process gocronInfra.ProcessFunc) *gocronInfra.Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processes = append(s.processes, process)
	return &gocronInfra.Task{}
}

func (s *stubRepository) latestProcess() gocronInfra.ProcessFunc {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.processes) == 0 {
		return nil
	}
	return s.processes[len(s.processes)-1]
}

type registeredJob struct {
	name        string
	expression  string
	withSeconds bool
}

type recordingCronScheduler struct {
	mu          sync.Mutex
	registerErr error
	jobs        []registeredJob
}

func (r *recordingCronScheduler) Start() {}

func (r *recordingCronScheduler) RegisterJob(name, expression string, withSeconds bool, task *gocronInfra.Task) (*gocronInfra.Job, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.registerErr != nil {
		return nil, r.registerErr
	}
	r.jobs = append(r.jobs, registeredJob{name: name, expression: expression, withSeconds: withSeconds})
	return &gocronInfra.Job{}, nil
}

func (r *recordingCronScheduler) Shutdown() error { return nil }

func (r *recordingCronScheduler) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.jobs)
}

func (r *recordingCronScheduler) snapshot() []registeredJob {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]registeredJob, len(r.jobs))
	copy(out, r.jobs)
	return out
}

func newSignalTestScheduler(shutdownErr error) *signalTestScheduler {
	return &signalTestScheduler{
		shutdownErr:    shutdownErr,
		shutdownBlock:  nil,
		shutdownCalled: make(chan struct{}, 1),
	}
}

type signalTestScheduler struct {
	shutdownErr    error
	shutdownBlock  <-chan struct{}
	shutdownCalled chan struct{}
}

func (s *signalTestScheduler) Start() {}

func (s *signalTestScheduler) RegisterJob(name, expression string, withSeconds bool, task *gocronInfra.Task) (*gocronInfra.Job, error) {
	return nil, nil
}

func (s *signalTestScheduler) Shutdown() error {
	select {
	case s.shutdownCalled <- struct{}{}:
	default:
	}
	if s.shutdownBlock != nil {
		<-s.shutdownBlock
	}
	return s.shutdownErr
}

type blockingSignalTestScheduler struct {
	signalTestScheduler
	unblock chan struct{}
}

func newBlockingSignalTestScheduler() *blockingSignalTestScheduler {
	unblock := make(chan struct{})
	return &blockingSignalTestScheduler{
		signalTestScheduler: signalTestScheduler{
			shutdownCalled: make(chan struct{}, 1),
			shutdownBlock:  unblock,
		},
		unblock: unblock,
	}
}

func (b *blockingSignalTestScheduler) release() {
	close(b.unblock)
}

func sendSignal(t *testing.T, sig syscall.Signal) {
	t.Helper()
	if err := syscall.Kill(os.Getpid(), sig); err != nil {
		t.Fatalf("failed to send signal: %v", err)
	}
}

func setSchedulerShutdownTimeout(t *testing.T, d time.Duration) {
	t.Helper()
	prev := schedulerShutdownTimeout
	schedulerShutdownTimeout = d
	t.Cleanup(func() {
		schedulerShutdownTimeout = prev
	})
}

func expectShutdownCalled(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatalf("scheduler Shutdown() was not invoked")
	}
}
