package usecases

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/landmaster135/devbox/internal/grpc_request/config"
	"github.com/landmaster135/devbox/internal/grpc_request/domain/models"
	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#
// MockGRPCRepository はGRPCRepositoryのモック実装
type MockGRPCRepository struct {
	SendRequestFunc    func(ctx context.Context, request *models.GRPCRequest) (*models.GRPCResponse, error)
	TestConnectionFunc func(ctx context.Context, serverAddress string, useTLS bool) error
	ListServicesFunc   func(ctx context.Context, serverAddress string, useTLS bool) ([]string, error)
}

func (m *MockGRPCRepository) SendRequest(ctx context.Context, request *models.GRPCRequest) (*models.GRPCResponse, error) {
	return m.SendRequestFunc(ctx, request)
}

func (m *MockGRPCRepository) TestConnection(ctx context.Context, serverAddress string, useTLS bool) error {
	return m.TestConnectionFunc(ctx, serverAddress, useTLS)
}

func (m *MockGRPCRepository) ListServices(ctx context.Context, serverAddress string, useTLS bool) ([]string, error) {
	return m.ListServicesFunc(ctx, serverAddress, useTLS)
}

func TestNewGRPCService_Normal(t *testing.T) {
	// Arrange
	mockRepo := &MockGRPCRepository{}
	cfg := config.NewConfig()

	// Act
	service := NewGRPCService(mockRepo, cfg)

	// Assert
	assert.NotNil(t, service)
}

func TestGRPCService_SendRequestWithData_Normal(t *testing.T) {
	// Arrange
	request := &models.GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "package.Service/Method",
		Data:          map[string]interface{}{"key": "value"},
		Metadata:      map[string]string{"authorization": "Bearer token"},
		UseTLS:        false,
		Timeout:       30 * time.Second,
	}

	expectedResponse := &models.GRPCResponse{
		Data:       map[string]interface{}{"result": "success"},
		StatusCode: 0,
		StatusMsg:  "OK",
		Duration:   100 * time.Millisecond,
	}

	mockRepo := &MockGRPCRepository{
		SendRequestFunc: func(ctx context.Context, req *models.GRPCRequest) (*models.GRPCResponse, error) {
			assert.Equal(t, request, req)
			return expectedResponse, nil
		},
	}

	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg)

	// Act
	response, err := service.SendRequestWithData(context.Background(), request)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestGRPCService_SendRequestWithData_WithDefaultTimeout_Normal(t *testing.T) {
	// Arrange
	cfg := config.NewConfig()

	request := &models.GRPCRequest{
		ServerAddress: "localhost:50051",
		Method:        "package.Service/Method",
		Data:          map[string]interface{}{"key": "value"},
		UseTLS:        false,
		Timeout:       0, // デフォルトタイムアウトを使用
	}

	expectedResponse := &models.GRPCResponse{
		Data:       map[string]interface{}{"result": "success"},
		StatusCode: 0,
		StatusMsg:  "OK",
		Duration:   100 * time.Millisecond,
	}

	mockRepo := &MockGRPCRepository{
		SendRequestFunc: func(ctx context.Context, req *models.GRPCRequest) (*models.GRPCResponse, error) {
			assert.Equal(t, cfg.DefaultTimeout, req.Timeout)
			return expectedResponse, nil
		},
	}

	service := NewGRPCService(mockRepo, cfg)

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

	mockRepo := &MockGRPCRepository{
		TestConnectionFunc: func(ctx context.Context, addr string, tls bool) error {
			assert.Equal(t, serverAddress, addr)
			assert.Equal(t, useTLS, tls)
			return nil
		},
	}

	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg)

	// Act
	err := service.TestConnection(context.Background(), serverAddress, useTLS)

	// Assert
	assert.NoError(t, err)
}

func TestGRPCService_ListServices_Normal(t *testing.T) {
	// Arrange
	serverAddress := "localhost:50051"
	useTLS := false
	expectedServices := []string{"package.Service1", "package.Service2"}

	mockRepo := &MockGRPCRepository{
		ListServicesFunc: func(ctx context.Context, addr string, tls bool) ([]string, error) {
			assert.Equal(t, serverAddress, addr)
			assert.Equal(t, useTLS, tls)
			return expectedServices, nil
		},
	}

	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg)

	// Act
	services, err := service.ListServices(context.Background(), serverAddress, useTLS)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedServices, services)
}

func TestGRPCService_FormatResponse_Normal(t *testing.T) {
	// Arrange
	mockRepo := &MockGRPCRepository{}
	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg)

	response := &models.GRPCResponse{
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
	mockRepo := &MockGRPCRepository{}
	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg)

	// テスト用のJSONファイルを作成
	tempFile, err := os.CreateTemp("", "test_*.json")
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	jsonContent := `{"key": "value", "number": 42, "nested": {"inner": "data"}}`
	_, err = tempFile.WriteString(jsonContent)
	assert.NoError(t, err)
	tempFile.Close()

	// Act
	result, err := service.LoadJSONFile(tempFile.Name())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "value", result["key"])
	assert.Equal(t, float64(42), result["number"]) // JSONの数値はfloat64になる
	nested := result["nested"].(map[string]interface{})
	assert.Equal(t, "data", nested["inner"])
}

func TestGRPCService_validateRequest_Normal(t *testing.T) {
	// Arrange
	mockRepo := &MockGRPCRepository{}
	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg).(*grpcService)

	request := &models.GRPCRequest{
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
	mockRepo := &MockGRPCRepository{}
	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg).(*grpcService)

	request := &models.GRPCRequest{
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
	mockRepo := &MockGRPCRepository{}
	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg).(*grpcService)

	request := &models.GRPCRequest{
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
	mockRepo := &MockGRPCRepository{}
	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg).(*grpcService)

	request := &models.GRPCRequest{
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
	mockRepo := &MockGRPCRepository{}
	cfg := config.NewConfig()
	service := NewGRPCService(mockRepo, cfg).(*grpcService)

	request := &models.GRPCRequest{
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
