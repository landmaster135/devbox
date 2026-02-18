package usecases

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestDistributeFiles(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[
				{"con_id":"CO0001","category":"software"},
				{"con_id":"CO0002","category":"tools/dev"},
				{"con_id":"CO0001","category":"software"}
			]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return []string{
				filepath.Join(dirPath, "CO0001.md"),
				filepath.Join(dirPath, "CO0002.md"),
				filepath.Join(dirPath, "CO9999.md"),
			}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return true, nil
		}

		service := NewService(mockRepo)
		result, err := service.DistributeFiles("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(result, "総件数=3") || !strings.Contains(result, "コピー成功=2") || !strings.Contains(result, "未検出=0") || !strings.Contains(result, "スキップ=1") {
			t.Fatalf("unexpected result: %s", result)
		}
		if !strings.Contains(result, "総md件数=3") || !strings.Contains(result, "JSON対応=2") || !strings.Contains(result, "JSON未対応=1") {
			t.Fatalf("unexpected src-body-dir result: %s", result)
		}
		if len(mockRepo.CopyFileCalls) != 2 {
			t.Fatalf("copy count = %d, want 2", len(mockRepo.CopyFileCalls))
		}

		expectedFirst := filepath.Join("/tmp/out", "software", "CO0001.md")
		if mockRepo.CopyFileCalls[0].DstPath != expectedFirst {
			t.Fatalf("first copy dst = %s, want %s", mockRepo.CopyFileCalls[0].DstPath, expectedFirst)
		}
		expectedSecond := filepath.Join("/tmp/out", "tools_dev", "CO0002.md")
		if mockRepo.CopyFileCalls[1].DstPath != expectedSecond {
			t.Fatalf("second copy dst = %s, want %s", mockRepo.CopyFileCalls[1].DstPath, expectedSecond)
		}
	})

	t.Run("skip and missing", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[
				{"con_id":"CO0001","category":""},
				{"con_id":"CO0002","category":"software"},
				{"con_id":"","category":"ignored"}
			]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return []string{
				filepath.Join(dirPath, "CO0001.md"),
				filepath.Join(dirPath, "CO0002.md"),
				filepath.Join(dirPath, "CO0003.md"),
			}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			if strings.HasSuffix(path, "CO0001.md") {
				return false, nil
			}
			return true, nil
		}

		service := NewService(mockRepo)
		result, err := service.DistributeFiles("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(result, "総件数=3") || !strings.Contains(result, "コピー成功=1") || !strings.Contains(result, "未検出=1") || !strings.Contains(result, "スキップ=1") {
			t.Fatalf("unexpected result: %s", result)
		}
		if !strings.Contains(result, "総md件数=3") || !strings.Contains(result, "JSON対応=2") || !strings.Contains(result, "JSON未対応=1") {
			t.Fatalf("unexpected src-body-dir result: %s", result)
		}
		if len(mockRepo.CopyFileCalls) != 1 {
			t.Fatalf("copy count = %d, want 1", len(mockRepo.CopyFileCalls))
		}
	})

	t.Run("unsupported page type", func(t *testing.T) {
		service := NewService(&filesystem.MockRepository{})
		_, err := service.DistributeFiles("artifact", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || err.Error() != "未対応のpage-typeです: artifact" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{
			ReadFileFunc: func(path string) ([]byte, error) {
				return []byte(`{bad json`), nil
			},
		}
		service := NewService(mockRepo)
		_, err := service.DistributeFiles("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "Content JSONの解析に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("file exists error", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO0001","category":"software"}]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return []string{
				filepath.Join(dirPath, "CO0001.md"),
			}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return false, errors.New("stat failed")
		}

		service := NewService(mockRepo)
		_, err := service.DistributeFiles("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "コピー元ファイルの確認に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("mkdir error", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO0001","category":"software"}]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return []string{
				filepath.Join(dirPath, "CO0001.md"),
			}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return true, nil
		}
		mockRepo.MkdirAllFunc = func(path string, perm os.FileMode) error {
			return errors.New("mkdir failed")
		}

		service := NewService(mockRepo)
		_, err := service.DistributeFiles("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "カテゴリディレクトリの作成に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("src body dir read error", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO0001","category":"software"}]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return nil, errors.New("read dir failed")
		}

		service := NewService(mockRepo)
		_, err := service.DistributeFiles("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "src-body-dir の読み取りに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestSanitizeCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category string
		want     string
	}{
		{name: "empty", category: "", want: "uncategorized"},
		{name: "normal", category: "software", want: "software"},
		{name: "with slash", category: "tools/dev", want: "tools_dev"},
		{name: "with traversal", category: "../tools/../dev", want: "tools_dev"},
		{name: "windows separator", category: `tools\dev`, want: "tools_dev"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := sanitizeCategory(tt.category)
			if got != tt.want {
				t.Fatalf("sanitizeCategory(%q) = %q, want %q", tt.category, got, tt.want)
			}
		})
	}
}

func TestCountSrcBodyMetrics(t *testing.T) {
	t.Parallel()

	srcBodyFiles := []string{
		"/tmp/body/CO0001.md",
		"/tmp/body/CO0002.md",
		"/tmp/body/CO0003.md",
	}
	jsonConIDSet := map[string]struct{}{
		"CO0001": {},
		"CO0003": {},
	}

	total, mapped, unmapped := countSrcBodyMetrics(srcBodyFiles, jsonConIDSet)
	if total != 3 || mapped != 2 || unmapped != 1 {
		t.Fatalf("metrics = (%d,%d,%d), want (3,2,1)", total, mapped, unmapped)
	}
}
