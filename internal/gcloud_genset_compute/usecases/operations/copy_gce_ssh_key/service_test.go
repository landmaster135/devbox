package copygcesshkey

import (
	"strings"
	"testing"
)

func TestServiceBuild_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName: "vm-1",
		Zone:         "us-central1-a",
		SSHKeyPath:   "$HOME/.ssh/google_compute_engine",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParts := []string{
		"if [ -z \"${SSH_AUTH_SOCK:-}\" ]; then eval \"$(ssh-agent -s)\" >/dev/null; fi && ssh-add \"$HOME/.ssh/google_compute_engine\"",
		"gcloud compute scp \"$HOME/.ssh/google_compute_engine\" 'vm-1:/tmp' --zone='us-central1-a' --tunnel-through-iap",
	}
	for _, expected := range expectedParts {
		if !strings.Contains(command, expected) {
			t.Fatalf("command should contain %q: %s", expected, command)
		}
	}
	if strings.Count(command, "&&") != 2 {
		t.Fatalf("command should include 2 conditional chains: %s", command)
	}
}

func TestServiceBuild_QuoteEscape_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName: "vm'o",
		Zone:         "us-central1-a",
		SSHKeyPath:   "/tmp/key'o",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "'/tmp/key'\"'\"'o'") {
		t.Fatalf("quote escaping is not applied: %s", command)
	}
	if !strings.Contains(command, "'vm'\"'\"'o:/tmp'") {
		t.Fatalf("instance quote escaping is not applied: %s", command)
	}
	if !strings.Contains(command, "ssh-add '/tmp/key'\"'\"'o'") {
		t.Fatalf("ssh-add key path quote escaping is not applied: %s", command)
	}
}

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []Params{
		{Zone: "us-central1-a", SSHKeyPath: "$HOME/.ssh/google_compute_engine"},
		{InstanceName: "vm-1", SSHKeyPath: "$HOME/.ssh/google_compute_engine"},
		{InstanceName: "vm-1", Zone: "us-central1-a"},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}
