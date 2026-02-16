package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestService_SplitHeadings_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "source.md")
	outputDir := filepath.Join(tmpDir, "out")

	content := "# overview\n\n## bbb\n\ntest1\n\n## ccc\n\ntest2\n"
	if err := os.WriteFile(sourcePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	result, err := service.SplitHeadings(sourcePath, 2, outputDir)
	if err != nil {
		t.Fatalf("SplitHeadings returned error: %v", err)
	}

	if !strings.Contains(result, "2 ファイルを出力しました") {
		t.Fatalf("unexpected result: %s", result)
	}

	file1, err := os.ReadFile(filepath.Join(outputDir, "001.md"))
	if err != nil {
		t.Fatalf("failed to read output file 001.md: %v", err)
	}
	file2, err := os.ReadFile(filepath.Join(outputDir, "002.md"))
	if err != nil {
		t.Fatalf("failed to read output file 002.md: %v", err)
	}

	if string(file1) != "## bbb\n\ntest1\n\n" {
		t.Fatalf("unexpected content for 001.md: %q", string(file1))
	}
	if string(file2) != "## ccc\n\ntest2\n" {
		t.Fatalf("unexpected content for 002.md: %q", string(file2))
	}
}

func TestService_SplitHeadings_HeadingNotFound(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	sourcePath := filepath.Join(tmpDir, "source.md")
	outputDir := filepath.Join(tmpDir, "out")
	if err := os.WriteFile(sourcePath, []byte("# title\n\ntext"), 0644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	service := NewService(nil)
	_, err := service.SplitHeadings(sourcePath, 2, outputDir)
	if err == nil {
		t.Fatal("expected error when heading level does not exist")
	}
}

func TestExtractSectionsByHeadingLevel_Normal(t *testing.T) {
	t.Parallel()

	content := "## a\ntext-a\n## b\ntext-b"
	sections := extractSectionsByHeadingLevel(content, 2)
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}

	if sections[0] != "## a\ntext-a\n" {
		t.Fatalf("unexpected first section: %q", sections[0])
	}
	if sections[1] != "## b\ntext-b" {
		t.Fatalf("unexpected second section: %q", sections[1])
	}
}
