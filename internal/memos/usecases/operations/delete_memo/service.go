package deletememo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は delete_memo operation を扱う。
type Service struct {
	client *common.JSONClient
}

func New(client *common.JSONClient) *Service {
	return &Service{client: client}
}

func (s *Service) Execute(ctx context.Context, memo string, force bool) (*common.DeleteMemoOutput, error) {
	memoID := common.NormalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	var query url.Values
	if force {
		query = url.Values{}
		query.Set("force", "true")
	}

	requestPath := path.Join("/memos", url.PathEscape(memoID))
	if err := s.client.DoJSON(ctx, http.MethodDelete, requestPath, query, nil, nil); err != nil {
		return nil, err
	}

	return &common.DeleteMemoOutput{}, nil
}
