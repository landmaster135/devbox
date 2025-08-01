package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"

	config "github.com/landmaster135/devbox/internal/mcp_remote/config"
)

// TestProxyService_NewProxyService は新しいProxyServiceの作成をテストする
func TestProxyService_NewProxyService_Normal(t *testing.T) {
	service := NewProxyService()

	assert.NotNil(t, service)
	assert.NotNil(t, service.logger)
}

// TestProxyService_validateServerURL はサーバーURL検証をテストする
func TestProxyService_validateServerURL_Normal(t *testing.T) {
	service := NewProxyService()
	cfg := &config.Config{
		ServerURL: "https://example.com/mcp/sse",
		Debug:     false,
	}

	err := service.validateServerURL(cfg)
	assert.NoError(t, err)
}

// TestAuthService_NewAuthService は新しいAuthServiceの作成をテストする
func TestAuthService_NewAuthService_Normal(t *testing.T) {
	service := NewAuthService()

	assert.NotNil(t, service)
	assert.NotNil(t, service.logger)
}

// TestAuthService_InitializeAuth はOAuth認証初期化をテストする
func TestAuthService_InitializeAuth_Normal(t *testing.T) {
	service := NewAuthService()
	cfg := &config.Config{
		ServerURL: "https://example.com/mcp/sse",
		Debug:     false,
	}

	err := service.InitializeAuth(cfg)
	assert.NoError(t, err)
}

// TestTransportService_NewTransportService は新しいTransportServiceの作成をテストする
func TestTransportService_NewTransportService_Normal(t *testing.T) {
	service := NewTransportService()

	assert.NotNil(t, service)
	assert.NotNil(t, service.logger)
}

// TestTransportService_CreateTransport はトランスポート作成をテストする
func TestTransportService_CreateTransport_Normal(t *testing.T) {
	service := NewTransportService()
	cfg := &config.Config{
		ServerURL:         "https://example.com/mcp/sse",
		TransportStrategy: config.TransportHTTPFirst,
		Debug:             false,
	}

	err := service.CreateTransport(cfg)
	assert.NoError(t, err)
}
