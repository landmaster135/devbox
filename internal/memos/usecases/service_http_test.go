package usecases

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

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
