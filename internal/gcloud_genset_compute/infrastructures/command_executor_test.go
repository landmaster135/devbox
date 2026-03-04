package infrastructures

import "testing"

func TestMockCommandExecutorExecute_Normal(t *testing.T) {
	t.Parallel()

	executor := &MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte(name + " " + args[0]), nil
		},
	}

	out, err := executor.Execute("gcloud", "version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "gcloud version" {
		t.Fatalf("unexpected output: %s", string(out))
	}
	if len(executor.Calls) != 1 {
		t.Fatalf("call count mismatch: %d", len(executor.Calls))
	}
	if executor.Calls[0].Name != "gcloud" {
		t.Fatalf("name mismatch: %s", executor.Calls[0].Name)
	}
}
