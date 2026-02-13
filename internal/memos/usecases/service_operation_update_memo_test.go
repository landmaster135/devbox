package usecases

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

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
		"",
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
		"",
		"ARCHIVED",
		nil,
		[]string{" visibility ", "content,content", "state"},
	)
	if err != nil {
		t.Fatalf("UpdateMemo() error = %v", err)
	}
}

func TestService_UpdateMemo_ContentFile_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			if got := r.URL.Query().Get("updateMask"); got != "content,state" {
				t.Fatalf("updateMask = %s, want content,state", got)
			}
			body := readBodyAsMap(t, r.Body)
			if got := body["content"]; got != "updated from file" {
				t.Fatalf("content = %v, want updated from file", got)
			}
			if got := body["state"]; got != "ARCHIVED" {
				t.Fatalf("state = %v, want ARCHIVED", got)
			}
			return jsonResponse(http.StatusOK, `{"name":"memos/memo-1","content":"updated from file","state":"ARCHIVED"}`), nil
		},
	}
	fileSystem := &mockFileSystem{
		readFileFunc: func(filePath string) ([]byte, error) {
			if filePath != "./updated.md" {
				t.Fatalf("filePath = %s, want ./updated.md", filePath)
			}
			return []byte("updated from file"), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
		FileSystem: fileSystem,
	})

	result, err := service.UpdateMemo(
		context.Background(),
		"memo-1",
		"",
		"./updated.md",
		"",
		"ARCHIVED",
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("UpdateMemo() error = %v", err)
	}
	if result.State != "ARCHIVED" {
		t.Fatalf("state = %s, want ARCHIVED", result.State)
	}
}

func TestService_UpdateMemo_EmptyMemo_Error(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: &mockHTTPClient{},
	})

	_, err := service.UpdateMemo(context.Background(), "", "value", "", "", "", nil, nil)
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

func TestBuildUpdateMask_Normal(t *testing.T) {
	pinned := true
	mask := buildUpdateMask("value", "PRIVATE", "ARCHIVED", &pinned, nil)
	got := strings.Join(mask, ",")
	if got != "content,visibility,state,pinned" {
		t.Fatalf("mask = %s, want content,visibility,state,pinned", got)
	}
}
