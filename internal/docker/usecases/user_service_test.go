package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncUserIntoComposeUpdatesExistingField(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.yml")
	composePath := filepath.Join(dir, "docker-compose.yml")

	envContent := "USER_VALUE: \"8888:8888\"\n"
	composeContent := strings.TrimSpace(`
services:
  devbox:
    image: "devbox-cron:latest"
    user: "1000:1000"
`) + "\n"

	if err := os.WriteFile(envPath, []byte(envContent), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.WriteFile(composePath, []byte(composeContent), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	service := NewEnvSyncService()
	if err := service.SyncUserIntoCompose(envPath, composePath, "USER_VALUE", "devbox"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}

	expected := strings.TrimSpace(`
services:
  devbox:
    image: "devbox-cron:latest"
    user: "8888:8888"
`) + "\n"

	if string(updated) != expected {
		t.Fatalf("compose file mismatch\nexpected:\n%s\nactual:\n%s", expected, string(updated))
	}
}

func TestSyncUserIntoComposeInsertsWhenMissing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.yml")
	composePath := filepath.Join(dir, "docker-compose.yml")

	envContent := "USER_VALUE: 7777\n"
	composeContent := strings.TrimSpace(`
services:
  devbox:
    image: devbox
    environment:
      - GIN_MODE=staging
`) + "\n"

	if err := os.WriteFile(envPath, []byte(envContent), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.WriteFile(composePath, []byte(composeContent), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	service := NewEnvSyncService()
	if err := service.SyncUserIntoCompose(envPath, composePath, "USER_VALUE", "devbox"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}

	expected := strings.TrimSpace(`
services:
  devbox:
    image: devbox
    environment:
      - GIN_MODE=staging
    user: "7777"
`) + "\n"

	if string(updated) != expected {
		t.Fatalf("compose file mismatch\nexpected:\n%s\nactual:\n%s", expected, string(updated))
	}
}

func TestSyncUserIntoComposeRequiresScalar(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.yml")
	composePath := filepath.Join(dir, "docker-compose.yml")

	envContent := strings.TrimSpace(`
USER_VALUE:
  foo: bar
`) + "\n"
	composeContent := strings.TrimSpace(`
services:
  devbox:
    image: devbox
`) + "\n"

	if err := os.WriteFile(envPath, []byte(envContent), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.WriteFile(composePath, []byte(composeContent), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	service := NewEnvSyncService()
	err := service.SyncUserIntoCompose(envPath, composePath, "USER_VALUE", "devbox")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
