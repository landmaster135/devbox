package filesystem

import (
	"errors"
	"os"
	"testing"
)

type TestMockRepository struct{}

func TestMockRepositoryReadFile_Normal(t *testing.T) {
	mockRepo := &MockRepository{
		ReadFileFunc: func(path string) ([]byte, error) {
			if path != "coverage.out" {
				t.Fatalf("path = %q, want coverage.out", path)
			}
			return []byte("ok"), nil
		},
	}

	got, err := mockRepo.ReadFile("coverage.out")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "ok" {
		t.Fatalf("ReadFile() = %q, want ok", string(got))
	}
	if mockRepo.LastReadPath != "coverage.out" {
		t.Fatalf("LastReadPath = %q, want coverage.out", mockRepo.LastReadPath)
	}
}

func TestMockRepositoryWriteFile_Normal(t *testing.T) {
	mockRepo := &MockRepository{
		WriteFileFunc: func(path string, data []byte, perm os.FileMode) error {
			if path != "README.md" {
				t.Fatalf("path = %q, want README.md", path)
			}
			if string(data) != "content" {
				t.Fatalf("data = %q, want content", string(data))
			}
			if perm != 0o644 {
				t.Fatalf("perm = %v, want 0644", perm)
			}
			return nil
		},
	}

	if err := mockRepo.WriteFile("README.md", []byte("content"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if mockRepo.WriteCallCount != 1 {
		t.Fatalf("WriteCallCount = %d, want 1", mockRepo.WriteCallCount)
	}
}

func TestMockRepositoryWriteFile_Error(t *testing.T) {
	mockRepo := &MockRepository{
		WriteFileFunc: func(path string, data []byte, perm os.FileMode) error {
			return errors.New("write failed")
		},
	}

	if err := mockRepo.WriteFile("README.md", []byte("content"), 0o644); err == nil {
		t.Fatal("WriteFile() error = nil, want error")
	}
}
