package updatememo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	"github.com/landmaster135/devbox/internal/memos/usecases/common"
)

// Service は update_memo operation を扱う。
type Service struct {
	client     *common.JSONClient
	fileSystem infrastructures.FileSystem
}

func New(client *common.JSONClient, fileSystem infrastructures.FileSystem) *Service {
	return &Service{client: client, fileSystem: fileSystem}
}

func (s *Service) Execute(
	ctx context.Context,
	memo string,
	content string,
	contentFile string,
	visibility string,
	state string,
	pinned *bool,
	updateMask []string,
) (*common.Memo, error) {
	memoID := common.NormalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	resolvedContent, err := common.ResolveContent(content, contentFile, s.fileSystem)
	if err != nil {
		return nil, err
	}

	finalMask := buildUpdateMask(resolvedContent, visibility, state, pinned, updateMask)
	if len(finalMask) == 0 {
		return nil, fmt.Errorf("updateMask が空です")
	}

	query := url.Values{}
	query.Set("updateMask", strings.Join(finalMask, ","))

	payload := common.MemoMutationRequest{
		Content:    resolvedContent,
		Visibility: visibility,
		State:      state,
		Pinned:     pinned,
	}

	var result common.Memo
	requestPath := path.Join("/memos", url.PathEscape(memoID))
	if err := s.client.DoJSON(ctx, http.MethodPatch, requestPath, query, payload, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func buildUpdateMask(content string, visibility string, state string, pinned *bool, updateMask []string) []string {
	if len(updateMask) > 0 {
		return cleanMaskFields(updateMask)
	}

	mask := make([]string, 0, 4)
	if content != "" {
		mask = append(mask, "content")
	}
	if visibility != "" {
		mask = append(mask, "visibility")
	}
	if state != "" {
		mask = append(mask, "state")
	}
	if pinned != nil {
		mask = append(mask, "pinned")
	}
	return mask
}

func cleanMaskFields(raw []string) []string {
	seen := make(map[string]struct{})
	fields := make([]string, 0, len(raw))
	for _, item := range raw {
		for _, token := range strings.Split(item, ",") {
			field := strings.TrimSpace(token)
			if field == "" {
				continue
			}
			if _, ok := seen[field]; ok {
				continue
			}
			seen[field] = struct{}{}
			fields = append(fields, field)
		}
	}
	return fields
}
