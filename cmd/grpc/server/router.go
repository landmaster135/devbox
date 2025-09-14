package server

import (
	"log"

	"google.golang.org/grpc"

	weatherHandler "github.com/landmaster135/devbox/cmd/grpc/handlers/weather_notificator"
	pb "github.com/landmaster135/devbox/cmd/grpc/proto/weather_notificator"
)

// setupRouter はgRPCサービスをルーターに登録する
func setupRouter(server *grpc.Server) {
	// WeatherNotificatorServiceを登録
	weatherNotificatorHandler := weatherHandler.NewWeatherNotificatorHandler()
	pb.RegisterWeatherNotificatorServiceServer(server, weatherNotificatorHandler)

	// 以降、ここにルートを登録

	log.Println("gRPCサービスが登録されました:")
	log.Println("  - WeatherNotificatorService")
}
