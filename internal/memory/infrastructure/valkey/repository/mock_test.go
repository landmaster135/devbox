package valkey

import (
	"context"

	infraValkey "github.com/landmaster135/devbox/internal/valkey/infrastructure/valkey"
)

// MockDataStore は DataStore のモック実装
type MockDataStore struct {
	GetKeysFunc        func(ctx context.Context, pattern string) ([]string, error)
	GetValueFunc       func(ctx context.Context, key string) (string, error)
	GetValueAsByteFunc func(ctx context.Context, key string) ([]byte, error)
	GetTypeFunc        func(ctx context.Context, key string) (string, error)
	SetValueFunc       func(ctx context.Context, key string, valueJSON []byte) error
	DeleteKeyFunc      func(ctx context.Context, key string) (bool, error)
}

// GetKeys はパターンに一致するすべてのキーを取得します
func (m *MockDataStore) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	if m.GetKeysFunc != nil {
		return m.GetKeysFunc(ctx, pattern)
	}
	return []string{}, nil
}

// GetValue はキーに対応する値を取得します
func (m *MockDataStore) GetValue(ctx context.Context, key string) (string, error) {
	if m.GetValueFunc != nil {
		return m.GetValueFunc(ctx, key)
	}
	return "", nil
}

// GetValueAsByte はJSON形式のトークン情報を取得します
func (m *MockDataStore) GetValueAsByte(ctx context.Context, key string) ([]byte, error) {
	if m.GetValueAsByteFunc != nil {
		return m.GetValueAsByteFunc(ctx, key)
	}
	return []byte{}, nil
}

// GetType はキーの型を取得します
func (m *MockDataStore) GetType(ctx context.Context, key string) (string, error) {
	if m.GetTypeFunc != nil {
		return m.GetTypeFunc(ctx, key)
	}
	return "string", nil
}

// SetValue はJSON形式のトークン情報をValkeyに保存します
func (m *MockDataStore) SetValue(ctx context.Context, key string, valueJSON []byte) error {
	if m.SetValueFunc != nil {
		return m.SetValueFunc(ctx, key, valueJSON)
	}
	return nil
}

// DeleteKey はキーを削除します
func (m *MockDataStore) DeleteKey(ctx context.Context, key string) (bool, error) {
	if m.DeleteKeyFunc != nil {
		return m.DeleteKeyFunc(ctx, key)
	}
	return true, nil
}

// MockServer は Server のモック実装
type MockServer struct {
	StartFunc           func() error
	StopFunc            func() error
	GetAddressFunc      func() string
	IsRunningFunc       func() bool
	CheckConnectionFunc func() error
}

// Start はValkeyサーバーを起動します
func (m *MockServer) Start() error {
	if m.StartFunc != nil {
		return m.StartFunc()
	}
	return nil
}

// Stop はValkeyサーバーを停止します
func (m *MockServer) Stop() error {
	if m.StopFunc != nil {
		return m.StopFunc()
	}
	return nil
}

// GetAddress はValkeyサーバーのアドレスを返します
func (m *MockServer) GetAddress() string {
	if m.GetAddressFunc != nil {
		return m.GetAddressFunc()
	}
	return "localhost:6379"
}

// IsRunning はサーバーが起動しているかどうかを返します
func (m *MockServer) IsRunning() bool {
	if m.IsRunningFunc != nil {
		return m.IsRunningFunc()
	}
	return true
}

// CheckConnection はValkeyサーバーへの接続を確認します
func (m *MockServer) CheckConnection() error {
	if m.CheckConnectionFunc != nil {
		return m.CheckConnectionFunc()
	}
	return nil
}

// MockDataStoreFactory は DataStore を作成するためのファクトリー関数のモック
type MockDataStoreFactory struct {
	NewDataStoreFunc func(addr string) (*infraValkey.DataStore, error)
}

// NewDataStore は新しいDataStoreを作成します
func (m *MockDataStoreFactory) NewDataStore(addr string) (*infraValkey.DataStore, error) {
	if m.NewDataStoreFunc != nil {
		return m.NewDataStoreFunc(addr)
	}
	return nil, nil
}

// MockServerFactory は Server を作成するためのファクトリー関数のモック
type MockServerFactory struct {
	NewServerFunc func(addr string) (*infraValkey.Server, error)
}

// NewServer は新しいServerを作成します
func (m *MockServerFactory) NewServer(addr string) (*infraValkey.Server, error) {
	if m.NewServerFunc != nil {
		return m.NewServerFunc(addr)
	}
	return nil, nil
}
