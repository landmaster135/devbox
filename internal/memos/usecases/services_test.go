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

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
)

type mockHTTPClient struct {
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.doFunc(req)
}

type mockFileSystem struct {
	readFileFunc           func(filePath string) ([]byte, error)
	readAttachmentFileFunc func(filePath string) (*infrastructures.AttachmentFile, error)
}

func (m *mockFileSystem) ReadFile(filePath string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(filePath)
	}
	return nil, nil
}

func (m *mockFileSystem) ReadAttachmentFile(filePath string) (*infrastructures.AttachmentFile, error) {
	if m.readAttachmentFileFunc != nil {
		return m.readAttachmentFileFunc(filePath)
	}
	return nil, nil
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

func TestService_CreateAttachment_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPost {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPost)
			}
			if r.URL.Path != "/api/v1/attachments" {
				t.Fatalf("path = %s, want /api/v1/attachments", r.URL.Path)
			}
			body := readBodyAsMap(t, r.Body)
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
			return jsonResponse(http.StatusOK, `{"name":"attachments/1","filename":"hello.txt","type":"text/plain","memo":"memos/memo-1"}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	result, err := service.CreateAttachment(context.Background(), "hello.txt", []byte("hello"), "text/plain", "memo-1")
	if err != nil {
		t.Fatalf("CreateAttachment() error = %v", err)
	}
	if result.Name != "attachments/1" {
		t.Fatalf("name = %s, want attachments/1", result.Name)
	}
}

func TestService_ListMemoAttachments_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
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
			return jsonResponse(http.StatusOK, `{"attachments":[{"name":"attachments/1"},{"name":"attachments/2"}],"nextPageToken":"next-2"}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	result, err := service.ListMemoAttachments(context.Background(), "memos/memo-123", 100, "next")
	if err != nil {
		t.Fatalf("ListMemoAttachments() error = %v", err)
	}
	if len(result.Attachments) != 2 {
		t.Fatalf("len(attachments) = %d, want 2", len(result.Attachments))
	}
	if result.NextPageToken != "next-2" {
		t.Fatalf("nextPageToken = %s, want next-2", result.NextPageToken)
	}
}

func TestService_SetMemoAttachments_Normal(t *testing.T) {
	client := &mockHTTPClient{
		doFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodPatch {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodPatch)
			}
			if r.URL.Path != "/api/v1/memos/memo-123/attachments" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-123/attachments", r.URL.Path)
			}
			body := readBodyAsMap(t, r.Body)
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
			return jsonResponse(http.StatusOK, `{"name":"memos/memo-123","attachments":[{"name":"attachments/1"},{"name":"attachments/2"}]}`), nil
		},
	}

	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: client,
	})

	result, err := service.SetMemoAttachments(context.Background(), "memo-123", []Attachment{
		{Name: "attachments/1"},
		{Name: "attachments/2"},
	})
	if err != nil {
		t.Fatalf("SetMemoAttachments() error = %v", err)
	}
	if result.Name != "memos/memo-123" {
		t.Fatalf("name = %s, want memos/memo-123", result.Name)
	}
}

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

func TestService_ListMemoAttachments_EmptyMemo_Error(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: &mockHTTPClient{},
	})

	_, err := service.ListMemoAttachments(context.Background(), "", 100, "")
	if err == nil {
		t.Fatal("ListMemoAttachments() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}

func TestService_SetMemoAttachments_EmptyMemo_Error(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		Timeout:    time.Second,
		HTTPClient: &mockHTTPClient{},
	})

	_, err := service.SetMemoAttachments(context.Background(), "", nil)
	if err == nil {
		t.Fatal("SetMemoAttachments() error = nil, want error")
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
	if service.fileSystem == nil {
		t.Fatal("fileSystem is nil")
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

func TestBuildMemoResourceName_Normal(t *testing.T) {
	if got := buildMemoResourceName("https://example.com/api/v1/memos/memo-1"); got != "memos/memo-1" {
		t.Fatalf("buildMemoResourceName() = %s, want memos/memo-1", got)
	}
	if got := buildMemoResourceName("  "); got != "" {
		t.Fatalf("buildMemoResourceName() = %q, want empty", got)
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
