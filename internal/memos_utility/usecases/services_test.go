package usecases

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
	"github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
)

type mockCreateClipOperation struct {
	executeFunc func(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error)
}

func (m *mockCreateClipOperation) Execute(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, input)
	}
	return nil, nil
}

type mockCreateClipsOperation struct {
	executeFunc func(ctx context.Context, input common.CreateClipsInput) (*common.CreateClipsOutput, error)
}

func (m *mockCreateClipsOperation) Execute(ctx context.Context, input common.CreateClipsInput) (*common.CreateClipsOutput, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, input)
	}
	return nil, nil
}

type mockCreateCommonMemosOperation struct {
	executeFunc func(ctx context.Context, input common.CreateCommonMemosInput) (*common.CreateCommonMemosOutput, error)
}

func (m *mockCreateCommonMemosOperation) Execute(ctx context.Context, input common.CreateCommonMemosInput) (*common.CreateCommonMemosOutput, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, input)
	}
	return nil, nil
}

func TestService_CreateClipDelegate_Normal(t *testing.T) {
	called := false
	svc := &Service{
		createClipOp: &mockCreateClipOperation{
			executeFunc: func(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error) {
				called = true
				if input.Operation != common.OperationCreateWebClip {
					t.Fatalf("operation = %s, want %s", input.Operation, common.OperationCreateWebClip)
				}
				return &common.CreateClipOutput{Operation: input.Operation}, nil
			},
		},
	}

	out, err := svc.CreateClip(context.Background(), CreateClipInput{Operation: common.OperationCreateWebClip})
	if err != nil {
		t.Fatalf("CreateClip() error = %v", err)
	}
	if !called {
		t.Fatal("createClipOp.Execute was not called")
	}
	if out == nil || out.Operation != common.OperationCreateWebClip {
		t.Fatalf("out = %#v, want operation %s", out, common.OperationCreateWebClip)
	}
}

func TestService_CreateClipsDelegate_Normal(t *testing.T) {
	called := false
	svc := &Service{
		createClipsOp: &mockCreateClipsOperation{
			executeFunc: func(ctx context.Context, input common.CreateClipsInput) (*common.CreateClipsOutput, error) {
				called = true
				if input.Operation != common.OperationCreateClips {
					t.Fatalf("operation = %s, want %s", input.Operation, common.OperationCreateClips)
				}
				return &common.CreateClipsOutput{Operation: input.Operation, Total: 0}, nil
			},
		},
	}

	out, err := svc.CreateClips(context.Background(), CreateClipsInput{Operation: common.OperationCreateClips})
	if err != nil {
		t.Fatalf("CreateClips() error = %v", err)
	}
	if !called {
		t.Fatal("createClipsOp.Execute was not called")
	}
	if out == nil || out.Operation != common.OperationCreateClips {
		t.Fatalf("out = %#v, want operation %s", out, common.OperationCreateClips)
	}
}

func TestService_CreateCommonMemosDelegate_Normal(t *testing.T) {
	called := false
	svc := &Service{
		createCommonMemosOp: &mockCreateCommonMemosOperation{
			executeFunc: func(ctx context.Context, input common.CreateCommonMemosInput) (*common.CreateCommonMemosOutput, error) {
				called = true
				if input.Operation != common.OperationCreateCommonMemos {
					t.Fatalf("operation = %s, want %s", input.Operation, common.OperationCreateCommonMemos)
				}
				return &common.CreateCommonMemosOutput{Operation: input.Operation, Total: 0}, nil
			},
		},
	}

	out, err := svc.CreateCommonMemos(context.Background(), CreateCommonMemosInput{Operation: common.OperationCreateCommonMemos})
	if err != nil {
		t.Fatalf("CreateCommonMemos() error = %v", err)
	}
	if !called {
		t.Fatal("createCommonMemosOp.Execute was not called")
	}
	if out == nil || out.Operation != common.OperationCreateCommonMemos {
		t.Fatalf("out = %#v, want operation %s", out, common.OperationCreateCommonMemos)
	}
}

func TestNewService_CreateClipUsesInjectedMemosService_Normal(t *testing.T) {
	createMemoCalled := false
	svc := NewService(ServiceOptions{
		MemosService: &MockMemosService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				createMemoCalled = true
				return &memos.Memo{Name: "memos/1"}, nil
			},
		},
		FileSystem: infrastructures.NewOSFileSystem(),
	})

	_, err := svc.CreateClip(context.Background(), CreateClipInput{
		Operation:   common.OperationCreateWebClip,
		ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
	})
	if err != nil {
		t.Fatalf("CreateClip() error = %v", err)
	}
	if !createMemoCalled {
		t.Fatal("CreateMemo was not called")
	}
}

func TestNewService_CreateClipsProgressReporter_Normal(t *testing.T) {
	contentDir := t.TempDir()
	contentFile := filepath.Join(contentDir, "web-summary-20241225-233435-daikokuyu-event-info.md")
	if err := os.WriteFile(contentFile, []byte("# content"), 0o644); err != nil {
		t.Fatalf("WriteFile(content) error = %v", err)
	}

	progressCalled := false
	svc := NewService(ServiceOptions{
		MemosService: &MockMemosService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				return &memos.Memo{Name: "memos/1"}, nil
			},
		},
		FileSystem: infrastructures.NewOSFileSystem(),
		CreateClipsProgressReporter: func(progress CreateClipsProgress) {
			progressCalled = true
			if progress.Current != 1 || progress.Total != 1 {
				t.Fatalf("progress = %#v, want current=1 total=1", progress)
			}
			if progress.Operation != common.OperationCreateWebClip {
				t.Fatalf("operation = %s, want %s", progress.Operation, common.OperationCreateWebClip)
			}
			if filepath.Base(progress.ContentFile) != filepath.Base(contentFile) {
				t.Fatalf("contentFile = %s, want %s", progress.ContentFile, contentFile)
			}
			if progress.AttachmentCount != 0 {
				t.Fatalf("attachmentCount = %d, want 0", progress.AttachmentCount)
			}
		},
	})

	_, err := svc.CreateClips(context.Background(), CreateClipsInput{
		Operation:  common.OperationCreateClips,
		ContentDir: contentDir,
	})
	if err != nil {
		t.Fatalf("CreateClips() error = %v", err)
	}
	if !progressCalled {
		t.Fatal("CreateClipsProgressReporter was not called")
	}
}

func TestNewService_CreateCommonMemosProgressReporter_Normal(t *testing.T) {
	contentDir := t.TempDir()
	contentFile := filepath.Join(contentDir, "20260316080301_01.md")
	if err := os.WriteFile(contentFile, []byte("# content"), 0o644); err != nil {
		t.Fatalf("WriteFile(content) error = %v", err)
	}

	progressCalled := false
	svc := NewService(ServiceOptions{
		MemosService: &MockMemosService{
			CreateMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				return &memos.Memo{Name: "memos/1"}, nil
			},
		},
		FileSystem: infrastructures.NewOSFileSystem(),
		CreateCommonMemosProgressReporter: func(progress CreateClipsProgress) {
			progressCalled = true
			if progress.Current != 1 || progress.Total != 1 {
				t.Fatalf("progress = %#v, want current=1 total=1", progress)
			}
			if progress.Operation != common.OperationCreateCommonMemos {
				t.Fatalf("operation = %s, want %s", progress.Operation, common.OperationCreateCommonMemos)
			}
			if filepath.Base(progress.ContentFile) != filepath.Base(contentFile) {
				t.Fatalf("contentFile = %s, want %s", progress.ContentFile, contentFile)
			}
			if progress.AttachmentCount != 0 {
				t.Fatalf("attachmentCount = %d, want 0", progress.AttachmentCount)
			}
		},
	})

	_, err := svc.CreateCommonMemos(context.Background(), CreateCommonMemosInput{
		Operation:  common.OperationCreateCommonMemos,
		ContentDir: contentDir,
	})
	if err != nil {
		t.Fatalf("CreateCommonMemos() error = %v", err)
	}
	if !progressCalled {
		t.Fatal("CreateCommonMemosProgressReporter was not called")
	}
}
