package usecases

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type replaceImagesMockRepository struct {
	readContent string
	readErr     error
	writeErr    error
	writeValue  string
}

func (r *replaceImagesMockRepository) ReadFile(_ string) (string, error) {
	if r.readErr != nil {
		return "", r.readErr
	}
	return r.readContent, nil
}

func (r *replaceImagesMockRepository) WriteFile(_ string, content string) error {
	r.writeValue = content
	return r.writeErr
}

func (r *replaceImagesMockRepository) CreateDir(_ string) error {
	return nil
}

func (r *replaceImagesMockRepository) ListMarkdownFiles(_ string) ([]string, error) {
	return nil, nil
}

func (r *replaceImagesMockRepository) RemoveFile(_ string) error {
	return nil
}

func TestService_ReplaceImages_MultipleImages_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "# Contents\nbefore\n![image](https://example.com/a.png)\nafter\n![image](https://example.com/b.png)\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	result, err := service.ReplaceImages(filePath, "(添付画像)")
	if err != nil {
		t.Fatalf("ReplaceImages returned error: %v", err)
	}
	if !strings.Contains(result, "画像記法 2 件") {
		t.Fatalf("unexpected result: %q", result)
	}

	got, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read updated file: %v", err)
	}
	expected := "# Contents\nbefore\n(添付画像)\nafter\n(添付画像)\n"
	if string(got) != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, string(got))
	}
}

func TestService_ReplaceImages_NoImage_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	content := "# Contents\nno image\n"
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	result, err := service.ReplaceImages(filePath, "(添付画像)")
	if err != nil {
		t.Fatalf("ReplaceImages returned error: %v", err)
	}
	if !strings.Contains(result, "画像記法 0 件") {
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

func TestService_ReplaceImages_EmptyReplacement(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "doc.md")
	if err := os.WriteFile(filePath, []byte("![image](https://example.com/a.png)\n"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	if _, err := service.ReplaceImages(filePath, "   "); err == nil {
		t.Fatal("expected error for empty replacement text")
	}
}

func TestService_ReplaceImages_ReadError(t *testing.T) {
	t.Parallel()

	service := NewService(&replaceImagesMockRepository{
		readErr: errors.New("read error"),
	})
	if _, err := service.ReplaceImages("dummy.md", "(添付画像)"); err == nil {
		t.Fatal("expected read error")
	}
}

func TestService_ReplaceImages_WriteError(t *testing.T) {
	t.Parallel()

	repo := &replaceImagesMockRepository{
		readContent: "![image](https://example.com/a.png)\n",
		writeErr:    errors.New("write error"),
	}
	service := NewService(repo)
	if _, err := service.ReplaceImages("dummy.md", "(添付画像)"); err == nil {
		t.Fatal("expected write error")
	}
	if repo.writeValue != "(添付画像)\n" {
		t.Fatalf("unexpected written content: %q", repo.writeValue)
	}
}
