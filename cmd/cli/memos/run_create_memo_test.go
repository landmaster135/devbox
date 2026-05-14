package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestRun_CreateMemo_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := run([]string{
		"-operation=create-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content=hello",
		"-visibility=PRIVATE",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*usecases.Memo, error) {
				called = true
				if content != "hello" {
					t.Fatalf("content = %s, want hello", content)
				}
				if contentFile != "" {
					t.Fatalf("contentFile = %s, want empty", contentFile)
				}
				return &usecases.Memo{Name: "memos/1", Content: content}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("createMemoFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/1\"") {
		t.Fatalf("stdout = %s, want memo json", stdout.String())
	}
}

func TestRun_CreateMemoWithContentFile_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	tempDir := t.TempDir()
	contentPath := filepath.Join(tempDir, "memo.md")
	fileContent := "# title\n\n- [ ] task\n"
	if err := os.WriteFile(contentPath, []byte(fileContent), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	exitCode := run([]string{
		"-operation=create-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=" + contentPath,
		"-visibility=PRIVATE",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*usecases.Memo, error) {
				called = true
				if content != "" {
					t.Fatalf("content = %q, want empty", content)
				}
				if contentFile != contentPath {
					t.Fatalf("contentFile = %q, want %q", contentFile, contentPath)
				}
				return &usecases.Memo{Name: "memos/2", Content: fileContent}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("CreateMemoFunc was not called")
	}
}

func TestRun_CreateMemoWithMissingContentFile_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{
		"-operation=create-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=/tmp/not-found-memo.md",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*usecases.Memo, error) {
				return nil, errors.New("content-file の読み込みに失敗しました")
			},
		}
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "content-file") {
		t.Fatalf("stderr = %s, want content-file", stderr.String())
	}
}
