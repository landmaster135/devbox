package runner

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

func TestRun_ListMemoRelations_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var called bool
	exitCode := Run([]string{
		"-operation=list-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			ListMemoRelationsFunc: func(ctx context.Context, memo string) (*usecases.ListMemoRelationsOutput, error) {
				called = true
				if memo != "memo-1" {
					t.Fatalf("memo = %s, want memo-1", memo)
				}
				return &usecases.ListMemoRelationsOutput{
					Relations: []common.MemoRelation{
						{
							Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
							RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-2"},
							Type:        common.MemoRelationTypeReference,
						},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("ListMemoRelationsFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"relatedMemo\": {") {
		t.Fatalf("stdout = %s, want relatedMemo", stdout.String())
	}
}

func TestRun_AddMemoRelations_ReplacesFalse_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var called bool
	exitCode := Run([]string{
		"-operation=add-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-related-memos=memo-2,memo-3",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			AddMemoRelationsFunc: func(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*usecases.AddMemoRelationsOutput, error) {
				called = true
				if memo != "memo-1" {
					t.Fatalf("memo = %s, want memo-1", memo)
				}
				if replaces {
					t.Fatalf("replaces = %v, want false", replaces)
				}
				if len(relatedMemos) != 2 || relatedMemos[0] != "memo-2" || relatedMemos[1] != "memo-3" {
					t.Fatalf("relatedMemos = %v, want [memo-2 memo-3]", relatedMemos)
				}
				return &usecases.AddMemoRelationsOutput{
					Memo: "memos/memo-1",
					AddedRelations: []common.MemoRelation{
						{
							Memo:        common.MemoRelationMemo{Name: "memos/memo-1"},
							RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-2"},
						},
					},
				}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("AddMemoRelationsFunc was not called")
	}
	if !strings.Contains(stdout.String(), "\"addedRelations\": [") {
		t.Fatalf("stdout = %s, want addedRelations", stdout.String())
	}
}

func TestRun_AddMemoRelations_ReplacesTrue_Normal(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var called bool
	exitCode := Run([]string{
		"-operation=add-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-related-memos=memo-2",
		"-replaces=true",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			AddMemoRelationsFunc: func(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*usecases.AddMemoRelationsOutput, error) {
				called = true
				if !replaces {
					t.Fatalf("replaces = %v, want true", replaces)
				}
				return &usecases.AddMemoRelationsOutput{Memo: "memos/memo-1"}, nil
			},
		}
	})

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0, stderr=%s", exitCode, stderr.String())
	}
	if !called {
		t.Fatal("AddMemoRelationsFunc was not called")
	}
}

func TestRun_AddMemoRelations_ServiceError_Error(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run([]string{
		"-operation=add-memo-relations",
		"-base-url=https://memos.example.com",
		"-api-token=test-token",
		"-memo=memo-1",
		"-related-memos=memo-2",
	}, &stdout, &stderr, func(conf *cfg.Config) usecases.MemoService {
		return &usecases.MockMemoService{
			AddMemoRelationsFunc: func(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*usecases.AddMemoRelationsOutput, error) {
				return nil, errors.New("add relations failed")
			},
		}
	})

	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr.String(), "add relations failed") {
		t.Fatalf("stderr = %s, want add relations failed", stderr.String())
	}
}
