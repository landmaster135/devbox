package usecases

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/genai"

	config "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/config"
)

// AIClient はAI OCRを実行するクライアントの共通インターフェース
type AIClient interface {
	Generate(ctx context.Context, req *AIRequest) (string, error)
	Close() error
}

// AIRequest はAIクライアントへの入力情報
type AIRequest struct {
	Prompt            string
	SystemInstruction string
	Temperature       float64
	MaxTokens         int
	Model             string
	MimeType          string
	ImageBase64       string
	ImageData         []byte
}

// #==============================================================#
// ##       Implementations for GeminiAIClient                   ##
// #==============================================================#
// GeminiAIClient はGemini API用クライアント
type GeminiAIClient struct {
	client *genai.Client
}

// NewGeminiAIClient はGeminiAIClientを生成する
func NewGeminiAIClient(cfg *config.Config) (AIClient, error) {
	clientConfig := &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	}

	client, err := genai.NewClient(context.Background(), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("geminiクライアントの作成に失敗しました: %w", err)
	}

	return &GeminiAIClient{client: client}, nil
}

// Generate はGemini APIでOCRを実行する
func (c *GeminiAIClient) Generate(ctx context.Context, req *AIRequest) (string, error) {
	return generateWithGenAI(ctx, c.client, req)
}

// Close はリソースを解放する（Geminiでは特に不要）
func (c *GeminiAIClient) Close() error {
	return nil
}

// #==============================================================#
// ##       Implementations for VertexAIClient                   ##
// #==============================================================#
// VertexAIClient はVertex AI用クライアント
type VertexAIClient struct {
	client *genai.Client
}

// NewVertexAIClient はVertexAIClientを生成する
func NewVertexAIClient(cfg *config.Config) (AIClient, error) {
	clientConfig := &genai.ClientConfig{
		Project:  cfg.Project,
		Location: cfg.Location,
		Backend:  genai.BackendVertexAI,
	}

	client, err := genai.NewClient(context.Background(), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("vertexクライアントの作成に失敗しました: %w", err)
	}

	return &VertexAIClient{client: client}, nil
}

// Generate はVertex AIでOCRを実行する
func (c *VertexAIClient) Generate(ctx context.Context, req *AIRequest) (string, error) {
	return generateWithGenAI(ctx, c.client, req)
}

// Close はリソースを解放する（Vertexでは特に不要）
func (c *VertexAIClient) Close() error {
	return nil
}

// #==============================================================#
// ##       Implementations for OllamaAIClient                   ##
// #==============================================================#
// OllamaAIClient はOllama API用クライアント
type OllamaAIClient struct {
	httpClient *http.Client
}

// NewOllamaAIClient はOllamaAIClientを生成する
func NewOllamaAIClient(_ *config.Config) (AIClient, error) {
	return &OllamaAIClient{
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}, nil
}

// Generate はOllama APIでOCRを実行する
func (c *OllamaAIClient) Generate(ctx context.Context, req *AIRequest) (string, error) {
	if c.httpClient == nil {
		return "", fmt.Errorf("Ollamaクライアントが初期化されていません")
	}

	payload := map[string]any{
		"model":  req.Model,
		"prompt": strings.TrimSpace(req.Prompt),
		"stream": true,
	}

	if system := strings.TrimSpace(req.SystemInstruction); system != "" {
		payload["system"] = system
	}
	if req.ImageBase64 != "" {
		payload["images"] = []string{req.ImageBase64}
	}

	options := map[string]any{}
	if req.Temperature >= 0 {
		options["temperature"] = req.Temperature
	}
	if req.MaxTokens > 0 {
		options["num_predict"] = req.MaxTokens
	}
	if len(options) > 0 {
		payload["options"] = options
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("リクエストボディの生成に失敗しました: %w", err)
	}

	reqHTTP, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaGenerateEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("ollamaリクエストの作成に失敗しました: %w", err)
	}
	reqHTTP.Header.Set("Content-Type", ollamaRequestContentType)

	resp, err := c.httpClient.Do(reqHTTP)
	if err != nil {
		return "", fmt.Errorf("ollama API呼び出しに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return "", fmt.Errorf("ollama APIがエラーを返しました: status=%d body=%s", resp.StatusCode, strings.TrimSpace(buf.String()))
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 1024), ollamaScannerBufferSize)

	var responseBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(strings.TrimSpace(string(line))) == 0 {
			continue
		}
		var chunk ollamaResponseChunk
		if err := json.Unmarshal(line, &chunk); err != nil {
			return "", fmt.Errorf("ollamaレスポンスの解析に失敗しました: %w", err)
		}
		if chunk.Error != "" {
			return "", fmt.Errorf("ollamaレスポンスエラー: %s", chunk.Error)
		}
		responseBuilder.WriteString(chunk.Response)
		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("ollamaレスポンスの読み取りに失敗しました: %w", err)
	}

	return responseBuilder.String(), nil
}

// Close はリソースを解放する（Ollamaでは特に不要）
func (c *OllamaAIClient) Close() error {
	return nil
}

const (
	ollamaGenerateEndpoint   = "http://localhost:11434/api/generate"
	ollamaScannerBufferSize  = 1 << 20 // 1 MiB
	ollamaRequestContentType = "application/json"
)

type ollamaResponseChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

func generateWithGenAI(ctx context.Context, client *genai.Client, req *AIRequest) (string, error) {
	if len(req.ImageData) == 0 {
		return "", fmt.Errorf("画像データが空です")
	}

	temperature := float32(req.Temperature)
	maxTokens := int32(req.MaxTokens)
	generateConfig := &genai.GenerateContentConfig{
		Temperature:     &temperature,
		MaxOutputTokens: maxTokens,
		SystemInstruction: &genai.Content{
			Role:  "system",
			Parts: []*genai.Part{{Text: req.SystemInstruction}},
		},
		SafetySettings: []*genai.SafetySetting{
			{
				Category:  "HARM_CATEGORY_HATE_SPEECH",
				Threshold: "OFF",
			},
			{
				Category:  "HARM_CATEGORY_DANGEROUS_CONTENT",
				Threshold: "OFF",
			},
			{
				Category:  "HARM_CATEGORY_SEXUALLY_EXPLICIT",
				Threshold: "OFF",
			},
			{
				Category:  "HARM_CATEGORY_HARASSMENT",
				Threshold: "OFF",
			},
		},
	}

	content := &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: req.Prompt},
			{
				InlineData: &genai.Blob{
					Data:     req.ImageData,
					MIMEType: req.MimeType,
				},
			},
		},
	}

	response, err := client.Models.GenerateContent(ctx, req.Model, []*genai.Content{content}, generateConfig)
	if err != nil {
		return "", fmt.Errorf("コンテンツ生成エラー: %w", err)
	}

	if len(response.Candidates) == 0 {
		return "", fmt.Errorf("レスポンスに候補が含まれていません")
	}

	candidate := response.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("レスポンスにコンテンツが含まれていません")
	}

	var builder strings.Builder
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			builder.WriteString(part.Text)
		}
	}

	return builder.String(), nil
}
