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

func TestRun_ListUserTags_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := Run([]string{
		"-operation=list-user-tags",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-user-id=1",
		"-output-dir=/tmp/out",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListUserTagsFunc: func(ctx context.Context, userID string, outputDir string) (*usecases.ListUserTagsOutput, error) {
				called = true
				if userID != "1" {
					t.Fatalf("userID = %s, want 1", userID)
				}
				if outputDir != "/tmp/out" {
					t.Fatalf("outputDir = %s, want /tmp/out", outputDir)
				}
				return &usecases.ListUserTagsOutput{
					UserID:     userID,
					OutputPath: "/tmp/out/user-tags_1.json",
					TagCount:   map[string]int{"work": 2},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListUserTagsFunc was not called")
	}
	if !strings.Contains(stdout.String(), `"outputPath": "/tmp/out/user-tags_1.json"`) {
		t.Fatalf("stdout = %s, want outputPath json", stdout.String())
	}
	if !strings.Contains(stdout.String(), `"work": 2`) {
		t.Fatalf("stdout = %s, want tagCount json", stdout.String())
	}
}

func TestRun_ListUserTagsServiceError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{
		"-operation=list-user-tags",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-user-id=1",
		"-output-dir=/tmp/out",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListUserTagsFunc: func(ctx context.Context, userID string, outputDir string) (*usecases.ListUserTagsOutput, error) {
				return nil, errors.New("service failed")
			},
		}
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "service failed") {
		t.Fatalf("stderr = %s, want service failed", stderr.String())
	}
}
