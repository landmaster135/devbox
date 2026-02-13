package usecases

import "context"

// CreateMemo は CreateMemo API を呼び出す。
func (s *Service) CreateMemo(
	ctx context.Context,
	memoID string,
	content string,
	contentFile string,
	visibility string,
	state string,
	pinned *bool,
	displayTime string,
) (*Memo, error) {
	return s.createMemoOp.Execute(ctx, memoID, content, contentFile, visibility, state, pinned, displayTime)
}

// GetMemo は GetMemo API を呼び出す。
func (s *Service) GetMemo(ctx context.Context, memo string) (*Memo, error) {
	return s.getMemoOp.Execute(ctx, memo)
}

// ListMemos は ListMemos API を呼び出す。
func (s *Service) ListMemos(
	ctx context.Context,
	pageSize int,
	pageToken string,
	state string,
	orderBy string,
) (*ListMemosOutput, error) {
	return s.listMemosOp.Execute(ctx, pageSize, pageToken, state, orderBy)
}

// UpdateMemo は UpdateMemo API を呼び出す。
func (s *Service) UpdateMemo(
	ctx context.Context,
	memo string,
	content string,
	contentFile string,
	visibility string,
	state string,
	pinned *bool,
	updateMask []string,
	displayTime string,
) (*Memo, error) {
	return s.updateMemoOp.Execute(ctx, memo, content, contentFile, visibility, state, pinned, updateMask, displayTime)
}

// PatchFiles はファイルを添付として作成し、メモに紐づける。
func (s *Service) PatchFiles(
	ctx context.Context,
	memo string,
	filePaths []string,
	replaces bool,
) (*SetMemoAttachmentsOutput, error) {
	return s.patchFilesOp.Execute(ctx, memo, filePaths, replaces)
}

// CreateAttachment は CreateAttachment API を呼び出す。
func (s *Service) CreateAttachment(
	ctx context.Context,
	filename string,
	content []byte,
	attachmentType string,
	memo string,
) (*Attachment, error) {
	return s.createAttachmentOp.Create(ctx, filename, content, attachmentType, memo)
}

// ListMemoAttachments は ListMemoAttachments API を呼び出す。
func (s *Service) ListMemoAttachments(
	ctx context.Context,
	memo string,
	pageSize int,
	pageToken string,
) (*ListMemoAttachmentsOutput, error) {
	return s.listMemoAttachmentsOp.List(ctx, memo, pageSize, pageToken)
}

// SetMemoAttachments は SetMemoAttachments API を呼び出す。
func (s *Service) SetMemoAttachments(
	ctx context.Context,
	memo string,
	attachments []Attachment,
) (*SetMemoAttachmentsOutput, error) {
	return s.setMemoAttachmentsOp.Set(ctx, memo, attachments)
}
