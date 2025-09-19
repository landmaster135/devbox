package usecases

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	steamAPI "github.com/landmaster135/devbox/internal/steam/infrastructure/steam_api"
)

// MockUsersService はUsersServiceのモック実装
type MockUsersService struct {
	GetOwnedGamesFunc              func(ctx context.Context, steamID string, includeAppInfo, includeFreeGames bool) ([]steamAPI.OwnedGame, error)
	GetUserRecentlyPlayedGamesFunc func(ctx context.Context, steamID string) ([]steamAPI.RecentlyPlayedGame, error)
}

// GetOwnedGames はモックの所有ゲーム取得メソッド
func (m *MockUsersService) GetOwnedGames(ctx context.Context, steamID string, includeAppInfo, includeFreeGames bool) ([]steamAPI.OwnedGame, error) {
	if m.GetOwnedGamesFunc != nil {
		return m.GetOwnedGamesFunc(ctx, steamID, includeAppInfo, includeFreeGames)
	}
	return []steamAPI.OwnedGame{}, nil
}

// GetUserRecentlyPlayedGames はモックの最近プレイしたゲーム取得メソッド
func (m *MockUsersService) GetUserRecentlyPlayedGames(ctx context.Context, steamID string) ([]steamAPI.RecentlyPlayedGame, error) {
	if m.GetUserRecentlyPlayedGamesFunc != nil {
		return m.GetUserRecentlyPlayedGamesFunc(ctx, steamID)
	}
	return []steamAPI.RecentlyPlayedGame{}, nil
}

// MockAppsService はAppsServiceのモック実装
type MockAppsService struct {
	GetUserStatsFunc        func(ctx context.Context, steamID string, appID int) (*steamAPI.UserStats, error)
	GetUserAchievementsFunc func(ctx context.Context, steamID string, appID int, language string) ([]steamAPI.Achievement, error)
}

// GetUserStats はモックのユーザー統計取得メソッド
func (m *MockAppsService) GetUserStats(ctx context.Context, steamID string, appID int) (*steamAPI.UserStats, error) {
	if m.GetUserStatsFunc != nil {
		return m.GetUserStatsFunc(ctx, steamID, appID)
	}
	return &steamAPI.UserStats{}, nil
}

// GetUserAchievements はモックのユーザー実績取得メソッド
func (m *MockAppsService) GetUserAchievements(ctx context.Context, steamID string, appID int, language string) ([]steamAPI.Achievement, error) {
	if m.GetUserAchievementsFunc != nil {
		return m.GetUserAchievementsFunc(ctx, steamID, appID, language)
	}
	return []steamAPI.Achievement{}, nil
}

// MockSteamClient はSteamClientのモック実装
type MockSteamClient struct {
	users *MockUsersService
	apps  *MockAppsService
}

// NewMockSteamClient は新しいモックSteamClientを作成します
func NewMockSteamClient() *MockSteamClient {
	return &MockSteamClient{
		users: &MockUsersService{},
		apps:  &MockAppsService{},
	}
}

// GetUsers はモックのUsersServiceを返します
func (m *MockSteamClient) GetUsers() *steamAPI.UsersService {
	// 実際のUsersServiceの代わりにモックを返すため、型変換が必要
	// テスト用の実装として、実際のUsersServiceのインターフェースを満たすモックを返す
	return nil // この実装は後で改善が必要
}

// GetApps はモックのAppsServiceを返します
func (m *MockSteamClient) GetApps() *steamAPI.AppsService {
	// 実際のAppsServiceの代わりにモックを返すため、型変換が必要
	// テスト用の実装として、実際のAppsServiceのインターフェースを満たすモックを返す
	return nil // この実装は後で改善が必要
}

// SetMockUsers はモックのUsersServiceを設定します
func (m *MockSteamClient) SetMockUsers(users *MockUsersService) {
	m.users = users
}

// SetMockApps はモックのAppsServiceを設定します
func (m *MockSteamClient) SetMockApps(apps *MockAppsService) {
	m.apps = apps
}

// MockFileWriter はFileWriterのモック実装
type MockFileWriter struct {
	WriteToFileFunc func(data any, filename string) error
	WrittenData     []any
	WrittenFiles    []string
}

// WriteToFile はモックのファイル書き込みメソッド
func (m *MockFileWriter) WriteToFile(data any, filename string) error {
	m.WrittenData = append(m.WrittenData, data)
	m.WrittenFiles = append(m.WrittenFiles, filename)

	if m.WriteToFileFunc != nil {
		return m.WriteToFileFunc(data, filename)
	}
	return nil
}

// MockLogger はLoggerのモック実装
type MockLogger struct {
	PrintfFunc     func(format string, v ...any)
	PrintlnFunc    func(v ...any)
	LoggedMessages []string
}

// Printf はモックのフォーマット付きログ出力メソッド
func (m *MockLogger) Printf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)
	m.LoggedMessages = append(m.LoggedMessages, message)

	if m.PrintfFunc != nil {
		m.PrintfFunc(format, v...)
	}
}

// Println はモックのログ出力メソッド
func (m *MockLogger) Println(v ...any) {
	message := fmt.Sprint(v...)
	m.LoggedMessages = append(m.LoggedMessages, message)

	if m.PrintlnFunc != nil {
		m.PrintlnFunc(v...)
	}
}

// MockConcurrencyController は並行処理制御のモック実装
type MockConcurrencyController struct {
	maxConcurrency int
	semaphore      chan struct{}
}

// NewMockConcurrencyController は新しいモック並行処理制御を作成します
func NewMockConcurrencyController(maxConcurrency int) *MockConcurrencyController {
	return &MockConcurrencyController{
		maxConcurrency: maxConcurrency,
		semaphore:      make(chan struct{}, maxConcurrency),
	}
}

// GetSemaphore はモックのセマフォチャネルを返します
func (m *MockConcurrencyController) GetSemaphore() chan struct{} {
	return m.semaphore
}

// GetMaxConcurrency はモックの最大並行数を返します
func (m *MockConcurrencyController) GetMaxConcurrency() int {
	return m.maxConcurrency
}

// TestNewSteamService_Normal はNewSteamServiceの正常系テスト
func TestNewSteamService_Normal(t *testing.T) {
	// Arrange
	mockClient := NewMockSteamClient()
	config := &SteamServiceConfig{
		FileWriter:            &MockFileWriter{},
		Logger:                &MockLogger{},
		ConcurrencyController: NewMockConcurrencyController(5),
	}

	// Act
	service := NewSteamService(mockClient, config)

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
	service := NewSteamService(mockClient, nil)

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
	filename := "test.json"

	// Act
	err := mockWriter.WriteToFile(testData, filename)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(mockWriter.WrittenData) != 1 {
		t.Errorf("Expected 1 written data, got %d", len(mockWriter.WrittenData))
	}
	if len(mockWriter.WrittenFiles) != 1 {
		t.Errorf("Expected 1 written file, got %d", len(mockWriter.WrittenFiles))
	}
	if mockWriter.WrittenFiles[0] != filename {
		t.Errorf("Expected filename %s, got %s", filename, mockWriter.WrittenFiles[0])
	}
}

// TestMockFileWriter_WriteToFileWithError はMockFileWriterのエラーテスト
func TestMockFileWriter_WriteToFileWithError(t *testing.T) {
	// Arrange
	expectedError := errors.New("write error")
	mockWriter := &MockFileWriter{
		WriteToFileFunc: func(data any, filename string) error {
			return expectedError
		},
	}
	testData := map[string]string{"test": "data"}
	filename := "test.json"

	// Act
	err := mockWriter.WriteToFile(testData, filename)

	// Assert
	if err != expectedError {
		t.Errorf("Expected error %v, got %v", expectedError, err)
	}
}

// TestMockLogger_Printf はMockLoggerのPrintfテスト
func TestMockLogger_Printf(t *testing.T) {
	// Arrange
	mockLogger := &MockLogger{}
	format := "Test message: %s"
	arg := "test"

	// Act
	mockLogger.Printf(format, arg)

	// Assert
	if len(mockLogger.LoggedMessages) != 1 {
		t.Errorf("Expected 1 logged message, got %d", len(mockLogger.LoggedMessages))
	}
	expectedMessage := "Test message: test"
	if mockLogger.LoggedMessages[0] != expectedMessage {
		t.Errorf("Expected message %s, got %s", expectedMessage, mockLogger.LoggedMessages[0])
	}
}

// TestMockLogger_Println はMockLoggerのPrintlnテスト
func TestMockLogger_Println(t *testing.T) {
	// Arrange
	mockLogger := &MockLogger{}
	message := "Test message"

	// Act
	mockLogger.Println(message)

	// Assert
	if len(mockLogger.LoggedMessages) != 1 {
		t.Errorf("Expected 1 logged message, got %d", len(mockLogger.LoggedMessages))
	}
	if mockLogger.LoggedMessages[0] != message {
		t.Errorf("Expected message %s, got %s", message, mockLogger.LoggedMessages[0])
	}
}

// TestMockConcurrencyController はMockConcurrencyControllerのテスト
func TestMockConcurrencyController(t *testing.T) {
	// Arrange
	maxConcurrency := 3
	controller := NewMockConcurrencyController(maxConcurrency)

	// Act & Assert
	if controller.GetMaxConcurrency() != maxConcurrency {
		t.Errorf("Expected max concurrency %d, got %d", maxConcurrency, controller.GetMaxConcurrency())
	}

	semaphore := controller.GetSemaphore()
	if semaphore == nil {
		t.Error("Expected semaphore to be created")
	}

	// セマフォの容量をテスト
	for i := 0; i < maxConcurrency; i++ {
		select {
		case semaphore <- struct{}{}:
			// 成功
		default:
			t.Errorf("Expected to be able to send %d items to semaphore", maxConcurrency)
		}
	}

	// 容量を超えた場合のテスト
	select {
	case semaphore <- struct{}{}:
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
	filename := "/tmp/test_default_writer.json"

	// Act
	err := writer.WriteToFile(testData, filename)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// ファイルが作成されたかチェック（クリーンアップも含む）
	if _, err := os.Stat(filename); err != nil {
		t.Errorf("Expected file to be created, got error: %v", err)
	} else {
		// クリーンアップ
		os.Remove(filename)
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
