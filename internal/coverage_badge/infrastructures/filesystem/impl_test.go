package filesystem

import (
	"path/filepath"
	"testing"
)

type TestOSRepository struct{}

func TestOSRepositoryReadAndWriteFile_Normal(t *testing.T) {
	repo := NewRepository()
	baseDir := t.TempDir()
	targetFile := filepath.Join(baseDir, "README.md")
	content := []byte("coverage badge")

	if err := repo.WriteFile(targetFile, content, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := repo.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(got) != string(content) {
		t.Fatalf("ReadFile() = %q, want %q", string(got), string(content))
	}
}

func TestOSRepositoryReadFile_EmptyPath(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.ReadFile("   "); err == nil {
		t.Fatal("ReadFile() error = nil, want error")
	}
}

func TestOSRepositoryWriteFile_EmptyPath(t *testing.T) {
	repo := NewRepository()
	if err := repo.WriteFile("   ", []byte("x"), 0o644); err == nil {
		t.Fatal("WriteFile() error = nil, want error")
	}
}
