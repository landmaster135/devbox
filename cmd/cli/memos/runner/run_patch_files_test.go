package runner

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestRun_PatchFiles_ReplacesTrue_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var called bool
	exitCode := Run([]string{
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
	exitCode := Run([]string{
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

	exitCode := Run([]string{
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
