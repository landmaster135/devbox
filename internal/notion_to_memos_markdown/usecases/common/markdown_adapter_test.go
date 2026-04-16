package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestNewMarkdownService(t *testing.T) {
	t.Parallel()

	service := NewMarkdownService(&filesystem.MockRepository{})
	if service == nil {
		t.Fatalf("service should not be nil")
	}
}

func TestMarkdownCrafterRepositoryAdapter_Delegation(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		if path != "/tmp/input.md" {
			t.Fatalf("unexpected read path: %s", path)
		}
		return []byte("hello"), nil
	}
	mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
		if dirPath != "/tmp/body" {
			t.Fatalf("unexpected list path: %s", dirPath)
		}
		return []string{"/tmp/body/a.md"}, nil
	}

	adapter := newMarkdownCrafterRepositoryAdapter(mockRepo)

	gotRead, err := adapter.ReadFile("/tmp/input.md")
	if err != nil {
		t.Fatalf("ReadFile unexpected error: %v", err)
	}
	if gotRead != "hello" {
		t.Fatalf("ReadFile = %q, want %q", gotRead, "hello")
	}

	if err := adapter.WriteFile("/tmp/out.md", "content"); err != nil {
		t.Fatalf("WriteFile unexpected error: %v", err)
	}
	if len(mockRepo.WriteFileCalls) != 1 {
		t.Fatalf("WriteFileCalls = %d, want 1", len(mockRepo.WriteFileCalls))
	}
	if mockRepo.WriteFileCalls[0].Path != "/tmp/out.md" {
		t.Fatalf("WriteFile path = %s", mockRepo.WriteFileCalls[0].Path)
	}
	if string(mockRepo.WriteFileCalls[0].Data) != "content" {
		t.Fatalf("WriteFile data = %q", string(mockRepo.WriteFileCalls[0].Data))
	}

	if err := adapter.CreateDir("/tmp/newdir"); err != nil {
		t.Fatalf("CreateDir unexpected error: %v", err)
	}
	if len(mockRepo.MkdirAllCalls) != 1 {
		t.Fatalf("MkdirAllCalls = %d, want 1", len(mockRepo.MkdirAllCalls))
	}
	if mockRepo.MkdirAllCalls[0].Path != "/tmp/newdir" {
		t.Fatalf("MkdirAll path = %s", mockRepo.MkdirAllCalls[0].Path)
	}
	if mockRepo.MkdirAllCalls[0].Perm != DefaultDirectoryPerm {
		t.Fatalf("MkdirAll perm = %v, want %v", mockRepo.MkdirAllCalls[0].Perm, DefaultDirectoryPerm)
	}

	gotList, err := adapter.ListMarkdownFiles("/tmp/body")
	if err != nil {
		t.Fatalf("ListMarkdownFiles unexpected error: %v", err)
	}
	if len(gotList) != 1 || gotList[0] != "/tmp/body/a.md" {
		t.Fatalf("ListMarkdownFiles = %#v", gotList)
	}
}

func TestMarkdownCrafterRepositoryAdapter_ReadFileError(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		return nil, errors.New("read failed")
	}

	adapter := newMarkdownCrafterRepositoryAdapter(mockRepo)
	_, err := adapter.ReadFile("/tmp/input.md")
	if err == nil || !strings.Contains(err.Error(), "read failed") {
		t.Fatalf("error = %v", err)
	}
}

func TestMarkdownCrafterRepositoryAdapter_RemoveFile(t *testing.T) {
	t.Parallel()

	adapter := newMarkdownCrafterRepositoryAdapter(&filesystem.MockRepository{})

	tmpDir := t.TempDir()
	target := filepath.Join(tmpDir, "target.md")
	if err := os.WriteFile(target, []byte("body"), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	if err := adapter.RemoveFile(target); err != nil {
		t.Fatalf("RemoveFile unexpected error: %v", err)
	}
	if _, err := os.Stat(target); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("file should be removed, stat error = %v", err)
	}
}

func TestMarkdownCrafterRepositoryAdapter_RemoveFileError(t *testing.T) {
	t.Parallel()

	adapter := newMarkdownCrafterRepositoryAdapter(&filesystem.MockRepository{})
	nonExistent := filepath.Join(t.TempDir(), "missing.md")

	err := adapter.RemoveFile(nonExistent)
	if err == nil || !strings.Contains(err.Error(), "ファイルの削除に失敗しました") {
		t.Fatalf("error = %v", err)
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("error should wrap os.ErrNotExist: %v", err)
	}
}
