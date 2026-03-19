package deletememo

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	testutil "github.com/landmaster135/devbox/internal/memos/infrastructures/testutil"
	common "github.com/landmaster135/devbox/internal/memos/usecases/common"
)

func TestServiceOperationDeleteMemo_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodDelete {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodDelete)
			}
			if r.URL.Path != "/api/v1/memos/memo-abc" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-abc", r.URL.Path)
			}
			if got := r.URL.Query().Get("force"); got != "true" {
				t.Fatalf("force query = %s, want true", got)
			}
			return testutil.JSONResponse(http.StatusOK, ``), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	result, err := service.Execute(context.Background(), "memos/memo-abc", true)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
}

func TestServiceOperationDeleteMemo_ForceFalse_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if got := r.URL.Query().Get("force"); got != "" {
				t.Fatalf("force query = %s, want empty", got)
			}
			return testutil.JSONResponse(http.StatusOK, ``), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	result, err := service.Execute(context.Background(), "memo-abc", false)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result == nil {
		t.Fatal("result = nil, want non-nil")
	}
}

func TestServiceOperationDeleteMemo_APIError_Error(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return testutil.JSONResponse(http.StatusUnauthorized, `{"message":"unauthorized"}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	_, err := service.Execute(context.Background(), "memo-1", false)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "status=401") {
		t.Fatalf("error = %v, want status=401", err)
	}
}

func TestServiceOperationDeleteMemo_DoJSONNetworkError_Error(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("network error")
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	_, err := service.Execute(context.Background(), "memo-1", false)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Fatalf("error = %v, want network error", err)
	}
}

func TestServiceOperationDeleteMemo_EmptyMemo_Error(t *testing.T) {
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: &testutil.MockHTTPClient{},
	}))

	_, err := service.Execute(context.Background(), "", false)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}
