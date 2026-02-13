package usecases

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
)

func TestService_PatchFiles_ReplacesTrue_Normal(t *testing.T) {
	callCount := 0
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			callCount++
			switch callCount {
			case 1:
				if r.Method != http.MethodPost {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
				}
				if r.URL.Path != "/api/v1/attachments" {
					t.Fatalf("path = %s, want /api/v1/attachments", r.URL.Path)
				}
				return jsonResponse(http.StatusOK, `{"name":"attachments/new-1","filename":"a.txt"}`), nil
			case 2:
				if r.Method != http.MethodPatch {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodPatch)
				}
				if r.URL.Path != "/api/v1/memos/memo-1/attachments" {
					t.Fatalf("path = %s, want /api/v1/memos/memo-1/attachments", r.URL.Path)
				}
				body := readBodyAsMap(t, r.Body)
				attachments, ok := body["attachments"].([]any)
				if !ok {
					t.Fatalf("attachments type = %T, want []any", body["attachments"])
				}
				if len(attachments) != 1 {
					t.Fatalf("len(attachments) = %d, want 1", len(attachments))
				}
				return jsonResponse(http.StatusOK, `{"name":"memos/memo-1","attachments":[{"name":"attachments/new-1"}]}`), nil
			default:
				t.Fatalf("unexpected HTTP call: %d", callCount)
				return nil, nil
			}
		},
	}
	fileSystem := &mockFileSystem{
		readAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
			if filePath != "./a.txt" {
				t.Fatalf("filePath = %s, want ./a.txt", filePath)
			}
			return &infrastructures.AttachmentFile{
				Filename:    "a.txt",
				Content:     []byte("hello"),
				ContentType: "text/plain",
			}, nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
		FileSystem: fileSystem,
	})

	result, err := service.PatchFiles(context.Background(), "memo-1", []string{"./a.txt"}, true)
	if err != nil {
		t.Fatalf("PatchFiles() error = %v", err)
	}
	if result.Name != "memos/memo-1" {
		t.Fatalf("name = %s, want memos/memo-1", result.Name)
	}
}

func TestService_PatchFiles_ReplacesFalse_Normal(t *testing.T) {
	callCount := 0
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			callCount++
			switch callCount {
			case 1:
				if r.Method != http.MethodPost {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
				}
				if r.URL.Path != "/api/v1/attachments" {
					t.Fatalf("path = %s, want /api/v1/attachments", r.URL.Path)
				}
				return jsonResponse(http.StatusOK, `{"name":"attachments/new","filename":"a.txt"}`), nil
			case 2:
				if r.Method != http.MethodGet {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
				}
				if r.URL.Path != "/api/v1/memos/memo-1/attachments" {
					t.Fatalf("path = %s, want /api/v1/memos/memo-1/attachments", r.URL.Path)
				}
				if got := r.URL.Query().Get("pageSize"); got != "100" {
					t.Fatalf("pageSize = %s, want 100", got)
				}
				if got := r.URL.Query().Get("pageToken"); got != "" {
					t.Fatalf("pageToken = %s, want empty", got)
				}
				return jsonResponse(http.StatusOK, `{"attachments":[{"name":"attachments/existing-1"}],"nextPageToken":"next"}`), nil
			case 3:
				if r.Method != http.MethodGet {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
				}
				if got := r.URL.Query().Get("pageToken"); got != "next" {
					t.Fatalf("pageToken = %s, want next", got)
				}
				return jsonResponse(http.StatusOK, `{"attachments":[{"name":"attachments/new"},{"name":"attachments/existing-2"}]}`), nil
			case 4:
				if r.Method != http.MethodPatch {
					t.Fatalf("method = %s, want %s", r.Method, http.MethodPatch)
				}
				body := readBodyAsMap(t, r.Body)
				attachments, ok := body["attachments"].([]any)
				if !ok {
					t.Fatalf("attachments type = %T, want []any", body["attachments"])
				}
				if len(attachments) != 3 {
					t.Fatalf("len(attachments) = %d, want 3", len(attachments))
				}
				first, ok := attachments[0].(map[string]any)
				if !ok || first["name"] != "attachments/existing-1" {
					t.Fatalf("attachments[0] = %v, want attachments/existing-1", attachments[0])
				}
				second, ok := attachments[1].(map[string]any)
				if !ok || second["name"] != "attachments/existing-2" {
					t.Fatalf("attachments[1] = %v, want attachments/existing-2", attachments[1])
				}
				third, ok := attachments[2].(map[string]any)
				if !ok || third["name"] != "attachments/new" {
					t.Fatalf("attachments[2] = %v, want attachments/new", attachments[2])
				}
				return jsonResponse(http.StatusOK, `{"name":"memos/memo-1","attachments":[{"name":"attachments/existing-1"},{"name":"attachments/existing-2"},{"name":"attachments/new"}]}`), nil
			default:
				t.Fatalf("unexpected HTTP call: %d", callCount)
				return nil, nil
			}
		},
	}
	fileSystem := &mockFileSystem{
		readAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
			return &infrastructures.AttachmentFile{
				Filename:    "a.txt",
				Content:     []byte("hello"),
				ContentType: "text/plain",
			}, nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
		FileSystem: fileSystem,
	})

	result, err := service.PatchFiles(context.Background(), "memo-1", []string{"./a.txt"}, false)
	if err != nil {
		t.Fatalf("PatchFiles() error = %v", err)
	}
	if result.Name != "memos/memo-1" {
		t.Fatalf("name = %s, want memos/memo-1", result.Name)
	}
}

func TestService_PatchFiles_ReadAttachmentError_Error(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			t.Fatal("HTTP should not be called")
			return nil, nil
		},
	}
	fileSystem := &mockFileSystem{
		readAttachmentFileFunc: func(filePath string) (*infrastructures.AttachmentFile, error) {
			return nil, fmt.Errorf("read failed")
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
		FileSystem: fileSystem,
	})

	_, err := service.PatchFiles(context.Background(), "memo-1", []string{"./a.txt"}, true)
	if err == nil {
		t.Fatal("PatchFiles() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "files の読み込みに失敗しました") {
		t.Fatalf("error = %v, want files の読み込みに失敗しました", err)
	}
}

func TestMergeAttachmentsByName_Normal(t *testing.T) {
	existing := []Attachment{
		{Name: "attachments/existing-1"},
		{Name: "attachments/new"},
		{Name: "attachments/existing-2"},
	}
	created := []Attachment{
		{Name: "attachments/new"},
	}

	got := mergeAttachmentsByName(existing, created)
	if len(got) != 3 {
		t.Fatalf("len(got) = %d, want 3", len(got))
	}
	if got[0].Name != "attachments/existing-1" || got[1].Name != "attachments/existing-2" || got[2].Name != "attachments/new" {
		t.Fatalf("merge result = %+v, want existing-1, existing-2, new", got)
	}
}
