package usecases

import (
	"errors"
	"testing"
)

type mockMigrateToMemosOperation struct {
	ExecuteFunc func(pageType, baseURL, apiToken, srcBodyDir, srcResourceDir string) (string, error)
}

func (m *mockMigrateToMemosOperation) Execute(pageType, baseURL, apiToken, srcBodyDir, srcResourceDir string) (string, error) {
	return m.ExecuteFunc(pageType, baseURL, apiToken, srcBodyDir, srcResourceDir)
}

func TestService_MigrateToMemos_Normal(t *testing.T) {
	t.Parallel()

	op := &mockMigrateToMemosOperation{
		ExecuteFunc: func(pageType, baseURL, apiToken, srcBodyDir, srcResourceDir string) (string, error) {
			if pageType != "content" ||
				baseURL != "https://memos.example.com" ||
				apiToken != "token" ||
				srcBodyDir != "/tmp/body" ||
				srcResourceDir != "/tmp/resources" {
				t.Fatalf(
					"unexpected args: pageType=%s baseURL=%s apiToken=%s srcBodyDir=%s srcResourceDir=%s",
					pageType,
					baseURL,
					apiToken,
					srcBodyDir,
					srcResourceDir,
				)
			}
			return "ok", nil
		},
	}

	service := newServiceWithOperations(nil, nil, nil, nil, nil, op)
	got, err := service.MigrateToMemos("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "ok" {
		t.Fatalf("got = %q, want ok", got)
	}
}

func TestService_MigrateToMemos_Error(t *testing.T) {
	t.Parallel()

	op := &mockMigrateToMemosOperation{
		ExecuteFunc: func(pageType, baseURL, apiToken, srcBodyDir, srcResourceDir string) (string, error) {
			return "", errors.New("failed")
		},
	}

	service := newServiceWithOperations(nil, nil, nil, nil, nil, op)
	_, err := service.MigrateToMemos("content", "https://memos.example.com", "token", "/tmp/body", "/tmp/resources")
	if err == nil || err.Error() != "failed" {
		t.Fatalf("error = %v", err)
	}
}
