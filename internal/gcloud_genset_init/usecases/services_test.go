package usecases

import "testing"

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
