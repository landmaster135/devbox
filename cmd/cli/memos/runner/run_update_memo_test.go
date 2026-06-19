package runner

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestRun_UpdateMemo_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{
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

	exitCode := Run([]string{
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
