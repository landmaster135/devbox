package usecases

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestBuildAuthLoginCommand(t *testing.T) {
	t.Parallel()

	service := NewService()

	command, err := service.BuildAuthLoginCommand(AuthLoginParams{ProjectID: "my-project"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud auth login 'my-project'"
	if command != expected {
		t.Fatalf("command mismatch: got %q, want %q", command, expected)
	}
}

func TestBuildAuthLoginCommand_WithAdditionalArgs(t *testing.T) {
	t.Parallel()

	service := NewService()

	command, err := service.BuildAuthLoginCommand(AuthLoginParams{
		ProjectID:      "project with spaces",
		AdditionalArgs: "--quiet --brief",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud auth login 'project with spaces' --quiet --brief"
	if command != expected {
		t.Fatalf("command mismatch: got %q, want %q", command, expected)
	}
}

func TestBuildAuthLoginCommand_Errors(t *testing.T) {
	t.Parallel()

	service := NewService()

	if _, err := service.BuildAuthLoginCommand(AuthLoginParams{}); err == nil {
		t.Fatal("expected error when project-id is empty")
	}
}

func TestBuildSetProjectConfigCommand(t *testing.T) {
	t.Parallel()

	service := NewService()

	command, err := service.BuildSetProjectConfigCommand(SetProjectConfigParams{ProjectID: "sample-project"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud config set project 'sample-project'"
	if command != expected {
		t.Fatalf("command mismatch: got %q, want %q", command, expected)
	}
}

func TestBuildSetProjectConfigCommand_WithAdditionalArgs(t *testing.T) {
	t.Parallel()

	service := NewService()

	command, err := service.BuildSetProjectConfigCommand(SetProjectConfigParams{
		ProjectID:      "proj'ect",
		AdditionalArgs: "--configuration=default",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "gcloud config set project 'proj'\"'\"'ect' --configuration=default"
	if command != expected {
		t.Fatalf("command mismatch: got %q, want %q", command, expected)
	}
}

func TestBuildSetProjectConfigCommand_Errors(t *testing.T) {
	t.Parallel()

	service := NewService()

	if _, err := service.BuildSetProjectConfigCommand(SetProjectConfigParams{}); err == nil {
		t.Fatal("expected error when project-id is empty")
	}
}

func TestPrintHighlightedCommand(t *testing.T) {
	t.Parallel()

	service := NewService()

	output := captureStdout(func() {
		service.PrintHighlightedCommand("gcloud auth login 'example'")
	})

	if !strings.Contains(output, "生成された gcloud コマンド") {
		t.Fatalf("expected header in output: %s", output)
	}
	if !strings.Contains(output, "gcloud auth login 'example'") {
		t.Fatalf("expected command in output: %s", output)
	}
}

func TestShellQuote(t *testing.T) {
	if got := shellQuote("proj'ect"); got != "'proj'\"'\"'ect'" {
		t.Fatalf("unexpected quoting: %s", got)
	}
	if got := shellQuote("plain"); got != "'plain'" {
		t.Fatalf("unexpected quoting for plain string: %s", got)
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
