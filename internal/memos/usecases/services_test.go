package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

func TestService_CreateMemo_Normal(t *testing.T) {
	pinned := true

	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
			}
			if r.URL.Path != "/api/v1/memos" {
				t.Fatalf("path = %s, want /api/v1/memos", r.URL.Path)
			}
			if got := r.URL.Query().Get("memoId"); got != "memo-1" {
				t.Fatalf("memoId = %s, want memo-1", got)
			}
			if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
				t.Fatalf("authorization = %s, want Bearer test-token", got)
			}
			body := readBodyAsMap(t, r.Body)
			if got := body["content"]; got != "hello memo" {
				t.Fatalf("content = %v, want hello memo", got)
			}
			if got := body["visibility"]; got != "PRIVATE" {
				t.Fatalf("visibility = %v, want PRIVATE", got)
			}
			if got := body["state"]; got != "NORMAL" {
				t.Fatalf("state = %v, want NORMAL", got)
			}
			if got := body["pinned"]; got != true {
				t.Fatalf("pinned = %v, want true", got)
			}
			return jsonResponse(http.StatusOK, `{"name":"memos/memo-1","content":"hello memo"}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "test-token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	result, err := service.CreateMemo(
		context.Background(),
		"memo-1",
		"hello memo",
		"PRIVATE",
		"NORMAL",
		&pinned,
		"2026-02-12T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("CreateMemo() error = %v", err)
	}
	if result.Name != "memos/memo-1" {
		t.Fatalf("name = %s, want memos/memo-1", result.Name)
	}
}

func TestService_GetMemo_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/memos/memo-abc" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-abc", r.URL.Path)
			}
			return jsonResponse(http.StatusOK, `{"name":"memos/memo-abc","content":"value"}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com/api/v1",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	result, err := service.GetMemo(context.Background(), "memos/memo-abc")
	if err != nil {
		t.Fatalf("GetMemo() error = %v", err)
	}
	if result.Content != "value" {
		t.Fatalf("content = %s, want value", result.Content)
	}
}

func TestService_ListMemos_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
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
			return jsonResponse(http.StatusOK, `{"memos":[{"name":"memos/1"},{"name":"memos/2"}],"nextPageToken":"next-token","totalSize":2}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	result, err := service.ListMemos(
		context.Background(),
		15,
		"next-token",
		"NORMAL",
		"update_time desc",
	)
	if err != nil {
		t.Fatalf("ListMemos() error = %v", err)
	}
	if len(result.Memos) != 2 {
		t.Fatalf("len(memos) = %d, want 2", len(result.Memos))
	}
}

func TestService_UpdateMemo_AutoUpdateMask_Normal(t *testing.T) {
	pinned := false

	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPatch {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPatch)
			}
			if r.URL.Path != "/api/v1/memos/memo-xyz" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-xyz", r.URL.Path)
			}
			if got := r.URL.Query().Get("updateMask"); got != "content,visibility,pinned" {
				t.Fatalf("updateMask = %s, want content,visibility,pinned", got)
			}
			body := readBodyAsMap(t, r.Body)
			if got := body["content"]; got != "updated text" {
				t.Fatalf("content = %v, want updated text", got)
			}
			if got := body["visibility"]; got != "PUBLIC" {
				t.Fatalf("visibility = %v, want PUBLIC", got)
			}
			if got := body["pinned"]; got != false {
				t.Fatalf("pinned = %v, want false", got)
			}
			return jsonResponse(http.StatusOK, `{"name":"memos/memo-xyz","content":"updated text","visibility":"PUBLIC","pinned":false}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	result, err := service.UpdateMemo(
		context.Background(),
		"memo-xyz",
		"updated text",
		"PUBLIC",
		"",
		&pinned,
		nil,
	)
	if err != nil {
		t.Fatalf("UpdateMemo() error = %v", err)
	}
	if result.Visibility != "PUBLIC" {
		t.Fatalf("visibility = %s, want PUBLIC", result.Visibility)
	}
}

func TestService_UpdateMemo_CustomUpdateMask_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			if got := r.URL.Query().Get("updateMask"); got != "visibility,content,state" {
				t.Fatalf("updateMask = %s, want visibility,content,state", got)
			}
			return jsonResponse(http.StatusOK, `{"name":"memos/memo-123","content":"new content","state":"ARCHIVED"}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	_, err := service.UpdateMemo(
		context.Background(),
		"memo-123",
		"new content",
		"",
		"ARCHIVED",
		nil,
		[]string{" visibility ", "content,content", "state"},
	)
	if err != nil {
		t.Fatalf("UpdateMemo() error = %v", err)
	}
}

func TestService_GetMemo_APIError(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusUnauthorized, `{"message":"unauthorized"}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	_, err := service.GetMemo(context.Background(), "memo-1")
	if err == nil {
		t.Fatal("GetMemo() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "status=401") {
		t.Fatalf("error = %v, want status=401", err)
	}
	if !strings.Contains(err.Error(), "unauthorized") {
		t.Fatalf("error = %v, want unauthorized", err)
	}
}

func TestService_DoJSONNetworkError_Error(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	_, err := service.GetMemo(context.Background(), "memo-1")
	if err == nil {
		t.Fatal("GetMemo() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Fatalf("error = %v, want network error", err)
	}
}

func TestService_GetMemo_EmptyMemo_Error(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: &mockHTTPClient{},
	})

	_, err := service.GetMemo(context.Background(), "")
	if err == nil {
		t.Fatal("GetMemo() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}

func TestService_UpdateMemo_EmptyMemo_Error(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: &mockHTTPClient{},
	})

	_, err := service.UpdateMemo(context.Background(), "", "value", "", "", nil, nil)
	if err == nil {
		t.Fatal("UpdateMemo() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}

func TestService_UpdateMemo_EmptyUpdateMask_Error(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			t.Fatal("HTTP call should not be executed")
			return nil, nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	_, err := service.UpdateMemo(
		context.Background(),
		"memo-1",
		"value",
		"",
		"",
		nil,
		[]string{" ", ","},
	)
	if err == nil {
		t.Fatal("UpdateMemo() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "updateMask が空") {
		t.Fatalf("error = %v, want updateMask が空", err)
	}
}

func TestService_CreateMemo_DecodeError_Error(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, "{invalid"), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	_, err := service.CreateMemo(context.Background(), "", "memo", "", "", nil, "")
	if err == nil {
		t.Fatal("CreateMemo() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "レスポンスデコード") {
		t.Fatalf("error = %v, want レスポンスデコード", err)
	}
}

func TestService_DoJSONOutNil_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			return jsonResponse(http.StatusOK, `{"ok":true}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	err := service.doJSON(context.Background(), http.MethodGet, "/memos", nil, nil, nil)
	if err != nil {
		t.Fatalf("doJSON() error = %v", err)
	}
}

func TestService_NewRequestEncodeError_Error(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: &mockHTTPClient{},
	})

	_, err := service.newRequest(context.Background(), http.MethodPost, "/memos", nil, map[string]any{
		"invalid": make(chan int),
	})
	if err == nil {
		t.Fatal("newRequest() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "エンコード") {
		t.Fatalf("error = %v, want エンコード", err)
	}
}

func TestBuildUpdateMask_Normal(t *testing.T) {
	pinned := true
	mask := buildUpdateMask("value", "PRIVATE", "ARCHIVED", &pinned, nil)
	got := strings.Join(mask, ",")
	if got != "content,visibility,state,pinned" {
		t.Fatalf("mask = %s, want content,visibility,state,pinned", got)
	}
}

func TestNewService_DefaultClient_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:  "https://memos.example.com",
		APIToken: "token",
		Timeout:  0,
	})
	if service.client == nil {
		t.Fatal("client is nil")
	}
	if service.baseURL != "https://memos.example.com/api/v1" {
		t.Fatalf("baseURL = %s, want https://memos.example.com/api/v1", service.baseURL)
	}
}

func TestNormalizeBaseURL_Normal(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "https://example.com", want: "https://example.com/api/v1"},
		{name: "with trailing slash", in: "https://example.com/", want: "https://example.com/api/v1"},
		{name: "with api v1", in: "https://example.com/api/v1", want: "https://example.com/api/v1"},
		{name: "empty", in: "", want: "/api/v1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeBaseURL(tc.in)
			if got != tc.want {
				t.Fatalf("normalizeBaseURL(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestNormalizeMemoIdentifier_Normal(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "memo-1", want: "memo-1"},
		{name: "with memos prefix", in: "memos/memo-1", want: "memo-1"},
		{name: "with api prefix", in: "api/v1/memos/memo-1", want: "memo-1"},
		{name: "with full url", in: "https://example.com/api/v1/memos/memo-1", want: "memo-1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeMemoIdentifier(tc.in)
			if got != tc.want {
				t.Fatalf("normalizeMemoIdentifier(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func readBodyAsMap(t *testing.T, body io.ReadCloser) map[string]any {
	t.Helper()

	defer body.Close()

	var m map[string]any
	if err := json.NewDecoder(body).Decode(&m); err != nil {
		t.Fatalf("decode body error: %v", err)
	}
	return m
}

func jsonResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
