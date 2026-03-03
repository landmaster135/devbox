package main

import (
	"bytes"
	"errors"
	"testing"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
)

type mockExecutor struct {
	result string
	err    error
}

func (m *mockExecutor) ExecuteByConfig(cfg *config.Config) (string, error) {
	return m.result, m.err
}

type errWriter struct{}

func (w *errWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write error")
}

func TestExecute_Normal(t *testing.T) {
	buffer := &bytes.Buffer{}
	executor := &mockExecutor{result: "ok\n"}

	err := execute(&config.Config{Operation: config.OperationAdd}, executor, buffer)
	if err != nil {
		t.Fatalf("execute returned error: %v", err)
	}
	if buffer.String() != "ok\n" {
		t.Fatalf("unexpected output: %q", buffer.String())
	}
}

func TestExecute_ExecutorError(t *testing.T) {
	executor := &mockExecutor{err: errors.New("execute error")}

	err := execute(&config.Config{Operation: config.OperationAdd}, executor, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "execute error" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecute_WriteError(t *testing.T) {
	executor := &mockExecutor{result: "ok\n"}

	err := execute(&config.Config{Operation: config.OperationAdd}, executor, &errWriter{})
	if err == nil {
		t.Fatal("expected write error")
	}
}
