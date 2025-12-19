package config

import (
	"runtime"
	"testing"
)

func TestConfigNormalizeDefaults(t *testing.T) {
	initialWorkers := -1
	cfg := Config{
		ContentID: "  MA  ",
		Suffix:    "",
		Delimiter: " _ ",
		Digits:    0,
		Start:     0,
		Workers:   initialWorkers,
		Operation: "  mackerel  ",
	}

	cfg.Normalize()

	if cfg.ContentID != "MA" {
		t.Fatalf("expected ContentID trimmed to 'MA', got %q", cfg.ContentID)
	}
	if cfg.Delimiter != "_" {
		t.Fatalf("expected Delimiter trimmed to '_', got %q", cfg.Delimiter)
	}
	if cfg.Suffix != "01" {
		t.Fatalf("expected default Suffix '01', got %q", cfg.Suffix)
	}
	if cfg.Digits != 4 {
		t.Fatalf("expected default Digits 4, got %d", cfg.Digits)
	}
	if cfg.Start != 1 {
		t.Fatalf("expected default Start 1, got %d", cfg.Start)
	}
	if cfg.Workers != DefaultWorkers() {
		t.Fatalf("expected default Workers %d, got %d", DefaultWorkers(), cfg.Workers)
	}
	if cfg.Operation != "mackerel" {
		t.Fatalf("expected Operation trimmed to 'mackerel', got %q", cfg.Operation)
	}
}

func TestConfigNormalizePreservesValues(t *testing.T) {
	cfg := Config{
		ContentID: "ID",
		Suffix:    "99",
		Delimiter: "-",
		Digits:    2,
		Start:     3,
		Workers:   8,
		Operation: "wine",
	}

	cfg.Normalize()

	if cfg.Suffix != "99" {
		t.Fatalf("expected suffix to remain '99', got %q", cfg.Suffix)
	}
	if cfg.Digits != 2 {
		t.Fatalf("expected digits to remain 2, got %d", cfg.Digits)
	}
	if cfg.Start != 3 {
		t.Fatalf("expected start to remain 3, got %d", cfg.Start)
	}
	if cfg.Workers != 8 {
		t.Fatalf("expected workers to remain 8, got %d", cfg.Workers)
	}
}

func TestDefaultWorkers(t *testing.T) {
	w := DefaultWorkers()
	cpu := runtime.NumCPU()
	expected := cpu - 1
	if expected < 1 {
		expected = 1
	}

	if w != expected {
		t.Fatalf("expected default workers %d, got %d", expected, w)
	}
	if w < 1 {
		t.Fatalf("default workers must be at least 1, got %d", w)
	}
}
