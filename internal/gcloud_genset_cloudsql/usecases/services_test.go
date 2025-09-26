package usecases

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBuildDeleteInstanceCommand(t *testing.T) {
	service := NewService()
	cmd, err := service.BuildDeleteInstanceCommand(DeleteInstanceParams{InstanceName: " prod-db "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"gcloud sql instances patch 'prod-db' --activation-policy=ALWAYS",
		"gcloud sql instances patch 'prod-db' --no-deletion-protection",
		"gcloud sql instances delete 'prod-db'",
	}, " && \\\n")

	if cmd != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, cmd)
	}
}

func TestBuildDeleteInstanceCommand_ErrorWhenInstanceMissing(t *testing.T) {
	service := NewService()
	if _, err := service.BuildDeleteInstanceCommand(DeleteInstanceParams{}); err == nil {
		t.Fatal("expected error when instance name missing")
	}
}

func TestBuildPatchDeletionProtectionCommand(t *testing.T) {
	service := NewService()

	enableCmd, err := service.BuildPatchDeletionProtectionCommand(PatchDeletionProtectionParams{
		InstanceName: "db",
		Mode:         "ENABLE",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enableCmd != "gcloud sql instances patch 'db' --deletion-protection" {
		t.Fatalf("unexpected enable command: %s", enableCmd)
	}

	disableCmd, err := service.BuildPatchDeletionProtectionCommand(PatchDeletionProtectionParams{
		InstanceName: "db",
		Mode:         "disable",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if disableCmd != "gcloud sql instances patch 'db' --no-deletion-protection" {
		t.Fatalf("unexpected disable command: %s", disableCmd)
	}
}

func TestBuildPatchDeletionProtectionCommand_Errors(t *testing.T) {
	service := NewService()

	if _, err := service.BuildPatchDeletionProtectionCommand(PatchDeletionProtectionParams{Mode: "enable"}); err == nil {
		t.Fatal("expected error when instance missing")
	}
	if _, err := service.BuildPatchDeletionProtectionCommand(PatchDeletionProtectionParams{InstanceName: "db"}); err == nil {
		t.Fatal("expected error when mode missing")
	}
	if _, err := service.BuildPatchDeletionProtectionCommand(PatchDeletionProtectionParams{InstanceName: "db", Mode: "unknown"}); err == nil {
		t.Fatal("expected error for invalid mode")
	}
}

func TestBuildPatchActivationPolicyCommand(t *testing.T) {
	service := NewService()

	alwaysCmd, err := service.BuildPatchActivationPolicyCommand(PatchActivationPolicyParams{
		InstanceName: "db",
		Policy:       "ALWAYS",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if alwaysCmd != "gcloud sql instances patch 'db' --activation-policy=ALWAYS" {
		t.Fatalf("unexpected always command: %s", alwaysCmd)
	}

	neverCmd, err := service.BuildPatchActivationPolicyCommand(PatchActivationPolicyParams{
		InstanceName: "db",
		Policy:       "never",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if neverCmd != "gcloud sql instances patch 'db' --activation-policy=never" {
		t.Fatalf("unexpected never command: %s", neverCmd)
	}
}

func TestBuildPatchActivationPolicyCommand_Errors(t *testing.T) {
	service := NewService()

	if _, err := service.BuildPatchActivationPolicyCommand(PatchActivationPolicyParams{Policy: "always"}); err == nil {
		t.Fatal("expected error when instance missing")
	}
	if _, err := service.BuildPatchActivationPolicyCommand(PatchActivationPolicyParams{InstanceName: "db"}); err == nil {
		t.Fatal("expected error when policy missing")
	}
	if _, err := service.BuildPatchActivationPolicyCommand(PatchActivationPolicyParams{InstanceName: "db", Policy: "sometimes"}); err == nil {
		t.Fatal("expected error for invalid policy")
	}
}

func TestBuildStartAndStopInstanceCommand(t *testing.T) {
	service := NewService()

	startCmd, err := service.BuildStartInstanceCommand(InstanceParams{InstanceName: "db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startCmd != "gcloud sql instances patch 'db' --activation-policy=ALWAYS" {
		t.Fatalf("unexpected start command: %s", startCmd)
	}

	stopCmd, err := service.BuildStopInstanceCommand(InstanceParams{InstanceName: "db"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stopCmd != "gcloud sql instances patch 'db' --activation-policy=never" {
		t.Fatalf("unexpected stop command: %s", stopCmd)
	}
}

func TestPrintHighlightedCommand(t *testing.T) {
	service := NewService()
	output := captureStdout(func() {
		service.PrintHighlightedCommand("gcloud sql instances list")
	})

	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "gcloud sql instances list") {
		t.Fatalf("expected command in output: %s", output)
	}
}

func captureStdout(fn func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}
