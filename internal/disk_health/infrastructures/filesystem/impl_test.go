package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSRepository_ReadFile_Normal(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "smart.log")
	expected := []byte("SMART overall-health self-assessment test result: PASSED\n")

	if err := os.WriteFile(filePath, expected, 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	repository := NewOSRepository()
	actual, err := repository.ReadFile(filePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(actual) != string(expected) {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestOSRepository_ReadFile_MissingFile(t *testing.T) {
	repository := NewOSRepository()

	_, err := repository.ReadFile(filepath.Join(t.TempDir(), "missing.log"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
