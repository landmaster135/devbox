package usecases

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// ListMemos は ListMemos API を呼び出す。
func (s *Service) ListMemos(
	ctx context.Context,
	pageSize int,
	pageToken string,
	state string,
	orderBy string,
) (*ListMemosOutput, error) {
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

	var result ListMemosOutput
	if err := s.doJSON(ctx, http.MethodGet, "/memos", query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
