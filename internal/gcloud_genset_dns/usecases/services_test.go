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

func TestBuildManagedZonesListCommand_ReturnsError(t *testing.T) {
	service := NewService()

	if _, err := service.BuildManagedZonesListCommand(ManagedZonesListParams{Limit: -1}); err == nil {
		t.Fatalf("expected error for negative limit, got nil")
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
