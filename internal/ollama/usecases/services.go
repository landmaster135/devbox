package usecases

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/ollama/domain"
)

const (
	defaultBaseURL      = "http://127.0.0.1:11434"
	scannerMaxTokenSize = 2 * 1024 * 1024 // 2MB
)

// HTTPClient は http.Client と互換性のあるインターフェース。
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ServiceOptions は Service の組み立て時に渡すオプション。
type ServiceOptions struct {
	BaseURL    string
	Timeout    time.Duration
	HTTPClient HTTPClient
}

// Service は Ollama API を呼び出すユースケース層。
type Service struct {
	baseURL string
	client  HTTPClient
}

// NewService は Service を生成する。
func NewService(opts ServiceOptions) *Service {
	baseURL := strings.TrimRight(strings.TrimSpace(opts.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	client := opts.HTTPClient
	if client == nil {
		timeout := opts.Timeout
		if timeout <= 0 {
			timeout = 30 * time.Second
		}
		client = &http.Client{Timeout: timeout}
	}

	return &Service{
		baseURL: baseURL,
		client:  client,
	}
}

// GetVersion は /api/version を呼び出す。
func (s *Service) GetVersion(ctx context.Context) (*domain.VersionResponse, error) {
	var resp domain.VersionResponse
	if err := s.doJSON(ctx, http.MethodGet, "/api/version", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListInstalledModels は /api/tags を呼び出す。
func (s *Service) ListInstalledModels(ctx context.Context) (*domain.TagsResponse, error) {
	var resp domain.TagsResponse
	if err := s.doJSON(ctx, http.MethodGet, "/api/tags", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ListRunningModels は /api/ps を呼び出す。
func (s *Service) ListRunningModels(ctx context.Context) (*domain.ProcessesResponse, error) {
	var resp domain.ProcessesResponse
	if err := s.doJSON(ctx, http.MethodGet, "/api/ps", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateEmbeddings は /api/embed を呼び出す。
func (s *Service) CreateEmbeddings(ctx context.Context, req domain.EmbedRequest) (*domain.EmbedResponse, error) {
	var resp domain.EmbedResponse
	if err := s.doJSON(ctx, http.MethodPost, "/api/embed", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Generate は /api/generate をストリーミングで処理する。
func (s *Service) Generate(ctx context.Context, req domain.GenerateRequest) (string, error) {
	if !req.Stream {
		req.Stream = true
	}
	resp, err := s.stream(ctx, "/api/generate", req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, scannerMaxTokenSize)

	var builder strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var chunk domain.GenerateChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			return "", fmt.Errorf("generate レスポンスの解析に失敗しました: %w", err)
		}
		if chunk.Error != "" {
			return "", fmt.Errorf("Ollama API エラー: %s", chunk.Error)
		}
		builder.WriteString(chunk.Response)
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("generate レスポンスの読み取りに失敗しました: %w", err)
	}
	return builder.String(), nil
}

// PullModel は /api/pull のストリーミングレスポンスを処理する。
func (s *Service) PullModel(ctx context.Context, req domain.PullRequest) (string, error) {
	if !req.Stream {
		req.Stream = true
	}
	resp, err := s.stream(ctx, "/api/pull", req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, scannerMaxTokenSize)

	var builder strings.Builder
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var chunk domain.PullChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			return "", fmt.Errorf("pull レスポンスの解析に失敗しました: %w", err)
		}
		if chunk.Error != "" {
			return "", fmt.Errorf("Ollama API エラー: %s", chunk.Error)
		}
		if chunk.Status != "" {
			builder.WriteString(formatPullStatus(chunk))
			builder.WriteByte('\n')
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("pull レスポンスの読み取りに失敗しました: %w", err)
	}
	return strings.TrimRight(builder.String(), "\n"), nil
}

func (s *Service) doJSON(ctx context.Context, method, path string, payload any, out any) error {
	req, err := s.newJSONRequest(ctx, method, path, payload)
	if err != nil {
		return err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("Ollama API %s %s の呼び出しに失敗しました: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Ollama API %s %s が失敗しました: status=%d body=%s", method, path, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if out == nil {
		io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("Ollama API %s %s レスポンスのデコードに失敗しました: %w", method, path, err)
	}
	return nil
}

func (s *Service) stream(ctx context.Context, path string, payload any) (*http.Response, error) {
	req, err := s.newJSONRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Ollama API POST %s の呼び出しに失敗しました: %w", path, err)
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("Ollama API POST %s が失敗しました: status=%d body=%s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return resp, nil
}

func (s *Service) newJSONRequest(ctx context.Context, method, path string, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return nil, fmt.Errorf("リクエストボディのエンコードに失敗しました: %w", err)
		}
		body = buf
	}

	requestURL := s.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("HTTP リクエストの組み立てに失敗しました: %w", err)
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func formatPullStatus(chunk domain.PullChunk) string {
	switch {
	case chunk.Total > 0 && chunk.Completed > 0:
		percent := float64(chunk.Completed) / float64(chunk.Total) * 100
		return fmt.Sprintf("%s %.1f%% (%d/%d)", chunk.Status, percent, chunk.Completed, chunk.Total)
	case chunk.Total > 0:
		return fmt.Sprintf("%s 0.0%% (0/%d)", chunk.Status, chunk.Total)
	default:
		return chunk.Status
	}
}
