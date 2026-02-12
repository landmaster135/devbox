package usecases

import "context"

// MemoService は Memos ユースケースの公開契約。
type MemoService interface {
	CreateMemo(ctx context.Context, memoID string, content string, visibility string, state string, pinned *bool, displayTime string) (*Memo, error)
	GetMemo(ctx context.Context, memo string) (*Memo, error)
	ListMemos(ctx context.Context, pageSize int, pageToken string, state string, orderBy string) (*ListMemosOutput, error)
	UpdateMemo(ctx context.Context, memo string, content string, visibility string, state string, pinned *bool, updateMask []string) (*Memo, error)
}

// MockMemoService はテスト用モック。
type MockMemoService struct {
	CreateMemoFunc func(ctx context.Context, memoID string, content string, visibility string, state string, pinned *bool, displayTime string) (*Memo, error)
	GetMemoFunc    func(ctx context.Context, memo string) (*Memo, error)
	ListMemosFunc  func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string) (*ListMemosOutput, error)
	UpdateMemoFunc func(ctx context.Context, memo string, content string, visibility string, state string, pinned *bool, updateMask []string) (*Memo, error)
}

func (m *MockMemoService) CreateMemo(ctx context.Context, memoID string, content string, visibility string, state string, pinned *bool, displayTime string) (*Memo, error) {
	if m.CreateMemoFunc != nil {
		return m.CreateMemoFunc(ctx, memoID, content, visibility, state, pinned, displayTime)
	}
	return nil, nil
}

func (m *MockMemoService) GetMemo(ctx context.Context, memo string) (*Memo, error) {
	if m.GetMemoFunc != nil {
		return m.GetMemoFunc(ctx, memo)
	}
	return nil, nil
}

func (m *MockMemoService) ListMemos(ctx context.Context, pageSize int, pageToken string, state string, orderBy string) (*ListMemosOutput, error) {
	if m.ListMemosFunc != nil {
		return m.ListMemosFunc(ctx, pageSize, pageToken, state, orderBy)
	}
	return nil, nil
}

func (m *MockMemoService) UpdateMemo(ctx context.Context, memo string, content string, visibility string, state string, pinned *bool, updateMask []string) (*Memo, error) {
	if m.UpdateMemoFunc != nil {
		return m.UpdateMemoFunc(ctx, memo, content, visibility, state, pinned, updateMask)
	}
	return nil, nil
}

var _ MemoService = (*Service)(nil)
var _ MemoService = (*MockMemoService)(nil)
