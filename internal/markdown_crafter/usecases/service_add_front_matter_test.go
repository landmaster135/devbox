package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestService_AddFrontMatter_New_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("hello\n"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddFrontMatter(filePath, []string{"title=New Doc", "author=nov"})
	if err != nil {
		t.Fatalf("AddFrontMatter returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "---\ntitle: New Doc\nauthor: nov\n---\n\nhello\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddFrontMatter_MergeOverwrite_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "---\ntitle: old\nstatus: draft\n---\n\nbody\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddFrontMatter(filePath, []string{"title=new", "category=dev"})
	if err != nil {
		t.Fatalf("AddFrontMatter returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "---\ntitle: new\nstatus: draft\ncategory: dev\n---\n\nbody\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddFrontMatter_InvalidKV(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("body"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddFrontMatter(filePath, []string{"invalid-kv"})
	if err == nil {
		t.Fatal("expected error for invalid kv format")
	}
	if !strings.Contains(err.Error(), "key=value") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_AddFrontMatter_OnlyFrontMatter_NoBody(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "---\ntitle: old\n---\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddFrontMatter(filePath, []string{"title=new", "author=nov"})
	if err != nil {
		t.Fatalf("AddFrontMatter returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}
	expected := "---\ntitle: new\nauthor: nov\n---\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}
