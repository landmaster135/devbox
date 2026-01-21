package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncPortsIntoCompose(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.yml")
	composePath := filepath.Join(dir, "docker-compose.yml")

	envContent := strings.TrimSpace(`
VITE_FRONT_URL_PORT: 3333
VITE_API_BASE_URL_PORT: 4444
`) + "\n"
	composeContent := strings.TrimSpace(`
services:
  dathub:
    image: "dathub-frontend:latest"
    ports:
      - "9999:9999"
    labels:
      tsdproxy.enable: "true"
      tsdproxy.container_port: 9999

  api:
    image: "dathub-backend:latest"
    ports:
      - "8888:8888"
`) + "\n"

	if err := os.WriteFile(envPath, []byte(envContent), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.WriteFile(composePath, []byte(composeContent), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	service := NewEnvSyncService()
	if err := service.SyncPortsIntoCompose(envPath, composePath, "VITE_FRONT_URL_PORT", "dathub"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}

	expected := strings.TrimSpace(`
services:
  dathub:
    image: "dathub-frontend:latest"
    ports:
      - "3333:3333"
    labels:
      tsdproxy.enable: "true"
      tsdproxy.container_port: 3333

  api:
    image: "dathub-backend:latest"
    ports:
      - "8888:8888"
`) + "\n"

	if string(updated) != expected {
		t.Fatalf("ports sync mismatch\nexpected:\n%s\nactual:\n%s", expected, string(updated))
	}
}

func TestSyncPortsIntoComposeMissingKey(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.yml")
	composePath := filepath.Join(dir, "docker-compose.yml")

	if err := os.WriteFile(envPath, []byte("OTHER: 1\n"), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.WriteFile(composePath, []byte("services:\n  sample:\n    ports:\n      - \"1:1\"\n"), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	service := NewEnvSyncService()
	err := service.SyncPortsIntoCompose(envPath, composePath, "MISSING", "sample")
	if err == nil {
		t.Fatal("expected error when port key is missing")
	}
}
