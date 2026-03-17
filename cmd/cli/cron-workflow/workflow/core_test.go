package workflow

import "testing"

func TestWorkflowHandler_List_Normal(t *testing.T) {
	workflows, err := List(nil)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	expectedFrequencyByDescription := map[string]string{
		"Daily Tokyo weather notification":                          "0 1 * * 0-6",
		"Daily heading Discord notification":                        "1 0 * * 0-6",
		"Daily PostgreSQL dump with notification":                   "0 2 * * 0-6",
		"Daily PostgreSQL dump for memos staging with notification": "5 2 * * 0-6",
		"Ubuntu PC info snapshot":                                   "*/10 * * * 0-6",
	}

	if len(workflows) != len(expectedFrequencyByDescription) {
		t.Fatalf("unexpected workflow count: got=%d want=%d", len(workflows), len(expectedFrequencyByDescription))
	}

	for _, wf := range workflows {
		expectedFrequency, ok := expectedFrequencyByDescription[wf.Description]
		if !ok {
			t.Fatalf("unexpected workflow description: %s", wf.Description)
		}
		if wf.Frequency != expectedFrequency {
			t.Fatalf("unexpected frequency for %s: got=%s want=%s", wf.Description, wf.Frequency, expectedFrequency)
		}
		if wf.Timezone != "Asia/Tokyo" {
			t.Fatalf("unexpected timezone for %s: got=%s want=%s", wf.Description, wf.Timezone, "Asia/Tokyo")
		}
		if wf.Process == nil {
			t.Fatalf("workflow process is nil: %s", wf.Description)
		}
	}
}
