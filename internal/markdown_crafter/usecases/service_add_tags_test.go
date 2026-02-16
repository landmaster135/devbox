package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestService_AddTags_WithoutFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("body\n"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddTags(filePath, "go,markdown")
	if err != nil {
		t.Fatalf("AddTags returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "#go #markdown\n\nbody\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddTags_WithFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "---\ntitle: mydoc\n---\n\nbody\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddTags(filePath, "go, markdown")
	if err != nil {
		t.Fatalf("AddTags returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "---\ntitle: mydoc\n---\n\n#go #markdown\n\nbody\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddTags_Deduplicate_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("body"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddTags(filePath, "go,#go,md, md ")
	if err != nil {
		t.Fatalf("AddTags returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}
	if !strings.HasPrefix(string(got), "#go #md\n\n") {
		t.Fatalf("unexpected tag line: %q", string(got))
	}
}

func TestService_AddTags_InvalidTags(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("body"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddTags(filePath, ",,")
	if err == nil {
		t.Fatal("expected error for empty tags")
	}
}

func TestService_AddTags_WithFrontMatter_NoBody(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "---\ntitle: mydoc\n---\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddTags(filePath, "go,markdown")
	if err != nil {
		t.Fatalf("AddTags returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "---\ntitle: mydoc\n---\n\n#go #markdown\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddTags_EmptyBody_NoFrontMatter(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.AddTags(filePath, "go")
	if err != nil {
		t.Fatalf("AddTags returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "#go\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}
