package usecases

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// GetMemo は GetMemo API を呼び出す。
func (s *Service) GetMemo(ctx context.Context, memo string) (*Memo, error) {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	var result Memo
	requestPath := path.Join("/memos", url.PathEscape(memoID))
	if err := s.doJSON(ctx, http.MethodGet, requestPath, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
