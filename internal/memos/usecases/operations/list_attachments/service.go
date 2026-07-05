package listattachments

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は list_attachments operation を扱う。
type Service struct {
	client *common.JSONClient
}

func New(client *common.JSONClient) *Service {
	return &Service{client: client}
}

func (s *Service) Execute(
	ctx context.Context,
	pageSize int,
	pageToken string,
	orderBy string,
	filter string,
) (*common.ListAttachmentsOutput, error) {
	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}
	if orderBy != "" {
		query.Set("orderBy", orderBy)
	}
	if filter != "" {
		query.Set("filter", filter)
	}

	var result common.ListAttachmentsOutput
	if err := s.client.DoJSON(ctx, http.MethodGet, "/attachments", query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
