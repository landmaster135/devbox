package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos_utility/config"
	usecases "github.com/landmaster135/devbox/internal/memos_utility/usecases"
)

type mockMemosUtilityService struct {
	createClipFunc        func(ctx context.Context, input usecases.CreateClipInput) (*usecases.CreateClipOutput, error)
	createClipsFunc       func(ctx context.Context, input usecases.CreateClipsInput) (*usecases.CreateClipsOutput, error)
	createCommonMemosFunc func(ctx context.Context, input usecases.CreateCommonMemosInput) (*usecases.CreateCommonMemosOutput, error)
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

func (m *mockMemosUtilityService) CreateCommonMemos(ctx context.Context, input usecases.CreateCommonMemosInput) (*usecases.CreateCommonMemosOutput, error) {
	if m.createCommonMemosFunc != nil {
		return m.createCommonMemosFunc(ctx, input)
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
	}, &stdout, &stderr, func(conf *cfg.Config, _ io.Writer) usecases.MemosUtilityService {
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
	}, &stdout, &stderr, func(conf *cfg.Config, _ io.Writer) usecases.MemosUtilityService {
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
	}, &stdout, &stderr, func(conf *cfg.Config, _ io.Writer) usecases.MemosUtilityService {
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

func TestRun_CreateCommonMemos_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	called := false

	exitCode := run([]string{
		"-operation=create-common-memos",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-content-dir=/tmp/common-memos",
		"-attachment-dir=/tmp/attachments",
	}, &stdout, &stderr, func(conf *cfg.Config, _ io.Writer) usecases.MemosUtilityService {
		return &mockMemosUtilityService{
			createCommonMemosFunc: func(ctx context.Context, input usecases.CreateCommonMemosInput) (*usecases.CreateCommonMemosOutput, error) {
				called = true
				if input.Operation != cfg.OperationCreateCommonMemos {
					t.Fatalf("operation = %s, want %s", input.Operation, cfg.OperationCreateCommonMemos)
				}
				if input.ContentDir != "/tmp/common-memos" {
					t.Fatalf("contentDir = %s, want /tmp/common-memos", input.ContentDir)
				}
				if input.AttachmentDir != "/tmp/attachments" {
					t.Fatalf("attachmentDir = %s, want /tmp/attachments", input.AttachmentDir)
				}
				return &usecases.CreateCommonMemosOutput{
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
		t.Fatal("CreateCommonMemos was not called")
	}
	if !strings.Contains(stdout.String(), "\"operation\": \"create-common-memos\"") {
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
	}, &stdout, &stderr, func(conf *cfg.Config, _ io.Writer) usecases.MemosUtilityService {
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

	exitCode := run([]string{"-help"}, &stdout, &bytes.Buffer{}, func(conf *cfg.Config, _ io.Writer) usecases.MemosUtilityService {
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

func TestNewCreateClipsProgressReporter_CreateClips_Normal(t *testing.T) {
	var stderr bytes.Buffer
	reporter := newCreateClipsProgressReporter(&cfg.Config{
		Operation: cfg.OperationCreateClips,
	}, &stderr)

	if reporter == nil {
		t.Fatal("reporter = nil, want non-nil")
	}

	reporter(usecases.CreateClipsProgress{
		Current:         1,
		Total:           3,
		Operation:       cfg.OperationCreateWebClip,
		ContentFile:     "/tmp/web-summary-20241225-233435-daikokuyu-event-info.md",
		AttachmentCount: 2,
	})

	got := stderr.String()
	if !strings.Contains(got, "進捗: 1/3 [create-web-clip] web-summary-20241225-233435-daikokuyu-event-info.md (attachments=2)") {
		t.Fatalf("stderr = %q, want progress line", got)
	}
}

func TestNewCreateClipsProgressReporter_NotCreateClips_Normal(t *testing.T) {
	reporter := newCreateClipsProgressReporter(&cfg.Config{
		Operation: cfg.OperationCreateWebClip,
	}, &bytes.Buffer{})

	if reporter != nil {
		t.Fatal("reporter != nil, want nil")
	}
}

func TestNewCreateCommonMemosProgressReporter_CreateCommonMemos_Normal(t *testing.T) {
	var stderr bytes.Buffer
	reporter := newCreateCommonMemosProgressReporter(&cfg.Config{
		Operation: cfg.OperationCreateCommonMemos,
	}, &stderr)

	if reporter == nil {
		t.Fatal("reporter = nil, want non-nil")
	}

	reporter(usecases.CreateClipsProgress{
		Current:         2,
		Total:           5,
		Operation:       cfg.OperationCreateCommonMemos,
		ContentFile:     "/tmp/20260316080301_03.md",
		AttachmentCount: 4,
	})

	got := stderr.String()
	if !strings.Contains(got, "進捗: 2/5 [create-common-memos] 20260316080301_03.md (attachments=4)") {
		t.Fatalf("stderr = %q, want progress line", got)
	}
}

func TestNewCreateCommonMemosProgressReporter_NotCreateCommonMemos_Normal(t *testing.T) {
	reporter := newCreateCommonMemosProgressReporter(&cfg.Config{
		Operation: cfg.OperationCreateClips,
	}, &bytes.Buffer{})

	if reporter != nil {
		t.Fatal("reporter != nil, want nil")
	}
}

func TestNewCreateCommonMemosRelationReporter_CreateCommonMemos_Normal(t *testing.T) {
	var stderr bytes.Buffer
	reporter := newCreateCommonMemosRelationReporter(&cfg.Config{
		Operation: cfg.OperationCreateCommonMemos,
	}, &stderr)

	if reporter == nil {
		t.Fatal("reporter = nil, want non-nil")
	}

	reporter(usecases.CreateCommonMemosRelationProgress{
		Phase:                        usecases.CreateCommonMemosRelationPhaseStart,
		ContentFile:                  "/tmp/20260316080301_03.md",
		CurrentMemoIdentifier:        "memos/current",
		CurrentMemoIdentifierSource:  "name",
		PreviousMemoIdentifier:       "memos/previous",
		PreviousMemoIdentifierSource: "uid",
		Attempt:                      1,
		TotalAttempts:                3,
		Retrying:                     false,
	})
	reporter(usecases.CreateCommonMemosRelationProgress{
		Phase:                  usecases.CreateCommonMemosRelationPhaseRetry,
		ContentFile:            "/tmp/20260316080301_03.md",
		CurrentMemoIdentifier:  "memos/current",
		PreviousMemoIdentifier: "memos/previous",
		Attempt:                1,
		TotalAttempts:          3,
		Retrying:               true,
		RetryAfter:             "1s",
		ErrorMessage:           "temporary failed to get memo",
	})
	reporter(usecases.CreateCommonMemosRelationProgress{
		Phase:                  usecases.CreateCommonMemosRelationPhaseOK,
		ContentFile:            "/tmp/20260316080301_03.md",
		CurrentMemoIdentifier:  "memos/current",
		PreviousMemoIdentifier: "memos/previous",
		Attempt:                2,
		TotalAttempts:          3,
	})
	reporter(usecases.CreateCommonMemosRelationProgress{
		Phase:                  usecases.CreateCommonMemosRelationPhaseError,
		ContentFile:            "/tmp/20260316080301_03.md",
		CurrentMemoIdentifier:  "memos/current",
		PreviousMemoIdentifier: "memos/previous",
		Attempt:                3,
		TotalAttempts:          3,
		ErrorMessage:           "relation failed",
	})

	got := stderr.String()
	if !strings.Contains(got, "relation: [start] 20260316080301_03.md current=memos/current(from=name) previous=memos/previous(from=uid) attempt=1/3 retrying=false") {
		t.Fatalf("stderr = %q, want relation start line", got)
	}
	if !strings.Contains(got, "relation: [retry] 20260316080301_03.md current=memos/current previous=memos/previous attempt=1/3 retry-after=1s err=temporary failed to get memo") {
		t.Fatalf("stderr = %q, want relation retry line", got)
	}
	if !strings.Contains(got, "relation: [ok] 20260316080301_03.md current=memos/current previous=memos/previous attempt=2/3") {
		t.Fatalf("stderr = %q, want relation ok line", got)
	}
	if !strings.Contains(got, "relation: [error] 20260316080301_03.md current=memos/current previous=memos/previous attempt=3/3 err=relation failed") {
		t.Fatalf("stderr = %q, want relation error line", got)
	}
}

func TestNewCreateCommonMemosRelationReporter_NotCreateCommonMemos_Normal(t *testing.T) {
	reporter := newCreateCommonMemosRelationReporter(&cfg.Config{
		Operation: cfg.OperationCreateClips,
	}, &bytes.Buffer{})

	if reporter != nil {
		t.Fatal("reporter != nil, want nil")
	}
}
