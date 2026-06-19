package attachments

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は attachments operation 群を扱う。
type Service struct {
	client *common.JSONClient
}

func New(client *common.JSONClient) *Service {
	return &Service{client: client}
}

func (s *Service) Create(
	ctx context.Context,
	filename string,
	content []byte,
	attachmentType string,
	memo string,
) (*common.Attachment, error) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return nil, fmt.Errorf("filename が空です")
	}

	attachmentType = strings.TrimSpace(attachmentType)
	if attachmentType == "" {
		return nil, fmt.Errorf("type が空です")
	}

	payload := common.CreateAttachmentRequest{
		Filename: filename,
		Content:  content,
		Type:     attachmentType,
	}
	if memoName := common.BuildMemoResourceName(memo); memoName != "" {
		payload.Memo = memoName
	}

	var result common.Attachment
	if err := s.client.DoJSON(ctx, http.MethodPost, "/attachments", nil, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) List(
	ctx context.Context,
	memo string,
	pageSize int,
	pageToken string,
) (*common.ListMemoAttachmentsOutput, error) {
	memoID := common.NormalizeMemoIdentifier(memo)
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
	var result common.ListMemoAttachmentsOutput
	if err := s.client.DoJSON(ctx, http.MethodGet, requestPath, query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) Set(
	ctx context.Context,
	memo string,
	attachments []common.Attachment,
) (*common.SetMemoAttachmentsOutput, error) {
	memoID := common.NormalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	requestPath := path.Join("/memos", url.PathEscape(memoID), "attachments")
	payload := common.SetMemoAttachmentsRequest{
		Name:        common.BuildMemoResourceName(memoID),
		Attachments: attachments,
	}

	var result common.SetMemoAttachmentsOutput
	if err := s.client.DoJSON(ctx, http.MethodPatch, requestPath, nil, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
