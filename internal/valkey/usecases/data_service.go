// Package valkey はValkeyデータ操作のためのサービスを提供します
package valkey

import (
	"context"
	"fmt"

	config "github.com/landmaster135/devbox/internal/valkey/config"
	loggerRepo "github.com/landmaster135/devbox/internal/valkey/infrastructure/logger/repository"
	valkeyRepo "github.com/landmaster135/devbox/internal/valkey/infrastructure/valkey/repository"
)

// DataService はValkeyデータ操作のためのサービス
type DataService struct {
	repo   valkeyRepo.DataRepository
	logger loggerRepo.Logger
}

// GetRepository はリポジトリを返します
func (s *DataService) GetRepository() valkeyRepo.DataRepository {
	return s.repo
}

// NewDataService は新しいDataServiceを作成します
func NewDataService(repo valkeyRepo.DataRepository, logger loggerRepo.Logger) *DataService {
	return &DataService{
		repo:   repo,
		logger: logger,
	}
}

// NewDataServiceWithConfig は新しいDataServiceを作成します
func NewDataServiceWithConfig(cfg *config.Config) (*DataService, error) {
	// Valkey接続URLを構築
	valkeyURL := cfg.BuildValkeyURL()

	// リポジトリを初期化
	repo, err := valkeyRepo.NewDataRepository(valkeyURL)
	if err != nil {
		return nil, fmt.Errorf("リポジトリの初期化に失敗しました: %v", err)
	}

	// ロガーを初期化
	logger := loggerRepo.NewDefaultLogger()

	return &DataService{
		repo:   repo,
		logger: logger,
	}, nil
}

// GetKeys はパターンに一致するすべてのキーを取得します
func (s *DataService) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	return s.GetRepository().GetKeys(ctx, pattern)
}

// GetValue はキーに対応する値を取得します
func (s *DataService) GetValue(ctx context.Context, key string) (string, error) {
	return s.GetRepository().GetValue(ctx, key)
}

// GetType はキーの型を取得します
func (s *DataService) GetType(ctx context.Context, key string) (string, error) {
	return s.GetRepository().GetType(ctx, key)
}

// SetValue はJSON形式のトークン情報をValkeyに保存します
func (s *DataService) SetValue(ctx context.Context, key string, valueJSON []byte) error {
	return s.GetRepository().SetValue(ctx, key, valueJSON)
}

// DeleteKey はキーを削除します
func (s *DataService) DeleteKey(ctx context.Context, key string) (bool, error) {
	return s.GetRepository().DeleteKey(ctx, key)
}

// DeleteKeys は複数のキーを削除します
func (s *DataService) DeleteKeys(ctx context.Context, keys []string) (map[string]bool, error) {
	results := make(map[string]bool)

	for _, key := range keys {
		deleted, err := s.DeleteKey(ctx, key)
		if err != nil {
			return results, err
		}
		results[key] = deleted
	}

	return results, nil
}

// SelectKeys はValkeyデータを選択して結果を返します
func (s *DataService) SelectKeys(ctx context.Context, key string, keys []string, pattern string, all bool) (any, error) {
	// 引数のバリデーション
	keyProvided := key != ""
	keysProvided := len(keys) > 0
	patternProvided := pattern != ""

	// keyとkeysとpatternは全てゼロ値ではあってはならない
	if !keyProvided && !keysProvided && !patternProvided && !all {
		return nil, fmt.Errorf("key、keys、pattern のいずれか、または all=true を指定してください")
	}

	// keyとkeysとpatternのいずれかがゼロ値でない場合、それ以外の引数はゼロ値でなければならない
	providedCount := 0
	if keyProvided {
		providedCount++
	}
	if keysProvided {
		providedCount++
	}
	if patternProvided {
		providedCount++
	}
	if all {
		providedCount++
	}

	if providedCount > 1 {
		return nil, fmt.Errorf("key、keys、pattern、all のうち、同時に指定できるのは1つだけです")
	}

	// allがtrueの時、keyとkeysとpatternは全てゼロ値でなければならない
	if all && (keyProvided || keysProvided || patternProvided) {
		return nil, fmt.Errorf("all=true の場合、key、keys、pattern は指定できません")
	}

	// キーが指定されている場合
	if key != "" {
		value, err := s.GetValue(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("キー '%s' の値の取得に失敗しました: %w", key, err)
		}

		keyType, err := s.GetType(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("キー '%s' の型の取得に失敗しました: %w", key, err)
		}

		// 結果を返す
		return map[string]any{
			"key":   key,
			"value": value,
			"type":  keyType,
		}, nil
	}

	// keysが指定されている場合
	if len(keys) > 0 {
		// s.GetKeys関数でキーを取得
		var retrievedKeys []string
		for _, keyPattern := range keys {
			matchedKeys, err := s.GetKeys(ctx, keyPattern)
			if err != nil {
				return nil, fmt.Errorf("パターン '%s' に一致するキーの取得に失敗しました: %w", keyPattern, err)
			}
			retrievedKeys = append(retrievedKeys, matchedKeys...)
		}

		// 結果を返す
		return map[string]any{
			"keys":  retrievedKeys,
			"count": len(retrievedKeys),
		}, nil
	}

	// パターンが指定されている場合またはすべてのキーを表示する場合
	searchPattern := pattern
	if all {
		searchPattern = "*"
	}

	matchedKeys, err := s.GetKeys(ctx, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("パターン '%s' に一致するキーの取得に失敗しました: %w", searchPattern, err)
	}

	// 結果を返す
	return map[string]any{
		"keys":  matchedKeys,
		"count": len(matchedKeys),
	}, nil
}

// GetAllValues はすべてのキーの値を取得して結果を返します
func (s *DataService) GetAllValues(ctx context.Context, keys []string) (any, error) {
	if len(keys) == 0 {
		return map[string]any{
			"values": []any{},
			"count":  0,
		}, nil
	}

	type KeyValue struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		Type  string `json:"type"`
	}

	var keyValues []KeyValue

	for _, key := range keys {
		value, err := s.GetValue(ctx, key)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("キー '%s' の値の取得に失敗しました", key))
			continue
		}

		keyType, err := s.GetType(ctx, key)
		if err != nil {
			s.logger.Warn(fmt.Sprintf("キー '%s' の型の取得に失敗しました", key))
			continue
		}

		keyValues = append(keyValues, KeyValue{
			Key:   key,
			Value: value,
			Type:  keyType,
		})
	}

	// KeyValueを[]anyに変換
	values := make([]any, len(keyValues))
	for i, kv := range keyValues {
		values[i] = map[string]any{
			"key":   kv.Key,
			"value": kv.Value,
			"type":  kv.Type,
		}
	}

	return map[string]any{
		"values": values,
		"count":  len(values),
	}, nil
}

// DeleteData はValkeyデータを削除します
func (s *DataService) DeleteData(ctx context.Context, key string, keys []string, pattern string, dryRun bool) (any, error) {
	// ドライランモードの場合
	if dryRun {
		s.logger.Info("ドライランモードで実行中 - 実際のデータ削除は行いません")

		// 単一のキーが指定されている場合
		if key != "" {
			return map[string]any{
				"key":     key,
				"deleted": false,
				"message": "ドライランモードで実行しました。実際のデータ削除は行われていません。",
			}, nil
		}

		// 複数のキーが直接指定されている場合
		if len(keys) > 0 {
			return map[string]any{
				"keys":    keys,
				"count":   len(keys),
				"deleted": 0,
				"message": "ドライランモードで実行しました。実際のデータ削除は行われていません。",
			}, nil
		}

		// パターンが指定されている場合
		if pattern != "" {
			// パターンに一致するキーを取得
			matchedKeys, err := s.GetKeys(ctx, pattern)
			if err != nil {
				return nil, fmt.Errorf("パターン '%s' に一致するキーの取得に失敗しました: %w", pattern, err)
			}

			// キーが見つからない場合
			if len(matchedKeys) == 0 {
				return map[string]any{
					"message": fmt.Sprintf("パターン '%s' に一致するキーが見つかりませんでした", pattern),
					"count":   0,
				}, nil
			}

			// 結果を返す
			return map[string]any{
				"keys":    matchedKeys,
				"count":   len(matchedKeys),
				"deleted": 0,
				"message": "ドライランモードで実行しました。実際のデータ削除は行われていません。",
			}, nil
		}

		return nil, fmt.Errorf("キーまたはパターンを指定してください")
	}

	// 単一のキーが指定されている場合
	if key != "" {
		deleted, err := s.DeleteKey(ctx, key)
		if err != nil {
			return nil, fmt.Errorf("キー '%s' の削除に失敗しました: %w", key, err)
		}

		// 結果を返す
		return map[string]any{
			"key":     key,
			"deleted": deleted,
		}, nil
	}

	// 複数のキーが直接指定されている場合
	if len(keys) > 0 {
		// 複数のキーを削除
		results, err := s.DeleteKeys(ctx, keys)
		if err != nil {
			return nil, fmt.Errorf("キーの削除に失敗しました: %w", err)
		}

		// 削除されたキーの数をカウント
		deletedCount := 0
		for _, deleted := range results {
			if deleted {
				deletedCount++
			}
		}

		// 結果を返す
		return map[string]any{
			"keys":    keys,
			"results": results,
			"count":   len(keys),
			"deleted": deletedCount,
		}, nil
	}

	// パターンが指定されている場合
	if pattern != "" {
		// パターンに一致するキーを取得
		matchedKeys, err := s.GetKeys(ctx, pattern)
		if err != nil {
			return nil, fmt.Errorf("パターン '%s' に一致するキーの取得に失敗しました: %w", pattern, err)
		}

		// キーが見つからない場合
		if len(matchedKeys) == 0 {
			return map[string]any{
				"message": fmt.Sprintf("パターン '%s' に一致するキーが見つかりませんでした", pattern),
				"count":   0,
			}, nil
		}

		// 複数のキーを削除
		results, err := s.DeleteKeys(ctx, matchedKeys)
		if err != nil {
			return nil, fmt.Errorf("キーの削除に失敗しました: %w", err)
		}

		// 削除されたキーの数をカウント
		deletedCount := 0
		for _, deleted := range results {
			if deleted {
				deletedCount++
			}
		}

		// 結果を返す
		return map[string]any{
			"keys":    matchedKeys,
			"results": results,
			"count":   len(matchedKeys),
			"deleted": deletedCount,
		}, nil
	}

	return nil, fmt.Errorf("キーまたはパターンを指定してください")
}
