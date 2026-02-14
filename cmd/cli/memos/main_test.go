package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestRun_ListMemosWithFilter_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := run([]string{
		"-operation=list-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-page-token=next-token",
		"-state=NORMAL",
		"-order-by=update_time desc",
		`-filter=created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`,
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemosFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string) (*usecases.ListMemosOutput, error) {
				called = true
				if pageSize != 10 {
					t.Fatalf("pageSize = %d, want 10", pageSize)
				}
				if pageToken != "next-token" {
					t.Fatalf("pageToken = %s, want next-token", pageToken)
				}
				if state != "NORMAL" {
					t.Fatalf("state = %s, want NORMAL", state)
				}
				if orderBy != "update_time desc" {
					t.Fatalf("orderBy = %s, want update_time desc", orderBy)
				}
				wantFilter := `created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`
				if filter != wantFilter {
					t.Fatalf("filter = %q, want %q", filter, wantFilter)
				}
				return &usecases.ListMemosOutput{
					Memos: []usecases.Memo{
						{Name: "memos/1"},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemosFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/1\"") {
		t.Fatalf("stdout = %s, want memo json", stdout.String())
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
			UpdateMemoFunc: func(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*usecases.Memo, error) {
				if strings.Join(updateMask, ",") != "content,visibility" {
					t.Fatalf("updateMask = %v, want [content visibility]", updateMask)
				}
				if displayTime != "" {
					t.Fatalf("displayTime = %s, want empty", displayTime)
				}
				if contentFile != "" {
					t.Fatalf("contentFile = %s, want empty", contentFile)
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

func TestRun_UpdateMemo_UpdatesTime_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{
		"-operation=update-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-content=updated",
		"-updates-time=true",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			UpdateMemoFunc: func(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*usecases.Memo, error) {
				if displayTime == "" {
					t.Fatal("displayTime = empty, want RFC3339 time")
				}
				if _, err := time.Parse(time.RFC3339, displayTime); err != nil {
					t.Fatalf("displayTime parse error = %v, displayTime=%s", err, displayTime)
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

func TestRun_PatchFiles_ReplacesTrue_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var called bool
	exitCode := run([]string{
		"-operation=patch-files",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-files=./a.txt,./b.png",
		"-replaces=true",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			PatchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*usecases.SetMemoAttachmentsOutput, error) {
				called = true
				if memo != "memo-1" {
					t.Fatalf("memo = %s, want memo-1", memo)
				}
				if len(filePaths) != 2 || filePaths[0] != "./a.txt" || filePaths[1] != "./b.png" {
					t.Fatalf("filePaths = %v, want [./a.txt ./b.png]", filePaths)
				}
				if !replaces {
					t.Fatalf("replaces = %v, want true", replaces)
				}
				return &usecases.SetMemoAttachmentsOutput{
					Name:        "memos/memo-1",
					Attachments: []usecases.Attachment{{Name: "attachments/1"}},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("PatchFilesFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"name\": \"memos/memo-1\"") {
		t.Fatalf("stdout = %s, want set result", stdout.String())
	}
}

func TestRun_PatchFiles_ReplacesFalse_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var called bool
	exitCode := run([]string{
		"-operation=patch-files",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-files=./new.txt",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			PatchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*usecases.SetMemoAttachmentsOutput, error) {
				called = true
				if replaces {
					t.Fatalf("replaces = %v, want false", replaces)
				}
				if len(filePaths) != 1 || filePaths[0] != "./new.txt" {
					t.Fatalf("filePaths = %v, want [./new.txt]", filePaths)
				}
				return &usecases.SetMemoAttachmentsOutput{
					Name:        "memos/memo-1",
					Attachments: []usecases.Attachment{{Name: "attachments/new"}},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("PatchFilesFunc was not called")
	}
}

func TestRun_PatchFiles_ServiceError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{
		"-operation=patch-files",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-files=./a,./b",
		"-replaces=true",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			PatchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*usecases.SetMemoAttachmentsOutput, error) {
				return nil, errors.New("patch failed")
			},
		}
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "patch failed") {
		t.Fatalf("stderr = %s, want patch failed", stderr.String())
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
