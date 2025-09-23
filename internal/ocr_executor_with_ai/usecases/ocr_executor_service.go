package usecases

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"path/filepath"
	"strings"
	"time"

	config "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/config"
	models "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/models"
)

// OcrExecutorService はAI OCRを実行するサービス
type OcrExecutorService struct {
	client        AIClient
	base64Service *Base64ExtractorService
	config        *config.Config
}

// NewOcrExecutorService は新しいOcrExecutorServiceを作成する
func NewOcrExecutorService(cfg *config.Config) (*OcrExecutorService, error) {
	var (
		client AIClient
		err    error
	)

	switch cfg.AiType {
	case "gemini":
		client, err = NewGeminiAIClient(cfg)
	case "vertex":
		client, err = NewVertexAIClient(cfg)
	case "ollama":
		client, err = NewOllamaAIClient(cfg)
	default:
		return nil, fmt.Errorf("無効なAIタイプです: %s", cfg.AiType)
	}

	if err != nil {
		return nil, fmt.Errorf("AIクライアントの作成に失敗しました: %w", err)
	}

	base64Service := NewBase64ExtractorService(cfg.Path, cfg.Recursive)

	return &OcrExecutorService{
		client:        client,
		base64Service: base64Service,
		config:        cfg,
	}, nil
}

// ProcessPath は指定されたパスの画像に対してOCRを実行する
func (s *OcrExecutorService) ProcessPath() (*models.OcrExecutionResult, error) {
	extractResult, err := s.base64Service.ExtractFromPath()
	if err != nil {
		return nil, fmt.Errorf("画像ファイルの検索に失敗しました: %v", err)
	}

	result := &models.OcrExecutionResult{}

	for _, imageResult := range extractResult.Images {
		ocrResult := s.processImage(imageResult)
		result.AddResult(ocrResult)

		if len(extractResult.Images) > 1 {
			time.Sleep(2 * time.Second)
		}
	}

	return result, nil
}

// processImage は単一画像に対してOCRを実行する
func (s *OcrExecutorService) processImage(imageResult ImageResult) models.OcrResult {
	ocrResult := models.OcrResult{FilePath: imageResult.FilePath}

	if imageResult.Error != "" {
		ocrResult.Error = imageResult.Error
		return ocrResult
	}

	mimeType := s.getMimeType(imageResult.FilePath)
	if mimeType == "" {
		ocrResult.Error = "サポートされていない画像形式です"
		return ocrResult
	}

	imageData, err := base64.StdEncoding.DecodeString(imageResult.Base64)
	if err != nil {
		ocrResult.Error = fmt.Sprintf("Base64デコードエラー: %v", err)
		return ocrResult
	}

	req := &AIRequest{
		Prompt:            s.config.Prompt,
		SystemInstruction: s.config.SystemInstruction,
		Temperature:       s.config.Temperature,
		MaxTokens:         s.config.MaxTokens,
		Model:             s.config.Model,
		MimeType:          mimeType,
		ImageBase64:       imageResult.Base64,
		ImageData:         imageData,
	}

	content, err := s.client.Generate(context.Background(), req)
	if err != nil {
		ocrResult.Error = fmt.Sprintf("AI API呼び出しエラー: %v", err)
		return ocrResult
	}

	ocrResult.Content = content
	return ocrResult
}

// getMimeType はファイル拡張子からMIMEタイプを取得する
func (s *OcrExecutorService) getMimeType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return mimeType
	}

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
	if s.client == nil {
		return nil
	}
	return s.client.Close()
}
