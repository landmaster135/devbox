package usecases

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

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
