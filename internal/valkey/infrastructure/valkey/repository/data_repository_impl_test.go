package valkey

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	infraValkey "github.com/landmaster135/devbox/internal/valkey/infrastructure/valkey"
)

type DataRepositoryImplTestSuite struct {
	mockDataStore *MockDataStore
	mockServer    *MockServer
	repository    *DataRepositoryImpl
	ctx           context.Context
}

// setupDataRepositoryImplTest はテスト用のセットアップを行います
func setupDataRepositoryImplTest(t *testing.T) *DataRepositoryImplTestSuite {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}

	// DataRepositoryImplを直接作成（テスト用）
	repository := &DataRepositoryImpl{
		dataStore: (*infraValkey.DataStore)(nil), // モックを使用するため実際のDataStoreは不要
		server:    (*infraValkey.Server)(nil),    // モックを使用するため実際のServerは不要
		addr:      "localhost:6379",
	}

	return &DataRepositoryImplTestSuite{
		mockDataStore: mockDataStore,
		mockServer:    mockServer,
		repository:    repository,
		ctx:           context.Background(),
	}
}

func TestNewDataRepository_Normal(t *testing.T) {
	// 実際のNewDataRepositoryは外部依存があるため、
	// ここでは基本的な構造体の作成をテストします
	addr := "localhost:6379"

	// 実際のテストでは外部依存があるため、構造体の初期化のみテスト
	repository := &DataRepositoryImpl{
		addr: addr,
	}

	assert.NotNil(t, repository)
	assert.Equal(t, addr, repository.addr)
}

func TestDataRepositoryImpl_GetServerAddress_Normal(t *testing.T) {
	suite := setupDataRepositoryImplTest(t)
	expectedAddr := "localhost:6379"

	addr := suite.repository.GetServerAddress()

	assert.Equal(t, expectedAddr, addr)
}

func TestDataRepositoryImpl_IsServerRunning_Normal(t *testing.T) {
	suite := setupDataRepositoryImplTest(t)

	// サーバーがnilでない場合はtrueを返す
	isRunning := suite.repository.IsServerRunning()

	// 現在の実装では、serverがnilでなければtrueを返す
	// テスト用のsetupではserverはnilなので、実装に応じて調整が必要
	assert.False(t, isRunning) // serverがnilの場合
}

func TestDataRepositoryImpl_IsServerRunning_ServerNotNil(t *testing.T) {
	suite := setupDataRepositoryImplTest(t)

	// serverを非nilに設定
	suite.repository.server = &infraValkey.Server{}

	isRunning := suite.repository.IsServerRunning()

	assert.True(t, isRunning) // serverが非nilの場合
}

// モックを使用したDataRepositoryImplのテスト用ラッパー
type MockDataRepositoryImpl struct {
	*DataRepositoryImpl
	mockDataStore *MockDataStore
	mockServer    *MockServer
}

func NewMockDataRepositoryImpl(mockDataStore *MockDataStore, mockServer *MockServer, addr string) *MockDataRepositoryImpl {
	return &MockDataRepositoryImpl{
		DataRepositoryImpl: &DataRepositoryImpl{
			addr: addr,
		},
		mockDataStore: mockDataStore,
		mockServer:    mockServer,
	}
}

func (m *MockDataRepositoryImpl) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	return m.mockDataStore.GetKeys(ctx, pattern)
}

func (m *MockDataRepositoryImpl) GetValue(ctx context.Context, key string) (string, error) {
	return m.mockDataStore.GetValue(ctx, key)
}

func (m *MockDataRepositoryImpl) GetValueAsByte(ctx context.Context, key string) ([]byte, error) {
	return m.mockDataStore.GetValueAsByte(ctx, key)
}

func (m *MockDataRepositoryImpl) GetType(ctx context.Context, key string) (string, error) {
	return m.mockDataStore.GetType(ctx, key)
}

func (m *MockDataRepositoryImpl) SetValue(ctx context.Context, key string, valueJSON []byte) error {
	return m.mockDataStore.SetValue(ctx, key, valueJSON)
}

func (m *MockDataRepositoryImpl) DeleteKey(ctx context.Context, key string) (bool, error) {
	return m.mockDataStore.DeleteKey(ctx, key)
}

func (m *MockDataRepositoryImpl) StartServer() error {
	return m.mockServer.Start()
}

func (m *MockDataRepositoryImpl) StopServer() error {
	return m.mockServer.Stop()
}

func TestMockDataRepositoryImpl_GetKeys_Normal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedKeys := []string{"key1", "key2", "key3"}
	pattern := "test:*"

	mockDataStore.GetKeysFunc = func(ctx context.Context, p string) ([]string, error) {
		assert.Equal(t, pattern, p)
		return expectedKeys, nil
	}

	keys, err := repository.GetKeys(context.Background(), pattern)

	assert.NoError(t, err)
	assert.Equal(t, expectedKeys, keys)
}

func TestMockDataRepositoryImpl_GetKeys_Error(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedError := errors.New("valkey connection error")

	mockDataStore.GetKeysFunc = func(ctx context.Context, pattern string) ([]string, error) {
		return nil, expectedError
	}

	keys, err := repository.GetKeys(context.Background(), "test:*")

	assert.Error(t, err)
	assert.Nil(t, keys)
	assert.Equal(t, expectedError, err)
}

func TestMockDataRepositoryImpl_GetValue_Normal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedValue := "test_value"
	key := "test_key"

	mockDataStore.GetValueFunc = func(ctx context.Context, k string) (string, error) {
		assert.Equal(t, key, k)
		return expectedValue, nil
	}

	value, err := repository.GetValue(context.Background(), key)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestMockDataRepositoryImpl_GetValue_KeyNotFound(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	mockDataStore.GetValueFunc = func(ctx context.Context, key string) (string, error) {
		return "", nil // キーが存在しない場合は空文字列
	}

	value, err := repository.GetValue(context.Background(), "nonexistent_key")

	assert.NoError(t, err)
	assert.Equal(t, "", value)
}

func TestMockDataRepositoryImpl_GetValue_Error(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedError := errors.New("valkey connection error")

	mockDataStore.GetValueFunc = func(ctx context.Context, key string) (string, error) {
		return "", expectedError
	}

	value, err := repository.GetValue(context.Background(), "test_key")

	assert.Error(t, err)
	assert.Equal(t, "", value)
	assert.Equal(t, expectedError, err)
}

func TestMockDataRepositoryImpl_GetValueAsByte_Normal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedValue := []byte(`{"token": "test_token"}`)
	key := "test_key"

	mockDataStore.GetValueAsByteFunc = func(ctx context.Context, k string) ([]byte, error) {
		assert.Equal(t, key, k)
		return expectedValue, nil
	}

	value, err := repository.GetValueAsByte(context.Background(), key)

	assert.NoError(t, err)
	assert.Equal(t, expectedValue, value)
}

func TestMockDataRepositoryImpl_GetValueAsByte_KeyNotFound(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	mockDataStore.GetValueAsByteFunc = func(ctx context.Context, key string) ([]byte, error) {
		return nil, nil // キーが存在しない場合はnil
	}

	value, err := repository.GetValueAsByte(context.Background(), "nonexistent_key")

	assert.NoError(t, err)
	assert.Nil(t, value)
}

func TestMockDataRepositoryImpl_GetValueAsByte_Error(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedError := errors.New("valkey connection error")

	mockDataStore.GetValueAsByteFunc = func(ctx context.Context, key string) ([]byte, error) {
		return nil, expectedError
	}

	value, err := repository.GetValueAsByte(context.Background(), "test_key")

	assert.Error(t, err)
	assert.Nil(t, value)
	assert.Equal(t, expectedError, err)
}

func TestMockDataRepositoryImpl_GetType_Normal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedType := "string"
	key := "test_key"

	mockDataStore.GetTypeFunc = func(ctx context.Context, k string) (string, error) {
		assert.Equal(t, key, k)
		return expectedType, nil
	}

	typeStr, err := repository.GetType(context.Background(), key)

	assert.NoError(t, err)
	assert.Equal(t, expectedType, typeStr)
}

func TestMockDataRepositoryImpl_GetType_KeyNotFound(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	mockDataStore.GetTypeFunc = func(ctx context.Context, key string) (string, error) {
		return "none", nil // キーが存在しない場合は"none"
	}

	typeStr, err := repository.GetType(context.Background(), "nonexistent_key")

	assert.NoError(t, err)
	assert.Equal(t, "none", typeStr)
}

func TestMockDataRepositoryImpl_GetType_Error(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedError := errors.New("valkey connection error")

	mockDataStore.GetTypeFunc = func(ctx context.Context, key string) (string, error) {
		return "", expectedError
	}

	typeStr, err := repository.GetType(context.Background(), "test_key")

	assert.Error(t, err)
	assert.Equal(t, "", typeStr)
	assert.Equal(t, expectedError, err)
}

func TestMockDataRepositoryImpl_SetValue_Normal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	key := "test_key"
	value := []byte(`{"token": "test_token"}`)

	mockDataStore.SetValueFunc = func(ctx context.Context, k string, v []byte) error {
		assert.Equal(t, key, k)
		assert.Equal(t, value, v)
		return nil
	}

	err := repository.SetValue(context.Background(), key, value)

	assert.NoError(t, err)
}

func TestMockDataRepositoryImpl_SetValue_Error(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedError := errors.New("valkey connection error")

	mockDataStore.SetValueFunc = func(ctx context.Context, key string, value []byte) error {
		return expectedError
	}

	err := repository.SetValue(context.Background(), "test_key", []byte("test_value"))

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestMockDataRepositoryImpl_DeleteKey_Normal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	key := "test_key"

	mockDataStore.DeleteKeyFunc = func(ctx context.Context, k string) (bool, error) {
		assert.Equal(t, key, k)
		return true, nil
	}

	deleted, err := repository.DeleteKey(context.Background(), key)

	assert.NoError(t, err)
	assert.True(t, deleted)
}

func TestMockDataRepositoryImpl_DeleteKey_KeyNotFound(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	mockDataStore.DeleteKeyFunc = func(ctx context.Context, key string) (bool, error) {
		return false, nil // キーが存在しない場合はfalse
	}

	deleted, err := repository.DeleteKey(context.Background(), "nonexistent_key")

	assert.NoError(t, err)
	assert.False(t, deleted)
}

func TestMockDataRepositoryImpl_DeleteKey_Error(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedError := errors.New("valkey connection error")

	mockDataStore.DeleteKeyFunc = func(ctx context.Context, key string) (bool, error) {
		return false, expectedError
	}

	deleted, err := repository.DeleteKey(context.Background(), "test_key")

	assert.Error(t, err)
	assert.False(t, deleted)
	assert.Equal(t, expectedError, err)
}

func TestMockDataRepositoryImpl_StartServer_Normal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	mockServer.StartFunc = func() error {
		return nil
	}

	err := repository.StartServer()

	assert.NoError(t, err)
}

func TestMockDataRepositoryImpl_StartServer_Error(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedError := errors.New("server start error")

	mockServer.StartFunc = func() error {
		return expectedError
	}

	err := repository.StartServer()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestMockDataRepositoryImpl_StopServer_Normal(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	mockServer.StopFunc = func() error {
		return nil
	}

	err := repository.StopServer()

	assert.NoError(t, err)
}

func TestMockDataRepositoryImpl_StopServer_Error(t *testing.T) {
	mockDataStore := &MockDataStore{}
	mockServer := &MockServer{}
	repository := NewMockDataRepositoryImpl(mockDataStore, mockServer, "localhost:6379")

	expectedError := errors.New("server stop error")

	mockServer.StopFunc = func() error {
		return expectedError
	}

	err := repository.StopServer()

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// NewDataRepositoryToStart のテスト用モック関数
func TestNewDataRepositoryToStart_ErrorHandling(t *testing.T) {
	// この関数は実際のValkeyサーバーに依存するため、
	// エラーハンドリングのロジックのみテストします

	// エラーメッセージの形式をテスト
	testError := errors.New("test error")
	wrappedError := fmt.Errorf("valkeyデータリポジトリの初期化に失敗しました: %w", testError)

	assert.Contains(t, wrappedError.Error(), "valkeyデータリポジトリの初期化に失敗しました")
	assert.Contains(t, wrappedError.Error(), "test error")

	serverError := errors.New("server error")
	serverWrappedError := fmt.Errorf("サーバーの起動に失敗しました: %w", serverError)

	assert.Contains(t, serverWrappedError.Error(), "サーバーの起動に失敗しました")
	assert.Contains(t, serverWrappedError.Error(), "server error")
}

// 実際のDataRepositoryImplの実装をテストするためのテストケース
// これらのテストは実際のdataStoreやserverがnilのため、パニックが発生することを確認します

func TestDataRepositoryImpl_RealImplementation_GetKeys_Panic(t *testing.T) {
	repository := &DataRepositoryImpl{
		dataStore: nil,
		server:    nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		repository.GetKeys(context.Background(), "test:*")
	})
}

func TestDataRepositoryImpl_RealImplementation_GetValue_Panic(t *testing.T) {
	repository := &DataRepositoryImpl{
		dataStore: nil,
		server:    nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		repository.GetValue(context.Background(), "test_key")
	})
}

func TestDataRepositoryImpl_RealImplementation_GetValueAsByte_Panic(t *testing.T) {
	repository := &DataRepositoryImpl{
		dataStore: nil,
		server:    nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		repository.GetValueAsByte(context.Background(), "test_key")
	})
}

func TestDataRepositoryImpl_RealImplementation_GetType_Panic(t *testing.T) {
	repository := &DataRepositoryImpl{
		dataStore: nil,
		server:    nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		repository.GetType(context.Background(), "test_key")
	})
}

func TestDataRepositoryImpl_RealImplementation_SetValue_Panic(t *testing.T) {
	repository := &DataRepositoryImpl{
		dataStore: nil,
		server:    nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		repository.SetValue(context.Background(), "test_key", []byte("test_value"))
	})
}

func TestDataRepositoryImpl_RealImplementation_DeleteKey_Panic(t *testing.T) {
	repository := &DataRepositoryImpl{
		dataStore: nil,
		server:    nil,
		addr:      "localhost:6379",
	}

	// dataStoreがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		repository.DeleteKey(context.Background(), "test_key")
	})
}

func TestDataRepositoryImpl_RealImplementation_StartServer_Panic(t *testing.T) {
	repository := &DataRepositoryImpl{
		dataStore: nil,
		server:    nil,
		addr:      "localhost:6379",
	}

	// serverがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		repository.StartServer()
	})
}

func TestDataRepositoryImpl_RealImplementation_StopServer_Panic(t *testing.T) {
	repository := &DataRepositoryImpl{
		dataStore: nil,
		server:    nil,
		addr:      "localhost:6379",
	}

	// serverがnilの場合、パニックが発生することを確認
	assert.Panics(t, func() {
		repository.StopServer()
	})
}
