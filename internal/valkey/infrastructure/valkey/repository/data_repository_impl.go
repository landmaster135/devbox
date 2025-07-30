// Package valkey はValkeyデータ操作のためのリポジトリ実装を提供します
package valkey

import (
	"context"
	"fmt"
	"strings"

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
func NewDataRepository(valkeyURL string) (DataRepository, error) {
	// URLからアドレス部分を抽出（サーバー用）
	// 簡易的にhost:portを抽出
	addr := extractAddrFromURL(valkeyURL)

	server, err := infraValkey.NewServer(addr)
	if err != nil {
		return nil, err
	}

	dataStore, err := infraValkey.NewDataStore(valkeyURL)
	if err != nil {
		return nil, err
	}

	return &DataRepositoryImpl{
		dataStore: dataStore,
		server:    server,
		addr:      addr,
	}, nil
}

func NewDataRepositoryToStart(valkeyURL string) (DataRepository, error) {
	repo, err := NewDataRepository(valkeyURL)
	if err != nil {
		return nil, fmt.Errorf("valkeyデータリポジトリの初期化に失敗しました: %w", err)
	}

	// サーバーの起動
	if err := repo.StartServer(); err != nil {
		return nil, fmt.Errorf("サーバーの起動に失敗しました: %w", err)
	}

	return repo, nil
}

// extractAddrFromURL はValkeyURLからhost:port部分を抽出します
func extractAddrFromURL(valkeyURL string) string {
	// 簡易的な実装：valkey://[user:pass@]host:port[/db] からhost:portを抽出
	// より厳密な実装が必要な場合はurl.Parseを使用
	if len(valkeyURL) > 9 && valkeyURL[:9] == "valkey://" {
		remaining := valkeyURL[9:]

		// @がある場合は認証情報をスキップ
		if atIndex := strings.Index(remaining, "@"); atIndex != -1 {
			remaining = remaining[atIndex+1:]
		}

		// /がある場合はデータベース番号をスキップ
		if slashIndex := strings.Index(remaining, "/"); slashIndex != -1 {
			remaining = remaining[:slashIndex]
		}

		return remaining
	}

	// フォールバック
	return "localhost:6379"
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
