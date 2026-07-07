package workflow

import (
	"reflect"
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

	for i := range workflows {
		wf := workflows[i]
		if wf.Description == "" {
			t.Fatalf("workflow description is empty at index=%d", i)
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
}

func TestNewPostgresDumpTarget_SplitAttachments_Normal(t *testing.T) {
	target := newPostgresDumpTarget("staging", "postgres://example/db", "/tmp/dump", true)

	if target.name != "staging" {
		t.Fatalf("name = %s, want staging", target.name)
	}
	if target.dbURL != "postgres://example/db" {
		t.Fatalf("dbURL = %s, want postgres://example/db", target.dbURL)
	}
	if target.outputDir != "/tmp/dump" {
		t.Fatalf("outputDir = %s, want /tmp/dump", target.outputDir)
	}
	wantTables := []string{postgresAttachmentsTable}
	if !reflect.DeepEqual(target.excludeTableData, wantTables) {
		t.Fatalf("excludeTableData = %#v, want %#v", target.excludeTableData, wantTables)
	}
	if !reflect.DeepEqual(target.extraSQLDumpTables, wantTables) {
		t.Fatalf("extraSQLDumpTables = %#v, want %#v", target.extraSQLDumpTables, wantTables)
	}
}

func TestNewPostgresDumpTarget_NoSplitAttachments_Normal(t *testing.T) {
	target := newPostgresDumpTarget("memos-staging", "postgres://example/memos", "/tmp/memos", false)

	if len(target.excludeTableData) != 0 {
		t.Fatalf("excludeTableData = %#v, want empty", target.excludeTableData)
	}
	if len(target.extraSQLDumpTables) != 0 {
		t.Fatalf("extraSQLDumpTables = %#v, want empty", target.extraSQLDumpTables)
	}
}
