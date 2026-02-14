package listmemos

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は list_memos operation を扱う。
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
	state string,
	orderBy string,
	filter string,
) (*common.ListMemosOutput, error) {
	query := url.Values{}
	if pageSize > 0 {
		query.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		query.Set("pageToken", pageToken)
	}
	if state != "" {
		query.Set("state", state)
	}
	if orderBy != "" {
		query.Set("orderBy", orderBy)
	}
	if filter != "" {
		query.Set("filter", filter)
	}

	var result common.ListMemosOutput
	if err := s.client.DoJSON(ctx, http.MethodGet, "/memos", query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
