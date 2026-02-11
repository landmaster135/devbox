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
	svc, err := NewService(Options{Client: client})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	cfg := &config.Config{
		Operation: config.OperationOllama,
		Model:     "abc",
		Inputs:    []string{"hello"},
	}

	got, err := svc.Run(context.Background(), cfg)
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
	if svc.client == nil {
		t.Fatal("client should be initialized")
	}
}

func TestServiceRun_Error(t *testing.T) {
	wantErr := errors.New("boom")
	client := &mockEmbeddingClient{err: wantErr}
	svc, _ := NewService(Options{Client: client})

	cfg := &config.Config{Operation: config.OperationOllama, Model: "abc", Inputs: []string{"x"}}
	if _, err := svc.Run(context.Background(), cfg); err == nil {
		t.Fatal("expected error but got nil")
	}
}
