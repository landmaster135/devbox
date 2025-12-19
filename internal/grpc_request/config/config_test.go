package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_Normal(t *testing.T) {
	// Act
	cfg := NewConfig()

	// Assert
	assert.NotNil(t, cfg)
	assert.Equal(t, 30*time.Second, cfg.DefaultTimeout)
	assert.Equal(t, 3, cfg.MaxRetries)
	assert.Equal(t, 1*time.Second, cfg.RetryDelay)
}

func TestConfig_Normal(t *testing.T) {
	// Arrange
	cfg := &Config{
		DefaultTimeout: 60 * time.Second,
		MaxRetries:     5,
		RetryDelay:     2 * time.Second,
	}

	// Act & Assert
	assert.Equal(t, 60*time.Second, cfg.DefaultTimeout)
	assert.Equal(t, 5, cfg.MaxRetries)
	assert.Equal(t, 2*time.Second, cfg.RetryDelay)
}

func TestNewConnectionConfig_Normal(t *testing.T) {
	// Arrange
	serverAddress := "localhost:50051"
	useTLS := true

	// Act
	connCfg := NewConnectionConfig(serverAddress, useTLS)

	// Assert
	assert.NotNil(t, connCfg)
	assert.Equal(t, "localhost:50051", connCfg.ServerAddress)
	assert.True(t, connCfg.UseTLS)
	assert.False(t, connCfg.Insecure)
	assert.Empty(t, connCfg.CertFile)
	assert.Empty(t, connCfg.KeyFile)
	assert.Empty(t, connCfg.CAFile)
}

func TestNewConnectionConfig_WithoutTLS_Normal(t *testing.T) {
	// Arrange
	serverAddress := "localhost:50051"
	useTLS := false

	// Act
	connCfg := NewConnectionConfig(serverAddress, useTLS)

	// Assert
	assert.NotNil(t, connCfg)
	assert.Equal(t, "localhost:50051", connCfg.ServerAddress)
	assert.False(t, connCfg.UseTLS)
	assert.True(t, connCfg.Insecure)
}

func TestConnectionConfig_Normal(t *testing.T) {
	// Arrange
	connCfg := &ConnectionConfig{
		ServerAddress: "example.com:443",
		UseTLS:        true,
		Insecure:      false,
		CertFile:      "/path/to/cert.pem",
		KeyFile:       "/path/to/key.pem",
		CAFile:        "/path/to/ca.pem",
	}

	// Act & Assert
	assert.Equal(t, "example.com:443", connCfg.ServerAddress)
	assert.True(t, connCfg.UseTLS)
	assert.False(t, connCfg.Insecure)
	assert.Equal(t, "/path/to/cert.pem", connCfg.CertFile)
	assert.Equal(t, "/path/to/key.pem", connCfg.KeyFile)
	assert.Equal(t, "/path/to/ca.pem", connCfg.CAFile)
}
