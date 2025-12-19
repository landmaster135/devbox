package main

import (
	"log"

	httpServer "github.com/landmaster135/devbox/cmd/http/server"
)

func main() {
	log.Println("Weather Notification API サーバーを初期化しています...")

	// サーバーを作成
	server := httpServer.NewHTTPServer()

	// グレースフルシャットダウンのゴルーチンを開始
	go server.GracefulShutdown()

	// サーバーを開始
	server.Start()
}
