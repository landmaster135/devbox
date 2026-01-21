package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepositoryWriteOverwrites(t *testing.T) {
	repo := NewRepository()
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "heartbeat.txt")

	if err := repo.Write(path, true, "first"); err != nil {
		t.Fatalf("first write: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after first write: %v", err)
	}

	if string(got) != "first" {
		t.Fatalf("want first, got %q", string(got))
	}

	if err := repo.Write(path, true, "second"); err != nil {
		t.Fatalf("second write: %v", err)
	}

	got, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("read after second write: %v", err)
	}

	if string(got) != "second" {
		t.Fatalf("want second, got %q", string(got))
	}
}

func TestRepositoryWriteAppend(t *testing.T) {
	repo := NewRepository()
	dir := t.TempDir()
	path := filepath.Join(dir, "heartbeat.txt")

	if err := repo.Write(path, false, "first"); err != nil {
		t.Fatalf("first append: %v", err)
	}

	if err := repo.Write(path, false, "second"); err != nil {
		t.Fatalf("second append: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read appended file: %v", err)
	}

	if string(got) != "firstsecond" {
		t.Fatalf("unexpected content: %q", string(got))
	}
}

func TestRepositoryWriteEmptyPath(t *testing.T) {
	t.Parallel()
	repo := NewRepository()
	if err := repo.Write(" \t\n", true, "noop"); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestRepositoryEnsureDirCreatesDirectory(t *testing.T) {
	repo := NewRepository()
	dir := filepath.Join(t.TempDir(), "nested", "dirs")

	if err := repo.EnsureDir(dir); err != nil {
		t.Fatalf("ensure dir: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat ensured dir: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("expected directory, got %v", info.Mode())
	}
}

func TestRepositoryEnsureDirEmptyPath(t *testing.T) {
	t.Parallel()
	repo := NewRepository()
	if err := repo.EnsureDir("  \t"); err == nil {
		t.Fatal("expected error for empty directory path")
	}
}

func TestWriteFileWrapperUsesDefaultRepository(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wrapper.txt")

	defaultRepository = NewRepository()
	if err := WriteFile(path, true, "wrapper"); err != nil {
		t.Fatalf("write via wrapper: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read wrapper file: %v", err)
	}

	if string(got) != "wrapper" {
		t.Fatalf("unexpected wrapper content: %q", string(got))
	}
}

func TestEnsureDirWrapperUsesDefaultRepository(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "wrapper", "dir")
	defaultRepository = NewRepository()
	if err := EnsureDir(dir); err != nil {
		t.Fatalf("ensure dir via wrapper: %v", err)
	}

	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat after ensure dir: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("expected directory, got %v", info.Mode())
	}
}
