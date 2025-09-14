package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	weatherHandler "github.com/landmaster135/devbox/cmd/grpc/handlers/weather_notificator"
	pb "github.com/landmaster135/devbox/cmd/grpc/proto/weather_notificator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	// DefaultPort はデフォルトのポート番号
	DefaultPort = 50051
	// ShutdownTimeout はシャットダウンタイムアウト
	ShutdownTimeout = 30 * time.Second
)

// GRPCServer はgRPCサーバーを管理する構造体
type GRPCServer struct {
	server *grpc.Server
	port   int
}

// NewGRPCServer は新しいGRPCServerインスタンスを作成する
func NewGRPCServer() *GRPCServer {
	grpcServer := &GRPCServer{
		port: getPort(),
	}

	// gRPCサーバーを作成
	grpcServer.server = grpc.NewServer()

	// サービスを登録
	grpcServer.registerServices()

	// リフレクションを有効化（開発用）
	reflection.Register(grpcServer.server)

	return grpcServer
}

// getPort は環境変数またはデフォルトからポート番号を取得する
func getPort() int {
	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		return DefaultPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("無効なGRPC_PORT環境変数: %s, デフォルトポート %d を使用します", portStr, DefaultPort)
		return DefaultPort
	}

	if port <= 0 || port > 65535 {
		log.Printf("無効なポート番号: %d, デフォルトポート %d を使用します", port, DefaultPort)
		return DefaultPort
	}

	return port
}

// registerServices はgRPCサービスを登録する
func (s *GRPCServer) registerServices() {
	// WeatherNotificatorServiceを登録
	weatherNotificatorHandler := weatherHandler.NewWeatherNotificatorHandler()
	pb.RegisterWeatherNotificatorServiceServer(s.server, weatherNotificatorHandler)

	log.Println("gRPCサービスが登録されました:")
	log.Println("  - WeatherNotificatorService")
}

// GracefulShutdown はグレースフルシャットダウンを実行する
func (s *GRPCServer) GracefulShutdown() {
	// シグナルチャネルを作成
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// シグナルを待機
	<-quit
	log.Println("シャットダウンシグナルを受信しました...")

	// グレースフルシャットダウンを実行
	log.Printf("gRPCサーバーをシャットダウンしています（タイムアウト: %v）...", ShutdownTimeout)

	// タイムアウト付きでグレースフルシャットダウン
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	// タイムアウトまたは完了を待機
	select {
	case <-done:
		log.Println("gRPCサーバーが正常にシャットダウンされました")
	case <-time.After(ShutdownTimeout):
		log.Println("シャットダウンタイムアウト、強制終了します")
		s.server.Stop()
	}
}

// Start はgRPCサーバーを開始する
func (s *GRPCServer) Start() {
	// リスナーを作成
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("gRPCサーバーのリスナー作成に失敗しました: %v", err)
	}

	log.Printf("Weather Notification gRPCサーバーを開始します")
	log.Printf("サーバーアドレス: %s", lis.Addr().String())
	log.Printf("リフレクション: 有効")

	// サーバーを開始
	if err := s.server.Serve(lis); err != nil {
		log.Fatalf("gRPCサーバーの開始に失敗しました: %v", err)
	}
}

// GetServer はgRPCサーバーインスタンスを取得する（テスト用）
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// GetPort はポート番号を取得する（テスト用）
func (s *GRPCServer) GetPort() int {
	return s.port
}
