package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/sqlite/config"
	"github.com/landmaster135/devbox/internal/sqlite/infrastructures"
)

// TableRepository はテーブル一覧取得の抽象です。
type TableRepository interface {
	ListTables(ctx context.Context, dbPath string) ([]string, error)
}

// SQLiteService は sqlite CLI の業務処理を提供します。
type SQLiteService struct {
	repository TableRepository
}

// NewSQLiteService は SQLiteService を生成します。
func NewSQLiteService(repository TableRepository) *SQLiteService {
	repo := repository
	if repo == nil {
		repo = infrastructures.NewSQLiteRepository()
	}
	return &SQLiteService{repository: repo}
}

// HandleListTables はテーブル一覧を指定形式で返します。
func (s *SQLiteService) HandleListTables(ctx context.Context, dbPath, format string) (string, error) {
	tables, err := s.repository.ListTables(ctx, dbPath)
	if err != nil {
		return "", err
	}

	switch strings.TrimSpace(format) {
	case config.FormatText:
		return strings.Join(tables, "\n"), nil
	case config.FormatJSON:
		encoded, marshalErr := json.MarshalIndent(tables, "", "  ")
		if marshalErr != nil {
			return "", fmt.Errorf("JSON への整形に失敗しました: %w", marshalErr)
		}
		return string(encoded), nil
	default:
		return "", fmt.Errorf("未対応の format です: %s", format)
	}
}
