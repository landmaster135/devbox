// Package valkey はValkeyデータ操作のためのリポジトリ実装を提供します
package valkey

import (
	"context"
	"fmt"

	infraValkey "github.com/landmaster135/devbox/internal/valkey/infrastructure/valkey"
)

// DataRepositoryImpl はValkeyデータリポジトリの実装
type DataRepositoryImpl struct {
	dataStore *infraValkey.DataStore
	server    *infraValkey.Server
	addr      string
}

// インターフェースを実装していることを確認
var _ DataRepository = (*DataRepositoryImpl)(nil)

// NewDataRepository は新しいDataRepositoryImplを作成します
func NewDataRepository(addr string) (DataRepository, error) {
	server, err := infraValkey.NewServer(addr)
	if err != nil {
		return nil, err
	}

	dataStore, err := infraValkey.NewDataStore(addr)
	if err != nil {
		return nil, err
	}

	return &DataRepositoryImpl{
		dataStore: dataStore,
		server:    server,
		addr:      addr,
	}, nil
}

func NewDataRepositoryToStart(addr string) (DataRepository, error) {
	repo, err := NewDataRepository(addr)
	if err != nil {
		return nil, fmt.Errorf("valkeyデータリポジトリの初期化に失敗しました: %w", err)
	}

	// サーバーの起動
	if err := repo.StartServer(); err != nil {
		return nil, fmt.Errorf("サーバーの起動に失敗しました: %w", err)
	}

	return repo, nil
}

// GetKeys はパターンに一致するすべてのキーを取得します
func (r *DataRepositoryImpl) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	return r.dataStore.GetKeys(ctx, pattern)
}

// GetValue はキーに対応する値を取得します
func (r *DataRepositoryImpl) GetValue(ctx context.Context, key string) (string, error) {
	return r.dataStore.GetValue(ctx, key)
}

// GetValueAsByte はJSON形式のトークン情報を取得します
func (r *DataRepositoryImpl) GetValueAsByte(ctx context.Context, key string) ([]byte, error) {
	return r.dataStore.GetValueAsByte(ctx, key)
}

// GetType はキーの型を取得します
func (r *DataRepositoryImpl) GetType(ctx context.Context, key string) (string, error) {
	return r.dataStore.GetType(ctx, key)
}

// SetValue はJSON形式のトークン情報をValkeyに保存します
func (r *DataRepositoryImpl) SetValue(ctx context.Context, key string, valueJSON []byte) error {
	return r.dataStore.SetValue(ctx, key, valueJSON)
}

// StartServer はDBサーバーを起動します
func (r *DataRepositoryImpl) StartServer() error {
	return r.server.Start()
}

// StopServer はDBサーバーを停止します
func (r *DataRepositoryImpl) StopServer() error {
	return r.server.Stop()
}

// IsServerRunning はDBサーバーが起動しているかどうかを返します
func (r *DataRepositoryImpl) IsServerRunning() bool {
	// サーバーの状態を確認するロジック
	// 現在のServer実装ではisRunningフィールドが公開されていないため、
	// サーバーが非nilであれば起動していると仮定
	return r.server != nil
}

// GetServerAddress はDBサーバーのアドレスを返します
func (r *DataRepositoryImpl) GetServerAddress() string {
	return r.addr
}

// DeleteKey はキーを削除します
func (r *DataRepositoryImpl) DeleteKey(ctx context.Context, key string) (bool, error) {
	return r.dataStore.DeleteKey(ctx, key)
}
