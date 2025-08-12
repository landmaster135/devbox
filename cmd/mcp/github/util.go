package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	mcp "github.com/mark3labs/mcp-go/mcp"
)

// HTTPClient インターフェースを定義
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// GitHubClient 構造体を修正（テスト用）
type GitHubClient struct {
	httpClient HTTPClient
	token      string
}

// NewGitHubClient は新しいGitHubクライアントを作成します
func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{},
		token:      token,
	}
}

// GitHubError はGitHub APIからのエラーを表します
type GitHubError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
	StatusCode       int
}

func (e *GitHubError) Error() string {
	return fmt.Sprintf("GitHub API Error: %s (Status: %d)", e.Message, e.StatusCode)
}

// ヘルパー関数: 結果をJSON形式で返却
func returnJSONResult(result interface{}) (*mcp.CallToolResult, error) {
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(jsonResult)), nil
}
