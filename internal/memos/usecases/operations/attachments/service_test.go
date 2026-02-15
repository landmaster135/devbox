package attachments

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/memos/usecases/common"
	"github.com/landmaster135/devbox/internal/memos/usecases/testutil"
)

func TestServiceOperationCreateAttachment_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
			}
			if r.URL.Path != "/api/v1/attachments" {
				t.Fatalf("path = %s, want /api/v1/attachments", r.URL.Path)
			}
			body := testutil.ReadBodyAsMap(t, r.Body)
			if got := body["filename"]; got != "hello.txt" {
				t.Fatalf("filename = %v, want hello.txt", got)
			}
			if got := body["type"]; got != "text/plain" {
				t.Fatalf("type = %v, want text/plain", got)
			}
			if got := body["memo"]; got != "memos/memo-1" {
				t.Fatalf("memo = %v, want memos/memo-1", got)
			}
			if got := body["content"]; got != "aGVsbG8=" {
				t.Fatalf("content = %v, want aGVsbG8=", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"name":"attachments/1","filename":"hello.txt","type":"text/plain","memo":"memos/memo-1"}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	result, err := service.Create(context.Background(), "hello.txt", []byte("hello"), "text/plain", "memo-1")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.Name != "attachments/1" {
		t.Fatalf("name = %s, want attachments/1", result.Name)
	}
}

func TestServiceOperationListMemoAttachments_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/memos/memo-123/attachments" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-123/attachments", r.URL.Path)
			}
			query := r.URL.Query()
			if got := query.Get("pageSize"); got != "100" {
				t.Fatalf("pageSize = %s, want 100", got)
			}
			if got := query.Get("pageToken"); got != "next" {
				t.Fatalf("pageToken = %s, want next", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"attachments":[{"name":"attachments/1"},{"name":"attachments/2"}],"nextPageToken":"next-2"}`), nil
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
	if len(result.Attachments) != 2 {
		t.Fatalf("len(attachments) = %d, want 2", len(result.Attachments))
	}
	if result.NextPageToken != "next-2" {
		t.Fatalf("nextPageToken = %s, want next-2", result.NextPageToken)
	}
}

func TestServiceOperationSetMemoAttachments_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPatch {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPatch)
			}
			if r.URL.Path != "/api/v1/memos/memo-123/attachments" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-123/attachments", r.URL.Path)
			}
			body := testutil.ReadBodyAsMap(t, r.Body)
			if got := body["name"]; got != "memos/memo-123" {
				t.Fatalf("name = %v, want memos/memo-123", got)
			}
			attachments, ok := body["attachments"].([]any)
			if !ok {
				t.Fatalf("attachments type = %T, want []any", body["attachments"])
			}
			if len(attachments) != 2 {
				t.Fatalf("len(attachments) = %d, want 2", len(attachments))
			}
			return testutil.JSONResponse(http.StatusOK, `{"name":"memos/memo-123","attachments":[{"name":"attachments/1"},{"name":"attachments/2"}]}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	result, err := service.Set(context.Background(), "memo-123", []common.Attachment{
		{Name: "attachments/1"},
		{Name: "attachments/2"},
	})
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}
	if result.Name != "memos/memo-123" {
		t.Fatalf("name = %s, want memos/memo-123", result.Name)
	}
}

func TestServiceOperationListMemoAttachments_EmptyMemo_Error(t *testing.T) {
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

func TestServiceOperationSetMemoAttachments_EmptyMemo_Error(t *testing.T) {
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: &testutil.MockHTTPClient{},
	}))

	_, err := service.Set(context.Background(), "", nil)
	if err == nil {
		t.Fatal("Set() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}
