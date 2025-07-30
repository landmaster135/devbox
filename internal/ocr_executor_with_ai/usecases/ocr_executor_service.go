package usecases

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
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
	base64Service *Base64ExtractorService
	config        *config.Config
}

// NewOcrExecutorService は新しいOcrExecutorServiceを作成する
func NewOcrExecutorService(cfg *config.Config) (*OcrExecutorService, error) {
	var client *genai.Client
	var err error

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
	default:
		return nil, fmt.Errorf("無効なAIタイプです: %s (gemini または vertex を指定してください)", cfg.AiType)
	}

	if err != nil {
		return nil, fmt.Errorf("gemini APIクライアントの作成に失敗しました: %v", err)
	}

	// Base64変換サービスを作成
	base64Service := NewBase64ExtractorService(cfg.Path, cfg.Recursive)

	return &OcrExecutorService{
		client:        client,
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

	// Base64データをデコード
	imageData, err := base64.StdEncoding.DecodeString(imageResult.Base64)
	if err != nil {
		ocrResult.Error = fmt.Sprintf("Base64デコードエラー: %v", err)
		return ocrResult
	}

	// Gemini APIを呼び出してOCRを実行
	content, err := s.callGeminiAPI(imageData, mimeType)
	if err != nil {
		ocrResult.Error = fmt.Sprintf("Gemini API呼び出しエラー: %v", err)
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
