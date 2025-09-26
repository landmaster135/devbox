package usecases

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBuildManagedZonesListCommand(t *testing.T) {
	service := NewService()

	command, err := service.BuildManagedZonesListCommand(ManagedZonesListParams{
		Project: "example-project",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud dns managed-zones list --project='example-project'"
	if command != expected {
		t.Fatalf("unexpected command. expected %q, got %q", expected, command)
	}
}

func TestBuildManagedZonesListCommand_WithAllOptions(t *testing.T) {
	service := NewService()

	command, err := service.BuildManagedZonesListCommand(ManagedZonesListParams{
		Project:        "my-project",
		Format:         "value(name)",
		Filter:         "name:example",
		Limit:          25,
		PageSize:       10,
		SortBy:         "NAME",
		Verbosity:      "debug",
		URI:            true,
		AdditionalArgs: "--account=my-account",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedParts := []string{
		"gcloud dns managed-zones list",
		"--project='my-project'",
		"--format='value(name)'",
		"--filter='name:example'",
		"--limit=25",
		"--page-size=10",
		"--sort-by='NAME'",
		"--verbosity='debug'",
		"--uri",
		"--account=my-account",
	}

	expected := strings.Join(expectedParts, " ")
	if command != expected {
		t.Fatalf("unexpected command. expected %q, got %q", expected, command)
	}
}

func TestBuildManagedZonesListCommand_AdditionalArgsTrim(t *testing.T) {
	service := NewService()

	command, err := service.BuildManagedZonesListCommand(ManagedZonesListParams{
		AdditionalArgs: "   --flags   ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud dns managed-zones list --flags"
	if command != expected {
		t.Fatalf("unexpected command. expected %q, got %q", expected, command)
	}
}

func TestBuildManagedZonesListCommand_EscapesSingleQuote(t *testing.T) {
	service := NewService()

	command, err := service.BuildManagedZonesListCommand(ManagedZonesListParams{
		Filter: "name:'example'",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedFilter := "--filter='name:'\\''example'\\'''"
	if !strings.Contains(command, expectedFilter) {
		t.Fatalf("expected command to contain %q, got %q", expectedFilter, command)
	}
}

func TestBuildManagedZonesListCommand_InvalidLimit(t *testing.T) {
	service := NewService()

	if _, err := service.BuildManagedZonesListCommand(ManagedZonesListParams{Limit: -1}); err == nil {
		t.Fatalf("expected error for negative limit, got nil")
	}
}

func TestBuildManagedZonesListCommand_InvalidPageSize(t *testing.T) {
	service := NewService()

	if _, err := service.BuildManagedZonesListCommand(ManagedZonesListParams{PageSize: -5}); err == nil {
		t.Fatalf("expected error for negative page-size, got nil")
	}
}

func TestShellQuote(t *testing.T) {
	input := "value 'with' quotes"
	expected := "'value '\\''with'\\'' quotes'"

	if got := shellQuote(input); got != expected {
		t.Fatalf("unexpected quoted string. expected %q, got %q", expected, got)
	}
}

func TestPrintHighlightedCommand(t *testing.T) {
	service := NewService()

	output := captureStdout(func() {
		service.PrintHighlightedCommand("gcloud dns managed-zones list")
	})

	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("expected output to contain header, got %q", output)
	}
}

func captureStdout(fn func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic(err)
	}

	return buf.String()
}
