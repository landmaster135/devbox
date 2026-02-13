package patchfiles

import (
	"context"
	"fmt"
	"strings"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	"github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// AttachmentCreator は添付作成契約。
type AttachmentCreator interface {
	Create(ctx context.Context, filename string, content []byte, attachmentType string, memo string) (*common.Attachment, error)
}

// AttachmentLister は添付一覧取得契約。
type AttachmentLister interface {
	List(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoAttachmentsOutput, error)
}

// AttachmentSetter は添付更新契約。
type AttachmentSetter interface {
	Set(ctx context.Context, memo string, attachments []common.Attachment) (*common.SetMemoAttachmentsOutput, error)
}

// Service は patch_files operation を扱う。
type Service struct {
	fileSystem        infrastructures.FileSystem
	attachmentCreator AttachmentCreator
	attachmentLister  AttachmentLister
	attachmentSetter  AttachmentSetter
}

func New(
	fileSystem infrastructures.FileSystem,
	attachmentCreator AttachmentCreator,
	attachmentLister AttachmentLister,
	attachmentSetter AttachmentSetter,
) *Service {
	return &Service{
		fileSystem:        fileSystem,
		attachmentCreator: attachmentCreator,
		attachmentLister:  attachmentLister,
		attachmentSetter:  attachmentSetter,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	memo string,
	filePaths []string,
	replaces bool,
) (*common.SetMemoAttachmentsOutput, error) {
	created := make([]common.Attachment, 0, len(filePaths))
	for _, filePath := range filePaths {
		attachmentFile, err := s.fileSystem.ReadAttachmentFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("files の読み込みに失敗しました (%s): %w", filePath, err)
		}
		attachment, err := s.attachmentCreator.Create(
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

	return s.attachmentSetter.Set(ctx, memo, finalAttachments)
}

func (s *Service) listAllMemoAttachments(ctx context.Context, memo string) ([]common.Attachment, error) {
	const pageSize = 100

	var all []common.Attachment
	pageToken := ""

	for {
		result, err := s.attachmentLister.List(ctx, memo, pageSize, pageToken)
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

func mergeAttachmentsByName(existing []common.Attachment, created []common.Attachment) []common.Attachment {
	createdByName := make(map[string]struct{}, len(created))
	for _, attachment := range created {
		name := strings.TrimSpace(attachment.Name)
		if name == "" {
			continue
		}
		createdByName[name] = struct{}{}
	}

	merged := make([]common.Attachment, 0, len(existing)+len(created))
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
