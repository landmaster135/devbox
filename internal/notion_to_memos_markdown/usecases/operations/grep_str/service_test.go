package grepstr

import (
	"errors"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestService_Execute_Normal(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
		return []string{
			"/tmp/body/a.md",
			"/tmp/body/b.md",
			"/tmp/body/c.md",
		}, nil
	}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		switch path {
		case "/tmp/body/a.md":
			return []byte("TODO: task A"), nil
		case "/tmp/body/b.md":
			return []byte("nothing to see"), nil
		case "/tmp/body/c.md":
			return []byte("line1\nTODO: task C"), nil
		default:
			return []byte(""), nil
		}
	}

	service := NewService(mockRepo)
	result, err := service.Execute("/tmp/body", "TODO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "対象ファイル総数=3") {
		t.Fatalf("unexpected total count result: %s", result)
	}
	if !strings.Contains(result, "該当ファイル総数=2") {
		t.Fatalf("unexpected matched count result: %s", result)
	}
	if !strings.Contains(result, "/tmp/body/a.md") || !strings.Contains(result, "/tmp/body/c.md") {
		t.Fatalf("unexpected matched file list result: %s", result)
	}
	if strings.Contains(result, "/tmp/body/b.md") {
		t.Fatalf("non-matching file should not be listed: %s", result)
	}
}

func TestService_Execute_NoMatch(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
		return []string{"/tmp/body/a.md"}, nil
	}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		return []byte("hello"), nil
	}

	service := NewService(mockRepo)
	result, err := service.Execute("/tmp/body", "TODO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "該当ファイル総数=0") {
		t.Fatalf("unexpected matched count result: %s", result)
	}
	if !strings.Contains(result, "(なし)") {
		t.Fatalf("unexpected empty list result: %s", result)
	}
}

func TestService_Execute_Error(t *testing.T) {
	t.Parallel()

	t.Run("empty src body dir", func(t *testing.T) {
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("", "TODO")
		if err == nil || err.Error() != "src-body-dir パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("empty target string", func(t *testing.T) {
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("/tmp/body", " ")
		if err == nil || err.Error() != "target-str パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("list files recursive error", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return nil, errors.New("read dir failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("/tmp/body", "TODO")
		if err == nil || !strings.Contains(err.Error(), "src-body-dir の読み取りに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read file error", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return []string{"/tmp/body/a.md"}, nil
		}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return nil, errors.New("read file failed")
		}

		service := NewService(mockRepo)
		_, err := service.Execute("/tmp/body", "TODO")
		if err == nil || !strings.Contains(err.Error(), "ファイルの読み込みに失敗しました (/tmp/body/a.md)") {
			t.Fatalf("error = %v", err)
		}
	})
}
