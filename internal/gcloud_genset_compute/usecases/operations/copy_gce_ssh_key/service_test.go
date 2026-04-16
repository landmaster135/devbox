package copygcesshkey

import (
	"strings"
	"testing"
)

func TestServiceBuild_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:  "vm-1",
		Zone:          "us-central1-a",
		SSHKeyPath:    "$HOME/.ssh/google_compute_engine",
		CreatesSSHKey: false,
		Forces:        false,
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

func TestServiceBuild_CreatesSSHKeyWithoutForce_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:  "vm-1",
		Zone:          "us-central1-a",
		SSHKeyPath:    "$HOME/.ssh/google_compute_engine",
		CreatesSSHKey: true,
		Forces:        false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParts := []string{
		"ssh_key_path=\"$HOME/.ssh/google_compute_engine\"",
		"if [ -f \"$ssh_key_path\" ]; then echo \"SSH秘密鍵は既に存在します: $ssh_key_path。上書きするには -forces=true を指定してください\" >&2; exit 1; fi",
		"ssh-keygen -t rsa -f \"$ssh_key_path\"",
	}
	for _, expected := range expectedParts {
		if !strings.Contains(command, expected) {
			t.Fatalf("command should contain %q: %s", expected, command)
		}
	}
}

func TestServiceBuild_CreatesSSHKeyWithForce_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:  "vm-1",
		Zone:          "us-central1-a",
		SSHKeyPath:    "$HOME/.ssh/google_compute_engine",
		CreatesSSHKey: true,
		Forces:        true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParts := []string{
		"echo \"既存SSH秘密鍵を上書きしました: $ssh_key_path\" >&2",
		"rm -f \"$ssh_key_path\" \"$ssh_key_path.pub\"",
		"ssh-keygen -t rsa -f \"$ssh_key_path\"",
	}
	for _, expected := range expectedParts {
		if !strings.Contains(command, expected) {
			t.Fatalf("command should contain %q: %s", expected, command)
		}
	}
}

func TestServiceBuild_QuoteEscape_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:  "vm'o",
		Zone:          "us-central1-a",
		SSHKeyPath:    "/tmp/key'o",
		CreatesSSHKey: false,
		Forces:        false,
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
		{InstanceName: "vm-1", Zone: "us-central1-a", SSHKeyPath: "$HOME/.ssh/google_compute_engine", Forces: true},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}
