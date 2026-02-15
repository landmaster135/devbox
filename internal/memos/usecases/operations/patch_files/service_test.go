package patchfiles

import (
	"context"
	"fmt"
	"strings"
	"testing"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	"github.com/landmaster135/devbox/internal/memos/usecases/common"
	"github.com/landmaster135/devbox/internal/memos/usecases/testutil"
)

type mockAttachmentCreator struct {
	createFunc func(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*common.Attachment, error)
}

func (m *mockAttachmentCreator) Create(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*common.Attachment, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, filename, content, attachmentType, memo)
	}
	return nil, nil
}

type mockAttachmentLister struct {
	listFunc func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoAttachmentsOutput, error)
}

func (m *mockAttachmentLister) List(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoAttachmentsOutput, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, memo, pageSize, pageToken)
	}
	return nil, nil
}

type mockAttachmentSetter struct {
	setFunc func(ctx context.Context, memo string, attachments []common.Attachment) (*common.SetMemoAttachmentsOutput, error)
}

func (m *mockAttachmentSetter) Set(ctx context.Context, memo string, attachments []common.Attachment) (*common.SetMemoAttachmentsOutput, error) {
	if m.setFunc != nil {
		return m.setFunc(ctx, memo, attachments)
	}
	return nil, nil
}

func TestServiceOperationPatchFiles_ReplacesTrue_Normal(t *testing.T) {
	createCalled := 0
	setCalled := 0

	fileSystem := &testutil.MockFileSystem{
		ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
			if filePath != "./a.txt" {
				t.Fatalf("filePath = %s, want ./a.txt", filePath)
			}
			return &infrastructures.AttachmentFile{
				Filename:    "a.txt",
				Content:     []byte("hello"),
				ContentType: "text/plain",
			}, nil
		},
	}

	service := New(
		fileSystem,
		&mockAttachmentCreator{createFunc: func(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*common.Attachment, error) {
			createCalled++
			if filename != "a.txt" {
				t.Fatalf("filename = %s, want a.txt", filename)
			}
			if attachmentType != "text/plain" {
				t.Fatalf("type = %s, want text/plain", attachmentType)
			}
			return &common.Attachment{Name: "attachments/new-1", Filename: filename}, nil
		}},
		&mockAttachmentLister{listFunc: func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoAttachmentsOutput, error) {
			t.Fatal("List should not be called when replaces=true")
			return nil, nil
		}},
		&mockAttachmentSetter{setFunc: func(ctx context.Context, memo string, attachments []common.Attachment) (*common.SetMemoAttachmentsOutput, error) {
			setCalled++
			if memo != "memo-1" {
				t.Fatalf("memo = %s, want memo-1", memo)
			}
			if len(attachments) != 1 || attachments[0].Name != "attachments/new-1" {
				t.Fatalf("attachments = %+v, want one new attachment", attachments)
			}
			return &common.SetMemoAttachmentsOutput{Name: "memos/memo-1", Attachments: attachments}, nil
		}},
	)

	result, err := service.Execute(context.Background(), "memo-1", []string{"./a.txt"}, true)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if createCalled != 1 {
		t.Fatalf("createCalled = %d, want 1", createCalled)
	}
	if setCalled != 1 {
		t.Fatalf("setCalled = %d, want 1", setCalled)
	}
	if result.Name != "memos/memo-1" {
		t.Fatalf("name = %s, want memos/memo-1", result.Name)
	}
}

func TestServiceOperationPatchFiles_MarkdownFile_Normal(t *testing.T) {
	createCalled := 0
	setCalled := 0

	service := New(
		&testutil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				if filePath != "./memo.md" {
					t.Fatalf("filePath = %s, want ./memo.md", filePath)
				}
				return &infrastructures.AttachmentFile{
					Filename:    "memo.md",
					Content:     []byte("# hello"),
					ContentType: "text/markdown",
				}, nil
			},
		},
		&mockAttachmentCreator{createFunc: func(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*common.Attachment, error) {
			createCalled++
			if attachmentType != "text/markdown" {
				t.Fatalf("type = %s, want text/markdown", attachmentType)
			}
			return &common.Attachment{Name: "attachments/new-md", Filename: filename}, nil
		}},
		&mockAttachmentLister{},
		&mockAttachmentSetter{setFunc: func(ctx context.Context, memo string, attachments []common.Attachment) (*common.SetMemoAttachmentsOutput, error) {
			setCalled++
			return &common.SetMemoAttachmentsOutput{Name: "memos/memo-1", Attachments: attachments}, nil
		}},
	)

	_, err := service.Execute(context.Background(), "memo-1", []string{"./memo.md"}, true)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if createCalled != 1 {
		t.Fatalf("createCalled = %d, want 1", createCalled)
	}
	if setCalled != 1 {
		t.Fatalf("setCalled = %d, want 1", setCalled)
	}
}

func TestServiceOperationPatchFiles_ReplacesFalse_Normal(t *testing.T) {
	listCalled := 0

	service := New(
		&testutil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				return &infrastructures.AttachmentFile{
					Filename:    "a.txt",
					Content:     []byte("hello"),
					ContentType: "text/plain",
				}, nil
			},
		},
		&mockAttachmentCreator{createFunc: func(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*common.Attachment, error) {
			return &common.Attachment{Name: "attachments/new", Filename: filename}, nil
		}},
		&mockAttachmentLister{listFunc: func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoAttachmentsOutput, error) {
			listCalled++
			if pageSize != 100 {
				t.Fatalf("pageSize = %d, want 100", pageSize)
			}
			switch listCalled {
			case 1:
				if pageToken != "" {
					t.Fatalf("pageToken = %s, want empty", pageToken)
				}
				return &common.ListMemoAttachmentsOutput{
					Attachments:   []common.Attachment{{Name: "attachments/existing-1"}},
					NextPageToken: "next",
				}, nil
			case 2:
				if pageToken != "next" {
					t.Fatalf("pageToken = %s, want next", pageToken)
				}
				return &common.ListMemoAttachmentsOutput{
					Attachments: []common.Attachment{{Name: "attachments/new"}, {Name: "attachments/existing-2"}},
				}, nil
			default:
				t.Fatalf("unexpected list call: %d", listCalled)
				return nil, nil
			}
		}},
		&mockAttachmentSetter{setFunc: func(ctx context.Context, memo string, attachments []common.Attachment) (*common.SetMemoAttachmentsOutput, error) {
			if len(attachments) != 3 {
				t.Fatalf("len(attachments) = %d, want 3", len(attachments))
			}
			if attachments[0].Name != "attachments/existing-1" || attachments[1].Name != "attachments/existing-2" || attachments[2].Name != "attachments/new" {
				t.Fatalf("attachments = %+v, want existing-1,existing-2,new", attachments)
			}
			return &common.SetMemoAttachmentsOutput{Name: "memos/memo-1", Attachments: attachments}, nil
		}},
	)

	result, err := service.Execute(context.Background(), "memo-1", []string{"./a.txt"}, false)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if listCalled != 2 {
		t.Fatalf("listCalled = %d, want 2", listCalled)
	}
	if result.Name != "memos/memo-1" {
		t.Fatalf("name = %s, want memos/memo-1", result.Name)
	}
}

func TestServiceOperationPatchFiles_ReadAttachmentError_Error(t *testing.T) {
	service := New(
		&testutil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				return nil, fmt.Errorf("read failed")
			},
		},
		&mockAttachmentCreator{},
		&mockAttachmentLister{},
		&mockAttachmentSetter{},
	)

	_, err := service.Execute(context.Background(), "memo-1", []string{"./a.txt"}, true)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if got := err.Error(); !strings.Contains(got, "files の読み込みに失敗しました") {
		t.Fatalf("error = %v, want files の読み込みに失敗しました", err)
	}
}

func TestServiceOperationPatchFiles_MIMEFailure_AbortWithoutUpload_Error(t *testing.T) {
	createCalled := 0
	setCalled := 0

	service := New(
		&testutil.MockFileSystem{
			ReadAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
				switch filePath {
				case "./a.json":
					return &infrastructures.AttachmentFile{
						Filename:    "a.json",
						Content:     []byte(`{"ok":true}`),
						ContentType: "application/json",
					}, nil
				case "./memo":
					return nil, fmt.Errorf("MIME type の判定に失敗しました (./memo): MIME type format が不正です: text/plain; charset=utf-8")
				default:
					return nil, fmt.Errorf("unexpected path: %s", filePath)
				}
			},
		},
		&mockAttachmentCreator{createFunc: func(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*common.Attachment, error) {
			createCalled++
			return &common.Attachment{Name: "attachments/new", Filename: filename}, nil
		}},
		&mockAttachmentLister{listFunc: func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoAttachmentsOutput, error) {
			t.Fatal("List should not be called when file precheck fails")
			return nil, nil
		}},
		&mockAttachmentSetter{setFunc: func(ctx context.Context, memo string, attachments []common.Attachment) (*common.SetMemoAttachmentsOutput, error) {
			setCalled++
			return &common.SetMemoAttachmentsOutput{}, nil
		}},
	)

	_, err := service.Execute(context.Background(), "memo-1", []string{"./a.json", "./memo"}, true)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if got := err.Error(); !strings.Contains(got, "files の読み込みに失敗しました (./memo)") {
		t.Fatalf("error = %v, want files の読み込みに失敗しました (./memo)", err)
	}
	if createCalled != 0 {
		t.Fatalf("createCalled = %d, want 0", createCalled)
	}
	if setCalled != 0 {
		t.Fatalf("setCalled = %d, want 0", setCalled)
	}
}

func TestMergeAttachmentsByName_Normal(t *testing.T) {
	existing := []common.Attachment{
		{Name: "attachments/existing-1"},
		{Name: "attachments/new"},
		{Name: "attachments/existing-2"},
	}
	created := []common.Attachment{
		{Name: "attachments/new"},
	}

	got := mergeAttachmentsByName(existing, created)
	if len(got) != 3 {
		t.Fatalf("len(got) = %d, want 3", len(got))
	}
	if got[0].Name != "attachments/existing-1" || got[1].Name != "attachments/existing-2" || got[2].Name != "attachments/new" {
		t.Fatalf("merge result = %+v, want existing-1, existing-2, new", got)
	}
}
