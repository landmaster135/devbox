package usecases

import (
	"encoding/json"
	"errors"
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

// #==============================================================#
// ##          GitHubSearchService                               ##
// #==============================================================#

func validateParams(token string) error {
	if token == "" {
		return errors.New("token must not be empty")
	}
	return nil
}

// GitHubSearchService はGitHub検索機能のビジネスロジックを提供します
type GitHubSearchService struct {
	clientService *GitHubClientService
}

// NewGitHubSearchService は新しいGitHubSearchServiceを作成します
func NewGitHubSearchService(token string) (*GitHubSearchService, error) {
	err := validateParams(token)
	if err != nil {
		return nil, err
	}
	return &GitHubSearchService{
		clientService: NewGitHubClientService(token),
	}, nil
}

// NewGitHubSearchServiceWithDependencies はテスト用に依存性を注入できるGitHubSearchServiceを作成します
func NewGitHubSearchServiceWithDependencies(clientService *GitHubClientService) *GitHubSearchService {
	return &GitHubSearchService{
		clientService: clientService,
	}
}

// SearchCode はGitHub全体でコードを検索します
func (s *GitHubSearchService) SearchCode(query string, options map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/search/code?q=%s", apiBaseURL, query)

	// クエリパラメータを追加
	for k, v := range options {
		url += fmt.Sprintf("&%s=%v", k, v)
	}

	data, err := s.clientService.DoRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// HandleToSearchCode はGitHub全体でコードを検索して、結果をJSON形式で返します
func (s *GitHubSearchService) HandleToSearchCode(query string, page, perPage int) (string, error) {
	options := make(map[string]interface{})

	// 数値オプションパラメータを追加
	if page != 1 {
		options["page"] = page
	}

	if perPage != 30 {
		options["per_page"] = perPage
	}

	result, err := s.SearchCode(query, options)
	if err != nil {
		return "", err
	}

	return s.clientService.ReturnJSONResult(result)
}
