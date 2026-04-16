package addheading1

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/landmaster135/devbox/internal/markdown_crafter/infrastructures/filesystem"
)

func TestService_AddHeading1_WithoutFrontMatterHead_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("body\n"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	if _, err := service.Execute(filePath, "Overview", "head"); err != nil {
		t.Fatalf("AddHeading1 returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "# Overview\n\nbody\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddHeading1_WithoutFrontMatterTail_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("body\n"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	if _, err := service.Execute(filePath, "Overview", "tail"); err != nil {
		t.Fatalf("AddHeading1 returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "body\n\n# Overview\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddHeading1_WithFrontMatterHead_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "---\ntitle: mydoc\n---\n\nbody\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	if _, err := service.Execute(filePath, "Overview", "head"); err != nil {
		t.Fatalf("AddHeading1 returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "---\ntitle: mydoc\n---\n\n# Overview\n\nbody\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddHeading1_WithFrontMatterTail_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "---\ntitle: mydoc\n---\n\nbody\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	if _, err := service.Execute(filePath, "Overview", "tail"); err != nil {
		t.Fatalf("AddHeading1 returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "---\ntitle: mydoc\n---\n\nbody\n\n# Overview\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddHeading1_EmptyBody_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	if _, err := service.Execute(filePath, "Overview", "tail"); err != nil {
		t.Fatalf("AddHeading1 returned error: %v", err)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}

	expected := "# Overview\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_AddHeading1_InvalidPosition(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("body"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	if _, err := service.Execute(filePath, "Overview", "middle"); err == nil {
		t.Fatal("expected error for invalid position")
	}
}
