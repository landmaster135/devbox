package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient は通常のHTTPリクエストを使用してリモートMCPサーバーと通信する
type HTTPClient struct {
	serverURL *url.URL
	headers   map[string]string
	logger    *log.Logger
	client    *http.Client
}

// NewHTTPClient は新しいHTTPClientを作成する
func NewHTTPClient(serverURL *url.URL, headers map[string]string, logger *log.Logger) (*HTTPClient, error) {
	return &HTTPClient{
		serverURL: serverURL,
		headers:   headers,
		logger:    logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Connect はHTTP接続を確認する（実際の接続は各リクエストで行う）
func (c *HTTPClient) Connect(ctx context.Context) error {
	c.logger.Printf("HTTP接続を確認中: %s", c.serverURL.String())

	// 接続テスト用のリクエストを送信
	req, err := http.NewRequestWithContext(ctx, "GET", c.serverURL.String(), nil)
	if err != nil {
		return fmt.Errorf("テストリクエストの作成に失敗しました: %v", err)
	}

	// ヘッダーを設定
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP接続テストに失敗しました: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP接続テストが失敗しました。ステータス: %d", resp.StatusCode)
	}

	c.logger.Printf("HTTP接続が確認されました")
	return nil
}

// SendMessage はメッセージをリモートサーバーに送信する
func (c *HTTPClient) SendMessage(ctx context.Context, message json.RawMessage) (json.RawMessage, error) {
	c.logger.Printf("HTTP経由でメッセージを送信中: %s", string(message))

	req, err := http.NewRequestWithContext(ctx, "POST", c.serverURL.String(), bytes.NewReader(message))
	if err != nil {
		return nil, fmt.Errorf("HTTPリクエストの作成に失敗しました: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTPメッセージの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTPメッセージ送信が失敗しました。ステータス: %d, レスポンス: %s", resp.StatusCode, string(body))
	}

	// レスポンスを読み取り
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTPレスポンスの読み取りに失敗しました: %v", err)
	}

	c.logger.Printf("HTTP経由でレスポンスを受信: %s", string(responseData))
	return json.RawMessage(responseData), nil
}

// Close はHTTPクライアントをクリーンアップする
func (c *HTTPClient) Close() error {
	c.logger.Printf("HTTPクライアントをクリーンアップしています")
	// HTTPクライアントは特別なクリーンアップは不要
	return nil
}
