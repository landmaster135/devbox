package usecases

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// UpdateMemo は UpdateMemo API を呼び出す。
func (s *Service) UpdateMemo(
	ctx context.Context,
	memo string,
	content string,
	contentFile string,
	visibility string,
	state string,
	pinned *bool,
	updateMask []string,
) (*Memo, error) {
	memoID := normalizeMemoIdentifier(memo)
	if memoID == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	resolvedContent, err := s.resolveContent(content, contentFile)
	if err != nil {
		return nil, err
	}

	finalMask := buildUpdateMask(resolvedContent, visibility, state, pinned, updateMask)
	if len(finalMask) == 0 {
		return nil, fmt.Errorf("updateMask が空です")
	}

	query := url.Values{}
	query.Set("updateMask", strings.Join(finalMask, ","))

	payload := memoMutationRequest{
		Content:    resolvedContent,
		Visibility: visibility,
		State:      state,
		Pinned:     pinned,
	}

	var result Memo
	requestPath := path.Join("/memos", url.PathEscape(memoID))
	if err := s.doJSON(ctx, http.MethodPatch, requestPath, query, payload, &result); err != nil {
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
