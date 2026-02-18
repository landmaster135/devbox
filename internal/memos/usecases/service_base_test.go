package usecases

import (
	"context"
	"testing"
)

func TestNewService_DefaultDependencies_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:  "https://memos.example.com",
		APIToken: "token",
		Timeout:  0,
	})

	if service.createMemoOp == nil {
		t.Fatal("createMemoOp is nil")
	}
	if service.getMemoOp == nil {
		t.Fatal("getMemoOp is nil")
	}
	if service.deleteMemoOp == nil {
		t.Fatal("deleteMemoOp is nil")
	}
	if service.listMemosOp == nil {
		t.Fatal("listMemosOp is nil")
	}
	if service.listAttachmentsOp == nil {
		t.Fatal("listAttachmentsOp is nil")
	}
	if service.updateMemoOp == nil {
		t.Fatal("updateMemoOp is nil")
	}
	if service.updateTagOp == nil {
		t.Fatal("updateTagOp is nil")
	}
	if service.patchFilesOp == nil {
		t.Fatal("patchFilesOp is nil")
	}
	if service.createAttachmentOp == nil {
		t.Fatal("createAttachmentOp is nil")
	}
	if service.listMemoAttachmentsOp == nil {
		t.Fatal("listMemoAttachmentsOp is nil")
	}
	if service.setMemoAttachmentsOp == nil {
		t.Fatal("setMemoAttachmentsOp is nil")
	}
}

type stubGetMemoOperation struct{}

func (s *stubGetMemoOperation) Execute(ctx context.Context, memo string) (*Memo, error) {
	return &Memo{Name: memo}, nil
}

func TestService_GetMemo_DelegatesOperation(t *testing.T) {
	service := &Service{getMemoOp: &stubGetMemoOperation{}}

	result, err := service.GetMemo(context.Background(), "memos/memo-1")
	if err != nil {
		t.Fatalf("GetMemo() error = %v", err)
	}
	if result.Name != "memos/memo-1" {
		t.Fatalf("name = %s, want memos/memo-1", result.Name)
	}
}

type stubDeleteMemoOperation struct{}

func (s *stubDeleteMemoOperation) Execute(ctx context.Context, memo string, force bool) (*DeleteMemoOutput, error) {
	return &DeleteMemoOutput{}, nil
}

func TestService_DeleteMemo_DelegatesOperation(t *testing.T) {
	service := &Service{deleteMemoOp: &stubDeleteMemoOperation{}}

	result, err := service.DeleteMemo(context.Background(), "memos/memo-1", true)
	if err != nil {
		t.Fatalf("DeleteMemo() error = %v", err)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
}

type stubListAttachmentsOperation struct{}

func (s *stubListAttachmentsOperation) Execute(ctx context.Context, pageSize int, pageToken string, orderBy string, filter string) (*ListAttachmentsOutput, error) {
	return &ListAttachmentsOutput{
		Attachments: []Attachment{{Name: "attachments/1"}},
		TotalSize:   1,
	}, nil
}

func TestService_ListAttachments_DelegatesOperation(t *testing.T) {
	service := &Service{listAttachmentsOp: &stubListAttachmentsOperation{}}

	result, err := service.ListAttachments(context.Background(), 20, "next-token", "create_time desc", `memo == "memos/memo-1"`)
	if err != nil {
		t.Fatalf("ListAttachments() error = %v", err)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
	if result.TotalSize != 1 {
		t.Fatalf("totalSize = %d, want 1", result.TotalSize)
	}
}

type stubUpdateTagOperation struct{}

func (s *stubUpdateTagOperation) Execute(ctx context.Context, srcTag string, destTag string) (*UpdateTagOutput, error) {
	return &UpdateTagOutput{
		SourceTag:      srcTag,
		DestinationTag: destTag,
		MatchedCount:   3,
		UpdatedCount:   2,
	}, nil
}

func TestService_UpdateTag_DelegatesOperation(t *testing.T) {
	service := &Service{updateTagOp: &stubUpdateTagOperation{}}

	result, err := service.UpdateTag(context.Background(), "work", "project")
	if err != nil {
		t.Fatalf("UpdateTag() error = %v", err)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
	if result.SourceTag != "work" {
		t.Fatalf("sourceTag = %s, want work", result.SourceTag)
	}
	if result.DestinationTag != "project" {
		t.Fatalf("destinationTag = %s, want project", result.DestinationTag)
	}
}
