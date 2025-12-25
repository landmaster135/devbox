package usecases

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/landmaster135/devbox/internal/shell/domain"
)

type mockCommandExecutor struct {
	output      *ExecutionOutput
	err         error
	lastCommand []string
	lastWorkDir string
	lastEnv     map[string]string
}

func (m *mockCommandExecutor) Execute(ctx context.Context, command []string, workDir string, env map[string]string) (*ExecutionOutput, error) {
	m.lastCommand = append([]string(nil), command...)
	m.lastWorkDir = workDir
	if len(env) > 0 {
		m.lastEnv = make(map[string]string, len(env))
		for k, v := range env {
			m.lastEnv[k] = v
		}
	}
	if m.err != nil {
		return nil, m.err
	}
	if m.output != nil {
		return m.output, nil
	}
	return &ExecutionOutput{}, nil
}

func TestShellServiceExecuteCommandSuccess(t *testing.T) {
	mock := &mockCommandExecutor{
		output: &ExecutionOutput{
			Stdout:   "done",
			Stderr:   "",
			ExitCode: 0,
		},
	}
	service := NewShellServiceWithExecutor(mock)

	baseDir := t.TempDir()
	workDir := filepath.Join(baseDir, "workspace")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatalf("failed to create workdir: %v", err)
	}

	input := &ExecuteCommandInput{
		Command:   []string{"bash", "-lc", "echo shell"},
		WorkDir:   "workspace",
		BaseDir:   baseDir,
		TimeoutMs: 1200,
		Env: map[string]string{
			"FOO": "bar",
		},
		SandboxPermissions: domain.SandboxPermissionsUseDefault,
	}

	result, err := service.ExecuteCommand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(mock.lastCommand, input.Command) {
		t.Fatalf("command mismatch: %+v vs %+v", mock.lastCommand, input.Command)
	}
	if mock.lastWorkDir != workDir {
		t.Fatalf("workdir mismatch: %s vs %s", mock.lastWorkDir, workDir)
	}
	if !result.Success {
		t.Fatalf("expected success result: %+v", result)
	}
	if result.Stdout != "done" {
		t.Fatalf("unexpected stdout: %s", result.Stdout)
	}
	if result.TimeoutMs != 1200 {
		t.Fatalf("timeout should reflect requested value, got %d", result.TimeoutMs)
	}
	if result.WorkDir != workDir {
		t.Fatalf("workdir should be absolute, got %s", result.WorkDir)
	}
	if result.BaseDir != baseDir {
		t.Fatalf("base dir mismatch: %s vs %s", result.BaseDir, baseDir)
	}
	if result.EscalationRequested {
		t.Fatalf("escalation should be false")
	}
}

func TestShellServiceExecuteCommandRequireJustification(t *testing.T) {
	service := NewShellServiceWithExecutor(&mockCommandExecutor{})

	input := &ExecuteCommandInput{
		Command:            []string{"echo", "hi"},
		BaseDir:            t.TempDir(),
		SandboxPermissions: domain.SandboxPermissionsRequireEscalated,
		Justification:      "",
	}

	if _, err := service.ExecuteCommand(input); err == nil {
		t.Fatalf("expected error for missing justification")
	}
}

func TestShellServiceExecuteCommandRejectsOutsideWorkdir(t *testing.T) {
	service := NewShellServiceWithExecutor(&mockCommandExecutor{})
	baseDir := t.TempDir()

	input := &ExecuteCommandInput{
		Command:            []string{"echo", "hi"},
		WorkDir:            "../outside",
		BaseDir:            baseDir,
		SandboxPermissions: domain.SandboxPermissionsUseDefault,
	}

	if _, err := service.ExecuteCommand(input); err == nil {
		t.Fatalf("expected error for workdir outside base")
	}
}

func TestShellServiceExecuteCommandHandlesTimeout(t *testing.T) {
	mock := &mockCommandExecutor{
		output: &ExecutionOutput{
			Stdout:   "partial",
			Stderr:   "",
			ExitCode: -1,
			TimedOut: true,
		},
	}
	service := NewShellServiceWithExecutor(mock)
	baseDir := t.TempDir()

	input := &ExecuteCommandInput{
		Command:            []string{"echo", "hi"},
		BaseDir:            baseDir,
		SandboxPermissions: domain.SandboxPermissionsUseDefault,
	}

	result, err := service.ExecuteCommand(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.TimedOut {
		t.Fatalf("expected timeout flag")
	}
	if result.Success {
		t.Fatalf("timeout should not be success")
	}
	if result.ExitCode != -1 {
		t.Fatalf("timeout exit code should be -1, got %d", result.ExitCode)
	}
}

func TestShellServiceListAllowedCommandsSorted(t *testing.T) {
	service := NewShellService()
	commands := service.ListAllowedCommands()

	if len(commands) == 0 {
		t.Fatalf("commands should not be empty")
	}

	sorted := append([]string(nil), commands...)
	sort.Strings(sorted)
	if !reflect.DeepEqual(commands, sorted) {
		t.Fatalf("commands must be sorted")
	}
}
