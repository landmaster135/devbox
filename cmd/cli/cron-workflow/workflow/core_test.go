package workflow

import (
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
