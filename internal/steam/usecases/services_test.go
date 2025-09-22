package usecases

import (
	"errors"
	"os"
	"testing"

	steamAPI "github.com/landmaster135/devbox/internal/steam/infrastructure/steam_api"
)

// NewMockSteamClient は新しいモックSteamClientを作成します
func NewMockSteamClient() *steamAPI.MockSteamClient {
	return &steamAPI.MockSteamClient{
		UsersService: &steamAPI.MockUsersService{},
		AppsService:  &steamAPI.MockAppsService{},
	}
}

// TestNewSteamService_Normal はNewSteamServiceの正常系テスト
func TestNewSteamService_Normal(t *testing.T) {
	// Arrange
	mockClient := NewMockSteamClient()
	mockFileWriter := &MockFileWriter{}
	mockLogger := &MockLogger{}
	mockConcurrencyController := &MockConcurrencyManager{
		GetMaxConcurrencyFunc: func() int { return 5 },
		GetSemaphoreFunc:      func() chan struct{} { return make(chan struct{}, 5) },
	}

	// Act
	service := NewSteamService(mockClient, mockFileWriter, mockLogger, mockConcurrencyController)

	// Assert
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
	if service.client != mockClient {
		t.Error("Expected client to be set correctly")
	}
	if service.fileWriter == nil {
		t.Error("Expected fileWriter to be set")
	}
	if service.logger == nil {
		t.Error("Expected logger to be set")
	}
	if service.concurrencyController == nil {
		t.Error("Expected concurrencyController to be set")
	}
}

// TestNewSteamService_WithNilConfig はNewSteamServiceのnilコンフィグテスト
func TestNewSteamService_WithNilConfig(t *testing.T) {
	// Arrange
	mockClient := NewMockSteamClient()

	// Act
	service := NewSteamService(mockClient, nil, nil, nil)

	// Assert
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
	if service.client != mockClient {
		t.Error("Expected client to be set correctly")
	}

	// デフォルト実装が設定されているかチェック
	if _, ok := service.fileWriter.(*FileWriter); !ok {
		t.Error("Expected default FileWriter to be set")
	}
	if _, ok := service.logger.(*Logger); !ok {
		t.Error("Expected default Logger to be set")
	}
	if _, ok := service.concurrencyController.(*ConcurrencyManager); !ok {
		t.Error("Expected default ConcurrencyController to be set")
	}
}

// TestNewSteamServiceWithAPIKey はNewSteamServiceWithAPIKeyのテスト
func TestNewSteamServiceWithAPIKey(t *testing.T) {
	// Arrange
	apiKey := "test-api-key"

	// Act
	service := NewSteamServiceWithAPIKey(apiKey)

	// Assert
	if service == nil {
		t.Fatal("Expected service to be created, got nil")
	}
	if service.client == nil {
		t.Error("Expected client to be set")
	}
}

// TestMockFileWriter_WriteToFile はMockFileWriterのテスト
func TestMockFileWriter_WriteToFile(t *testing.T) {
	// Arrange
	mockWriter := &MockFileWriter{}
	testData := map[string]string{"test": "data"}
	outputDir := "."
	filename := "test.json"

	// Act
	filepath, err := mockWriter.WriteToFile(testData, outputDir, filename)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if filepath != "" {
		t.Errorf("Expected empty filepath from default mock, got %s", filepath)
	}
}

// TestMockFileWriter_WriteToFileWithError はMockFileWriterのエラーテスト
func TestMockFileWriter_WriteToFileWithError(t *testing.T) {
	// Arrange
	expectedError := errors.New("write error")
	mockWriter := &MockFileWriter{
		WriteToFileFunc: func(data any, outputDir, filename string) (string, error) {
			return "", expectedError
		},
	}
	testData := map[string]string{"test": "data"}
	outputDir := "."
	filename := "test.json"

	// Act
	filepath, err := mockWriter.WriteToFile(testData, outputDir, filename)

	// Assert
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
	if filepath != "" {
		t.Errorf("Expected empty filepath on error, got %s", filepath)
	}
}

// TestMockLogger_Printf はMockLoggerのPrintfテスト
func TestMockLogger_Printf(t *testing.T) {
	// Arrange
	var capturedFormat string
	var capturedArgs []any
	mockLogger := &MockLogger{
		PrintfFunc: func(format string, v ...any) {
			capturedFormat = format
			capturedArgs = v
		},
	}
	format := "Test message: %s"
	arg := "test"

	// Act
	mockLogger.Printf(format, arg)

	// Assert
	if capturedFormat != format {
		t.Errorf("Expected format %s, got %s", format, capturedFormat)
	}
	if len(capturedArgs) != 1 || capturedArgs[0] != arg {
		t.Errorf("Expected args [%s], got %v", arg, capturedArgs)
	}
}

// TestMockLogger_Println はMockLoggerのPrintlnテスト
func TestMockLogger_Println(t *testing.T) {
	// Arrange
	var capturedArgs []any
	mockLogger := &MockLogger{
		PrintlnFunc: func(v ...any) {
			capturedArgs = v
		},
	}
	message := "Test message"

	// Act
	mockLogger.Println(message)

	// Assert
	if len(capturedArgs) != 1 || capturedArgs[0] != message {
		t.Errorf("Expected args [%s], got %v", message, capturedArgs)
	}
}

// TestMockConcurrencyManager はMockConcurrencyManagerのテスト
func TestMockConcurrencyManager(t *testing.T) {
	// Arrange
	maxConcurrency := 3
	semaphore := make(chan struct{}, maxConcurrency)
	controller := &MockConcurrencyManager{
		GetMaxConcurrencyFunc: func() int { return maxConcurrency },
		GetSemaphoreFunc:      func() chan struct{} { return semaphore },
	}

	// Act & Assert
	if controller.GetMaxConcurrency() != maxConcurrency {
		t.Errorf("Expected max concurrency %d, got %d", maxConcurrency, controller.GetMaxConcurrency())
	}

	returnedSemaphore := controller.GetSemaphore()
	if returnedSemaphore != semaphore {
		t.Error("Expected same semaphore instance to be returned")
	}

	// セマフォの容量をテスト
	for i := 0; i < maxConcurrency; i++ {
		select {
		case returnedSemaphore <- struct{}{}:
			// 成功
		default:
			t.Errorf("Expected to be able to send %d items to semaphore", maxConcurrency)
		}
	}

	// 容量を超えた場合のテスト
	select {
	case returnedSemaphore <- struct{}{}:
		t.Error("Expected semaphore to be full")
	default:
		// 期待される動作
	}
}

// TestDefaultFileWriter_WriteToFile はDefaultFileWriterのテスト
func TestDefaultFileWriter_WriteToFile(t *testing.T) {
	// Arrange
	writer := &FileWriter{}
	testData := map[string]string{"test": "data"}
	outputDir := "/tmp"
	filename := "test_default_writer.json"

	// Act
	filepath, err := writer.WriteToFile(testData, outputDir, filename)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	expectedPath := "/tmp/test_default_writer.json"
	if filepath != expectedPath {
		t.Errorf("Expected filepath %s, got %s", expectedPath, filepath)
	}

	// ファイルが作成されたかチェック（クリーンアップも含む）
	if _, err := os.Stat(filepath); err != nil {
		t.Errorf("Expected file to be created, got error: %v", err)
	} else {
		// クリーンアップ
		os.Remove(filepath)
	}
}

// TestDefaultLogger はDefaultLoggerのテスト
func TestDefaultLogger(t *testing.T) {
	// Arrange
	logger := &Logger{}

	// Act & Assert - パニックしないことを確認
	logger.Printf("Test message: %s", "test")
	logger.Println("Test message")

	// ログ出力のテストは標準出力をキャプチャする必要があるが、
	// ここでは単純にパニックしないことを確認
}

// TestDefaultConcurrencyController はDefaultConcurrencyControllerのテスト
func TestDefaultConcurrencyController(t *testing.T) {
	// Arrange
	maxConcurrency := 5
	controller := &ConcurrencyManager{maxConcurrency: maxConcurrency}

	// Act & Assert
	if controller.GetMaxConcurrency() != maxConcurrency {
		t.Errorf("Expected max concurrency %d, got %d", maxConcurrency, controller.GetMaxConcurrency())
	}

	semaphore := controller.GetSemaphore()
	if semaphore == nil {
		t.Error("Expected semaphore to be created")
	}

	// 2回目の呼び出しで同じセマフォが返されることを確認
	semaphore2 := controller.GetSemaphore()
	if semaphore != semaphore2 {
		t.Error("Expected same semaphore instance to be returned")
	}
}
