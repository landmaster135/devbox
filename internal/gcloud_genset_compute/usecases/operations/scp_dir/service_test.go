package scpdir

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
		SrcDir:       "/tmp/src-dir",
		DestDir:      "workspace",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute scp --recurse '/tmp/src-dir' 'vm-1:workspace/' --zone='us-central1-a'"
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
		SrcDir:       "/tmp/src'o",
		DestDir:      "/dest'o/",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "'/tmp/src'\"'\"'o'") {
		t.Fatalf("src-dir quote escaping is not applied: %s", command)
	}
	if !strings.Contains(command, "'vm'\"'\"'o:/dest'\"'\"'o/'") {
		t.Fatalf("destination quote escaping is not applied: %s", command)
	}
}

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()

	tests := []Params{
		{Zone: "us-central1-a", SrcDir: "/tmp/src-dir", DestDir: "workspace"},
		{InstanceName: "vm-1", SrcDir: "/tmp/src-dir", DestDir: "workspace"},
		{InstanceName: "vm-1", Zone: "us-central1-a", DestDir: "workspace"},
		{InstanceName: "vm-1", Zone: "us-central1-a", SrcDir: "/tmp/src-dir"},
		{InstanceName: "vm-1", Zone: "us-central1-a", SrcDir: "/tmp/src-dir", DestDir: "/"},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}
