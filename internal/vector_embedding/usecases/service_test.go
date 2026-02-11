package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/vector_embedding/config"
)

type mockEmbeddingClient struct {
	vectors [][]float64
	err     error
}

func (m *mockEmbeddingClient) CreateEmbeddings(ctx context.Context, model string, inputs []string) ([][]float64, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.vectors, nil
}

func TestServiceRun_Ollama(t *testing.T) {
	client := &mockEmbeddingClient{vectors: [][]float64{{1, 2, 3}}}
	svc, err := NewService(Options{OllamaClient: client})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	cfg := &config.Config{
		Operation: config.OperationOllama,
		Model:     "abc",
		Inputs:    []string{"hello"},
	}

	got, err := svc.Embed(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got.Provider != config.OperationOllama {
		t.Fatalf("unexpected provider: %s", got.Provider)
	}
	if got.Dimensions != 3 {
		t.Fatalf("unexpected dimensions: %d", got.Dimensions)
	}
}

func TestNewService_DefaultClientInit(t *testing.T) {
	svc, err := NewService(Options{Host: "127.0.0.1", Port: 11434, Timeout: time.Second})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	if svc.ollamaClient == nil {
		t.Fatal("ollama client should be initialized")
	}
}

func TestServiceRun_Error(t *testing.T) {
	wantErr := errors.New("boom")
	client := &mockEmbeddingClient{err: wantErr}
	svc, _ := NewService(Options{OllamaClient: client})

	cfg := &config.Config{Operation: config.OperationOllama, Model: "abc", Inputs: []string{"x"}}
	if _, err := svc.Embed(context.Background(), cfg); err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestServiceRun_OpenAI(t *testing.T) {
	client := &mockEmbeddingClient{vectors: [][]float64{{0.1, 0.2}}}
	svc, err := NewService(Options{OpenAIClient: client})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	cfg := &config.Config{
		Operation: config.OperationOpenAI,
		Model:     "text-embedding-3-small",
		Inputs:    []string{"hello"},
		APIKey:    "sk-test",
	}

	result, err := svc.Embed(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Provider != config.OperationOpenAI {
		t.Fatalf("unexpected provider: %s", result.Provider)
	}
}

func TestServiceRun_OpenAIWithoutClient(t *testing.T) {
	svc, err := NewService(Options{})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	cfg := &config.Config{
		Operation: config.OperationOpenAI,
		Model:     "text-embedding-3-small",
		Inputs:    []string{"hello"},
	}

	if _, err := svc.Embed(context.Background(), cfg); err == nil {
		t.Fatal("expected error due to missing api key")
	}
}
