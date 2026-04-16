package creategceinstancewithstartupscript

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
		MachineType:       "e2-medium",
		BootDiskSize:      "100GB",
		BootDiskType:      "pd-balanced",
		StartupScriptPath: "cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParts := []string{
		"gcloud compute instances create 'vm-1' --zone='us-central1-a' --machine-type='e2-medium' --no-address --boot-disk-size='100GB' --boot-disk-type='pd-balanced'",
		"gcloud compute instances add-metadata 'vm-1' --zone='us-central1-a' --metadata-from-file startup-script='cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh'",
	}
	for _, expected := range expectedParts {
		if !strings.Contains(command, expected) {
			t.Fatalf("command should contain %q: %s", expected, command)
		}
	}
	if strings.Count(command, "&&") != 1 {
		t.Fatalf("command should include 1 conditional chain: %s", command)
	}
}

func TestServiceBuild_QuoteEscape_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:      "vm'o",
		Zone:              "us-central1-a",
		MachineType:       "e2-medium",
		BootDiskSize:      "100GB",
		BootDiskType:      "pd-balanced",
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
	if _, err := service.Build(Params{
		Zone:              "us-central1-a",
		MachineType:       "e2-medium",
		BootDiskSize:      "100GB",
		BootDiskType:      "pd-balanced",
		StartupScriptPath: "startup-script.sh",
	}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestServiceBuild_AddStartupScriptError(t *testing.T) {
	t.Parallel()

	service := NewService()
	if _, err := service.Build(Params{
		InstanceName: "vm-1",
		Zone:         "us-central1-a",
		MachineType:  "e2-medium",
		BootDiskSize: "100GB",
		BootDiskType: "pd-balanced",
	}); err == nil {
		t.Fatal("expected startup-script-path validation error")
	}
}
