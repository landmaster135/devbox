package creategceinstanceandconfigure

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
	if err := os.WriteFile(yamlPath, []byte("FOO: bar\n"), 0o644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:      "vm-1",
		Zone:              "us-central1-a",
		MachineType:       "e2-medium",
		BootDiskSize:      "100GB",
		BootDiskType:      "pd-balanced",
		MetadataYAMLPath:  yamlPath,
		StartupScriptPath: "cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParts := []string{
		"gcloud compute instances create 'vm-1' --zone='us-central1-a' --machine-type='e2-medium' --no-address --boot-disk-size='100GB' --boot-disk-type='pd-balanced'",
		"gcloud compute instances add-metadata 'vm-1' --zone='us-central1-a' --metadata='FOO=bar'",
		"gcloud compute instances add-metadata 'vm-1' --zone='us-central1-a' --metadata-from-file startup-script='cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh'",
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

	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "env.yml")
	if err := os.WriteFile(yamlPath, []byte("ENV: qa'1\n"), 0o644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	service := NewService()
	command, err := service.Build(Params{
		InstanceName:      "vm'o",
		Zone:              "us-central1-a",
		MachineType:       "e2-medium",
		BootDiskSize:      "100GB",
		BootDiskType:      "pd-balanced",
		MetadataYAMLPath:  yamlPath,
		StartupScriptPath: "/tmp/startup'script.sh",
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
		MetadataYAMLPath:  "env.yml",
		StartupScriptPath: "startup-script.sh",
	}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestServiceBuild_MetadataReadError(t *testing.T) {
	t.Parallel()

	service := NewService()
	if _, err := service.Build(Params{
		InstanceName:      "vm-1",
		Zone:              "us-central1-a",
		MachineType:       "e2-medium",
		BootDiskSize:      "100GB",
		BootDiskType:      "pd-balanced",
		MetadataYAMLPath:  "/tmp/not-found.yml",
		StartupScriptPath: "startup-script.sh",
	}); err == nil {
		t.Fatal("expected metadata read error")
	}
}

func TestServiceBuild_AddStartupScriptError(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "env.yml")
	if err := os.WriteFile(yamlPath, []byte("FOO: bar\n"), 0o644); err != nil {
		t.Fatalf("failed to write yaml file: %v", err)
	}

	service := NewService()
	if _, err := service.Build(Params{
		InstanceName:     "vm-1",
		Zone:             "us-central1-a",
		MachineType:      "e2-medium",
		BootDiskSize:     "100GB",
		BootDiskType:     "pd-balanced",
		MetadataYAMLPath: yamlPath,
	}); err == nil {
		t.Fatal("expected startup-script-path validation error")
	}
}
