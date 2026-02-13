package usecases

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

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
		"",
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

func TestService_CreateMemo_ContentFile_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			body := readBodyAsMap(t, r.Body)
			if got := body["content"]; got != "memo from file" {
				t.Fatalf("content = %v, want memo from file", got)
			}
			return jsonResponse(http.StatusOK, `{"name":"memos/memo-file","content":"memo from file"}`), nil
		},
	}
	fileSystem := &mockFileSystem{
		readFileFunc: func(filePath string) ([]byte, error) {
			if filePath != "./memo.md" {
				t.Fatalf("filePath = %s, want ./memo.md", filePath)
			}
			return []byte("memo from file"), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "test-token",
		Timeout:    time.Second,
		HTTPClient: client,
		FileSystem: fileSystem,
	})

	result, err := service.CreateMemo(
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
		t.Fatalf("CreateMemo() error = %v", err)
	}
	if result.Name != "memos/memo-file" {
		t.Fatalf("name = %s, want memos/memo-file", result.Name)
	}
}

func TestService_CreateMemo_ContentFileReadError_Error(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			t.Fatal("HTTP call should not be executed")
			return nil, nil
		},
	}
	fileSystem := &mockFileSystem{
		readFileFunc: func(filePath string) ([]byte, error) {
			return nil, fmt.Errorf("read failed")
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "test-token",
		Timeout:    time.Second,
		HTTPClient: client,
		FileSystem: fileSystem,
	})

	_, err := service.CreateMemo(
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
		t.Fatal("CreateMemo() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content-file の読み込みに失敗しました") {
		t.Fatalf("error = %v, want content-file の読み込みに失敗しました", err)
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

	_, err := service.CreateMemo(context.Background(), "", "memo", "", "", "", nil, "")
	if err == nil {
		t.Fatal("CreateMemo() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "レスポンスデコード") {
		t.Fatalf("error = %v, want レスポンスデコード", err)
	}
}
