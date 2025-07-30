// Package valkey はValkeyを使用したデータ操作機能を提供します
package valkey

import (
	"context"
	"fmt"

	"github.com/valkey-io/valkey-go"
)

// DataStore はValkeyデータ操作の機能を提供する構造体
type DataStore struct {
	client valkey.Client
}

// NewDataStore は新しいDataStoreを作成します
func NewDataStore(valkeyURL string) (*DataStore, error) {
	opt, err := valkey.ParseURL(valkeyURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse valkey URL: %w", err)
	}

	client, err := valkey.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}

	return &DataStore{
		client: client,
	}, nil
}

// GetKeys はパターンに一致するすべてのキーを取得します
func (s *DataStore) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	// KEYSコマンドを実行
	keys, err := s.client.Do(ctx, s.client.B().Keys().Pattern(pattern).Build()).AsStrSlice()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from valkey: %w", err)
	}

	return keys, nil
}

// GetValue はキーに対応する値を取得します
func (s *DataStore) GetValue(ctx context.Context, key string) (string, error) {
	// GETコマンドを実行
	value, err := s.client.Do(ctx, s.client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "", nil // キーが存在しない場合は空文字列を返す
		}
		return "", fmt.Errorf("failed to get value from valkey: %w", err)
	}

	return value, nil
}

// GetValueAsByte はValkeyからJSON形式のトークン情報を取得します
func (s *DataStore) GetValueAsByte(ctx context.Context, key string) ([]byte, error) {
	tokenStr, err := s.client.Do(ctx, s.client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, nil // トークンが存在しない
		}
		return nil, fmt.Errorf("failed to get token from valkey: %w", err)
	}

	return []byte(tokenStr), nil
}


// GetType はキーの型を取得します
func (s *DataStore) GetType(ctx context.Context, key string) (string, error) {
	// TYPEコマンドを実行
	typeStr, err := s.client.Do(ctx, s.client.B().Type().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "none", nil // キーが存在しない場合は"none"を返す
		}
		return "", fmt.Errorf("failed to get type from valkey: %w", err)
	}

	return typeStr, nil
}

// SetValue はJSON形式のトークン情報をValkeyに保存します
func (s *DataStore) SetValue(ctx context.Context, key string, valueJSON []byte) error {
	err := s.client.Do(ctx, s.client.B().Set().Key(key).Value(string(valueJSON)).Build()).Error()
	if err != nil {
		return fmt.Errorf("failed to save token to valkey: %w", err)
	}

	return nil
}

// DeleteKey はキーを削除します
func (s *DataStore) DeleteKey(ctx context.Context, key string) (bool, error) {
	// DELコマンドを実行し、削除されたキーの数を返す
	result, err := s.client.Do(ctx, s.client.B().Del().Key(key).Build()).AsInt64()
	if err != nil {
		return false, fmt.Errorf("failed to delete key from valkey: %w", err)
	}

	// 削除されたキーの数が0より大きければtrue、そうでなければfalseを返す
	return result > 0, nil
}
