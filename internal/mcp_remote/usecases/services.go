package usecases

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/landmaster135/devbox/internal/mcp_remote/config"
)

// ProxyService はMCPプロキシサービスを提供する
type ProxyService struct {
	logger *log.Logger
}

// NewProxyService は新しいProxyServiceを作成する
func NewProxyService() *ProxyService {
	return &ProxyService{
		logger: log.New(os.Stderr, "[mcp-remote] ", log.LstdFlags),
	}
}

// RunProxy はプロキシを実行する
func (s *ProxyService) RunProxy(cfg *config.Config) error {
	s.logger.Printf("MCP Remote Proxy を開始します")
	s.logger.Printf("サーバーURL: %s", cfg.ServerURL)
	s.logger.Printf("トランスポート戦略: %s", cfg.TransportStrategy)

	if cfg.Debug {
		s.logger.Printf("デバッグモードが有効です")
	}

	// URL検証
	if err := s.validateServerURL(cfg); err != nil {
		return fmt.Errorf("サーバーURL検証エラー: %v", err)
	}

	// シグナルハンドリングの設定
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// プロキシの初期化と実行
	if err := s.initializeProxy(cfg); err != nil {
		return fmt.Errorf("プロキシ初期化エラー: %v", err)
	}

	s.logger.Printf("プロキシが正常に開始されました。終了するにはCtrl+Cを押してください")

	// シグナル待機
	<-sigChan
	s.logger.Printf("シャットダウンシグナルを受信しました")

	// クリーンアップ
	if err := s.cleanup(); err != nil {
		s.logger.Printf("クリーンアップ中にエラーが発生しました: %v", err)
	}

	s.logger.Printf("プロキシを正常に終了しました")
	return nil
}

// validateServerURL はサーバーURLを検証する
func (s *ProxyService) validateServerURL(cfg *config.Config) error {
	// TODO: URL検証の実装
	// - HTTPSまたはlocalhostのHTTPのみ許可
	// - allow-httpフラグの考慮
	s.logger.Printf("サーバーURL検証: %s", cfg.ServerURL)
	return nil
}

// initializeProxy はプロキシを初期化する
func (s *ProxyService) initializeProxy(cfg *config.Config) error {
	// TODO: プロキシ初期化の実装
	// 1. OAuth認証の設定
	// 2. トランスポートの選択と初期化
	// 3. 双方向プロキシの開始
	s.logger.Printf("プロキシを初期化中...")
	return nil
}

// cleanup はリソースのクリーンアップを行う
func (s *ProxyService) cleanup() error {
	// TODO: クリーンアップの実装
	// - 接続の切断
	// - リソースの解放
	s.logger.Printf("リソースをクリーンアップ中...")
	return nil
}

// AuthService はOAuth認証サービスを提供する
type AuthService struct {
	logger *log.Logger
}

// NewAuthService は新しいAuthServiceを作成する
func NewAuthService() *AuthService {
	return &AuthService{
		logger: log.New(os.Stderr, "[auth] ", log.LstdFlags),
	}
}

// InitializeAuth はOAuth認証を初期化する
func (s *AuthService) InitializeAuth(cfg *config.Config) error {
	// TODO: OAuth認証初期化の実装
	s.logger.Printf("OAuth認証を初期化中...")
	return nil
}

// TransportService はトランスポート管理サービスを提供する
type TransportService struct {
	logger *log.Logger
}

// NewTransportService は新しいTransportServiceを作成する
func NewTransportService() *TransportService {
	return &TransportService{
		logger: log.New(os.Stderr, "[transport] ", log.LstdFlags),
	}
}

// CreateTransport はトランスポートを作成する
func (s *TransportService) CreateTransport(cfg *config.Config) error {
	// TODO: トランスポート作成の実装
	s.logger.Printf("トランスポートを作成中: %s", cfg.TransportStrategy)
	return nil
}
