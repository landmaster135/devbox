package usecases

import (
	"context"
	"fmt"
	"strings"
)

// PatchFiles はファイルを添付として作成し、メモに紐づける。
func (s *Service) PatchFiles(
	ctx context.Context,
	memo string,
	filePaths []string,
	replaces bool,
) (*SetMemoAttachmentsOutput, error) {
	created := make([]Attachment, 0, len(filePaths))
	for _, filePath := range filePaths {
		attachmentFile, err := s.fileSystem.ReadAttachmentFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("files の読み込みに失敗しました (%s): %w", filePath, err)
		}
		attachment, err := s.CreateAttachment(
			ctx,
			attachmentFile.Filename,
			attachmentFile.Content,
			attachmentFile.ContentType,
			memo,
		)
		if err != nil {
			return nil, err
		}
		if attachment == nil {
			return nil, fmt.Errorf("CreateAttachment の結果が空です")
		}
		created = append(created, *attachment)
	}

	finalAttachments := created
	if !replaces {
		existing, err := s.listAllMemoAttachments(ctx, memo)
		if err != nil {
			return nil, err
		}
		finalAttachments = mergeAttachmentsByName(existing, created)
	}

	return s.SetMemoAttachments(ctx, memo, finalAttachments)
}

func (s *Service) listAllMemoAttachments(ctx context.Context, memo string) ([]Attachment, error) {
	const pageSize = 100

	var all []Attachment
	pageToken := ""

	for {
		result, err := s.ListMemoAttachments(ctx, memo, pageSize, pageToken)
		if err != nil {
			return nil, err
		}
		if result == nil {
			return all, nil
		}
		all = append(all, result.Attachments...)
		if strings.TrimSpace(result.NextPageToken) == "" {
			return all, nil
		}
		pageToken = result.NextPageToken
	}
}

func mergeAttachmentsByName(existing []Attachment, created []Attachment) []Attachment {
	createdByName := make(map[string]struct{}, len(created))
	for _, attachment := range created {
		name := strings.TrimSpace(attachment.Name)
		if name == "" {
			continue
		}
		createdByName[name] = struct{}{}
	}

	merged := make([]Attachment, 0, len(existing)+len(created))
	for _, attachment := range existing {
		name := strings.TrimSpace(attachment.Name)
		if name != "" {
			if _, exists := createdByName[name]; exists {
				continue
			}
		}
		merged = append(merged, attachment)
	}
	merged = append(merged, created...)
	return merged
}
