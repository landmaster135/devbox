package repositories

import (
	"context"
	"testing"

	config "github.com/landmaster135/devbox/internal/grpc_request/config"
	"github.com/stretchr/testify/assert"
)

func TestGRPCService_TestConnection_Normal(t *testing.T) {
	t.Skip("実際にgRPCサーバを起動しないとテスト出来ない")

	// Arrange
	serverAddress := "localhost:50051"
	useTLS := false

	cfg := config.NewConfig()
	service := NewGRPCClient(cfg)

	// Act
	err := service.TestConnection(context.Background(), serverAddress, useTLS)

	// Assert
	assert.NoError(t, err)
}

func TestGRPCService_ListServices_Normal(t *testing.T) {
	t.Skip("実際にgRPCサーバを起動しないとテスト出来ない")

	// Arrange
	serverAddress := "localhost:50051"
	useTLS := false
	expectedServices := []string{"package.Service1", "package.Service2"}

	cfg := config.NewConfig()
	service := NewGRPCClient(cfg)

	// Act
	services, err := service.ListServices(context.Background(), serverAddress, useTLS)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedServices, services)
}
