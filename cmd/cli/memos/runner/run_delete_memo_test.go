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

func TestRun_DeleteMemo_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var called bool
	exitCode := Run([]string{
		"-operation=delete-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-force=true",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			DeleteMemoFunc: func(ctx context.Context, memo string, force bool) (*usecases.DeleteMemoOutput, error) {
				called = true
				if memo != "memo-1" {
					t.Fatalf("memo = %s, want memo-1", memo)
				}
				if !force {
					t.Fatalf("force = %v, want true", force)
				}
				return &usecases.DeleteMemoOutput{}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("DeleteMemoFunc was not called")
	}
	if strings.TrimSpace(stdout.String()) != "{}" {
		t.Fatalf("stdout = %q, want {}", stdout.String())
	}
}

func TestRun_DeleteMemo_ServiceError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{
		"-operation=delete-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			DeleteMemoFunc: func(ctx context.Context, memo string, force bool) (*usecases.DeleteMemoOutput, error) {
				return nil, errors.New("delete failed")
			},
		}
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "delete failed") {
		t.Fatalf("stderr = %s, want delete failed", stderr.String())
	}
}

func TestRun_DeleteMemo_DefaultForceFalse_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{
		"-operation=delete-memo",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			DeleteMemoFunc: func(ctx context.Context, memo string, force bool) (*usecases.DeleteMemoOutput, error) {
				if force {
					t.Fatalf("force = %v, want false", force)
				}
				return &usecases.DeleteMemoOutput{}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
}
