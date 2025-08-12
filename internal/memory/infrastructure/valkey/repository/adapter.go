package valkey

import (
	"context"
	"fmt"

	infraValkey "github.com/landmaster135/devbox/internal/valkey/infrastructure/valkey"
)

// DataStoreAdapter は *infraValkey.DataStore を DataRepository に適応させるアダプター
type DataStoreAdapter struct {
	dataStore *infraValkey.DataStore
	addr      string
}

// NewDataStoreAdapter は新しい DataStoreAdapter を作成します
// dataStoreがnilの場合はエラーを返します
func NewDataStoreAdapter(dataStore interface{}, addr string) (DataRepository, error) {
	if dataStore == nil {
		return nil, fmt.Errorf("dataStore is nil")
	}
	if addr == "" {
		return nil, fmt.Errorf("addr is empty")
	}

	// DataStoreの型をチェック
	ds, ok := dataStore.(*infraValkey.DataStore)
	if !ok {
		return nil, fmt.Errorf("dataStore is not *infraValkey.DataStore")
	}

	return &DataStoreAdapter{
		dataStore: ds,
		addr:      addr,
	}, nil
}

// GetKeys はパターンに一致するすべてのキーを取得します
func (a *DataStoreAdapter) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	return a.dataStore.GetKeys(ctx, pattern)
}

// GetValue はキーに対応する値を取得します
func (a *DataStoreAdapter) GetValue(ctx context.Context, key string) (string, error) {
	return a.dataStore.GetValue(ctx, key)
}

// GetValueAsByte はJSON形式のトークン情報を取得します
func (a *DataStoreAdapter) GetValueAsByte(ctx context.Context, key string) ([]byte, error) {
	return a.dataStore.GetValueAsByte(ctx, key)
}

// GetType はキーの型を取得します
func (a *DataStoreAdapter) GetType(ctx context.Context, key string) (string, error) {
	return a.dataStore.GetType(ctx, key)
}

// SetValue はJSON形式のトークン情報をValkeyに保存します
func (a *DataStoreAdapter) SetValue(ctx context.Context, key string, valueJSON []byte) error {
	return a.dataStore.SetValue(ctx, key, valueJSON)
}

// StartServer はDBサーバーを起動します（アダプターでは何もしない）
func (a *DataStoreAdapter) StartServer() error {
	// このアダプターではサーバー操作は行わない
	return nil
}

// StopServer はDBサーバーを停止します（アダプターでは何もしない）
func (a *DataStoreAdapter) StopServer() error {
	// このアダプターではサーバー操作は行わない
	return nil
}

// IsServerRunning はDBサーバーが起動しているかどうかを返します
func (a *DataStoreAdapter) IsServerRunning() bool {
	// このアダプターではサーバーは常に起動していると仮定
	return true
}

// GetServerAddress はDBサーバーのアドレスを返します
func (a *DataStoreAdapter) GetServerAddress() string {
	return a.addr
}

// DeleteKey はキーを削除します
func (a *DataStoreAdapter) DeleteKey(ctx context.Context, key string) (bool, error) {
	return a.dataStore.DeleteKey(ctx, key)
}
