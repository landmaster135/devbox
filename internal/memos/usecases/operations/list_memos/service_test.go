package listmemos

import (
	"context"
	"net/http"
	"testing"

	"github.com/landmaster135/devbox/internal/memos/usecases/common"
	"github.com/landmaster135/devbox/internal/memos/usecases/testutil"
)

func TestServiceOperationListMemos_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/memos" {
				t.Fatalf("path = %s, want /api/v1/memos", r.URL.Path)
			}
			query := r.URL.Query()
			if got := query.Get("pageSize"); got != "15" {
				t.Fatalf("pageSize = %s, want 15", got)
			}
			if got := query.Get("pageToken"); got != "next-token" {
				t.Fatalf("pageToken = %s, want next-token", got)
			}
			if got := query.Get("state"); got != "NORMAL" {
				t.Fatalf("state = %s, want NORMAL", got)
			}
			if got := query.Get("orderBy"); got != "update_time desc" {
				t.Fatalf("orderBy = %s, want update_time desc", got)
			}
			if got := query.Get("filter"); got != `created_ts > 1672578000 && visibility == "PUBLIC"` {
				t.Fatalf("filter = %s, want filter condition", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/1"},{"name":"memos/2"}],"nextPageToken":"next-token","totalSize":2}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	result, err := service.Execute(
		context.Background(),
		15,
		"next-token",
		"NORMAL",
		"update_time desc",
		`created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`,
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Memos) != 2 {
		t.Fatalf("len(memos) = %d, want 2", len(result.Memos))
	}
}

func TestServiceOperationListMemos_InvalidTimestampFilter_Error(t *testing.T) {
	called := false
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			called = true
			return testutil.JSONResponse(http.StatusOK, `{"memos":[]}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	_, err := service.Execute(
		context.Background(),
		15,
		"",
		"NORMAL",
		"",
		`created_ts > "2023-01-01T13:00:00"`,
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if called {
		t.Fatal("HTTP client was called, want not called")
	}
}
