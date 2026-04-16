package memorelations

import (
	"context"
	"net/http"
	"strings"
	"testing"

	testutil "github.com/landmaster135/devbox/internal/memos/infrastructures/testutil"
	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

func TestServiceOperationListMemoRelations_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/memos/memo-123/relations" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-123/relations", r.URL.Path)
			}
			query := r.URL.Query()
			if got := query.Get("pageSize"); got != "100" {
				t.Fatalf("pageSize = %s, want 100", got)
			}
			if got := query.Get("pageToken"); got != "next" {
				t.Fatalf("pageToken = %s, want next", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"relations":[{"memo":{"name":"memos/memo-123"},"relatedMemo":{"name":"memos/memo-456"},"type":"REFERENCE"}],"nextPageToken":"next-2"}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	result, err := service.List(context.Background(), "memos/memo-123", 100, "next")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(result.Relations) != 1 {
		t.Fatalf("len(relations) = %d, want 1", len(result.Relations))
	}
	if result.Relations[0].RelatedMemo.Name != "memos/memo-456" {
		t.Fatalf("relatedMemo = %s, want memos/memo-456", result.Relations[0].RelatedMemo.Name)
	}
	if result.NextPageToken != "next-2" {
		t.Fatalf("nextPageToken = %s, want next-2", result.NextPageToken)
	}
}

func TestServiceOperationSetMemoRelations_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPatch {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPatch)
			}
			if r.URL.Path != "/api/v1/memos/memo-123/relations" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-123/relations", r.URL.Path)
			}
			body := testutil.ReadBodyAsMap(t, r.Body)
			if got := body["name"]; got != "memos/memo-123" {
				t.Fatalf("name = %v, want memos/memo-123", got)
			}
			relations, ok := body["relations"].([]any)
			if !ok {
				t.Fatalf("relations type = %T, want []any", body["relations"])
			}
			if len(relations) != 1 {
				t.Fatalf("len(relations) = %d, want 1", len(relations))
			}
			return testutil.JSONResponse(http.StatusOK, `{}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	err := service.Set(context.Background(), "memo-123", []common.MemoRelation{
		{
			Memo:        common.MemoRelationMemo{Name: "memos/memo-123"},
			RelatedMemo: common.MemoRelationMemo{Name: "memos/memo-456"},
			Type:        "REFERENCE",
		},
	})
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}
}

func TestServiceOperationListMemoRelations_EmptyMemo_Error(t *testing.T) {
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: &testutil.MockHTTPClient{},
	}))

	_, err := service.List(context.Background(), "", 100, "")
	if err == nil {
		t.Fatal("List() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}

func TestServiceOperationSetMemoRelations_EmptyMemo_Error(t *testing.T) {
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: &testutil.MockHTTPClient{},
	}))

	err := service.Set(context.Background(), "", nil)
	if err == nil {
		t.Fatal("Set() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}
