package updatetag

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

type mockMemoLister struct {
	executeFunc func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*common.ListMemosOutput, error)
}

func (m *mockMemoLister) Execute(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*common.ListMemosOutput, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, pageSize, pageToken, state, orderBy, filter, anyContents, allContents, allTags)
	}
	return nil, nil
}

type mockMemoUpdater struct {
	executeFunc func(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*common.Memo, error)
}

func (m *mockMemoUpdater) Execute(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*common.Memo, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, memo, content, contentFile, visibility, state, pinned, updateMask, displayTime)
	}
	return nil, nil
}

func TestServiceOperationUpdateTag_Normal(t *testing.T) {
	listCalled := 0
	updatedContents := map[string]string{}
	var updatedContentsMu sync.Mutex

	service := New(
		&mockMemoLister{
			executeFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*common.ListMemosOutput, error) {
				listCalled++
				if len(anyContents) != 0 || len(allContents) != 0 || len(allTags) != 0 {
					t.Fatalf("anyContents/allContents/allTags = %v/%v/%v, want empty", anyContents, allContents, allTags)
				}
				if pageSize != listMemosPageSize {
					t.Fatalf("pageSize = %d, want %d", pageSize, listMemosPageSize)
				}
				if state != "" {
					t.Fatalf("state = %q, want empty", state)
				}
				if orderBy != "" {
					t.Fatalf("orderBy = %q, want empty", orderBy)
				}
				if filter != `"work" in tags` {
					t.Fatalf("filter = %q, want %q", filter, `"work" in tags`)
				}
				switch listCalled {
				case 1:
					if pageToken != "" {
						t.Fatalf("pageToken = %q, want empty", pageToken)
					}
					return &common.ListMemosOutput{
						Memos: []common.Memo{
							{Name: "memos/1", Content: "alpha #work bravo"},
							{Name: "memos/2", Content: "#workshop #work #work-item"},
						},
						NextPageToken: "next-token",
					}, nil
				case 2:
					if pageToken != "next-token" {
						t.Fatalf("pageToken = %q, want next-token", pageToken)
					}
					return &common.ListMemosOutput{
						Memos: []common.Memo{
							{Name: "memos/3", Content: "no replace"},
						},
					}, nil
				default:
					t.Fatalf("unexpected list call count: %d", listCalled)
					return nil, nil
				}
			},
		},
		&mockMemoUpdater{
			executeFunc: func(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*common.Memo, error) {
				if contentFile != "" {
					t.Fatalf("contentFile = %q, want empty", contentFile)
				}
				if visibility != "" || state != "" {
					t.Fatalf("visibility/state = %q/%q, want empty", visibility, state)
				}
				if pinned != nil {
					t.Fatalf("pinned = %v, want nil", pinned)
				}
				if displayTime != "" {
					t.Fatalf("displayTime = %q, want empty", displayTime)
				}
				if len(updateMask) != 1 || updateMask[0] != "content" {
					t.Fatalf("updateMask = %v, want [content]", updateMask)
				}
				updatedContentsMu.Lock()
				updatedContents[memo] = content
				updatedContentsMu.Unlock()
				return &common.Memo{Name: memo}, nil
			},
		},
	)

	result, err := service.Execute(context.Background(), "work", "project")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if listCalled != 2 {
		t.Fatalf("listCalled = %d, want 2", listCalled)
	}
	if result.SourceTag != "work" {
		t.Fatalf("sourceTag = %q, want work", result.SourceTag)
	}
	if result.DestinationTag != "project" {
		t.Fatalf("destinationTag = %q, want project", result.DestinationTag)
	}
	if result.MatchedCount != 3 {
		t.Fatalf("matchedCount = %d, want 3", result.MatchedCount)
	}
	if result.UpdatedCount != 2 {
		t.Fatalf("updatedCount = %d, want 2", result.UpdatedCount)
	}
	if len(result.UpdatedMemoNames) != 2 {
		t.Fatalf("updatedMemoNames = %v, want 2 items", result.UpdatedMemoNames)
	}

	if got := updatedContents["memos/1"]; got != "alpha #project bravo" {
		t.Fatalf("updated content for memos/1 = %q, want %q", got, "alpha #project bravo")
	}
	if got := updatedContents["memos/2"]; got != "#workshop #project #work-item" {
		t.Fatalf("updated content for memos/2 = %q, want %q", got, "#workshop #project #work-item")
	}
}

func TestServiceOperationUpdateTag_SourceTagWithVariationSelector_Normal(t *testing.T) {
	crossWithVS := "001-todo-status/2-wip-❌️"
	crossWithoutVS := "001-todo-status/2-wip-❌"

	filterChecked := false
	updated := false

	service := New(
		&mockMemoLister{
			executeFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*common.ListMemosOutput, error) {
				filterChecked = true
				if len(anyContents) != 0 || len(allContents) != 0 || len(allTags) != 0 {
					t.Fatalf("anyContents/allContents/allTags = %v/%v/%v, want empty", anyContents, allContents, allTags)
				}
				if !strings.Contains(filter, `"`+crossWithVS+`" in tags`) {
					t.Fatalf("filter = %q, want to include VS form", filter)
				}
				if !strings.Contains(filter, `"`+crossWithoutVS+`" in tags`) {
					t.Fatalf("filter = %q, want to include non-VS form", filter)
				}
				return &common.ListMemosOutput{
					Memos: []common.Memo{
						{Name: "memos/1", Content: "before #" + crossWithoutVS + " after"},
					},
				}, nil
			},
		},
		&mockMemoUpdater{
			executeFunc: func(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*common.Memo, error) {
				updated = true
				want := "before #001-todo-status/2-wip after"
				if content != want {
					t.Fatalf("content = %q, want %q", content, want)
				}
				return &common.Memo{Name: memo}, nil
			},
		},
	)

	result, err := service.Execute(context.Background(), crossWithVS, "001-todo-status/2-wip")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !filterChecked {
		t.Fatal("filter was not checked")
	}
	if !updated {
		t.Fatal("UpdateMemo was not called")
	}
	if result.UpdatedCount != 1 {
		t.Fatalf("updatedCount = %d, want 1", result.UpdatedCount)
	}
}

func TestServiceOperationUpdateTag_ParallelUpdate_Normal(t *testing.T) {
	started := make(chan struct{}, 3)
	release := make(chan struct{})
	var current int32
	var maxConcurrent int32

	service := New(
		&mockMemoLister{
			executeFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*common.ListMemosOutput, error) {
				if len(anyContents) != 0 || len(allContents) != 0 || len(allTags) != 0 {
					t.Fatalf("anyContents/allContents/allTags = %v/%v/%v, want empty", anyContents, allContents, allTags)
				}
				return &common.ListMemosOutput{
					Memos: []common.Memo{
						{Name: "memos/1", Content: "#work"},
						{Name: "memos/2", Content: "#work"},
						{Name: "memos/3", Content: "#work"},
					},
				}, nil
			},
		},
		&mockMemoUpdater{
			executeFunc: func(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*common.Memo, error) {
				working := atomic.AddInt32(&current, 1)
				for {
					prev := atomic.LoadInt32(&maxConcurrent)
					if working <= prev {
						break
					}
					if atomic.CompareAndSwapInt32(&maxConcurrent, prev, working) {
						break
					}
				}

				started <- struct{}{}
				<-release
				atomic.AddInt32(&current, -1)
				return &common.Memo{Name: memo}, nil
			},
		},
	)

	resultCh := make(chan *common.UpdateTagOutput, 1)
	errCh := make(chan error, 1)
	go func() {
		result, err := service.Execute(context.Background(), "work", "project")
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- result
	}()

	for i := 0; i < 3; i++ {
		select {
		case <-started:
		case err := <-errCh:
			t.Fatalf("Execute() error = %v", err)
		case <-time.After(2 * time.Second):
			t.Fatal("timeout waiting for updater starts")
		}
	}

	close(release)

	var result *common.UpdateTagOutput
	select {
	case err := <-errCh:
		t.Fatalf("Execute() error = %v", err)
	case result = <-resultCh:
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for Execute()")
	}

	if result.UpdatedCount != 3 {
		t.Fatalf("updatedCount = %d, want 3", result.UpdatedCount)
	}
	if atomic.LoadInt32(&maxConcurrent) < 2 {
		t.Fatalf("maxConcurrent = %d, want >= 2", atomic.LoadInt32(&maxConcurrent))
	}
}

func TestServiceOperationUpdateTag_MissingSrcTag_Error(t *testing.T) {
	service := New(&mockMemoLister{}, &mockMemoUpdater{})

	_, err := service.Execute(context.Background(), "", "dest")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "src-tag") {
		t.Fatalf("error = %v, want src-tag", err)
	}
}

func TestServiceOperationUpdateTag_MissingDestTag_Error(t *testing.T) {
	service := New(&mockMemoLister{}, &mockMemoUpdater{})

	_, err := service.Execute(context.Background(), "src", "")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "dest-tag") {
		t.Fatalf("error = %v, want dest-tag", err)
	}
}

func TestServiceOperationUpdateTag_UpdateError_Error(t *testing.T) {
	service := New(
		&mockMemoLister{
			executeFunc: func(ctx context.Context, pageSize int, pageToken string, state string, orderBy string, filter string, anyContents []string, allContents []string, allTags []string) (*common.ListMemosOutput, error) {
				if len(anyContents) != 0 || len(allContents) != 0 || len(allTags) != 0 {
					t.Fatalf("anyContents/allContents/allTags = %v/%v/%v, want empty", anyContents, allContents, allTags)
				}
				return &common.ListMemosOutput{
					Memos: []common.Memo{
						{Name: "memos/1", Content: "#src"},
					},
				}, nil
			},
		},
		&mockMemoUpdater{
			executeFunc: func(ctx context.Context, memo string, content string, contentFile string, visibility string, state string, pinned *bool, updateMask []string, displayTime string) (*common.Memo, error) {
				return nil, fmt.Errorf("update failed")
			},
		},
	)

	_, err := service.Execute(context.Background(), "src", "dest")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "update failed") {
		t.Fatalf("error = %v, want update failed", err)
	}
}

func TestTagReplacerReplace_Normal(t *testing.T) {
	replacer := newTagReplacer(normalizeTagForComparison("tag"), "new-tag")

	input := "A #tag B #tag-keep C #tagged D #tag\n#tag!"
	got, changed := replacer.Replace(input)
	if !changed {
		t.Fatal("changed = false, want true")
	}
	want := "A #new-tag B #tag-keep C #tagged D #new-tag\n#new-tag!"
	if got != want {
		t.Fatalf("replace result = %q, want %q", got, want)
	}
}

func TestTagReplacerReplace_TagWithSlash_Normal(t *testing.T) {
	replacer := newTagReplacer(normalizeTagForComparison("001-status/1-todo"), "001-todo-status/1-planned")

	input := "before #001-status/1-todo after\n#001-status/1-todo/x"
	got, changed := replacer.Replace(input)
	if !changed {
		t.Fatal("changed = false, want true")
	}
	want := "before #001-todo-status/1-planned after\n#001-status/1-todo/x"
	if got != want {
		t.Fatalf("replace result = %q, want %q", got, want)
	}
}

func TestTagReplacerReplace_TagWithEmoji_Normal(t *testing.T) {
	pill := "\U0001F48A"
	sourceTag := "001-todo-status/2-wip-" + pill
	targetTag := "001-todo-status/2-wip"
	replacer := newTagReplacer(normalizeTagForComparison(sourceTag), targetTag)

	input := "before #" + sourceTag + " after\n#" + sourceTag + "-extra"
	got, changed := replacer.Replace(input)
	if !changed {
		t.Fatal("changed = false, want true")
	}
	want := "before #" + targetTag + " after\n#" + sourceTag + "-extra"
	if got != want {
		t.Fatalf("replace result = %q, want %q", got, want)
	}
}

func TestBuildTagFilter_WithVariationSelector_Normal(t *testing.T) {
	filter := buildTagFilter("001-todo-status/2-wip-❌️")
	if !strings.Contains(filter, `"001-todo-status/2-wip-❌️" in tags`) {
		t.Fatalf("filter = %q, want VS form", filter)
	}
	if !strings.Contains(filter, `"001-todo-status/2-wip-❌" in tags`) {
		t.Fatalf("filter = %q, want non-VS form", filter)
	}
}
