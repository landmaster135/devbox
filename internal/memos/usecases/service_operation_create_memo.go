package usecases

import (
	"context"
	"net/http"
	"net/url"
)

// CreateMemo は CreateMemo API を呼び出す。
func (s *Service) CreateMemo(
	ctx context.Context,
	memoID string,
	content string,
	contentFile string,
	visibility string,
	state string,
	pinned *bool,
	displayTime string,
) (*Memo, error) {
	resolvedContent, err := s.resolveContent(content, contentFile)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if memoID != "" {
		query.Set("memoId", memoID)
	}

	payload := memoMutationRequest{
		Content:     resolvedContent,
		Visibility:  visibility,
		State:       state,
		Pinned:      pinned,
		DisplayTime: displayTime,
	}

	var result Memo
	if err := s.doJSON(ctx, http.MethodPost, "/memos", query, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
