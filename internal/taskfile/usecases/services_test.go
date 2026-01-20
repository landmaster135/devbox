package usecases

import (
	"os"
	"path/filepath"
	"testing"
)

func TestServiceInspectRootTaskfileComplete(t *testing.T) {
	t.Parallel()

	referencePath, err := resolveReferencePath(taskTypeRoot)
	if err != nil {
		t.Fatalf("failed to resolve reference path: %v", err)
	}

	data, err := os.ReadFile(referencePath)
	if err != nil {
		t.Fatalf("failed to read reference Taskfile: %v", err)
	}

	tempTarget := writeTempTaskfile(t, string(data))

	service := NewService()
	result, err := service.Inspect(taskTypeRoot, tempTarget)
	if err != nil {
		t.Fatalf("inspect returned error: %v", err)
	}

	if result.HasMissingFields() {
		t.Fatalf("expected no missing fields, but got: %v", result.MissingFields)
	}
}

func TestServiceInspectRootTaskfileMissingFields(t *testing.T) {
	t.Parallel()

	content := `version: "3"

tasks:
  default:
    desc: ""
`

	tempTarget := writeTempTaskfile(t, content)

	service := NewService()
	result, err := service.Inspect(taskTypeRoot, tempTarget)
	if err != nil {
		t.Fatalf("inspect returned error: %v", err)
	}

	if !result.HasMissingFields() {
		t.Fatalf("expected missing fields but got none")
	}

	if !contains(result.MissingFields, "tasks.default.cmds") {
		t.Fatalf("expected tasks.default.cmds to be missing, got: %v", result.MissingFields)
	}

	if !contains(result.MissingFields, "tasks.alias") {
		t.Fatalf("expected tasks.alias to be missing, got: %v", result.MissingFields)
	}
}

func TestServiceInspectUnknownTaskType(t *testing.T) {
	t.Parallel()

	service := NewService()
	if _, err := service.Inspect("unknown", "Taskfile.yml"); err == nil {
		t.Fatalf("expected error for unsupported task type")
	}
}

func writeTempTaskfile(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "Taskfile.yml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp Taskfile: %v", err)
	}

	return path
}

func contains(collection []string, value string) bool {
	for _, v := range collection {
		if v == value {
			return true
		}
	}
	return false
}
