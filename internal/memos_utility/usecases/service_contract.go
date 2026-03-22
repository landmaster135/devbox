package usecases

import (
	"context"

	memos "github.com/landmaster135/devbox/internal/memos/usecases"
	"github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
)

// CreateClipInput は create-web-clip / create-movie-clip の入力。
type CreateClipInput = common.CreateClipInput

// CreateClipOutput は create-web-clip / create-movie-clip の出力。
type CreateClipOutput = common.CreateClipOutput

// CreateClipsInput は create-clips の入力。
type CreateClipsInput = common.CreateClipsInput

// CreateClipsOutput は create-clips の出力。
type CreateClipsOutput = common.CreateClipsOutput

// CreateClipsProgress は create-clips の進捗通知情報。
type CreateClipsProgress = common.CreateClipsProgress

// CreateCommonMemosInput は create-common-memos の入力。
type CreateCommonMemosInput = common.CreateCommonMemosInput

// CreateCommonMemosOutput は create-common-memos の出力。
type CreateCommonMemosOutput = common.CreateCommonMemosOutput

// MemosService は memos-utility が利用する Memos サービスの契約。
type MemosService interface {
	CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
	AddMemoRelations(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error)
}

// MemosUtilityService は memos-utility の公開契約。
type MemosUtilityService interface {
	CreateClip(ctx context.Context, input CreateClipInput) (*CreateClipOutput, error)
	CreateClips(ctx context.Context, input CreateClipsInput) (*CreateClipsOutput, error)
	CreateCommonMemos(ctx context.Context, input CreateCommonMemosInput) (*CreateCommonMemosOutput, error)
}

type createClipOperation interface {
	Execute(ctx context.Context, input common.CreateClipInput) (*common.CreateClipOutput, error)
}

type createClipsOperation interface {
	Execute(ctx context.Context, input common.CreateClipsInput) (*common.CreateClipsOutput, error)
}

type createCommonMemosOperation interface {
	Execute(ctx context.Context, input common.CreateCommonMemosInput) (*common.CreateCommonMemosOutput, error)
}

// MockMemosService はテスト用モック。
type MockMemosService struct {
	CreateMemoFunc       func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	PatchFilesFunc       func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
	AddMemoRelationsFunc func(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error)
}

// MockMemosUtilityService は MemosUtilityService のテストダブル。
type MockMemosUtilityService struct {
	CreateClipFunc        func(ctx context.Context, input CreateClipInput) (*CreateClipOutput, error)
	CreateClipsFunc       func(ctx context.Context, input CreateClipsInput) (*CreateClipsOutput, error)
	CreateCommonMemosFunc func(ctx context.Context, input CreateCommonMemosInput) (*CreateCommonMemosOutput, error)
}

func (m *MockMemosService) CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error) {
	if m.CreateMemoFunc != nil {
		return m.CreateMemoFunc(ctx, memoID, content, contentFile, visibility, state, pinned, displayTime)
	}
	return nil, nil
}

func (m *MockMemosService) PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error) {
	if m.PatchFilesFunc != nil {
		return m.PatchFilesFunc(ctx, memo, filePaths, replaces)
	}
	return nil, nil
}

func (m *MockMemosService) AddMemoRelations(ctx context.Context, memo string, relatedMemos []string, replaces bool) (*memos.AddMemoRelationsOutput, error) {
	if m.AddMemoRelationsFunc != nil {
		return m.AddMemoRelationsFunc(ctx, memo, relatedMemos, replaces)
	}
	return nil, nil
}

func (m *MockMemosUtilityService) CreateClip(ctx context.Context, input CreateClipInput) (*CreateClipOutput, error) {
	if m.CreateClipFunc != nil {
		return m.CreateClipFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockMemosUtilityService) CreateClips(ctx context.Context, input CreateClipsInput) (*CreateClipsOutput, error) {
	if m.CreateClipsFunc != nil {
		return m.CreateClipsFunc(ctx, input)
	}
	return nil, nil
}

func (m *MockMemosUtilityService) CreateCommonMemos(ctx context.Context, input CreateCommonMemosInput) (*CreateCommonMemosOutput, error) {
	if m.CreateCommonMemosFunc != nil {
		return m.CreateCommonMemosFunc(ctx, input)
	}
	return nil, nil
}

var _ MemosService = (*memos.Service)(nil)
var _ MemosService = (*MockMemosService)(nil)
var _ MemosUtilityService = (*Service)(nil)
var _ MemosUtilityService = (*MockMemosUtilityService)(nil)
