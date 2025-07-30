package valkey

import (
	"context"
)

// DataRepository はValkeyデータ操作のためのリポジトリインターフェース
type DataRepository interface {
	// GetKeys はパターンに一致するすべてのキーを取得します
	GetKeys(ctx context.Context, pattern string) ([]string, error)

	// GetValue はキーに対応する値を取得します
	// キーが存在しない場合は空文字列とnilを返します
	GetValue(ctx context.Context, key string) (string, error)

	// GetTokenAsByte はJSON形式のトークン情報を取得します
	// キーが存在しない場合は空文字列とnilを返します
	GetValueAsByte(ctx context.Context, key string) ([]byte, error)

	// GetType はキーの型を取得します
	// キーが存在しない場合は"none"とnilを返します
	GetType(ctx context.Context, key string) (string, error)

	// SetValue はJSON形式のトークン情報をValkeyに保存します
	SetValue(ctx context.Context, key string, valueJSON []byte) error

	// StartServer はDBサーバーを起動します
	StartServer() error

	// StopServer はDBサーバーを停止します
	StopServer() error

	// IsServerRunning はDBサーバーが起動しているかどうかを返します
	IsServerRunning() bool

	// GetServerAddress はDBサーバーのアドレスを返します
	GetServerAddress() string

	// DeleteKey はキーを削除します
	// 削除に成功した場合はtrue、キーが存在しない場合はfalseを返します
	DeleteKey(ctx context.Context, key string) (bool, error)
}
