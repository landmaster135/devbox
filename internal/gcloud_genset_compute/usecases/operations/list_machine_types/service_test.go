package listmachinetypes

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
  {"name":"e2-medium","zone":"https://www.googleapis.com/compute/v1/projects/test/zones/asia-southeast3-a","guestCpus":2,"memoryMb":4096,"maximumPersistentDisksSizeGb":257},
  {"name":"c3-highmem-4","zone":"https://www.googleapis.com/compute/v1/projects/test/zones/asia-southeast3-a","guestCpus":4,"memoryMb":32768,"maximumPersistentDisksSizeGb":65536}
]`
			return []byte(output), nil
		},
	}
	service := newServiceWithCommandExecutor(mockExecutor)

	result, err := service.Execute(Params{
		MinSizeGiB: 1024,
		MaxSizeGiB: 70000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mockExecutor.Calls) != 1 {
		t.Fatalf("call count mismatch: %d", len(mockExecutor.Calls))
	}
	if !strings.Contains(result, "GUEST_CPUS") {
		t.Fatalf("header missing: %s", result)
	}
	if strings.Contains(result, "e2-medium") {
		t.Fatalf("e2-medium should be filtered out: %s", result)
	}
	if !strings.Contains(result, "c3-highmem-4") {
		t.Fatalf("filtered row missing: %s", result)
	}
	if !strings.Contains(result, "asia-southeast3-a") {
		t.Fatalf("zone basename missing: %s", result)
	}
}

func TestServiceExecute_StringDiskSizeField_Normal(t *testing.T) {
	t.Parallel()

	mockExecutor := &infrastructures.MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			output := `[
  {"name":"c4-standard-16","zone":"asia-southeast3-a","guestCpus":16,"memoryMb":61440,"maximumPersistentDisksSizeGb":"524288"}
]`
			return []byte(output), nil
		},
	}
	service := newServiceWithCommandExecutor(mockExecutor)

	result, err := service.Execute(Params{
		MinSizeGiB: 1024,
		MaxSizeGiB: 600000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "c4-standard-16") {
		t.Fatalf("expected row missing: %s", result)
	}
}

func TestServiceExecute_CommandError(t *testing.T) {
	t.Parallel()

	mockExecutor := &infrastructures.MockCommandExecutor{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("error"), errors.New("failed")
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
	if _, err := service.Execute(Params{MinSizeGiB: 1000, MaxSizeGiB: 500}); err == nil {
		t.Fatal("expected validation error")
	}
}
