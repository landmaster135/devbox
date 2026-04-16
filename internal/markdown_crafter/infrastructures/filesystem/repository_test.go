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

func TestRepository_ListMarkdownFiles_Normal(t *testing.T) {
	t.Parallel()

	repo := NewRepository()
	dirPath := t.TempDir()

	files := map[string]string{
		"a.md":  "a",
		"b.MD":  "b",
		"c.txt": "c",
	}
	for fileName, content := range files {
		if err := os.WriteFile(filepath.Join(dirPath, fileName), []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}
	if err := os.Mkdir(filepath.Join(dirPath, "nested"), 0755); err != nil {
		t.Fatalf("failed to create nested dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dirPath, "nested", "in_nested.md"), []byte("x"), 0644); err != nil {
		t.Fatalf("failed to write nested file: %v", err)
	}

	markdownFiles, err := repo.ListMarkdownFiles(dirPath)
	if err != nil {
		t.Fatalf("ListMarkdownFiles returned error: %v", err)
	}
	if len(markdownFiles) != 2 {
		t.Fatalf("unexpected markdown file count: %d", len(markdownFiles))
	}
	if markdownFiles[0] != filepath.Join(dirPath, "a.md") {
		t.Fatalf("unexpected first file: %s", markdownFiles[0])
	}
	if markdownFiles[1] != filepath.Join(dirPath, "b.MD") {
		t.Fatalf("unexpected second file: %s", markdownFiles[1])
	}
}

func TestRepository_RemoveFile_Normal(t *testing.T) {
	t.Parallel()

	repo := NewRepository()
	dirPath := t.TempDir()
	filePath := filepath.Join(dirPath, "target.md")
	if err := os.WriteFile(filePath, []byte("x"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := repo.RemoveFile(filePath); err != nil {
		t.Fatalf("RemoveFile returned error: %v", err)
	}
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Fatalf("file should be deleted: %s", filePath)
	}
}
