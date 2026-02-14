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
	if service.updateMemoOp == nil {
		t.Fatal("updateMemoOp is nil")
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
