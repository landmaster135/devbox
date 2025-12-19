package usecases

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"google.golang.org/genai"

	config "github.com/landmaster135/devbox/internal/ocr_executor_with_ai/config"
)

// --- Constructors -----------------------------------------------------------

func TestNewGeminiAIClient_UsesInjectedClient(t *testing.T) {
	cfg := &config.Config{}
	injected := &genai.Client{Models: &genai.Models{}}

	client, err := NewGeminiAIClient(cfg, injected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gClient, ok := client.(*GeminiAIClient)
	if !ok {
		t.Fatalf("expected *GeminiAIClient, got %T", client)
	}
	if gClient.client != injected {
		t.Fatal("injected client should be preserved")
	}
}

func TestNewGeminiAIClient_MissingAPIKey_ReturnsError(t *testing.T) {
	cfg := &config.Config{}
	client, err := NewGeminiAIClient(cfg, nil)
	if err == nil {
		t.Fatal("expected error when API key is missing")
	}
	if client != nil {
		t.Fatal("client should be nil on error")
	}
}

func TestNewGeminiAIClient_MissingModels_ReturnsError(t *testing.T) {
	cfg := &config.Config{}
	injected := &genai.Client{}
	client, err := NewGeminiAIClient(cfg, injected)
	if err == nil || !strings.Contains(err.Error(), "モデルが初期化されていません") {
		t.Fatalf("expected models initialization error, got %v", err)
	}
	if client != nil {
		t.Fatal("client should be nil when models are missing")
	}
}

func TestNewVertexAIClient_UsesInjectedClient(t *testing.T) {
	cfg := &config.Config{}
	injected := &genai.Client{Models: &genai.Models{}}
	client, err := NewVertexAIClient(cfg, injected)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vClient, ok := client.(*VertexAIClient)
	if !ok {
		t.Fatalf("expected *VertexAIClient, got %T", client)
	}
	if vClient.client != injected {
		t.Fatal("injected client should be preserved")
	}
}

func TestNewVertexAIClient_MissingProject_ReturnsError(t *testing.T) {
	cfg := &config.Config{}
	client, err := NewVertexAIClient(cfg, nil)
	if err == nil {
		t.Fatal("expected project/location error")
	}
	if client != nil {
		t.Fatal("client should be nil on error")
	}
}

func TestNewOllamaAIClient_DefaultClient(t *testing.T) {
	cfg := &config.Config{}
	client, err := NewOllamaAIClient(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	oClient, ok := client.(*OllamaAIClient)
	if !ok {
		t.Fatalf("expected *OllamaAIClient, got %T", client)
	}
	if oClient.httpClient == nil {
		t.Fatal("expected default http client")
	}
	if oClient.httpClient.Timeout != 120*time.Second {
		t.Fatalf("unexpected timeout: %s", oClient.httpClient.Timeout)
	}
}

// --- genAIBaseClient --------------------------------------------------------

type stubContentGenerator struct {
	resp      *genai.GenerateContentResponse
	err       error
	seenModel string
}

func (s *stubContentGenerator) GenerateContent(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	s.seenModel = model
	return s.resp, s.err
}

func TestGenAIBaseGenerate_Success(t *testing.T) {
	stub := &stubContentGenerator{resp: &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{{Content: &genai.Content{Parts: []*genai.Part{{Text: "Hello"}, {Text: " World"}}}}},
	}}
	base := genAIBaseClient{generator: stub}

	result, err := base.generate(context.Background(), &AIRequest{
		Prompt:            "prompt",
		SystemInstruction: "system",
		Temperature:       0.5,
		MaxTokens:         123,
		Model:             "model",
		MimeType:          "image/png",
		ImageData:         []byte{1, 2, 3},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello World" {
		t.Fatalf("unexpected result: %s", result)
	}
	if stub.seenModel != "model" {
		t.Fatalf("model not forwarded, got %s", stub.seenModel)
	}
}

func TestGenAIBaseGenerate_ErrorPropagation(t *testing.T) {
	stub := &stubContentGenerator{err: errors.New("boom")}
	base := genAIBaseClient{generator: stub}
	_, err := base.generate(context.Background(), &AIRequest{MimeType: "image/png", ImageData: []byte{1}})
	if err == nil || !strings.Contains(err.Error(), "コンテンツ生成エラー") {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestGenAIBaseGenerate_NoCandidates(t *testing.T) {
	stub := &stubContentGenerator{resp: &genai.GenerateContentResponse{}}
	base := genAIBaseClient{generator: stub}
	_, err := base.generate(context.Background(), &AIRequest{MimeType: "image/png", ImageData: []byte{1}})
	if err == nil || !strings.Contains(err.Error(), "候補") {
		t.Fatalf("expected candidate error, got %v", err)
	}
}

func TestGenAIBaseGenerate_NoContent(t *testing.T) {
	stub := &stubContentGenerator{resp: &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{{Content: &genai.Content{Parts: []*genai.Part{}}}},
	}}
	base := genAIBaseClient{generator: stub}
	_, err := base.generate(context.Background(), &AIRequest{MimeType: "image/png", ImageData: []byte{1}})
	if err == nil || !strings.Contains(err.Error(), "コンテンツ") {
		t.Fatalf("expected content error, got %v", err)
	}
}

func TestGenAIBaseGenerate_MultipleCandidates(t *testing.T) {
	stub := &stubContentGenerator{resp: &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{Content: &genai.Content{Parts: []*genai.Part{{Text: "Primary"}}}},
			{Content: &genai.Content{Parts: []*genai.Part{{Text: "Secondary"}}}},
		},
	}}
	base := genAIBaseClient{generator: stub}
	res, err := base.generate(context.Background(), &AIRequest{MimeType: "image/png", ImageData: []byte{1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "Primary" {
		t.Fatalf("expected first candidate text, got %s", res)
	}
}

func TestGenAIBaseGenerate_NilContentCandidate(t *testing.T) {
	stub := &stubContentGenerator{resp: &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{{Content: nil}},
	}}
	base := genAIBaseClient{generator: stub}
	_, err := base.generate(context.Background(), &AIRequest{MimeType: "image/png", ImageData: []byte{1}})
	if err == nil || !strings.Contains(err.Error(), "コンテンツ") {
		t.Fatalf("expected nil content error, got %v", err)
	}
}

func TestGenAIBaseGenerate_NilGenerator(t *testing.T) {
	base := genAIBaseClient{}
	_, err := base.generate(context.Background(), &AIRequest{MimeType: "image/png", ImageData: []byte{1}})
	if err == nil || !strings.Contains(err.Error(), "初期化") {
		t.Fatalf("expected initialization error, got %v", err)
	}
}

func TestGenAIBaseGenerate_EmptyImageData(t *testing.T) {
	base := genAIBaseClient{generator: &stubContentGenerator{}}
	_, err := base.generate(context.Background(), &AIRequest{})
	if err == nil || !strings.Contains(err.Error(), "画像データが空です") {
		t.Fatalf("expected image data error, got %v", err)
	}
}

// --- Gemini / Vertex Clients ------------------------------------------------

type recordingGenerator struct {
	called bool
}

func (r *recordingGenerator) GenerateContent(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	r.called = true
	return &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{{Content: &genai.Content{Parts: []*genai.Part{{Text: "text"}}}}},
	}, nil
}

func TestGeminiAIClientGenerate_Delegates(t *testing.T) {
	gen := &recordingGenerator{}
	client := &GeminiAIClient{base: genAIBaseClient{generator: gen}}
	result, err := client.Generate(context.Background(), &AIRequest{MimeType: "image/png", ImageData: []byte{1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "text" {
		t.Fatalf("unexpected result: %s", result)
	}
	if !gen.called {
		t.Fatal("expected generator to be called")
	}
}

func TestVertexAIClientGenerate_PropagatesError(t *testing.T) {
	gen := &failingGenerator{err: errors.New("fail")}
	client := &VertexAIClient{base: genAIBaseClient{generator: gen}}
	_, err := client.Generate(context.Background(), &AIRequest{MimeType: "image/png", ImageData: []byte{1}})
	if err == nil || !strings.Contains(err.Error(), "fail") {
		t.Fatalf("expected propagated error, got %v", err)
	}
}

type failingGenerator struct {
	err error
}

func (f *failingGenerator) GenerateContent(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	return nil, f.err
}

func TestGeminiAndVertexClose(t *testing.T) {
	if err := (&GeminiAIClient{}).Close(); err != nil {
		t.Fatalf("Gemini close returned error: %v", err)
	}
	if err := (&VertexAIClient{}).Close(); err != nil {
		t.Fatalf("Vertex close returned error: %v", err)
	}
	if err := (&OllamaAIClient{}).Close(); err != nil {
		t.Fatalf("Ollama close returned error: %v", err)
	}
}

// --- Ollama client ----------------------------------------------------------

func TestOllamaAIClient_Generate_Success(t *testing.T) {
	trip := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", req.Method)
		}
		body := strings.NewReader(`{"response":"Hello","done":false}` + "\n" + `{"response":" World","done":true}`)
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(body),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewOllamaAIClient(&config.Config{}, &http.Client{Transport: trip})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := client.Generate(context.Background(), &AIRequest{Prompt: "hello", ImageBase64: "abc"})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if result != "Hello World" {
		t.Fatalf("unexpected concatenated result: %s", result)
	}
}

func TestOllamaAIClient_Generate_HTTPError(t *testing.T) {
	trip := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("bad")),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewOllamaAIClient(&config.Config{}, &http.Client{Transport: trip})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.Generate(context.Background(), &AIRequest{Prompt: "hello", ImageBase64: "abc"})
	if err == nil || !strings.Contains(err.Error(), "ollama APIがエラーを返しました") {
		t.Fatalf("expected HTTP error propagation, got %v", err)
	}
}

func TestOllamaAIClient_Generate_InvalidChunks(t *testing.T) {
	trip := roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"response":"broken"`)),
			Header:     make(http.Header),
		}, nil
	})

	client, err := NewOllamaAIClient(&config.Config{}, &http.Client{Transport: trip})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.Generate(context.Background(), &AIRequest{Prompt: "hello", ImageBase64: "abc"})
	if err == nil || !strings.Contains(err.Error(), "ollamaレスポンスの解析に失敗しました") {
		t.Fatalf("expected JSON parse error, got %v", err)
	}
}

func TestOllamaAIClient_Generate_NoHTTPClient(t *testing.T) {
	client := &OllamaAIClient{}
	_, err := client.Generate(context.Background(), &AIRequest{Prompt: ""})
	if err == nil || !strings.Contains(err.Error(), "初期化されていません") {
		t.Fatalf("expected initialization error, got %v", err)
	}
}

// --- Helpers ----------------------------------------------------------------

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
