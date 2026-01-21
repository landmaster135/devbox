package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncVolumesIntoCompose(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.yml")
	composePath := filepath.Join(dir, "docker-compose.yml")

	envContent := strings.TrimSpace(`
MOUNT_VOLUME:
  - type: bind
    source: /home/user/cron_output
    target: /app/volume
`) + "\n"

	composeContent := strings.TrimSpace(`
services:
  devbox:
    image: "devbox-cron:latest"
    volumes:
      - type: bind
        source: /tmp
        target: /app/tmp
`) + "\n"

	if err := os.WriteFile(envPath, []byte(envContent), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.WriteFile(composePath, []byte(composeContent), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	service := NewEnvSyncService()
	if err := service.SyncVolumesIntoCompose(envPath, composePath, "MOUNT_VOLUME", "devbox"); err != nil {
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
    volumes:
      - type: bind
        source: /home/user/cron_output
        target: /app/volume
`) + "\n"

	if string(updated) != expected {
		t.Fatalf("compose file mismatch\nexpected:\n%s\nactual:\n%s", expected, string(updated))
	}
}

func TestSyncVolumesIntoComposeRequiresSequence(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.yml")
	composePath := filepath.Join(dir, "docker-compose.yml")

	envContent := strings.TrimSpace(`
MOUNT_VOLUME: invalid
`) + "\n"

	composeContent := strings.TrimSpace(`
services:
  devbox:
    image: devbox
    volumes:
      - type: bind
        source: /tmp
        target: /app/tmp
`) + "\n"

	if err := os.WriteFile(envPath, []byte(envContent), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.WriteFile(composePath, []byte(composeContent), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	service := NewEnvSyncService()
	err := service.SyncVolumesIntoCompose(envPath, composePath, "MOUNT_VOLUME", "devbox")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}
