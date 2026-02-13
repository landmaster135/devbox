package common

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/memos/usecases/testutil"
)

func TestJSONClient_DoJSONOutNil_Normal(t *testing.T) {
	client := &testutil.MockHTTPClient{
		DoFunc: func(r *http.Request) (*http.Response, error) {
			return testutil.JSONResponse(http.StatusOK, `{"ok":true}`), nil
		},
	}

	jsonClient := NewJSONClient(JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: client,
	})

	err := jsonClient.DoJSON(context.Background(), http.MethodGet, "/memos", nil, nil, nil)
	if err != nil {
		t.Fatalf("DoJSON() error = %v", err)
	}
}

func TestJSONClient_NewRequestEncodeError_Error(t *testing.T) {
	jsonClient := NewJSONClient(JSONClientOptions{
		BaseURL:    "https://memos.example.com",
		APIToken:   "token",
		HTTPClient: &testutil.MockHTTPClient{},
	})

	_, err := jsonClient.NewRequest(context.Background(), http.MethodPost, "/memos", nil, map[string]any{
		"invalid": make(chan int),
	})
	if err == nil {
		t.Fatal("NewRequest() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "エンコード") {
		t.Fatalf("error = %v, want エンコード", err)
	}
}
