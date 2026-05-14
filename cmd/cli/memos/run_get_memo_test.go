package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

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
