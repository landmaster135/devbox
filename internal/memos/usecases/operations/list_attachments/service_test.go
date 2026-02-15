package listattachments

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/landmaster135/devbox/internal/memos/usecases/common"
	"github.com/landmaster135/devbox/internal/memos/usecases/testutil"
)

func TestServiceOperationListAttachments_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/attachments" {
				t.Fatalf("path = %s, want /api/v1/attachments", r.URL.Path)
			}
			query := r.URL.Query()
			if got := query.Get("pageSize"); got != "50" {
				t.Fatalf("pageSize = %s, want 50", got)
			}
			if got := query.Get("pageToken"); got != "next-token" {
				t.Fatalf("pageToken = %s, want next-token", got)
			}
			if got := query.Get("orderBy"); got != "create_time desc" {
				t.Fatalf("orderBy = %s, want create_time desc", got)
			}
			if got := query.Get("filter"); got != `memo == "memos/memo-1"` {
				t.Fatalf("filter = %s, want memo filter", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"attachments":[{"name":"attachments/1"},{"name":"attachments/2"}],"nextPageToken":"next-2","totalSize":2}`), nil
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
		50,
		"next-token",
		"create_time desc",
		`memo == "memos/memo-1"`,
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Attachments) != 2 {
		t.Fatalf("len(attachments) = %d, want 2", len(result.Attachments))
	}
	if result.NextPageToken != "next-2" {
		t.Fatalf("nextPageToken = %s, want next-2", result.NextPageToken)
	}
	if result.TotalSize != 2 {
		t.Fatalf("totalSize = %d, want 2", result.TotalSize)
	}
}

func TestServiceOperationListAttachments_HTTPError_Error(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, errors.New("request failed")
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	_, err := service.Execute(context.Background(), 20, "", "", "")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
}

func TestServiceOperationListAttachments_EmptyOptionalParams_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/attachments" {
				t.Fatalf("path = %s, want /api/v1/attachments", r.URL.Path)
			}
			if got := r.URL.RawQuery; got != "" {
				t.Fatalf("rawQuery = %q, want empty", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"attachments":[]}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	result, err := service.Execute(context.Background(), 0, "", "", "")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Attachments) != 0 {
		t.Fatalf("len(attachments) = %d, want 0", len(result.Attachments))
	}
}
