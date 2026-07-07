package workflow

import (
	"reflect"
	"strings"
	"testing"
)

func TestWorkflowHandler_List_Normal(t *testing.T) {
	workflows, err := List(nil)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	if len(workflows) == 0 {
		t.Fatalf("expected at least one workflow")
	}

	wantDescriptions := map[string]bool{
		"Daily PostgreSQL dump with notification":                  false,
		"Daily PostgreSQL extra tables SQL dump with notification": false,
		"Daily PostgreSQL dump for memos with notification":        false,
	}

	for i := range workflows {
		wf := workflows[i]
		if wf.Description == "" {
			t.Fatalf("workflow description is empty at index=%d", i)
		}
		if _, ok := wantDescriptions[wf.Description]; ok {
			wantDescriptions[wf.Description] = true
		}
		if wf.Timezone == "" {
			t.Fatalf("workflow timezone is empty at index=%d", i)
		}
		if _, _, err := wf.GetCronDefinition(); err != nil {
			t.Fatalf("workflow cron definition is invalid at index=%d description=%s: %v", i, wf.Description, err)
		}
		if wf.Process == nil {
			t.Fatalf("workflow process is nil at index=%d description=%s", i, wf.Description)
		}
	}

	for description, found := range wantDescriptions {
		if !found {
			t.Fatalf("workflow %q was not registered", description)
		}
	}
}

func TestNewPostgresDumpTarget_WithExcludeTableData_Normal(t *testing.T) {
	wantTables := []string{postgresAttachmentsTable}
	target := newPostgresDumpTarget("staging", "postgres://example/db", "/tmp/dump", wantTables)

	if target.name != "staging" {
		t.Fatalf("name = %s, want staging", target.name)
	}
	if target.dbURL != "postgres://example/db" {
		t.Fatalf("dbURL = %s, want postgres://example/db", target.dbURL)
	}
	if target.outputDir != "/tmp/dump" {
		t.Fatalf("outputDir = %s, want /tmp/dump", target.outputDir)
	}
	if !reflect.DeepEqual(target.excludeTableData, wantTables) {
		t.Fatalf("excludeTableData = %#v, want %#v", target.excludeTableData, wantTables)
	}
}

func TestNewPostgresDumpTarget_NoExcludeTableData_Normal(t *testing.T) {
	target := newPostgresDumpTarget("memos-staging", "postgres://example/memos", "/tmp/memos", nil)

	if len(target.excludeTableData) != 0 {
		t.Fatalf("excludeTableData = %#v, want empty", target.excludeTableData)
	}
}

func TestNewPostgresExtraSQLDumpTarget_Normal(t *testing.T) {
	wantTables := []string{postgresAttachmentsTable}
	target := newPostgresExtraSQLDumpTarget("staging", "postgres://example/db", "/tmp/dump", wantTables)

	if target.name != "staging" {
		t.Fatalf("name = %s, want staging", target.name)
	}
	if target.dbURL != "postgres://example/db" {
		t.Fatalf("dbURL = %s, want postgres://example/db", target.dbURL)
	}
	if target.outputDir != "/tmp/dump" {
		t.Fatalf("outputDir = %s, want /tmp/dump", target.outputDir)
	}
	if !reflect.DeepEqual(target.tableNames, wantTables) {
		t.Fatalf("tableNames = %#v, want %#v", target.tableNames, wantTables)
	}
}

func TestFormatExtraSQLDumpSummary_Normal(t *testing.T) {
	got := formatExtraSQLDumpSummary("staging", postgresAttachmentsTable)

	if !strings.Contains(got, "## staging extra SQL dump") {
		t.Fatalf("summary = %s, want target heading", got)
	}
	if !strings.Contains(got, "`public.attachments` | sql | completed") {
		t.Fatalf("summary = %s, want table status row", got)
	}
}
