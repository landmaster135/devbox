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
}

func TestParseFlagsCustomValues(t *testing.T) {
	args := []string{"--operation=ubuntu", "--network-interface=enp0s3", "--output-dir=/tmp/machine"}
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
}

func TestParseFlagsInvalidOperation(t *testing.T) {
	_, err := ParseFlags([]string{"--operation= "})
	if err == nil {
		t.Fatalf("expected error when operation is empty")
	}
}
