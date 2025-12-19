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

// #==============================================================#
// ##       Mocks for HTTPClient                                 ##
// #==============================================================#
// MockHTTPClientForClient はHTTPClientのモック実装（client.go用）
type MockHTTPClientForClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do はHTTPリクエストを実行します
func (m *MockHTTPClientForClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}

// #==============================================================#
// ##       Interfaces for HTTPClient                            ##
// #==============================================================#
// HTTPClient はHTTPリクエストを実行するためのインターフェース
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// #==============================================================#
// ##       Implementations for HTTPClient                       ##
// #==============================================================#
// "net/http"

// #==============================================================#
// ##       Mocks for Client                                     ##
// #==============================================================#
// MockClient はClientのモック実装
type MockClient struct {
	RequestFunc           func(ctx context.Context, method, endpoint string, params map[string]any) (any, error)
	RequestWithoutKeyFunc func(ctx context.Context, method, url string, params map[string]any) (any, error)
	GetHTTPClientFunc     func() HTTPClient
	GetAPIKeyFunc         func() string
	GetHeadersFunc        func() map[string]string
}

// GetHTTPClient はモックのGetHTTPClientメソッド
func (m *MockClient) GetHTTPClient() HTTPClient {
	if m.GetHTTPClientFunc != nil {
		return m.GetHTTPClientFunc()
	}
	return nil
}

// GetAPIKey はモックのGetAPIKeyメソッド
func (m *MockClient) GetAPIKey() string {
	if m.GetAPIKeyFunc != nil {
		return m.GetAPIKeyFunc()
	}
	return ""
}

// GetHeaders はモックのGetHeadersメソッド
func (m *MockClient) GetHeaders() map[string]string {
	if m.GetHeadersFunc != nil {
		return m.GetHeadersFunc()
	}
	return nil
}

// Request はモックのRequestメソッド
func (m *MockClient) Request(ctx context.Context, method, endpoint string, params map[string]any) (any, error) {
	if m.RequestFunc != nil {
		return m.RequestFunc(ctx, method, endpoint, params)
	}
	return nil, nil
}

// RequestWithoutKey はモックのRequestWithoutKeyメソッド
func (m *MockClient) RequestWithoutKey(ctx context.Context, method, url string, params map[string]any) (any, error) {
	if m.RequestWithoutKeyFunc != nil {
		return m.RequestWithoutKeyFunc(ctx, method, url, params)
	}
	return nil, nil
}

// #==============================================================#
// ##       Interfaces for Client                                ##
// #==============================================================#
// ClientInterface はClientを抽象化するインターフェース
type ClientInterface interface {
	Request(ctx context.Context, method, endpoint string, params map[string]any) (any, error)
	RequestWithoutKey(ctx context.Context, method, url string, params map[string]any) (any, error)
	GetHTTPClient() HTTPClient
	GetAPIKey() string
	GetHeaders() map[string]string
}

// #==============================================================#
// ##       Implementations for Client                           ##
// #==============================================================#
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
		map[string]any{"Content-Type": "application/json", "Accept": "application/json"},
		map[string]any{},
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

// GetHTTPClient はHTTPクライアントを返します
func (c *Client) GetHTTPClient() HTTPClient {
	return c.httpClient
}

// GetAPIKey はAPIキーを返します
func (c *Client) GetAPIKey() string {
	return c.apiKey
}

// GetHeaders はヘッダーを返します
func (c *Client) GetHeaders() map[string]string {
	return c.headers
}

// Request はHTTPリクエストを実行し、レスポンスを返します
func (c *Client) Request(ctx context.Context, method, endpoint string, params map[string]any) (any, error) {
	url := buildURLWithParams(APIBaseURL+endpoint, c.GetAPIKey(), params)

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// ヘッダーを設定
	for key, value := range c.GetHeaders() {
		req.Header.Set(key, value)
	}

	// リトライ機能付きでリクエストを実行
	return c.executeWithRetry(req, 3)
}

// RequestWithoutKey はAPIキーなしでHTTPリクエストを実行します（Store APIなど用）
func (c *Client) RequestWithoutKey(ctx context.Context, method, url string, params map[string]any) (any, error) {
	if len(params) > 0 {
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
	for key, value := range c.GetHeaders() {
		req.Header.Set(key, value)
	}

	return c.executeWithRetry(req, 3)
}

// executeWithRetry はリトライ機能付きでHTTPリクエストを実行します
func (c *Client) executeWithRetry(req *http.Request, maxRetries int) (any, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			// リトライ前に少し待機
			time.Sleep(time.Duration(attempt) * time.Second)
			log.Printf("Retrying request to %s, attempt %d/%d", req.URL.String(), attempt, maxRetries)
		}

		resp, err := c.GetHTTPClient().Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		result, err := c.validateResponse(resp)
		if err != nil {
			lastErr = err
			// HTTPエラーまたはSteamAPIエラーの場合はリトライしない
			if resp.StatusCode >= HTTPStatusCode400 && resp.StatusCode < HTTPStatusCode500 {
				return nil, err
			}
			// SteamErrorの場合もリトライしない（直接返す）
			if _, isSteamError := err.(*SteamError); isSteamError {
				return nil, err
			}
			continue
		}

		return result, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

// validateResponse はHTTPレスポンスを検証し、適切な形式で返します
func (c *Client) validateResponse(resp *http.Response) (any, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// HTTPステータスコードをチェック
	if resp.StatusCode >= HTTPStatusCode400 {
		return nil, fmt.Errorf("HTTP error: %d %s - %s", resp.StatusCode, resp.Status, string(body))
	}

	// レスポンスが空の場合
	if len(body) == 0 {
		return "OK", nil
	}

	// JSONとしてパース
	var result any
	if err := json.Unmarshal(body, &result); err != nil {
		// JSONパースに失敗した場合は文字列として返す
		return string(body), nil
	}

	// Steam APIのエラーレスポンスをチェック
	if resultMap, ok := result.(map[string]any); ok {
		if code, exists := resultMap["code"]; exists {
			if description, exists := resultMap["description"]; exists {
				// 安全な型変換
				var codeInt int
				switch v := code.(type) {
				case float64:
					codeInt = int(v)
				case int:
					codeInt = v
				default:
					// 型変換に失敗した場合はデフォルト値を使用
					codeInt = 0
				}

				var descStr string
				if desc, ok := description.(string); ok {
					descStr = desc
				} else {
					descStr = "Unknown error"
				}

				return nil, &SteamError{
					Code:        codeInt,
					Description: descStr,
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

// #==============================================================#
// ##       Mocks for SteamClient                                ##
// #==============================================================#
// MockSteamClient はSteamClientのモック実装
type MockSteamClient struct {
	UsersService *MockUsersService
	AppsService  *MockAppsService
}

// GetUsers はモックのUsersServiceを返します
func (m *MockSteamClient) GetUsers() UsersServiceInterface {
	if m.UsersService != nil {
		return m.UsersService
	}
	return &MockUsersService{}
}

// GetApps はモックのAppsServiceを返します
func (m *MockSteamClient) GetApps() AppsServiceInterface {
	if m.AppsService != nil {
		return m.AppsService
	}
	return &MockAppsService{}
}

// #==============================================================#
// ##       Interfaces for SteamClient                           ##
// #==============================================================#
// SteamClientInterface はSteamクライアント全体を抽象化
type SteamClientInterface interface {
	GetUsers() UsersServiceInterface
	GetApps() AppsServiceInterface
}

// #==============================================================#
// ##       Implementations for SteamClient                      ##
// #==============================================================#
// SteamClient はSteam APIのメインクライアント
type SteamClient struct {
	Users UsersServiceInterface
	Apps  AppsServiceInterface
}

// NewSteamClient は新しいSteam APIクライアントを作成します
func NewSteamClient(apiKey string, headers map[string]string) *SteamClient {
	client := NewClient(apiKey, headers)

	return &SteamClient{
		Users: NewUsersService(client),
		Apps:  NewAppsService(client),
	}
}

// GetUsers はUsersServiceを返します
func (sc *SteamClient) GetUsers() UsersServiceInterface {
	return sc.Users
}

// GetApps はAppsServiceを返します
func (sc *SteamClient) GetApps() AppsServiceInterface {
	return sc.Apps
}
