package addstartupscripttogceinstance

import (
	"strings"
	"testing"
)

func TestServiceBuild_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:      "vm-1",
		Zone:              "us-central1-a",
		StartupScriptPath: "cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute instances add-metadata 'vm-1' --zone='us-central1-a' --metadata-from-file startup-script='cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuild_QuoteEscape_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:      "vm'o",
		Zone:              "us-central1-a",
		StartupScriptPath: "/tmp/startup'script.sh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "'vm'\"'\"'o'") {
		t.Fatalf("instance-name quote escaping is not applied: %s", command)
	}
	if !strings.Contains(command, "'/tmp/startup'\"'\"'script.sh'") {
		t.Fatalf("startup-script-path quote escaping is not applied: %s", command)
	}
}

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()
	tests := []Params{
		{Zone: "us-central1-a", StartupScriptPath: "startup.sh"},
		{InstanceName: "vm-1", StartupScriptPath: "startup.sh"},
		{InstanceName: "vm-1", Zone: "us-central1-a"},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}
