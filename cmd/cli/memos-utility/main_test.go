package main

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos_utility/config"
	usecases "github.com/landmaster135/devbox/internal/memos_utility/usecases"
)

type mockMemosUtilityService struct {
	createClipFunc  func(ctx context.Context, input usecases.CreateClipInput) (*usecases.CreateClipOutput, error)
	createClipsFunc func(ctx context.Context, input usecases.CreateClipsInput) (*usecases.CreateClipsOutput, error)
}

func (m *mockMemosUtilityService) CreateClip(ctx context.Context, input usecases.CreateClipInput) (*usecases.CreateClipOutput, error) {
	if m.createClipFunc != nil {
		return m.createClipFunc(ctx, input)
	}
	return nil, nil
}

func (m *mockMemosUtilityService) CreateClips(ctx context.Context, input usecases.CreateClipsInput) (*usecases.CreateClipsOutput, error) {
	if m.createClipsFunc != nil {
		return m.createClipsFunc(ctx, input)
	}
	return nil, nil
}

func TestRun_CreateWebClip_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := run([]string{
		"-operation=create-web-clip",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=/tmp/web-summary-20240719-231059-palworld.md",
		"-attachments= ./a.png , ./b.txt ",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemosUtilityService {
		return &mockMemosUtilityService{
			createClipFunc: func(ctx context.Context, input usecases.CreateClipInput) (*usecases.CreateClipOutput, error) {
				called = true
				if input.Operation != cfg.OperationCreateWebClip {
					t.Fatalf("operation = %s, want %s", input.Operation, cfg.OperationCreateWebClip)
				}
				if input.ContentFile != "/tmp/web-summary-20240719-231059-palworld.md" {
					t.Fatalf("contentFile = %s, want expected path", input.ContentFile)
				}
				if len(input.Attachments) != 2 || input.Attachments[0] != "./a.png" || input.Attachments[1] != "./b.txt" {
					t.Fatalf("attachments = %#v, want [./a.png ./b.txt]", input.Attachments)
				}
				return &usecases.CreateClipOutput{Operation: input.Operation, DisplayTime: "2024-07-19T23:10:59+09:00"}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("CreateClip was not called")
	}
	if !strings.Contains(stdout.String(), "\"operation\": \"create-web-clip\"") {
		t.Fatalf("stdout = %s, want operation json", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %s, want empty", stderr.String())
	}
}

func TestRun_ServiceError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{
		"-operation=create-movie-clip",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=/tmp/movie-summary-20260319-055716-sample.md",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemosUtilityService {
		return &mockMemosUtilityService{
			createClipFunc: func(ctx context.Context, input usecases.CreateClipInput) (*usecases.CreateClipOutput, error) {
				return nil, errors.New("create failed")
			},
		}
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "create failed") {
		t.Fatalf("stderr = %s, want create failed", stderr.String())
	}
}

func TestRun_CreateClips_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := run([]string{
		"-operation=create-clips",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-dir=/tmp/clips",
		"-attachment-dir=/tmp/attachments",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemosUtilityService {
		return &mockMemosUtilityService{
			createClipsFunc: func(ctx context.Context, input usecases.CreateClipsInput) (*usecases.CreateClipsOutput, error) {
				called = true
				if input.Operation != cfg.OperationCreateClips {
					t.Fatalf("operation = %s, want %s", input.Operation, cfg.OperationCreateClips)
				}
				if input.ContentDir != "/tmp/clips" {
					t.Fatalf("contentDir = %s, want /tmp/clips", input.ContentDir)
				}
				if input.AttachmentDir != "/tmp/attachments" {
					t.Fatalf("attachmentDir = %s, want /tmp/attachments", input.AttachmentDir)
				}
				return &usecases.CreateClipsOutput{
					Operation:  input.Operation,
					ContentDir: input.ContentDir,
					Total:      0,
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("CreateClips was not called")
	}
	if !strings.Contains(stdout.String(), "\"operation\": \"create-clips\"") {
		t.Fatalf("stdout = %s, want operation json", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %s, want empty", stderr.String())
	}
}

func TestRun_ParseError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{
		"-operation=create-web-clip",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-file=/tmp/invalid.md",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemosUtilityService {
		t.Fatal("factory should not be called")
		return nil
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "形式のみ指定できます") {
		t.Fatalf("stderr = %s, want validation error", stderr.String())
	}
}

func TestRun_Help_Normal(t *testing.T) {
	var stdout bytes.Buffer
	factoryCalled := false

	exitCode := run([]string{"-help"}, &stdout, &bytes.Buffer{}, func(conf *cfg.Config) usecases.MemosUtilityService {
		factoryCalled = true
		return &mockMemosUtilityService{}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0", exitCode)
	}
	if factoryCalled {
		t.Fatal("factory should not be called")
	}
}
