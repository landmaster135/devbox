package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSRepositoryEnsureDirAndWriteFile(t *testing.T) {
	repo := NewRepository()
	baseDir := t.TempDir()
	targetDir := filepath.Join(baseDir, "nested", "dir")

	if err := repo.EnsureDir(targetDir, 0o755); err != nil {
		t.Fatalf("EnsureDir returned error: %v", err)
	}

	info, err := os.Stat(targetDir)
	if err != nil {
		t.Fatalf("failed to stat created directory: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %s to be a directory", targetDir)
	}

	filePath := repo.Join(targetDir, "log.json")
	content := []byte("hello, machine-info")
	if err := repo.WriteFile(filePath, content, 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(data) != string(content) {
		t.Fatalf("expected %q, got %q", string(content), string(data))
	}
}

func TestOSRepositoryWriteFileValidatesPath(t *testing.T) {
	repo := NewRepository()
	if err := repo.WriteFile("   ", []byte("noop"), 0o644); err == nil {
		t.Fatalf("expected error for empty path")
	}
}
