package addtags

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/markdown_crafter/infrastructures/filesystem"
)

func TestService_AddTags_WithoutFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("body\n"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	_, err := service.ExecuteByFile(filePath, "go,markdown")
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

	service := NewService(filesystem.NewRepository())
	_, err := service.ExecuteByFile(filePath, "go, markdown")
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

	service := NewService(filesystem.NewRepository())
	_, err := service.ExecuteByFile(filePath, "go,#go,md, md ")
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

	service := NewService(filesystem.NewRepository())
	_, err := service.ExecuteByFile(filePath, ",,")
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

	service := NewService(filesystem.NewRepository())
	_, err := service.ExecuteByFile(filePath, "go,markdown")
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

	service := NewService(filesystem.NewRepository())
	_, err := service.ExecuteByFile(filePath, "go")
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

func TestService_AddTagsByDir_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePathA := filepath.Join(tmpDir, "a.md")
	filePathB := filepath.Join(tmpDir, "b.md")
	filePathC := filepath.Join(tmpDir, "c.txt")
	if err := os.WriteFile(filePathA, []byte("body-a\n"), 0644); err != nil {
		t.Fatalf("failed to write file a: %v", err)
	}
	if err := os.WriteFile(filePathB, []byte("---\ntitle: b\n---\n\nbody-b\n"), 0644); err != nil {
		t.Fatalf("failed to write file b: %v", err)
	}
	if err := os.WriteFile(filePathC, []byte("plain text"), 0644); err != nil {
		t.Fatalf("failed to write file c: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.ExecuteByDir(tmpDir, "go,markdown")
	if err != nil {
		t.Fatalf("AddTagsByDir returned error: %v", err)
	}
	if !strings.Contains(result, "2 ファイルにタグを追加しました") {
		t.Fatalf("unexpected result: %q", result)
	}
	if !strings.Contains(result, filePathA) || !strings.Contains(result, filePathB) {
		t.Fatalf("result does not contain updated file paths: %q", result)
	}

	gotA, err := os.ReadFile(filePathA)
	if err != nil {
		t.Fatalf("failed to read file a: %v", err)
	}
	expectedA := "#go #markdown\n\nbody-a\n"
	if string(gotA) != expectedA {
		t.Fatalf("unexpected content a:\nexpected:\n%q\ngot:\n%q", expectedA, string(gotA))
	}

	gotB, err := os.ReadFile(filePathB)
	if err != nil {
		t.Fatalf("failed to read file b: %v", err)
	}
	expectedB := "---\ntitle: b\n---\n\n#go #markdown\n\nbody-b\n"
	if string(gotB) != expectedB {
		t.Fatalf("unexpected content b:\nexpected:\n%q\ngot:\n%q", expectedB, string(gotB))
	}

	gotC, err := os.ReadFile(filePathC)
	if err != nil {
		t.Fatalf("failed to read file c: %v", err)
	}
	if string(gotC) != "plain text" {
		t.Fatalf("unexpected content c: %q", string(gotC))
	}
}

func TestService_AddTagsByDir_NoMarkdownFiles_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	txtPath := filepath.Join(tmpDir, "notes.txt")
	if err := os.WriteFile(txtPath, []byte("plain text"), 0644); err != nil {
		t.Fatalf("failed to write txt file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.ExecuteByDir(tmpDir, "go")
	if err != nil {
		t.Fatalf("AddTagsByDir returned error: %v", err)
	}
	if !strings.Contains(result, "0件") {
		t.Fatalf("unexpected result: %q", result)
	}
}
