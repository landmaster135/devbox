package usecases

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"google.golang.org/genai"

	config "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/config"
	models "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/models"
)

// OcrExecutorService はAI OCRを実行するサービス
type OcrExecutorService struct {
	client        *genai.Client
	httpClient    *http.Client
	base64Service *Base64ExtractorService
	config        *config.Config
}

// NewOcrExecutorService は新しいOcrExecutorServiceを作成する
func NewOcrExecutorService(cfg *config.Config) (*OcrExecutorService, error) {
	var (
		client     *genai.Client
		err        error
		httpClient *http.Client
	)

	// AIタイプに応じてクライアントを作成
	switch cfg.AiType {
	case "gemini":
		clientConfig := &genai.ClientConfig{
			APIKey:  cfg.APIKey,
			Backend: genai.BackendGeminiAPI,
		}
		client, err = genai.NewClient(context.Background(), clientConfig)
	case "vertex":
		clientConfig := &genai.ClientConfig{
			Project:  cfg.Project,
			Location: cfg.Location,
			Backend:  genai.BackendVertexAI,
		}
		client, err = genai.NewClient(context.Background(), clientConfig)
	case "ollama":
		httpClient = &http.Client{
			Timeout: 120 * time.Second,
		}
	default:
		return nil, fmt.Errorf("無効なAIタイプです: %s", cfg.AiType)
	}

	if err != nil {
		return nil, fmt.Errorf("AIクライアントの作成に失敗しました: %v", err)
	}

	// Base64変換サービスを作成
	base64Service := NewBase64ExtractorService(cfg.Path, cfg.Recursive)

	return &OcrExecutorService{
		client:        client,
		httpClient:    httpClient,
		base64Service: base64Service,
		config:        cfg,
	}, nil
}

// ProcessPath は指定されたパスの画像に対してOCRを実行する
func (s *OcrExecutorService) ProcessPath() (*models.OcrExecutionResult, error) {
	// 画像ファイルを検索してBase64に変換
	extractResult, err := s.base64Service.ExtractFromPath()
	if err != nil {
		return nil, fmt.Errorf("画像ファイルの検索に失敗しました: %v", err)
	}

	result := &models.OcrExecutionResult{}

	// 各画像に対してOCRを実行
	for _, imageResult := range extractResult.Images {
		ocrResult := s.processImage(imageResult)
		result.AddResult(ocrResult)

		// API呼び出し間隔を空ける（レート制限対策）
		if len(extractResult.Images) > 1 {
			time.Sleep(2 * time.Second)
		}
	}

	return result, nil
}

// processImage は単一画像に対してOCRを実行する
func (s *OcrExecutorService) processImage(imageResult ImageResult) models.OcrResult {
	ocrResult := models.OcrResult{
		FilePath: imageResult.FilePath,
	}

	// Base64変換でエラーが発生している場合
	if imageResult.Error != "" {
		ocrResult.Error = imageResult.Error
		return ocrResult
	}

	// MIMEタイプを取得
	mimeType := s.getMimeType(imageResult.FilePath)
	if mimeType == "" {
		ocrResult.Error = "サポートされていない画像形式です"
		return ocrResult
	}

	if s.config.AiType == "ollama" {
		content, err := s.callOllamaAPI(imageResult.Base64)
		if err != nil {
			ocrResult.Error = fmt.Sprintf("Ollama API呼び出しエラー: %v", err)
			return ocrResult
		}
		ocrResult.Content = content
		return ocrResult
	}

	// Base64データをデコード
	imageData, err := base64.StdEncoding.DecodeString(imageResult.Base64)
	if err != nil {
		ocrResult.Error = fmt.Sprintf("Base64デコードエラー: %v", err)
		return ocrResult
	}

	// Google AIを呼び出してOCRを実行
	content, err := s.callGeminiAPI(imageData, mimeType)
	if err != nil {
		ocrResult.Error = fmt.Sprintf("AI API呼び出しエラー: %v", err)
		return ocrResult
	}

	ocrResult.Content = content
	return ocrResult
}

// callGeminiAPI はGemini APIを呼び出してOCRを実行する
func (s *OcrExecutorService) callGeminiAPI(imageData []byte, mimeType string) (string, error) {
	ctx := context.Background()

	// 生成設定を作成
	temperature := float32(s.config.Temperature)
	maxTokens := int32(s.config.MaxTokens)
	generateConfig := &genai.GenerateContentConfig{
		Temperature:     &temperature,
		MaxOutputTokens: maxTokens,
		SystemInstruction: &genai.Content{
			Role:  "system",
			Parts: []*genai.Part{{Text: s.config.SystemInstruction}},
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

	// コンテンツを作成（テキストと画像）
	content := &genai.Content{
		Role: "user",
		Parts: []*genai.Part{
			{Text: s.config.Prompt},
			{
				InlineData: &genai.Blob{
					Data:     imageData,
					MIMEType: mimeType,
				},
			},
		},
	}

	// Gemini APIを呼び出し
	result, err := s.client.Models.GenerateContent(ctx, s.config.Model, []*genai.Content{content}, generateConfig)
	if err != nil {
		return "", fmt.Errorf("コンテンツ生成エラー: %v", err)
	}

	// レスポンスからテキストを抽出
	if len(result.Candidates) == 0 {
		return "", fmt.Errorf("レスポンスに候補が含まれていません")
	}

	candidate := result.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return "", fmt.Errorf("レスポンスにコンテンツが含まれていません")
	}

	var responseText strings.Builder
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			responseText.WriteString(part.Text)
		}
	}

	return responseText.String(), nil
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

// callOllamaAPI はOllama APIを呼び出してOCRを実行する
func (s *OcrExecutorService) callOllamaAPI(imageBase64 string) (string, error) {
	if s.httpClient == nil {
		return "", fmt.Errorf("Ollamaクライアントが初期化されていません")
	}

	payload := map[string]any{
		"model":  s.config.Model,
		"prompt": strings.TrimSpace(s.config.Prompt),
		"stream": true,
	}

	if system := strings.TrimSpace(s.config.SystemInstruction); system != "" {
		payload["system"] = system
	}
	if imageBase64 != "" {
		payload["images"] = []string{imageBase64}
	}

	options := map[string]any{}
	if s.config.Temperature >= 0 {
		options["temperature"] = s.config.Temperature
	}
	if s.config.MaxTokens > 0 {
		options["num_predict"] = s.config.MaxTokens
	}
	if len(options) > 0 {
		payload["options"] = options
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("リクエストボディの生成に失敗しました: %w", err)
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ollamaGenerateEndpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("Ollamaリクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", ollamaRequestContentType)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Ollama API呼び出しに失敗しました: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		return "", fmt.Errorf("Ollama APIがエラーを返しました: status=%d body=%s", resp.StatusCode, strings.TrimSpace(buf.String()))
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
			return "", fmt.Errorf("Ollamaレスポンスの解析に失敗しました: %w", err)
		}
		if chunk.Error != "" {
			return "", fmt.Errorf("Ollamaレスポンスエラー: %s", chunk.Error)
		}
		responseBuilder.WriteString(chunk.Response)
		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Ollamaレスポンスの読み取りに失敗しました: %w", err)
	}

	return responseBuilder.String(), nil
}

// getMimeType はファイル拡張子からMIMEタイプを取得する
func (s *OcrExecutorService) getMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	// mime.TypeByExtensionを使用してMIMEタイプを取得
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return mimeType
	}

	// 手動でマッピング（mime.TypeByExtensionで取得できない場合）
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".bmp":
		return "image/bmp"
	case ".webp":
		return "image/webp"
	default:
		return ""
	}
}

// Close はリソースをクリーンアップする
func (s *OcrExecutorService) Close() error {
	// genai.Clientにはクローズメソッドがないため、何もしない
	return nil
}
