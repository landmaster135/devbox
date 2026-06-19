package runner

import (
	"bytes"
	"context"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestRun_UpdateTag_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var called bool
	exitCode := Run([]string{
		"-operation=update-tag",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-src-tag=work",
		"-dest-tag=project",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			UpdateTagFunc: func(ctx context.Context, srcTag string, destTag string) (*usecases.UpdateTagOutput, error) {
				called = true
				if srcTag != "work" {
					t.Fatalf("srcTag = %s, want work", srcTag)
				}
				if destTag != "project" {
					t.Fatalf("destTag = %s, want project", destTag)
				}
				return &usecases.UpdateTagOutput{
					SourceTag:        "work",
					DestinationTag:   "project",
					MatchedCount:     2,
					UpdatedCount:     2,
					UpdatedMemoNames: []string{"memos/1", "memos/2"},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("UpdateTagFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"sourceTag\": \"work\"") {
		t.Fatalf("stdout = %s, want sourceTag", stdout.String())
	}
	if !strings.Contains(stdout.String(), "\"updatedMemoNames\": [") {
		t.Fatalf("stdout = %s, want updatedMemoNames", stdout.String())
	}
}
