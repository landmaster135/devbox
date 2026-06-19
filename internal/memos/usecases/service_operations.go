package usecases

import (
	"context"
	"strings"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

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

// DeleteMemo は DeleteMemo API を呼び出す。
func (s *Service) DeleteMemo(ctx context.Context, memo string, force bool) (*DeleteMemoOutput, error) {
	return s.deleteMemoOp.Execute(ctx, memo, force)
}

// ListMemos は ListMemos API を呼び出す。
func (s *Service) ListMemos(
	ctx context.Context,
	pageSize int,
	pageToken string,
	state string,
	orderBy string,
	filter string,
	anyContents []string,
	allContents []string,
	allTags []string,
) (*ListMemosOutput, error) {
	return s.listMemosOp.Execute(ctx, pageSize, pageToken, state, orderBy, filter, anyContents, allContents, allTags)
}

// ListAttachments は ListAttachments API を呼び出す。
func (s *Service) ListAttachments(
	ctx context.Context,
	pageSize int,
	pageToken string,
	orderBy string,
	filter string,
) (*ListAttachmentsOutput, error) {
	return s.listAttachmentsOp.Execute(ctx, pageSize, pageToken, orderBy, filter)
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

// UpdateTag は srcTag を含むメモのタグを destTag に置換して更新する。
func (s *Service) UpdateTag(
	ctx context.Context,
	srcTag string,
	destTag string,
) (*UpdateTagOutput, error) {
	return s.updateTagOp.Execute(ctx, srcTag, destTag)
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

// ListMemoRelations は ListMemoRelations API を呼び出し、全ページを取得する。
func (s *Service) ListMemoRelations(
	ctx context.Context,
	memo string,
) (*ListMemoRelationsOutput, error) {
	const pageSize = 100

	all := make([]common.MemoRelation, 0)
	pageToken := ""

	for {
		result, err := s.listMemoRelationsOp.List(ctx, memo, pageSize, pageToken)
		if err != nil {
			return nil, err
		}
		if result == nil {
			return &ListMemoRelationsOutput{Relations: all}, nil
		}

		all = append(all, result.Relations...)
		if strings.TrimSpace(result.NextPageToken) == "" {
			return &ListMemoRelationsOutput{Relations: all}, nil
		}
		pageToken = result.NextPageToken
	}
}

// AddMemoRelations は指定メモへ related memos を追加/置換する。
func (s *Service) AddMemoRelations(
	ctx context.Context,
	memo string,
	relatedMemos []string,
	replaces bool,
) (*AddMemoRelationsOutput, error) {
	return s.addMemoRelationsOp.Execute(ctx, memo, relatedMemos, replaces)
}

// ListUserTags は GetUserStats API から tagCount を取得し、JSON ファイルへ保存する。
func (s *Service) ListUserTags(
	ctx context.Context,
	userID string,
	outputDir string,
) (*ListUserTagsOutput, error) {
	return s.listUserTagsOp.Execute(ctx, userID, outputDir)
}
