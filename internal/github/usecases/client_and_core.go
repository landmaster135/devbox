package usecases

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiBaseURL = "https://api.github.com"
	version    = "1.0.0"
)

// #==============================================================#
// ##          Interfaces                                        ##
// #==============================================================#

// HTTPClient インターフェースを定義
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// JSONMarshaler インターフェースを定義
type JSONMarshaler interface {
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
}

// DefaultJSONMarshaler は標準のjson.MarshalIndentを使用する実装
type DefaultJSONMarshaler struct{}

func (m *DefaultJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// #==============================================================#
// ##          GitHubError                                       ##
// #==============================================================#

// GitHubError はGitHub APIからのエラーを表します
type GitHubError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	StatusCode       int
}

func (e *GitHubError) Error() string {
	return fmt.Sprintf("GitHub API Error: %s (Status: %d)", e.Message, e.StatusCode)
}

// #==============================================================#
// ##          GitHubClientService                               ##
// #==============================================================#

// GitHubClientService はGitHub API操作のビジネスロジックを提供します
type GitHubClientService struct {
	httpClient    HTTPClient
	token         string
	jsonMarshaler JSONMarshaler
}

// NewGitHubClientService は新しいGitHubClientServiceを作成します
func NewGitHubClientService(token string) *GitHubClientService {
	return &GitHubClientService{
		httpClient:    &http.Client{},
		token:         token,
		jsonMarshaler: &DefaultJSONMarshaler{},
	}
}

// NewGitHubClientServiceWithDependencies はテスト用に依存性を注入できるGitHubClientServiceを作成します
func NewGitHubClientServiceWithDependencies(httpClient HTTPClient, token string, jsonMarshaler JSONMarshaler) *GitHubClientService {
	return &GitHubClientService{
		httpClient:    httpClient,
		token:         token,
		jsonMarshaler: jsonMarshaler,
	}
}

// DoRequest はHTTPリクエストを実行し、レスポンスを処理します
func (c *GitHubClientService) DoRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	if method == "POST" || method == "PATCH" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var ghError GitHubError
		if err := json.Unmarshal(respBody, &ghError); err != nil {
			return nil, fmt.Errorf("HTTP error: %d - %s", resp.StatusCode, string(respBody))
		}
		ghError.StatusCode = resp.StatusCode
		return nil, &ghError
	}

	return respBody, nil
}

// ReturnJSONResult は結果をJSON形式で返却するヘルパー関数です
func (c *GitHubClientService) ReturnJSONResult(result interface{}) (string, error) {
	jsonResult, err := c.jsonMarshaler.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonResult), nil
}
