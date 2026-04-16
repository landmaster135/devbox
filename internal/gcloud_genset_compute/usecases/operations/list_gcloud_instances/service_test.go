package listgcloudinstances

import "testing"

func TestServiceBuild_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{
		Filter: "zone:us-central1-a",
		Format: "table(name,status)",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute instances list --filter='zone:us-central1-a' --format='table(name,status)'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuild_WithoutFilter_Normal(t *testing.T) {
	t.Parallel()

	service := NewService()
	command, err := service.Build(Params{Format: "json"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "gcloud compute instances list --format='json'"
	if command != expected {
		t.Fatalf("command mismatch:\nexpected: %s\nactual:   %s", expected, command)
	}
}

func TestServiceBuild_ValidationError(t *testing.T) {
	t.Parallel()

	service := NewService()
	if _, err := service.Build(Params{}); err == nil {
		t.Fatal("expected validation error")
	}
}
