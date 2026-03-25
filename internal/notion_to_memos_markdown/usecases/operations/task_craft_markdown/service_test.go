package taskcraftmarkdown

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	domain "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
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
			"powered_artifacts":[
				{"page_title":"devbox"},
				{"page_title":"ワインを飲んだり勉強するムーヴメント"},
				{"page_title":"devbox"},
				{"page_title":"Unknown Artifact"}
			],
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
	if !strings.Contains(result, "対象件数=1") || !strings.Contains(result, "加工成功=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	const outputFile = "20260418105147_01.md"
	data, readErr := os.ReadFile(filepath.Join(outDir, outputFile))
	if readErr != nil {
		t.Fatalf("read output failed (%s): %v", outputFile, readErr)
	}
	text := string(data)
	if !strings.Contains(text, "# Daily Task") {
		t.Fatalf("heading missing (%s): %s", outputFile, text)
	}
	if !strings.Contains(text, "#01-p/todo-status/2-wip") {
		t.Fatalf("status tag missing (%s): %s", outputFile, text)
	}
	if !strings.Contains(text, "#01-p/todo-prior/3-mid") {
		t.Fatalf("priority tag missing (%s): %s", outputFile, text)
	}
	if !strings.Contains(text, "#91-backup/tool-migration/202602_notion") {
		t.Fatalf("required tag missing (%s): %s", outputFile, text)
	}
	if !strings.Contains(text, "#06-af/system/devbox") {
		t.Fatalf("artifact tag missing (%s): %s", outputFile, text)
	}
	if !strings.Contains(text, "#06-af/diary/hobby") || !strings.Contains(text, "#wine") {
		t.Fatalf("artifact mapped tags missing (%s): %s", outputFile, text)
	}
	if strings.Count(text, "#06-af/system/devbox") != 1 {
		t.Fatalf("artifact duplicate tags should be removed (%s): %s", outputFile, text)
	}
	if strings.Contains(text, "#javascript") {
		t.Fatalf("task tags should be ignored (%s): %s", outputFile, text)
	}
	if strings.Contains(text, "#unknown") {
		t.Fatalf("unknown artifact title should be ignored (%s): %s", outputFile, text)
	}

	if _, statErr := os.Stat(filepath.Join(outDir, "20260418105147_02.md")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("unexpected second output: %v", statErr)
	}
}

func TestService_Execute_UsesSplitSourcesWhenExactMissing_Normal(t *testing.T) {
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
			"con_id":"TK000002949",
			"page_title":"Split Source Task",
			"status_id":"0198d2e4-9a5e-7165-9b5b-80fde571d270",
			"priority":3,
			"updated_at":"2026-04-18T10:51:47+09:00",
			"tags":[]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "TK000002949_01.md"), []byte("first body\n"), 0644); err != nil {
		t.Fatalf("write split body failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "TK000002949_02.md"), []byte("second body\n"), 0644); err != nil {
		t.Fatalf("write split body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.Execute("task", "", false, 2949, 2949, srcJSONFile, srcBodyDir, outDir)
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
		if !strings.Contains(text, "# Split Source Task") {
			t.Fatalf("heading missing (%s): %s", fileName, text)
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

func TestService_Execute_RenameCollisionIncrementsSeconds(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "tasks.json")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	const baseName = "20260418105147_01.md"
	if err := os.WriteFile(filepath.Join(outDir, baseName), []byte("existing\n"), 0644); err != nil {
		t.Fatalf("write existing file failed: %v", err)
	}

	if err := os.WriteFile(srcJSONFile, []byte(`[
		{
			"con_id":"TK000002947",
			"page_title":"Collision Task",
			"status_id":"0198d2e4-9a5e-7165-9b5b-80fde571d270",
			"priority":3,
			"updated_at":"2026-04-18T10:51:47+09:00",
			"tags":[]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "TK000002947.md"), []byte("plain body\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	if _, err := service.Execute("task", "", false, 2947, 2947, srcJSONFile, srcBodyDir, outDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	existing, err := os.ReadFile(filepath.Join(outDir, baseName))
	if err != nil {
		t.Fatalf("read existing file failed: %v", err)
	}
	if string(existing) != "existing\n" {
		t.Fatalf("existing file was overwritten: %q", string(existing))
	}

	if _, err := os.Stat(filepath.Join(outDir, "20260418105148_01.md")); err != nil {
		t.Fatalf("expected collided output not found: %v", err)
	}
}

func TestService_Execute_RenameCollisionIncrementsUntilAvailable(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "tasks.json")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	for _, fileName := range []string{
		"20260418105147_01.md",
		"20260418105148_01.md",
		"20260418105149_01.md",
	} {
		if err := os.WriteFile(filepath.Join(outDir, fileName), []byte("existing\n"), 0644); err != nil {
			t.Fatalf("write existing file failed (%s): %v", fileName, err)
		}
	}

	if err := os.WriteFile(srcJSONFile, []byte(`[
		{
			"con_id":"TK000002948",
			"page_title":"Collision Chain Task",
			"status_id":"0198d2e4-9a5e-7165-9b5b-80fde571d270",
			"priority":3,
			"updated_at":"2026-04-18T10:51:47+09:00",
			"tags":[]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "TK000002948.md"), []byte("plain body\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	if _, err := service.Execute("task", "", false, 2948, 2948, srcJSONFile, srcBodyDir, outDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outDir, "20260418105150_01.md")); err != nil {
		t.Fatalf("expected collided output not found: %v", err)
	}
}

func TestService_resolveTaskRenamePath_FileExistsError(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{
		FileExistsFunc: func(path string) (bool, error) {
			return false, errors.New("file exists failed")
		},
	}
	service := NewService(mockRepo)

	_, err := service.resolveTaskRenamePath(
		"/tmp/out",
		time.Date(2026, 4, 18, 10, 51, 47, 0, time.UTC),
		"01",
	)
	if err == nil || !strings.Contains(err.Error(), "候補ファイルの確認に失敗しました") || !strings.Contains(err.Error(), "file exists failed") {
		t.Fatalf("unexpected error: %v", err)
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

func TestBuildTagsForTask_PoweredArtifactsIgnoresUnknownAndDeduplicates(t *testing.T) {
	t.Parallel()

	task := domain.Task{
		PoweredArtifacts: []domain.TaskPoweredArtifact{
			{PageTitle: "devbox"},
			{PageTitle: "devbox"},
			{PageTitle: "  devbox  "},
			{PageTitle: "Unknown"},
		},
	}

	tags := buildTagsForTask(task)
	joined := "#" + strings.Join(tags, " #")

	if !strings.Contains(joined, "#91-backup/tool-migration/202602_notion") {
		t.Fatalf("required backup tag missing: %s", joined)
	}
	if strings.Count(joined, "#06-af/system/devbox") != 1 {
		t.Fatalf("artifact tags should be deduplicated: %s", joined)
	}
	if strings.Contains(joined, "#unknown") {
		t.Fatalf("unknown artifact title should be ignored: %s", joined)
	}
}
