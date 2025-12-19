package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	config "github.com/landmaster135/devbox/internal/grpc_request/config"
	grpcDomain "github.com/landmaster135/devbox/internal/grpc_request/domain"
	grpcInfra "github.com/landmaster135/devbox/internal/grpc_request/infrastructure"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

func TestNewGRPCService_Normal(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()

	// Act
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Assert
	assert.NotNil(t, service)
}

func TestGRPCService_SendRequestWithData_Normal(t *testing.T) {
	// Arrange
	request := &grpcDomain.GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "package.Service/Method",
		Data:          map[string]interface{}{"key": "value"},
		Metadata:      map[string]string{"authorization": "Bearer token"},
		UseTLS:        false,
		Timeout:       30 * time.Second,
	}

	expectedResponse := &grpcDomain.GRPCResponse{
		Data:       map[string]interface{}{"result": "success"},
		StatusCode: 0,
		StatusMsg:  "OK",
		Duration:   100 * time.Millisecond,
	}

	mockRepo := &grpcInfra.MockGRPCClient{
		SendRequestFunc: func(ctx context.Context, req *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error) {
			assert.Equal(t, request, req)
			return expectedResponse, nil
		},
	}

	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	response, err := service.SendRequestWithData(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestGRPCService_SendRequestWithData_WithDefaultTimeout_Normal(t *testing.T) {
	// Arrange
	cfg := config.NewConfig()

	request := &grpcDomain.GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "package.Service/Method",
		Data:          map[string]interface{}{"key": "value"},
		UseTLS:        false,
		Timeout:       0, // デフォルトタイムアウトを使用
	}

	expectedResponse := &grpcDomain.GRPCResponse{
		Data:       map[string]interface{}{"result": "success"},
		StatusCode: 0,
		StatusMsg:  "OK",
		Duration:   100 * time.Millisecond,
	}

	mockRepo := &grpcInfra.MockGRPCClient{
		SendRequestFunc: func(ctx context.Context, req *grpcDomain.GRPCRequest) (*grpcDomain.GRPCResponse, error) {
			assert.Equal(t, cfg.DefaultTimeout, req.Timeout)
			return expectedResponse, nil
		},
	}

	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	response, err := service.SendRequestWithData(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
	assert.Equal(t, cfg.DefaultTimeout, request.Timeout)
}

func TestGRPCService_TestConnection_Normal(t *testing.T) {
	// Arrange
	serverAddress := "localhost:50051"
	useTLS := false

	mockRepo := &grpcInfra.MockGRPCClient{
		TestConnectionFunc: func(ctx context.Context, addr string, tls bool) error {
			assert.Equal(t, serverAddress, addr)
			assert.Equal(t, useTLS, tls)
			return nil
		},
	}

	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	err := service.GetRepository().TestConnection(context.Background(), serverAddress, useTLS)

	// Assert
	assert.NoError(t, err)
}

func TestGRPCService_ListServices_Normal(t *testing.T) {
	// Arrange
	serverAddress := "localhost:50051"
	useTLS := false
	expectedServices := []string{"package.Service1", "package.Service2"}

	mockRepo := &grpcInfra.MockGRPCClient{
		ListServicesFunc: func(ctx context.Context, addr string, tls bool) ([]string, error) {
			assert.Equal(t, serverAddress, addr)
			assert.Equal(t, useTLS, tls)
			return expectedServices, nil
		},
	}

	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	services, err := service.GetRepository().ListServices(context.Background(), serverAddress, useTLS)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedServices, services)
}

func TestGRPCService_FormatResponse_Normal(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	response := &grpcDomain.GRPCResponse{
		Data:       map[string]interface{}{"result": "success", "count": 42},
		Metadata:   map[string][]string{"content-type": {"application/grpc"}},
		StatusCode: 0,
		StatusMsg:  "OK",
		Duration:   100 * time.Millisecond,
	}

	// Act
	result, err := service.FormatResponse(response)

	// Assert
	assert.NoError(t, err)
	assert.Contains(t, result, "Status: 0 - OK")
	assert.Contains(t, result, "Duration: 100ms")
	assert.Contains(t, result, "content-type: application/grpc")
	assert.Contains(t, result, "Response Data:")
	assert.Contains(t, result, "\"result\": \"success\"")
	assert.Contains(t, result, "\"count\": 42")
}

func TestGRPCService_LoadJSONFile_Normal(t *testing.T) {
	// Arrange
	expectedData := map[string]interface{}{
		"key":    "value",
		"number": float64(42),
		"nested": map[string]interface{}{"inner": "data"},
	}

	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{
		LoadJSONFileFunc: func(filePath string) (map[string]interface{}, error) {
			assert.Equal(t, "test.json", filePath)
			return expectedData, nil
		},
	}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options)

	// Act
	result, err := service.GetFileReader().LoadJSONFile("test.json")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func TestGRPCService_validateRequest_Normal(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options).(*GRPCService)

	request := &grpcDomain.GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "package.Service/Method",
		Data:          map[string]interface{}{"key": "value"},
	}

	// Act
	err := service.validateRequest(request)

	// Assert
	assert.NoError(t, err)
}

func TestGRPCService_validateRequest_EmptyServerAddress(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options).(*GRPCService)

	request := &grpcDomain.GRPCRequest{
		ServerAddress: "",
		Method:        "package.Service/Method",
		Data:          map[string]interface{}{"key": "value"},
	}

	// Act
	err := service.validateRequest(request)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サーバーアドレスが指定されていません")
}

func TestGRPCService_validateRequest_EmptyMethod(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options).(*GRPCService)

	request := &grpcDomain.GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "",
		Data:          map[string]interface{}{"key": "value"},
	}

	// Act
	err := service.validateRequest(request)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "メソッドが指定されていません")
}

func TestGRPCService_validateRequest_InvalidMethodFormat(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options).(*GRPCService)

	request := &grpcDomain.GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "InvalidMethod",
		Data:          map[string]interface{}{"key": "value"},
	}

	// Act
	err := service.validateRequest(request)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "メソッド名の形式が正しくありません")
}

func TestGRPCService_validateRequest_InvalidServiceFormat(t *testing.T) {
	// Arrange
	mockRepo := &grpcInfra.MockGRPCClient{}
	cfg := config.NewConfig()
	mockFileReader := &grpcInfra.MockFileReader{}
	options := NewCLIOptions()
	service := NewGRPCService(mockRepo, cfg, mockFileReader, options).(*GRPCService)

	request := &grpcDomain.GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "Service/Method",
		Data:          map[string]interface{}{"key": "value"},
	}

	// Act
	err := service.validateRequest(request)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サービス名の形式が正しくありません")
}
