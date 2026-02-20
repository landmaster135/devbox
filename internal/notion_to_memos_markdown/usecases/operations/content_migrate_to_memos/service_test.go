package migratetomemos

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
	memos "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/memos"
)

func TestService_Execute_Normal(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
		if dirPath == "/tmp/body" {
			return []string{
				"/tmp/body/CO000000001.md",
				"/tmp/body/CO000000002.md",
				"/tmp/body/ignore.txt",
			}, nil
		}
		return []string{
			"/tmp/resources/CO000000001_1.png",
			"/tmp/resources/CO000000001_2.webp",
			"/tmp/resources/CO000000009_1.png",
		}, nil
	}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		switch path {
		case "/tmp/body/CO000000001.md":
			return []byte("memo-1"), nil
		case "/tmp/body/CO000000002.md":
			return []byte("memo-2"), nil
		default:
			return nil, errors.New("unexpected read path")
		}
	}

	mockClient := &memos.MockClient{}
	mockClient.CreateMemoFunc = func(_ context.Context, content string) (string, error) {
		switch content {
		case "memo-1":
			return "memos/aaa", nil
		case "memo-2":
			return "memos/bbb", nil
		default:
			return "", errors.New("unexpected content")
		}
	}

	service := NewService(mockRepo, func(baseURL, apiToken string) memos.Client {
		if baseURL != "https://memos.example.com" {
			t.Fatalf("baseURL = %q, want https://memos.example.com", baseURL)
		}
		if apiToken != "token" {
			t.Fatalf("apiToken = %q, want token", apiToken)
		}
		return mockClient
	}, nil)

	result, err := service.Execute("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "対象body件数=2") ||
		!strings.Contains(result, "メモ作成成功=2") ||
		!strings.Contains(result, "添付ファイル総数=2") ||
		!strings.Contains(result, "添付スキップ(リソースなし)=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	if len(mockClient.CreateMemoCalls) != 2 {
		t.Fatalf("CreateMemo call count = %d, want 2", len(mockClient.CreateMemoCalls))
	}
	if len(mockClient.PatchFilesCalls) != 1 {
		t.Fatalf("PatchFiles call count = %d, want 1", len(mockClient.PatchFilesCalls))
	}
	if mockClient.PatchFilesCalls[0].Memo != "memos/aaa" {
		t.Fatalf("PatchFiles memo = %s, want memos/aaa", mockClient.PatchFilesCalls[0].Memo)
	}
	if len(mockClient.PatchFilesCalls[0].FilePaths) != 2 {
		t.Fatalf("PatchFiles files count = %d, want 2", len(mockClient.PatchFilesCalls[0].FilePaths))
	}
}

func TestService_Execute_Error(t *testing.T) {
	t.Parallel()

	t.Run("unsupported page type", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{}, func(baseURL, apiToken string) memos.Client {
			return &memos.MockClient{}
		}, nil)
		_, err := service.Execute("artifact", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
		if err == nil || err.Error() != "未対応のpage-typeです: artifact" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("missing base url", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{}, nil, nil)
		_, err := service.Execute("content", " ", "token", "/tmp/body", "/tmp/resources")
		if err == nil || err.Error() != "base-url パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("missing api token", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{}, nil, nil)
		_, err := service.Execute("content", "https://memos.example.com", " ", "/tmp/body", "/tmp/resources")
		if err == nil || err.Error() != "api-token パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("missing src body dir", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{}, nil, nil)
		_, err := service.Execute("content", "https://memos.example.com", "token", " ", "/tmp/resources")
		if err == nil || err.Error() != "src-body-dir パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("missing src resource dir", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{}, nil, nil)
		_, err := service.Execute("content", "https://memos.example.com", "token", "/tmp/body", " ")
		if err == nil || err.Error() != "src-resource-dir パラメータは必須です" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read body dir error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			return nil, errors.New("read body failed")
		}

		service := NewService(mockRepo, nil, nil)
		_, err := service.Execute("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
		if err == nil || !strings.Contains(err.Error(), "src-body-dir の読み取りに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read resource dir error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			if dirPath == "/tmp/body" {
				return []string{"/tmp/body/CO000000001.md"}, nil
			}
			return nil, errors.New("read resources failed")
		}

		service := NewService(mockRepo, nil, nil)
		_, err := service.Execute("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
		if err == nil || !strings.Contains(err.Error(), "src-resource-dir の読み取りに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid con id path", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			if dirPath == "/tmp/body" {
				return []string{"/tmp/body/.md"}, nil
			}
			return []string{}, nil
		}

		service := NewService(mockRepo, nil, nil)
		_, err := service.Execute("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
		if err == nil || !strings.Contains(err.Error(), "con_id の抽出に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read body file error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			if dirPath == "/tmp/body" {
				return []string{"/tmp/body/CO000000001.md"}, nil
			}
			return []string{}, nil
		}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return nil, errors.New("read file failed")
		}

		service := NewService(mockRepo, nil, nil)
		_, err := service.Execute("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
		if err == nil || !strings.Contains(err.Error(), "body ファイルの読み込みに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("create memo error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			if dirPath == "/tmp/body" {
				return []string{"/tmp/body/CO000000001.md"}, nil
			}
			return []string{}, nil
		}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte("memo"), nil
		}
		mockClient := &memos.MockClient{
			CreateMemoFunc: func(_ context.Context, content string) (string, error) {
				return "", errors.New("create failed")
			},
		}

		service := NewService(mockRepo, func(baseURL, apiToken string) memos.Client {
			return mockClient
		}, nil)
		_, err := service.Execute("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
		if err == nil || !strings.Contains(err.Error(), "メモ作成に失敗しました (con_id=CO000000001)") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("patch files error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ListFilesRecursiveFunc = func(dirPath string) ([]string, error) {
			if dirPath == "/tmp/body" {
				return []string{"/tmp/body/CO000000001.md"}, nil
			}
			return []string{"/tmp/resources/CO000000001_1.png"}, nil
		}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte("memo"), nil
		}
		mockClient := &memos.MockClient{
			CreateMemoFunc: func(_ context.Context, content string) (string, error) {
				return "memos/aaa", nil
			},
			PatchFilesFunc: func(_ context.Context, memo string, filePaths []string) error {
				return errors.New("patch failed")
			},
		}

		service := NewService(mockRepo, func(baseURL, apiToken string) memos.Client {
			return mockClient
		}, nil)
		_, err := service.Execute("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
		if err == nil || !strings.Contains(err.Error(), "添付ファイルの登録に失敗しました (con_id=CO000000001)") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestCollectResourceFilesByConID(t *testing.T) {
	t.Parallel()

	resourceFiles := []string{
		filepath.Join("/tmp/resources", "CO0001_a.png"),
		filepath.Join("/tmp/resources", "CO0001_b.png"),
		filepath.Join("/tmp/resources", "CO0002_a.png"),
	}

	got := collectResourceFilesByConID(resourceFiles, "CO0001")
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
}
