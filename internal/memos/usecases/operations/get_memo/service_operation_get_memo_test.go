package getmemo

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/memos/usecases/common"
	"github.com/landmaster135/devbox/internal/memos/usecases/testutil"
)

func TestServiceOperationGetMemo_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			if r.Method != http.MethodGet {
				t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
			}
			if r.URL.Path != "/api/v1/memos/memo-abc" {
				t.Fatalf("path = %s, want /api/v1/memos/memo-abc", r.URL.Path)
			}
			return testutil.JSONResponse(http.StatusOK, `{"name":"memos/memo-abc","content":"value"}`), nil
		},
	}

	jsonClient := common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com/api/v1",
		APIToken:   "token",
		HTTPClient: client,
	})
	service := New(jsonClient)

	result, err := service.Execute(context.Background(), "memos/memo-abc")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result.Content != "value" {
		t.Fatalf("content = %s, want value", result.Content)
	}
}

func TestServiceOperationGetMemo_APIError(t *testing.T) {
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

	_, err := service.Execute(context.Background(), "memo-1")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "status=401") {
		t.Fatalf("error = %v, want status=401", err)
	}
	if !strings.Contains(err.Error(), "unauthorized") {
		t.Fatalf("error = %v, want unauthorized", err)
	}
}

func TestServiceOperationGetMemo_DoJSONNetworkError_Error(t *testing.T) {
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

	_, err := service.Execute(context.Background(), "memo-1")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "network error") {
		t.Fatalf("error = %v, want network error", err)
	}
}

func TestServiceOperationGetMemo_EmptyMemo_Error(t *testing.T) {
	service := New(common.NewJSONClient(common.JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: &testutil.MockHTTPClient{},
	}))

	_, err := service.Execute(context.Background(), "")
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "memo が空") {
		t.Fatalf("error = %v, want memo が空", err)
	}
}
