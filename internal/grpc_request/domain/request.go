package models

import (
	"time"
)

// GRPCRequest はgRPCリクエストの情報を表します
type GRPCRequest struct {
	ServerAddress string            `json:"server_address"`
	Method        string            `json:"method"`
	Data          map[string]interface{} `json:"data"`
	Metadata      map[string]string `json:"metadata"`
	UseTLS        bool              `json:"use_tls"`
	Timeout       time.Duration     `json:"timeout"`
}

// GRPCResponse はgRPCレスポンスの情報を表します
type GRPCResponse struct {
	Data       map[string]interface{} `json:"data"`
	Metadata   map[string][]string    `json:"metadata"`
	StatusCode int                    `json:"status_code"`
	StatusMsg  string                 `json:"status_message"`
	Duration   time.Duration          `json:"duration"`
}

// GRPCError はgRPCエラーの情報を表します
type GRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

// Error はerrorインターフェースを実装します
func (e *GRPCError) Error() string {
	return e.Message
}
