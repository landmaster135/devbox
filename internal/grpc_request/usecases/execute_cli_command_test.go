package usecases

import (
	"context"
	"strings"
	"testing"
	"time"

	config "github.com/landmaster135/devbox/internal/grpc_request/config"
	grpcDomain "github.com/landmaster135/devbox/internal/grpc_request/domain"
	grpcInfra "github.com/landmaster135/devbox/internal/grpc_request/infrastructure"
)

func TestGRPCService_ExecuteCLICommand_TestConnection_Normal(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{
		TestConnectionFunc: func(ctx context.Context, serverAddress string, useTLS bool) error {
			return nil
		},
	}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := &CLIOptions{
		Server:   "localhost:50051",
		TestConn: true,
	}
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	result, err := service.ExecuteCLICommand(context.Background())

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expected := "接続テスト成功"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGRPCService_ExecuteCLICommand_ListServices_Normal(t *testing.T) {
	// Arrange
	expectedServices := []string{"service1", "service2", "service3"}
	mockRepo := &grpcInfra.MockGRPCClient{
		ListServicesFunc: func(ctx context.Context, serverAddress string, useTLS bool) ([]string, error) {
			return expectedServices, nil
		},
	}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := &CLIOptions{
		Server:       "localhost:50051",
		ListServices: true,
	}
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	result, err := service.ExecuteCLICommand(context.Background())

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "利用可能なサービス:") {
		t.Errorf("Expected result to contain service list header")
	}

	for _, service := range expectedServices {
		if !strings.Contains(result, service) {
			t.Errorf("Expected result to contain service %s", service)
		}
	}
}

func TestGRPCService_ExecuteCLICommand_SendRequest_Normal(t *testing.T) {
	// Arrange
	mockData := map[string]interface{}{
		"message": "test",
	}
	mockResponse := &grpcDomain.GRPCResponse{
		Data:       mockData,
		StatusCode: 200,
		StatusMsg:  "OK",
		Duration:   100 * time.Millisecond,
		Metadata:   make(map[string][]string),
	}

	mockRepo := &grpcInfra.MockGRPCClient{
		SendRequestFunc: func(ctx context.Context, request *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error) {
			return mockResponse, nil
		},
	}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{
		LoadJSONFileFunc: func(filename string) (map[string]interface{}, error) {
			return mockData, nil
		},
	}
	options := &CLIOptions{
		Server:   "localhost:50051",
		Method:   "test.Service/TestMethod",
		JSONFile: "test.json",
		Token:    "test-token",
		Timeout:  30 * time.Second,
	}
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	result, err := service.ExecuteCLICommand(context.Background())

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !strings.Contains(result, "Status: 200 - OK") {
		t.Errorf("Expected result to contain status information")
	}

	if !strings.Contains(result, "Duration:") {
		t.Errorf("Expected result to contain duration information")
	}

	if !strings.Contains(result, "Response Data:") {
		t.Errorf("Expected result to contain response data")
	}
}

func TestGRPCService_ExecuteCLICommand_ValidationError(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := &CLIOptions{
		// Server is empty, should cause validation error
	}
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	result, err := service.ExecuteCLICommand(context.Background())

	// Assert
	if err == nil {
		t.Error("Expected validation error, got nil")
	}

	if result != "" {
		t.Errorf("Expected empty result on error, got %s", result)
	}

	expectedError := "サーバーアドレスが指定されていません"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %s, got %s", expectedError, err.Error())
	}
}

func TestGRPCService_ExecuteCLICommand_TestConnectionError(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{
		TestConnectionFunc: func(ctx context.Context, serverAddress string, useTLS bool) error {
			return &grpcDomain.GRPCError{
				Code:    1,
				Message: "connection failed",
			}
		},
	}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := &CLIOptions{
		Server:   "localhost:50051",
		TestConn: true,
	}
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	result, err := service.ExecuteCLICommand(context.Background())

	// Assert
	if err == nil {
		t.Error("Expected connection error, got nil")
	}

	if result != "" {
		t.Errorf("Expected empty result on error, got %s", result)
	}

	expectedError := "接続テストに失敗しました"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %s, got %s", expectedError, err.Error())
	}
}
