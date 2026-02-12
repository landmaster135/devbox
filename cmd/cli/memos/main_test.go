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
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, visibility string, state string, pinned *bool, displayTime string) (*usecases.Memo, error) {
				called = true
				if content != "hello" {
					t.Fatalf("content = %s, want hello", content)
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
	contentFile := filepath.Join(tempDir, "memo.md")
	fileContent := "# title\n\n- [ ] task\n"
	if err := os.WriteFile(contentFile, []byte(fileContent), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	exitCode := run([]string{
		"-operation=create-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=" + contentFile,
		"-visibility=PRIVATE",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, visibility string, state string, pinned *bool, displayTime string) (*usecases.Memo, error) {
				called = true
				if content != fileContent {
					t.Fatalf("content = %q, want %q", content, fileContent)
				}
				return &usecases.Memo{Name: "memos/2", Content: content}, nil
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

func TestRun_GetMemoError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{
		"-operation=get-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			GetMemoFunc: func(ctx context.Context, memo string) (*usecases.Memo, error) {
				return nil, errors.New("get failed")
			},
		}
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "get failed") {
		t.Fatalf("stderr = %s, want get failed", stderr.String())
	}
}

func TestRun_ParseError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{
		"-operation=get-memo",
		"-base-url=https://memos.example.com",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		t.Fatal("factory should not be called")
		return nil
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "エラー") {
		t.Fatalf("stderr = %s, want エラー", stderr.String())
	}
}

func TestRun_UpdateMemo_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{
		"-operation=update-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-content=updated",
		"-update-mask=content,visibility",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			UpdateMemoFunc: func(ctx context.Context, memo string, content string, visibility string, state string, pinned *bool, updateMask []string) (*usecases.Memo, error) {
				if strings.Join(updateMask, ",") != "content,visibility" {
					t.Fatalf("updateMask = %v, want [content visibility]", updateMask)
				}
				return &usecases.Memo{Name: "memos/memo-1", Content: content}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "updated") {
		t.Fatalf("stdout = %s, want updated", stdout.String())
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
		return &usecases.MockMemoService{}
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "content-file") {
		t.Fatalf("stderr = %s, want content-file", stderr.String())
	}
}

func TestSplitByComma_Normal(t *testing.T) {
	got := splitByComma(" content, visibility ,state ")
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if got[0] != "content" || got[1] != "visibility" || got[2] != "state" {
		t.Fatalf("split result = %v, want [content visibility state]", got)
	}
}

func TestSplitByComma_Empty_Normal(t *testing.T) {
	got := splitByComma("  ")
	if got != nil {
		t.Fatalf("got = %v, want nil", got)
	}
}

func TestBoolPointer_Normal(t *testing.T) {
	if got := boolPointer(true, false); got != nil {
		t.Fatalf("got = %v, want nil", got)
	}
	got := boolPointer(false, true)
	if got == nil || *got != false {
		t.Fatalf("got = %v, want pointer to false", got)
	}
}

func TestPrintJSON_Normal(t *testing.T) {
	var out bytes.Buffer
	if err := printJSON(&out, map[string]string{"hello": "world"}); err != nil {
		t.Fatalf("printJSON() error = %v", err)
	}
	if !strings.Contains(out.String(), "\"hello\": \"world\"") {
		t.Fatalf("output = %s, want JSON content", out.String())
	}
}

func TestPrintJSON_MarshalError_Error(t *testing.T) {
	var out bytes.Buffer
	err := printJSON(&out, map[string]any{"invalid": make(chan int)})
	if err == nil {
		t.Fatal("printJSON() error = nil, want error")
	}
}

func TestResolveContent_Normal(t *testing.T) {
	got, err := resolveContent("hello", "")
	if err != nil {
		t.Fatalf("resolveContent() error = %v", err)
	}
	if got != "hello" {
		t.Fatalf("content = %s, want hello", got)
	}
}

func TestResolveContent_FromFile_Normal(t *testing.T) {
	tempDir := t.TempDir()
	contentFile := filepath.Join(tempDir, "memo.md")
	want := "# memo\n\nbody\n"
	if err := os.WriteFile(contentFile, []byte(want), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := resolveContent("", contentFile)
	if err != nil {
		t.Fatalf("resolveContent() error = %v", err)
	}
	if got != want {
		t.Fatalf("content = %q, want %q", got, want)
	}
}

func TestResolveContent_BothSpecified_Error(t *testing.T) {
	_, err := resolveContent("x", "/tmp/memo.md")
	if err == nil {
		t.Fatal("resolveContent() error = nil, want error")
	}
}

func TestResolveContent_EmptyFile_Error(t *testing.T) {
	tempDir := t.TempDir()
	contentFile := filepath.Join(tempDir, "empty.md")
	if err := os.WriteFile(contentFile, []byte("  \n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := resolveContent("", contentFile)
	if err == nil {
		t.Fatal("resolveContent() error = nil, want error")
	}
}
