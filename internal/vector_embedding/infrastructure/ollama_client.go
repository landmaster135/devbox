package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"time"

	ollamaDomain "github.com/landmaster135/devbox/internal/ollama/domain"
	ollamaUsecases "github.com/landmaster135/devbox/internal/ollama/usecases"
)

// EmbeddingClient は埋め込み生成クライアントを表す。
type EmbeddingClient interface {
	CreateEmbeddings(ctx context.Context, model string, inputs []string) ([][]float64, error)
}

// OllamaClientOptions は OllamaClient を初期化するためのオプション。
type OllamaClientOptions struct {
	BaseURL string
	Timeout time.Duration
}

// OllamaClient は Ollama API を呼び出して埋め込みを取得する。
type OllamaClient struct {
	service *ollamaUsecases.Service
}

// NewOllamaClient は OllamaClient を作成する。
func NewOllamaClient(opts OllamaClientOptions) (*OllamaClient, error) {
	baseURL := strings.TrimSpace(opts.BaseURL)
	if baseURL == "" {
		return nil, fmt.Errorf("Ollama の baseURL は必須です")
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	service := ollamaUsecases.NewService(ollamaUsecases.ServiceOptions{
		BaseURL: baseURL,
		Timeout: timeout,
	})

	return &OllamaClient{service: service}, nil
}

// CreateEmbeddings は Ollama の /api/embed を呼び出す。
func (c *OllamaClient) CreateEmbeddings(ctx context.Context, model string, inputs []string) ([][]float64, error) {
	if c == nil {
		return nil, fmt.Errorf("OllamaClient が初期化されていません")
	}
	req := ollamaDomain.EmbedRequest{Model: model, Input: inputs}
	resp, err := c.service.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Embeddings, nil
}
