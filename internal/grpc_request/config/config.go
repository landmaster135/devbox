package config

import (
	"time"
)

// Config はgRPCクライアントの設定を表します
type Config struct {
	DefaultTimeout time.Duration `json:"default_timeout"`
	MaxRetries     int           `json:"max_retries"`
	RetryDelay     time.Duration `json:"retry_delay"`
}

// NewConfig は新しい設定インスタンスを作成します
func NewConfig() *Config {
	return &Config{
		DefaultTimeout: 30 * time.Second,
		MaxRetries:     3,
		RetryDelay:     1 * time.Second,
	}
}

// ConnectionConfig はgRPC接続の設定を表します
type ConnectionConfig struct {
	ServerAddress string
	UseTLS        bool
	Insecure      bool
	CertFile      string
	KeyFile       string
	CAFile        string
}

// NewConnectionConfig は新しい接続設定インスタンスを作成します
func NewConnectionConfig(serverAddress string, useTLS bool) *ConnectionConfig {
	return &ConnectionConfig{
		ServerAddress: serverAddress,
		UseTLS:        useTLS,
		Insecure:      !useTLS,
	}
}
