package taskcraftmarkdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestService_Execute_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "tasks.json")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(srcJSONFile, []byte(`[
		{
			"con_id":"TK000002945",
			"page_title":"Daily Task",
			"status_id":"0198d2e4-9a5e-7165-9b5b-80fde571d270",
			"priority":{"page_title":"3-mid"},
			"updated_at":"2026-04-18T10:51:47.182772+09:00",
			"tags":[{"page_title":"JavaScript"}]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "TK000002945.md"), []byte("## t1\nA\n\n## t2\nB\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.Execute("task", "", false, 2945, 2945, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=2") || !strings.Contains(result, "加工成功=2") {
		t.Fatalf("unexpected result: %s", result)
	}

	for _, fileName := range []string{"20260418105147_01.md", "20260418105147_02.md"} {
		data, readErr := os.ReadFile(filepath.Join(outDir, fileName))
		if readErr != nil {
			t.Fatalf("read output failed (%s): %v", fileName, readErr)
		}
		text := string(data)
		if !strings.Contains(text, "# Daily Task") {
			t.Fatalf("heading missing (%s): %s", fileName, text)
		}
		if !strings.Contains(text, "#01-p/todo-status/2-wip") {
			t.Fatalf("status tag missing (%s): %s", fileName, text)
		}
		if !strings.Contains(text, "#01-p/todo-prior/3-mid") {
			t.Fatalf("priority tag missing (%s): %s", fileName, text)
		}
		if !strings.Contains(text, "#91-backup/tool-migration/202602_notion") {
			t.Fatalf("required tag missing (%s): %s", fileName, text)
		}
		if strings.Contains(text, "#javascript") {
			t.Fatalf("task tags should be ignored (%s): %s", fileName, text)
		}
	}
}

func TestService_Execute_MissingSourceCanSkip(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "tasks.json")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.WriteFile(srcJSONFile, []byte(`[
		{
			"con_id":"TK000002946",
			"page_title":"Skip Missing",
			"status_id":"0198d2e4-9a5e-7165-9b5b-80fde571d270",
			"priority":1,
			"updated_at":"2026-04-18T10:51:47.182772+09:00",
			"tags":[]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	service := NewService(filesystem.NewRepository())
	result, err := service.Execute("task", "", true, 2946, 2946, srcJSONFile, filepath.Join(tmpDir, "missing"), outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=1") || !strings.Contains(result, "加工成功=0") || !strings.Contains(result, "スキップ=1") {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestService_Execute_Error(t *testing.T) {
	t.Parallel()

	service := NewService(&filesystem.MockRepository{})
	_, err := service.Execute("content", "", false, 1, 1, "/tmp/tasks.json", "/tmp/body", "/tmp/out")
	if err == nil || err.Error() != "未対応のpage-typeです: content" {
		t.Fatalf("error = %v", err)
	}
}
