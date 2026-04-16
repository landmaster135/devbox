package creategceinstance

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
		MachineType:  "e2-medium",
		BootDiskSize: "100GB",
		BootDiskType: "pd-balanced",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute instances create 'vm-1' --zone='us-central1-a' --machine-type='e2-medium' --no-address --boot-disk-size='100GB' --boot-disk-type='pd-balanced'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuild_QuoteEscape_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		InstanceName: "vm'o",
		Zone:         "us-central1-a",
		MachineType:  "e2-medium",
		BootDiskSize: "100GB",
		BootDiskType: "pd-balanced",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "'vm'\"'\"'o'") {
		t.Fatalf("quote escaping is not applied: %s", command)
	}
}

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []Params{
		{Zone: "us-central1-a", MachineType: "e2-medium", BootDiskSize: "100GB", BootDiskType: "pd-balanced"},
		{InstanceName: "vm-1", MachineType: "e2-medium", BootDiskSize: "100GB", BootDiskType: "pd-balanced"},
		{InstanceName: "vm-1", Zone: "us-central1-a", BootDiskSize: "100GB", BootDiskType: "pd-balanced"},
		{InstanceName: "vm-1", Zone: "us-central1-a", MachineType: "e2-medium", BootDiskType: "pd-balanced"},
		{InstanceName: "vm-1", Zone: "us-central1-a", MachineType: "e2-medium", BootDiskSize: "100GB"},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}
