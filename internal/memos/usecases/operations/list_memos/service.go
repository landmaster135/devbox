package listmemos

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
	commonfilter "github.com/landmaster135/devbox/internal/memos/usecases/common/filter"
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
	anyContents []string,
	allContents []string,
	allTags []string,
) (*common.ListMemosOutput, error) {
	normalizedAllTags := normalizeTags(allTags)
	if len(normalizedAllTags) > 0 {
		return s.executeAllTags(ctx, pageSize, pageToken, state, orderBy, filter, normalizedAllTags)
	}

	normalizedAllContents := normalizeContents(allContents)
	if len(normalizedAllContents) > 0 {
		return s.executeAllContents(ctx, pageSize, pageToken, state, orderBy, filter, normalizedAllContents)
	}

	normalizedAnyContents := normalizeContents(anyContents)
	if len(normalizedAnyContents) == 0 {
		return s.executeSingle(ctx, pageSize, pageToken, state, orderBy, filter)
	}

	mergedMemos := make([]common.Memo, 0)
	seenMemoKeys := make(map[string]struct{})
	nextPageToken := ""

	for i, term := range normalizedAnyContents {
		combinedFilter := buildContentContainsFilter(filter, term)
		result, err := s.executeSingle(ctx, pageSize, pageToken, state, orderBy, combinedFilter)
		if err != nil {
			return nil, err
		}
		if result == nil {
			continue
		}

		mergedMemos = appendDedupMemos(mergedMemos, result.Memos, seenMemoKeys)
		if len(normalizedAnyContents) == 1 && i == 0 {
			nextPageToken = strings.TrimSpace(result.NextPageToken)
		}
	}

	return &common.ListMemosOutput{
		Memos:         mergedMemos,
		NextPageToken: nextPageToken,
		TotalSize:     int64(len(mergedMemos)),
	}, nil
}

func (s *Service) executeAllTags(
	ctx context.Context,
	pageSize int,
	pageToken string,
	state string,
	orderBy string,
	filter string,
	allTags []string,
) (*common.ListMemosOutput, error) {
	requiredCount := len(allTags)
	if requiredCount == 0 {
		return &common.ListMemosOutput{
			Memos:         []common.Memo{},
			NextPageToken: "",
			TotalSize:     0,
		}, nil
	}

	countByMemoKey := make(map[string]int)
	firstMemoByKey := make(map[string]common.Memo)
	memoOrder := make([]string, 0)

	for _, tag := range allTags {
		combinedFilter := buildTagContainsFilter(filter, tag)
		result, err := s.executeSingle(ctx, pageSize, pageToken, state, orderBy, combinedFilter)
		if err != nil {
			return nil, err
		}
		if result == nil {
			continue
		}

		seenInTag := make(map[string]struct{})
		for _, memo := range result.Memos {
			memoKey := memoDedupKey(memo)
			if memoKey == "" {
				continue
			}
			if _, exists := seenInTag[memoKey]; exists {
				continue
			}
			seenInTag[memoKey] = struct{}{}

			if _, exists := countByMemoKey[memoKey]; !exists {
				firstMemoByKey[memoKey] = memo
				memoOrder = append(memoOrder, memoKey)
			}
			countByMemoKey[memoKey]++
		}
	}

	overlappedMemos := make([]common.Memo, 0)
	for _, memoKey := range memoOrder {
		if countByMemoKey[memoKey] < requiredCount {
			continue
		}
		overlappedMemos = append(overlappedMemos, firstMemoByKey[memoKey])
	}

	return &common.ListMemosOutput{
		Memos:         overlappedMemos,
		NextPageToken: "",
		TotalSize:     int64(len(overlappedMemos)),
	}, nil
}

func (s *Service) executeAllContents(
	ctx context.Context,
	pageSize int,
	pageToken string,
	state string,
	orderBy string,
	filter string,
	allContents []string,
) (*common.ListMemosOutput, error) {
	requiredCount := len(allContents)
	if requiredCount == 0 {
		return &common.ListMemosOutput{
			Memos:         []common.Memo{},
			NextPageToken: "",
			TotalSize:     0,
		}, nil
	}

	countByMemoKey := make(map[string]int)
	firstMemoByKey := make(map[string]common.Memo)
	memoOrder := make([]string, 0)

	for _, term := range allContents {
		combinedFilter := buildContentContainsFilter(filter, term)
		result, err := s.executeSingle(ctx, pageSize, pageToken, state, orderBy, combinedFilter)
		if err != nil {
			return nil, err
		}
		if result == nil {
			continue
		}

		seenInTerm := make(map[string]struct{})
		for _, memo := range result.Memos {
			memoKey := memoDedupKey(memo)
			if memoKey == "" {
				continue
			}
			if _, exists := seenInTerm[memoKey]; exists {
				continue
			}
			seenInTerm[memoKey] = struct{}{}

			if _, exists := countByMemoKey[memoKey]; !exists {
				firstMemoByKey[memoKey] = memo
				memoOrder = append(memoOrder, memoKey)
			}
			countByMemoKey[memoKey]++
		}
	}

	overlappedMemos := make([]common.Memo, 0)
	for _, memoKey := range memoOrder {
		if countByMemoKey[memoKey] < requiredCount {
			continue
		}
		overlappedMemos = append(overlappedMemos, firstMemoByKey[memoKey])
	}

	return &common.ListMemosOutput{
		Memos:         overlappedMemos,
		NextPageToken: "",
		TotalSize:     int64(len(overlappedMemos)),
	}, nil
}

func (s *Service) executeSingle(
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
		normalizedFilter, err := commonfilter.NormalizeFilter(filter)
		if err != nil {
			return nil, err
		}
		query.Set("filter", normalizedFilter)
	}

	var result common.ListMemosOutput
	if err := s.client.DoJSON(ctx, http.MethodGet, "/memos", query, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func normalizeContents(contents []string) []string {
	out := make([]string, 0, len(contents))
	for _, content := range contents {
		term := strings.TrimSpace(content)
		if term == "" {
			continue
		}
		out = append(out, term)
	}
	return out
}

func buildContentContainsFilter(baseFilter string, content string) string {
	escapedContent := strings.ReplaceAll(content, `\`, `\\`)
	escapedContent = strings.ReplaceAll(escapedContent, `"`, `\"`)
	containsCondition := fmt.Sprintf(`content.contains("%s")`, escapedContent)
	trimmedBaseFilter := strings.TrimSpace(baseFilter)
	if trimmedBaseFilter == "" {
		return containsCondition
	}
	return fmt.Sprintf("(%s) && %s", trimmedBaseFilter, containsCondition)
}

func normalizeTags(tags []string) []string {
	out := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

func buildTagContainsFilter(baseFilter string, tag string) string {
	escapedTag := strings.ReplaceAll(tag, `\`, `\\`)
	escapedTag = strings.ReplaceAll(escapedTag, `'`, `\'`)
	tagFilter := fmt.Sprintf("tag in ['%s']", escapedTag)

	trimmedBaseFilter := strings.TrimSpace(baseFilter)
	if trimmedBaseFilter == "" {
		return tagFilter
	}
	return fmt.Sprintf("(%s) && (%s)", trimmedBaseFilter, tagFilter)
}

func appendDedupMemos(target []common.Memo, incoming []common.Memo, seen map[string]struct{}) []common.Memo {
	for _, memo := range incoming {
		dedupKey := memoDedupKey(memo)
		if dedupKey == "" {
			target = append(target, memo)
			continue
		}
		if _, exists := seen[dedupKey]; exists {
			continue
		}
		seen[dedupKey] = struct{}{}
		target = append(target, memo)
	}
	return target
}

func memoDedupKey(memo common.Memo) string {
	memoName := common.BuildMemoResourceName(memo.Name)
	if memoName != "" {
		return memoName
	}

	memoUID := common.BuildMemoResourceName(memo.UID)
	if memoUID != "" {
		return memoUID
	}

	if memo.ID > 0 {
		return fmt.Sprintf("memo-id:%d", memo.ID)
	}

	return ""
}
