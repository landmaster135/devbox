package listusertags

import (
	"context"
	"errors"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	testutil "github.com/landmaster135/devbox/internal/memos/infrastructures/testutil"
	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

func TestServiceOperationListUserTags_Normal(t *testing.T) {
	var writtenPath string
	var writtenContent string

	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/users/1:getStats" {
				t.Fatalf("path = %s, want /api/v1/users/1:getStats", r.URL.Path)
			}
			if got := r.Header.Get("Authorization"); got != "Bearer token" {
				t.Fatalf("Authorization = %s, want Bearer token", got)
			}
			return testutil.JSONResponse(http.StatusOK, `{"name":"users/1","tagCount":{"work":2,"book":1},"totalMemoCount":3}`), nil
		},
	}
	fileSystem := &testutil.MockFileSystem{
		EnsureDirectoryFunc: func(dirPath string) (string, error) {
			if dirPath != "/tmp/out" {
				t.Fatalf("dirPath = %s, want /tmp/out", dirPath)
			}
			return "/tmp/out", nil
		},
		WriteFileFunc: func(filePath string, content []byte) error {
			writtenPath = filePath
			writtenContent = string(content)
			return nil
		},
	}
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	}), fileSystem)

	result, err := service.Execute(context.Background(), " 1 ", "/tmp/out")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	wantPath := filepath.Join("/tmp/out", "user-tags_1.json")
	if result.OutputPath != wantPath {
		t.Fatalf("outputPath = %s, want %s", result.OutputPath, wantPath)
	}
	if writtenPath != wantPath {
		t.Fatalf("writtenPath = %s, want %s", writtenPath, wantPath)
	}
	if !strings.Contains(writtenContent, `"work": 2`) {
		t.Fatalf("writtenContent = %s, want work count", writtenContent)
	}
	if result.TagCount["book"] != 1 {
		t.Fatalf("tagCount[book] = %d, want 1", result.TagCount["book"])
	}
}

func TestServiceOperationListUserTags_EscapedUserID_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.URL.EscapedPath() != "/api/v1/users/team%2F1:getStats" {
				t.Fatalf("escapedPath = %s, want /api/v1/users/team%%2F1:getStats", r.URL.EscapedPath())
			}
			return testutil.JSONResponse(http.StatusOK, `{"tagCount":{}}`), nil
		},
	}
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	}), &testutil.MockFileSystem{
		EnsureDirectoryFunc: func(dirPath string) (string, error) {
			return "/tmp/out", nil
		},
	})

	result, err := service.Execute(context.Background(), "team/1", "/tmp/out")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.HasSuffix(result.OutputPath, "user-tags_team_1.json") {
		t.Fatalf("outputPath = %s, want sanitized filename", result.OutputPath)
	}
}

func TestServiceOperationListUserTags_EmptyTagCount_Normal(t *testing.T) {
	var writtenContent string
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return testutil.JSONResponse(http.StatusOK, `{"name":"users/1"}`), nil
		},
	}
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	}), &testutil.MockFileSystem{
		EnsureDirectoryFunc: func(dirPath string) (string, error) {
			return "/tmp/out", nil
		},
		WriteFileFunc: func(filePath string, content []byte) error {
			writtenContent = string(content)
			return nil
		},
	})

	result, err := service.Execute(context.Background(), "1", "/tmp/out")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.TagCount) != 0 {
		t.Fatalf("len(tagCount) = %d, want 0", len(result.TagCount))
	}
	if writtenContent != "{}\n" {
		t.Fatalf("writtenContent = %q, want {} newline", writtenContent)
	}
}

func TestServiceOperationListUserTags_FileSystemError_Error(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return testutil.JSONResponse(http.StatusOK, `{"tagCount":{"work":2}}`), nil
		},
	}
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	}), &testutil.MockFileSystem{
		EnsureDirectoryFunc: func(dirPath string) (string, error) {
			return "", errors.New("directory error")
		},
	})

	_, err := service.Execute(context.Background(), "1", "/tmp/out")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "directory error") {
		t.Fatalf("error = %v, want directory error", err)
	}
}
