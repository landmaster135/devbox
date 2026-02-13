package creatememo

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/memos/usecases/common"
	"github.com/landmaster135/devbox/internal/memos/usecases/testutil"
)

func TestServiceOperationCreateMemo_Normal(t *testing.T) {
	pinned := true

	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
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
			body := testutil.ReadBodyAsMap(t, r.Body)
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
			return testutil.JSONResponse(http.StatusOK, `{"name":"memos/memo-1","content":"hello memo"}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "test-token",
		HTTPClient: client,
	})
	service := New(jsonClient, &testutil.MockFileSystem{})

	result, err := service.Execute(
		context.Background(),
		"memo-1",
		"hello memo",
		"",
		"PRIVATE",
		"NORMAL",
		&pinned,
		"2026-02-12T00:00:00Z",
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Name != "memos/memo-1" {
		t.Fatalf("name = %s, want memos/memo-1", result.Name)
	}
}

func TestServiceOperationCreateMemo_ContentFile_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			body := testutil.ReadBodyAsMap(t, r.Body)
			if got := body["content"]; got != "memo from file" {
				t.Fatalf("content = %v, want memo from file", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"name":"memos/memo-file","content":"memo from file"}`), nil
		},
	}
	fileSystem := &testutil.MockFileSystem{
		ReadFileFunc: func(filePath string) ([]byte, error) {
			if filePath != "./memo.md" {
				t.Fatalf("filePath = %s, want ./memo.md", filePath)
			}
			return []byte("memo from file"), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "test-token",
		HTTPClient: client,
	})
	service := New(jsonClient, fileSystem)

	result, err := service.Execute(
		context.Background(),
		"memo-file",
		"",
		"./memo.md",
		"PRIVATE",
		"NORMAL",
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Name != "memos/memo-file" {
		t.Fatalf("name = %s, want memos/memo-file", result.Name)
	}
}

func TestServiceOperationCreateMemo_ContentFileReadError_Error(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			t.Fatal("HTTP call should not be executed")
			return nil, nil
		},
	}
	fileSystem := &testutil.MockFileSystem{
		ReadFileFunc: func(filePath string) ([]byte, error) {
			return nil, fmt.Errorf("read failed")
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "test-token",
		HTTPClient: client,
	})
	service := New(jsonClient, fileSystem)

	_, err := service.Execute(
		context.Background(),
		"memo-1",
		"",
		"./memo.md",
		"",
		"",
		nil,
		"",
	)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content-file の読み込みに失敗しました") {
		t.Fatalf("error = %v, want content-file の読み込みに失敗しました", err)
	}
}

func TestServiceOperationCreateMemo_DecodeError_Error(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return testutil.JSONResponse(http.StatusOK, "{invalid"), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient, &testutil.MockFileSystem{})

	_, err := service.Execute(context.Background(), "", "memo", "", "", "", nil, "")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "レスポンスデコード") {
		t.Fatalf("error = %v, want レスポンスデコード", err)
	}
}
