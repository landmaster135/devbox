package main

import (
	"bytes"
	"context"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

func TestRun_ListAttachments_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := run([]string{
		"-operation=list-attachments",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-page-size=10",
		"-page-token=next-token",
		"-order-by=create_time desc",
		`-filter=memo == "memos/memo-1"`,
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListAttachmentsFunc: func(ctx context.Context, pageSize int, pageToken string, orderBy string, filter string) (*usecases.ListAttachmentsOutput, error) {
				called = true
				if pageSize != 10 {
					t.Fatalf("pageSize = %d, want 10", pageSize)
				}
				if pageToken != "next-token" {
					t.Fatalf("pageToken = %s, want next-token", pageToken)
				}
				if orderBy != "create_time desc" {
					t.Fatalf("orderBy = %s, want create_time desc", orderBy)
				}
				if filter != `memo == "memos/memo-1"` {
					t.Fatalf("filter = %q, want memo filter", filter)
				}
				return &usecases.ListAttachmentsOutput{
					Attachments: []usecases.Attachment{
						{Name: "attachments/1"},
					},
					TotalSize: 1,
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListAttachmentsFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"name\": \"attachments/1\"") {
		t.Fatalf("stdout = %s, want attachment json", stdout.String())
	}
}
