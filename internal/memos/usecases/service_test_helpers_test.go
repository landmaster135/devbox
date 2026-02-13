package usecases

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

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
