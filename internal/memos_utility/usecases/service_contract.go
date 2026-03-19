package usecases

import (
	"context"

	memos "github.com/landmaster135/devbox/internal/memos/usecases"
)

// MemosService は memos-utility が利用する Memos サービスの契約。
type MemosService interface {
	CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
}

// MemosUtilityService は memos-utility の公開契約。
type MemosUtilityService interface {
	CreateClip(ctx context.Context, input CreateClipInput) (*CreateClipOutput, error)
	CreateClips(ctx context.Context, input CreateClipsInput) (*CreateClipsOutput, error)
}

// MockMemosService はテスト用モック。
type MockMemosService struct {
	CreateMemoFunc func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	PatchFilesFunc func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
}

// MockMemosUtilityService は MemosUtilityService のテストダブル。
type MockMemosUtilityService struct {
	CreateClipFunc  func(ctx context.Context, input CreateClipInput) (*CreateClipOutput, error)
	CreateClipsFunc func(ctx context.Context, input CreateClipsInput) (*CreateClipsOutput, error)
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

var _ MemosService = (*memos.Service)(nil)
var _ MemosService = (*MockMemosService)(nil)
var _ MemosUtilityService = (*Service)(nil)
var _ MemosUtilityService = (*MockMemosUtilityService)(nil)
