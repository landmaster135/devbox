package listmemos

import (
	"context"
	"net/http"
	"testing"

	testutil "github.com/landmaster135/devbox/internal/memos/infrastructures/testutil"
	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
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
			return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/1","attachments":[{"name":"attachments/1"}]},{"name":"memos/2"}],"nextPageToken":"next-token","totalSize":2}`), nil
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
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Memos) != 2 {
		t.Fatalf("len(memos) = %d, want 2", len(result.Memos))
	}
	if len(result.Memos[0].Attachments) != 1 {
		t.Fatalf("len(memos[0].attachments) = %d, want 1", len(result.Memos[0].Attachments))
	}
	if result.Memos[0].Attachments[0].Name != "attachments/1" {
		t.Fatalf("memos[0].attachments[0].name = %q, want attachments/1", result.Memos[0].Attachments[0].Name)
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
		nil,
		nil,
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if called {
		t.Fatal("HTTP client was called, want not called")
	}
}

func TestServiceOperationListMemos_WithAnyContents_DedupByMemoID(t *testing.T) {
	callIndex := 0
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			callIndex++
			query := r.URL.Query()

			if got := query.Get("pageSize"); got != "20" {
				t.Fatalf("pageSize = %s, want 20", got)
			}
			if got := query.Get("filter"); got == "" {
				t.Fatal("filter is empty, want content.contains filter")
			}

			switch callIndex {
			case 1:
				if got := query.Get("filter"); got != `content.contains("meeting")` {
					t.Fatalf("filter(1st) = %s, want content filter", got)
				}
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/1"},{"name":"memos/2"}],"nextPageToken":"first-token","totalSize":2}`), nil
			case 2:
				if got := query.Get("filter"); got != `content.contains("study")` {
					t.Fatalf("filter(2nd) = %s, want content filter", got)
				}
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/2"},{"name":"memos/3"}],"nextPageToken":"second-token","totalSize":2}`), nil
			default:
				t.Fatalf("unexpected call index: %d", callIndex)
				return nil, nil
			}
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
		20,
		"",
		"",
		"",
		"",
		[]string{"meeting", "study"},
		nil,
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if callIndex != 2 {
		t.Fatalf("callIndex = %d, want 2", callIndex)
	}
	if len(result.Memos) != 3 {
		t.Fatalf("len(memos) = %d, want 3", len(result.Memos))
	}
	if result.Memos[0].Name != "memos/1" || result.Memos[1].Name != "memos/2" || result.Memos[2].Name != "memos/3" {
		t.Fatalf("memos order = %#v, want memos/1,memos/2,memos/3", result.Memos)
	}
	if result.NextPageToken != "" {
		t.Fatalf("nextPageToken = %q, want empty", result.NextPageToken)
	}
	if result.TotalSize != 3 {
		t.Fatalf("totalSize = %d, want 3", result.TotalSize)
	}
}

func TestServiceOperationListMemos_WithAnyContentsAndFilter_CombineConditions(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			query := r.URL.Query()
			got := query.Get("filter")
			want := `(created_ts > 1672578000 && visibility == "PUBLIC") && content.contains("meeting")`
			if got != want {
				t.Fatalf("filter = %q, want %q", got, want)
			}
			return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/1"}],"nextPageToken":"next-token","totalSize":1}`), nil
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
		20,
		"",
		"",
		"",
		`created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`,
		[]string{"meeting"},
		nil,
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.NextPageToken != "next-token" {
		t.Fatalf("nextPageToken = %q, want next-token", result.NextPageToken)
	}
	if result.TotalSize != 1 {
		t.Fatalf("totalSize = %d, want 1", result.TotalSize)
	}
}

func TestServiceOperationListMemos_WithMultipleAnyContentsAndFilter_CombineEachCondition(t *testing.T) {
	callIndex := 0
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			callIndex++
			query := r.URL.Query()

			switch callIndex {
			case 1:
				want := `(created_ts > 1672578000 && visibility == "PUBLIC") && content.contains("meeting")`
				if got := query.Get("filter"); got != want {
					t.Fatalf("filter(1st) = %q, want %q", got, want)
				}
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/1"},{"name":"memos/2"}],"totalSize":2}`), nil
			case 2:
				want := `(created_ts > 1672578000 && visibility == "PUBLIC") && content.contains("study")`
				if got := query.Get("filter"); got != want {
					t.Fatalf("filter(2nd) = %q, want %q", got, want)
				}
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/2"},{"name":"memos/3"}],"totalSize":2}`), nil
			default:
				t.Fatalf("unexpected call index: %d", callIndex)
				return nil, nil
			}
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
		20,
		"",
		"",
		"",
		`created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`,
		[]string{"meeting", "study"},
		nil,
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if callIndex != 2 {
		t.Fatalf("callIndex = %d, want 2", callIndex)
	}
	if len(result.Memos) != 3 {
		t.Fatalf("len(memos) = %d, want 3", len(result.Memos))
	}
	if result.Memos[0].Name != "memos/1" || result.Memos[1].Name != "memos/2" || result.Memos[2].Name != "memos/3" {
		t.Fatalf("memos order = %#v, want memos/1,memos/2,memos/3", result.Memos)
	}
}

func TestServiceOperationListMemos_WithAllContents_OnlyOverlapsRemain(t *testing.T) {
	callIndex := 0
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			callIndex++
			query := r.URL.Query()

			switch callIndex {
			case 1:
				if got := query.Get("filter"); got != `content.contains("meeting")` {
					t.Fatalf("filter(1st) = %s, want content filter", got)
				}
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/1"},{"name":"memos/2"}],"totalSize":2}`), nil
			case 2:
				if got := query.Get("filter"); got != `content.contains("study")` {
					t.Fatalf("filter(2nd) = %s, want content filter", got)
				}
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/2"},{"name":"memos/3"}],"totalSize":2}`), nil
			default:
				t.Fatalf("unexpected call index: %d", callIndex)
				return nil, nil
			}
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
		20,
		"",
		"",
		"",
		"",
		nil,
		[]string{"meeting", "study"},
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if callIndex != 2 {
		t.Fatalf("callIndex = %d, want 2", callIndex)
	}
	if len(result.Memos) != 1 {
		t.Fatalf("len(memos) = %d, want 1", len(result.Memos))
	}
	if result.Memos[0].Name != "memos/2" {
		t.Fatalf("memo name = %q, want memos/2", result.Memos[0].Name)
	}
	if result.NextPageToken != "" {
		t.Fatalf("nextPageToken = %q, want empty", result.NextPageToken)
	}
	if result.TotalSize != 1 {
		t.Fatalf("totalSize = %d, want 1", result.TotalSize)
	}
}

func TestServiceOperationListMemos_WithAllContentsAndFilter_CombineEachCondition(t *testing.T) {
	callIndex := 0
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			callIndex++
			query := r.URL.Query()

			switch callIndex {
			case 1:
				want := `(created_ts > 1672578000 && visibility == "PUBLIC") && content.contains("meeting")`
				if got := query.Get("filter"); got != want {
					t.Fatalf("filter(1st) = %q, want %q", got, want)
				}
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/1"},{"name":"memos/2"}],"totalSize":2}`), nil
			case 2:
				want := `(created_ts > 1672578000 && visibility == "PUBLIC") && content.contains("study")`
				if got := query.Get("filter"); got != want {
					t.Fatalf("filter(2nd) = %q, want %q", got, want)
				}
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/2"}],"totalSize":1}`), nil
			default:
				t.Fatalf("unexpected call index: %d", callIndex)
				return nil, nil
			}
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
		20,
		"",
		"",
		"",
		`created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"`,
		nil,
		[]string{"meeting", "study"},
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Memos) != 1 || result.Memos[0].Name != "memos/2" {
		t.Fatalf("memos = %#v, want only memos/2", result.Memos)
	}
}

func TestServiceOperationListMemos_WithAllContents_NoOverlap(t *testing.T) {
	callIndex := 0
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			callIndex++
			switch callIndex {
			case 1:
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/1"}],"totalSize":1}`), nil
			case 2:
				return testutil.JSONResponse(http.StatusOK, `{"memos":[{"name":"memos/2"}],"totalSize":1}`), nil
			default:
				t.Fatalf("unexpected call index: %d", callIndex)
				return nil, nil
			}
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
		20,
		"",
		"",
		"",
		"",
		nil,
		[]string{"meeting", "study"},
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Memos) != 0 {
		t.Fatalf("len(memos) = %d, want 0", len(result.Memos))
	}
	if result.TotalSize != 0 {
		t.Fatalf("totalSize = %d, want 0", result.TotalSize)
	}
}

func TestServiceOperationListMemos_WithSingleAllContents_ReturnsMatchedMemos(t *testing.T) {
	callIndex := 0
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			callIndex++
			if callIndex != 1 {
				t.Fatalf("callIndex = %d, want 1", callIndex)
			}
			query := r.URL.Query()
			if got := query.Get("filter"); got != `content.contains("any1")` {
				t.Fatalf("filter = %q, want %q", got, `content.contains("any1")`)
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
		20,
		"",
		"",
		"",
		"",
		nil,
		[]string{"any1"},
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Memos) != 2 {
		t.Fatalf("len(memos) = %d, want 2", len(result.Memos))
	}
	if result.Memos[0].Name != "memos/1" || result.Memos[1].Name != "memos/2" {
		t.Fatalf("memos = %#v, want memos/1,memos/2", result.Memos)
	}
	if result.NextPageToken != "" {
		t.Fatalf("nextPageToken = %q, want empty", result.NextPageToken)
	}
	if result.TotalSize != 2 {
		t.Fatalf("totalSize = %d, want 2", result.TotalSize)
	}
}
