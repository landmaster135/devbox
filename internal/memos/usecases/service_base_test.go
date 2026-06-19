package usecases

import (
	"context"
	"testing"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
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
	if service.listMemoRelationsOp == nil {
		t.Fatal("listMemoRelationsOp is nil")
	}
	if service.addMemoRelationsOp == nil {
		t.Fatal("addMemoRelationsOp is nil")
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

type stubListMemoRelationsOperation struct{}

func (s *stubListMemoRelationsOperation) List(ctx context.Context, memo string, pageSize int, pageToken string) (*ListMemoRelationsOutput, error) {
	return &ListMemoRelationsOutput{
		Relations: []common.MemoRelation{
			{
				Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
				RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-2"},
			},
		},
	}, nil
}

func TestService_ListMemoRelations_DelegatesOperation(t *testing.T) {
	service := &Service{listMemoRelationsOp: &stubListMemoRelationsOperation{}}

	result, err := service.ListMemoRelations(context.Background(), "memos/memo-1")
	if err != nil {
		t.Fatalf("ListMemoRelations() error = %v", err)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
	if len(result.Relations) != 1 {
		t.Fatalf("len(relations) = %d, want 1", len(result.Relations))
	}
}

type stubAddMemoRelationsOperation struct{}

func (s *stubAddMemoRelationsOperation) Execute(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*AddMemoRelationsOutput, error) {
	return &AddMemoRelationsOutput{
		Memo: memo,
		AddedRelations: []common.MemoRelation{
			{
				Memo:        common.MemoRelationMemo{Name: memo},
				RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-2"},
			},
		},
	}, nil
}

func TestService_AddMemoRelations_DelegatesOperation(t *testing.T) {
	service := &Service{addMemoRelationsOp: &stubAddMemoRelationsOperation{}}

	result, err := service.AddMemoRelations(context.Background(), "memos/memo-1", []string{"memo-2"}, false)
	if err != nil {
		t.Fatalf("AddMemoRelations() error = %v", err)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
	if len(result.AddedRelations) != 1 {
		t.Fatalf("len(addedRelations) = %d, want 1", len(result.AddedRelations))
	}
}
