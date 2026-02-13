package creatememo

import (
	"context"
	"net/http"
	"net/url"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	"github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は create_memo operation を扱う。
type Service struct {
	client     *common.JSONClient
	fileSystem infrastructures.FileSystem
}

func New(client *common.JSONClient, fileSystem infrastructures.FileSystem) *Service {
	return &Service{client: client, fileSystem: fileSystem}
}

func (s *Service) Execute(
	ctx context.Context,
	memoID string,
	content string,
	contentFile string,
	visibility string,
	state string,
	pinned *bool,
	displayTime string,
) (*common.Memo, error) {
	resolvedContent, err := common.ResolveContent(content, contentFile, s.fileSystem)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	if memoID != "" {
		query.Set("memoId", memoID)
	}

	payload := common.MemoMutationRequest{
		Content:     resolvedContent,
		Visibility:  visibility,
		State:       state,
		Pinned:      pinned,
		DisplayTime: displayTime,
	}

	var result common.Memo
	if err := s.client.DoJSON(ctx, http.MethodPost, "/memos", query, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
