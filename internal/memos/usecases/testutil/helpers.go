package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
)

// MockHTTPClient は HTTP クライアントのテストダブル。
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}

// MockFileSystem は FileSystem のテストダブル。
type MockFileSystem struct {
	ReadFileFunc           func(filePath string) ([]byte, error)
	ReadAttachmentFileFunc func(filePath string) (*infrastructures.AttachmentFile, error)
}

func (m *MockFileSystem) ReadFile(filePath string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(filePath)
	}
	return nil, nil
}

func (m *MockFileSystem) ReadAttachmentFile(filePath string) (*infrastructures.AttachmentFile, error) {
	if m.ReadAttachmentFileFunc != nil {
		return m.ReadAttachmentFileFunc(filePath)
	}
	return nil, nil
}

func ReadBodyAsMap(t *testing.T, body io.ReadCloser) map[string]any {
	t.Helper()
	defer body.Close()

	var m map[string]any
	if err := json.NewDecoder(body).Decode(&m); err != nil {
		t.Fatalf("decode body error: %v", err)
	}
	return m
}

func JSONResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}
