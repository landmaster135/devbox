package updatetag

import (
	"context"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/landmaster135/devbox/internal/memos/usecases/common"
)

const listMemosPageSize = 100

// MemoLister は対象メモ取得の契約。
type MemoLister interface {
	Execute(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string) (*common.ListMemosOutput, error)
}

// MemoUpdater はメモ更新の契約。
type MemoUpdater interface {
	Execute(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*common.Memo, error)
}

// Service は update-tag operation を扱う。
type Service struct {
	memoLister  MemoLister
	memoUpdater MemoUpdater
}

func New(memoLister MemoLister, memoUpdater MemoUpdater) *Service {
	return &Service{
		memoLister:  memoLister,
		memoUpdater: memoUpdater,
	}
}

func (s *Service) Execute(ctx context.Context, srcTag string, destTag string) (*common.UpdateTagOutput, error) {
	sourceTag := normalizeTagValue(srcTag)
	if sourceTag == "" {
		return nil, fmt.Errorf("src-tag が空です")
	}

	sourceTagForCompare := normalizeTagForComparison(sourceTag)
	if sourceTagForCompare == "" {
		return nil, fmt.Errorf("src-tag が不正です")
	}

	destinationTag := normalizeTagValue(destTag)
	if destinationTag == "" {
		return nil, fmt.Errorf("dest-tag が空です")
	}

	filter := buildTagFilter(sourceTag)
	replacer := newTagReplacer(sourceTagForCompare, destinationTag)
	output := &common.UpdateTagOutput{
		SourceTag:      sourceTag,
		DestinationTag: destinationTag,
	}

	pageToken := ""
	for {
		listResult, err := s.memoLister.Execute(ctx, listMemosPageSize, pageToken, "", "", filter)
		if err != nil {
			return nil, err
		}
		if listResult == nil {
			return output, nil
		}

		for _, memo := range listResult.Memos {
			output.MatchedCount++
			if strings.TrimSpace(memo.Name) == "" {
				return nil, fmt.Errorf("更新対象の memo name が空です")
			}

			updatedContent, changed := replacer.Replace(memo.Content)
			if !changed {
				continue
			}

			updatedMemo, err := s.memoUpdater.Execute(
				ctx,
				memo.Name,
				updatedContent,
				"",
				"",
				"",
				nil,
				[]string{"content"},
				"",
			)
			if err != nil {
				return nil, err
			}

			updatedName := memo.Name
			if updatedMemo != nil {
				name := strings.TrimSpace(updatedMemo.Name)
				if name != "" {
					updatedName = name
				}
			}

			output.UpdatedCount++
			output.UpdatedMemoNames = append(output.UpdatedMemoNames, updatedName)
		}

		if strings.TrimSpace(listResult.NextPageToken) == "" {
			return output, nil
		}
		pageToken = listResult.NextPageToken
	}
}

func normalizeTagValue(raw string) string {
	return strings.TrimPrefix(strings.TrimSpace(raw), "#")
}

func buildTagFilter(tag string) string {
	var conditions []string
	for _, candidate := range buildTagSearchCandidates(tag) {
		escaped := strings.ReplaceAll(candidate, `\`, `\\`)
		escaped = strings.ReplaceAll(escaped, `"`, `\"`)
		conditions = append(conditions, fmt.Sprintf(`"%s" in tags`, escaped))
	}
	if len(conditions) == 0 {
		return ""
	}
	if len(conditions) == 1 {
		return conditions[0]
	}
	return "(" + strings.Join(conditions, " || ") + ")"
}

func buildTagSearchCandidates(tag string) []string {
	raw := strings.TrimSpace(tag)
	base := normalizeTagForComparison(raw)
	withEmojiVariation := addEmojiVariationSelector(base)

	candidates := make([]string, 0, 3)
	seen := make(map[string]struct{}, 3)
	for _, candidate := range []string{raw, base, withEmojiVariation} {
		if candidate == "" {
			continue
		}
		if _, exists := seen[candidate]; exists {
			continue
		}
		seen[candidate] = struct{}{}
		candidates = append(candidates, candidate)
	}
	return candidates
}

func normalizeTagForComparison(tag string) string {
	trimmed := strings.TrimSpace(tag)
	if trimmed == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(trimmed))
	for _, r := range trimmed {
		if isVariationSelector(r) {
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

func addEmojiVariationSelector(tag string) string {
	if tag == "" {
		return ""
	}

	runes := []rune(tag)
	var builder strings.Builder
	builder.Grow(len(tag) * 2)
	for i := 0; i < len(runes); i++ {
		current := runes[i]
		builder.WriteRune(current)
		if !unicode.IsSymbol(current) {
			continue
		}

		next := rune(0)
		if i+1 < len(runes) {
			next = runes[i+1]
		}
		if isVariationSelector(next) {
			continue
		}
		builder.WriteRune('\ufe0f')
	}
	return builder.String()
}

func isVariationSelector(r rune) bool {
	return r == '\ufe0e' || r == '\ufe0f'
}

type tagReplacer struct {
	sourceTagForCompare string
	targetTag           string
}

func newTagReplacer(sourceTagForCompare string, targetTag string) *tagReplacer {
	return &tagReplacer{
		sourceTagForCompare: sourceTagForCompare,
		targetTag:           targetTag,
	}
}

func (r *tagReplacer) Replace(content string) (string, bool) {
	if content == "" || r.sourceTagForCompare == "" {
		return content, false
	}

	var builder strings.Builder
	builder.Grow(len(content))

	changed := false
	index := 0
	var previousRune rune
	hasPreviousRune := false

	for index < len(content) {
		currentRune, currentSize := utf8.DecodeRuneInString(content[index:])
		if currentRune != '#' || !isTagPrefixBoundary(previousRune, hasPreviousRune) {
			builder.WriteString(content[index : index+currentSize])
			previousRune = currentRune
			hasPreviousRune = true
			index += currentSize
			continue
		}

		tagStart := index + currentSize
		tagEnd := tagStart
		for tagEnd < len(content) {
			candidateRune, candidateSize := utf8.DecodeRuneInString(content[tagEnd:])
			if !isTagRune(candidateRune) {
				break
			}
			tagEnd += candidateSize
		}

		if tagEnd == tagStart {
			builder.WriteString(content[index : index+currentSize])
			previousRune = currentRune
			hasPreviousRune = true
			index += currentSize
			continue
		}

		tagValue := content[tagStart:tagEnd]
		if normalizeTagForComparison(tagValue) == r.sourceTagForCompare {
			builder.WriteRune('#')
			builder.WriteString(r.targetTag)
			changed = true

			if r.targetTag == "" {
				previousRune = '#'
			} else {
				previousRune, _ = utf8.DecodeLastRuneInString(r.targetTag)
			}
		} else {
			builder.WriteString(content[index:tagEnd])
			previousRune, _ = utf8.DecodeLastRuneInString(content[index:tagEnd])
		}
		hasPreviousRune = true
		index = tagEnd
	}

	if !changed {
		return content, false
	}
	return builder.String(), true
}

func isTagPrefixBoundary(r rune, hasRune bool) bool {
	if !hasRune {
		return true
	}
	return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '#'
}

func isTagRune(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSymbol(r) || unicode.IsMark(r) {
		return true
	}
	return r == '_' || r == '-' || r == '/' || r == '\u200d' || isVariationSelector(r)
}
