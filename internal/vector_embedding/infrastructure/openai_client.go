package infrastructure

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

// OpenAIClientOptions は OpenAIClient の初期化時に指定するオプション。
type OpenAIClientOptions struct {
	APIKey  string
	Timeout time.Duration
	BaseURL string
}

// OpenAIClient は OpenAI Embeddings API を呼び出すクライアント。
type OpenAIClient struct {
	client *openai.Client
}

// NewOpenAIClient は OpenAIClient を生成する。
func NewOpenAIClient(opts OpenAIClientOptions) (*OpenAIClient, error) {
	apiKey := strings.TrimSpace(opts.APIKey)
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI の API キーは必須です")
	}

	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	httpClient := &http.Client{Timeout: timeout}

	requestOptions := []option.RequestOption{
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
	}

	if baseURL := strings.TrimSpace(opts.BaseURL); baseURL != "" {
		requestOptions = append(requestOptions, option.WithBaseURL(baseURL))
	}

	client := openai.NewClient(requestOptions...)
	return &OpenAIClient{client: &client}, nil
}

// CreateEmbeddings は OpenAI の embeddings エンドポイントを呼び出してベクトルを生成する。
func (c *OpenAIClient) CreateEmbeddings(ctx context.Context, model string, inputs []string) ([][]float64, error) {
	if c == nil || c.client == nil {
		return nil, fmt.Errorf("OpenAIClient が初期化されていません")
	}
	trimmedModel := strings.TrimSpace(model)
	if trimmedModel == "" {
		return nil, fmt.Errorf("model は必須です")
	}
	if len(inputs) == 0 {
		return nil, fmt.Errorf("input は 1 件以上が必要です")
	}

	params := openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel(trimmedModel),
		Input: openai.EmbeddingNewParamsInputUnion{OfArrayOfStrings: inputs},
	}

	resp, err := c.client.Embeddings.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("OpenAI Embeddings API の呼び出しに失敗しました: %w", err)
	}

	vectors := make([][]float64, len(resp.Data))
	for i, data := range resp.Data {
		vec := make([]float64, len(data.Embedding))
		copy(vec, data.Embedding)
		vectors[i] = vec
	}
	return vectors, nil
}
