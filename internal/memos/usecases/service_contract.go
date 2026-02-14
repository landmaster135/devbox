package usecases

import "context"

// MemoService は Memos ユースケースの公開契約。
type MemoService interface {
	CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*Memo, error)
	GetMemo(ctx context.Context, memo string) (*Memo, error)
	DeleteMemo(ctx context.Context, memo string, force bool) (*DeleteMemoOutput, error)
	ListMemos(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string) (*ListMemosOutput, error)
	ListAttachments(ctx context.Context, pageSize int, pageToken string, orderBy string, filter string) (*ListAttachmentsOutput, error)
	UpdateMemo(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*Memo, error)
	PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*SetMemoAttachmentsOutput, error)
	CreateAttachment(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*Attachment, error)
	ListMemoAttachments(ctx context.Context, memo string, pageSize int, pageToken string) (*ListMemoAttachmentsOutput, error)
	SetMemoAttachments(ctx context.Context, memo string, attachments []Attachment) (*SetMemoAttachmentsOutput, error)
}

// MockMemoService はテスト用モック。
type MockMemoService struct {
	CreateMemoFunc      func(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*Memo, error)
	GetMemoFunc         func(ctx context.Context, memo string) (*Memo, error)
	DeleteMemoFunc      func(ctx context.Context, memo string, force bool) (*DeleteMemoOutput, error)
	ListMemosFunc       func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string) (*ListMemosOutput, error)
	ListAttachmentsFunc func(ctx context.Context, pageSize int, pageToken string, orderBy string, filter string) (*ListAttachmentsOutput, error)
	UpdateMemoFunc      func(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*Memo, error)
	PatchFilesFunc      func(ctx context.Context, memo string, filePaths []string, replaces bool) (*SetMemoAttachmentsOutput, error)

	CreateAttachmentFunc    func(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*Attachment, error)
	ListMemoAttachmentsFunc func(ctx context.Context, memo string, pageSize int, pageToken string) (*ListMemoAttachmentsOutput, error)
	SetMemoAttachmentsFunc  func(ctx context.Context, memo string, attachments []Attachment) (*SetMemoAttachmentsOutput, error)
}

func (m *MockMemoService) CreateMemo(ctx context.Context, memoID string, content string, contentFile string, visibility string, state string, pinned *bool, displayTime string) (*Memo, error) {
	if m.CreateMemoFunc != nil {
		return m.CreateMemoFunc(ctx, memoID, content, contentFile, visibility, state, pinned, displayTime)
	}
	return nil, nil
}

func (m *MockMemoService) GetMemo(ctx context.Context, memo string) (*Memo, error) {
	if m.GetMemoFunc != nil {
		return m.GetMemoFunc(ctx, memo)
	}
	return nil, nil
}

func (m *MockMemoService) DeleteMemo(ctx context.Context, memo string, force bool) (*DeleteMemoOutput, error) {
	if m.DeleteMemoFunc != nil {
		return m.DeleteMemoFunc(ctx, memo, force)
	}
	return nil, nil
}

func (m *MockMemoService) ListMemos(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string) (*ListMemosOutput, error) {
	if m.ListMemosFunc != nil {
		return m.ListMemosFunc(ctx, pageSize, pageToken, state, orderBy, filter)
	}
	return nil, nil
}

func (m *MockMemoService) ListAttachments(ctx context.Context, pageSize int, pageToken string, orderBy string, filter string) (*ListAttachmentsOutput, error) {
	if m.ListAttachmentsFunc != nil {
		return m.ListAttachmentsFunc(ctx, pageSize, pageToken, orderBy, filter)
	}
	return nil, nil
}

func (m *MockMemoService) UpdateMemo(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*Memo, error) {
	if m.UpdateMemoFunc != nil {
		return m.UpdateMemoFunc(ctx, memo, content, contentFile, visibility, state, pinned, updateMask, displayTime)
	}
	return nil, nil
}

func (m *MockMemoService) PatchFiles(ctx context.Context, memo string, filePaths []string, replaces bool) (*SetMemoAttachmentsOutput, error) {
	if m.PatchFilesFunc != nil {
		return m.PatchFilesFunc(ctx, memo, filePaths, replaces)
	}
	return nil, nil
}

func (m *MockMemoService) CreateAttachment(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*Attachment, error) {
	if m.CreateAttachmentFunc != nil {
		return m.CreateAttachmentFunc(ctx, filename, content, attachmentType, memo)
	}
	return nil, nil
}

func (m *MockMemoService) ListMemoAttachments(ctx context.Context, memo string, pageSize int, pageToken string) (*ListMemoAttachmentsOutput, error) {
	if m.ListMemoAttachmentsFunc != nil {
		return m.ListMemoAttachmentsFunc(ctx, memo, pageSize, pageToken)
	}
	return nil, nil
}

func (m *MockMemoService) SetMemoAttachments(ctx context.Context, memo string, attachments []Attachment) (*SetMemoAttachmentsOutput, error) {
	if m.SetMemoAttachmentsFunc != nil {
		return m.SetMemoAttachmentsFunc(ctx, memo, attachments)
	}
	return nil, nil
}

var _ MemoService = (*Service)(nil)
var _ MemoService = (*MockMemoService)(nil)
