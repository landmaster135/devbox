package infrastructures

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSFileSystem_ReadFile_Normal(t *testing.T) {
	fs := NewOSFileSystem()
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "memo.txt")

	if err := os.WriteFile(target, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := fs.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "hello" {
		t.Fatalf("ReadFile() = %q, want %q", string(got), "hello")
	}
}

func TestOSFileSystem_ReadFile_EmptyPath_Error(t *testing.T) {
	fs := NewOSFileSystem()
	_, err := fs.ReadFile("  ")
	if err == nil {
		t.Fatal("ReadFile() error = nil, want error")
	}
}

func TestOSFileSystem_ReadAttachmentFile_Normal(t *testing.T) {
	fs := NewOSFileSystem()
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "memo.txt")

	if err := os.WriteFile(target, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := fs.ReadAttachmentFile(target)
	if err != nil {
		t.Fatalf("ReadAttachmentFile() error = %v", err)
	}
	if got.Filename != "memo.txt" {
		t.Fatalf("filename = %s, want memo.txt", got.Filename)
	}
	if got.ContentType != "text/plain" {
		t.Fatalf("contentType = %s, want text/plain", got.ContentType)
	}
	if string(got.Content) != "hello" {
		t.Fatalf("content = %q, want %q", string(got.Content), "hello")
	}
}

func TestOSFileSystem_ReadAttachmentFile_NotFound_Error(t *testing.T) {
	fs := NewOSFileSystem()
	_, err := fs.ReadAttachmentFile("/tmp/does-not-exist-memos-attachment")
	if err == nil {
		t.Fatal("ReadAttachmentFile() error = nil, want error")
	}
}
