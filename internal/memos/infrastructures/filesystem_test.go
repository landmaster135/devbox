package infrastructures

import (
	"os"
	"path/filepath"
	"strings"
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

func TestOSFileSystem_ReadAttachmentFile_DirectoryPath_Error(t *testing.T) {
	fs := NewOSFileSystem()
	dirPath := t.TempDir()

	_, err := fs.ReadAttachmentFile(dirPath)
	if err == nil {
		t.Fatal("ReadAttachmentFile() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "ファイルの読み込みに失敗しました") {
		t.Fatalf("error = %v, want ファイルの読み込みに失敗しました", err)
	}
}

func TestOSFileSystem_ReadAttachmentFile_Markdown_Normal(t *testing.T) {
	fs := NewOSFileSystem()
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "memo.md")

	if err := os.WriteFile(target, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := fs.ReadAttachmentFile(target)
	if err != nil {
		t.Fatalf("ReadAttachmentFile() error = %v", err)
	}
	if got.ContentType != "text/markdown" {
		t.Fatalf("contentType = %s, want text/markdown", got.ContentType)
	}
}

func TestOSFileSystem_ReadAttachmentFile_TextWithoutExtension_Normal(t *testing.T) {
	fs := NewOSFileSystem()
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "memo")

	if err := os.WriteFile(target, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := fs.ReadAttachmentFile(target)
	if err != nil {
		t.Fatalf("ReadAttachmentFile() error = %v", err)
	}
	if got.ContentType != "text/plain" {
		t.Fatalf("contentType = %s, want text/plain", got.ContentType)
	}
}

func TestOSFileSystem_ReadAttachmentFile_EmptyFile_DefaultContentType_Normal(t *testing.T) {
	fs := NewOSFileSystem()
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "empty")

	if err := os.WriteFile(target, []byte{}, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := fs.ReadAttachmentFile(target)
	if err != nil {
		t.Fatalf("ReadAttachmentFile() error = %v", err)
	}
	if got.ContentType != "application/octet-stream" {
		t.Fatalf("contentType = %s, want application/octet-stream", got.ContentType)
	}
}

func TestOSFileSystem_EnsureDirectory_Normal(t *testing.T) {
	fs := NewOSFileSystem()
	tempDir := t.TempDir()

	got, err := fs.EnsureDirectory(tempDir)
	if err != nil {
		t.Fatalf("EnsureDirectory() error = %v", err)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("EnsureDirectory() = %s, want absolute path", got)
	}
}

func TestOSFileSystem_EnsureDirectory_EmptyPath_Error(t *testing.T) {
	fs := NewOSFileSystem()
	_, err := fs.EnsureDirectory("  ")
	if err == nil {
		t.Fatal("EnsureDirectory() error = nil, want error")
	}
}

func TestOSFileSystem_EnsureDirectory_FilePath_Error(t *testing.T) {
	fs := NewOSFileSystem()
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "memo.txt")
	if err := os.WriteFile(target, []byte("hello"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := fs.EnsureDirectory(target)
	if err == nil {
		t.Fatal("EnsureDirectory() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "ディレクトリではありません") {
		t.Fatalf("error = %v, want ディレクトリではありません", err)
	}
}

func TestOSFileSystem_WriteFile_Normal(t *testing.T) {
	fs := NewOSFileSystem()
	tempDir := t.TempDir()
	target := filepath.Join(tempDir, "out.json")

	if err := fs.WriteFile(target, []byte("{}\n")); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "{}\n" {
		t.Fatalf("content = %q, want %q", string(got), "{}\n")
	}
}

func TestOSFileSystem_WriteFile_EmptyPath_Error(t *testing.T) {
	fs := NewOSFileSystem()
	err := fs.WriteFile("  ", []byte("{}"))
	if err == nil {
		t.Fatal("WriteFile() error = nil, want error")
	}
}

func TestNormalizeContentType_Empty_Error(t *testing.T) {
	_, err := normalizeContentType("   ")
	if err == nil {
		t.Fatal("normalizeContentType() error = nil, want error")
	}
}

func TestNormalizeContentType_ParseError_Error(t *testing.T) {
	_, err := normalizeContentType("text/plain extra")
	if err == nil {
		t.Fatal("normalizeContentType() error = nil, want error")
	}
}

func TestNormalizeContentType_WithParams_Normal(t *testing.T) {
	got, err := normalizeContentType("text/plain; charset=utf-8")
	if err != nil {
		t.Fatalf("normalizeContentType() error = %v", err)
	}
	if got != "text/plain" {
		t.Fatalf("normalizeContentType() = %s, want text/plain", got)
	}
}
