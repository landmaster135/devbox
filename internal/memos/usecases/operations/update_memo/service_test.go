package updatememo

import (
	"context"
	"net/http"
	"strings"
	"testing"

	testutil "github.com/landmaster135/devbox/internal/memos/infrastructures/testutil"
	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

func TestServiceOperationUpdateMemo_AutoUpdateMask_Normal(t *testing.T) {
	pinned := false

	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPatch {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPatch)
			}
			if r.URL.Path != "/api/v1/memos/memo-xyz" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-xyz", r.URL.Path)
			}
			if got := r.URL.Query().Get("updateMask"); got != "content,visibility,pinned" {
				t.Fatalf("updateMask = %s, want content,visibility,pinned", got)
			}
			body := testutil.ReadBodyAsMap(t, r.Body)
			if got := body["content"]; got != "updated text" {
				t.Fatalf("content = %v, want updated text", got)
			}
			if got := body["visibility"]; got != "PUBLIC" {
				t.Fatalf("visibility = %v, want PUBLIC", got)
			}
			if got := body["pinned"]; got != false {
				t.Fatalf("pinned = %v, want false", got)
			}
			if _, ok := body["displayTime"]; ok {
				t.Fatalf("displayTime = %v, want absent", body["displayTime"])
			}
			return testutil.JSONResponse(http.StatusOK, `{"name":"memos/memo-xyz","content":"updated text","visibility":"PUBLIC","pinned":false}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient, &testutil.MockFileSystem{})

	result, err := service.Execute(
		context.Background(),
		"memo-xyz",
		"updated text",
		"",
		"PUBLIC",
		"",
		&pinned,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Visibility != "PUBLIC" {
		t.Fatalf("visibility = %s, want PUBLIC", result.Visibility)
	}
}

func TestServiceOperationUpdateMemo_CustomUpdateMask_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if got := r.URL.Query().Get("updateMask"); got != "visibility,content,state" {
				t.Fatalf("updateMask = %s, want visibility,content,state", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"name":"memos/memo-123","content":"new content","state":"ARCHIVED"}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient, &testutil.MockFileSystem{})

	_, err := service.Execute(
		context.Background(),
		"memo-123",
		"new content",
		"",
		"",
		"ARCHIVED",
		nil,
		[]string{" visibility ", "content,content", "state"},
		"",
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func TestServiceOperationUpdateMemo_ContentFile_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if got := r.URL.Query().Get("updateMask"); got != "content,state,display_time" {
				t.Fatalf("updateMask = %s, want content,state,display_time", got)
			}
			body := testutil.ReadBodyAsMap(t, r.Body)
			if got := body["content"]; got != "updated from file" {
				t.Fatalf("content = %v, want updated from file", got)
			}
			if got := body["state"]; got != "ARCHIVED" {
				t.Fatalf("state = %v, want ARCHIVED", got)
			}
			if got := body["displayTime"]; got != "2026-02-14T01:23:45Z" {
				t.Fatalf("displayTime = %v, want 2026-02-14T01:23:45Z", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"name":"memos/memo-1","content":"updated from file","state":"ARCHIVED"}`), nil
		},
	}
	fileSystem := &testutil.MockFileSystem{
		ReadFileFunc: func(filePath string) ([]byte, error) {
			if filePath != "./updated.md" {
				t.Fatalf("filePath = %s, want ./updated.md", filePath)
			}
			return []byte("updated from file"), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient, fileSystem)

	result, err := service.Execute(
		context.Background(),
		"memo-1",
		"",
		"./updated.md",
		"",
		"ARCHIVED",
		nil,
		nil,
		"2026-02-14T01:23:45Z",
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.State != "ARCHIVED" {
		t.Fatalf("state = %s, want ARCHIVED", result.State)
	}
}

func TestServiceOperationUpdateMemo_EmptyMemo_Error(t *testing.T) {
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: &testutil.MockHTTPClient{},
	}), &testutil.MockFileSystem{})

	_, err := service.Execute(context.Background(), "", "value", "", "", "", nil, nil, "")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}

func TestServiceOperationUpdateMemo_EmptyUpdateMask_Error(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			t.Fatal("HTTP call should not be executed")
			return nil, nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient, &testutil.MockFileSystem{})

	_, err := service.Execute(
		context.Background(),
		"memo-1",
		"value",
		"",
		"",
		"",
		nil,
		[]string{" ", ","},
		"",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "updateMask が空") {
		t.Fatalf("error = %v, want updateMask が空", err)
	}
}

func TestBuildUpdateMask_Normal(t *testing.T) {
	pinned := true
	mask := buildUpdateMask("value", "PRIVATE", "ARCHIVED", &pinned, nil, "")
	got := strings.Join(mask, ",")
	if got != "content,visibility,state,pinned" {
		t.Fatalf("mask = %s, want content,visibility,state,pinned", got)
	}
}

func TestBuildUpdateMask_WithDisplayTime_Normal(t *testing.T) {
	mask := buildUpdateMask("value", "", "", nil, []string{"content"}, "2026-02-14T01:23:45Z")
	got := strings.Join(mask, ",")
	if got != "content,display_time" {
		t.Fatalf("mask = %s, want content,display_time", got)
	}
}
