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
  {"name":"e2-medium","zone":"https://www.googleapis.com/compute/v1/projects/test/zones/asia-southeast3-a","guestCpus":2,"memoryMb":4096,"maximumPersistentDisksSizeGb":257,"deprecated":null},
  {"name":"c3-highmem-4","zone":"https://www.googleapis.com/compute/v1/projects/test/zones/asia-southeast3-a","guestCpus":4,"memoryMb":32768,"maximumPersistentDisksSizeGb":65536,"deprecated":{"state":"DEPRECATED"}}
]`
			return []byte(output), nil
		},
	}
	service := newServiceWithCommandExecutor(mockExecutor)

	result, err := service.Execute(Params{
		MinMemorySizeMiB: 8192,
		MaxMemorySizeMiB: 65536,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mockExecutor.Calls) != 1 {
		t.Fatalf("call count mismatch: %d", len(mockExecutor.Calls))
	}
	if len(mockExecutor.Calls[0].Args) == 0 {
		t.Fatal("args should not be empty")
	}
	if !strings.Contains(strings.Join(mockExecutor.Calls[0].Args, " "), "deprecated") {
		t.Fatalf("deprecated format missing: %v", mockExecutor.Calls[0].Args)
	}
	if !strings.Contains(result, "GUEST_CPUS") {
		t.Fatalf("header missing: %s", result)
	}
	if !strings.Contains(result, "DEPRECATED") {
		t.Fatalf("deprecated header missing: %s", result)
	}
	if strings.Contains(result, "e2-medium") {
		t.Fatalf("e2-medium should be filtered out: %s", result)
	}
	if !strings.Contains(result, "c3-highmem-4") {
		t.Fatalf("filtered row missing: %s", result)
	}
	if !strings.Contains(result, "true") {
		t.Fatalf("deprecated value should be true: %s", result)
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
		MinDiskSizeGiB:   1024,
		MaxDiskSizeGiB:   600000,
		MinMemorySizeMiB: 32000,
		MaxMemorySizeMiB: 70000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "c4-standard-16") {
		t.Fatalf("expected row missing: %s", result)
	}
	if !strings.Contains(result, "false") {
		t.Fatalf("default deprecated value should be false: %s", result)
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
	tests := []Params{
		{MinDiskSizeGiB: 1000, MaxDiskSizeGiB: 500},
		{MinMemorySizeMiB: 8192, MaxMemorySizeMiB: 4096},
		{MinMemorySizeMiB: -1},
	}

	for _, params := range tests {
		if _, err := service.Execute(params); err == nil {
			t.Fatalf("expected validation error for params: %+v", params)
		}
	}
}
