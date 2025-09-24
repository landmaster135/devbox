package usecases

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/genai"

	config "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/config"
)

type stubAIClient struct {
	reqs      []*AIRequest
	responses []string
	errs      []error
	callCount int
}

func (s *stubAIClient) Generate(ctx context.Context, req *AIRequest) (string, error) {
	s.reqs = append(s.reqs, req)
	var resp string
	var err error
	if s.callCount < len(s.responses) {
		resp = s.responses[s.callCount]
	}
	if s.callCount < len(s.errs) {
		err = s.errs[s.callCount]
	}
	s.callCount++
	return resp, err
}

func (s *stubAIClient) Close() error {
	return nil
}

type closeTrackingClient struct {
	closeFn func() error
}

func (c *closeTrackingClient) Generate(ctx context.Context, req *AIRequest) (string, error) {
	return "", nil
}

func (c *closeTrackingClient) Close() error {
	if c.closeFn != nil {
		return c.closeFn()
	}
	return nil
}

func TestOcrExecutorService_ProcessImage_UsesInjectedClient(t *testing.T) {
	imageData := []byte("dummy data")
	encoded := base64.StdEncoding.EncodeToString(imageData)

	stubClient := &stubAIClient{responses: []string{"ok"}}
	svc := &OcrExecutorService{
		client:        stubClient,
		base64Service: nil,
		config: &config.Config{
			Prompt:            "prompt",
			SystemInstruction: "system",
			Temperature:       0.5,
			MaxTokens:         123,
			Model:             "model",
		},
	}

	result := svc.processImage(ImageResult{FilePath: "sample.png", Base64: encoded})

	if result.Error != "" {
		t.Fatalf("unexpected error: %s", result.Error)
	}

	if result.Content != "ok" {
		t.Fatalf("unexpected content: %s", result.Content)
	}

	if len(stubClient.reqs) == 0 {
		t.Fatal("expected Generate to be called")
	}

	if stubClient.reqs[0].ImageBase64 != encoded {
		t.Fatalf("expected encoded image to be forwarded, got %s", stubClient.reqs[0].ImageBase64)
	}

	if string(stubClient.reqs[0].ImageData) != string(imageData) {
		t.Fatalf("expected decoded data to be forwarded, got %q", stubClient.reqs[0].ImageData)
	}
}

func TestOcrExecutorService_ProcessPath_Integration(t *testing.T) {
	tempDir := t.TempDir()
	imagePath := filepath.Join(tempDir, "sample.png")
	imageBytes := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	if err := os.WriteFile(imagePath, imageBytes, 0o644); err != nil {
		t.Fatalf("failed to create sample image: %v", err)
	}

	stubClient := &stubAIClient{responses: []string{"stub result"}}
	service := &OcrExecutorService{
		client:        stubClient,
		base64Service: NewBase64ExtractorService(tempDir, false),
		config: &config.Config{
			Prompt:            "prompt",
			SystemInstruction: "system",
			Temperature:       0.7,
			MaxTokens:         999,
			Model:             "model",
		},
	}

	result, err := service.ProcessPath()
	if err != nil {
		t.Fatalf("ProcessPath() returned error: %v", err)
	}

	if result.Total != 1 || result.Success != 1 || result.Failed != 0 {
		t.Fatalf("unexpected aggregate result: %+v", result)
	}

	if len(result.Results) != 1 {
		t.Fatalf("unexpected results length: %d", len(result.Results))
	}

	if result.Results[0].Content != "stub result" {
		t.Fatalf("unexpected OCR content: %s", result.Results[0].Content)
	}

	if len(stubClient.reqs) == 0 {
		t.Fatal("expected Generate to be called in integration test")
	}

	if stubClient.reqs[0].MimeType != "image/png" {
		t.Fatalf("unexpected MIME type passed to client: %s", stubClient.reqs[0].MimeType)
	}

	encoded := base64.StdEncoding.EncodeToString(imageBytes)
	if stubClient.reqs[0].ImageBase64 != encoded {
		t.Fatalf("expected base64 payload to be forwarded")
	}

	if string(stubClient.reqs[0].ImageData) != string(imageBytes) {
		t.Fatalf("expected raw image data to be forwarded")
	}
}

func TestOcrExecutorService_ProcessPath_IntegrationWithClientErrors(t *testing.T) {
	tempDir := t.TempDir()
	image1 := filepath.Join(tempDir, "img1.png")
	image2 := filepath.Join(tempDir, "img2.png")
	content := []byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}
	if err := os.WriteFile(image1, content, 0o644); err != nil {
		t.Fatalf("failed to create sample image1: %v", err)
	}
	if err := os.WriteFile(image2, content, 0o644); err != nil {
		t.Fatalf("failed to create sample image2: %v", err)
	}

	stubClient := &stubAIClient{
		responses: []string{"success", ""},
		errs:      []error{nil, fmt.Errorf("upstream error")},
	}

	service := &OcrExecutorService{
		client:        stubClient,
		base64Service: NewBase64ExtractorService(tempDir, false),
		config: &config.Config{
			Prompt:            "prompt",
			SystemInstruction: "system",
			Temperature:       0.1,
			MaxTokens:         256,
			Model:             "model",
		},
	}

	result, err := service.ProcessPath()
	if err != nil {
		t.Fatalf("ProcessPath() returned error: %v", err)
	}

	if result.Total != 2 || result.Success != 1 || result.Failed != 1 {
		t.Fatalf("unexpected aggregate result: %+v", result)
	}

	if stubClient.callCount != 2 {
		t.Fatalf("expected AI client to be called twice, got %d", stubClient.callCount)
	}

	if result.Results[1].Error == "" || !strings.Contains(result.Results[1].Error, "upstream error") {
		t.Fatalf("expected upstream error to be recorded, got %q", result.Results[1].Error)
	}
}

func TestNewOcrExecutorService_UsesFactory(t *testing.T) {
	originalGemini := newGeminiClientFactory
	defer func() { newGeminiClientFactory = originalGemini }()

	called := false
	newGeminiClientFactory = func(cfg *config.Config, client *genai.Client) (AIClient, error) {
		called = true
		if cfg.AiType != "gemini" {
			t.Fatalf("unexpected ai type: %s", cfg.AiType)
		}
		return &stubAIClient{}, nil
	}

	cfg := &config.Config{AiType: "gemini", Path: t.TempDir()}
	service, err := NewOcrExecutorService(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !called {
		t.Fatal("expected gemini factory to be called")
	}

	if service == nil {
		t.Fatal("expected service to be created")
	}
}

func TestNewOcrExecutorService_InvalidAiType(t *testing.T) {
	cfg := &config.Config{AiType: "invalid", Path: t.TempDir()}
	service, err := NewOcrExecutorService(cfg)
	if err == nil {
		t.Fatal("expected error for invalid ai type")
	}
	if service != nil {
		t.Fatal("service should be nil when construction fails")
	}
}

func TestNewOcrExecutorService_FactoryErrorPropagates(t *testing.T) {
	originalVertex := newVertexClientFactory
	defer func() { newVertexClientFactory = originalVertex }()

	expectedErr := fmt.Errorf("factory error")
	newVertexClientFactory = func(cfg *config.Config, client *genai.Client) (AIClient, error) {
		return nil, expectedErr
	}

	cfg := &config.Config{AiType: "vertex", Path: t.TempDir()}
	service, err := NewOcrExecutorService(cfg)
	if err == nil || err.Error() != "AIクライアントの作成に失敗しました: factory error" {
		t.Fatalf("expected wrapped factory error, got %v", err)
	}
	if service != nil {
		t.Fatal("service should be nil when factory fails")
	}
}

func TestOcrExecutorService_CloseDelegates(t *testing.T) {
	closed := false
	client := &closeTrackingClient{closeFn: func() error { closed = true; return nil }}
	svc := &OcrExecutorService{client: client}
	if err := svc.Close(); err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}
	if !closed {
		t.Fatal("expected close to be delegated to client")
	}

	// nil client should be safe
	svc.client = nil
	if err := svc.Close(); err != nil {
		t.Fatalf("unexpected close error when client nil: %v", err)
	}
}

func TestGetMimeType(t *testing.T) {
	svc := &OcrExecutorService{}

	tests := []struct {
		path     string
		expected string
	}{
		{"file.JPG", "image/jpeg"},
		{"file.unknown", ""},
	}

	for _, tc := range tests {
		if got := svc.getMimeType(tc.path); got != tc.expected {
			t.Fatalf("mime mismatch for %s: expected %s got %s", tc.path, tc.expected, got)
		}
	}
}
