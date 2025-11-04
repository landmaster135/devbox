package main

import (
	"bytes"
	"testing"
)

func TestParseFlags_StoresOperation(t *testing.T) {
	stderr := &bytes.Buffer{}
	config, err := parseFlags([]string{"-operation", "mackerel", "-name", "-suffix", "05"}, stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v (stderr: %s)", err, stderr.String())
	}

	if config.Operation != "mackerel" {
		t.Fatalf("expected operation to be mackerel, got %s", config.Operation)
	}

	if !config.SortByName {
		t.Fatalf("expected SortByName to be true")
	}

	if config.Suffix != "05" {
		t.Fatalf("expected suffix to be 05, got %s", config.Suffix)
	}
}

func TestParseFlags_AllowsMissingOperation(t *testing.T) {
	stderr := &bytes.Buffer{}
	config, err := parseFlags([]string{"-time"}, stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.Operation != "" {
		t.Fatalf("expected operation to be empty when not provided")
	}

	if !config.SortByTime {
		t.Fatalf("expected SortByTime to be true")
	}
}
