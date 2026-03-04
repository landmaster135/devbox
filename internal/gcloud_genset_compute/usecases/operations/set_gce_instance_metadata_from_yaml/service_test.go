package setgceinstancemetadatafromyaml

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestServiceBuild_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "env.yml")
	content := "# comment\nFOO: bar\nEMPTY:    \nINVALID\nBAZ: qux\n"
	if err := os.WriteFile(yamlPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:     "vm-1",
		Zone:             "us-central1-a",
		MetadataYAMLPath: yamlPath,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute instances add-metadata 'vm-1' --zone='us-central1-a' --metadata='FOO=bar,EMPTY=,BAZ=qux'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuild_QuoteEscape_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "env.yml")
	if err := os.WriteFile(yamlPath, []byte("ENV: qa'1\n"), 0o644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:     "vm'o",
		Zone:             "us-central1-a",
		MetadataYAMLPath: yamlPath,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(command, "'vm'\"'\"'o'") {
		t.Fatalf("instance-name quote escaping is not applied: %s", command)
	}
	if !strings.Contains(command, "'ENV=qa'\"'\"'1'") {
		t.Fatalf("metadata quote escaping is not applied: %s", command)
	}
}

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()
	tests := []Params{
		{Zone: "us-central1-a", MetadataYAMLPath: "/tmp/env.yml"},
		{InstanceName: "vm-1", MetadataYAMLPath: "/tmp/env.yml"},
		{InstanceName: "vm-1", Zone: "us-central1-a"},
	}

	for i, params := range tests {
		if _, err := service.Build(params); err == nil {
			t.Fatalf("expected validation error at case %d", i)
		}
	}
}

func TestServiceBuild_YAMLReadError(t *testing.T) {
	t.Parallel()

	service := NewService()
	if _, err := service.Build(Params{
		InstanceName:     "vm-1",
		Zone:             "us-central1-a",
		MetadataYAMLPath: "/tmp/not-found.yml",
	}); err == nil {
		t.Fatal("expected read error")
	}
}

func TestServiceBuild_EmptyMetadataNoOp_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "env.yml")
	if err := os.WriteFile(yamlPath, []byte("# comment\n\nINVALID\n"), 0o644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:     "vm-1",
		Zone:             "us-central1-a",
		MetadataYAMLPath: yamlPath,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(command, "No valid metadata found") {
		t.Fatalf("expected no-op info command, got: %s", command)
	}
}
