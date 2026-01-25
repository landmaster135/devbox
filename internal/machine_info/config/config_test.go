package config

import "testing"

func TestParseFlagsDefaults(t *testing.T) {
	cfg, err := ParseFlags([]string{})
	if err != nil {
		t.Fatalf("ParseFlags returned error: %v", err)
	}

	if cfg.Operation != OperationUbuntu {
		t.Fatalf("expected default operation %q, got %q", OperationUbuntu, cfg.Operation)
	}
	if cfg.NetworkInterface != "eth0" {
		t.Fatalf("expected default network interface eth0, got %s", cfg.NetworkInterface)
	}
	if cfg.OutputDir != "." {
		t.Fatalf("expected default output dir '.', got %s", cfg.OutputDir)
	}
	if cfg.MemoryManufacturers != "" {
		t.Fatalf("expected default memory manufacturers '', got %s", cfg.MemoryManufacturers)
	}
	if cfg.MemoryNames != "" {
		t.Fatalf("expected default memory names '', got %s", cfg.MemoryNames)
	}
}

func TestParseFlagsCustomValues(t *testing.T) {
	args := []string{
		"--operation=ubuntu",
		"--network-interface=enp0s3",
		"--output-dir=/tmp/machine",
		"--memory-manufacturers=MakerA,MakerB",
		"--memory-names=Part1,Part2",
	}
	cfg, err := ParseFlags(args)
	if err != nil {
		t.Fatalf("ParseFlags returned error: %v", err)
	}

	if cfg.NetworkInterface != "enp0s3" {
		t.Fatalf("expected enp0s3, got %s", cfg.NetworkInterface)
	}
	if cfg.OutputDir != "/tmp/machine" {
		t.Fatalf("expected /tmp/machine, got %s", cfg.OutputDir)
	}
	if cfg.MemoryManufacturers != "MakerA,MakerB" {
		t.Fatalf("expected memory manufacturers MakerA,MakerB got %s", cfg.MemoryManufacturers)
	}
	if cfg.MemoryNames != "Part1,Part2" {
		t.Fatalf("expected memory names Part1,Part2 got %s", cfg.MemoryNames)
	}
}

func TestParseFlagsInvalidOperation(t *testing.T) {
	_, err := ParseFlags([]string{"--operation= "})
	if err == nil {
		t.Fatalf("expected error when operation is empty")
	}
}
