package config

import (
	"os"
	"testing"
)

func TestParseFlags_OllamaSuccess(t *testing.T) {
	t.Setenv("GOFLAGS", "")
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{
		"vector-embedding",
		"-operation", "ollama",
		"-host", "localhost",
		"-port", "1234",
		"-model", "test-model",
		"-input", "foo",
		"-input", "bar",
		"-timeout", "90",
	}

	cfg, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}

	if cfg.Operation != OperationOllama {
		t.Fatalf("unexpected operation: %s", cfg.Operation)
	}
	if cfg.Host != "localhost" {
		t.Fatalf("unexpected host: %s", cfg.Host)
	}
	if cfg.Port != 1234 {
		t.Fatalf("unexpected port: %d", cfg.Port)
	}
	if cfg.Model != "test-model" {
		t.Fatalf("unexpected model: %s", cfg.Model)
	}
	if len(cfg.Inputs) != 2 {
		t.Fatalf("unexpected inputs: %#v", cfg.Inputs)
	}
	if cfg.TimeoutSeconds != 90 {
		t.Fatalf("unexpected timeout: %d", cfg.TimeoutSeconds)
	}
}

func TestParseFlags_ValidationError(t *testing.T) {
	t.Setenv("GOFLAGS", "")
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"vector-embedding", "-operation", "ollama", "-model", "", "-input", ""}

	if _, err := ParseFlags(); err == nil {
		t.Fatal("expected error but got nil")
	}
}
