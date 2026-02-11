package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/ollama/domain"
)

func TestService_GetVersion(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/version", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"version":"0.1.30"}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	svc := NewService(ServiceOptions{BaseURL: server.URL, HTTPClient: server.Client(), Timeout: time.Second})
	got, err := svc.GetVersion(context.Background())
	if err != nil {
		t.Fatalf("GetVersion returned error: %v", err)
	}
	if got.Version != "0.1.30" {
		t.Fatalf("unexpected version: %s", got.Version)
	}
}

func TestService_CreateEmbeddings(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/embed", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var req domain.EmbedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Model != "nomic-embed" {
			t.Fatalf("unexpected model: %s", req.Model)
		}
		if len(req.Input) != 2 || req.Input[0] != "hello" {
			t.Fatalf("unexpected input: %#v", req.Input)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"model":"nomic-embed","embeddings":[[1.0,2.0]]}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	svc := NewService(ServiceOptions{BaseURL: server.URL, HTTPClient: server.Client(), Timeout: time.Second})
	resp, err := svc.CreateEmbeddings(context.Background(), domain.EmbedRequest{
		Model: "nomic-embed",
		Input: []string{"hello", "world"},
	})
	if err != nil {
		t.Fatalf("CreateEmbeddings returned error: %v", err)
	}
	if len(resp.Embeddings) != 1 {
		t.Fatalf("unexpected embeddings count: %d", len(resp.Embeddings))
	}
}

func TestService_Generate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var req domain.GenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if !req.Stream {
			t.Fatalf("stream flag must be true")
		}
		if req.Prompt != "Hi" {
			t.Fatalf("unexpected prompt: %s", req.Prompt)
		}
		flusher, _ := w.(http.Flusher)
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Write([]byte(`{"response":"Hello ","done":false}` + "\n"))
		if flusher != nil {
			flusher.Flush()
		}
		w.Write([]byte(`{"response":"world!","done":true}` + "\n"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	svc := NewService(ServiceOptions{BaseURL: server.URL, HTTPClient: server.Client(), Timeout: time.Second})
	out, err := svc.Generate(context.Background(), domain.GenerateRequest{
		Model:  "some-model",
		Prompt: "Hi",
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if out != "Hello world!" {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestService_StreamPull(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/pull", func(w http.ResponseWriter, r *http.Request) {
		var req domain.PullRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Model != "llama3" {
			t.Fatalf("unexpected model: %s", req.Model)
		}
		if !req.Stream {
			t.Fatalf("expected stream true but got false")
		}
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Write([]byte(`{"status":"downloading","total":100,"completed":50}` + "\n"))
		w.Write([]byte(`{"status":"success"}` + "\n"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	svc := NewService(ServiceOptions{BaseURL: server.URL, HTTPClient: server.Client(), Timeout: time.Second})
	var buf bytes.Buffer
	if err := svc.StreamPull(context.Background(), domain.PullRequest{Model: "llama3"}, &buf); err != nil {
		t.Fatalf("StreamPull returned error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "downloading 50.0% (50/100)") {
		t.Fatalf("unexpected progress output: %s", out)
	}
	if !strings.Contains(out, "success") {
		t.Fatalf("missing success status: %s", out)
	}
}

func TestService_ListRunningModels_APIError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/ps", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("bad gateway"))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	svc := NewService(ServiceOptions{BaseURL: server.URL, HTTPClient: server.Client(), Timeout: time.Second})
	if _, err := svc.ListRunningModels(context.Background()); err == nil {
		t.Fatal("expected error but got nil")
	}
}
