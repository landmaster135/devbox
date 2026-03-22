package addmemorelations

import (
	"context"
	"errors"
	"testing"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

type mockRelationLister struct {
	listFunc func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoRelationsOutput, error)
}

func (m *mockRelationLister) List(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoRelationsOutput, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx, memo, pageSize, pageToken)
	}
	return nil, nil
}

type mockRelationSetter struct {
	setFunc func(ctx context.Context, memo string, relations []common.MemoRelation) error
}

func (m *mockRelationSetter) Set(ctx context.Context, memo string, relations []common.MemoRelation) error {
	if m.setFunc != nil {
		return m.setFunc(ctx, memo, relations)
	}
	return nil
}

func TestServiceOperationAddMemoRelations_ReplacesFalse_Normal(t *testing.T) {
	listCalled := 0
	setCalled := 0

	service := New(
		&mockRelationLister{
			listFunc: func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoRelationsOutput, error) {
				listCalled++
				if memo != "memos/memo-1" {
					t.Fatalf("memo = %s, want memos/memo-1", memo)
				}
				if pageSize != defaultListPageSize {
					t.Fatalf("pageSize = %d, want %d", pageSize, defaultListPageSize)
				}
				switch listCalled {
				case 1:
					if pageToken != "" {
						t.Fatalf("pageToken = %s, want empty", pageToken)
					}
					return &common.ListMemoRelationsOutput{
						Relations: []common.MemoRelation{
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-2"},
							},
						},
						NextPageToken: "next",
					}, nil
				case 2:
					if pageToken != "next" {
						t.Fatalf("pageToken = %s, want next", pageToken)
					}
					return &common.ListMemoRelationsOutput{
						Relations: []common.MemoRelation{
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-9"},
								Type:        common.MemoRelationTypeUnspecified,
							},
						},
					}, nil
				case 3:
					if pageToken != "" {
						t.Fatalf("pageToken = %s, want empty", pageToken)
					}
					return &common.ListMemoRelationsOutput{
						Relations: []common.MemoRelation{
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-2"},
							},
						},
						NextPageToken: "next",
					}, nil
				case 4:
					if pageToken != "next" {
						t.Fatalf("pageToken = %s, want next", pageToken)
					}
					return &common.ListMemoRelationsOutput{
						Relations: []common.MemoRelation{
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-9"},
								Type:        common.MemoRelationTypeUnspecified,
							},
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-3"},
								Type:        common.MemoRelationTypeReference,
							},
						},
					}, nil
				default:
					t.Fatalf("unexpected list call: %d", listCalled)
					return nil, nil
				}
			},
		},
		&mockRelationSetter{
			setFunc: func(ctx context.Context, memo string, relations []common.MemoRelation) error {
				setCalled++
				if memo != "memos/memo-1" {
					t.Fatalf("memo = %s, want memos/memo-1", memo)
				}
				if len(relations) != 3 {
					t.Fatalf("len(relations) = %d, want 3", len(relations))
				}
				if relations[1].Type != common.MemoRelationTypeUnspecified {
					t.Fatalf("relations[1].type = %s, want %s", relations[1].Type, common.MemoRelationTypeUnspecified)
				}
				if relations[2].Type != common.MemoRelationTypeReference {
					t.Fatalf("relations[2].type = %s, want %s", relations[2].Type, common.MemoRelationTypeReference)
				}
				return nil
			},
		},
	)

	result, err := service.Execute(
		context.Background(),
		"memo-1",
		[]string{"memo-2", "memos/memo-3", "memos/memo-3", "memo-1"},
		false,
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if listCalled != 4 {
		t.Fatalf("listCalled = %d, want 4", listCalled)
	}
	if setCalled != 1 {
		t.Fatalf("setCalled = %d, want 1", setCalled)
	}
	if result.Memo != "memos/memo-1" {
		t.Fatalf("memo = %s, want memos/memo-1", result.Memo)
	}
	if len(result.DiscardedRelations) != 0 {
		t.Fatalf("discarded = %+v, want empty", result.DiscardedRelations)
	}
	if len(result.AddedRelations) != 1 || result.AddedRelations[0].RelatedMemo.Name != "memos/memo-3" {
		t.Fatalf("added = %+v, want only memos/memo-3", result.AddedRelations)
	}
	if len(result.FinalRelations) != 3 {
		t.Fatalf("len(final) = %d, want 3", len(result.FinalRelations))
	}
}

func TestServiceOperationAddMemoRelations_ReplacesTrue_Normal(t *testing.T) {
	listCalled := 0

	service := New(
		&mockRelationLister{
			listFunc: func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoRelationsOutput, error) {
				listCalled++
				switch listCalled {
				case 1:
					return &common.ListMemoRelationsOutput{
						Relations: []common.MemoRelation{
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-2"},
							},
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-9"},
							},
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-x"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-1"},
							},
						},
					}, nil
				case 2:
					return &common.ListMemoRelationsOutput{
						Relations: []common.MemoRelation{
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-2"},
								Type:        common.MemoRelationTypeReference,
							},
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-3"},
								Type:        common.MemoRelationTypeReference,
							},
							{
								Memo:        common.MemoRelationMemo{Name: "memos/memo-x"},
								RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-1"},
								Type:        common.MemoRelationTypeReference,
							},
						},
					}, nil
				default:
					t.Fatalf("unexpected list call: %d", listCalled)
					return nil, nil
				}
			},
		},
		&mockRelationSetter{
			setFunc: func(ctx context.Context, memo string, relations []common.MemoRelation) error {
				if len(relations) != 2 {
					t.Fatalf("len(relations) = %d, want 2", len(relations))
				}
				if relations[0].RelatedMemo.Name != "memos/memo-2" {
					t.Fatalf("relations[0] = %+v, want memo-2", relations[0])
				}
				if relations[1].RelatedMemo.Name != "memos/memo-3" {
					t.Fatalf("relations[1] = %+v, want memo-3", relations[1])
				}
				if relations[0].Type != common.MemoRelationTypeReference || relations[1].Type != common.MemoRelationTypeReference {
					t.Fatalf("relation types = [%s, %s], want both %s", relations[0].Type, relations[1].Type, common.MemoRelationTypeReference)
				}
				return nil
			},
		},
	)

	result, err := service.Execute(context.Background(), "memo-1", []string{"memo-2", "memo-3"}, true)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if listCalled != 2 {
		t.Fatalf("listCalled = %d, want 2", listCalled)
	}
	if len(result.DiscardedRelations) != 1 || result.DiscardedRelations[0].RelatedMemo.Name != "memos/memo-9" {
		t.Fatalf("discarded = %+v, want only memo-9", result.DiscardedRelations)
	}
	if len(result.AddedRelations) != 1 || result.AddedRelations[0].RelatedMemo.Name != "memos/memo-3" {
		t.Fatalf("added = %+v, want only memo-3", result.AddedRelations)
	}
	if len(result.FinalRelations) != 3 {
		t.Fatalf("len(final) = %d, want 3", len(result.FinalRelations))
	}
	if result.FinalRelations[2].Memo.Name != "memos/memo-x" || result.FinalRelations[2].RelatedMemo.Name != "memos/memo-1" {
		t.Fatalf("final[2] = %+v, want inbound relation to remain", result.FinalRelations[2])
	}
}

func TestServiceOperationAddMemoRelations_EmptyRelatedMemos_Error(t *testing.T) {
	setCalled := 0
	service := New(
		&mockRelationLister{
			listFunc: func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoRelationsOutput, error) {
				t.Fatal("List should not be called when related-memos is empty")
				return nil, nil
			},
		},
		&mockRelationSetter{
			setFunc: func(ctx context.Context, memo string, relations []common.MemoRelation) error {
				setCalled++
				return nil
			},
		},
	)

	_, err := service.Execute(context.Background(), "memo-1", []string{" ", "memos/memo-1"}, false)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if err.Error() != "related-memos が空です" {
		t.Fatalf("error = %v, want related-memos が空です", err)
	}
	if setCalled != 0 {
		t.Fatalf("setCalled = %d, want 0", setCalled)
	}
}

func TestServiceOperationAddMemoRelations_ListError_Error(t *testing.T) {
	service := New(
		&mockRelationLister{
			listFunc: func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoRelationsOutput, error) {
				return nil, errors.New("list failed")
			},
		},
		&mockRelationSetter{},
	)

	_, err := service.Execute(context.Background(), "memo-1", []string{"memo-2"}, false)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if err.Error() != "list failed" {
		t.Fatalf("error = %v, want list failed", err)
	}
}

func TestServiceOperationAddMemoRelations_ListAfterSetError_Error(t *testing.T) {
	listCalled := 0
	service := New(
		&mockRelationLister{
			listFunc: func(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoRelationsOutput, error) {
				listCalled++
				if listCalled == 1 {
					return &common.ListMemoRelationsOutput{}, nil
				}
				return nil, errors.New("list after set failed")
			},
		},
		&mockRelationSetter{
			setFunc: func(ctx context.Context, memo string, relations []common.MemoRelation) error {
				return nil
			},
		},
	)

	_, err := service.Execute(context.Background(), "memo-1", []string{"memo-2"}, true)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if err.Error() != "list after set failed" {
		t.Fatalf("error = %v, want list after set failed", err)
	}
}
