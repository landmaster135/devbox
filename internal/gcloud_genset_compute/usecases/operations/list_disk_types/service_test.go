package listdisktypes

import (
	"errors"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/infrastructures"
)

func TestServiceExecute_Normal(t *testing.T) {
	t.Parallel()

	mockExecutor := &infrastructures.MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			output := `[
  {"name":"hyperdisk-balanced","zone":"https://www.googleapis.com/compute/v1/projects/test/zones/asia-southeast3-a","validDiskSize":"4GB-65536GB"},
  {"name":"hyperdisk-extreme","zone":"https://www.googleapis.com/compute/v1/projects/test/zones/asia-southeast3-a","validDiskSize":"64GB-65536GB"},
  {"name":"pd-balanced","zone":"https://www.googleapis.com/compute/v1/projects/test/zones/asia-southeast3-a","validDiskSize":"10GB-65536GB"}
]`
			return []byte(output), nil
		},
	}
	service := newServiceWithCommandExecutor(mockExecutor)

	result, err := service.Execute(Params{
		Zones:      []string{"asia-southeast3-a"},
		MinSizeGiB: 4,
		MaxSizeGiB: 65536,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mockExecutor.Calls) != 1 {
		t.Fatalf("call count mismatch: %d", len(mockExecutor.Calls))
	}
	if mockExecutor.Calls[0].Name != "gcloud" {
		t.Fatalf("command name mismatch: %s", mockExecutor.Calls[0].Name)
	}
	if !strings.Contains(result, "NAME") {
		t.Fatalf("header missing: %s", result)
	}
	if !strings.Contains(result, "hyperdisk-extreme") {
		t.Fatalf("hyperdisk-extreme should be included: %s", result)
	}
	if !strings.Contains(result, "asia-southeast3-a") {
		t.Fatalf("zone basename missing: %s", result)
	}
}

func TestServiceExecute_CommandError(t *testing.T) {
	t.Parallel()

	mockExecutor := &infrastructures.MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("permission denied"), errors.New("exec failed")
		},
	}
	service := newServiceWithCommandExecutor(mockExecutor)

	if _, err := service.Execute(Params{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestServiceExecute_InvalidJSONError(t *testing.T) {
	t.Parallel()

	mockExecutor := &infrastructures.MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("{invalid"), nil
		},
	}
	service := newServiceWithCommandExecutor(mockExecutor)

	if _, err := service.Execute(Params{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestServiceExecute_ValidationError(t *testing.T) {
	t.Parallel()

	service := newServiceWithCommandExecutor(&infrastructures.MockCommandExecutor{})
	if _, err := service.Execute(Params{MinSizeGiB: 1024, MaxSizeGiB: 100}); err == nil {
		t.Fatal("expected validation error")
	}
}
