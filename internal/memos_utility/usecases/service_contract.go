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
}

// MockMemosService はテスト用モック。
type MockMemosService struct {
	CreateMemoFunc func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*memos.Memo, error)
	PatchFilesFunc func(ctx context.Context, memo string, filePaths []string, replaces bool) (*memos.SetMemoAttachmentsOutput, error)
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

var _ MemosService = (*memos.Service)(nil)
var _ MemosService = (*MockMemosService)(nil)
