package usecases

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

// CreateAttachment は CreateAttachment API を呼び出す。
func (s *Service) CreateAttachment(
	ctx context.Context,
	filename string,
	content []byte,
	attachmentType string,
	memo string,
) (*Attachment, error) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, fmt.Errorf("filename が空です")
	}

	attachmentType = strings.TrimSpace(attachmentType)
	if attachmentType == "" {
		return nil, fmt.Errorf("type が空です")
	}

	payload := createAttachmentRequest{
		Filename: filename,
		Content:  content,
		Type:     attachmentType,
	}
	if memoName := buildMemoResourceName(memo); memoName != "" {
		payload.Memo = memoName
	}

	var result Attachment
	if err := s.doJSON(ctx, http.MethodPost, "/attachments", nil, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMemoAttachments は ListMemoAttachments API を呼び出す。
func (s *Service) ListMemoAttachments(
	ctx context.Context,
	memo string,
	pageSize int,
	pageToken string,
) (*ListMemoAttachmentsOutput, error) {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}

	requestPath := path.Join("/memos", url.PathEscape(memoID), "attachments")
	var result ListMemoAttachmentsOutput
	if err := s.doJSON(ctx, http.MethodGet, requestPath, query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetMemoAttachments は SetMemoAttachments API を呼び出す。
func (s *Service) SetMemoAttachments(
	ctx context.Context,
	memo string,
	attachments []Attachment,
) (*SetMemoAttachmentsOutput, error) {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	requestPath := path.Join("/memos", url.PathEscape(memoID), "attachments")
	payload := setMemoAttachmentsRequest{
		Name:        buildMemoResourceName(memoID),
		Attachments: attachments,
	}

	var result SetMemoAttachmentsOutput
	if err := s.doJSON(ctx, http.MethodPatch, requestPath, nil, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
