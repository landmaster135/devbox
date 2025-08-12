package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	
	config "github.com/landmaster135/devbox/internal/mcp_remote/config"
)

// ProxyTransport はSTDIOとリモートサーバー間の双方向プロキシを提供する
type ProxyTransport struct {
	config        *config.Config
	logger        *log.Logger
	remoteClient  RemoteClient
	localServer   *server.MCPServer
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// RemoteClient はリモートMCPサーバーとの通信を抽象化するインターフェース
type RemoteClient interface {
	Connect(ctx context.Context) error
	SendMessage(ctx context.Context, message json.RawMessage) (json.RawMessage, error)
	Close() error
}

// NewProxyTransport は新しいProxyTransportを作成する
func NewProxyTransport(cfg *config.Config) *ProxyTransport {
	ctx, cancel := context.WithCancel(context.Background())

	return &ProxyTransport{
		config: cfg,
		logger: log.New(os.Stderr, "[proxy] ", log.LstdFlags),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start はプロキシを開始する
func (p *ProxyTransport) Start() error {
	p.logger.Printf("プロキシを開始します: %s", p.config.ServerURL)

	// リモートクライアントを作成
	remoteClient, err := p.createRemoteClient()
	if err != nil {
		return fmt.Errorf("リモートクライアントの作成に失敗しました: %v", err)
	}
	p.remoteClient = remoteClient

	// リモートサーバーに接続
	if err := p.remoteClient.Connect(p.ctx); err != nil {
		return fmt.Errorf("リモートサーバーへの接続に失敗しました: %v", err)
	}

	// ローカルMCPサーバーを作成（プロキシとして動作）
	p.localServer = p.createLocalProxyServer()

	// STDIOサーバーを開始
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		if err := server.ServeStdio(p.localServer); err != nil {
			p.logger.Printf("STDIOサーバーエラー: %v", err)
		}
	}()

	p.logger.Printf("プロキシが正常に開始されました")
	return nil
}

// Stop はプロキシを停止する
func (p *ProxyTransport) Stop() error {
	p.logger.Printf("プロキシを停止中...")

	p.cancel()

	if p.remoteClient != nil {
		if err := p.remoteClient.Close(); err != nil {
			p.logger.Printf("リモートクライアントの終了エラー: %v", err)
		}
	}

	p.wg.Wait()
	p.logger.Printf("プロキシが正常に停止されました")
	return nil
}

// createRemoteClient はトランスポート戦略に基づいてリモートクライアントを作成する
func (p *ProxyTransport) createRemoteClient() (RemoteClient, error) {
	serverURL, err := url.Parse(p.config.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("無効なサーバーURL: %v", err)
	}

	switch p.config.TransportStrategy {
	case config.TransportSSEOnly, config.TransportSSEFirst:
		p.logger.Printf("SSEクライアントを作成中...")
		return NewSSEClient(serverURL, p.config.Headers, p.logger)
	case config.TransportHTTPOnly, config.TransportHTTPFirst:
		p.logger.Printf("HTTPクライアントを作成中...")
		return NewHTTPClient(serverURL, p.config.Headers, p.logger)
	default:
		return nil, fmt.Errorf("未対応のトランスポート戦略: %s", p.config.TransportStrategy)
	}
}

// createLocalProxyServer はローカルプロキシサーバーを作成する
func (p *ProxyTransport) createLocalProxyServer() *server.MCPServer {
	s := server.NewMCPServer(
		"MCP Remote Proxy",
		"1.0.0",
		server.WithLogging(),
	)

	// プロキシツールを追加（リモートサーバーのツールを中継）
	p.addProxyTools(s)

	return s
}

// addProxyTools はプロキシツールを追加する
func (p *ProxyTransport) addProxyTools(s *server.MCPServer) {
	// リモートサーバーからツール一覧を取得して動的に追加
	p.logger.Printf("リモートサーバーからツール一覧を取得中...")

	// tools/listリクエストを作成
	toolsListRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "tools-list-1",
		"method":  "tools/list",
		"params":  map[string]interface{}{},
	}

	requestJSON, err := json.Marshal(toolsListRequest)
	if err != nil {
		p.logger.Printf("tools/listリクエストの作成に失敗しました: %v", err)
		return
	}

	// リモートサーバーにツール一覧を要求
	responseJSON, err := p.remoteClient.SendMessage(p.ctx, requestJSON)
	if err != nil {
		p.logger.Printf("リモートサーバーからのツール一覧取得に失敗しました: %v", err)
		return
	}

	// レスポンスを解析
	var response struct {
		Result struct {
			Tools []struct {
				Name        string                 `json:"name"`
				Description string                 `json:"description"`
				InputSchema map[string]interface{} `json:"inputSchema"`
			} `json:"tools"`
		} `json:"result"`
	}

	if err := json.Unmarshal(responseJSON, &response); err != nil {
		p.logger.Printf("ツール一覧レスポンスの解析に失敗しました: %v", err)
		return
	}

	p.logger.Printf("リモートサーバーから%d個のツールを取得しました", len(response.Result.Tools))

	// 各ツールをプロキシとして追加
	for _, tool := range response.Result.Tools {
		p.logger.Printf("ツールを追加中: %s", tool.Name)

		// 動的にツールを作成
		proxyTool := mcp.NewTool(tool.Name,
			mcp.WithDescription(tool.Description),
		)

		// ツール名をキャプチャしてハンドラーを作成
		toolName := tool.Name
		handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return p.handleRemoteToolCall(ctx, toolName, request)
		}

		s.AddTool(proxyTool, handler)
	}
}
// handleRemoteToolCall はリモートツール呼び出しを処理する
func (p *ProxyTransport) handleRemoteToolCall(ctx context.Context, toolName string, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	p.logger.Printf("リモートツール呼び出し: %s", toolName)

	// tools/callリクエストを作成
	toolsCallRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      fmt.Sprintf("tool-call-%s", toolName),
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": request.Params.Arguments,
		},
	}

	requestJSON, err := json.Marshal(toolsCallRequest)
	if err != nil {
		return nil, fmt.Errorf("tools/callリクエストの作成に失敗しました: %v", err)
	}

	// リモートサーバーに送信
	responseJSON, err := p.remoteClient.SendMessage(ctx, requestJSON)
	if err != nil {
		return nil, fmt.Errorf("リモートサーバーへの送信に失敗しました: %v", err)
	}

	// レスポンスを解析
	var response struct {
		Result mcp.CallToolResult `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(responseJSON, &response); err != nil {
		return nil, fmt.Errorf("レスポンスの解析に失敗しました: %v", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("リモートツールエラー: %s", response.Error.Message)
	}

	p.logger.Printf("リモートツール呼び出し成功: %s", toolName)
	return &response.Result, nil
}

// handleProxyCall はプロキシ呼び出しを処理する
func (p *ProxyTransport) handleProxyCall(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	p.logger.Printf("プロキシ呼び出しを処理中: %+v", request)

	// リクエストをJSONに変換
	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("リクエストのJSON変換に失敗しました: %v", err)
	}

	// リモートサーバーに送信
	responseJSON, err := p.remoteClient.SendMessage(ctx, requestJSON)
	if err != nil {
		return nil, fmt.Errorf("リモートサーバーへの送信に失敗しました: %v", err)
	}

	// レスポンスを解析
	var result mcp.CallToolResult
	if err := json.Unmarshal(responseJSON, &result); err != nil {
		return nil, fmt.Errorf("レスポンスの解析に失敗しました: %v", err)
	}

	return &result, nil
}

// Wait はプロキシの終了を待機する
func (p *ProxyTransport) Wait() {
	p.wg.Wait()
}
