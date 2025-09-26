package usecases

import "testing"

func TestBuildCommand_ListDeployments(t *testing.T) {
	svc := NewService()

	cmd, err := svc.BuildCommand(BuildRequest{
		Operation: OperationListDeployments,
		ListDeployments: &ListDeploymentsOptions{
			Project: "my-project",
			Filter:  "name:sample",
			Format:  "table(name,insertTime)",
			Limit:   "5",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{
		"gcloud",
		"deployment-manager",
		"deployments",
		"list",
		"--project=my-project",
		"--filter=name:sample",
		"--format=table(name,insertTime)",
		"--limit=5",
	}

	got := cmd.ArgsWithBinary()
	if len(got) != len(expected) {
		t.Fatalf("unexpected args length: %v", got)
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("arg %d mismatch: expected %s, got %s", i, expected[i], got[i])
		}
	}
}

func TestBuildCommand_ListDeploymentsRequiresOptions(t *testing.T) {
	svc := NewService()
	if _, err := svc.BuildCommand(BuildRequest{Operation: OperationListDeployments}); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestBuildCommand_UnsupportedOperation(t *testing.T) {
	svc := NewService()
	if _, err := svc.BuildCommand(BuildRequest{Operation: Operation("unknown")}); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestCommandStringQuotesArguments(t *testing.T) {
	cmd := Command{
		Binary: "gcloud",
		Args: []string{
			"deployment-manager",
			"deployments",
			"list",
			"--filter=name:sample deployment",
		},
	}

	expected := "gcloud deployment-manager deployments list \"--filter=name:sample deployment\""
	if cmd.String() != expected {
		t.Fatalf("unexpected command string: %s", cmd.String())
	}
}
