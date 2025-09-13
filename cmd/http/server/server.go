package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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

// getPort は環境変数またはデフォルトからポート番号を取得する
func getPort() int {
	portStr := os.Getenv("PORT")
	if portStr == "" {
		return DefaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("無効なPORT環境変数: %s, デフォルトポート %d を使用します", portStr, DefaultPort)
		return DefaultPort
	}

	if port <= 0 || port > 65535 {
		log.Printf("無効なポート番号: %d, デフォルトポート %d を使用します", port, DefaultPort)
		return DefaultPort
	}

	return port
}

// SetupServer はHTTPサーバーを設定する
func SetupServer() *http.Server {
	// ポート番号を取得
	port := getPort()

	// ルーターを初期化
	mux := router()

	// サーバーを設定
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
	}

	return server
}

// GracefulShutdown はグレースフルシャットダウンを実行する
func GracefulShutdown(server *http.Server) {
	// シグナルチャネルを作成
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// シグナルを待機
	<-quit
	log.Println("シャットダウンシグナルを受信しました...")

	// シャットダウンのコンテキストを作成
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	// サーバーをグレースフルシャットダウン
	log.Printf("サーバーをシャットダウンしています（タイムアウト: %v）...", ShutdownTimeout)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("サーバーのシャットダウンでエラーが発生しました: %v", err)
	} else {
		log.Println("サーバーが正常にシャットダウンされました")
	}
}

// StartServer はHTTPサーバーを開始する
func StartServer(server *http.Server) {
	log.Printf("Weather Notification APIサーバーを開始します")
	log.Printf("サーバーアドレス: %s", server.Addr)
	log.Printf("読み取りタイムアウト: %v", server.ReadTimeout)
	log.Printf("書き込みタイムアウト: %v", server.WriteTimeout)
	log.Printf("アイドルタイムアウト: %v", server.IdleTimeout)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("サーバーの開始に失敗しました: %v", err)
	}
}
