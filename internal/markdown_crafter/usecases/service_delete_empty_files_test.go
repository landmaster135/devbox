package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestService_DeleteEmptyFiles_Normal(t *testing.T) {
	t.Parallel()

	dirPath := t.TempDir()

	emptyMarkdown := filepath.Join(dirPath, "empty.md")
	if err := os.WriteFile(emptyMarkdown, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write empty markdown: %v", err)
	}
	miscMarkdown := filepath.Join(dirPath, "misc.md")
	if err := os.WriteFile(miscMarkdown, []byte("# Miscellaneous notes\n- "), 0644); err != nil {
		t.Fatalf("failed to write misc markdown: %v", err)
	}
	miscWithTrailingNewline := filepath.Join(dirPath, "misc_newline.md")
	if err := os.WriteFile(miscWithTrailingNewline, []byte("# Miscellaneous notes\n- \n"), 0644); err != nil {
		t.Fatalf("failed to write misc markdown with trailing newline: %v", err)
	}
	keepMarkdown := filepath.Join(dirPath, "keep.md")
	if err := os.WriteFile(keepMarkdown, []byte("# Miscellaneous notes\n- keep"), 0644); err != nil {
		t.Fatalf("failed to write keep markdown: %v", err)
	}
	otherExt := filepath.Join(dirPath, "empty.txt")
	if err := os.WriteFile(otherExt, []byte(""), 0644); err != nil {
		t.Fatalf("failed to write txt file: %v", err)
	}

	service := NewService(nil)

	result, err := service.DeleteEmptyFiles(dirPath)
	if err != nil {
		t.Fatalf("DeleteEmptyFiles returned error: %v", err)
	}

	if !strings.Contains(result, "delete-empty-files: 3 ファイルを削除しました") {
		t.Fatalf("unexpected result: %s", result)
	}

	for _, deleted := range []string{emptyMarkdown, miscMarkdown, miscWithTrailingNewline} {
		if _, err := os.Stat(deleted); !os.IsNotExist(err) {
			t.Fatalf("file should be deleted: %s", deleted)
		}
	}
	if _, err := os.Stat(keepMarkdown); err != nil {
		t.Fatalf("keep markdown should remain: %v", err)
	}
	if _, err := os.Stat(otherExt); err != nil {
		t.Fatalf("txt file should remain: %v", err)
	}
}

func TestService_DeleteEmptyFiles_NoTarget_Normal(t *testing.T) {
	t.Parallel()

	dirPath := t.TempDir()
	keepMarkdown := filepath.Join(dirPath, "keep.md")
	if err := os.WriteFile(keepMarkdown, []byte("content"), 0644); err != nil {
		t.Fatalf("failed to write keep markdown: %v", err)
	}

	service := NewService(nil)
	result, err := service.DeleteEmptyFiles(dirPath)
	if err != nil {
		t.Fatalf("DeleteEmptyFiles returned error: %v", err)
	}
	if !strings.Contains(result, "delete-empty-files: 0 ファイルを削除しました") {
		t.Fatalf("unexpected result: %s", result)
	}
	if _, err := os.Stat(keepMarkdown); err != nil {
		t.Fatalf("keep markdown should remain: %v", err)
	}
}
