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
	APIKey  string

	OllamaClient infrastructure.EmbeddingClient
	OpenAIClient infrastructure.EmbeddingClient
}

// Service は vector-embedding CLI のユースケースを提供する。
type Service struct {
	host    string
	port    int
	timeout time.Duration
	apiKey  string

	ollamaClient infrastructure.EmbeddingClient
	openAIClient infrastructure.EmbeddingClient
}

// NewService は Service を生成する。
func NewService(opts Options) (*Service, error) {
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	svc := &Service{
		host:         strings.TrimSpace(opts.Host),
		port:         opts.Port,
		timeout:      timeout,
		apiKey:       strings.TrimSpace(opts.APIKey),
		ollamaClient: opts.OllamaClient,
		openAIClient: opts.OpenAIClient,
	}

	if svc.ollamaClient == nil && svc.host != "" && svc.port > 0 && svc.port <= 65535 {
		client, err := svc.initOllamaClient()
		if err != nil {
			return nil, err
		}
		svc.ollamaClient = client
	}

	if svc.openAIClient == nil && svc.apiKey != "" {
		client, err := svc.initOpenAIClient(svc.apiKey)
		if err != nil {
			return nil, err
		}
		svc.openAIClient = client
	}

	return svc, nil
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
	case config.OperationOpenAI:
		return s.handleOpenAI(ctx, request, cfg.APIKey)
	default:
		return nil, fmt.Errorf("未対応の operation です: %s", cfg.Operation)
	}
}

func (s *Service) handleOllama(ctx context.Context, req domain.EmbedRequest) (*domain.EmbedResult, error) {
	client, err := s.ensureOllamaClient()
	if err != nil {
		return nil, err
	}
	return s.executeEmbedding(ctx, client, req, config.OperationOllama, "Ollama への埋め込み要求に失敗しました")
}

func (s *Service) handleOpenAI(ctx context.Context, req domain.EmbedRequest, apiKey string) (*domain.EmbedResult, error) {
	client, err := s.ensureOpenAIClient(apiKey)
	if err != nil {
		return nil, err
	}
	return s.executeEmbedding(ctx, client, req, config.OperationOpenAI, "OpenAI への埋め込み要求に失敗しました")
}

func (s *Service) executeEmbedding(
	ctx context.Context,
	client infrastructure.EmbeddingClient,
	req domain.EmbedRequest,
	provider string,
	errMsg string,
) (*domain.EmbedResult, error) {
	if client == nil {
		return nil, fmt.Errorf("%s クライアントが初期化されていません", provider)
	}

	vectors, err := client.CreateEmbeddings(ctx, req.Model, req.Inputs)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", errMsg, err)
	}

	result := &domain.EmbedResult{
		Provider:   provider,
		Model:      req.Model,
		Embeddings: vectors,
		InputCount: len(req.Inputs),
	}
	if len(vectors) > 0 {
		result.Dimensions = len(vectors[0])
	}
	return result, nil
}

func (s *Service) ensureOllamaClient() (infrastructure.EmbeddingClient, error) {
	if s.ollamaClient != nil {
		return s.ollamaClient, nil
	}
	client, err := s.initOllamaClient()
	if err != nil {
		return nil, err
	}
	s.ollamaClient = client
	return client, nil
}

func (s *Service) initOllamaClient() (infrastructure.EmbeddingClient, error) {
	host := strings.TrimSpace(s.host)
	if host == "" {
		return nil, fmt.Errorf("host が指定されていません")
	}
	port := s.port
	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("port が不正です: %d", port)
	}

	baseURL := fmt.Sprintf("http://%s:%d", host, port)
	client, err := infrastructure.NewOllamaClient(infrastructure.OllamaClientOptions{
		BaseURL: baseURL,
		Timeout: s.timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("Ollama クライアントの初期化に失敗しました: %w", err)
	}
	return client, nil
}

func (s *Service) ensureOpenAIClient(apiKey string) (infrastructure.EmbeddingClient, error) {
	if s.openAIClient != nil {
		return s.openAIClient, nil
	}

	key := strings.TrimSpace(apiKey)
	if key == "" {
		key = s.apiKey
	}
	if key == "" {
		return nil, fmt.Errorf("OpenAI の API キーが指定されていません")
	}

	client, err := s.initOpenAIClient(key)
	if err != nil {
		return nil, err
	}
	s.apiKey = key
	s.openAIClient = client
	return client, nil
}

func (s *Service) initOpenAIClient(apiKey string) (infrastructure.EmbeddingClient, error) {
	client, err := infrastructure.NewOpenAIClient(infrastructure.OpenAIClientOptions{
		APIKey:  apiKey,
		Timeout: s.timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("OpenAI クライアントの初期化に失敗しました: %w", err)
	}
	return client, nil
}
