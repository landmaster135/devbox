package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// HTTPClient は http.Client と互換なインターフェース。
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// JSONClientOptions は JSONClient 生成時の入力。
type JSONClientOptions struct {
	BaseURL    string
	APIToken   string
	HTTPClient HTTPClient
}

// JSONClient は Memos API への JSON リクエストを扱う。
type JSONClient struct {
	baseURL  string
	apiToken string
	client   HTTPClient
}

func NewJSONClient(opts JSONClientOptions) *JSONClient {
	return &JSONClient{
		baseURL:  NormalizeBaseURL(opts.BaseURL),
		apiToken: strings.TrimSpace(opts.APIToken),
		client:   opts.HTTPClient,
	}
}

func (c *JSONClient) DoJSON(ctx context.Context, method, requestPath string, query url.Values, payload any, out any) error {
	req, err := c.NewRequest(ctx, method, requestPath, query, payload)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("Memos API %s %s の呼び出しに失敗しました: %w", method, requestPath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Memos API %s %s が失敗しました: status=%d body=%s", method, requestPath, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if out == nil {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("Memos API %s %s のレスポンスデコードに失敗しました: %w", method, requestPath, err)
	}
	return nil
}

func (c *JSONClient) NewRequest(ctx context.Context, method, requestPath string, query url.Values, payload any) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		buf := bytes.NewBuffer(nil)
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return nil, fmt.Errorf("リクエストボディのエンコードに失敗しました: %w", err)
		}
		body = buf
	}

	requestURL := c.baseURL + requestPath
	if len(query) > 0 {
		requestURL = requestURL + "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, fmt.Errorf("HTTP リクエストの作成に失敗しました: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiToken)
	}
	return req, nil
}
