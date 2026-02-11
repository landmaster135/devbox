package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/vector_embedding/config"
	"github.com/landmaster135/devbox/internal/vector_embedding/domain"
	"github.com/landmaster135/devbox/internal/vector_embedding/infrastructure"
)

// Options は Service の設定値。
type Options struct {
	Host    string
	Port    int
	Timeout time.Duration

	Client infrastructure.EmbeddingClient
}

// Service は vector-embedding CLI のユースケースを提供する。
type Service struct {
	client infrastructure.EmbeddingClient
}

// NewService は Service を生成する。
func NewService(opts Options) (*Service, error) {
	client := opts.Client
	if client == nil {
		host := strings.TrimSpace(opts.Host)
		if host == "" {
			return nil, fmt.Errorf("host が指定されていません")
		}
		port := opts.Port
		if port <= 0 || port > 65535 {
			return nil, fmt.Errorf("port が不正です: %d", port)
		}
		timeout := opts.Timeout
		if timeout <= 0 {
			timeout = 60 * time.Second
		}
		baseURL := fmt.Sprintf("http://%s:%d", host, port)
		var err error
		client, err = infrastructure.NewOllamaClient(infrastructure.OllamaClientOptions{
			BaseURL: baseURL,
			Timeout: timeout,
		})
		if err != nil {
			return nil, fmt.Errorf("Ollama クライアントの初期化に失敗しました: %w", err)
		}
	}

	return &Service{client: client}, nil
}

// Run は設定に基づき埋め込み処理を実行する。
func (s *Service) Run(ctx context.Context, cfg *config.Config) (*domain.EmbedResult, error) {
	if s == nil {
		return nil, fmt.Errorf("Service が初期化されていません")
	}
	if cfg == nil {
		return nil, fmt.Errorf("設定情報が渡されていません")
	}

	request := domain.EmbedRequest{
		Operation: cfg.Operation,
		Model:     cfg.Model,
		Inputs:    cfg.Inputs,
	}

	switch cfg.Operation {
	case config.OperationOllama:
		return s.handleOllama(ctx, request)
	default:
		return nil, fmt.Errorf("未対応の operation です: %s", cfg.Operation)
	}
}

func (s *Service) handleOllama(ctx context.Context, req domain.EmbedRequest) (*domain.EmbedResult, error) {
	vectors, err := s.client.CreateEmbeddings(ctx, req.Model, req.Inputs)
	if err != nil {
		return nil, fmt.Errorf("Ollama への埋め込み要求に失敗しました: %w", err)
	}

	result := &domain.EmbedResult{
		Provider:   config.OperationOllama,
		Model:      req.Model,
		Embeddings: vectors,
		InputCount: len(req.Inputs),
	}
	if len(vectors) > 0 {
		result.Dimensions = len(vectors[0])
	}
	return result, nil
}
