package usecases

import (
	"errors"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestService_CheckBodyLength_Normal(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
		return []string{
			"/tmp/body/a.txt",
			"/tmp/body/b.txt",
			"/tmp/body/c.txt",
		}, nil
	}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		switch path {
		case "/tmp/body/a.txt":
			return []byte("ab"), nil
		case "/tmp/body/b.txt":
			return []byte("あい"), nil
		case "/tmp/body/c.txt":
			return []byte("😀😀😀"), nil
		default:
			return []byte(""), nil
		}
	}

	service := NewService(mockRepo)
	result, err := service.CheckBodyLength("/tmp/body", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "対象ファイル総数=3") {
		t.Fatalf("unexpected total count result: %s", result)
	}
	if !strings.Contains(result, "閾値超過ファイル総数=1") {
		t.Fatalf("unexpected exceeded total result: %s", result)
	}
	if !strings.Contains(result, "/tmp/body/c.txt: 3") {
		t.Fatalf("unexpected exceeded file list result: %s", result)
	}
	if strings.Contains(result, "/tmp/body/a.txt") || strings.Contains(result, "/tmp/body/b.txt") {
		t.Fatalf("threshold boundary should not include files with rune count <= threshold: %s", result)
	}
}

func TestService_CheckBodyLength_NoExceededFiles(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
		return []string{"/tmp/body/a.txt"}, nil
	}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		return []byte("a"), nil
	}

	service := NewService(mockRepo)
	result, err := service.CheckBodyLength("/tmp/body", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "閾値超過ファイル総数=0") {
		t.Fatalf("unexpected exceeded count result: %s", result)
	}
	if !strings.Contains(result, "(なし)") {
		t.Fatalf("unexpected empty list result: %s", result)
	}
}

func TestService_CheckBodyLength_Error(t *testing.T) {
	t.Parallel()

	t.Run("empty src body dir", func(t *testing.T) {
		service := NewService(&filesystem.MockRepository{})
		_, err := service.CheckBodyLength("", 1)
		if err == nil || err.Error() != "src-body-dir パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("negative threshold", func(t *testing.T) {
		service := NewService(&filesystem.MockRepository{})
		_, err := service.CheckBodyLength("/tmp/body", -1)
		if err == nil || err.Error() != "threshold パラメータは0以上で必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("list files recursive error", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return nil, errors.New("read dir failed")
		}

		service := NewService(mockRepo)
		_, err := service.CheckBodyLength("/tmp/body", 1)
		if err == nil || !strings.Contains(err.Error(), "src-body-dir の読み取りに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read file error", func(t *testing.T) {
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return []string{"/tmp/body/a.txt"}, nil
		}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return nil, errors.New("read file failed")
		}

		service := NewService(mockRepo)
		_, err := service.CheckBodyLength("/tmp/body", 1)
		if err == nil || !strings.Contains(err.Error(), "ファイルの読み込みに失敗しました (/tmp/body/a.txt)") {
			t.Fatalf("error = %v", err)
		}
	})
}
