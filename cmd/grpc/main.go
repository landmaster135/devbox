package main

import (
	"log"

	grpcServer "github.com/landmaster135/devbox/cmd/grpc/server"
)

func main() {
	log.Println("Weather Notification gRPCサーバーを初期化しています...")

	// サーバーを作成
	server := grpcServer.NewGRPCServer()

	// グレースフルシャットダウンのゴルーチンを開始
	go server.GracefulShutdown()

	// サーバーを開始
	server.Start()
}
