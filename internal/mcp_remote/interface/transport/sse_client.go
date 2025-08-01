package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SSEClient はServer-Sent Eventsを使用してリモートMCPサーバーと通信する
type SSEClient struct {
	serverURL *url.URL
	headers   map[string]string
	logger    *log.Logger
	client    *http.Client
	conn      *http.Response
}

// NewSSEClient は新しいSSEClientを作成する
func NewSSEClient(serverURL *url.URL, headers map[string]string, logger *log.Logger) (*SSEClient, error) {
	return &SSEClient{
		serverURL: serverURL,
		headers:   headers,
		logger:    logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Connect はSSE接続を確立する
func (c *SSEClient) Connect(ctx context.Context) error {
	c.logger.Printf("SSE接続を確立中: %s", c.serverURL.String())

	req, err := http.NewRequestWithContext(ctx, "GET", c.serverURL.String(), nil)
	if err != nil {
		return fmt.Errorf("リクエストの作成に失敗しました: %v", err)
	}

	// ヘッダーを設定
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("SSE接続に失敗しました: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("SSE接続が失敗しました。ステータス: %d", resp.StatusCode)
	}

	c.conn = resp
	c.logger.Printf("SSE接続が確立されました")
	return nil
}

// SendMessage はメッセージをリモートサーバーに送信する
func (c *SSEClient) SendMessage(ctx context.Context, message json.RawMessage) (json.RawMessage, error) {
	c.logger.Printf("SSE経由でメッセージを送信中: %s", string(message))

	// SSEは通常一方向なので、POSTリクエストでメッセージを送信
	postURL := c.serverURL.String()
	if !strings.HasSuffix(postURL, "/") {
		postURL += "/"
	}
	postURL += "message"

	req, err := http.NewRequestWithContext(ctx, "POST", postURL, strings.NewReader(string(message)))
	if err != nil {
		return nil, fmt.Errorf("POSTリクエストの作成に失敗しました: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("メッセージの送信に失敗しました: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("メッセージ送信が失敗しました。ステータス: %d", resp.StatusCode)
	}

	// レスポンスを読み取り
	scanner := bufio.NewScanner(resp.Body)
	var responseData []byte
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data != "" && data != "[DONE]" {
				responseData = append(responseData, []byte(data)...)
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("レスポンスの読み取りに失敗しました: %v", err)
	}

	c.logger.Printf("SSE経由でレスポンスを受信: %s", string(responseData))
	return json.RawMessage(responseData), nil
}

// Close はSSE接続を閉じる
func (c *SSEClient) Close() error {
	c.logger.Printf("SSE接続を閉じています")
	if c.conn != nil {
		return c.conn.Body.Close()
	}
	return nil
}
