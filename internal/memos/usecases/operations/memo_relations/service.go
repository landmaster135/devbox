package memorelations

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は memo relations API を扱う。
type Service struct {
	client *common.JSONClient
}

func New(client *common.JSONClient) *Service {
	return &Service{client: client}
}

func (s *Service) List(
	ctx context.Context,
	memo string,
	pageSize int,
	pageToken string,
) (*common.ListMemoRelationsOutput, error) {
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

	requestPath := path.Join("/memos", url.PathEscape(memoID), "relations")
	var result common.ListMemoRelationsOutput
	if err := s.client.DoJSON(ctx, http.MethodGet, requestPath, query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) Set(
	ctx context.Context,
	memo string,
	relations []common.MemoRelation,
) error {
	memoID := common.NormalizeMemoIdentifier(memo)
	if memoID == "" {
		return fmt.Errorf("memo が空です")
	}

	requestPath := path.Join("/memos", url.PathEscape(memoID), "relations")
	payload := common.SetMemoRelationsRequest{
		Name:      common.BuildMemoResourceName(memoID),
		Relations: relations,
	}

	return s.client.DoJSON(ctx, http.MethodPatch, requestPath, nil, payload, nil)
}
