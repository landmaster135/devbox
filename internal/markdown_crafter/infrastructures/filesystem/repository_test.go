package filesystem

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepository_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	repo := NewRepository()

	targetDir := filepath.Join(tmpDir, "nested")
	if err := repo.CreateDir(targetDir); err != nil {
		t.Fatalf("CreateDir returned error: %v", err)
	}

	filePath := filepath.Join(targetDir, "sample.md")
	if err := repo.WriteFile(filePath, "hello"); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	got, err := repo.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if got != "hello" {
		t.Fatalf("unexpected file content: %q", got)
	}
}

func TestRepository_ReadFile_NotFound(t *testing.T) {
	t.Parallel()

	repo := NewRepository()
	_, err := repo.ReadFile(filepath.Join(t.TempDir(), "missing.md"))
	if err == nil {
		t.Fatal("expected ReadFile error for missing file")
	}
	if !strings.Contains(err.Error(), "ファイルの読み込みに失敗しました") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRepository_CreateDir_InvalidPath(t *testing.T) {
	t.Parallel()

	repo := NewRepository()

	file := filepath.Join(t.TempDir(), "plain.txt")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatalf("failed to prepare file: %v", err)
	}

	err := repo.CreateDir(filepath.Join(file, "child"))
	if err == nil {
		t.Fatal("expected CreateDir error for invalid path")
	}
}

func TestRepository_WriteFile_InvalidPath(t *testing.T) {
	t.Parallel()

	repo := NewRepository()
	err := repo.WriteFile(filepath.Join(t.TempDir(), "missing", "sample.md"), "hello")
	if err == nil {
		t.Fatal("expected WriteFile error for missing directory")
	}
}
