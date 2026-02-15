package filesystem

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestOSRepositoryReadFileAndWriteFile_Normal(t *testing.T) {
	repo := NewRepository()
	baseDir := t.TempDir()
	filePath := filepath.Join(baseDir, "status.md")
	content := []byte("hello")

	if err := repo.WriteFile(filePath, content, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := repo.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(got) != string(content) {
		t.Fatalf("ReadFile() = %q, want %q", string(got), string(content))
	}
}

func TestOSRepositoryListDirectories_Normal(t *testing.T) {
	repo := NewRepository()
	baseDir := t.TempDir()

	if err := os.MkdirAll(filepath.Join(baseDir, "cli"), 0o755); err != nil {
		t.Fatalf("failed to create cli dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(baseDir, "mcp"), 0o755); err != nil {
		t.Fatalf("failed to create mcp dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(baseDir, "README.md"), []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	got, err := repo.ListDirectories(baseDir)
	if err != nil {
		t.Fatalf("ListDirectories() error = %v", err)
	}

	want := []string{"cli", "mcp"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListDirectories() = %v, want %v", got, want)
	}
}

func TestOSRepositoryListDirectories_NonExistentDirectory(t *testing.T) {
	repo := NewRepository()

	got, err := repo.ListDirectories(filepath.Join(t.TempDir(), "not-found"))
	if err != nil {
		t.Fatalf("ListDirectories() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("ListDirectories() len = %d, want 0", len(got))
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

func TestOSRepositoryListDirectories_EmptyPath(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.ListDirectories("   "); err == nil {
		t.Fatal("ListDirectories() error = nil, want error")
	}
}

func TestOSRepositoryJoin_Normal(t *testing.T) {
	repo := NewRepository()
	got := repo.Join("a", "b", "c")
	want := filepath.Join("a", "b", "c")
	if got != want {
		t.Fatalf("Join() = %s, want %s", got, want)
	}
}
