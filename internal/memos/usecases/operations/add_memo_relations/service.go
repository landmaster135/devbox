package addmemorelations

import (
	"context"
	"fmt"
	"strings"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

const (
	defaultListPageSize = 100
	addRelationType     = common.MemoRelationTypeReference
)

// RelationLister はリレーション一覧取得契約。
type RelationLister interface {
	List(ctx context.Context, memo string, pageSize int, pageToken string) (*common.ListMemoRelationsOutput, error)
}

// RelationSetter はリレーション更新契約。
type RelationSetter interface {
	Set(ctx context.Context, memo string, relations []common.MemoRelation) error
}

// Service は add-memo-relations operation を扱う。
type Service struct {
	lister RelationLister
	setter RelationSetter
}

func New(lister RelationLister, setter RelationSetter) *Service {
	return &Service{
		lister: lister,
		setter: setter,
	}
}

func (s *Service) Execute(
	ctx context.Context,
	memo string,
	relatedMemos []string,
	replaces bool,
) (*common.AddMemoRelationsOutput, error) {
	baseMemoName := common.BuildMemoResourceName(memo)
	if baseMemoName == "" {
		return nil, fmt.Errorf("memo が空です")
	}

	newRelatedMemoNames := normalizeRelatedMemoNames(baseMemoName, relatedMemos)
	if len(newRelatedMemoNames) == 0 {
		return nil, fmt.Errorf("related-memos が空です")
	}

	existingRelations, err := s.listAllMemoRelations(ctx, baseMemoName)
	if err != nil {
		return nil, err
	}

	newRelations := buildRelations(baseMemoName, newRelatedMemoNames)
	relationsToSet := newRelations
	if !replaces {
		relationsToSet = mergeRelations(existingRelations, newRelations)
	}

	if err := s.setter.Set(ctx, baseMemoName, relationsToSet); err != nil {
		return nil, err
	}

	finalRelations, err := s.listAllMemoRelations(ctx, baseMemoName)
	if err != nil {
		return nil, err
	}

	return &common.AddMemoRelationsOutput{
		Memo:               baseMemoName,
		DiscardedRelations: diffRelations(existingRelations, finalRelations),
		AddedRelations:     diffRelations(finalRelations, existingRelations),
		FinalRelations:     finalRelations,
	}, nil
}

func (s *Service) listAllMemoRelations(ctx context.Context, memo string) ([]common.MemoRelation, error) {
	all := make([]common.MemoRelation, 0)
	pageToken := ""

	for {
		result, err := s.lister.List(ctx, memo, defaultListPageSize, pageToken)
		if err != nil {
			return nil, err
		}
		if result == nil {
			return all, nil
		}

		all = append(all, result.Relations...)
		if strings.TrimSpace(result.NextPageToken) == "" {
			return all, nil
		}
		pageToken = result.NextPageToken
	}
}

func normalizeRelatedMemoNames(baseMemoName string, relatedMemos []string) []string {
	out := make([]string, 0, len(relatedMemos))
	seen := make(map[string]struct{}, len(relatedMemos))
	for _, relatedMemo := range relatedMemos {
		relatedMemoName := common.BuildMemoResourceName(relatedMemo)
		if relatedMemoName == "" || relatedMemoName == baseMemoName {
			continue
		}
		if _, ok := seen[relatedMemoName]; ok {
			continue
		}
		seen[relatedMemoName] = struct{}{}
		out = append(out, relatedMemoName)
	}
	return out
}

func buildRelations(baseMemoName string, relatedMemoNames []string) []common.MemoRelation {
	relations := make([]common.MemoRelation, 0, len(relatedMemoNames))
	for _, relatedMemoName := range relatedMemoNames {
		relations = append(relations, common.MemoRelation{
			Memo:        common.MemoRelationMemo{Name: baseMemoName},
			RelatedMemo: common.MemoRelationMemo{Name: relatedMemoName},
			Type:        addRelationType,
		})
	}
	return relations
}

func mergeRelations(existing []common.MemoRelation, added []common.MemoRelation) []common.MemoRelation {
	merged := make([]common.MemoRelation, 0, len(existing)+len(added))
	seen := make(map[string]struct{}, len(existing)+len(added))

	for _, relation := range existing {
		key := relationKey(relation)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		merged = append(merged, relation)
	}

	for _, relation := range added {
		key := relationKey(relation)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		merged = append(merged, relation)
	}
	return merged
}

func diffRelations(left []common.MemoRelation, right []common.MemoRelation) []common.MemoRelation {
	rightMap := make(map[string]struct{}, len(right))
	for _, relation := range right {
		rightMap[relationKey(relation)] = struct{}{}
	}

	diff := make([]common.MemoRelation, 0, len(left))
	for _, relation := range left {
		if _, ok := rightMap[relationKey(relation)]; ok {
			continue
		}
		diff = append(diff, relation)
	}
	return diff
}

func relationKey(relation common.MemoRelation) string {
	memoName := common.BuildMemoResourceName(relation.Memo.Name)
	relatedMemoName := common.BuildMemoResourceName(relation.RelatedMemo.Name)
	return memoName + "|" + relatedMemoName
}
