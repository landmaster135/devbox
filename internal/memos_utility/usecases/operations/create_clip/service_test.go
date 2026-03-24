package createclip

import (
	"context"
	"errors"
	"strings"
	"testing"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	testUtil "github.com/landmaster135/devbox/internal/memos/infrastructures/testutil"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
	common "github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
)

type mockMemosService struct {
	createMemoFunc func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	patchFilesFunc func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
}

func (m *mockMemosService) CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
	if m.createMemoFunc != nil {
		return m.createMemoFunc(ctx, memoID, content, contentFile, visibility, state, pinned, displayTime)
	}
	return nil, nil
}

func (m *mockMemosService) PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
	if m.patchFilesFunc != nil {
		return m.patchFilesFunc(ctx, memo, filePaths, replaces)
	}
	return nil, nil
}

func TestService_ExecuteWithAttachments_Normal(t *testing.T) {
	createMemoCalled := false
	patchFilesCalled := false

	service := NewService(ServiceOptions{
		FileSystem: &testUtil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				return &infrastructures.AttachmentFile{Path: filePath}, nil
			},
		},
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				createMemoCalled = true
				if memoID != "" {
					t.Fatalf("memoID = %s, want empty", memoID)
				}
				if content != "" {
					t.Fatalf("content = %s, want empty", content)
				}
				if contentFile != "/tmp/web-summary-20240719-231059-palworld.md" {
					t.Fatalf("contentFile = %s, want expected path", contentFile)
				}
				if visibility != "PRIVATE" {
					t.Fatalf("visibility = %s, want PRIVATE", visibility)
				}
				if state != "NORMAL" {
					t.Fatalf("state = %s, want NORMAL", state)
				}
				if pinned == nil || *pinned {
					t.Fatalf("pinned = %#v, want pointer(false)", pinned)
				}
				if displayTime != "2024-07-19T23:10:59+09:00" {
					t.Fatalf("displayTime = %s, want 2024-07-19T23:10:59+09:00", displayTime)
				}
				return &memos.Memo{Name: "memos/1"}, nil
			},
			patchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
				patchFilesCalled = true
				if memo != "memos/1" {
					t.Fatalf("memo = %s, want memos/1", memo)
				}
				if len(filePaths) != 2 || filePaths[0] != "./a.png" || filePaths[1] != "./b.txt" {
					t.Fatalf("filePaths = %#v, want [./a.png ./b.txt]", filePaths)
				}
				if replaces {
					t.Fatal("replaces = true, want false")
				}
				return &memos.SetMemoAttachmentsOutput{Name: "memos/1"}, nil
			},
		},
	})

	result, err := service.Execute(context.Background(), common.CreateClipInput{
		Operation:   common.OperationCreateWebClip,
		ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
		Attachments: []string{" ./a.png ", "", "./b.txt"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !createMemoCalled {
		t.Fatal("CreateMemo was not called")
	}
	if !patchFilesCalled {
		t.Fatal("PatchFiles was not called")
	}
	if result.DisplayTime != "2024-07-19T23:10:59+09:00" {
		t.Fatalf("displayTime = %s, want 2024-07-19T23:10:59+09:00", result.DisplayTime)
	}
	if result.SetMemoAttachments == nil {
		t.Fatal("setMemoAttachments = nil, want non-nil")
	}
}

func TestService_ExecuteWithoutAttachments_Normal(t *testing.T) {
	patchFilesCalled := false

	service := NewService(ServiceOptions{
		FileSystem: &testUtil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				return &infrastructures.AttachmentFile{Path: filePath}, nil
			},
		},
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				if displayTime != "2026-03-19T05:57:16+09:00" {
					t.Fatalf("displayTime = %s, want 2026-03-19T05:57:16+09:00", displayTime)
				}
				return &memos.Memo{Name: "memos/2"}, nil
			},
			patchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
				patchFilesCalled = true
				return nil, nil
			},
		},
	})

	result, err := service.Execute(context.Background(), common.CreateClipInput{
		Operation:   common.OperationCreateMovieClip,
		ContentFile: "/tmp/movie-summary-20260319-055716-trump-masako-diplomacy.md",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if patchFilesCalled {
		t.Fatal("PatchFiles was called, want not called")
	}
	if result.Memo == nil || result.Memo.Name != "memos/2" {
		t.Fatalf("memo = %#v, want name memos/2", result.Memo)
	}
}

func TestService_ExecuteUseUIDAsIdentifier_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		FileSystem: &testUtil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				return &infrastructures.AttachmentFile{Path: filePath}, nil
			},
		},
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				return &memos.Memo{UID: "memo-uid-1"}, nil
			},
			patchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
				if memo != "memo-uid-1" {
					t.Fatalf("memo = %s, want memo-uid-1", memo)
				}
				return &memos.SetMemoAttachmentsOutput{Name: "memos/uid"}, nil
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateClipInput{
		Operation:   common.OperationCreateWebClip,
		ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
		Attachments: []string{"./a.png"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestService_ExecuteAttachmentDirectory_NoMemoCreated_Error(t *testing.T) {
	directoryPath := t.TempDir()

	service := NewService(ServiceOptions{
		FileSystem: infrastructures.NewOSFileSystem(),
		MemosService: &mockMemosService{
			createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
				t.Fatal("CreateMemo should not be called when attachment path is a directory")
				return nil, nil
			},
		},
	})

	_, err := service.Execute(context.Background(), common.CreateClipInput{
		Operation:   common.OperationCreateWebClip,
		ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
		Attachments: []string{directoryPath},
	})
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "メモは作成されませんでした") {
		t.Fatalf("error = %v, want メモは作成されませんでした", err)
	}
}

func TestService_ExecuteError_Error(t *testing.T) {
	tests := []struct {
		name       string
		input      common.CreateClipInput
		memos      *mockMemosService
		fileSystem infrastructures.FileSystem
		wantErrSub string
	}{
		{
			name: "OperationUnsupported",
			input: common.CreateClipInput{
				Operation:   "create",
				ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
			},
			memos:      &mockMemosService{},
			wantErrSub: "未対応の operation",
		},
		{
			name: "ContentFilePatternInvalid",
			input: common.CreateClipInput{
				Operation:   common.OperationCreateMovieClip,
				ContentFile: "/tmp/invalid.md",
			},
			memos:      &mockMemosService{},
			wantErrSub: "形式が不正",
		},
		{
			name: "AttachmentPrecheckFailed_NoMemoCreated",
			input: common.CreateClipInput{
				Operation:   common.OperationCreateWebClip,
				ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
				Attachments: []string{"./missing.png"},
			},
			fileSystem: &testUtil.MockFileSystem{
				ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
					return nil, errors.New("file does not exist")
				},
			},
			memos: &mockMemosService{
				createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
					t.Fatal("CreateMemo should not be called when attachment precheck fails")
					return nil, nil
				},
			},
			wantErrSub: "メモは作成されませんでした",
		},
		{
			name: "CreateMemoFailed",
			input: common.CreateClipInput{
				Operation:   common.OperationCreateWebClip,
				ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
			},
			fileSystem: &testUtil.MockFileSystem{},
			memos: &mockMemosService{
				createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
					return nil, errors.New("create failed")
				},
			},
			wantErrSub: "メモの作成に失敗しました",
		},
		{
			name: "PatchFilesFailed",
			input: common.CreateClipInput{
				Operation:   common.OperationCreateWebClip,
				ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
				Attachments: []string{"./a.png"},
			},
			fileSystem: &testUtil.MockFileSystem{
				ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
					return &infrastructures.AttachmentFile{Path: filePath}, nil
				},
			},
			memos: &mockMemosService{
				createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
					return &memos.Memo{Name: "memos/1"}, nil
				},
				patchFilesFunc: func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
					return nil, errors.New("patch failed")
				},
			},
			wantErrSub: "メモの作成には成功しましたが、添付の追加に失敗しました",
		},
		{
			name: "MemoIdentifierMissing",
			input: common.CreateClipInput{
				Operation:   common.OperationCreateWebClip,
				ContentFile: "/tmp/web-summary-20240719-231059-palworld.md",
				Attachments: []string{"./a.png"},
			},
			fileSystem: &testUtil.MockFileSystem{
				ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
					return &infrastructures.AttachmentFile{Path: filePath}, nil
				},
			},
			memos: &mockMemosService{
				createMemoFunc: func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
					return &memos.Memo{}, nil
				},
			},
			wantErrSub: "添付対象メモの識別子を取得できません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(ServiceOptions{MemosService: tt.memos, FileSystem: tt.fileSystem})
			_, err := service.Execute(context.Background(), tt.input)
			if err == nil {
				t.Fatal("Execute() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErrSub) {
				t.Fatalf("error = %v, want %s", err, tt.wantErrSub)
			}
		})
	}
}
