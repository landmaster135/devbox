package usecases

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
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

func TestServiceFillPopulatesBlankValues(t *testing.T) {
	t.Parallel()

	content := `version: "3"

tasks:
  default:
    desc: ""
    cmds: []
  alias:
    desc: ""
`

	tempTarget := writeTempTaskfile(t, content)
	service := NewService()
	updated, err := service.Fill(taskTypeRoot, tempTarget)
	if err != nil {
		t.Fatalf("fill returned error: %v", err)
	}
	if !updated {
		t.Fatalf("expected fill to report updates")
	}

	var parsed struct {
		Tasks map[string]map[string]interface{} `yaml:"tasks"`
	}
	data, err := os.ReadFile(tempTarget)
	if err != nil {
		t.Fatalf("failed to read filled Taskfile: %v", err)
	}
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to parse filled Taskfile: %v", err)
	}

	defaultTask, ok := parsed.Tasks["default"]
	if !ok {
		t.Fatalf("expected default task to exist")
	}
	if desc, ok := defaultTask["desc"].(string); !ok || desc != "Default task" {
		t.Fatalf("expected default desc to be filled, got %#v", defaultTask["desc"])
	}

	if cmds, ok := defaultTask["cmds"].([]interface{}); !ok || len(cmds) != 3 {
		t.Fatalf("expected default cmds to be filled with 3 entries, got %#v", defaultTask["cmds"])
	} else {
		if taskMap, ok := cmds[0].(map[string]interface{}); !ok || taskMap["task"] != "test:all:cov" {
			t.Fatalf("unexpected first default cmd: %#v", cmds[0])
		}
	}

	aliasTask, ok := parsed.Tasks["alias"]
	if !ok {
		t.Fatalf("expected alias task to exist")
	}
	if desc, ok := aliasTask["desc"].(string); !ok || desc != "List tasks defined in the taskfile.yml" {
		t.Fatalf("expected alias desc to be filled, got %#v", aliasTask["desc"])
	}
	if cmds, ok := aliasTask["cmds"].([]interface{}); !ok || len(cmds) != 1 {
		t.Fatalf("expected alias cmds to include one entry, got %#v", aliasTask["cmds"])
	} else if cmd, ok := cmds[0].(string); !ok || cmd != "task --list-all" {
		t.Fatalf("unexpected alias cmd value: %#v", cmds[0])
	}
}

func TestServiceFillNoChanges(t *testing.T) {
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
	updated, err := service.Fill(taskTypeRoot, tempTarget)
	if err != nil {
		t.Fatalf("fill returned error: %v", err)
	}
	if updated {
		t.Fatalf("expected no updates for complete Taskfile")
	}

	filled, err := os.ReadFile(tempTarget)
	if err != nil {
		t.Fatalf("failed to read Taskfile: %v", err)
	}
	if string(filled) != string(data) {
		t.Fatalf("expected Taskfile to remain unchanged")
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
