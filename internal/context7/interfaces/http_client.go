package interfaces

import "net/http"

// HTTPClient はHTTPリクエストを実行するためのインターフェースです
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// DefaultHTTPClient は標準のHTTPクライアントを提供します
type DefaultHTTPClient struct {
	client *http.Client
}

// NewDefaultHTTPClient は新しいDefaultHTTPClientインスタンスを作成します
func NewDefaultHTTPClient() *DefaultHTTPClient {
	return &DefaultHTTPClient{
		client: &http.Client{},
	}
}

// Do はHTTPリクエストを実行します
func (c *DefaultHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
