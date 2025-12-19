package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGRPCRequest_Normal(t *testing.T) {
	// Arrange
	request := &GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "package.Service/Method",
		Data:          map[string]interface{}{"key": "value"},
		Metadata:      map[string]string{"authorization": "Bearer token"},
		UseTLS:        false,
		Timeout:       30 * time.Second,
	}

	// Act & Assert
	assert.Equal(t, "localhost:50051", request.ServerAddress)
	assert.Equal(t, "package.Service/Method", request.Method)
	assert.Equal(t, map[string]interface{}{"key": "value"}, request.Data)
	assert.Equal(t, map[string]string{"authorization": "Bearer token"}, request.Metadata)
	assert.False(t, request.UseTLS)
	assert.Equal(t, 30*time.Second, request.Timeout)
}

func TestGRPCResponse_Normal(t *testing.T) {
	// Arrange
	response := &GRPCResponse{
		Data:       map[string]interface{}{"result": "success"},
		Metadata:   map[string][]string{"content-type": {"application/grpc"}},
		StatusCode: 0,
		StatusMsg:  "OK",
		Duration:   100 * time.Millisecond,
	}

	// Act & Assert
	assert.Equal(t, map[string]interface{}{"result": "success"}, response.Data)
	assert.Equal(t, map[string][]string{"content-type": {"application/grpc"}}, response.Metadata)
	assert.Equal(t, 0, response.StatusCode)
	assert.Equal(t, "OK", response.StatusMsg)
	assert.Equal(t, 100*time.Millisecond, response.Duration)
}

func TestGRPCError_Error_Normal(t *testing.T) {
	// Arrange
	grpcError := &GRPCError{
		Code:    5,
		Message: "Not Found",
		Details: "The requested resource was not found",
	}

	// Act
	result := grpcError.Error()

	// Assert
	assert.Equal(t, "Not Found", result)
}

func TestGRPCError_Normal(t *testing.T) {
	// Arrange
	grpcError := &GRPCError{
		Code:    14,
		Message: "Unavailable",
		Details: "Service temporarily unavailable",
	}

	// Act & Assert
	assert.Equal(t, 14, grpcError.Code)
	assert.Equal(t, "Unavailable", grpcError.Message)
	assert.Equal(t, "Service temporarily unavailable", grpcError.Details)
}
