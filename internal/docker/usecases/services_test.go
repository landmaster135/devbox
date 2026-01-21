package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSyncEnvIntoCompose(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	envPath := filepath.Join(dir, "env.yml")
	composePath := filepath.Join(dir, "docker-compose.yml")

	envContent := strings.TrimSpace(`
VITE_GIN_MODE: production
SIMPLE_VALUE: simple
VALUE_WITH_SPACE: "value with space"
`) + "\n"
	composeContent := strings.TrimSpace(`
services:
  dathub:
    image: "dathub-frontend:latest"
    environment: &dathub-env
      - OLD=value
    labels:
      sample.label: "true"

  dathub-backend:
    image: "dathub-backend:latest"
    environment: *dathub-env
`) + "\n"

	if err := os.WriteFile(envPath, []byte(envContent), 0o644); err != nil {
		t.Fatalf("failed to write env file: %v", err)
	}
	if err := os.WriteFile(composePath, []byte(composeContent), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	service := NewEnvSyncService()
	count, err := service.SyncEnvIntoCompose(envPath, composePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 env entries, got %d", count)
	}

	updated, err := os.ReadFile(composePath)
	if err != nil {
		t.Fatalf("failed to read compose file: %v", err)
	}

	expected := strings.TrimSpace(`
services:
  dathub:
    image: "dathub-frontend:latest"
    environment: &dathub-env
      - VITE_GIN_MODE=production
      - SIMPLE_VALUE=simple
      - VALUE_WITH_SPACE="value with space"
    labels:
      sample.label: "true"

  dathub-backend:
    image: "dathub-backend:latest"
    environment: *dathub-env
`) + "\n"

	if string(updated) != expected {
		t.Fatalf("compose file mismatch\nexpected:\n%s\nactual:\n%s", expected, string(updated))
	}
}


func TestParseEnvEntries(t *testing.T) {
	t.Parallel()

	content := strings.TrimSpace(`
# comment line
VITE_FOO: foo # inline
BAR: "bar #ignored"
BAZ: 'value with # hash'
EMPTY:
INVALID-LINE: should_skip
`) + "\n"

	entries, err := parseEnvEntries(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	cases := []struct {
		idx   int
		key   string
		value string
	}{
		{0, "VITE_FOO", "foo"},
		{1, "BAR", "bar #ignored"},
		{2, "BAZ", "value with # hash"},
		{3, "EMPTY", ""},
	}

	for _, c := range cases {
		if entries[c.idx].Key != c.key {
			t.Fatalf("entry %d key mismatch: expected %s got %s", c.idx, c.key, entries[c.idx].Key)
		}
		if entries[c.idx].Value != c.value {
			t.Fatalf("entry %d value mismatch: expected %s got %s", c.idx, c.value, entries[c.idx].Value)
		}
	}
}

func TestInjectEnvironmentBlockErrorsWithoutSection(t *testing.T) {
	t.Parallel()

	_, err := injectEnvironmentBlock("services:\n  foo:\n    image: test\n", []EnvEntry{{Key: "KEY", Value: "VALUE"}})
	if err == nil {
		t.Fatal("expected error when environment section is missing")
	}
}

func TestInjectEnvironmentBlockSkipsAnchorReference(t *testing.T) {
	t.Parallel()

	compose := strings.TrimSpace(`
services:
  foo:
    environment: *shared-env
  bar:
    environment: &shared-env
      - SAMPLE=1
`) + "\n"

	updated, err := injectEnvironmentBlock(compose, []EnvEntry{{Key: "KEY", Value: "VALUE"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Count(updated, "KEY=VALUE") != 1 {
		t.Fatalf("expected environment block to be updated only once, got:\n%s", updated)
	}
}
