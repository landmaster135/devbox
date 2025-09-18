package steam_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient はHTTPリクエストを実行するためのインターフェース
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client はSteam Web APIのHTTPクライアント
type Client struct {
	httpClient HTTPClient
	apiKey     string
	headers    map[string]string
}

// NewClient は新しいHTTPクライアントを作成します
func NewClient(apiKey string, headers map[string]string) *Client {
	if headers == nil {
		headers = make(map[string]string)
	}

	// デフォルトヘッダーを設定
	defaultHeaders := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// ユーザー指定のヘッダーとマージ
	mergedHeaders := mergeParams(
		map[string]interface{}{"Content-Type": "application/json", "Accept": "application/json"},
		map[string]interface{}{},
	)
	for k, v := range headers {
		mergedHeaders[k] = v
	}

	finalHeaders := make(map[string]string)
	for k, v := range defaultHeaders {
		finalHeaders[k] = v
	}
	for k, v := range headers {
		finalHeaders[k] = v
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:  apiKey,
		headers: finalHeaders,
	}
}

// Request はHTTPリクエストを実行し、レスポンスを返します
func (c *Client) Request(ctx context.Context, method, endpoint string, params map[string]interface{}) (interface{}, error) {
	url := buildURLWithParams(APIBaseURL+endpoint, c.apiKey, params)

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// ヘッダーを設定
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// リトライ機能付きでリクエストを実行
	return c.executeWithRetry(req, 3)
}

// RequestWithoutKey はAPIキーなしでHTTPリクエストを実行します（Store APIなど用）
func (c *Client) RequestWithoutKey(ctx context.Context, method, url string, params map[string]interface{}) (interface{}, error) {
	if params != nil && len(params) > 0 {
		cleanedParams := cleanParams(params)
		if len(cleanedParams) > 0 {
			u, err := buildURLWithQueryParams(url, cleanedParams)
			if err != nil {
				return nil, fmt.Errorf("failed to build URL: %w", err)
			}
			url = u
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// ヘッダーを設定
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	return c.executeWithRetry(req, 3)
}

// executeWithRetry はリトライ機能付きでHTTPリクエストを実行します
func (c *Client) executeWithRetry(req *http.Request, maxRetries int) (interface{}, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// リトライ前に少し待機
			time.Sleep(time.Duration(attempt) * time.Second)
			log.Printf("Retrying request to %s, attempt %d/%d", req.URL.String(), attempt, maxRetries)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		result, err := c.validateResponse(resp)
		if err != nil {
			lastErr = err
			// HTTPエラーの場合はリトライしない
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				break
			}
			continue
		}

		return result, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

// validateResponse はHTTPレスポンスを検証し、適切な形式で返します
func (c *Client) validateResponse(resp *http.Response) (interface{}, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// HTTPステータスコードをチェック
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error: %d %s - %s", resp.StatusCode, resp.Status, string(body))
	}

	// レスポンスが空の場合
	if len(body) == 0 {
		return "OK", nil
	}

	// JSONとしてパース
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// JSONパースに失敗した場合は文字列として返す
		return string(body), nil
	}

	// Steam APIのエラーレスポンスをチェック
	if resultMap, ok := result.(map[string]interface{}); ok {
		if code, exists := resultMap["code"]; exists {
			if description, exists := resultMap["description"]; exists {
				return nil, &SteamError{
					Code:        int(code.(float64)),
					Description: description.(string),
					StatusCode:  resp.StatusCode,
				}
			}
		}
	}

	return result, nil
}

// buildURLWithQueryParams はクエリパラメータ付きのURLを構築します
func buildURLWithQueryParams(baseURL string, params map[string]string) (string, error) {
	u, err := parseURL(baseURL)
	if err != nil {
		return "", err
	}

	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}

	u.RawQuery = query.Encode()
	return u.String(), nil
}

// parseURL はURL文字列をパースします
func parseURL(rawURL string) (*url.URL, error) {
	return url.Parse(rawURL)
}

// SteamClient はSteam APIのメインクライアント
type SteamClient struct {
	client *Client
	Users  *UsersService
	Apps   *AppsService
}

// NewSteamClient は新しいSteam APIクライアントを作成します
func NewSteamClient(apiKey string, headers map[string]string) *SteamClient {
	client := NewClient(apiKey, headers)

	return &SteamClient{
		client: client,
		Users:  NewUsersService(client),
		Apps:   NewAppsService(client),
	}
}

// GetClient はHTTPクライアントを返します
func (sc *SteamClient) GetClient() *Client {
	return sc.client
}
