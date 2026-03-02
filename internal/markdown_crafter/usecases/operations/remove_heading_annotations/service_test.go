package removeheadingannotations

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/markdown_crafter/infrastructures/filesystem"
)

type removeHeadingAnnotationsMockRepository struct {
	readContent string
	readErr     error
	writeErr    error
	writeValue  string
}

func (r *removeHeadingAnnotationsMockRepository) ReadFile(_ string) (string, error) {
	if r.readErr != nil {
		return "", r.readErr
	}
	return r.readContent, nil
}

func (r *removeHeadingAnnotationsMockRepository) WriteFile(_ string, content string) error {
	r.writeValue = content
	return r.writeErr
}

func (r *removeHeadingAnnotationsMockRepository) CreateDir(_ string) error {
	return nil
}

func (r *removeHeadingAnnotationsMockRepository) ListMarkdownFiles(_ string) ([]string, error) {
	return nil, nil
}

func (r *removeHeadingAnnotationsMockRepository) RemoveFile(_ string) error {
	return nil
}

func TestService_RemoveHeadingAnnotations_LevelScoped_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "## **変更しない**\n\n### **対象の見出し**\n本文\n\n### そのまま\n#### **変更しない**\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.Execute(filePath, 3)
	if err != nil {
		t.Fatalf("RemoveHeadingAnnotations returned error: %v", err)
	}
	if !strings.Contains(result, "見出し注釈 1 件") {
		t.Fatalf("unexpected result: %q", result)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}
	expected := "## **変更しない**\n\n### 対象の見出し\n本文\n\n### そのまま\n#### **変更しない**\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_RemoveHeadingAnnotations_NoMatch_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "### 見出し\n本文\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.Execute(filePath, 3)
	if err != nil {
		t.Fatalf("RemoveHeadingAnnotations returned error: %v", err)
	}
	if !strings.Contains(result, "見出し注釈 0 件") {
		t.Fatalf("unexpected result: %q", result)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}
	if string(got) != content {
		t.Fatalf("content should remain unchanged: %q", string(got))
	}
}

func TestService_RemoveHeadingAnnotations_InvalidHeadingLevel(t *testing.T) {
	t.Parallel()

	service := NewService(&removeHeadingAnnotationsMockRepository{})
	if _, err := service.Execute("dummy.md", 0); err == nil {
		t.Fatal("expected error for invalid heading level")
	}
}

func TestService_RemoveHeadingAnnotations_ReadError(t *testing.T) {
	t.Parallel()

	service := NewService(&removeHeadingAnnotationsMockRepository{
		readErr: errors.New("read error"),
	})
	if _, err := service.Execute("dummy.md", 3); err == nil {
		t.Fatal("expected read error")
	}
}

func TestService_RemoveHeadingAnnotations_WriteError(t *testing.T) {
	t.Parallel()

	repo := &removeHeadingAnnotationsMockRepository{
		readContent: "### **見出し**\n",
		writeErr:    errors.New("write error"),
	}
	service := NewService(repo)
	if _, err := service.Execute("dummy.md", 3); err == nil {
		t.Fatal("expected write error")
	}
	if repo.writeValue != "### 見出し\n" {
		t.Fatalf("unexpected written content: %q", repo.writeValue)
	}
}

func TestRemoveHeadingAnnotations_WithoutTrailingNewline_Normal(t *testing.T) {
	t.Parallel()

	content := "### **見出し**"
	got, replaced := removeHeadingAnnotations(content, 3)
	if replaced != 1 {
		t.Fatalf("unexpected replaced count: %d", replaced)
	}
	if got != "### 見出し" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestRemoveHeadingAnnotationFromLine_IndentedHeading_Normal(t *testing.T) {
	t.Parallel()

	got, replaced := removeHeadingAnnotationFromLine("   ### **見出し**\n", "###")
	if !replaced {
		t.Fatal("expected replacement")
	}
	if got != "   ### 見出し\n" {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestRemoveHeadingAnnotationFromLine_FourLeadingSpaces_Normal(t *testing.T) {
	t.Parallel()

	line := "    ### **見出し**\n"
	got, replaced := removeHeadingAnnotationFromLine(line, "###")
	if replaced {
		t.Fatal("unexpected replacement")
	}
	if got != line {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestRemoveHeadingAnnotationFromLine_EmptyUnwrappedText_Normal(t *testing.T) {
	t.Parallel()

	line := "### ****\n"
	got, replaced := removeHeadingAnnotationFromLine(line, "###")
	if replaced {
		t.Fatal("unexpected replacement")
	}
	if got != line {
		t.Fatalf("unexpected content: %q", got)
	}
}

func TestLeadingSpaceCount_Normal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  int
	}{
		{name: "NoLeadingSpace", value: "abc", want: 0},
		{name: "OnlySpaces", value: "   ", want: 3},
		{name: "TwoLeadingSpaces", value: "  abc", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := leadingSpaceCount(tt.value)
			if got != tt.want {
				t.Fatalf("leadingSpaceCount() = %d, want %d", got, tt.want)
			}
		})
	}
}
