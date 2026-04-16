package distributefiles

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestService_Execute_Normal(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		return []byte(`[
			{"con_id":"CO0001","category":"software"},
			{"con_id":"CO0002","category":"tools/dev"},
			{"con_id":"CO0001","category":"software"},
			{"con_id":"","category":"ignored"},
			{"con_id":"CO0003","category":"../books"}
		]`), nil
	}
	mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
		return []string{
			filepath.Join(dirPath, "CO0001.md"),
			filepath.Join(dirPath, "CO0002.md"),
			filepath.Join(dirPath, "CO0003.md"),
			filepath.Join(dirPath, "CO9999.md"),
		}, nil
	}
	mockRepo.FileExistsFunc = func(path string) (bool, error) {
		switch filepath.Base(path) {
		case "CO0001.md", "CO0002.md":
			return true, nil
		case "CO0003.md":
			return false, nil
		default:
			return false, nil
		}
	}

	service := NewService(mockRepo)
	result, err := service.Execute("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "総件数=5") ||
		!strings.Contains(result, "コピー成功=2") ||
		!strings.Contains(result, "未検出=1") ||
		!strings.Contains(result, "スキップ=2") {
		t.Fatalf("unexpected JSON summary: %s", result)
	}
	if !strings.Contains(result, "総md件数=4") ||
		!strings.Contains(result, "JSON対応=3") ||
		!strings.Contains(result, "JSON未対応=1") {
		t.Fatalf("unexpected src-body summary: %s", result)
	}

	if len(mockRepo.CopyFileCalls) != 2 {
		t.Fatalf("copy count = %d, want 2", len(mockRepo.CopyFileCalls))
	}
	if mockRepo.CopyFileCalls[0].DstPath != filepath.Join("/tmp/out", "software", "CO0001.md") {
		t.Fatalf("first copy dst = %s", mockRepo.CopyFileCalls[0].DstPath)
	}
	if mockRepo.CopyFileCalls[1].DstPath != filepath.Join("/tmp/out", "tools_dev", "CO0002.md") {
		t.Fatalf("second copy dst = %s", mockRepo.CopyFileCalls[1].DstPath)
	}

	if len(mockRepo.MkdirAllCalls) != 2 {
		t.Fatalf("mkdir call count = %d, want 2", len(mockRepo.MkdirAllCalls))
	}
}

func TestService_Execute_Error(t *testing.T) {
	t.Parallel()

	t.Run("unsupported page type", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("artifact", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || err.Error() != "未対応のpage-typeです: artifact" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read json error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return nil, errors.New("read failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "src-json-file の読み込みに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte("{invalid"), nil
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "Content JSONの解析に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("list markdown files error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO0001","category":"software"}]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return nil, errors.New("list failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "src-body-dir の読み取りに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("file exists error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO0001","category":"software"}]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return []string{filepath.Join(dirPath, "CO0001.md")}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return false, errors.New("stat failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "コピー元ファイルの確認に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("mkdir error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO0001","category":"software"}]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return []string{filepath.Join(dirPath, "CO0001.md")}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return true, nil
		}
		mockRepo.MkdirAllFunc = func(path string, perm os.FileMode) error {
			return errors.New("mkdir failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "カテゴリディレクトリの作成に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("copy error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`[{"con_id":"CO0001","category":"software"}]`), nil
		}
		mockRepo.ListMarkdownFilesFunc = func(dirPath string) ([]string, error) {
			return []string{filepath.Join(dirPath, "CO0001.md")}, nil
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return true, nil
		}
		mockRepo.CopyFileFunc = func(srcPath, dstPath string) error {
			return errors.New("copy failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("content", "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "Markdownのコピーに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})
}
