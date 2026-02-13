package getmemo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は get_memo operation を扱う。
type Service struct {
	client *common.JSONClient
}

func New(client *common.JSONClient) *Service {
	return &Service{client: client}
}

func (s *Service) Execute(ctx context.Context, memo string) (*common.Memo, error) {
	memoID := common.NormalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	var result common.Memo
	requestPath := path.Join("/memos", url.PathEscape(memoID))
	if err := s.client.DoJSON(ctx, http.MethodGet, requestPath, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
