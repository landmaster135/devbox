package main

import (
	"log"

	httpServer "github.com/landmaster135/devbox/cmd/http/server"
)

func main() {
	log.Println("Weather Notification API サーバーを初期化しています...")

	// サーバーを設定
	server := httpServer.SetupServer()

	// グレースフルシャットダウンのゴルーチンを開始
	go httpServer.GracefulShutdown(server)

	// サーバーを開始
	httpServer.StartServer(server)
}
