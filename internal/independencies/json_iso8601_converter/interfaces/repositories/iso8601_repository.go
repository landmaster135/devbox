package repositories

import (
	"fmt"
	"time"

	domainRepo "github.com/landmaster135/devbox/internal/independencies/json_iso8601_converter/domain/repositories"
)

// ISO8601RepositoryImpl はISO8601Repositoryインターフェースの実装です
type ISO8601RepositoryImpl struct{}

// NewISO8601Repository は新しいISO8601RepositoryImplインスタンスを作成します
func NewISO8601Repository() domainRepo.ISO8601Repository {
	return &ISO8601RepositoryImpl{}
}

// ParseISO8601 はISO8601形式の日時文字列をUNIXタイムスタンプに変換します
func (r *ISO8601RepositoryImpl) ParseISO8601(dateStr string) (int64, error) {
	// RFC3339形式（ISO8601の一種）でパース
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// RFC3339Nanoでも試す
		t, err = time.Parse(time.RFC3339Nano, dateStr)
		if err != nil {
			return 0, fmt.Errorf("日時文字列のパースに失敗しました: %w", err)
		}
	}

	// UNIXタイムスタンプに変換
	return t.Unix(), nil
}
