package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	logging "github.com/landmaster135/devbox/internal/logging"
)

const (
	// DefaultPort はデフォルトのポート番号
	DefaultPort = 8080
	// DefaultReadTimeout は読み取りタイムアウト
	DefaultReadTimeout = 30 * time.Second
	// DefaultWriteTimeout は書き込みタイムアウト
	DefaultWriteTimeout = 30 * time.Second
	// DefaultIdleTimeout はアイドルタイムアウト
	DefaultIdleTimeout = 60 * time.Second
	// ShutdownTimeout はシャットダウンタイムアウト
	ShutdownTimeout = 30 * time.Second
)

// HTTPServer はHTTPサーバーを管理する構造体
type HTTPServer struct {
	server *http.Server
	port   int
	logger *logging.StructuredLogger
}

// NewHTTPServer は新しいHTTPServerインスタンスを作成する
func NewHTTPServer(logger *logging.StructuredLogger) *HTTPServer {
	appLogger := logging.Ensure(logger)
	httpServer := &HTTPServer{
		port:   getPort(appLogger),
		logger: appLogger,
	}

	// ルーターを初期化
	mux := setupRouter(appLogger)

	// サーバーを設定
	httpServer.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", httpServer.port),
		Handler:      mux,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
	}

	return httpServer
}

// getPort は環境変数またはデフォルトからポート番号を取得する
func getPort(logger *logging.StructuredLogger) int {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return DefaultPort
	}

	port, err := strconv.Atoi(portStr)
	taggedLogger := logging.Ensure(logger).WithTags("config")
	if err != nil {
		taggedLogger.Warnf("無効なPORT環境変数: %s, デフォルトポート %d を使用します", portStr, DefaultPort)
		return DefaultPort
	}

	if port <= 0 || port > 65535 {
		taggedLogger.Warnf("無効なポート番号: %d, デフォルトポート %d を使用します", port, DefaultPort)
		return DefaultPort
	}

	return port
}

// GracefulShutdown はグレースフルシャットダウンを実行する
func (s *HTTPServer) GracefulShutdown() {
	// シグナルチャネルを作成
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// シグナルを待機
	<-quit
	taggedLogger := s.logger.WithTags("shutdown")
	taggedLogger.Infof("シャットダウンシグナルを受信しました...")

	// シャットダウンのコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	// サーバーをグレースフルシャットダウン
	taggedLogger.Infof("サーバーをシャットダウンしています（タイムアウト: %v）...", ShutdownTimeout)
	if err := s.server.Shutdown(ctx); err != nil {
		taggedLogger.Errorf("サーバーのシャットダウンでエラーが発生しました: %v", err)
	} else {
		taggedLogger.Infof("サーバーが正常にシャットダウンされました")
	}
}

// Start はHTTPサーバーを開始する
func (s *HTTPServer) Start() {
	lifecycleLogger := s.logger.WithTags("lifecycle")
	lifecycleLogger.Infof("Weather Notification APIサーバーを開始します")
	lifecycleLogger.Infof("サーバーアドレス: %s", s.server.Addr)
	lifecycleLogger.Infof("読み取りタイムアウト: %v", s.server.ReadTimeout)
	lifecycleLogger.Infof("書き込みタイムアウト: %v", s.server.WriteTimeout)
	lifecycleLogger.Infof("アイドルタイムアウト: %v", s.server.IdleTimeout)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		lifecycleLogger.Errorf("サーバーの開始に失敗しました: %v", err)
		os.Exit(1)
	}
}
